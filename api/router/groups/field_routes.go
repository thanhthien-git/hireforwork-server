package groups

import (
	"hireforwork-server/api/router/types"
)

// FieldRoutes trả về tất cả các route liên quan đến lĩnh vực
func FieldRoutes() []types.RouteConfig {
	return []types.RouteConfig{
		{
			Path:         "/field",
			Methods:      []string{"GET"},
			RequiresAuth: false,
			Handler:      "field",
		},
		{
			Path:         "/field",
			Methods:      []string{"POST"},
			RequiresAuth: true,
			Handler:      "field",
		},
		{
			Path:         "/field/{id}",
			Methods:      []string{"PUT"},
			RequiresAuth: true,
			Handler:      "field",
		},
		{
			Path:         "/field/{id}",
			Methods:      []string{"DELETE"},
			RequiresAuth: true,
			Handler:      "field",
		},
	}
}
