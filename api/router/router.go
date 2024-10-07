package api

import (
	"hireforwork-server/api/handlers"
	"hireforwork-server/middleware"
	"hireforwork-server/service"
	"os"

	"github.com/gorilla/mux"
)

var authService *service.AuthService

func SetUpRouter() *mux.Router {
	router := mux.NewRouter()

	authService = &service.AuthService{JwtSecret: []byte(os.Getenv("SECRET_KEY"))}

	handler := &handlers.Handler{
		AuthService: authService,
	}

	//---USER ROUTER---//
	//PUBLIC ROUTER//
	router.HandleFunc("/careers/auth/login", handler.Login).Methods("POST")
	router.HandleFunc("/careers/create", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/careers", handlers.GetUser).Methods("GET")
	router.HandleFunc("/careers/{id}", handlers.GetUserByID).Methods("GET")
	router.HandleFunc("/careers/{id}", handlers.DeleteUserByID).Methods("DELETE")
	//AUTH ROUTER
	careerRouter := router.PathPrefix("/careers").Subrouter()
	careerRouter.Use(middleware.JWTMiddleware(authService))

	//Post router
	router.HandleFunc("/careers/savejob", handlers.SaveJob).Methods("POST")
	router.HandleFunc("/careers/viewedjob", handlers.CareerViewedJob).Methods("POST")

	//Job Router
	router.HandleFunc("/jobs", handlers.GetJob).Methods("GET")
	router.HandleFunc("/jobs/{id}/apply", handlers.ApplyJob).Methods("POST")

	//Company Router
	router.HandleFunc("/companies/auth/login", handler.LoginCompany).Methods("POST")
	router.HandleFunc("/companies", handlers.GetCompaniesHandler).Methods("GET")
	router.HandleFunc("/companies/{id}", handlers.GetCompanyByID).Methods("GET")
	router.HandleFunc("/companies/create", handlers.CreateCompany).Methods("POST")
	router.HandleFunc("/companies/{id}", handlers.DeleteCompanyByID).Methods("DELETE")
	router.HandleFunc("/companies/update/{id}", handlers.UpdateCompanyByID).Methods("PUT")
	router.HandleFunc("/companies/{companyId}/jobs/{id}", handlers.GetCareersByJobID).Methods("GET")
	router.HandleFunc("/companies/{id}/jobs", handlers.GetJobsByCompany).Methods("GET")

	//Update User Router
	router.HandleFunc("/careers/{id}", handlers.UpdateUser).Methods("PUT")
	return router
}
