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

	// Use all predictions that match our shop name
	var results []*CoffeeShopDetails
	shopNameLower := strings.ToLower(shopName)
	seenPlaceIDs := make(map[string]bool)

	// First try predictions if we have any
	if len(predictions.Predictions) > 0 {
		for _, prediction := range predictions.Predictions {
			predictionName := strings.ToLower(prediction.StructuredFormatting.MainText)
			if strings.Contains(predictionName, shopNameLower) {
				// Get details for this prediction
				detailsRequest := &maps.PlaceDetailsRequest{
					PlaceID: prediction.PlaceID,
				}

				details, err := m.client.PlaceDetails(ctx, detailsRequest)
				if err != nil {
					log.Printf("Warning: failed to get details for place %s: %v", prediction.Description, err)
					continue
				}

				// Skip if we've already seen this place
				if seenPlaceIDs[details.PlaceID] {
					continue
				}
				seenPlaceIDs[details.PlaceID] = true

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
		}
	}

	// Then do a text search to find any additional locations
	query := fmt.Sprintf("%s in %s", shopName, location)
	request := &maps.TextSearchRequest{
		Query: query,
	}

	response, err := m.client.TextSearch(ctx, request)
	if err != nil {
		if len(results) > 0 {
			// If we already have results from predictions, don't fail
			log.Printf("Warning: text search failed: %v", err)
		} else {
			return nil, fmt.Errorf("failed to search for coffee shop: %w", err)
		}
	}

	// Process text search results
	if len(response.Results) > 0 {
		for _, place := range response.Results {
			// Skip if we've already seen this place
			if seenPlaceIDs[place.PlaceID] {
				continue
			}

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

			seenPlaceIDs[details.PlaceID] = true

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

func (m* MapsClient) ReverseGeocode(ctx context.Context, lat, lng float64) (string, error) {
	location := &maps.LatLng{
		Lat: lat,
		Lng: lng,
	}

	resp, err := m.client.ReverseGeocode(ctx, &maps.GeocodingRequest{
		LatLng: location,
	})

	if err != nil {
		return "", fmt.Errorf("failed to reverse geocode: %w", err)
	}

	if len(resp) == 0 {
		return "", fmt.Errorf("no results found for location: %f,%f", lat, lng)
	}

	// format the location string
	var locality, adminArea, postalCode string
	for _, component := range resp[0].AddressComponents {
		for _, typ := range component.Types {
			switch typ {
			case "locality":
				locality = component.LongName
			case "administrative_area_level_1":
				adminArea = component.ShortName
			case "postal_code":
				postalCode = component.ShortName
			}
		}
	}

	if locality != "" && adminArea != "" && postalCode != "" {
		return fmt.Sprintf("%s, %s %s", locality, adminArea, postalCode), nil
	}

	// fallback to formatted address if components not found
	return resp[0].FormattedAddress, nil
}
