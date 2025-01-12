package main

import (
	"log"
	"net/http"
	"os"

	"github.com/johnnynu/Coffeehaus/internal/claude"
	"github.com/johnnynu/Coffeehaus/internal/config"
	"github.com/johnnynu/Coffeehaus/internal/database"
	handlers "github.com/johnnynu/Coffeehaus/internal/handlers"
	"github.com/johnnynu/Coffeehaus/internal/maps"
	jwtauth "github.com/johnnynu/Coffeehaus/internal/middleware"

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

	// initialize auth middleware
	authMiddleware := jwtauth.NewAuthMiddleware(dbConfig.SupabaseURL, dbConfig.SupabaseKey)

	if err != nil {
		log.Fatalf("Failed to initialize auth middleware: %v", err)
	}

	// Initialize maps client
	mapsClient, err := maps.NewMapsClient()
	if err != nil {
		log.Fatalf("Failed to initialize maps client: %v", err)
	}

	// Initialize claude service
	claudeService := claude.NewService(os.Getenv("CLAUDE_API_KEY"))

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(db)
	userHandler := handlers.NewUserHandler(db)

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
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	// protected routes
	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Authenticate)

		// User routes
		r.Route("/user", func(r chi.Router) {
			r.Get("/", authHandler.GetUser)
			r.Put("/{username}", userHandler.UpdateProfile)
		})
	})

	log.Printf("Server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
} 