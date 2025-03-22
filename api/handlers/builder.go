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

type GenericHandler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
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
type HandlerConfig struct {
	HandlerType    reflect.Type
	ServiceName    string
	ServiceType    reflect.Type
	RequiresAuth   bool
	FallbackCreate func(*db.DB) interface{}
}

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

// Simple Factory  Pattern
/*
1. function creates handlers based on configuration
2. uses reflection to set up handler dependencies
3. uses dependency injection to inject services into handlers
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

func (b *HandlerBuilder) Build() GenericHandler {
	return *b.handler
}
