package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/johnnynu/Coffeehaus/internal/constants"
	"github.com/johnnynu/Coffeehaus/internal/database"
	"github.com/supabase-community/auth-go/types"
)

type UserHandler struct {
	db *database.Client
}

type UpdateProfileRequest struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
}

func NewUserHandler(db *database.Client) *UserHandler {
	return &UserHandler{db: db}
}

// UpdateProfile handles updating a user's profile information
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	// Get username from URL parameter
	username := chi.URLParam(r, "username")
	log.Printf("Updating profile for username: %s", username)
	
	if username == "" {
		log.Println("Username parameter is empty")
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	// Get user from context
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
	log.Printf("Authenticated user ID: %s", userResp.User.ID.String())

	// Parse request body
	var req UpdateProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	log.Printf("Update request data: %+v", req)

	// Verify user owns this profile
	log.Printf("Verifying ownership - Querying user with username: %s", username)
	
	// First, let's log the full query we're about to execute
	query := `id, username`  // Let's include username in the response for better logging
	log.Printf("Executing Supabase query: SELECT %s FROM users WHERE username = '%s'", query, username)
	
	res, status, err := h.db.From("users").
		Select(query, "", false).
		Eq("username", username).
		Single().
		Execute()
	
	log.Printf("Query response - Raw result: %s, Status: %d, Error: %v", string(res), status, err)
	
	if err != nil {
		log.Printf("Database error: %v", err)
		http.Error(w, "Failed to fetch user profile", http.StatusInternalServerError)
		return
	}
	
	if len(res) == 0 || string(res) == "null" {
		log.Printf("No user found with username: %s", username)
		http.Error(w, "Profile not found", http.StatusNotFound)
		return
	}

	var userData struct {
		ID       string `json:"id"`
		Username string `json:"username"`
	}
	if err := json.Unmarshal(res, &userData); err != nil {
		log.Printf("Failed to unmarshal user data: %v. Raw response: %s", err, string(res))
		http.Error(w, "Failed to parse user data", http.StatusInternalServerError)
		return
	}
	log.Printf("Found user - ID: %s, Username: %s", userData.ID, userData.Username)

	if userData.ID != userResp.User.ID.String() {
		log.Printf("Unauthorized - Profile ID: %s, User ID: %s", userData.ID, userResp.User.ID.String())
		http.Error(w, "Unauthorized to update this profile", http.StatusForbidden)
		return
	}

	// Check if new username is available (if username is being changed)
	if req.Username != username {
		log.Printf("Checking availability for new username: %s", req.Username)
		
		// Only select the username column for the check
		checkRes, status, err := h.db.From("users").
			Select("username", "", false).
			Eq("username", req.Username).
			Execute()
		
		if err != nil {
			log.Printf("Failed to check username availability: %v", err)
			http.Error(w, "Failed to check username availability", http.StatusInternalServerError)
			return
		}

		// Log the full response and status for debugging
		log.Printf("Username check response - Status: %d, Response: %s", status, string(checkRes))

		// Parse the response into a slice of user objects
		var users []struct {
			Username string `json:"username"`
		}
		if err := json.Unmarshal(checkRes, &users); err != nil {
			log.Printf("Failed to parse username check response: %v", err)
			http.Error(w, "Failed to check username availability", http.StatusInternalServerError)
			return
		}

		// If we found any users with this username, it's taken
		if len(users) > 0 {
			log.Printf("Username '%s' is already taken (found %d matches)", req.Username, len(users))
			http.Error(w, "Username already taken", http.StatusConflict)
			return
		}

		log.Printf("New username '%s' is available", req.Username)
	}

	// Update profile
	updateData := map[string]interface{}{
		"username":     req.Username,
		"display_name": req.DisplayName,
		"bio":         req.Bio,
	}
	log.Printf("Updating profile with data: %+v", updateData)

	updateRes, _, err := h.db.From("users").Update(updateData, "", "").Eq("id", userResp.User.ID.String()).Execute()
	if err != nil {
		log.Printf("Failed to update profile - Error: %v", err)
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	if len(updateRes) == 0 || string(updateRes) == "null" {
		log.Printf("Update returned no data")
		http.Error(w, "Failed to update profile", http.StatusInternalServerError)
		return
	}

	log.Printf("Profile updated successfully")

	// Return updated profile
	w.Header().Set("Content-Type", "application/json")
	w.Write(updateRes)
} 