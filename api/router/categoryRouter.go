package api

import (
	"hireforwork-server/api/handlers"

	"github.com/gorilla/mux"
)

func setUpCategoryRouter(router *mux.Router, handler *handlers.Handler) {
	router.HandleFunc("/category", handlers.GetCategory).Methods("GET")
	router.HandleFunc("/category/{id}", handlers.GetCategoryByID).Methods("GET")
	router.HandleFunc("/category", handlers.CreateCategory).Methods("POST")
	router.HandleFunc("/category/{id}", handlers.DeleteCategoryByID).Methods("DELETE")
	router.HandleFunc("/category/update/{id}", handlers.UpdateCategoryByID).Methods("POST")
}
