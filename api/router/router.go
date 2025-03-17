package api

import (
	"hireforwork-server/api/handlers"
	"hireforwork-server/api/router/decorator"
	"hireforwork-server/api/router/groups"
	"hireforwork-server/api/router/types"
	"hireforwork-server/db"
	"hireforwork-server/middleware"
	service "hireforwork-server/service"
	auth "hireforwork-server/service/modules/auth"
	"net/http"

	"github.com/gorilla/mux"
)

/*
Design Patterns Used in Router Setup:

1. Decorator Pattern
   - Routes are defined using decorator-like functions
   - Makes route configuration more readable and maintainable
   - Provides type safety for HTTP methods

2. Builder Pattern
   - Separates router construction from its representation
   - Allows for different representations of the same construction process
   - Makes it easy to add new features to the router construction

3. Facade Pattern
   - Provides a simplified interface to the complex router setup
   - Hides the complexity of route registration and handler creation
   - Makes it easy to use the router system
*/

// RouterBuilder implements the Builder Pattern
type RouterBuilder struct {
	router   *mux.Router
	services *service.AppServices
	db       *db.DB
}

// NewRouterBuilder creates a new RouterBuilder instance
func NewRouterBuilder(services *service.AppServices, db *db.DB) *RouterBuilder {
	return &RouterBuilder{
		router:   mux.NewRouter(),
		services: services,
		db:       db,
	}
}

// BuildRoutes constructs and configures the router
func (b *RouterBuilder) BuildRoutes() *mux.Router {
	// Get all route groups
	routes := []types.RouteConfig{}
	routes = append(routes, groups.CareerRoutes()...)
	routes = append(routes, groups.JobRoutes()...)
	routes = append(routes, groups.CompanyRoutes()...)

	// Create auth service
	authService := auth.NewAuthService(b.db)

	// Apply global middleware and decorators
	b.router.Use(middleware.GlobalMiddleware(authService))
	b.router.Use(decorator.WithJSONResponse)
	b.router.Use(decorator.WithSecurityHeaders)
	b.router.Use(decorator.WithCORS)

	// Register all routes
	for _, route := range routes {
		handler := handlers.NewHandlerBuilder(b.services, route.Handler, b.db).Build()
		if handler != nil {
			var finalHandler http.Handler = handler

			// Apply JWT middleware only if route requires auth
			if route.RequiresAuth {
				finalHandler = middleware.JWTMiddleware(authService)(handler)
			}

			// Create route with methods
			r := b.router.Handle(route.Path, finalHandler)
			if len(route.Methods) > 0 {
				r.Methods(route.Methods...)
			}
		}
	}

	return b.router
}

// SetUpRouter is the facade for router setup
func SetUpRouter(services *service.AppServices, dbInstance *db.DB) *mux.Router {
	builder := NewRouterBuilder(services, dbInstance)
	return builder.BuildRoutes()
}

// Wrap http.ResponseWriter để thêm chức năng mới
type ResponseWriter struct {
	http.ResponseWriter
	status int
}
