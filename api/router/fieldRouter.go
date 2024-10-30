package api

import (
	"hireforwork-server/api/handlers"

	"github.com/gorilla/mux"
)

func setUpFieldRouter(router *mux.Router, handler *handlers.Handler) {
	// Public Routes
	router.HandleFunc("/company-field", handlers.GetField).Methods("GET")
}
