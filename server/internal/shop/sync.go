package shop

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"
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

func (s *SyncManager) BatchSyncShopData(ctx context.Context, inputs []SyncInput) error {
	if len(inputs) == 0 {
		return nil
	}

	// extract all place ids to check if a shop exists
	placeIDs := make([]string, len(inputs))
	for i, input := range inputs {
		placeIDs[i] = input.PlaceID
	}

	// find all existing shops
	existingShops, err := s.getExistingShops(ctx, placeIDs)
	if err != nil {
		return fmt.Errorf("failed to check existing shops: %w", err)
	}
	fmt.Printf("Existing shops count: %d\n", len(existingShops))

	// separate shops that need to be created vs updated
	var toCreate []SyncInput
	var toUpdate []SyncInput

	// map of placeID -> existing shop
	existingMap := make(map[string]ExistingShop)
	for _, shop := range existingShops {
		existingMap[shop.GooglePlaceID] = shop
	}

	for _, input := range inputs {
		existing, exists := existingMap[input.PlaceID]
		if !exists {
			toCreate = append(toCreate, input)
			continue
		}

		// check if shop needs update
		if needsUpdate(existing, input) {
			toUpdate = append(toUpdate, input)
		}
	}

	fmt.Printf("Shops to create: %d\n", len(toCreate))
	fmt.Printf("Shops to update: %d\n", len(toUpdate))
	// create new shops in batch
	if len(toCreate) > 0 {
		if err := s.batchCreateShops(ctx, toCreate); err != nil {
			return fmt.Errorf("failed to batch create shops: %w", err)
		}
	}

	// update existing shops in batch
	if len(toUpdate) > 0  {
		if err := s.batchUpdateShops(ctx, toUpdate); err != nil {
			return fmt.Errorf("failed to batch update shops: %w", err)
		}
	}

	return nil
}

func (s *SyncManager) getExistingShops(ctx context.Context, placeIDs []string) ([]ExistingShop, error) {
	_ = ctx

	// this query uses the "in" filter to find all shops with the given place IDs
	placeIDsStr := fmt.Sprintf("(%s)", strings.Join(placeIDs, ","))

	query := "id, google_place_id, name, formatted_address, vicinity, google_rating, ratings_total, price_level, website, formatted_phone, business_status"

	res, _, err := s.db.From("shops").Select(query, "", false).Filter("google_place_id", "in", placeIDsStr).Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to get existing shops: %w", err)
	}

	var existingShops []ExistingShop
	if err := json.Unmarshal(res, &existingShops); err != nil {
		return nil, fmt.Errorf("failed to unmarshal existing shops: %w", err)
	}

	return existingShops, nil
}

// needsUpdate checks if a shop needs to be updated
func needsUpdate(existing ExistingShop, input SyncInput) bool {
	return (input.Name != "" && existing.Name != input.Name) ||
	(input.FormattedAddress != "" && existing.FormattedAddress != input.FormattedAddress) ||
	(input.Vicinity != "" && existing.Vicinity != input.Vicinity) ||
	(input.Rating > 0 && existing.GoogleRating != input.Rating) ||
	(input.UserRatingsTotal > 0 && existing.RatingsTotal != input.UserRatingsTotal) ||
	(input.PriceLevel > 0 && existing.PriceLevel != input.PriceLevel) ||
	(input.Website != "" && existing.Website != input.Website) ||
	(input.FormattedPhone != "" && existing.FormattedPhone != input.FormattedPhone) ||
	(input.BusinessStatus != "" && existing.BusinessStatus != input.BusinessStatus)
}

func (s *SyncManager) batchCreateShops(ctx context.Context, inputs []SyncInput) error {
	_ = ctx

	if len(inputs) == 0 {
		return nil
	}

	shopDataBatch := make([]map[string]interface{}, len(inputs))

	for i, input := range inputs {
		photoRefs := make([]string, len(input.Photos))
		for j, photo := range input.Photos {
			photoRefs[j] = photo.PhotoReference
		}

		// Create point from location
		point := fmt.Sprintf("(%f,%f)", input.Location.Lat, input.Location.Lng)
		
		shopDataBatch[i] = map[string]interface{}{
			"google_place_id":    input.PlaceID,
			"name":               input.Name,
			"formatted_address":  input.FormattedAddress,
			"vicinity":           input.Vicinity,
			"location":           point,
			"google_rating":      input.Rating,
			"ratings_total":      input.UserRatingsTotal,
			"price_level":        input.PriceLevel,
			"types":              input.Types,
			"photo_refs":         photoRefs,
			"hours":              input.OpeningHours,
			"website":            input.Website,
			"formatted_phone":    input.FormattedPhone,
			"business_status":    input.BusinessStatus,
			"last_sync":          time.Now(),
			"coffeehaus_rating":  nil,
			"verified":           false,
		}		
	}

	// batch insert
	log.Printf("Batch inserting %d shops", len(shopDataBatch))
	_, _, err := s.db.From("shops").Insert(shopDataBatch, false, "", "", "").Execute()
	if err != nil {
		return fmt.Errorf("failed to batch create shops: %w", err)
	}
	return nil
}

// batchUpdateShops updates multiple shops
// Note: Postgrest doesnt support true batch updates so multiple requests are required
func (s *SyncManager) batchUpdateShops(ctx context.Context, inputs []SyncInput) error {
	_ = ctx

	if len(inputs) == 0 {
		return nil
	}

	// Group shops by update pattern to minimize requests
	log.Printf("Processing batch update for %d shops", len(inputs))

	batchSize := 25
	for i := 0; i < len(inputs); i += batchSize {
		end := i + batchSize
		if end > len(inputs) {
			end = len(inputs)
		}
		batch := inputs[i:end]

		// process for each batch
		for _, input := range batch {
			photoRefs := make([]string, len(input.Photos))
			for j, photo := range input.Photos {
				photoRefs[j] = photo.PhotoReference
			}

			point := fmt.Sprintf("(%f,%f)", input.Location.Lat, input.Location.Lng)

			updateData := map[string]interface{}{
				"name":              input.Name,
				"formatted_address": input.FormattedAddress,
				"vicinity":          input.Vicinity,
				"location":          point,
				"google_rating":     input.Rating,
				"ratings_total":     input.UserRatingsTotal,
				"price_level":       input.PriceLevel,
				"types":             input.Types,
				"photo_refs":        photoRefs,
				"hours":             input.OpeningHours,
				"website":           input.Website,
				"formatted_phone":   input.FormattedPhone,
				"business_status":   input.BusinessStatus,
				"last_sync":         time.Now(),
			}

			_, _, err := s.db.From("shops").Update(updateData, "", "").Eq("google_place_id", input.PlaceID).Execute()
			if err != nil {
				return fmt.Errorf("failed to update shop %s: %w", input.PlaceID, err)
			}
		}
	}
	return nil
}

/*


	Old methods kept for backward compatibility


*/

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