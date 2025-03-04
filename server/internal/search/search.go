package search

import (
	"context"
	"fmt"
	"log"

	"github.com/johnnynu/Coffeehaus/internal/claude"
	"github.com/johnnynu/Coffeehaus/internal/database"
	"github.com/johnnynu/Coffeehaus/internal/maps"
	"github.com/johnnynu/Coffeehaus/internal/shop"
)

type SearchService struct {
	maps *maps.MapsClient
	db *database.Client
	claude *claude.Service
	shops *shop.SyncManager
}

func NewSearchService(maps *maps.MapsClient, db *database.Client, claude *claude.Service, shops *shop.SyncManager) *SearchService {
	return &SearchService{
		maps: maps,
		db: db,
		claude: claude,
		shops: shops,
	}
}

// Search is the entry point for the search service
func (s *SearchService) Search(ctx context.Context, opts SearchOptions) (*SearchResult, error) {
	// default values
	if opts.Limit == 0 {
		opts.Limit = 10
	}

	if opts.Radius == 0 {
		opts.Radius = 10000 // 10km
	}

	// get user location
	userLocation := "unknown"
	if opts.Lat != 0 && opts.Lng != 0 {
		userLocation = fmt.Sprintf("%f,%f", opts.Lat, opts.Lng)
	}

	// use claude to analyze search query
	userIntent, err := s.claude.AnalyzeSearchQuery(ctx, opts.Query, userLocation)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze search query: %w", err)
	}

	// handle search based on intent type
	switch userIntent.SearchType {
	case "specific":
		return s.handleSpecificSearch(ctx, userIntent, opts)
	case "area":
		return s.handleAreaSearch(ctx, opts)
	default:
		return s.handleProximitySearch(ctx, opts)
	}
}

func (s *SearchService) handleSpecificSearch(ctx context.Context, userIntent *claude.SearchIntent, opts SearchOptions) (*SearchResult, error) {
	// try to find shop in db
	dbShop, err := s.db.FindShopsByName(ctx, userIntent.Terms.Shop)
	if err == nil && len(dbShop) > 0 {
		return &SearchResult{
			Shops: dbShop,
		}, nil
	}

	// determine location context for the search
	var locationContext string

	// Case 1: user provided location in SearchOptions
	if opts.Lat != 0 && opts.Lng != 0 {
		// use reverse geocoding to get location name
		location, err := s.maps.ReverseGeocode(ctx, opts.Lat, opts.Lng)
		if err != nil {
			fmt.Printf("failed to reverse geocode: %v\n", err)
		} else {
			locationContext = location
		}
	}

	// Case 2: fall back to location from claude's intent analysis
	if locationContext == "" && userIntent.Location != nil && userIntent.Location.Name != "" {
		locationContext = userIntent.Location.Name
	}

	// if not found in DB, search using places api
	shops, err := s.maps.SearchSpecificCoffeeShop(ctx, opts.Query, locationContext)
	if err != nil {
		return nil, fmt.Errorf("failed to search for specific coffee shop: %w", err)
	}

	s.backgroundSyncShops(shops)

	return &SearchResult{
		Shops: shops,
	}, nil
}

func (s *SearchService) handleAreaSearch(ctx context.Context, opts SearchOptions) (*SearchResult, error) {
	shops, err := s.maps.SearchCoffeeShopsByArea(ctx, opts.Query)
	if err != nil {
		return nil, fmt.Errorf("area search failed: %w", err)
	}

	s.backgroundSyncShops(shops)

	return &SearchResult{
		Shops: shops,
	}, nil
}

func (s *SearchService) handleProximitySearch(ctx context.Context, opts SearchOptions) (*SearchResult, error) {
	// Try database first
	dbShops, err := s.db.FindShopsByLocation(ctx, opts.Lat, opts.Lng, opts.Radius)
	if err == nil && len(dbShops) > 0 {
		return &SearchResult{
			Shops: dbShops,
		}, nil
	}

	// Fallback to Google Maps API if no results in database
	shops, err := s.maps.SearchCoffeeShops(ctx, opts.Lat, opts.Lng, opts.Radius)
	if err != nil {
		return nil, fmt.Errorf("proximity search failed: %w", err)
	}

	s.backgroundSyncShops(shops)

	return &SearchResult{
		Shops: shops,
	}, nil
}

// backgroundSyncShops starts a goroutine to sync shop data to the database
func (s *SearchService) backgroundSyncShops(shops []*maps.CoffeeShopDetails) {
	if len(shops) == 0 {
		return
	}

	go func() {
		inputs := make([]shop.SyncInput, 0, len(shops))
		
		for _, placeShop := range shops {
			input := convertToSyncInput(placeShop)
			inputs = append(inputs, input)
		}
		ctx := context.Background()
		log.Printf("Starting batch sync of %d shops", len(inputs))
		if err := s.shops.BatchSyncShopData(ctx, inputs); err != nil {
			fmt.Printf("failed to batch sync shop data: %v\n", err)
		}
	}()
}

// convertToSyncInput converts a Google Maps CoffeeShopDetails to a SyncInput
func convertToSyncInput(placeShop *maps.CoffeeShopDetails) shop.SyncInput {
	photos := make([]maps.Photo, len(placeShop.Photos))
	for i, photo := range placeShop.Photos {
		photos[i] = maps.Photo{
			PhotoReference: photo.PhotoReference,
			Height: photo.Height,
			Width: photo.Width,
			HTMLAttributions: photo.HTMLAttributions,
		}
	}

	// convert opening hours from Google Maps type to internal type
	var openingHours *maps.OpeningHours
	if placeShop.OpeningHours != nil {
		periods := make([]maps.Period, len(placeShop.OpeningHours.Periods))
		for i, p := range placeShop.OpeningHours.Periods {
			periods[i] = maps.Period{
				Open: maps.TimeOfDay{
					Day: p.Open.Day,
					Time: p.Open.Time,
				},
				Close: maps.TimeOfDay{
					Day: p.Close.Day,
					Time: p.Close.Time,
				},
			}
		}
		openingHours = &maps.OpeningHours{
			WeekdayText: placeShop.OpeningHours.WeekdayText,
			Periods: periods,
		}
	}

	return shop.SyncInput{
		PlaceID: placeShop.PlaceID,
		Name: placeShop.Name,
		FormattedAddress: placeShop.FormattedAddress,
		Vicinity: placeShop.Vicinity,
		Location: maps.LatLng{
			Lat: placeShop.Location.Lat,
			Lng: placeShop.Location.Lng,
		},
		Rating: placeShop.Rating,
		UserRatingsTotal: placeShop.UserRatingsTotal,
		PriceLevel: placeShop.PriceLevel,
		Types: placeShop.Types,
		Photos: photos,
		OpeningHours: openingHours,
		Website: placeShop.Website,
		FormattedPhone: placeShop.FormattedPhone,
		BusinessStatus: placeShop.BusinessStatus,
	}
}
