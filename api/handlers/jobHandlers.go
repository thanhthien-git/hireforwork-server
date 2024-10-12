package handlers

import (
	"encoding/json"
	"hireforwork-server/models"
	"hireforwork-server/service"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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

func ApplyJob(w http.ResponseWriter, r *http.Request) {

	params := mux.Vars(r)
	jobID := params["id"]

	var input struct {
		IDCareer   string `json:"idCareer"`
		IsAccepted string `json:"isAccepted"`
		CreateAt   string `json:"createAt"`
	}

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		log.Printf("Error decoding JSON: %v", err)
		return
	}

	log.Printf("Received ApplyJob input: %+v", input)

	userID, err := primitive.ObjectIDFromHex(input.IDCareer)
	if err != nil {
		http.Error(w, "Invalid career ID", http.StatusBadRequest)
		log.Printf("Invalid career ID: %v", err)
		return
	}

	createAt, err := time.Parse(time.RFC3339, input.CreateAt)
	if err != nil {
		http.Error(w, "Invalid date format", http.StatusBadRequest)
		log.Printf("Invalid date format: %v", err)
		return
	}

	userInfo := models.UserInfo{
		UserId:     userID,
		IsAccepted: input.IsAccepted,
		CreateAt:   primitive.NewDateTimeFromTime(createAt),
	}

	updatedJob, err := service.ApplyForJob(jobID, userInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("Error applying for job: %v", err)
		return
	}

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

	jobApplyHistory, err := service.GetJobApplyHistoryByCareerID(vars["careerID"])
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "No job apply history found for this user", http.StatusNotFound)
		} else {
			http.Error(w, "Error retrieving job apply history", http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jobApplyHistory); err != nil {
		http.Error(w, "Error encoding JSON response", http.StatusInternalServerError)
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
