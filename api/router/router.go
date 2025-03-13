package api

import (
	"hireforwork-server/api/handlers"
	"hireforwork-server/db"
	"hireforwork-server/middleware"
	service "hireforwork-server/service"
	auth "hireforwork-server/service/modules/auth"

	"github.com/gorilla/mux"
)

/*
Design Patterns Used in Router Setup:

1. Registry Pattern
  - Centralizes route configuration in one place
  - Makes it easy to manage and modify routes
  - Provides a single source of truth for route definitions
*/
type RouteConfig struct {
	Path         string
	Handler      string
	Methods      []string
	RequiresAuth bool
}

// RouteRegistry implements the Registry Pattern
// - Holds all route configurations in one place
// - Makes it easy to add/modify routes without changing code
// - Provides a centralized configuration point
type RouteRegistry struct {
	routes []RouteConfig
}

// NewRouteRegistry creates a new route registry with predefined routes
// - Implements the Factory Pattern for creating the registry
// - Routes are defined declaratively rather than imperatively
func NewRouteRegistry() *RouteRegistry {
	return &RouteRegistry{
		routes: []RouteConfig{
			{
				Path:    "/jobs",
				Handler: "job",
				Methods: []string{"GET"},
			},
			{
				Path:         "/jobs",
				Handler:      "job",
				Methods:      []string{"POST", "PUT"},
				RequiresAuth: true,
			},
			{
				Path:    "/companies",
				Handler: "company",
			},
			{
				Path:    "/careers",
				Handler: "career",
			},
			{
				Path:    "/tech",
				Handler: "tech",
			},
			{
				Path:    "/field",
				Handler: "field",
			},
			{
				Path:    "/category",
				Handler: "category",
			},
			{
				Path:    "/admin",
				Handler: "admin",
			},
		},
	}
}

/*
2. Builder Pattern
  - Separates the construction of the router from its representation
  - Allows for different representations of the same construction process
  - Makes it easy to add new features to the router construction
*/
type RouterBuilder struct {
	router   *mux.Router
	services *service.AppServices
	db       *db.DB
}

// NewRouterBuilder implements the Factory Pattern
// - Creates a new RouterBuilder instance
// - Injects dependencies (services and db)
func NewRouterBuilder(services *service.AppServices, db *db.DB) *RouterBuilder {
	return &RouterBuilder{
		router:   mux.NewRouter(),
		services: services,
		db:       db,
	}
}

// BuildRoutes implements the Builder Pattern
// - Constructs the router using the registry
// - Handles the creation and configuration of handlers
// - Makes the construction process flexible and extensible
func (b *RouterBuilder) BuildRoutes(registry *RouteRegistry) *mux.Router {
	for _, route := range registry.routes {
		// Uses the HandlerBuilder (another Builder Pattern implementation)
		handler := handlers.NewHandlerBuilder(b.services, route.Handler, b.db).Build()
		if handler != nil {
			// Create a subrouter for this specific route
			r := b.router.PathPrefix(route.Path).Subrouter()

			// Apply auth middleware if required
			if route.RequiresAuth {
				authService := auth.NewAuthService(b.db)
				r.Use(middleware.JWTMiddleware(authService))
			}

			// Handle the route with specified methods
			if len(route.Methods) > 0 {
				r.Handle("", handler).Methods(route.Methods...)
			} else {
				r.Handle("", handler)
			}

			// Handle the root path as well
			if route.Path != "/" {
				r.Handle("/", handler)
			}
		}
	}
	return b.router
}

/*
3. Facade Pattern
  - Provides a simplified interface to the complex router setup
  - Hides the complexity of route registration and handler creation
  - Makes it easy to use the router system
*/
func SetUpRouter(services *service.AppServices, dbInstance *db.DB) *mux.Router {
	registry := NewRouteRegistry()
	builder := NewRouterBuilder(services, dbInstance)
	return builder.BuildRoutes(registry)
}
