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

// Controller interface defines the base for all controllers
type Controller interface {
	Routes() []RouteMetadata
}

// Route is a decorator-like function that creates route metadata
func Route(path string, method HTTPMethod, requiresAuth bool) RouteMetadata {
	return RouteMetadata{
		Path:         path,
		Method:       method,
		RequiresAuth: requiresAuth,
	}
}

// Get decorator for GET routes
func Get(path string, requiresAuth bool) RouteMetadata {
	return Route(path, GET, requiresAuth)
}

// Post decorator for POST routes
func Post(path string, requiresAuth bool) RouteMetadata {
	return Route(path, POST, requiresAuth)
}

// Put decorator for PUT routes
func Put(path string, requiresAuth bool) RouteMetadata {
	return Route(path, PUT, requiresAuth)
}

// Delete decorator for DELETE routes
func Delete(path string, requiresAuth bool) RouteMetadata {
	return Route(path, DELETE, requiresAuth)
}

// Patch decorator for PATCH routes
func Patch(path string, requiresAuth bool) RouteMetadata {
	return Route(path, PATCH, requiresAuth)
}

// Options decorator for OPTIONS routes
func Options(path string, requiresAuth bool) RouteMetadata {
	return Route(path, OPTIONS, requiresAuth)
}
