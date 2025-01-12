package shop

import (
	"context"
	"testing"
	"time"

	"github.com/johnnynu/Coffeehaus/internal/config"
	"github.com/johnnynu/Coffeehaus/internal/database"
	"github.com/johnnynu/Coffeehaus/internal/maps"
	"github.com/joho/godotenv"
)

func init() {
	// Load the .env file before running tests
	if err := godotenv.Load("../../.env"); err != nil {
		println("Warning: Error loading .env file:", err)
	}
}

func setupTestDB(t *testing.T) *database.Client {
	dbConfig, err := config.NewDatabaseConfig()
	if err != nil {
		t.Fatalf("Failed to load database config: %v", err)
	}

	db, err := database.NewClient(dbConfig)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	return db
}

func TestSyncShopData(t *testing.T) {
	db := setupTestDB(t)
	syncManager := NewSyncManager(db)

	tests := []struct {
		name    string
		input   SyncInput
		wantErr bool
	}{
		{
			name: "successful sync new shop",
			input: SyncInput{
				PlaceID:          "test_place_id_" + time.Now().Format("20060102150405"),
				Name:             "Test Coffee Shop",
				FormattedAddress: "123 Test St, Test City, TC 12345",
				Vicinity:         "123 Test St",
				Location: maps.LatLng{
					Lat: 34.0522,
					Lng: -118.2437,
				},
				Rating:           4.5,
				UserRatingsTotal: 100,
				PriceLevel:       2,
				Types:           []string{"cafe", "restaurant"},
				Photos: []maps.Photo{
					{
						PhotoReference: "test_photo_ref_1",
						Height:        1000,
						Width:         1000,
					},
				},
				OpeningHours: &maps.OpeningHours{
					WeekdayText: []string{"Monday: 9:00 AM – 5:00 PM"},
				},
				Website:        "https://testcoffeeshop.com",
				FormattedPhone: "+1 (555) 123-4567",
				BusinessStatus: "OPERATIONAL",
			},
			wantErr: false,
		},
		{
			name: "sync existing shop with updates",
			input: SyncInput{
				PlaceID:          "test_place_id_existing",
				Name:             "Updated Coffee Shop " + time.Now().Format("20060102150405"),
				FormattedAddress: "456 Test Ave, Test City, TC 12345",
				Vicinity:         "456 Test Ave",
				Location: maps.LatLng{
					Lat: 34.0522,
					Lng: -118.2437,
				},
				Rating:           4.8,
				UserRatingsTotal: 200,
				PriceLevel:       2,
				Types:           []string{"cafe"},
				Photos:          []maps.Photo{},
				OpeningHours:    nil,
				Website:         "https://updatedshop.com",
				FormattedPhone: "+1 (555) 999-8888",
				BusinessStatus: "OPERATIONAL",
			},
			wantErr: false,
		},
	}

	// First create the "existing" shop that we'll update in the second test
	existingShop := tests[1].input
	existingShop.Name = "Original Coffee Shop"
	err := syncManager.SyncShopData(context.Background(), existingShop)
	if err != nil {
		t.Fatalf("Failed to create existing shop for update test: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := syncManager.SyncShopData(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("SyncShopData() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Verify the shop exists after sync
			exists, err := syncManager.shopExists(context.Background(), tt.input.PlaceID)
			if err != nil {
				t.Errorf("Failed to check if shop exists: %v", err)
			}
			if !exists {
				t.Error("Shop should exist after sync")
			}
		})
	}
}

func TestShopExists(t *testing.T) {
	db := setupTestDB(t)
	syncManager := NewSyncManager(db)

	// Create a test shop first
	testInput := SyncInput{
		PlaceID:          "test_exists_" + time.Now().Format("20060102150405"),
		Name:             "Test Shop for Exists Check",
		FormattedAddress: "789 Test Rd, Test City, TC 12345",
		Location: maps.LatLng{
			Lat: 34.0522,
			Lng: -118.2437,
		},
	}

	err := syncManager.SyncShopData(context.Background(), testInput)
	if err != nil {
		t.Fatalf("Failed to create test shop: %v", err)
	}

	tests := []struct {
		name     string
		placeID  string
		want     bool
		wantErr  bool
	}{
		{
			name:    "existing shop",
			placeID: testInput.PlaceID,
			want:    true,
			wantErr: false,
		},
		{
			name:    "non-existing shop",
			placeID: "non_existing_id",
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := syncManager.shopExists(context.Background(), tt.placeID)
			if (err != nil) != tt.wantErr {
				t.Errorf("shopExists() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("shopExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCheckForUpdate(t *testing.T) {
	db := setupTestDB(t)
	syncManager := NewSyncManager(db)

	// Create initial shop
	initialShop := SyncInput{
		PlaceID:          "test_update_" + time.Now().Format("20060102150405"),
		Name:             "Initial Shop Name",
		FormattedAddress: "Initial Address",
		Rating:           4.0,
		UserRatingsTotal: 100,
	}

	err := syncManager.SyncShopData(context.Background(), initialShop)
	if err != nil {
		t.Fatalf("Failed to create initial shop: %v", err)
	}

	tests := []struct {
		name    string
		input   SyncInput
		want    bool
		wantErr bool
	}{
		{
			name: "shop needs update",
			input: SyncInput{
				PlaceID:          initialShop.PlaceID,
				Name:             "Updated Shop Name",
				FormattedAddress: "Updated Address",
				Rating:           4.5,
				UserRatingsTotal: 150,
			},
			want:    true,
			wantErr: false,
		},
		{
			name: "no update needed",
			input: SyncInput{
				PlaceID:          initialShop.PlaceID,
				Name:             "Updated Shop Name",
				FormattedAddress: "Updated Address",
				Rating:           4.5,
				UserRatingsTotal: 150,
			},
			want:    false,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := syncManager.checkForUpdate(context.Background(), tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("checkForUpdate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("checkForUpdate() = %v, want %v", got, tt.want)
			}

			// If an update was needed, apply it for the next test
			if got {
				err = syncManager.SyncShopData(context.Background(), tt.input)
				if err != nil {
					t.Errorf("Failed to apply update: %v", err)
				}
			}
		})
	}
} 