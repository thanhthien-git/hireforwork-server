package handlers

import (
	"encoding/json"
	"hireforwork-server/models"
	service "hireforwork-server/service/modules"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CategoryHandler struct {
	CategoryService *service.CategoryService
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{
		CategoryService: nil,
	}
}

func (h *CategoryHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/category" && r.Method == http.MethodGet:
			h.GetCategory(w, r)
		case r.URL.Path == "/category/create" && r.Method == http.MethodPost:
			h.CreateCategory(w, r)
		case r.URL.Path == "/category":
			if r.Method == http.MethodPut {
				h.UpdateCategoryByID(w, r)
			}
			if r.Method == http.MethodDelete {
				h.DeleteCategoryByID(w, r)
			}
		default:
			h.GetCategory(w, r)
		}
	})

	// Áp dụng decorator nếu có
	// if h.decorator != nil {
	// 	handlerFunc = h.decorator(handlerFunc)
	// }

	handlerFunc.ServeHTTP(w, r)
}

func (c *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)

	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	CategoryName := r.URL.Query().Get("categoryName")
	techList, err := c.CategoryService.GetCategory(page, pageSize, CategoryName)
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

func (c *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {

	var category models.Category
	body, _ := ioutil.ReadAll(r.Body)
	if err := json.Unmarshal(body, &category); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	response, err := c.CategoryService.CreateCategory(category)
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

func (c *CategoryHandler) UpdateCategoryByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID := vars["id"]

	var updatedData models.Category
	if err := json.NewDecoder(r.Body).Decode(&updatedData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedCategory, err := c.CategoryService.UpdateCategoryByID(categoryID, updatedData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(updatedCategory); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (c *CategoryHandler) DeleteCategoryByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	response := c.CategoryService.DeleteCategoryByID(vars["id"])
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(response.StatusCode)
}
