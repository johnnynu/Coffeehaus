package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/johnnynu/Coffeehaus/internal/shop"
	"github.com/redis/go-redis/v9"
)

type SearchResult struct {
	NormalizedQuery string `json:"normalized_query"`
	SearchType      string `json:"search_type"`
	Location        *SearchLocation `json:"location,omitempty"`
	Shops           []*shop.Shop `json:"shops"`
	Timestamp       time.Time `json:"timestamp"`
}

type SearchLocation struct {
	Name     string `json:"name"`
	Radius   float64 `json:"radius"`
}

const (
	searchKeyPrefix = "search:"
	searchTTL = 6 * time.Hour
)

// InitializeSearchIndex creates the search index for search results
func (r *RedisClient) InitializeSearchIndex() error {
	// Create search index for JSON documents
	ctx := context.Background()

	normalizedQueryField := &redis.FieldSchema{
		FieldName: "$.normalized_query",
		FieldType: redis.SearchFieldTypeText,
		Sortable: true,
		As: "normalized_query",
	}

	err := r.client.FTCreate(ctx, "searchIdx", &redis.FTCreateOptions{
		OnJSON: true,
		Prefix: []interface{}{searchKeyPrefix},
	}, normalizedQueryField).Err(); 
	
	if err != nil {
		log.Printf("Failed to create search index or already exists: %v", err)
	}

	return nil
}

// CacheSearchResults stores search results using RedisJSON
func (r *RedisClient) CacheSearchResults(ctx context.Context, normalizedQuery string, result *SearchResult) error {
	key := searchKeyPrefix + normalizedQuery

	err := r.client.JSONSet(ctx, key, "$", result).Err()
	if err != nil {
		fmt.Printf("Failed to cache search result: %v", err)
	}

	// set ttl
	err = r.client.Expire(ctx, key, searchTTL).Err()
	if err != nil {
		fmt.Printf("Failed to set TTL for search result: %v", err)
	}

	return nil
}

func (r *RedisClient) GetSearchResults(ctx context.Context, normalizedQuery string) (*SearchResult, error) {
	key := searchKeyPrefix + normalizedQuery

	res, err := r.client.JSONGet(ctx, key, "$").Result()
	if err != nil {
		return nil, fmt.Errorf("failed to get search result: %w", err)
	}

	var result SearchResult
	err = json.Unmarshal([]byte(res), &result)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal search result: %w", err)
	}

	return &result, nil
}