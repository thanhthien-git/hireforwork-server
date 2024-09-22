package handlers

import (
	"encoding/json"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"hireforwork-server/service"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

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
