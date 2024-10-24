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
	companies.HandleFunc("/get-static/{id}", handlers.GetStatics).Methods("GET")
	companies.HandleFunc("/{id}", handlers.DeleteCompanyByID).Methods("DELETE")
	companies.HandleFunc("/{id}/update", handlers.UpdateCompanyByID).Methods("POST")
	companies.HandleFunc("/{id}/upload-cover", handlers.UploadCompanyCover).Methods("POST")
	companies.HandleFunc("/{id}/upload-img", handlers.UploadCompanyIMG).Methods("POST")
}
