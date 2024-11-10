package api

import (
	"hireforwork-server/api/handlers"

	"github.com/gorilla/mux"
)

func setUpCategoryRouter(router *mux.Router, handler *handlers.Handler) {
	router.HandleFunc("/category", handlers.GetCategory).Methods("GET")
	router.HandleFunc("/category/create", handlers.CreateCategory).Methods("POST")
	router.HandleFunc("/category/{id}/update", handlers.UpdateCategoryByID).Methods("PUT")
	router.HandleFunc("/category/{id}", handlers.DeleteCategoryByID).Methods("DELETE")
}
