package api

import (
	"hireforwork-server/api/handlers"
	"hireforwork-server/middleware"

	"github.com/gorilla/mux"
)

func setUpCompanyRoutes(router *mux.Router, handler *handlers.Handler) {
	router.HandleFunc("/companies/auth/login", handler.LoginCompany).Methods("POST")
	router.HandleFunc("/companies/{id}", handlers.GetCompanyByID).Methods("GET")
	router.HandleFunc("/companies/create", handlers.CreateCompany).Methods("POST")
	router.HandleFunc("/companies/get-job/{id}", handlers.GetJobsByCompany).Methods("GET")

	companies := router.PathPrefix("/companies").Subrouter()
	companies.Use(middleware.JWTMiddleware(handler.AuthService))
	companies.HandleFunc("", handlers.GetCompaniesHandler).Methods("GET")
	companies.HandleFunc("/get-applier/{id}", handlers.GetCareerApply).Methods("GET")
	companies.HandleFunc("/get-statis/{id}", handlers.GetStatics).Methods("GET")
	companies.HandleFunc("/{id}", handlers.DeleteCompanyByID).Methods("DELETE")
	// router.HandleFunc("/companies/update/{id}", handlers.UpdateCompanyByID).Methods("PUT")
	// router.HandleFunc("/companies/{companyId}/jobs/{jobId}", handlers.DeleteJobByID).Methods("DELETE")
}
