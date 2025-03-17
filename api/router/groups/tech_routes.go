package groups

import (
	"hireforwork-server/api/router/types"
)

// TechRoutes trả về tất cả các route liên quan đến công nghệ
func TechRoutes() []types.RouteConfig {
	return []types.RouteConfig{
		{
			Path:         "/tech",
			Methods:      []string{"GET"},
			RequiresAuth: false,
			Handler:      "tech",
		},
		{
			Path:         "/tech",
			Methods:      []string{"POST"},
			RequiresAuth: true,
			Handler:      "tech",
		},
		{
			Path:         "/tech/{id}",
			Methods:      []string{"PUT"},
			RequiresAuth: true,
			Handler:      "tech",
		},
		{
			Path:         "/tech/{id}",
			Methods:      []string{"DELETE"},
			RequiresAuth: true,
			Handler:      "tech",
		},
	}
}
