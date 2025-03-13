package handlers

import (
	handlers "hireforwork-server/api/handlers/items"
	"hireforwork-server/db"
	"hireforwork-server/service"
	modules "hireforwork-server/service/modules"
	"hireforwork-server/service/modules/auth"
	"hireforwork-server/service/modules/jobs"
	"net/http"
	"reflect"
)

/*
Design Patterns Used in Handler Setup:

1. Interface Pattern
  - Defines a common contract for all handlers
  - Enables polymorphism and loose coupling
  - Makes it easy to swap implementations
*/
type GenericHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

/*
2. Base Handler Pattern
  - Provides common functionality for all handlers
  - Reduces code duplication
  - Makes it easy to add shared features
*/
type BaseHandler struct {
	services *service.AppServices
	db       *db.DB
}

/*
3. Builder Pattern
  - Separates handler construction from its representation
  - Makes it easy to create complex handlers
  - Provides a flexible way to configure handlers
*/
type HandlerBuilder struct {
	handler *GenericHandler
	service *service.AppServices
}

/*
4. Configuration Pattern
  - Separates configuration from implementation
  - Makes it easy to modify handler behavior
  - Provides a declarative way to define handlers
*/
type HandlerConfig struct {
	HandlerType    reflect.Type
	ServiceName    string
	ServiceType    reflect.Type
	RequiresAuth   bool
	FallbackCreate func(*db.DB) interface{}
}

/*
5. Registry Pattern
  - Centralizes handler configurations
  - Makes it easy to manage handler types
  - Provides a single source of truth for handler definitions
*/
var handlerConfigs = map[string]HandlerConfig{
	"job": {
		HandlerType: reflect.TypeOf(&handlers.JobHandler{}),
		ServiceName: "job",
		ServiceType: reflect.TypeOf(&jobs.JobService{}),
	},
	"company": {
		HandlerType:  reflect.TypeOf(&handlers.CompanyHandler{}),
		ServiceName:  "company",
		ServiceType:  reflect.TypeOf(&modules.CompanyService{}),
		RequiresAuth: true,
	},
	"career": {
		HandlerType:  reflect.TypeOf(&handlers.UserHandler{}),
		ServiceName:  "career",
		ServiceType:  reflect.TypeOf(&modules.UserService{}),
		RequiresAuth: true,
	},
	"tech": {
		HandlerType:    reflect.TypeOf(&handlers.TechHandler{}),
		ServiceName:    "tech",
		ServiceType:    reflect.TypeOf(&modules.TechService{}),
		FallbackCreate: func(db *db.DB) interface{} { return modules.NewTechService(db) },
	},
	"field": {
		HandlerType: reflect.TypeOf(&handlers.FieldHandler{}),
		ServiceName: "field",
		ServiceType: reflect.TypeOf(&modules.FieldService{}),
	},
	"category": {
		HandlerType:    reflect.TypeOf(&handlers.CategoryHandler{}),
		ServiceName:    "category",
		ServiceType:    reflect.TypeOf(&modules.CategoryService{}),
		FallbackCreate: func(db *db.DB) interface{} { return modules.NewCategoryService(db) },
	},
}

/*
6. Factory Pattern with Reflection
  - Creates handlers dynamically based on configuration
  - Uses reflection to set up handler dependencies
  - Makes it easy to add new handler types
*/
func createHandler(config HandlerConfig, services *service.AppServices, db *db.DB) GenericHandler {
	handlerValue := reflect.New(config.HandlerType.Elem())
	handler := handlerValue.Interface().(GenericHandler)

	// Get the service instance using Dependency Injection
	var serviceInstance interface{}
	if services != nil {
		serviceInstance = services.GetService(config.ServiceName)
		if serviceInstance == nil && config.FallbackCreate != nil {
			serviceInstance = config.FallbackCreate(db)
		}
	}

	// Set the service field using reflection
	if serviceInstance != nil {
		serviceField := handlerValue.Elem().FieldByName(config.ServiceType.Name())
		if serviceField.IsValid() && serviceField.CanSet() {
			serviceField.Set(reflect.ValueOf(serviceInstance))
		} else {
			// If the field is not found or cannot be set, try to find it by type
			for i := 0; i < handlerValue.Elem().NumField(); i++ {
				field := handlerValue.Elem().Field(i)
				if field.Type() == config.ServiceType {
					field.Set(reflect.ValueOf(serviceInstance))
					break
				}
			}
		}
	}

	// Set auth service if required (Conditional Dependency Injection)
	if config.RequiresAuth {
		authField := handlerValue.Elem().FieldByName("AuthService")
		if authField.IsValid() && authField.CanSet() {
			authField.Set(reflect.ValueOf(auth.NewAuthService(db)))
		}
	}

	return handler
}

/*
7. Factory Method Pattern
  - Creates handler instances based on type
  - Uses the configuration registry
  - Provides a flexible way to create handlers
*/
func NewHandlerBuilder(services *service.AppServices, handlerType string, dbInstance *db.DB) *HandlerBuilder {
	var handler GenericHandler
	if config, exists := handlerConfigs[handlerType]; exists {
		handler = createHandler(config, services, dbInstance)
	}
	return &HandlerBuilder{
		handler: &handler,
		service: services,
	}
}

/*
8. Builder Pattern Method
  - Returns the constructed handler
  - Completes the builder pattern implementation
*/
func (b *HandlerBuilder) Build() GenericHandler {
	return *b.handler
}
