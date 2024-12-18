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

	// Protected Routes (with JWT middleware)
	careerRouter := router.PathPrefix("/careers").Subrouter()
	careerRouter.Use(middleware.JWTMiddleware(handler.AuthService))
	careerRouter.HandleFunc("", handlers.GetUser).Methods("GET")
	careerRouter.HandleFunc("/{id}", handlers.GetUserByID).Methods("GET")
	careerRouter.HandleFunc("/{id}/save-job", handlers.GetSavedJobs).Methods("GET")
	careerRouter.HandleFunc("/{id}/applied-job", handlers.GetAppliedJob).Methods("GET")
	careerRouter.HandleFunc("/{id}/save", handlers.SaveJob).Methods("POST")
	careerRouter.HandleFunc("/{id}/remove-save", handlers.RemoveSaveJob).Methods("POST")
	careerRouter.HandleFunc("/{id}", handlers.DeleteUserByID).Methods("DELETE")
	careerRouter.HandleFunc("/{id}/upload-image", handlers.UploadImage).Methods("POST")
	careerRouter.HandleFunc("/{id}/upload-resume", handlers.UploadResume).Methods("POST")
	careerRouter.HandleFunc("/{id}/remove-resume", handlers.RemoveResume).Methods("POST")
	careerRouter.HandleFunc("/{id}/update", handlers.UpdateUser).Methods("POST")

	// Additional Routes
	router.HandleFunc("/careers/viewedjobs/{id}", handlers.GetViewedJobs).Methods("GET")
	router.HandleFunc("/careers/viewedjob", handlers.CareerViewedJob).Methods("POST")
	router.HandleFunc("/request-password-reset", handlers.RequestPasswordResetHandler).Methods("POST")
	router.HandleFunc("/reset-password", handlers.ResetPasswordHandler).Methods("POST")
}
