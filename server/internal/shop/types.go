package shop

import (
	"time"

	"github.com/johnnynu/Coffeehaus/internal/maps"
)

// Shop represents a coffee shop in our database, combining Places API data
// with Coffeehaus-specific information
type Shop struct {
    ID               string     `json:"id"`
    GooglePlaceID    string     `json:"google_place_id"`
    Name             string     `json:"name"`
    FormattedAddress string     `json:"formatted_address"`
    Vicinity         string     `json:"vicinity"`
    Location         string     `json:"location"` // PostGIS POINT type
    GoogleRating     float32    `json:"google_rating"`
    RatingsTotal     int        `json:"ratings_total"`
    PriceLevel       int        `json:"price_level"`
    Types            []string   `json:"types"`
    PhotoRefs        []string   `json:"photo_refs"`
    OpeningHours     *maps.OpeningHours `json:"opening_hours,omitempty"`
    Website          string     `json:"website"`
    FormattedPhone   string     `json:"formatted_phone"`
    BusinessStatus   string     `json:"business_status"`
    
    // Coffeehaus-specific fields
    CoffeehausRating *float32   `json:"coffeehaus_rating"`
    LastSync         time.Time  `json:"last_sync"`
    Verified         bool       `json:"verified"`
}

// SyncInput represents the data we receive from the Places API
type SyncInput struct {
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

// ExistingShop represents the minimal data needed to check for updates
type ExistingShop struct {
	ID              string  `json:"id"`
	GooglePlaceID   string  `json:"google_place_id"`
	Name            string  `json:"name"`
	FormattedAddress string  `json:"formatted_address"`
	Vicinity        string  `json:"vicinity"`
	GoogleRating    float32 `json:"google_rating"`
	RatingsTotal    int     `json:"ratings_total"`
	PriceLevel      int     `json:"price_level"`
	Website         string  `json:"website"`
	FormattedPhone  string  `json:"formatted_phone"`
	BusinessStatus  string  `json:"business_status"`
}