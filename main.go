package main

import (
	api "hireforwork-server/api/router"
	_ "hireforwork-server/service"
	"log"
	"net/http"
)

func main() {
	// Create router
	router := api.SetUpRouter()

	// Run server
	log.Printf("Server is running on port 8080")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatal(err)
	}
}
