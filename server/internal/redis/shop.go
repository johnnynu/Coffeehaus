package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/RediSearch/redisearch-go/redisearch"
	"github.com/johnnynu/Coffeehaus/internal/shop"
	redisClient "github.com/redis/go-redis/v9"
)

const shopKeyPrefix = "shop:"

func (r *RedisClient) InitializeShopIndex() error {
	shopSchema := redisearch.NewSchema(redisearch.DefaultOptions).
		AddField(redisearch.NewTextField("name")).
		AddField(redisearch.NewTextField("formatted_address")).
		AddField(redisearch.NewTextField("vicinity")).
		AddField(redisearch.NewGeoField("location")).
		AddField(redisearch.NewNumericField("google_rating")).
		AddField(redisearch.NewNumericField("price_level"))

	if err := r.searchClient.CreateIndex(shopSchema); err != nil {
		log.Printf("error creating index: %s", err)
		return err
	}

	log.Printf("Successfully created shop index")
	return nil
}

func (r *RedisClient) CacheShop(ctx context.Context, shop *shop.Shop) error {
    // Handle nil shop
    if shop == nil {
        return fmt.Errorf("cannot cache nil shop")
    }

    // Validate shop ID
    if shop.ID == "" {
        return fmt.Errorf("shop ID cannot be empty")
    }

    key := shopKeyPrefix + shop.ID
    log.Printf("Attempting to cache shop with key: %s", key)

    // Store the shop object directly as JSON
    err := r.client.JSONSet(ctx, key, "$", shop).Err()
    if err != nil {
        return fmt.Errorf("failed to cache shop: %w", err)
    }

	doc := redisearch.NewDocument(shop.ID, 1.0)
	doc.Set("name", shop.Name).
		Set("formatted_address", shop.FormattedAddress).
		Set("vicinity", shop.Vicinity).
		Set("location", shop.Location).
		Set("google_rating", shop.GoogleRating).
		Set("price_level", shop.PriceLevel)

	if err = r.searchClient.Index([]redisearch.Document{doc}...); err != nil {
		fmt.Printf("Failed to index shop: %s", err)
	}

	log.Printf("Successfully cached and indexed shop with ID: %s", shop.ID)
    return nil
}

func (r *RedisClient) GetCachedShop(ctx context.Context, shopID string) (*shop.Shop, error) {
    key := shopKeyPrefix + shopID

    // Get the JSON data using root path "."
    jsonStr, err := r.client.JSONGet(ctx, key, "$").Result()
    if err != nil {
        if err == redisClient.Nil {
            return nil, fmt.Errorf("shop not found: %s", shopID)
        }
        return nil, fmt.Errorf("failed to get cached shop: %w", err)
    }

	if jsonStr == "" {
		return nil, fmt.Errorf("shop not found: %s", shopID)
	}

	// Redis JSON.GET returns an array with a single object
    // Remove the array brackets
    jsonStr = strings.TrimPrefix(jsonStr, "[")
    jsonStr = strings.TrimSuffix(jsonStr, "]")

    var shop shop.Shop
    err = json.Unmarshal([]byte(jsonStr), &shop)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal shop: %w", err)
    }

    return &shop, nil
}

func (r *RedisClient) SearchSpecific(ctx context.Context, shopName string) ([]*shop.Shop, error) {
	query := redisearch.NewQuery(fmt.Sprintf(`@name:"%s"`, shopName)).Limit(0, 10)

	docs, total, err := r.searchClient.Search(query)
	if err != nil {
		return nil, fmt.Errorf("search failed: %v", err)
	}

	log.Printf("Found %d shops matching the query", total)

	var shops []*shop.Shop
	for _, doc := range docs {
		shopID := doc.Id
		shop, err := r.GetCachedShop(ctx, shopID)
		if err != nil {
			log.Printf("Failed to retrieve shop from cache: %s", err)
			continue
		}
		shops = append(shops, shop)
	}

	return shops, nil
}