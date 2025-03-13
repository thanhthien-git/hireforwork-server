package service

import (
	"hireforwork-server/db"
	"hireforwork-server/service"
	modules "hireforwork-server/service/modules"
	auth "hireforwork-server/service/modules/auth"
	job "hireforwork-server/service/modules/jobs"
)

// ServiceDependencies holds all dependencies that services might need
type ServiceDependencies struct {
	DB *db.DB
	// Add other dependencies here as needed
}

// ServiceFactory handles service creation and dependency injection
type ServiceFactory struct {
	deps *ServiceDependencies
}

// ServiceCreator is a function type that creates a service with dependencies
type ServiceCreator func(*ServiceDependencies) interface{}

// serviceCreators maps service types to their creation functions
var serviceCreators = map[string]ServiceCreator{
	"auth":     func(deps *ServiceDependencies) interface{} { return auth.NewAuthService(deps.DB) },
	"job":      func(deps *ServiceDependencies) interface{} { return job.NewJobService(deps.DB) },
	"company":  func(deps *ServiceDependencies) interface{} { return modules.NewCompanyService(deps.DB) },
	"career":   func(deps *ServiceDependencies) interface{} { return modules.NewUserService(deps.DB) },
	"tech":     func(deps *ServiceDependencies) interface{} { return modules.NewTechService(deps.DB) },
	"category": func(deps *ServiceDependencies) interface{} { return modules.NewCategoryService(deps.DB) },
	"field":    func(deps *ServiceDependencies) interface{} { return modules.NewFieldService(deps.DB) },
}

// NewServiceFactory creates a new service factory with all required dependencies
func NewServiceFactory(deps *ServiceDependencies) *ServiceFactory {
	return &ServiceFactory{deps: deps}
}

// CreateService creates a new service instance with injected dependencies
func (f *ServiceFactory) CreateService(serviceType string) interface{} {
	if creator, exists := serviceCreators[serviceType]; exists {
		return creator(f.deps)
	}
	return nil
}

// RegisterAllServices registers all available services in the container
func (f *ServiceFactory) RegisterAllServices(container *service.ServiceContainer) {
	for serviceType := range serviceCreators {
		container.Register(serviceType, f.CreateService(serviceType))
	}
}
