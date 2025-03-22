package groups

import (
	"hireforwork-server/api/router/decorator"
	"hireforwork-server/api/router/types"
)

// JobRoutes returns all job-related routes using decorator pattern
func JobRoutes() []types.RouteConfig {
	routes := []decorator.RouteMetadata{
		decorator.Get("/jobs", false),
		decorator.Post("/jobs", true),
		decorator.Put("/jobs", true),
		decorator.Get("/jobs/{id}", false),
		decorator.Post("/jobs/{id}/apply", true),
		decorator.Post("/jobs/{id}/save", true),
		decorator.Post("/jobs/{id}/unsave", true),
		decorator.Post("/jobs/{id}/apply", true),
		decorator.Put("/jobs/{id}", true),
		decorator.Delete("/jobs/{id}", true),
	}

	// Convert decorator metadata to RouteConfig
	configs := make([]types.RouteConfig, len(routes))
	for i, route := range routes {
		configs[i] = types.RouteConfig{
			Path:         route.Path,
			Handler:      "job",
			Methods:      []string{string(route.Method)},
			RequiresAuth: route.RequiresAuth,
		}
	}

	return configs
}
