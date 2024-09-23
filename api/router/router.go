package api

import (
	"hireforwork-server/api/handlers"

	"github.com/gorilla/mux"
)

func SetUpRouter() *mux.Router {
	router := mux.NewRouter()

	//---USER ROUTER---//
	//Get router
	router.HandleFunc("/careers", handlers.GetUser).Methods("GET")
	router.HandleFunc("/careers/{id}", handlers.GetUserByID).Methods("GET")
	//Delete router
	router.HandleFunc("/careers/{id}", handlers.DeleteUserByID).Methods("DELETE")
	//Post router
	router.HandleFunc("/careers/create", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/careers/uploadImage", handlers.UploadImage).Methods("POST")
	router.HandleFunc("/careers/uploadResume", handlers.UploadResume).Methods("POST")
	//---USER ROUTER---//

	return router
}
