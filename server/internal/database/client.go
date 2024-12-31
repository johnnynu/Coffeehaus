package database

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/johnnynu/Coffeehaus/internal/config"

	"github.com/supabase-community/postgrest-go"
)

type Client struct {
	*postgrest.Client
	config *config.DatabaseConfig
}

func NewClient(cfg *config.DatabaseConfig) (*Client, error) {
	// Create headers with Supabase service role key
	headers := map[string]string{
		"apikey":        cfg.ServiceRoleKey,
		"Authorization": "Bearer " + cfg.ServiceRoleKey,
		"Content-Type": "application/json",
		"Accept":       "application/json",
		"Prefer":       "return=representation",
	}

	// Initialize postgrest client
	client := postgrest.NewClient(cfg.RestURL, "public", headers)
	if client.ClientError != nil {
		return nil, fmt.Errorf("failed to initialize postgrest client: %w", client.ClientError)
	}

	db := &Client{
		Client: client,
		config: cfg,
	}

	// Test connection
	if err := db.TestConnection(); err != nil {
		return nil, fmt.Errorf("database connection test failed: %w", err)
	}

	return db, nil
}

// TestConnection verifies we can connect to Supabase
func (c *Client) TestConnection() error {
	resp, _, err := c.From("users").
		Select("count", "", false).
		Execute()
	
	if err != nil {
		return fmt.Errorf("connection test query failed: %w", err)
	}

	var result []map[string]interface{}
	if err := json.Unmarshal(resp, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result) == 0 {
		return fmt.Errorf("unexpected empty response from database")
	}

	log.Printf("Successfully connected to Supabase")
	return nil
}