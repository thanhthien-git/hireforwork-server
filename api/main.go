package api

import (
	api "hireforwork-server/api/router"
	"net/http"

	"github.com/rs/cors"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	router := api.SetUpRouter()

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	})

	c.Handler(router).ServeHTTP(w, r)
}
