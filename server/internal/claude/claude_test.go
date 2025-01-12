package claude

import (
	"context"
	"encoding/json"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func init() {
	// Load the .env file before running tests
	if err := godotenv.Load("../../.env"); err != nil {
		println("Warning: Error loading .env file:", err)
	}
}

func TestAnalyzeSearchQuery(t *testing.T) {
	apiKey := os.Getenv("CLAUDE_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test: CLAUDE_API_KEY not set")
	}

	service := NewService(apiKey)
	ctx := context.Background()

	tests := []struct {
		name          string
		query         string
		userLocation  string
		wantType      string
		wantError     bool
	}{
		{
			name:         "specific shop search",
			query:        "File Systems of Coffee",
			userLocation: "Garden Grove, CA",
			wantType:     "specific",
			wantError:    false,
		},
		{
			name:         "proximity search",
			query:        "coffee shops near me",
			userLocation: "Long Beach, CA",
			wantType:     "proximity",
			wantError:    false,
		},
		{
			name:         "area search with filter",
			query:        "matcha lattes in LA",
			userLocation: "Los Angeles, CA",
			wantType:     "area",
			wantError:    false,
		},
		{
			name:         "area search with filter - variant 1",
			query:        "coffee shops that offer matcha in LA",
			userLocation: "Los Angeles, CA",
			wantType:     "area",
			wantError:    false,
		},
		{
			name:         "area search with filter - variant 2",
			query:        "coffee shops with matcha in los Angeles",
			userLocation: "Los Angeles, CA",
			wantType:     "area",
			wantError:    false,
		},
		{
			name:         "area search with filter - variant 3",
			query:        "coffee and matcha in LA",
			userLocation: "Los Angeles, CA",
			wantType:     "area",
			wantError:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.AnalyzeSearchQuery(ctx, tt.query, tt.userLocation)
			
			if (err != nil) != tt.wantError {
				t.Errorf("AnalyzeSearchQuery() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if err == nil {
				// Pretty print the entire result object
				resultJSON, _ := json.MarshalIndent(result, "", "    ")
				t.Logf("Claude Response:\n%s", string(resultJSON))

				if result.SearchType != tt.wantType {
					t.Errorf("AnalyzeSearchQuery() searchType = %v, want %v", result.SearchType, tt.wantType)
				}

				if result.NormalizedQuery == "" {
					t.Error("AnalyzeSearchQuery() normalizedQuery is empty")
				}

				if result.Location == nil {
					t.Error("AnalyzeSearchQuery() location is nil")
				}

				// Additional checks for specific search type
				if tt.wantType == "specific" {
					if result.Terms.Shop == "" {
						t.Error("AnalyzeSearchQuery() specific search should have shop name in terms")
					}
				}
			}
		})
	}
} 