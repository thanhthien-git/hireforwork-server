package api

import (
	"hireforwork-server/api/handlers"

	"github.com/gorilla/mux"
)

func SetUpRouter() *mux.Router {
	router := mux.NewRouter()

	//User Router
	router.HandleFunc("/careers", handlers.GetUser).Methods("GET")
	router.HandleFunc("/careers/{id}", handlers.GetUserByID).Methods("GET")
	router.HandleFunc("/careers/{id}", handlers.DeleteUserByID).Methods("DELETE")
	router.HandleFunc("/careers/create", handlers.CreateUser).Methods("POST")

	//Company Router
	router.HandleFunc("/companies", handlers.GetCompaniesHandler).Methods("GET")
	return router
}
