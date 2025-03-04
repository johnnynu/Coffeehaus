package redis

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/johnnynu/Coffeehaus/internal/shop"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
)

var testClient *RedisClient

func TestMain(m *testing.M) {
	// Load environment variables from .env file
	if err := godotenv.Load("../../.env"); err != nil {
		log.Printf("Warning: Error loading .env file: %v", err)
	}

	// Setup test Redis client using environment variables
	client, err := NewRedisClient(
        os.Getenv("REDIS_ADDR"),
        os.Getenv("REDIS_PASSWORD"),
        0,
    )

    if err != nil {
        log.Fatalf("Failed to create redis client: %v", err)
    }

    testClient = client

	testClient = client

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func TestCacheShop(t *testing.T) {
    ctx := context.Background()

    testCases := []struct {
        name        string
        shop        *shop.Shop
        expectError bool
    }{
        {
            name: "Stereoscope Newport",
            shop: &shop.Shop{
                ID:               "stereoscope-newport",
                Name:             "Stereoscope Coffee",
                FormattedAddress: "100 S Coast Hwy, Newport Beach, CA 92660",
                Vicinity:         "Newport Beach",
                Location:         "-117.9294,33.6186",
                GoogleRating:     4.7,
                PriceLevel:       2,
            },
            expectError: false,
        },
        {
            name: "Stereoscope Buena Park",
            shop: &shop.Shop{
                ID:               "stereoscope-buena-park",
                Name:             "Stereoscope Coffee",
                FormattedAddress: "7621 Valley View St, Buena Park, CA 90620",
                Vicinity:         "Buena Park",
                Location:         "-118.0241,33.8442",
                GoogleRating:     4.8,
                PriceLevel:       2,
            },
            expectError: false,
        },
        {
            name: "Stereoscope Blue",
            shop: &shop.Shop{
                ID:               "stereoscope-blue",
                Name:             "Stereoscope Coffee at Blue",
                FormattedAddress: "2525 Main St, Irvine, CA 92614",
                Vicinity:         "Irvine",
                Location:         "-117.8520,33.6866",
                GoogleRating:     4.6,
                PriceLevel:       2,
            },
            expectError: false,
        },
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            err := testClient.CacheShop(ctx, tc.shop)
            
            if tc.expectError {
                assert.Error(t, err, "Expected error for test case: %s", tc.name)
                return
            }
            
            assert.NoError(t, err, "Failed to cache shop")

            // Verify the shop was cached by checking if it exists in Redis
            if !tc.expectError {
                exists, err := testClient.client.Exists(ctx, shopKeyPrefix+tc.shop.ID).Result()
                assert.NoError(t, err, "Failed to check if shop exists")
                assert.Equal(t, int64(1), exists, "Shop should exist in Redis")
            }
        })
    }
}

func TestGetCachedShop(t *testing.T) {
    ctx := context.Background()

    t.Run("Get existing shop", func(t *testing.T) {
        retrieved, err := testClient.GetCachedShop(ctx, "test1")
        assert.NoError(t, err, "Failed to get cached shop")
        assert.NotNil(t, retrieved, "Retrieved shop is nil")

        // Verify all fields match
        assert.Equal(t, "test1", retrieved.ID, "ID mismatch")
        assert.Equal(t, "Test Coffee Shop", retrieved.Name, "Name mismatch")
        assert.Equal(t, "123 Test St", retrieved.FormattedAddress, "Address mismatch")
        assert.Equal(t, "Test Area", retrieved.Vicinity, "Vicinity mismatch")
        assert.Equal(t, "-122.4194,37.7749", retrieved.Location, "Location mismatch")
        assert.Equal(t, float32(4.5), retrieved.GoogleRating, "Rating mismatch")
        assert.Equal(t, 2, retrieved.PriceLevel, "Price level mismatch")
    })

    t.Run("Get non-existent shop", func(t *testing.T) {
        shop, err := testClient.GetCachedShop(ctx, "non-existent-id")
        assert.Error(t, err, "Expected error when getting non-existent shop")
        assert.Nil(t, shop, "Expected nil shop when getting non-existent shop")
        assert.Contains(t, err.Error(), "shop not found", "Expected 'shop not found' error message")
    })
}

func TestCreateIndex(t *testing.T) {
    t.Run("Creating index", func(t *testing.T) {
        err := testClient.InitializeShopIndex()
        assert.NoError(t, err, "Index created without errors")
    })
}

func TestSearchSpecific(t *testing.T) {
    ctx := context.Background()
    t.Run("Searching for a specific shop", func(t *testing.T) {
        shops, err := testClient.SearchSpecific(ctx, "Stereoscope")
        assert.NoError(t, err, "failed to search for shop")
        assert.NotNil(t, shops, "search results are nil")
        assert.Len(t, shops, 3, "expected one shop in the result")

        expectedIDs := map[string]bool {
            "stereoscope-newport": true,
            "stereoscope-buena-park": true,
            "stereoscope-blue": true,
        }

        if len(shops) > 0 {
            for _, shop := range shops {
                assert.True(t, expectedIDs[shop.ID], shop.ID, "unexpected shop ID: %s", shop.ID)
                delete(expectedIDs, shop.ID)
            }
        }
    })

    t.Run("Searching for a specific shop with typos", func (t *testing.T) {
        shops, err := testClient.SearchSpecific(ctx, "Ster")
        assert.NoError(t, err, "failed to search shop")
        assert.NotNil(t, shops, "search results are nil")
        assert.Len(t, shops, 3, "expected 3 shops in search results")

        expectedIDs := map[string]bool{
            "stereoscope-newport": true,
            "stereoscope-buena-park": true,
            "stereoscope-blue": true,
        }

        if len(shops) > 0 {
            for _, shop := range shops {
                assert.True(t, expectedIDs[shop.ID], shop.ID, "unexpected shop ID: %s", shop.ID)
                delete(expectedIDs, shop.ID)
            }
        }
    })
}