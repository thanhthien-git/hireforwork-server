package handlers

import (
	"encoding/json"
	"fmt"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"hireforwork-server/service"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func GetCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	filter := interfaces.ICompanyFilter{
		CompanyName:  r.URL.Query().Get("companyName"),
		CompanyEmail: r.URL.Query().Get("companyEmail"),
		StartDate:    getPointer(r.URL.Query().Get("startDate")),
		EndDate:      getPointer(r.URL.Query().Get("endDate")),
	}

	companies, err := service.GetCompanies(page, pageSize, filter)
	if err != nil {
		log.Printf("Error getting companies: %v", err)
		http.Error(w, "Failed to get companies", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}

func GetRandomCompanyHandler(w http.ResponseWriter, r *http.Request) {
	company, err := service.GetRandomCompany()
	if err != nil {
		http.Error(w, "Không có công ty nào", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(company)
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
		http.Error(w, err.Error(), http.StatusBadRequest)
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
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

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
	err = service.DeleteJobByID(resBody.JobIds)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Printf("Error deleting job: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"message": "Xóa thành công"}`))
}

func GetCareerApply(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	res, err := service.GetCareersApplyJob(vars["id"])
	if err != nil {
		http.Error(w, "Lỗi xảy ra", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// Hàm phụ để chuyển chuỗi thành con trỏ (nếu giá trị không rỗng)
func getPointer(value string) *string {
	if value == "" {
		return nil
	}
	return &value
}

func GetStatics(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	objID, _ := primitive.ObjectIDFromHex(id)
	res, err := service.GetStatics(objID)
	if err != nil {
		http.Error(w, "Lỗi không xác định", http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(res)
}

func UploadCompanyCover(w http.ResponseWriter, r *http.Request) {
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

	if err := service.UploadCompanyCover(url, vars["id"]); err != nil {
		http.Error(w, "Lỗi khi cập nhập hình ảnh", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"url": "%s"}`, url)
}

func UploadCompanyIMG(w http.ResponseWriter, r *http.Request) {
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

	if err := service.UploadCompanyImage(url, vars["id"]); err != nil {
		http.Error(w, "Lỗi khi cập nhập hình ảnh", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"url": "%s"}`, url)
}

func ChangeResumeStatusHandler(w http.ResponseWriter, r *http.Request) {
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
	err = service.ChangeResumeStatus(req.ResumeID, req.Status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Cập nhập thành công!"))

}
