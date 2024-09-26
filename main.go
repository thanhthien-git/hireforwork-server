package main

import (
	api "hireforwork-server/api/router"
	_ "hireforwork-server/service"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	// Create router
	router := api.SetUpRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	// Run server
	log.Printf("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
