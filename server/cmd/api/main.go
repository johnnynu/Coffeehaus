package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/johnnynu/Coffeehaus/internal/config"
	"github.com/johnnynu/Coffeehaus/internal/database"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Initialize database
	dbConfig, err := config.NewDatabaseConfig()
	if err != nil {
		log.Fatalf("Failed to load database config: %v", err)
	}

	log.Printf("Connecting to supabase rest url: %s", dbConfig.RestURL)

	db, err := database.NewClient(dbConfig)
	if err != nil {
		log.Printf("db connection details: %+v", err)
		log.Fatalf("Failed to initialize database: %v", err)
	}

	r := chi.NewRouter()

	// Middleware
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:5173"}, // Vite's default port
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	r.Route("/api/v1", func(r chi.Router) {

		r.Get("/users/{id}", func(w http.ResponseWriter, r *http.Request) {
			id := chi.URLParam(r, "id")
			user, err := db.GetUserByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			json.NewEncoder(w).Encode(user)
		})
	})

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
} 