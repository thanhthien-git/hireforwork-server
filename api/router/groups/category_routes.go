package groups

import (
	"hireforwork-server/api/router/types"
)

// CategoryRoutes trả về tất cả các route liên quan đến danh mục
func CategoryRoutes() []types.RouteConfig {
	return []types.RouteConfig{
		{
			Path:         "/category",
			Methods:      []string{"GET"},
			RequiresAuth: false,
			Handler:      "category",
		},
		{
			Path:         "/category",
			Methods:      []string{"POST"},
			RequiresAuth: true,
			Handler:      "category",
		},
		{
			Path:         "/category/{id}",
			Methods:      []string{"PUT"},
			RequiresAuth: true,
			Handler:      "category",
		},
		{
			Path:         "/category/{id}",
			Methods:      []string{"DELETE"},
			RequiresAuth: true,
			Handler:      "category",
		},
	}
}
