package api

import (
	"hireforwork-server/api/handlers"

	"github.com/gorilla/mux"
)

func SetUpRouter() *mux.Router {
	router := mux.NewRouter()

	//User Router
	router.HandleFunc("/careers", handlers.GetUser).Methods("GET")

	return router
}
