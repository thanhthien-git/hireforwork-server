package middleware

import (
	"context"
	"hireforwork-server/service"
	"net/http"
	"strings"
)

func JWTMiddleware(authService *service.AuthService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Yêu cầu phải có header Authorization", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimSpace(strings.Replace(authHeader, "Bearer ", "", 1))

			// Validate the token
			claims, err := authService.ValidateToken(tokenString)
			if err != nil {
				http.Error(w, "Token không hợp lệ: "+err.Error(), http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "username", claims.Username)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
