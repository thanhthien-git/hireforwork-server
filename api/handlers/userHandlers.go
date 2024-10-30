package handlers

import (
	"encoding/json"
	"fmt"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"hireforwork-server/service"
	"hireforwork-server/utils"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var imageAllowedType = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/jpg":  true,
}

var resumeAllowFile = map[string]bool{
	"application/pdf": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": true,
}

type Handler struct {
	AuthService *service.AuthService
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)

	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize, _ := strconv.Atoi(pageSizeStr)

	careerFirstName := r.URL.Query().Get("careerFirstName")
	lastName := r.URL.Query().Get("lastName")
	careerEmail := r.URL.Query().Get("careerEmail")
	careerPhone := r.URL.Query().Get("careerPhone")

	users, err := service.GetUser(page, pageSize, careerFirstName, lastName, careerEmail, careerPhone)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, _ := service.GetUserByID(vars["id"])
	response := interfaces.IResponse[models.User]{
		Doc: user,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response := service.DeleteUserByID(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {

	var user models.User
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := service.CreateUser(user)
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

func UploadImage(w http.ResponseWriter, r *http.Request) {
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

	if err := service.UpdateCareerImage(url, vars["id"]); err != nil {
		http.Error(w, "Lỗi khi cập nhập hình ảnh", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"url": "%s"}`, url)
}

func UploadResume(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if _, ok := resumeAllowFile[contentType]; !ok {
		http.Error(w, "Only DOCX, PDF are allowed", http.StatusBadRequest)
		return
	}

	url, err := service.UploadResume(file, header, contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"url": "%s"}`, url)
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var credential service.Credentials

	err := json.NewDecoder(r.Body).Decode(&credential)
	if err != nil {
		http.Error(w, "Invaild request", http.StatusBadRequest)
	}
	if credential.Role == "CAREER" {
		response, err := h.AuthService.LoginForCareer(credential)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(response)
	}
}

func RegisterCareer(w http.ResponseWriter, r *http.Request) {
	type RegisterRequest struct {
		FirstName   string `json:"firstName"`
		LastName    string `json:"lastName"`
		CareerEmail string `json:"careerEmail"`
		CareerPhone string `json:"careerPhone"`
		Password    string `json:"password"`
	}

	var req RegisterRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	existingUser, err := service.GetUserByEmail(req.CareerEmail)
	if err == nil && existingUser.CareerEmail != "" {
		http.Error(w, "Career email already exists", http.StatusConflict)
		return
	}

	hashedPassword := utils.EncodeToSHA(req.Password)

	newUser := models.User{
		Id:            primitive.NewObjectID(),
		CareerEmail:   req.CareerEmail,
		Password:      hashedPassword,
		CreateAt:      primitive.NewDateTimeFromTime(time.Now()),
		IsDeleted:     false,
		Role:          "Career",
		FirstName:     req.FirstName,
		LastName:      req.LastName,
		CareerPhone:   req.CareerPhone,
		CareerPicture: "",
		Languages:     nil,
		Profile:       models.Profile{},
	}

	createdUser, err := service.CreateUser(newUser)
	if err != nil {
		if err.Error() == "duplicate_email" {
			http.Error(w, "Career email already exists", http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("Error creating user: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "User registered successfully",
		"_id":     createdUser.Id.Hex(),
	})
}

// UpdateUser là handler để cập nhật user theo ID
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Lấy ID từ URL
	params := mux.Vars(r)
	id := params["id"]

	log.Printf("User ID: %s", id) // Log ID để kiểm tra
	// Giải mã JSON từ request body thành models.User
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Gọi service để cập nhật user
	updatedUser, err := service.UpdateUserByID(id, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Trả về user đã cập nhật dưới dạng JSON
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}

func SaveJob(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		CareerID string `json:"careerID"`
		JobID    string `json:"jobID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	savedJob, err := service.SaveJob(payload.CareerID, payload.JobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(savedJob)
}

func CareerViewedJob(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		CareerID string `json:"careerID"`
		JobID    string `json:"jobID"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Debugging payload
	fmt.Printf("Received CareerID: %s, JobID: %s\n", payload.CareerID, payload.JobID)

	viewedJob, err := service.CareerViewedJob(payload.CareerID, payload.JobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(viewedJob)
}

func RemoveSaveJobHandler(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	careerID := vars["careerID"]
	jobID := vars["jobID"]

	updatedCareerSaveJob, err := service.RemoveSaveJob(careerID, jobID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedCareerSaveJob)
}

func GetSavedJobs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	savedJobs, err := service.GetSavedJobByCareerID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(savedJobs)
}

func GetViewedJobs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	viewedJobs, err := service.GetViewedJobByCareerID(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(viewedJobs)
}
