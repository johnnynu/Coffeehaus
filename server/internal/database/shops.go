package database

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/johnnynu/Coffeehaus/internal/maps"
	"github.com/supabase-community/postgrest-go"
)

// FindShopsByName searches for coffee shops by name in the database
func (c *Client) FindShopsByName(ctx context.Context, name string) ([]*maps.CoffeeShopDetails, error) {
	resp, _, err := c.From("shops").
		Select("*", "", false).
		Like("name", "%"+name+"%").
		Execute()
	
	if err != nil {
		return nil, fmt.Errorf("failed to find shops by name: %w", err)
	}

	var shop []*maps.CoffeeShopDetails
	if err := json.Unmarshal(resp, &shop); err != nil {
		return nil, fmt.Errorf("failed to parse shops: %w", err)
	}

	return shop, nil
}

// FindShopsByLocation searches for coffee shops within a radius of a point
func (c *Client) FindShopsByLocation(ctx context.Context, lat, lng float64, radiusMeters uint) ([]*maps.CoffeeShopDetails, error) {
	// Create point from coordinates and convert radius to meters
	query := fmt.Sprintf(`ST_DWithin(location::geometry, ST_SetSRID(ST_MakePoint(%f, %f), 4326)::geometry, %d)`, lng, lat, radiusMeters)
	
	resp, _, err := c.From("shops").
		Select("*, ST_Distance(location::geometry, ST_SetSRID(ST_MakePoint(" + fmt.Sprintf("%f, %f", lng, lat) + "), 4326)::geometry) as distance", "", false).
		Filter(query, "eq", "true").
		Order("distance", &postgrest.OrderOpts{Ascending: true}).
		Execute()
	
	if err != nil {
		return nil, fmt.Errorf("failed to find shops by location: %w", err)
	}

	var shops []*maps.CoffeeShopDetails
	if err := json.Unmarshal(resp, &shops); err != nil {
		return nil, fmt.Errorf("failed to parse shops: %w", err)
	}

	return shops, nil
} 

func (c *Client) FindShopByPlaceID(ctx context.Context, placeID string) (*maps.CoffeeShopDetails, error) {
	resp, _, err := c.From("shops").Select("*", "", false).Eq("google_place_id", placeID).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to find shop by place id: %w", err)
	}

	var shop *maps.CoffeeShopDetails
	if err := json.Unmarshal(resp, &shop); err != nil {
		return nil, fmt.Errorf("failed to parse shop: %w", err)
	}

	return shop, nil
}