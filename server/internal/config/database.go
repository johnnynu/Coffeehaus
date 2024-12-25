package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type DatabaseConfig struct {
	SupabaseURL      string
	SupabaseKey      string
	ServiceRoleKey   string
	RestURL          string
}

func NewDatabaseConfig() (*DatabaseConfig, error) {
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found")
	}

	url := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_KEY")
	serviceKey := os.Getenv("SUPABASE_SERVICE_KEY")

	if url == "" || key == "" {
		return nil, fmt.Errorf("SUPABASE_URL and SUPABASE_KEY must be set")
	}

	restURL := fmt.Sprintf("%s/rest/v1", url)

	return &DatabaseConfig{
		SupabaseURL:    url,
		SupabaseKey:    key,
		ServiceRoleKey: serviceKey,
		RestURL:        restURL,
	}, nil
} 