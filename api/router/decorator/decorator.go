package decorator

import (
	"net/http"
)

// HTTPMethod represents supported HTTP methods
type HTTPMethod string

const (
	GET     HTTPMethod = "GET"
	POST    HTTPMethod = "POST"
	PUT     HTTPMethod = "PUT"
	DELETE  HTTPMethod = "DELETE"
	PATCH   HTTPMethod = "PATCH"
	OPTIONS HTTPMethod = "OPTIONS"
)

// RouteMetadata holds the metadata for a route
type RouteMetadata struct {
	Path         string
	Method       HTTPMethod
	RequiresAuth bool
	Handler      http.HandlerFunc
}
type Controller interface {
	Routes() []RouteMetadata
}

func Route(path string, method HTTPMethod, requiresAuth bool) RouteMetadata {
	return RouteMetadata{
		Path:         path,
		Method:       method,
		RequiresAuth: requiresAuth,
	}
}

func Get(path string, requiresAuth bool) RouteMetadata {
	return Route(path, GET, requiresAuth)
}

func Post(path string, requiresAuth bool) RouteMetadata {
	return Route(path, POST, requiresAuth)
}

func Put(path string, requiresAuth bool) RouteMetadata {
	return Route(path, PUT, requiresAuth)
}

func Delete(path string, requiresAuth bool) RouteMetadata {
	return Route(path, DELETE, requiresAuth)
}

func Patch(path string, requiresAuth bool) RouteMetadata {
	return Route(path, PATCH, requiresAuth)
}

func Options(path string, requiresAuth bool) RouteMetadata {
	return Route(path, OPTIONS, requiresAuth)
}
