package main

import (
	api "hireforwork-server/api/router"
	"hireforwork-server/db"
	"hireforwork-server/service"
	factory "hireforwork-server/service/modules/factory"
	"net/http"

	"github.com/rs/cors"
)

// enableCORS creates a CORS middleware with default settings
func enableCORS() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})
}

// Handler handles all requests
func Handler(w http.ResponseWriter, r *http.Request) {
	// Create database
	database := db.GetInstance()
	defer database.Close()

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
	handler := enableCORS().Handler(routerInstance)

	// Serve the request
	handler.ServeHTTP(w, r)
}

// For local development only
func main() {
	http.HandleFunc("/", Handler)
	port := ":8080"
	println("Server is running on port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		panic(err)
	}
}
