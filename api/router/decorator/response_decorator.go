package decorator

import (
	"net/http"
)

/*
1. Decorator Pattern - use to add functionality to the response writer
*/
// ResponseWriter wraps http.ResponseWriter để thêm các chức năng mới
type ResponseWriter struct {
	http.ResponseWriter
	status int
}

// WriteHeader ghi đè phương thức của http.ResponseWriter để theo dõi status code
func (rw *ResponseWriter) WriteHeader(status int) {
	rw.status = status
	rw.ResponseWriter.WriteHeader(status)
}

// WithJSONResponse là decorator để tự động thêm JSON Content-Type header
func WithJSONResponse(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Wrap ResponseWriter
		wrappedWriter := &ResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}

		// Thêm JSON Content-Type header
		wrappedWriter.Header().Set("Content-Type", "application/json")

		// Gọi handler tiếp theo với wrapped writer
		next.ServeHTTP(wrappedWriter, r)
	})
}

// WithSecurityHeaders là decorator để thêm các security headers
func WithSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Thêm các security headers
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")

		next.ServeHTTP(w, r)
	})
}

// WithCORS là decorator để thêm CORS headers
func WithCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
