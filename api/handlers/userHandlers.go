package handlers

import (
	"encoding/json"
	dbHelper "hireforwork-server/db"
	"hireforwork-server/service"
	"net/http"
)

func GetUser(w http.ResponseWriter, r *http.Request) {
	client, ctx, err := dbHelper.ConnectDB()

	if err != nil {
		http.Error(w, "Error when connect to database", http.StatusInternalServerError)
		return
	}
	defer client.Disconnect(ctx)

	collection := client.Database("hideforwork").Collection("Career")
	users, err := service.GetUser(ctx, collection)

	if err != nil {
		http.Error(w, "Error when fetching data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Error encoding JSON", http.StatusInternalServerError)
	}

}
