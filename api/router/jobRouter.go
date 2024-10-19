package api

import (
	"hireforwork-server/api/handlers"
	"hireforwork-server/service"

	"github.com/gorilla/mux"
)

var jobHandler *handlers.JobHandler
var jobService *service.JobService

func setUpJobRouter(router *mux.Router, handler *handlers.Handler) {
	// Public Routes
	jobService = &service.JobService{Collection: service.JobCollection}
	jobHandler = &handlers.JobHandler{JobService: jobService}
	router.HandleFunc("/jobs", handlers.GetJob).Methods("GET")
	router.HandleFunc("/jobs/{id}", handlers.GetJobByID).Methods("GET")
	router.HandleFunc("/FilteredJobs/filter", jobHandler.GetFilteredJobs).Methods("GET")
	router.HandleFunc("/suggest", jobHandler.GetSuggestJobs).Methods("GET")
	router.HandleFunc("/jobs/create", handlers.CreateJobHandler).Methods("POST")
}
