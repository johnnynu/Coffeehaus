package handlers

import (
	"log"
	"net/http"

	"github.com/johnnynu/Coffeehaus/internal/constants"
	"github.com/johnnynu/Coffeehaus/internal/database"
	"github.com/supabase-community/auth-go/types"
)

type AuthHandler struct {
	db *database.Client
}

func NewAuthHandler(db *database.Client) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) GetUser(w http.ResponseWriter, r *http.Request) {

	// get user from context
	val := r.Context().Value(constants.UserKey)
	if val == nil {
		log.Println("Context value is nil")
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

	userResp, ok := val.(*types.UserResponse)
	if !ok {
		log.Printf("Type assertion failed. Got type: %T, want: *types.UserResponse", val)
		http.Error(w, "User not found in context", http.StatusInternalServerError)
		return
	}

		query := `username, display_name, profile_photo_id, photos!profile_photo_id(versions)`
	res, status, err := h.db.From("users").Select(query, "", false).Eq("id", userResp.User.ID.String()).Single().Execute()

	log.Printf("Query result - Status: %d, Error: %v, Response: %s", status, err, string(res))

	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Failed to fetch user data", http.StatusInternalServerError)
		return
	}

	// Check if we got a valid response
	if len(res) == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// write user data to response
	w.Header().Set("Content-Type", "application/json")
	w.Write(res)
}
