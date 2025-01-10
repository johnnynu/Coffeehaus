package maps

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"googlemaps.github.io/maps"
)

type MapsClient struct {
	client *maps.Client
}

func NewMapsClient() (*MapsClient, error) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("GOOGLE_MAPS_API_KEY is not set")
	}

	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create maps client: %w", err)
	}

	return &MapsClient{client: client}, nil
}

func (m *MapsClient) TestConnection(ctx context.Context) error {
	request := &maps.GeocodingRequest{Address: "Los Angeles, CA"}

	_, err := m.client.Geocode(ctx, request)
	if err != nil {
		return fmt.Errorf("failed to test maps connection: %w", err)
	}

	return nil
}

// CoffeeShopDetails combines data from Places Search and Place Details
type CoffeeShopDetails struct {
	PlaceID          string
	Name             string
	FormattedAddress string
	Vicinity         string
	Location         maps.LatLng
	Rating           float32
	UserRatingsTotal int
	PriceLevel       int
	Types            []string
	Photos           []maps.Photo
	OpeningHours     *maps.OpeningHours
	Website          string
	FormattedPhone   string
	BusinessStatus   string
}

// search for coffee shops near the given coordinates
func (m *MapsClient) SearchCoffeeShops(ctx context.Context, lat, lng float64, radiusMeters uint) ([]*CoffeeShopDetails, error) {
	location := &maps.LatLng{
		Lat: lat,
		Lng: lng,
	}

	request := &maps.NearbySearchRequest{
		Location: location,
		Radius:   radiusMeters,
		Type:     "cafe",
		Keyword:  "coffee shop",
	}

	response, err := m.client.NearbySearch(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to search for coffee shops: %w", err)
	}

	// Limit results to 10
	maxResults := 10
	if len(response.Results) > maxResults {
		response.Results = response.Results[:maxResults]
	}

	var results []*CoffeeShopDetails
	for _, place := range response.Results {
		// Get additional details for each place
		detailsRequest := &maps.PlaceDetailsRequest{
			PlaceID: place.PlaceID,
		}

		details, err := m.client.PlaceDetails(ctx, detailsRequest)
		if err != nil {
			log.Printf("Warning: failed to get details for place %s: %v", place.Name, err)
			continue
		}

		result := &CoffeeShopDetails{
			PlaceID:          details.PlaceID,
			Name:             details.Name,
			FormattedAddress: details.FormattedAddress,
			Vicinity:         details.Vicinity,
			Location:         details.Geometry.Location,
			Rating:           details.Rating,
			UserRatingsTotal: details.UserRatingsTotal,
			PriceLevel:       details.PriceLevel,
			Types:            details.Types,
			Photos:           details.Photos,
			OpeningHours:     details.OpeningHours,
			Website:          details.Website,
			FormattedPhone:   details.InternationalPhoneNumber,
			BusinessStatus:   details.BusinessStatus,
		}

		results = append(results, result)
	}

	return results, nil
}

// SearchSpecificCoffeeShop searches for a specific coffee shop by name and returns all matching locations
func (m *MapsClient) SearchSpecificCoffeeShop(ctx context.Context, shopName string, location string) ([]*CoffeeShopDetails, error) {
	// Try autocomplete first to handle typos better
	autoCompleteRequest := &maps.PlaceAutocompleteRequest{
		Input: fmt.Sprintf("%s %s", shopName, location),
		Types: maps.AutocompletePlaceTypeEstablishment,
	}

	predictions, err := m.client.PlaceAutocomplete(ctx, autoCompleteRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get autocomplete predictions: %w", err)
	}

	// If we got predictions, use them to find the correct shop name
	var correctedShopName string
	if len(predictions.Predictions) > 0 {
		// Use the first prediction's name as it's likely the most relevant match
		correctedShopName = predictions.Predictions[0].StructuredFormatting.MainText
	} else {
		correctedShopName = shopName
	}

	// Now use text search with the corrected name to find all locations
	query := fmt.Sprintf("%s in %s", correctedShopName, location)
	request := &maps.TextSearchRequest{
		Query: query,
	}

	response, err := m.client.TextSearch(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to search for coffee shop: %w", err)
	}

	if len(response.Results) == 0 {
		return nil, fmt.Errorf("no coffee shop found with name: %s in %s", shopName, location)
	}

	// Take only first 10 results
	maxResults := 10
	if len(response.Results) > maxResults {
		response.Results = response.Results[:maxResults]
	}

	// Filter results to only include places that have the shop name in their name
	var results []*CoffeeShopDetails
	shopNameLower := strings.ToLower(correctedShopName)

	for _, place := range response.Results {
		// Check if this place's name contains the shop name we're looking for
		if !strings.Contains(strings.ToLower(place.Name), shopNameLower) {
			continue
		}

		// Get more details using Place Details
		detailsRequest := &maps.PlaceDetailsRequest{
			PlaceID: place.PlaceID,
		}

		details, err := m.client.PlaceDetails(ctx, detailsRequest)
		if err != nil {
			log.Printf("Warning: failed to get details for place %s: %v", place.Name, err)
			continue
		}

		result := &CoffeeShopDetails{
			PlaceID:          details.PlaceID,
			Name:             details.Name,
			FormattedAddress: details.FormattedAddress,
			Vicinity:         details.Vicinity,
			Location:         details.Geometry.Location,
			Rating:           details.Rating,
			UserRatingsTotal: details.UserRatingsTotal,
			PriceLevel:       details.PriceLevel,
			Types:            details.Types,
			Photos:           details.Photos,
			OpeningHours:     details.OpeningHours,
			Website:          details.Website,
			FormattedPhone:   details.InternationalPhoneNumber,
			BusinessStatus:   details.BusinessStatus,
		}

		results = append(results, result)
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no locations found for: %s in %s", shopName, location)
	}

	return results, nil
}

// SearchCoffeeShopsByArea searches for coffee shops or related places in a specific area using Text Search
func (m *MapsClient) SearchCoffeeShopsByArea(ctx context.Context, query string) ([]*CoffeeShopDetails, error) {
	// Use the query directly instead of formatting it
	request := &maps.TextSearchRequest{
		Query: query,
		Type:  "cafe",
	}

	response, err := m.client.TextSearch(ctx, request)
	if err != nil {
		return nil, fmt.Errorf("failed to search for places: %w", err)
	}

	if len(response.Results) == 0 {
		return nil, fmt.Errorf("no places found for query: %s", query)
	}

	// Limit results to 10
	maxResults := 10
	if len(response.Results) > maxResults {
		response.Results = response.Results[:maxResults]
	}

	// Get additional details for each result
	var results []*CoffeeShopDetails
	for _, place := range response.Results {
		// Get more details using Place Details
		detailsRequest := &maps.PlaceDetailsRequest{
			PlaceID: place.PlaceID,
		}

		details, err := m.client.PlaceDetails(ctx, detailsRequest)
		if err != nil {
			log.Printf("Warning: failed to get details for place %s: %v", place.Name, err)
			continue
		}

		// Create result with all available details
		result := &CoffeeShopDetails{
			PlaceID:          details.PlaceID,
			Name:             details.Name,
			FormattedAddress: details.FormattedAddress,
			Vicinity:         details.Vicinity,
			Location:         details.Geometry.Location,
			Rating:           details.Rating,
			UserRatingsTotal: details.UserRatingsTotal,
			PriceLevel:       details.PriceLevel,
			Types:            details.Types,
			Photos:           details.Photos,
			OpeningHours:     details.OpeningHours,
			Website:          details.Website,
			FormattedPhone:   details.InternationalPhoneNumber,
			BusinessStatus:   details.BusinessStatus,
		}

		results = append(results, result)
	}

	return results, nil
}