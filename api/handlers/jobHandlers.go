package handlers

import (
	"encoding/json"
	"fmt"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"hireforwork-server/service"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func GetJob(w http.ResponseWriter, r *http.Request) {
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

	jobs, err := service.GetJob(page, pageSize, filter)
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

func UpdateJobHandler(w http.ResponseWriter, r *http.Request) {
	var job models.Jobs
	err := json.NewDecoder(r.Body).Decode(&job)
	updateJob, err := service.UpdateJob(job)
	fmt.Println(updateJob)
	if err != nil {
		http.Error(w, fmt.Sprintln("Có lỗi xảy ra khi cập nhập!"), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(updateJob)
	if err != nil {
		http.Error(w, "Có gì đó không ổn", http.StatusInternalServerError)
	}
}

func ApplyJob(w http.ResponseWriter, r *http.Request) {

	request := interfaces.IJobApply{}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := service.ApplyForJob(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Ứng tuyển thành công, vui lòng kiểm tra email!"})

}

func GetSuggestJobs(w http.ResponseWriter, r *http.Request) {
	jobs, err := service.GetLatestJobs()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetJobByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	authHeader := r.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")

	job, err := service.GetJobByID(vars["id"], tokenString)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(job); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
