package handlers

import (
	"encoding/json"
	dbHelper "hireforwork-server/db"
	"hireforwork-server/service"
	"hireforwork-server/utils"
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) {

	client, ctx, err := dbHelper.ConnectDB()
	dbHelper.ValidateError(err, w)
	defer client.Disconnect(ctx)

	users, err := service.GetUser(ctx)
	dbHelper.ValidateError(err, w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}
}

func GetUserByID(w http.ResponseWriter, r *http.Request) {
	client, ctx, err := dbHelper.ConnectDB()
	dbHelper.ValidateError(err, w)
	defer client.Disconnect(ctx)

	user, err := service.GetUserByID(r.Context(), utils.GetQueryID(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}

}

func DeleteUserByID(w http.ResponseWriter, r *http.Request) {
	client, ctx, err := dbHelper.ConnectDB()
	dbHelper.ValidateError(err, w)
	defer client.Disconnect(ctx)

	_, err = service.DeleteUserByID(r.Context(), utils.GetQueryID(r))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}
