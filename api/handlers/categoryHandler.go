package handlers

import (
	"encoding/json"
	"hireforwork-server/service"
	"net/http"
	"strconv"
)

func GetCategory(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)

	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize, _ := strconv.Atoi(pageSizeStr)

	categoryName := r.URL.Query().Get("categoryName") 

	categoryList, err := service.GetCategory(page, pageSize, categoryName)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(categoryList); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}
