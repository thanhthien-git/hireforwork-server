package groups

import (
	"hireforwork-server/api/router/decorator"
	"hireforwork-server/api/router/types"
)

// CareerRoutes returns all career-related routes using decorator pattern
func CareerRoutes() []types.RouteConfig {
	routes := []decorator.RouteMetadata{
		decorator.Post("/careers/auth/login", false),
		decorator.Post("/careers/register", false),
		decorator.Post("/careers/create", false),
		decorator.Get("/careers", true),
		decorator.Get("/careers/{id}", true),
		decorator.Delete("/careers/{id}", true),
		decorator.Get("/careers/{id}/save-job", true),
		decorator.Get("/careers/{id}/applied-job", true),
		decorator.Post("/careers/{id}/upload-image", true),
		decorator.Post("/careers/{id}/upload-resume", true),
		decorator.Post("/careers/{id}/remove-resume", true),
		decorator.Post("/careers/{id}/update", true),
	}

	// Convert decorator metadata to RouteConfig
	configs := make([]types.RouteConfig, len(routes))
	for i, route := range routes {
		configs[i] = types.RouteConfig{
			Path:         route.Path,
			Handler:      "career",
			Methods:      []string{string(route.Method)},
			RequiresAuth: route.RequiresAuth,
		}
	}

	return configs
}
