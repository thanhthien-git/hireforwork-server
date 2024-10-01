package api

import (
	"hireforwork-server/api/handlers"

	"github.com/gorilla/mux"
)

func SetUpRouter() *mux.Router {
	router := mux.NewRouter()

	//---USER ROUTER---//
	//Get router
	router.HandleFunc("/careers", handlers.GetUser).Methods("GET")
	router.HandleFunc("/careers/{id}", handlers.GetUserByID).Methods("GET")
	//Delete router
	router.HandleFunc("/careers/{id}", handlers.DeleteUserByID).Methods("DELETE")
	//Post router
	router.HandleFunc("/careers/create", handlers.CreateUser).Methods("POST")

	//Company Router
	router.HandleFunc("/companies", handlers.GetCompaniesHandler).Methods("GET")
	router.HandleFunc("/companies/{id}", handlers.GetCompanyByID).Methods("GET")
	router.HandleFunc("/companies/create", handlers.CreateCompany).Methods("POST")
	router.HandleFunc("/companies/{id}", handlers.DeleteCompanyByID).Methods("DELETE")
	router.HandleFunc("/companies/update/{id}", handlers.UpdateCompanyByID).Methods("PUT")
	return router
}
