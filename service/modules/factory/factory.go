package service

import (
	"fmt"
	"hireforwork-server/db"
	"hireforwork-server/service"
	modules "hireforwork-server/service/modules"
	auth "hireforwork-server/service/modules/auth"
	job "hireforwork-server/service/modules/jobs"
	observe "hireforwork-server/service/observe"
)

/*
1. Factory pattern
2. ServiceFactory is a factory for creating services
3. ServiceDependencies holds all dependencies that services might need
4. ServiceCreator is a function type that creates a service with dependencies
5. serviceCreators maps service types to their creation functions
*/

type ServiceDependencies struct {
	DB *db.DB
}

type ServiceFactory struct {
	deps *ServiceDependencies
}

type ServiceCreator func(*ServiceDependencies) interface{}

var serviceCreators = map[string]ServiceCreator{
	"auth":     func(deps *ServiceDependencies) interface{} { return auth.NewAuthService(deps.DB) },
	"job":      func(deps *ServiceDependencies) interface{} { return job.NewJobService(deps.DB) },
	"company":  func(deps *ServiceDependencies) interface{} { return modules.NewCompanyService(deps.DB) },
	"career":   func(deps *ServiceDependencies) interface{} { return modules.NewUserService(deps.DB) },
	"tech":     func(deps *ServiceDependencies) interface{} { return modules.NewTechService(deps.DB) },
	"category": func(deps *ServiceDependencies) interface{} { return modules.NewCategoryService(deps.DB) },
	"field":    func(deps *ServiceDependencies) interface{} { return modules.NewFieldService(deps.DB) },
	"observe": func(deps *ServiceDependencies) interface{} {
		return observe.NewJobEventManager()
	},
}

func NewServiceFactory(deps *ServiceDependencies) *ServiceFactory {
	return &ServiceFactory{deps: deps}
}

func (f *ServiceFactory) CreateService(serviceType string) interface{} {
	if creator, exists := serviceCreators[serviceType]; exists {
		return creator(f.deps)
	}
	return nil
}

func (f *ServiceFactory) RegisterAllServices(container *service.ServiceContainer) {
	fmt.Println("Registering all services...")
	for serviceType := range serviceCreators {
		fmt.Printf("Registering service: %s\n", serviceType)
		container.Register(serviceType, f.CreateService(serviceType))
	}
}
