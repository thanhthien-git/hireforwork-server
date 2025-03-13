package middleware

import (
	"context"
	auth "hireforwork-server/service/modules/auth"
	"net/http"
	"strings"
)

func JWTMiddleware(authService *auth.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Authorization header is required", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer ", "", 1))

			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Invalid token: "+err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "userID", claims.Subject)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
