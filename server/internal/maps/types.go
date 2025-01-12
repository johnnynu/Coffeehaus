package maps

import (
	"time"

	"googlemaps.github.io/maps"
)

type OpeningHours struct {
    WeekdayText []string
    Periods     []Period
}

type Period struct {
    Open  TimeOfDay
    Close TimeOfDay
}

type TimeOfDay struct {
    Day  time.Weekday
    Time string // Format: "HHMM", e.g. "0900"
}

type LatLng struct {
    Lat float64
    Lng float64
}

type Photo struct {
    PhotoReference string
    Height         int
    Width          int
    HTMLAttributions []string
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