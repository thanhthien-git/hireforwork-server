package api

import (
	"hireforwork-server/api/handlers"
	"hireforwork-server/middleware"

	"github.com/gorilla/mux"
)

func setUpCareerRoutes(router *mux.Router, handler *handlers.Handler) {
	// Public Routes
	router.HandleFunc("/careers/auth/login", handler.Login).Methods("POST")
	router.HandleFunc("/careers/register", handlers.RegisterCareer).Methods("POST")
	router.HandleFunc("/careers/create", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/{id}/change-password", handler.ChangePassword).Methods("PUT")
	// Protected Routes (with JWT middleware)
	careerRouter := router.PathPrefix("/careers").Subrouter()
	careerRouter.Use(middleware.JWTMiddleware(handler.AuthService))

	careerRouter.HandleFunc("", handlers.GetUser).Methods("GET")
	careerRouter.HandleFunc("/{id}", handlers.GetUserByID).Methods("GET")
	careerRouter.HandleFunc("/{id}", handlers.DeleteUserByID).Methods("DELETE")
	router.HandleFunc("/careers/{id}/upload-image", handlers.UploadImage).Methods("POST")
	careerRouter.HandleFunc("/{id}", handlers.UpdateUser).Methods("PUT")
	// Additional Routes
	router.HandleFunc("/careers/savedjobs/{id}", handlers.GetSavedJobs).Methods("GET")
	router.HandleFunc("/careers/viewedjobs/{id}", handlers.GetViewedJobs).Methods("GET")
	router.HandleFunc("/careers/savejob", handlers.SaveJob).Methods("POST")
	router.HandleFunc("/careers/viewedjob", handlers.CareerViewedJob).Methods("POST")
	router.HandleFunc("/careers/{careerID}/saved-jobs/{jobID}", handlers.RemoveSaveJobHandler).Methods("DELETE")
	router.HandleFunc("/apply-jobs/status/{applyJobID}", handler.ChangeStatus).Methods("PUT")
}


