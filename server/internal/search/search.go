package search

import (
	"github.com/johnnynu/Coffeehaus/internal/claude"
	"github.com/johnnynu/Coffeehaus/internal/database"
	"github.com/johnnynu/Coffeehaus/internal/maps"
)

type SearchService struct {
	maps *maps.MapsClient
	db *database.Client
	claude *claude.Service
}

func NewSearchService(maps *maps.MapsClient, db *database.Client, claude *claude.Service) *SearchService {
	return &SearchService{
		maps: maps,
		db: db,
		claude: claude,
	}
}
