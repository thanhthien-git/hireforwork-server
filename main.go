package handler

import (
	"context"
	api "hireforwork-server/api/router"
	"hireforwork-server/db"
	"hireforwork-server/service"
	factory "hireforwork-server/service/modules/factory"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/rs/cors"
)

// enableCORS creates a CORS middleware with default settings
func enableCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://your-production-domain.com"}, // Replace with your actual domains
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           43200, // 12 hours in seconds
	})
}

// Handler handles all requests
func Handler(handler http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create context with timeout
		ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
		defer cancel()

		// Create a new request with context
		newReq := r.WithContext(ctx)

		// Serve the request
		handler.ServeHTTP(w, newReq)
	}
}

// For local development only
func main() {
	// Create database
	database := db.GetInstance()

	// Create service container
	container := service.NewServiceContainer()

	// Create service dependencies
	deps := &factory.ServiceDependencies{
		DB: database,
	}

	// Create service factory with dependencies
	serviceFactory := factory.NewServiceFactory(deps)

	// Register all services at once using factory
	serviceFactory.RegisterAllServices(container)

	// Create app services
	appServices := service.NewAppServices(container)

	// Create router
	routerInstance := api.SetUpRouter(appServices, database)

	// Enable CORS
	corsHandler := enableCORS().Handler(routerInstance)

	// Create the final handler with context
	finalHandler := Handler(corsHandler)

	port := ":8080"
	log.Printf("Server is running on port %s", port)
	if err := http.ListenAndServe(port, finalHandler); err != nil {
		log.Printf("Server error: %v", err)
		os.Exit(1)
	}
}
