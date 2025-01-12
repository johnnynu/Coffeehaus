package search

import (
	"github.com/johnnynu/Coffeehaus/internal/maps"
)

type SearchOptions struct {
    Query     string   `json:"query"` // search text from the user
    Lat       float64  `json:"lat,omitempty"` // latitude of the user for location-based search
    Lng       float64  `json:"lng,omitempty"` // longitude of the user for location-based search
    Radius    uint     `json:"radius,omitempty"` // radius of the search in meters
    Limit     int      `json:"limit,omitempty"` // number of results to return
    Offset    int      `json:"offset,omitempty"` // offset of the results to return
}

type SearchResult struct {
    Shops     []maps.CoffeeShopDetails `json:"shops"`
    PageToken string       `json:"nextPageToken,omitempty"`
}

