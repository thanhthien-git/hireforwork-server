package api

import (
	"hireforwork-server/api/handlers"
	"hireforwork-server/middleware"

	"github.com/gorilla/mux"
)

func setUpCompanyRoutes(router *mux.Router, handler *handlers.Handler) {
	// Public routes
	router.HandleFunc("/companies", handlers.GetCompaniesHandler).Methods("GET")
	router.HandleFunc("/companies/random", handlers.GetRandomCompanyHandler).Methods("GET")
	router.HandleFunc("/companies/auth/login", handler.LoginCompany).Methods("POST")
	router.HandleFunc("/companies/create", handlers.CreateCompany).Methods("POST")
	router.HandleFunc("/companies/{id}", handlers.GetCompanyByID).Methods("GET")
	router.HandleFunc("/companies/get-job/{id}", handlers.GetJobsByCompany).Methods("GET")
	router.HandleFunc("/request-password-reset-company", handlers.RequestPasswordCompanyResetHandler).Methods("POST")
	router.HandleFunc("/reset-password-company", handlers.ResetPasswordCompanyHandler).Methods("POST")

	companies := router.PathPrefix("/companies").Subrouter()
	companies.Use(middleware.JWTMiddleware(handler.AuthService))

	companies.HandleFunc("/{id}/get-applier", handlers.GetCareerApply).Methods("GET")
	companies.HandleFunc("/{id}/get-static", handlers.GetStatics).Methods("GET")
	companies.HandleFunc("/{id}/update", handlers.UpdateCompanyByID).Methods("POST")
	companies.HandleFunc("/{id}/upload-cover", handlers.UploadCompanyCover).Methods("POST")
	companies.HandleFunc("/{id}/upload-img", handlers.UploadCompanyIMG).Methods("POST")
	companies.HandleFunc("/change-application-status", handlers.ChangeResumeStatusHandler).Methods("POST")
	companies.HandleFunc("/{id}", handlers.DeleteCompanyByID).Methods("DELETE")
}
