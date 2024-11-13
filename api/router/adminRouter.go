package api

import (
	"hireforwork-server/api/handlers"

	"github.com/gorilla/mux"
)

func setAdminRouter(router *mux.Router, handler *handlers.Handler) {
	router.HandleFunc("/admin/static", handlers.GetStaticHandler).Methods("GET")
}
