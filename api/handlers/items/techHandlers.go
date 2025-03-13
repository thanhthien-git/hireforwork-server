package handlers

import (
	"encoding/json"
	"hireforwork-server/db"
	"hireforwork-server/models"
	modules "hireforwork-server/service/modules"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type TechHandler struct {
	TechService *modules.TechService
}

func NewTechHandler(dbInstance *db.DB) *TechHandler {
	return &TechHandler{
		TechService: modules.NewTechService(dbInstance),
	}
}

func (h *TechHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		switch r.URL.Path {
		case "/tech":
			h.GetTech(w, r)
		case "/tech/create":
			h.CreateTech(w, r)
		case "/tech/" + vars["id"]:
			if r.Method == http.MethodPost {
				h.UpdateTechByID(w, r)
				return
			}
			if r.Method == http.MethodDelete {
				h.DeleteTechByID(w, r)
				return
			}
		default:
			h.GetTech(w, r)
		}
	})

	// Áp dụng decorator nếu có
	// if h.decorator != nil {
	// 	handlerFunc = h.decorator(handlerFunc)
	// }

	handlerFunc.ServeHTTP(w, r)
}

func (t *TechHandler) GetTech(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)

	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	TechName := r.URL.Query().Get("technology")
	techList, err := t.TechService.GetTech(page, pageSize, TechName)
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

func (t *TechHandler) CreateTech(w http.ResponseWriter, r *http.Request) {

	var technology models.Tech
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &technology); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := t.TechService.CreateTech(technology)
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

func (t *TechHandler) UpdateTechByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	techID := vars["id"]

	var updatedData models.Tech
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedTech, err := t.TechService.UpdateTechByID(techID, updatedData)
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

func (t *TechHandler) DeleteTechByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response := t.TechService.DeleteTechByID(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
}
