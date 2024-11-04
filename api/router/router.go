package api

import (
	"hireforwork-server/api/handlers"
	"hireforwork-server/service"
	"os"

	"github.com/gorilla/mux"
)

func SetUpRouter() *mux.Router {
	router := mux.NewRouter()

	authService := &service.AuthService{JwtSecret: []byte(os.Getenv("SECRET_KEY"))}

	handler := &handlers.Handler{
		AuthService: authService,
	}
	setUpCategoryRouter(router, handler)
	setUpCareerRoutes(router, handler)
	setUpCompanyRoutes(router, handler)
	setUpJobRouter(router, handler)
	setUpTechRouter(router, handler)
	setUpCategoryRouter(router, handler)
	setUpFieldRouter(router, handler)
	return router
}
