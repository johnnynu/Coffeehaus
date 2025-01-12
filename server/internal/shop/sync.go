package shop

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/johnnynu/Coffeehaus/internal/database"
)

// SyncManager handles synchronization of shop data between Plces API and db
type SyncManager struct {
	db *database.Client
}

// NewSyncManager creates a new SyncManager
func NewSyncManager(db *database.Client) *SyncManager {
	return &SyncManager{db: db}
}

// SyncShopData syncs shop data from Places API to db
func (s *SyncManager) SyncShopData(ctx context.Context, input SyncInput) error {
	// check if shop exists in db
	exists, err := s.shopExists(ctx, input.PlaceID)
	if err != nil {
		return fmt.Errorf("failed to check if shop exists: %w", err)
	}

	if !exists {
		// create new shop
		return s.createShop(ctx, input)
	}

	// check if shop needs to be updated
	update, err := s.checkForUpdate(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to check for update: %w", err)
	}

	if update {
		return s.updateShop(ctx, input)
	}

	return nil
}

func (s *SyncManager) shopExists(ctx context.Context, placeID string) (bool, error) {
	_ = ctx
	res, _, err := s.db.From("shops").Select("id", "", false).Eq("google_place_id", placeID).Execute()

	if err != nil {
		return false, fmt.Errorf("querying shop: %w", err)
	}

	// Parse the response to check if it's an empty array
	var results []map[string]interface{}
	if err := json.Unmarshal(res, &results); err != nil {
		return false, fmt.Errorf("parsing response: %w", err)
	}

	return len(results) > 0, nil
}

func (s *SyncManager) createShop(ctx context.Context, input SyncInput) error {
	_ = ctx
	// extract photo references from Photos slice
	photoRefs := make([]string, len(input.Photos))
	for i, photo := range input.Photos {
		photoRefs[i] = photo.PhotoReference
	}

	// create the point from the Location (PostgreSQL point format)
	point := fmt.Sprintf("(%f,%f)", input.Location.Lat, input.Location.Lng)

	shopData := map[string]interface{}{
		"google_place_id": input.PlaceID,
		"name": input.Name,
		"formatted_address": input.FormattedAddress,
		"vicinity": input.Vicinity,
		"location": point,
		"google_rating": input.Rating,
		"ratings_total": input.UserRatingsTotal,
		"price_level": input.PriceLevel,
		"types": input.Types,
		"photo_refs": photoRefs,
		"hours": input.OpeningHours,
		"website": input.Website,
		"formatted_phone": input.FormattedPhone,
		"business_status": input.BusinessStatus,
		"last_sync": time.Now(),
		
		// coffeehaus specific fields
		"coffeehaus_rating": nil,
		"verified": false,
	}

	_, _, err := s.db.From("shops").Insert(shopData, false, "", "", "").Execute()

	if err != nil {
		return fmt.Errorf("failed to create shop: %w", err)
	}

	return nil
}

func (s *SyncManager) checkForUpdate(ctx context.Context, input SyncInput) (bool, error) {
	_ = ctx
	res, _, err := s.db.From("shops").Select(`name, formatted_address, vicinity, google_rating, ratings_total, price_level, website, formatted_phone, business_status`, "", false).Eq("google_place_id", input.PlaceID).Execute()

	if err != nil {
		return false, fmt.Errorf("failed to check for update: %w", err)
	}

    var existingShops []struct {
        Name            string  `json:"name"`
        FormattedAddress string `json:"formatted_address"`
        Vicinity        string  `json:"vicinity"`
        GoogleRating    float32 `json:"google_rating"`
        RatingsTotal    int     `json:"ratings_total"`
        PriceLevel      int     `json:"price_level"`
        Website         string  `json:"website"`
        FormattedPhone  string  `json:"formatted_phone"`
        BusinessStatus  string  `json:"business_status"`
    }

	if err := json.Unmarshal(res, &existingShops); err != nil {
		return false, fmt.Errorf("failed to unmarshal existing shop: %w", err)
	}

	// If no shops found, return false without error
	if len(existingShops) == 0 {
		return false, nil
	}

	shop := existingShops[0]

	return (input.Name != "" && shop.Name != input.Name) ||
		(input.FormattedAddress != "" && shop.FormattedAddress != input.FormattedAddress) ||
		(input.Vicinity != "" && shop.Vicinity != input.Vicinity) ||
		(input.Rating > 0 && shop.GoogleRating != input.Rating) ||
		(input.UserRatingsTotal > 0 && shop.RatingsTotal != input.UserRatingsTotal) ||
		(input.PriceLevel > 0 && shop.PriceLevel != input.PriceLevel) ||
		(input.Website != "" && shop.Website != input.Website) ||
		(input.FormattedPhone != "" && shop.FormattedPhone != input.FormattedPhone) ||
		(input.BusinessStatus != "" && shop.BusinessStatus != input.BusinessStatus), nil
}

func (s *SyncManager) updateShop(ctx context.Context, input SyncInput) error {
	_ = ctx
	// extract photo refs from Photos slice
	photoRefs := make([]string, len(input.Photos))
	for i, photo := range input.Photos {
		photoRefs[i] = photo.PhotoReference
	}

	// create the point from the Location
	point := fmt.Sprintf("(%f,%f)", input.Location.Lat, input.Location.Lng)

	updateData := map[string]interface{}{
		"name": input.Name,
		"formatted_address": input.FormattedAddress,
		"vicinity": input.Vicinity,
		"location": point,
		"google_rating": input.Rating,
		"ratings_total": input.UserRatingsTotal,
		"price_level": input.PriceLevel,
		"types": input.Types,
		"photo_refs": photoRefs,
		"hours": input.OpeningHours,
		"website": input.Website,
		"formatted_phone": input.FormattedPhone,
		"business_status": input.BusinessStatus,
		"last_sync": time.Now(),
	}

	_, _, err := s.db.From("shops").Update(updateData, "", "").Eq("google_place_id", input.PlaceID).Execute()

	if err != nil {
		return fmt.Errorf("failed to update shop: %w", err)
	}

	return nil
}