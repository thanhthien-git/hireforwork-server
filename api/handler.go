package api

import (
	api "hireforwork-server/api/router"
	"net/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	router := api.SetUpRouter()
	router.ServeHTTP(w, r)
}
