package decorator

import (
	"fmt"
	"net/http"
)

// HandlerFunc is a type for HTTP handler functions
type HandlerFunc func(http.ResponseWriter, *http.Request)

// LoggingDecorator adds logging to a handler
func LoggingDecorator(h HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Request received:", r.URL.Path)
		h(w, r)
	}
}

// AuthDecorator adds authentication to a handler
func AuthDecorator(h HandlerFunc) HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Checking authentication...")
		h(w, r)
	}
}
