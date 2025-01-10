package maps

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/joho/godotenv"
)

func init() {
	// Load the .env file before running tests
	if err := godotenv.Load("../../.env"); err != nil {
		println("Warning: Error loading .env file:", err)
	}
}

func TestNewMapsClient(t *testing.T) {
	// Store original API key
	originalKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	defer os.Setenv("GOOGLE_MAPS_API_KEY", originalKey) // Restore after tests

	tests := []struct {
		name          string
		setupEnv      func()
		expectError   bool
		errorContains string
	}{
		{
			name: "successful client creation",
			setupEnv: func() {
				os.Setenv("GOOGLE_MAPS_API_KEY", "test-api-key")
			},
			expectError: false,
		},
		{
			name: "missing API key",
			setupEnv: func() {
				os.Unsetenv("GOOGLE_MAPS_API_KEY")
			},
			expectError:   true,
			errorContains: "GOOGLE_MAPS_API_KEY is not set",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test environment
			tt.setupEnv()
			defer os.Setenv("GOOGLE_MAPS_API_KEY", originalKey) // Restore after each test case

			// Run the test
			client, err := NewMapsClient()

			// Check results
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				} else if tt.errorContains != "" && !strings.Contains(err.Error(), tt.errorContains) {
					t.Errorf("expected error containing %q but got %q", tt.errorContains, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				if client == nil {
					t.Error("expected client but got nil")
				}
			}
		})
	}
}

func TestMapsClient_SearchCoffeeShops(t *testing.T) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test because GOOGLE_MAPS_API_KEY is not set")
	}

	client, err := NewMapsClient()
	if err != nil {
		t.Fatalf("Failed to create maps client: %v", err)
	}

	ctx := context.Background()
	// Test coordinates (Westminster, CA)
	lat := 33.7514
	lng := -117.9940
	radius := uint(3000) // 3km radius

	results, err := client.SearchCoffeeShops(ctx, lat, lng, radius)
	if err != nil {
		t.Errorf("SearchCoffeeShops failed: %v", err)
		return
	}

	if len(results) == 0 {
		t.Error("expected to find coffee shops but got none")
		return
	}

	t.Logf("\nFound %d coffee shops near Westminster", len(results))
				
	// Log details for first few results
	maxToShow := 3
	if len(results) < maxToShow {
		maxToShow = len(results)
	}
	
	for i := 0; i < maxToShow; i++ {
		place := results[i]
		t.Logf("\n=== Coffee Shop %d ===", i+1)
		t.Logf("Name: %s", place.Name)
		t.Logf("Place ID: %s", place.PlaceID)
		t.Logf("Address: %s", place.FormattedAddress)
		t.Logf("Vicinity: %s", place.Vicinity)
		t.Logf("Location: %v", place.Location)
		
		if place.Rating > 0 {
			t.Logf("Rating: %.1f (%d reviews)", place.Rating, place.UserRatingsTotal)
		}
		if place.PriceLevel > 0 {
			t.Logf("Price Level: %d", place.PriceLevel)
		}
		if len(place.Types) > 0 {
			t.Logf("Types: %v", place.Types)
		}
		if len(place.Photos) > 0 {
			t.Logf("Number of Photos: %d", len(place.Photos))
			t.Logf("First Photo Reference: %s", place.Photos[0].PhotoReference)
		}
		if place.OpeningHours != nil {
			t.Logf("Open Now: %v", place.OpeningHours.OpenNow)
			if place.OpeningHours.WeekdayText != nil {
				t.Logf("Hours: %v", place.OpeningHours.WeekdayText)
			}
		}
		if place.Website != "" {
			t.Logf("Website: %s", place.Website)
		}
		if place.FormattedPhone != "" {
			t.Logf("Phone: %s", place.FormattedPhone)
		}
		if place.BusinessStatus != "" {
			t.Logf("Business Status: %s", place.BusinessStatus)
		}
	}
}

func TestMapsClient_TestConnection(t *testing.T) {
	// Skip this test if we don't have a real API key
	if os.Getenv("GOOGLE_MAPS_API_KEY") == "" {
		t.Skip("Skipping test because GOOGLE_MAPS_API_KEY is not set")
	}

	client, err := NewMapsClient()
	if err != nil {
		t.Fatalf("Failed to create maps client: %v", err)
	}

	ctx := context.Background()
	err = client.TestConnection(ctx)
	if err != nil {
		t.Errorf("TestConnection failed: %v", err)
	}
}

func TestMapsClient_SearchSpecificCoffeeShop(t *testing.T) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test because GOOGLE_MAPS_API_KEY is not set")
	}

	client, err := NewMapsClient()
	if err != nil {
		t.Fatalf("Failed to create maps client: %v", err)
	}

	ctx := context.Background()
	
	tests := []struct {
		name        string
		shopName    string
		location    string
		expectError bool
		minResults  int // minimum number of locations expected
	}{
		{
			name:        "Search for chain with multiple locations",
			shopName:    "Stereoscope Coffee",
			location:    "CA",
			expectError: false,
		},
		{
			name:        "Search with misspelling",
			shopName:    "File Systm of Coffee",
			location:    "Los Angeles",
			expectError: false,
		},
		{
			name:        "Search with partial name",
			shopName:    "Stereoscope",
			location:    "Newport Beach",
			expectError: false,
		},
		{
			name:        "Search with location typo",
			shopName:    "Stereoscope Coffee",
			location:    "NewportBeach",  // Missing space
			expectError: false,
		},
		{
			name:        "Search for restaurant that sells coffee",
			shopName:    "Nep Cafe",
			location:    "Orange County",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add a small delay between tests to avoid rate limiting
			time.Sleep(time.Second)
			
			results, err := client.SearchSpecificCoffeeShop(ctx, tt.shopName, tt.location)
			
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(results) < tt.minResults {
				t.Errorf("expected at least %d results, but got %d", tt.minResults, len(results))
				return
			}

			t.Logf("\nFound %d locations for %q in %s", len(results), tt.shopName, tt.location)
			
			// Log details for each location
			for i, result := range results {
				t.Logf("\n=== Location %d ===", i+1)
				t.Logf("Name: %s", result.Name)
				t.Logf("Place ID: %s", result.PlaceID)
				t.Logf("Address: %s", result.FormattedAddress)
				t.Logf("Vicinity: %s", result.Vicinity)
				t.Logf("Location: %v", result.Location)
				
				if result.Rating > 0 {
					t.Logf("Rating: %.1f (%d reviews)", result.Rating, result.UserRatingsTotal)
				}
				if result.PriceLevel > 0 {
					t.Logf("Price Level: %d", result.PriceLevel)
				}
				if len(result.Types) > 0 {
					t.Logf("Types: %v", result.Types)
				}
				if len(result.Photos) > 0 {
					t.Logf("Number of Photos: %d", len(result.Photos))
					t.Logf("First Photo Reference: %s", result.Photos[0].PhotoReference)
				}
				if result.OpeningHours != nil {
					t.Logf("Open Now: %v", result.OpeningHours.OpenNow)
					if result.OpeningHours.WeekdayText != nil {
						t.Logf("Hours: %v", result.OpeningHours.WeekdayText)
					}
				}
				if result.Website != "" {
					t.Logf("Website: %s", result.Website)
				}
				if result.FormattedPhone != "" {
					t.Logf("Phone: %s", result.FormattedPhone)
				}
				if result.BusinessStatus != "" {
					t.Logf("Business Status: %s", result.BusinessStatus)
				}
			}
		})
	}
}

func TestMapsClient_SearchCoffeeShopsByArea(t *testing.T) {
	apiKey := os.Getenv("GOOGLE_MAPS_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping test because GOOGLE_MAPS_API_KEY is not set")
	}

	client, err := NewMapsClient()
	if err != nil {
		t.Fatalf("Failed to create maps client: %v", err)
	}

	ctx := context.Background()
	
	tests := []struct {
		name        string
		query       string
		expectError bool
	}{
		{
			name:        "Search for matcha places",
			query:       "matcha latte places in Little Tokyo, LA",
			expectError: false,
		},
		{
			name:        "Search for coffee shops that sell matcha places",
			query:       "coffee shops that offer matcha LA",
			expectError: false,
		},
		{
			name:        "Search in non-existent location",
			query:       "coffee in NonExistentCity12345, XY",
			expectError: true,
		},
		{
			name:        "Search with specific requirements",
			query:       "24 hour coffee shops with wifi in Costa Mesa",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Add a small delay between tests to avoid rate limiting
			time.Sleep(time.Second)
			
			results, err := client.SearchCoffeeShopsByArea(ctx, tt.query)
			
			if tt.expectError {
				if err == nil {
					t.Error("expected error but got none")
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
					return
				}
				if len(results) == 0 {
					t.Error("expected results but got none")
					return
				}

				t.Logf("\nFound %d places for query: %s", len(results), tt.query)
				
				// Log details for first few results
				maxToShow := 3
				if len(results) < maxToShow {
					maxToShow = len(results)
				}
				
				for i := 0; i < maxToShow; i++ {
					place := results[i]
					t.Logf("\n=== Place %d ===", i+1)
					t.Logf("Name: %s", place.Name)
					t.Logf("Place ID: %s", place.PlaceID)
					t.Logf("Address: %s", place.FormattedAddress)
					t.Logf("Vicinity: %s", place.Vicinity)
					t.Logf("Location: %v", place.Location)
					
					if place.Rating > 0 {
						t.Logf("Rating: %.1f (%d reviews)", place.Rating, place.UserRatingsTotal)
					}
					if place.PriceLevel > 0 {
						t.Logf("Price Level: %d", place.PriceLevel)
					}
					if len(place.Types) > 0 {
						t.Logf("Types: %v", place.Types)
					}
					if len(place.Photos) > 0 {
						t.Logf("Number of Photos: %d", len(place.Photos))
						t.Logf("First Photo Reference: %s", place.Photos[0].PhotoReference)
					}
					if place.OpeningHours != nil {
						t.Logf("Open Now: %v", place.OpeningHours.OpenNow)
						if place.OpeningHours.WeekdayText != nil {
							t.Logf("Hours: %v", place.OpeningHours.WeekdayText)
						}
					}
					if place.Website != "" {
						t.Logf("Website: %s", place.Website)
					}
					if place.FormattedPhone != "" {
						t.Logf("Phone: %s", place.FormattedPhone)
					}
					if place.BusinessStatus != "" {
						t.Logf("Business Status: %s", place.BusinessStatus)
					}
				}
			}
		})
	}
}