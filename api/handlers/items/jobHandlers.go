package handlers

import (
	"encoding/json"
	"fmt"
	"hireforwork-server/db"
	"hireforwork-server/interfaces"
	"hireforwork-server/middleware"
	"hireforwork-server/models"
	"hireforwork-server/service/modules/jobs"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

type JobHandler struct {
	JobService *jobs.JobService
}

func NewJobHandler(dbInstance *db.DB) *JobHandler {
	return &JobHandler{
		JobService: jobs.NewJobService(dbInstance),
	}
}

func (h *JobHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		vars := mux.Vars(r)
		// Remove trailing slash if present
		if len(path) > 1 && path[len(path)-1] == '/' {
			path = path[:len(path)-1]
		}
		// Handle the base path
		if path == "/jobs" || path == "/jobs/" || path == "/" || path == "" || path == "/api/jobs" || path == "/api/jobs/" {
			switch r.Method {
			case http.MethodGet:
				h.GetJob(w, r)
			case http.MethodPost:
				h.CreateJobHandler(w, r)
			case http.MethodPut:
				h.UpdateJobHandler(w, r)
			default:
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
			return
		}

		// Handle other paths
		switch {
		case strings.HasSuffix(path, "save"):
			if r.Method == http.MethodPost {
				h.SaveJob(w, r)
				return
			}
		case strings.HasSuffix(path, "apply"):
			if r.Method == http.MethodPost {
				h.ApplyJob(w, r)
				return
			}
		case strings.HasSuffix(path, "suggest"):
			if r.Method == http.MethodGet {
				h.GetSuggestJobs(w, r)
			} else {
				http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			}
		default:
			if _, ok := vars["id"]; ok && r.Method == http.MethodGet {
				h.GetJobByID(w, r)
			} else {
				http.Error(w, "Not Found", http.StatusNotFound)
			}
		}
	})
	handlerFunc.ServeHTTP(w, r)
}

func (h *JobHandler) GetJob(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)

	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize, _ := strconv.Atoi(pageSizeStr)

	isHotStr := r.URL.Query().Get("isHot")
	isHot := false
	if isHotStr == "true" || isHotStr == "1" {
		isHot = true
	}

	salaryFromStr := r.URL.Query().Get("salaryFrom")
	salaryFrom, _ := strconv.ParseInt(salaryFromStr, 10, 64)

	salaryToStr := r.URL.Query().Get("salaryTo")
	salaryTo, _ := strconv.ParseInt(salaryToStr, 10, 64)

	filter := interfaces.IJobFilter{
		JobTitle:        r.URL.Query().Get("jobTitle"),
		CompanyName:     r.URL.Query().Get("companyName"),
		DateCreateFrom:  r.URL.Query().Get("dateCreateFrom"),
		DateCreateTo:    r.URL.Query().Get("dateCreateTo"),
		EndDateFrom:     r.URL.Query().Get("endDateFrom"),
		EndDateTo:       r.URL.Query().Get("endDateTo"),
		SalaryFrom:      salaryFrom,
		SalaryTo:        salaryTo,
		WorkingLocation: r.URL.Query()["workingLocation"],
		JobRequirement:  r.URL.Query()["jobRequirement"],
		JobCategory:     r.URL.Query()["jobCategory"],
		JobLevel:        r.URL.Query().Get("jobLevel"),
		IsHot:           isHot,
		Query:           r.URL.Query().Get("query"),
	}

	jobs, err := h.JobService.GetJob(page, pageSize, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func (h *JobHandler) CreateJobHandler(w http.ResponseWriter, r *http.Request) {
	var job models.Jobs
	err := json.NewDecoder(r.Body).Decode(&job)
	if err != nil {
		http.Error(w, fmt.Sprintf("Lỗi khi tạo mới bài đăng"), http.StatusInternalServerError)
		return
	}
	createJob, err := h.JobService.CreateJob(job)
	if err != nil {
		http.Error(w, fmt.Sprintf("Lỗi khi tạo mới bài đăng"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createJob)
}

func (h *JobHandler) UpdateJobHandler(w http.ResponseWriter, r *http.Request) {
	var job models.Jobs
	err := json.NewDecoder(r.Body).Decode(&job)
	updateJob, err := h.JobService.UpdateJob(job)
	fmt.Println(updateJob)
	if err != nil {
		http.Error(w, fmt.Sprintln("Có lỗi xảy ra khi cập nhập!"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(updateJob)
	if err != nil {
		http.Error(w, "Có gì đó không ổn", http.StatusInternalServerError)
	}
}

func (h *JobHandler) SaveJob(w http.ResponseWriter, r *http.Request) {
	careerId := middleware.GetUserID(r)
	jobId := mux.Vars(r)["id"]

	savedJob, err := h.JobService.SaveJob(careerId, jobId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(savedJob)
}

func (h *JobHandler) ApplyJob(w http.ResponseWriter, r *http.Request) {

	userID := middleware.GetUserID(r)
	var request interfaces.IJobApply

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request.IDCareer = userID
	request.JobID = mux.Vars(r)["id"]

	err := h.JobService.Apply(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"message": "Ứng tuyển thành công, vui lòng kiểm tra email!"})

}

func (h *JobHandler) GetSuggestJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := h.JobService.GetLatestJobs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *JobHandler) GetJobByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var userID string
	if r.Header.Get("Authorization") != "" {
		userID = middleware.GetUserID(r)
	}
	job, err := h.JobService.GetJobByID(vars["id"], userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
