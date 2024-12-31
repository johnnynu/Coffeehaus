package middleware

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/johnnynu/Coffeehaus/internal/constants"
	"github.com/supabase-community/auth-go"
)

type AuthMiddleware struct {
	client auth.Client
}

// NewAuthMiddleware creates a new AuthMiddleware instance
func NewAuthMiddleware(supabaseURL, supabaseKey string) *AuthMiddleware {
	// Ensure URL is properly formatted
	if !strings.HasPrefix(supabaseURL, "https://") {
		supabaseURL = fmt.Sprintf("https://%s", supabaseURL)
	}
	// Remove any trailing slashes
	supabaseURL = strings.TrimSuffix(supabaseURL, "/")

	// Create auth client with custom auth URL
	client := auth.New(supabaseURL, supabaseKey).WithCustomAuthURL(fmt.Sprintf("%s/auth/v1", supabaseURL))
	log.Printf("Initializing auth middleware with URL: %s", supabaseURL)

	return &AuthMiddleware{client: client}
}

// Authenticate middleware verifies the JWT token from the Authorization header
func (am *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract JWT token
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Println("Authorization header is missing")
			http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
			return
		}

		// Get token from Authorization header
		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == "" {
			log.Println("Invalid token format")
			http.Error(w, "Invalid token format", http.StatusUnauthorized)
			return
		}

		log.Printf("Attempting to verify token: %s...", token[:10])

		// Create authenticated client with token
		authedClient := am.client.WithToken(token)

		// Get user data
		user, err := authedClient.GetUser()
		if err != nil {
			log.Printf("Error getting user: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), constants.UserKey, user)
		
		// Verify context was set correctly
		if val := ctx.Value(constants.UserKey); val != nil {
			log.Printf("Context value set successfully, type: %T", val)
		} else {
			log.Printf("Failed to set context value")
		}
		
		// Call the next handler with the authenticated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
