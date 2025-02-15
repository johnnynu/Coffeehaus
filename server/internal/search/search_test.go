package search

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/johnnynu/Coffeehaus/internal/claude"
	"github.com/johnnynu/Coffeehaus/internal/config"
	"github.com/johnnynu/Coffeehaus/internal/database"
	"github.com/johnnynu/Coffeehaus/internal/maps"
	"github.com/johnnynu/Coffeehaus/internal/shop"
	"github.com/joho/godotenv"
)

func init() {
	// Load the .env file before running tests
	if err := godotenv.Load("../../.env"); err != nil {
		println("Warning: Error loading .env file:", err)
	}
}

func setupTestService(t *testing.T) *SearchService {
	// Initialize Maps client
	mapsClient, err := maps.NewMapsClient()
	if err != nil {
		t.Fatalf("Failed to create maps client: %v", err)
	}

	// Initialize DB client
	dbConfig, err := config.NewDatabaseConfig()
	if err != nil {
		t.Fatalf("Failed to create database config: %v", err)
	}
	
	dbClient, err := database.NewClient(dbConfig)
	if err != nil {
		t.Fatalf("Failed to create database client: %v", err)
	}

	// Initialize Claude service
	claudeService := claude.NewService(os.Getenv("CLAUDE_API_KEY"))

	// Initialize Sync manager
	syncManager := shop.NewSyncManager(dbClient)

	return NewSearchService(mapsClient, dbClient, claudeService, syncManager)
}

func TestSearch_Specific(t *testing.T) {
	if os.Getenv("GOOGLE_MAPS_API_KEY") == "" {
		t.Skip("Skipping test because GOOGLE_MAPS_API_KEY is not set")
	}

	service := setupTestService(t)

	tests := []struct {
		name          string
		query         string
		lat           float64
		lng           float64
		wantErr       bool
		minLocations  int // minimum number of locations expected
	}{
		{
			name:         "search for specific shop",
			query:        "Stereoscope Coffee",
			lat:         33.6189,
			lng:         -117.9289,
			minLocations: 2, // Stereoscope has multiple locations
		},
		{
			name:         "search with location context",
			query:        "Nep Cafe",
			lat:         33.7514,
			lng:         -117.9940,
			minLocations: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Search(context.Background(), SearchOptions{
				Query: tt.query,
				Lat:   tt.lat,
				Lng:   tt.lng,
			})

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("expected result but got nil")
				return
			}

			if len(result.Shops) == 0 {
				t.Error("expected shops but got none")
				return
			}

			if len(result.Shops) < tt.minLocations {
				t.Errorf("expected at least %d locations but got %d", tt.minLocations, len(result.Shops))
			}

			// Log all locations found
			t.Logf("Found %d locations for query: %s", len(result.Shops), tt.query)
			for i, shop := range result.Shops {
				t.Logf("Location %d: %s at %s", i+1, shop.Name, shop.FormattedAddress)
			}

			// Wait for background sync to complete
			time.Sleep(5 * time.Second)
			t.Log("Waited for background sync to complete")
		})
	}
}

func TestSearch_Area(t *testing.T) {
	if os.Getenv("GOOGLE_MAPS_API_KEY") == "" {
		t.Skip("Skipping test because GOOGLE_MAPS_API_KEY is not set")
	}

	service := setupTestService(t)

	tests := []struct {
		name         string
		query        string
		minLocations int
		wantErr      bool
	}{
		{
			name:         "search in area",
			query:        "coffee shops that offer strawberry matcha in OC",
			minLocations: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Search(context.Background(), SearchOptions{
				Query: tt.query,
			})

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("expected result but got nil")
				return
			}

			if len(result.Shops) == 0 {
				t.Error("expected shops but got none")
				return
			}

			if len(result.Shops) < tt.minLocations {
				t.Errorf("expected at least %d locations but got %d", tt.minLocations, len(result.Shops))
			}

			// Log all locations found
			t.Logf("Found %d locations for query: %s", len(result.Shops), tt.query)
			for i, shop := range result.Shops {
				t.Logf("Location %d: %s", i+1, shop.Name)
				t.Logf("  Address: %s", shop.FormattedAddress)
				if shop.Rating > 0 {
					t.Logf("  Rating: %.1f (%d reviews)", shop.Rating, shop.UserRatingsTotal)
				}
				if shop.OpeningHours != nil && len(shop.OpeningHours.WeekdayText) > 0 {
					t.Log("  Hours:")
					for _, hours := range shop.OpeningHours.WeekdayText {
						t.Logf("    %s", hours)
					}
				}
				t.Log("") // Empty line between shops
			}

			// Wait for background sync to complete
			time.Sleep(5 * time.Second)
			t.Log("Waited for background sync to complete")
		})
	}
}

func TestSearch_Proximity(t *testing.T) {
	if os.Getenv("GOOGLE_MAPS_API_KEY") == "" {
		t.Skip("Skipping test because GOOGLE_MAPS_API_KEY is not set")
	}

	service := setupTestService(t)

	tests := []struct {
		name         string
		lat          float64
		lng          float64
		radius       uint
		minLocations int
		wantErr      bool
	}{
		{
			name:         "search near Westminster",
			lat:         33.7514,
			lng:         -117.9940,
			radius:      3000,
			minLocations: 5, // Should find several coffee shops within 3km
		},
		{
			name:         "search near Newport Beach",
			lat:         33.6189,
			lng:         -117.9289,
			radius:      5000,
			minLocations: 8, // Should find more shops with larger radius
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := service.Search(context.Background(), SearchOptions{
				Lat:    tt.lat,
				Lng:    tt.lng,
				Radius: tt.radius,
			})

			if tt.wantErr {
				if err == nil {
					t.Error("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Error("expected result but got nil")
				return
			}

			if len(result.Shops) == 0 {
				t.Error("expected shops but got none")
				return
			}

			if len(result.Shops) < tt.minLocations {
				t.Errorf("expected at least %d locations but got %d", tt.minLocations, len(result.Shops))
			}

			// Log all locations found
			t.Logf("Found %d locations within %.1f km of (%.4f, %.4f)", 
				len(result.Shops), float64(tt.radius)/1000, tt.lat, tt.lng)
			for i, shop := range result.Shops {
				t.Logf("Location %d: %s", i+1, shop.Name)
				t.Logf("  Address: %s", shop.FormattedAddress)
				if shop.Rating > 0 {
					t.Logf("  Rating: %.1f (%d reviews)", shop.Rating, shop.UserRatingsTotal)
				}
				if shop.OpeningHours != nil && len(shop.OpeningHours.WeekdayText) > 0 {
					t.Log("  Hours:")
					for _, hours := range shop.OpeningHours.WeekdayText {
						t.Logf("    %s", hours)
					}
				}
				t.Log("") // Empty line between shops
			}

			// Wait for background sync to complete
			time.Sleep(5 * time.Second)
			t.Log("Waited for background sync to complete")
		})
	}
} 