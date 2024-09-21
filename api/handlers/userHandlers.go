package handlers

import (
	"encoding/json"
	"hireforwork-server/service"
	"net/http"
	"strconv"
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

// func GetUserByID(w http.ResponseWriter, r *http.Request) {
// 	client, ctx, err := dbHelper.ConnectDB()
// 	dbHelper.ValidateError(err, w)
// 	defer client.Disconnect(ctx)

// 	user, err := service.GetUserByID(r.Context(), utils.GetQueryID(r))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusNotFound)
// 		return
// 	}

// 	w.Header().Set("Content-Type", "application/json")
// 	err = json.NewEncoder(w).Encode(user)
// 	if err != nil {
// 		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
// 	}

// }

// func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
// 	client, ctx, err := dbHelper.ConnectDB()
// 	dbHelper.ValidateError(err, w)
// 	defer client.Disconnect(ctx)

// 	_, err = service.DeleteUserByID(r.Context(), utils.GetQueryID(r))
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusOK)
// }
