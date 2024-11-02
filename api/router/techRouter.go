package api

import (
	"hireforwork-server/api/handlers"

	"github.com/gorilla/mux"
)

func setUpTechRouter(router *mux.Router, handler *handlers.Handler) {
	router.HandleFunc("/tech", handlers.GetTech).Methods("GET")
}
