package api

import (
	"hireforwork-server/api/handlers"

	"github.com/gorilla/mux"
)

func setUpCategoryRouter(router *mux.Router, handler *handlers.Handler) {
	router.HandleFunc("/category", handlers.GetCategory).Methods("GET")
}
