package handlers

import (
	"encoding/json"
	"fmt"
	"hireforwork-server/db"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	service "hireforwork-server/service/modules"
	auth "hireforwork-server/service/modules/auth"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type CompanyHandler struct {
	CompanyService *service.CompanyService
	AuthService    *auth.AuthService
}

func NewCompanyHandler(dbInstance *db.DB) *CompanyHandler {
	return &CompanyHandler{
		CompanyService: service.NewCompanyService(dbInstance),
		AuthService:    auth.NewAuthService(dbInstance),
	}
}
func (h *CompanyHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	fmt.Println(r.URL)
	// Public routes
	publicRoutes := map[string]map[string]http.HandlerFunc{
		"GET": {
			"/companies":              h.GetCompaniesHandler,
			"/companies/random":       h.GetRandomCompanyHandler,
			"/companies/{id}":         h.GetCompanyByID,
			"/companies/get-job/{id}": h.GetJobsByCompany,
		},
		"POST": {
			"/companies/auth/login":           h.LoginCompany,
			"/companies/create":               h.CreateCompany,
			"/request-password-reset-company": h.RequestPasswordCompanyResetHandler,
			"/reset-password-company":         h.ResetPasswordCompanyHandler,
		},
	}

	// Protected routes
	protectedRoutes := map[string]map[string]http.HandlerFunc{
		"GET": {
			"/companies/" + vars["id"] + "/get-applier": h.GetCareerApply,
			"/companies/" + vars["id"] + "/get-static":  h.GetStatics,
		},
		"POST": {
			"/companies/" + vars["id"] + "/update":       h.UpdateCompanyByID,
			"/companies/" + vars["id"] + "/upload-cover": h.UploadCompanyCover,
			"/companies/" + vars["id"] + "/upload-img":   h.UploadCompanyIMG,
			"/companies/change-application-status":       h.ChangeResumeStatusHandler,
		},
		"DELETE": {
			"/companies/" + vars["id"]: h.DeleteCompanyByID,
		},
	}

	// Check public routes
	if handler, ok := publicRoutes[r.Method][r.URL.Path]; ok {
		handler(w, r)
		return
	}

	// Check protected routes
	if handler, ok := protectedRoutes[r.Method][r.URL.Path]; ok {
		handler(w, r)
		return
	}

	http.Error(w, "Not Found", http.StatusNotFound)
}

func (h *CompanyHandler) GetCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	filter := interfaces.ICompanyFilter{
		CompanyName:  r.URL.Query().Get("companyName"),
		CompanyEmail: r.URL.Query().Get("companyEmail"),
		StartDate:    h.getPointer(r.URL.Query().Get("startDate")),
		EndDate:      h.getPointer(r.URL.Query().Get("endDate")),
	}

	companies, err := h.CompanyService.GetCompanies(page, pageSize, filter)
	if err != nil {
		log.Printf("Error getting companies: %v", err)
		http.Error(w, "Failed to get companies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}

func (h *CompanyHandler) GetRandomCompanyHandler(w http.ResponseWriter, r *http.Request) {
	company, err := h.CompanyService.GetRandomCompany()
	if err != nil {
		http.Error(w, "Không có công ty nào", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(company)
}
func (h *CompanyHandler) GetCompanyByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	company, _ := h.CompanyService.GetCompanyByID(vars["id"])
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

func (h *CompanyHandler) CreateCompany(w http.ResponseWriter, r *http.Request) {

	var company models.Company
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &company); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := h.CompanyService.CreateCompany(company)
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

func (h *CompanyHandler) DeleteCompanyByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response := h.CompanyService.DeleteCompanyByID(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
}

func (h *CompanyHandler) UpdateCompanyByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyID := vars["id"]

	var updatedData models.Company
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedCompany, err := h.CompanyService.UpdateCompanyByID(companyID, updatedData)
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

func (h *CompanyHandler) LoginCompany(w http.ResponseWriter, r *http.Request) {
	var credential auth.Credentials

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

func (h *CompanyHandler) GetCareersByJobID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	jobID := vars["id"]
	companyID := vars["companyId"]

	applicants, err := h.CompanyService.GetCareersByJobID(jobID, companyID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(applicants)
}

func (h *CompanyHandler) GetJobsByCompany(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	companyID := vars["id"]
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	filter := interfaces.IJobFilter{
		JobTitle:       r.URL.Query().Get("jobTitle"),
		DateCreateFrom: r.URL.Query().Get("dateCreateFrom"),
		DateCreateTo:   r.URL.Query().Get("dateCreateTo"),
		EndDateFrom:    r.URL.Query().Get("endDateFrom"),
		EndDateTo:      r.URL.Query().Get("endDateTo"),
	}

	jobs, err := h.CompanyService.GetJobsByCompanyID(companyID, page, pageSize, filter)

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

func (h *CompanyHandler) DeleteJobByID(w http.ResponseWriter, r *http.Request) {
	var resBody struct {
		JobIds []string `json:"ids"`
	}
	body, _ := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	err := json.Unmarshal(body, &resBody)
	if err != nil {
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		log.Printf("Error unmarshalling request body: %v", err)
		return
	}
	err = h.CompanyService.DeleteJobByID(resBody.JobIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error deleting job: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Xóa thành công"}`))
}

func (h *CompanyHandler) GetCareerApply(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	filter := interfaces.IJobApplicationFilter{
		Page:        page,
		PageSize:    pageSize,
		CareerEmail: r.URL.Query().Get("careerEmail"),
		JobLevel:    r.URL.Query().Get("jobLevel"),
		JobTitle:    r.URL.Query().Get("jobTitle"),
		Status:      r.URL.Query().Get("status"),
		CreateFrom:  r.URL.Query().Get("createFrom"),
		CreateTo:    r.URL.Query().Get("createTo"),
	}

	res, err := h.CompanyService.GetCareersApplyJob(vars["id"], filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Hàm phụ để chuyển chuỗi thành con trỏ (nếu giá trị không rỗng)
func (h *CompanyHandler) getPointer(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func (h *CompanyHandler) GetStatics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	objID, _ := primitive.ObjectIDFromHex(id)
	res, err := h.CompanyService.GetStatics(objID)
	if err != nil {
		http.Error(w, "Lỗi không xác định", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func (h *CompanyHandler) UploadCompanyCover(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)

	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if _, ok := imageAllowedType[contentType]; !ok {
		http.Error(w, "Chỉ được dùng JPEG, JPG, and PNG.", http.StatusBadRequest)
		return
	}

	url, err := service.UploadImage(file, header, contentType)
	if err != nil {
		http.Error(w, "Lỗi khi upload hình ảnh", http.StatusInternalServerError)
		return
	}

	if err := h.CompanyService.UploadCompanyCover(url, vars["id"]); err != nil {
		http.Error(w, "Lỗi khi cập nhập hình ảnh", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"url": "%s"}`, url)
}

func (h *CompanyHandler) UploadCompanyIMG(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "Unable to parse form", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)

	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if _, ok := imageAllowedType[contentType]; !ok {
		http.Error(w, "Chỉ được dùng JPEG, JPG, and PNG.", http.StatusBadRequest)
		return
	}

	url, err := service.UploadImage(file, header, contentType)
	if err != nil {
		http.Error(w, "Lỗi khi upload hình ảnh", http.StatusInternalServerError)
		return
	}

	if err := h.CompanyService.UploadCompanyImage(url, vars["id"]); err != nil {
		http.Error(w, "Lỗi khi cập nhập hình ảnh", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"url": "%s"}`, url)
}

func (h *CompanyHandler) ChangeResumeStatusHandler(w http.ResponseWriter, r *http.Request) {
	type ChangeStatusRequest struct {
		ResumeID string `json:"_id"`
		Status   string `json:"status"`
	}

	var req ChangeStatusRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Vui lòng thử lại sau", http.StatusBadRequest)
		return
	}
	err = h.CompanyService.ChangeResumeStatus(req.ResumeID, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cập nhập thành công!"))

}
func (h *CompanyHandler) RequestPasswordCompanyResetHandler(w http.ResponseWriter, r *http.Request) {
	var req interfaces.PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	code, err := h.CompanyService.RequestPasswordResetCompany(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"code": code})
}

func (h *CompanyHandler) ResetPasswordCompanyHandler(w http.ResponseWriter, r *http.Request) {
	var req interfaces.PasswordReset
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.CompanyService.ResetPasswordCompany(req.Email, req.Code, req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
