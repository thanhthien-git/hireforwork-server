package api

import (
	"hireforwork-server/api/handlers"
	"hireforwork-server/service"

	"github.com/gorilla/mux"
)

var authService *service.AuthService

func SetUpRouter() *mux.Router {
	router := mux.NewRouter()
	authService = &service.AuthService{}

	handler := &handlers.Handler{
		AuthService: authService,
	}

	//---USER ROUTER---//
	//Get router
	router.HandleFunc("/careers", handlers.GetUser).Methods("GET")
	router.HandleFunc("/careers/{id}", handlers.GetUserByID).Methods("GET")
	//Delete router
	router.HandleFunc("/careers/{id}", handlers.DeleteUserByID).Methods("DELETE")
	//Post router
	router.HandleFunc("/careers/auth/login", handler.Login).Methods("POST")
	router.HandleFunc("/careers/create", handlers.CreateUser).Methods("POST")

	//Job Router
	router.HandleFunc("/jobs", handlers.GetJob).Methods("GET")

	//Company Router
	router.HandleFunc("/companies", handlers.GetCompaniesHandler).Methods("GET")

	//Update User Router
	router.HandleFunc("/careers/{id}", handlers.UpdateUser).Methods("PUT")
	return router
}
