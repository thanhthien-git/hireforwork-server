package api

import (
	"hireforwork-server/api/handlers"

	"github.com/gorilla/mux"
)

func setUpJobRouter(router *mux.Router, handler *handlers.Handler) {
	// Public Routes
	router.HandleFunc("/jobs", handlers.GetJob).Methods("GET")
	router.HandleFunc("/jobs/{id}", handlers.GetJobByID).Methods("GET")
	router.HandleFunc("/jobs/create", handlers.CreateJobHandler).Methods("POST")
	router.HandleFunc("/jobs/{id}", handlers.UpdateJobHandler).Methods("POST")
	router.HandleFunc("/jobs", handlers.DeleteJobByID).Methods("DELETE")
}
