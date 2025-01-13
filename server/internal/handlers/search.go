package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/johnnynu/Coffeehaus/internal/search"
)

type SearchHandler struct {
	service *search.SearchService
}

func NewSearchHandler(service *search.SearchService) *SearchHandler {
	return &SearchHandler{
		service: service,
	}
}

func (h *SearchHandler) HandleSearch(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// parse query params
	userQuery := r.URL.Query().Get("q")
	if userQuery == "" {
		http.Error(w, "Query parameter 'q' is required", http.StatusBadRequest)
		return
	}

	opts := search.SearchOptions{
		Query: userQuery,
	}

	// parse other optional params
	if lat := r.URL.Query().Get("lat"); lat != "" {
		if parsedLat, err := strconv.ParseFloat(lat, 64); err == nil {
			opts.Lat = parsedLat
		}
	}
	if lng := r.URL.Query().Get("lng"); lng != "" {
		if parsedLng, err := strconv.ParseFloat(lng, 64); err == nil {
			opts.Lng = parsedLng
		}
	}
	if radius := r.URL.Query().Get("radius"); radius != "" {
		if parsedRadius, err := strconv.ParseUint(radius, 10, 64); err == nil {
			opts.Radius = uint(parsedRadius)
		}
	}
	if limit := r.URL.Query().Get("limit"); limit != "" {
		if parsedLimit, err := strconv.Atoi(limit); err == nil {
			opts.Limit = parsedLimit
		}
	}
	if offset := r.URL.Query().Get("offset"); offset != "" {
		if parsedOffset, err := strconv.Atoi(offset); err == nil {
			opts.Offset = parsedOffset
		}
	}

	// perform search
	results, err := h.service.Search(r.Context(), opts)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// return results
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(results); err != nil {
		http.Error(w, "failed to encode results", http.StatusInternalServerError)
		return
	}
}