package handlers

import (
	"encoding/json"
	"fmt"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	service "hireforwork-server/service/modules"
	"hireforwork-server/service/modules/auth"
	"hireforwork-server/utils"
	"io/ioutil"
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

type UserHandler struct {
	UserService *service.UserService
	AuthService *auth.AuthService
}

func (h *UserHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Public routes
		switch r.URL.Path {
		case "/careers/auth/login":
			if r.Method == http.MethodPost {
				h.Login(w, r)
				return
			}
		case "/careers/register":
			if r.Method == http.MethodPost {
				h.RegisterCareer(w, r)
				return
			}
		case "/careers/create":
			if r.Method == http.MethodPost {
				h.CreateUser(w, r)
				return
			}
		case "/request-password-reset":
			if r.Method == http.MethodPost {
				h.RequestPasswordResetHandler(w, r)
				return
			}
		case "/reset-password":
			if r.Method == http.MethodPost {
				h.ResetPasswordHandler(w, r)
				return
			}
		}

		// Protected routes (middleware JWTMiddleware sẽ được áp dụng trong router)
		vars := mux.Vars(r)
		switch r.URL.Path {
		case "/careers":
			if r.Method == http.MethodGet {
				h.GetUser(w, r)
				return
			}
		case "/careers/" + vars["id"]:
			if r.Method == http.MethodGet {
				h.GetUserByID(w, r)
				return
			}
			if r.Method == http.MethodDelete {
				h.DeleteUserByID(w, r)
				return
			}
		case "/careers/" + vars["id"] + "/save-job":
			if r.Method == http.MethodGet {
				h.GetSavedJobs(w, r)
				return
			}
		case "/careers/" + vars["id"] + "/applied-job":
			if r.Method == http.MethodGet {
				h.GetAppliedJob(w, r)
				return
			}
		// case "/careers/" + vars["id"] + "/save":
		// 	if r.Method == http.MethodPost {
		// 		h.SaveJob(w, r)
		// 		return
		// 	}
		// case "/careers/" + vars["id"] + "/remove-save":
		// 	if r.Method == http.MethodPost {
		// 		h.RemoveSaveJob(w, r)
		// 		return
		// 	}
		case "/careers/" + vars["id"] + "/upload-image":
			if r.Method == http.MethodPost {
				h.UploadImage(w, r)
				return
			}
		case "/careers/" + vars["id"] + "/upload-resume":
			if r.Method == http.MethodPost {
				h.UploadResume(w, r)
				return
			}
		case "/careers/" + vars["id"] + "/remove-resume":
			if r.Method == http.MethodPost {
				h.RemoveResume(w, r)
				return
			}
		case "/careers/" + vars["id"] + "/update":
			if r.Method == http.MethodPost {
				h.UpdateUser(w, r)
				return
			}
		}

		http.Error(w, "Not Found", http.StatusNotFound)
	})

	// // Áp dụng decorator nếu có
	// if h.decorator != nil {
	// 	handlerFunc = h.decorator(handlerFunc)
	// }
	handlerFunc.ServeHTTP(w, r)
}

func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)

	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize, _ := strconv.Atoi(pageSizeStr)

	careerFirstName := r.URL.Query().Get("careerFirstName")
	lastName := r.URL.Query().Get("lastName")
	careerEmail := r.URL.Query().Get("careerEmail")
	careerPhone := r.URL.Query().Get("careerPhone")

	users, err := h.UserService.GetUser(page, pageSize, careerFirstName, lastName, careerEmail, careerPhone)
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

func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	user, _ := h.UserService.GetUserByID(vars["id"])
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

func (h *UserHandler) DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response := h.UserService.DeleteUserByID(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {

	var user models.User
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := h.UserService.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(map[string]string{"message": "User created successfully"}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *UserHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
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

	if err := h.UserService.UpdateCareerImage(url, vars["id"]); err != nil {
		http.Error(w, "Lỗi khi cập nhập hình ảnh", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"url": "%s"}`, url)
}

func (h *UserHandler) UploadResume(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("resume")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()
	vars := mux.Vars(r)

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
	if err := h.UserService.UpdateCareerResume(url, vars["id"]); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, `{"url": "%s"}`, url)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credential auth.Credentials

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

func (h *UserHandler) RegisterCareer(w http.ResponseWriter, r *http.Request) {
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

	existingUser, err := h.UserService.GetUserByEmail(req.CareerEmail)
	if err == nil && existingUser.CareerEmail != "" {
		http.Error(w, "Career email already exists", http.StatusConflict)
		return
	}

	hashedPassword := utils.EncodeToSHA(req.Password)

	newUser := models.User{
		Id:               primitive.NewObjectID(),
		CareerEmail:      req.CareerEmail,
		Password:         hashedPassword,
		CreateAt:         primitive.NewDateTimeFromTime(time.Now()),
		IsDeleted:        false,
		Role:             "Career",
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		CareerPhone:      req.CareerPhone,
		CareerPicture:    "",
		Languages:        nil,
		Profile:          models.Profile{},
		VerificationCode: "",
	}

	err = h.UserService.CreateUser(newUser)
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
	})
}

func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	updatedUser, err := h.UserService.UpdateUserByID(id, user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(updatedUser)
}

func (h *UserHandler) GetSavedJobs(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	savedJobs := h.UserService.GetSavedJobByCareerID(id)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(savedJobs)
}

func (h *UserHandler) RemoveResume(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	var data map[string]interface{}
	json.NewDecoder(r.Body).Decode(&data)

	resumeURL := data["resumeURL"].(string)
	err := h.UserService.RemoveResume(id, resumeURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}
func (h *UserHandler) RequestPasswordResetHandler(w http.ResponseWriter, r *http.Request) {
	var req interfaces.PasswordResetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	code, err := h.UserService.RequestPasswordReset(req.Email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(map[string]string{"code": code})
}

func (h *UserHandler) ResetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	var req interfaces.PasswordReset
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	if err := h.UserService.ResetPassword(req.Email, req.Code, req.NewPassword); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// func (h *UserHandler) GetStaticHandler(w http.ResponseWriter, r *http.Request) {
// 	staticData := h.UserService.GetStatic()
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	if err := json.NewEncoder(w).Encode(staticData); err != nil {
// 		// Handle error if encoding fails
// 		w.WriteHeader(http.StatusInternalServerError)
// 		log.Println("Error encoding response:", err)
// 		return
// 	}
// }

func (h *UserHandler) GetAppliedJob(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	query := r.URL.Query()
	page, err := strconv.Atoi(query.Get("page"))
	pageSize, err := strconv.Atoi(query.Get("pageSize"))

	result, err := h.UserService.GetAppliedJob(id, page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
