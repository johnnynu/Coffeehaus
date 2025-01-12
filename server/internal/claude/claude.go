package claude

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/liushuangls/go-anthropic/v2"
)

type ClaudeConfig struct {
	APIKey string
	MaxTokens int
	Temperature float64
}

type Service struct {
	client *anthropic.Client
}

func NewClaudeConfig(apiKey string) *ClaudeConfig {
	return &ClaudeConfig{
		APIKey: os.Getenv("CLAUDE_API_KEY"),
		MaxTokens: 1000,
		Temperature: 0.5,
	}
}

func NewService(apiKey string) *Service {
	return &Service{
		client: anthropic.NewClient(apiKey),
	}
}

func (s* Service) AnalyzeSearchQuery(ctx context.Context, query string, userLocation string) (*SearchIntent, error) {
	systemPrompt := `You are a search query analyzer for a coffee shop discovery app. 
	You must correctly identify coffee shop names, even if they are unique or use technical terms.
	Your task is to analyze search queries and return structured JSON that matches the SearchIntent type.
	Only respond with valid JSON, no other text.`
	
	userPrompt := fmt.Sprintf(`Analyze this coffee shop search query and return a JSON object. The user's current location is: %s
	
	Query: "%s"
	
	Rules for classification:
	1. IMPORTANT: Treat multi-word phrases followed by "coffee", "cafe", "roasters", or appearing before "in", "near", "at" as potential shop names
	2. If a potential shop name is detected, ALWAYS classify as "specific" and store the full name in terms.shop
	3. If the query uses "near", "around", or current location context WITHOUT a shop name, classify as "proximity"
	4. If the query mentions a location/area WITHOUT a shop name, classify as "area"
	5. For proximity searches, include a reasonable radius in km
	6. Extract relevant filters (e.g., "matcha", "pour-over") but do NOT classify shop names as filters
	7. Only remove standalone filler words
	8. Preserve all product/drink names intact as single filters (e.g., "matcha latte" is one filter)
	
	Location Name Standardization:
	- Always use full city names (e.g., "Los Angeles" not "LA")
	- For well-known areas within cities, use format "Area, City" (e.g., "Little Tokyo, Los Angeles")
	- Use standard US state abbreviations (CA, NY, etc.)
	- For current location context, use the provided user location as is
	- Common standardizations:
	  * "LA" -> "Los Angeles"
	  * "OC" -> "Orange County"
	  * "NYC" -> "New York City"
	  * "SF" -> "San Francisco"
	  * "DTLA" -> "Downtown Los Angeles"
	
	Examples of specific shop queries:
	- "File Systems of Coffee in LA" -> specific (shop: "File Systems of Coffee", location: "Los Angeles")
	- "Stereoscope Coffee near me" -> specific (shop: "Stereoscope Coffee")
	- "Coffee at Blue Bottle DTLA" -> specific (shop: "Blue Bottle", location: "Downtown Los Angeles")
	
	Example format:
	{
		"searchType": "proximity|area|specific",
		"normalizedQuery": "normalized search terms",
		"location": {
			"name": "location name",
			"radius": radiusInKm
		},
		"terms": {
			"shop": "shop name",
			"filters": ["filter1", "filter2"]
		}
	}`, userLocation, query)

	// create message request
	resp, err := s.client.CreateMessages(ctx, anthropic.MessagesRequest{
		Model: anthropic.ModelClaude3Dot5Sonnet20241022,
		Messages: []anthropic.Message{
			anthropic.NewUserTextMessage(userPrompt),
		},
		MaxTokens: 1000,
		System: systemPrompt,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to analyze search query: %w", err)
	}

	// extract the response text
	responseText := resp.Content[0].Text

	// parse the JSON response
	var searchIntent SearchIntent
	if err := json.Unmarshal([]byte(*responseText), &searchIntent); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %w", err)
	}

	return &searchIntent, nil
}