package api

import (
	"hireforwork-server/api/handlers"
	"hireforwork-server/service"

	"github.com/gorilla/mux"
)

var authService *service.AuthService
var jobHandler *handlers.JobHandler // Khai báo JobHandler
var jobService *service.JobService  // Đảm bảo biến này được khai báo

func SetUpRouter() *mux.Router {
	router := mux.NewRouter()

	// authService = &service.AuthService{JwtSecret: []byte(os.Getenv("SECRET_KEY"))}

	// Khởi tạo jobService và jobHandler
	jobService = &service.JobService{Collection: service.JobCollection} // Sử dụng JobCollection từ package service
	jobHandler = &handlers.JobHandler{JobService: jobService}           // Khởi tạo JobHandler với JobService
	// handler := &handlers.Handler{
	// 	AuthService: authService,
	// }

	//---USER ROUTER---//
	//PUBLIC ROUTER//
	// router.HandleFunc("/careers/auth/login", handler.Login).Methods("POST")
	router.HandleFunc("/careers/create", handlers.CreateUser).Methods("POST")
	router.HandleFunc("/careers", handlers.GetUser).Methods("GET")
	router.HandleFunc("/careers/{id}", handlers.GetUserByID).Methods("GET")
	router.HandleFunc("/careers/{id}", handlers.DeleteUserByID).Methods("DELETE")
	router.HandleFunc("/careers/create", handlers.CreateUser).Methods("POST")
	//AUTH ROUTER
	// careerRouter := router.PathPrefix("/careers").Subrouter()
	// careerRouter.Use(middleware.JWTMiddleware(authService))
	// careerRouter.HandleFunc("", handlers.GetUser).Methods("GET")
	// careerRouter.HandleFunc("/{id}", handlers.GetUserByID).Methods("GET")
	// careerRouter.HandleFunc("/{id}", handlers.DeleteUserByID).Methods("DELETE")

	//Post router
	router.HandleFunc("/careers/savejob", handlers.SaveJob).Methods("POST")
	router.HandleFunc("/careers/viewedjob", handlers.CareerViewedJob).Methods("POST")

	//Job Router
	router.HandleFunc("/jobs", handlers.GetJob).Methods("GET")
	router.HandleFunc("/jobs/{id}/apply", handlers.ApplyJob).Methods("POST")

	//Company Router
	// router.HandleFunc("/companies/auth/login", handler.LoginCompany).Methods("POST")
	router.HandleFunc("/companies", handlers.GetCompaniesHandler).Methods("GET")

	//Update User Router
	router.HandleFunc("/careers/{id}", handlers.UpdateUser).Methods("PUT")
	router.HandleFunc("/companies/{id}", handlers.GetCompanyByID).Methods("GET")
	router.HandleFunc("/companies/create", handlers.CreateCompany).Methods("POST")
	router.HandleFunc("/companies/{id}", handlers.DeleteCompanyByID).Methods("DELETE")
	router.HandleFunc("/companies/update/{id}", handlers.UpdateCompanyByID).Methods("PUT")
	router.HandleFunc("/suggest", jobHandler.GetSuggestJobs).Methods("GET") // Sử dụng JobHandler để gọi GetSuggestJobs
	return router
}
