package api

import (
	"hireforwork-server/api/handlers"

	"github.com/gorilla/mux"
)

func setUpTechRouter(router *mux.Router, handler *handlers.Handler) {
	router.HandleFunc("/tech", handlers.GetTech).Methods("GET")
	router.HandleFunc("/tech/{id}", handlers.GetTechByID).Methods("GET")
	router.HandleFunc("/tech/create", handlers.CreateTech).Methods("POST")
	router.HandleFunc("/tech/{id}", handlers.DeleteTechByID).Methods("DELETE")
	router.HandleFunc("/tech/{id}/update", handlers.UpdateTechByID).Methods("POST")
}
