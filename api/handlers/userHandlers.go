package handlers

import (
	"encoding/json"
	"fmt"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"hireforwork-server/service"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
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

	users, err := service.GetUser(page, pageSize)
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
	r.ParseMultipartForm(10 << 20)

	file, header, err := r.FormFile("avatar")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if _, ok := imageAllowedType[contentType]; !ok {
		http.Error(w, "Only JPEG, JPG, and PNG are allowed.", http.StatusBadRequest)
		return
	}

	url, err := service.UploadImage(file, header, contentType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
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
	if credential.Role == "career" {
		token, err := h.AuthService.LoginForCareer(credential)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]string{"token": token})
	}
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
