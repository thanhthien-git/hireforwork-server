package handlers

import (
	"encoding/json"
	service "hireforwork-server/service/modules"
	"net/http"
	"strconv"
)

type FieldHandler struct {
	FieldService *service.FieldService
}

func (f *FieldHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handlerFunc := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		default:
			f.GetField(w, r)
		}
	})

	// Áp dụng decorator nếu có
	// if h.decorator != nil {
	// 	handlerFunc = h.decorator(handlerFunc)
	// }

	handlerFunc.ServeHTTP(w, r)
}

func (f *FieldHandler) GetField(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	page, _ := strconv.Atoi(pageStr)

	pageSizeStr := r.URL.Query().Get("pageSize")
	pageSize, _ := strconv.Atoi(pageSizeStr)
	fieldList, err := f.FieldService.GetField(page, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(fieldList); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}
