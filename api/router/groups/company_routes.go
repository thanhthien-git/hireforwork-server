package groups

import (
	"hireforwork-server/api/router/decorator"
	"hireforwork-server/api/router/types"
)

// CompanyRoutes returns all company-related routes using decorator pattern
func CompanyRoutes() []types.RouteConfig {
	routes := []decorator.RouteMetadata{
		decorator.Get("/companies", false),
		decorator.Get("/companies/{id}", false),
		decorator.Post("/companies", true),
		decorator.Put("/companies/{id}", true),
		decorator.Delete("/companies/{id}", true),
	}

	// Convert decorator metadata to RouteConfig
	configs := make([]types.RouteConfig, len(routes))
	for i, route := range routes {
		configs[i] = types.RouteConfig{
			Path:         route.Path,
			Handler:      "company",
			Methods:      []string{string(route.Method)},
			RequiresAuth: route.RequiresAuth,
		}
	}

	return configs
}
