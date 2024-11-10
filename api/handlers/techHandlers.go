package handlers

import (
	"encoding/json"
	"fmt"
	"hireforwork-server/interfaces"
	"hireforwork-server/models"
	"hireforwork-server/service"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func GetTech(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)

	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize, _ := strconv.Atoi(pageSizeStr)

	filter := interfaces.ITechnologyFilter{
		Technology: r.URL.Query().Get("technology")}
	techList, err := service.GetTech(page, pageSize, filter)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(techList); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}
func CreateTech(w http.ResponseWriter, r *http.Request) {
	var technology models.Tech
	if err := json.NewDecoder(r.Body).Decode(&technology); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	createTechnology, err := service.CreateTech(technology)
	if err != nil {
		http.Error(w, fmt.Sprintf("Lỗi khi tạo kĩ năng"), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(createTechnology)
	if err != nil {
		http.Error(w, "Kĩ năng có gì đó không ổn", http.StatusInternalServerError)
	}
}
func DeleteTechByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response := service.DeleteTechByID(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
}

func UpdateTechByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	techID := vars["id"]

	var technology models.Tech
	if err := json.NewDecoder(r.Body).Decode(&technology); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTech, err := service.UpdateTechByID(techID, technology)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedTech); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func GetTechByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	technology, _ := service.GetTechByID(vars["id"])
	response := interfaces.IResponse[models.Tech]{
		Doc: technology,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
