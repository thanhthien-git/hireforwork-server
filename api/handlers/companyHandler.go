package handlers

import (
	"encoding/json"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"hireforwork-server/service"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	companyName := r.URL.Query().Get("companyName")
	companyEmail := r.URL.Query().Get("companyEmail")

	companies, err := service.GetCompanies(page, pageSize, companyName, companyEmail, nil, nil)
	if err != nil {
		log.Printf("Error getting companies: %v", err)
		http.Error(w, "Failed to get companies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}

func GetCompanyByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	company, _ := service.GetCompanyByID(vars["id"])
	response := interfaces.IResponse[models.Company]{
		Doc: company,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func CreateCompany(w http.ResponseWriter, r *http.Request) {

	var company models.Company
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &company); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := service.CreateCompany(company)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteCompanyByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response := service.DeleteCompanyByID(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
}

func UpdateCompanyByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyID := vars["id"]

	var updatedData models.Company
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	updatedCompany, err := service.UpdateCompanyByID(companyID, updatedData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedCompany); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) LoginCompany(w http.ResponseWriter, r *http.Request) {
	var credential service.Credentials

	err := json.NewDecoder(r.Body).Decode(&credential)

	if err != nil {
		http.Error(w, "Invaild request", http.StatusBadRequest)
	}

	if credential.Role == "COMPANY" {
		response, err := h.AuthService.LoginForCompany(credential)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}

func GetCareersByJobID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["id"]
	companyID := vars["companyId"]

	applicants, err := service.GetCareersByJobID(jobID, companyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applicants)
}

func GetJobsByCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyID := vars["id"]

	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	jobs, err := service.GetJobsByCompanyID(companyID, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(jobs); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteJobByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyID := vars["companyId"]
	jobID := vars["jobId"]

	err := service.DeleteJobByID(companyID, jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error deleting job: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Job deleted successfully"}`))
}

func FilterCompaniesByDateHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	companyName := r.URL.Query().Get("companyName")
	companyEmail := r.URL.Query().Get("companyEmail")

	startDateStr := r.URL.Query().Get("startDate")
	endDateStr := r.URL.Query().Get("endDate")

	var endDate *string

	if endDateStr != "" {
		endDate = &endDateStr
	}

	companies, _ := service.GetCompanies(page, pageSize, companyName, companyEmail, &startDateStr, endDate)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}
