package main

import (
	api "hireforwork-server/api/router"
	dbHelper "hireforwork-server/db"
	"log"
	"net/http"

	"github.com/rs/cors"
)

func main() {
	client, ctx, err := dbHelper.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

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
	log.Printf("Server is running on port asdasdasdasdasdasd")
	if err := http.ListenAndServe(":8080", handler); err != nil {
		log.Fatal(err)
	}
}
