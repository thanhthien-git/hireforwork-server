package types

// RouteConfig defines the structure for route configuration
type RouteConfig struct {
	Path         string
	Handler      string
	Methods      []string
	RequiresAuth bool
}

// RouteGroup defines a group of routes with a common prefix
type RouteGroup struct {
	Prefix string
	Routes []RouteConfig
}
