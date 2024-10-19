package handlers

import (
	"encoding/json"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"hireforwork-server/service"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
func (h *JobHandler) GetFilteredJobs(w http.ResponseWriter, r *http.Request) {
	createDateStr := r.URL.Query().Get("createAt")
	expireDateStr := r.URL.Query().Get("expireDate")

	if createDateStr == "" {
		http.Error(w, "createAt is required", http.StatusBadRequest)
		return
	}
	if expireDateStr == "" {
		http.Error(w, "expireDate is required", http.StatusBadRequest)
		return
	}

	// Chuyển đổi chuỗi thành time.Time
	createDate, err := time.Parse("2006-01-02", createDateStr)
	createDate = createDate.UTC()
	if err != nil {
		http.Error(w, "Invalid createAt date", http.StatusBadRequest)
		return
	}

	expireDate, err := time.Parse("2006-01-02", expireDateStr)
	expireDate = expireDate.Add(24 * time.Hour).Add(-time.Nanosecond)
	if err != nil {
		http.Error(w, "Invalid expireDate date", http.StatusBadRequest)
		return
	}

	// Chuyển đổi sang primitive.DateTime
	createDatePrimitive := primitive.NewDateTimeFromTime(createDate)
	expireDatePrimitive := primitive.NewDateTimeFromTime(expireDate)

	jobs, err := h.JobService.GetFilteredJobs(r.Context(), createDatePrimitive, expireDatePrimitive)
	if err != nil {
		http.Error(w, "Error fetching jobs", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(jobs)
}

