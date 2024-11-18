package api

import (
	"hireforwork-server/api/handlers"

	"github.com/gorilla/mux"
)

func setUpJobRouter(router *mux.Router, handler *handlers.Handler) {
	router.HandleFunc("/jobs/suggest", handlers.GetSuggestJobs).Methods("GET")
	router.HandleFunc("/jobs/{id}", handlers.GetJobByID).Methods("GET")
	router.HandleFunc("/jobs", handlers.GetJob).Methods("GET")
	router.HandleFunc("/jobs/apply", handlers.ApplyJob).Methods("POST")
	router.HandleFunc("/jobs/create", handlers.CreateJobHandler).Methods("POST")
	router.HandleFunc("/jobs/update", handlers.UpdateJobHandler).Methods("POST")
	router.HandleFunc("/jobs", handlers.DeleteJobByID).Methods("DELETE")
}
