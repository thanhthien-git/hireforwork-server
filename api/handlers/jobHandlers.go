package handlers

import (
	"encoding/json"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"hireforwork-server/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetJob(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)

	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize, _ := strconv.Atoi(pageSizeStr)

	jobs, err := service.GetJob(page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func ApplyJob(w http.ResponseWriter, r *http.Request) {

	request := interfaces.IJobApply{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding JSON: %v", err)
		return
	}

	updatedJob, _ := service.ApplyForJob(request)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedJob); err != nil {
		http.Error(w, "Error encoding response JSON", http.StatusInternalServerError)
	}
}

type JobHandler struct {
	JobService *service.JobService
}

func (h *JobHandler) GetSuggestJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.JobService.GetLatestJobs()
	if err != nil {
		http.Error(w, "Error fetching jobs", http.StatusInternalServerError)
		log.Printf("Error fetching jobs: %v", err) // Ghi lại lỗi chi tiết
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetJobByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	job, err := service.GetJobByID(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	response := interfaces.IResponse[models.Jobs]{
		Doc: job,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

