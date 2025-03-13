package main

import (
	api "hireforwork-server/api/router"
	"hireforwork-server/db"
	"hireforwork-server/service"
	factory "hireforwork-server/service/modules/factory"
	"log"
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

// Handler is the exported function that Vercel will use to handle requests
func Handler(w http.ResponseWriter, r *http.Request) {
	//create database
	database := db.GetInstance()
	defer database.Close()

	//create service container
	container := service.NewServiceContainer()

	//create service dependencies
	deps := &factory.ServiceDependencies{
		DB: database,
	}

	//create service factory with dependencies
	serviceFactory := factory.NewServiceFactory(deps)

	// Register all services at once using factory
	serviceFactory.RegisterAllServices(container)

	// Create app services
	appServices := service.NewAppServices(container)

	// Create router
	router := api.SetUpRouter(appServices, database)

	// Enable CORS with a single line
	handler := enableCORS().Handler(router)

	// Serve the request
	handler.ServeHTTP(w, r)
}

func main() {
	// For local development
	http.HandleFunc("/", Handler)
	log.Printf("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}
