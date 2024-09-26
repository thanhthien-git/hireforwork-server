package handlers

import (
	"encoding/json"
	"hireforwork-server/service"
	"log"
	"net/http"
	"strconv"
)

// Lấy danh sách công ty với phân trang
func GetCompaniesHandler(w http.ResponseWriter, r *http.Request) {
	// Lấy các tham số phân trang từ query string
	pageStr := r.URL.Query().Get("page")
	pageSizeStr := r.URL.Query().Get("pageSize")

	// Chuyển đổi từ chuỗi sang số
	page, _ := strconv.Atoi(pageStr)
	pageSize, _ := strconv.Atoi(pageSizeStr)

	// Gọi service để lấy danh sách công ty
	companies, err := service.GetCompanies(page, pageSize)
	if err != nil {
		log.Printf("Error getting companies: %v", err)
		http.Error(w, "Failed to get companies", http.StatusInternalServerError)
		return
	}

	// Trả về danh sách công ty dưới dạng JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(companies)
}
