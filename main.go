package main

import (
	"log"
	"net/http"
	"os"
	"sports-backend/database"
	"sports-backend/handlers"
	"sports-backend/repositories"
	"sports-backend/services"

	"github.com/gorilla/mux"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer database.CloseDB()

	// Run migrations
	if err := database.RunMigrations(); err != nil {
		log.Fatal("Failed to run migrations:", err)
	}

	// Initialize repositories
	teamRepo := repositories.NewTeamRepository(database.DB)
	playerRepo := repositories.NewPlayerRepository(database.DB)

	// Initialize services
	teamService := services.NewTeamService(teamRepo)
	playerService := services.NewPlayerService(playerRepo, teamRepo)

	// Initialize handlers
	teamHandler := handlers.NewTeamHandler(teamService)
	playerHandler := handlers.NewPlayerHandler(playerService)

	// Create router
	router := mux.NewRouter()

	// Add CORS middleware
	router.Use(corsMiddleware)

	// API routes
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Teams routes
	apiRouter.HandleFunc("/teams", teamHandler.GetTeams).Methods("GET")
	apiRouter.HandleFunc("/teams", teamHandler.CreateTeam).Methods("POST")
	apiRouter.HandleFunc("/teams/{id}", teamHandler.GetTeam).Methods("GET")
	apiRouter.HandleFunc("/teams/{id}", teamHandler.UpdateTeam).Methods("PUT")
	apiRouter.HandleFunc("/teams/{id}", teamHandler.DeleteTeam).Methods("DELETE")
	apiRouter.HandleFunc("/teams/{id}/stats", teamHandler.GetTeamStats).Methods("GET")
	apiRouter.HandleFunc("/teams/{id}/stats", teamHandler.CreateTeamStats).Methods("POST")

	// Players routes
	apiRouter.HandleFunc("/players", playerHandler.GetPlayers).Methods("GET")
	apiRouter.HandleFunc("/players", playerHandler.CreatePlayer).Methods("POST")
	apiRouter.HandleFunc("/players/{id}", playerHandler.GetPlayer).Methods("GET")
	apiRouter.HandleFunc("/players/{id}", playerHandler.UpdatePlayer).Methods("PUT")
	apiRouter.HandleFunc("/players/{id}", playerHandler.DeletePlayer).Methods("DELETE")
	apiRouter.HandleFunc("/players/{id}/stats", playerHandler.GetPlayerStats).Methods("GET")
	apiRouter.HandleFunc("/players/{id}/stats", playerHandler.CreatePlayerStats).Methods("POST")
	apiRouter.HandleFunc("/players/{id}/stats/{stats_id}", playerHandler.UpdatePlayerStats).Methods("PUT")
	apiRouter.HandleFunc("/players/{id}/stats/{stats_id}", playerHandler.DeletePlayerStats).Methods("DELETE")

	// Health check endpoint
	router.HandleFunc("/health", func(responseWriter http.ResponseWriter, request *http.Request) {
		responseWriter.Header().Set("Content-Type", "application/json")
		responseWriter.WriteHeader(http.StatusOK)
		responseWriter.Write([]byte(`{"status": "healthy"}`))
	}).Methods("GET")

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Printf("API endpoints available at http://localhost:%s/api", port)
	log.Printf("Health check available at http://localhost:%s/health", port)

	if serverError := http.ListenAndServe(":"+port, router); serverError != nil {
		log.Fatal("Server failed to start:", serverError)
	}
}

// corsMiddleware adds CORS headers to allow frontend connections
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
