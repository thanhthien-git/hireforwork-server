package main

import (
	api "hireforwork-server/api/router"
	dbHelper "hireforwork-server/db"
	"log"
	"net/http"
)

func main() {
	client, ctx, err := dbHelper.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Create router
	router := api.SetUpRouter()

	// Run server
	log.Printf("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
