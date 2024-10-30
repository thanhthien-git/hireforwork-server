package api

import (
	"hireforwork-server/api/handlers"

	"github.com/gorilla/mux"
)

func setUpTechRouter(router *mux.Router, handler *handlers.Handler) {
	// Public Routes
	router.HandleFunc("/tech", handlers.GetTech).Methods("GET")
}
