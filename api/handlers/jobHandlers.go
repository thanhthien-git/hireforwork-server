package handlers

import (
	"encoding/json"
	"fmt"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"hireforwork-server/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
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

func CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	var job models.Jobs
	err := json.NewDecoder(r.Body).Decode(&job)
	createJob, err := service.CreateJob(job)
	if err != nil {
		http.Error(w, fmt.Sprintf("Lỗi khi tạo mới bài đăng"), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createJob)
	if err != nil {
		http.Error(w, "Có gì đó không ổn", http.StatusInternalServerError)
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

func GetSavedJobs(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	// Gọi service để lấy danh sách công việc đã lưu
	savedJobs, err := service.GetSavedJobsByCareerID(vars["careerID"])
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No saved jobs found for this user", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving saved jobs", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(savedJobs)
}

func GetJobApplyHistoryByCareerID(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)

	jobApplyHistory, err := service.GetJobApplyHistoryByCareerID(vars["id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(jobApplyHistory); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
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
