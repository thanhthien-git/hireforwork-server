package middleware

import (
	"context"
	auth "hireforwork-server/service/modules/auth"
	"log"
	"net/http"
	"strings"
	"time"
)

// Key for context values
type contextKey string

const UserIDKey contextKey = "userID"

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ip := r.RemoteAddr
		userAgent := r.Header.Get("User-Agent")

		log.Printf(
			"Started %s %s from %s (User-Agent: %s)",
			r.Method,
			r.URL.Path,
			ip,
			userAgent,
		)

		next.ServeHTTP(w, r)

		log.Printf(
			"Completed %s %s in %v",
			r.Method,
			r.URL.Path,
			time.Since(start),
		)
	})
}

// GlobalMiddleware checks for authorization header and decodes it if present
func GlobalMiddleware(authService *auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check for authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader != "" {
				// Try to decode token
				tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer ", "", 1))
				claims, err := authService.ValidateToken(tokenString)

				if err == nil {
					// Token is valid, add user info to context
					ctx := context.WithValue(r.Context(), UserIDKey, claims.Subject)
					r = r.WithContext(ctx)
					log.Printf("User authenticated: %s", claims.Subject)
				}
			}
			next.ServeHTTP(w, r)
		})
	}
}

// JWTMiddleware verifies that the request has valid authentication
func JWTMiddleware(authService *auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			userID := GetUserID(r)
			if userID == "" {
				http.Error(w, "Authentication required", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// GetUserID helper function to get userID from context
func GetUserID(r *http.Request) string {
	if userID, ok := r.Context().Value(UserIDKey).(string); ok {
		return userID
	}
	return ""
}
