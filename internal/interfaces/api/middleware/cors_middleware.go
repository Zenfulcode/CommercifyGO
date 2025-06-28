package middleware

import (
	"fmt"
	"net/http"

	"slices"

	"github.com/zenfulcode/commercify/config"
)

// CorsMiddleware handles CORS (Cross-Origin Resource Sharing)
type CorsMiddleware struct {
	config *config.Config
}

// NewCorsMiddleware creates a new CorsMiddleware
func NewCorsMiddleware(config *config.Config) *CorsMiddleware {
	return &CorsMiddleware{
		config: config,
	}
}

// ApplyCors adds CORS headers to responses
func (m *CorsMiddleware) ApplyCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get allowed origins from config or use default
		allowedOrigins := m.getAllowedOrigins()

		if len(allowedOrigins) == 0 {
			// If no origins are configured, allow all origins
			allowedOrigins = []string{"*"}
		}

		fmt.Println("Allowed Origins:", allowedOrigins)

		// Get origin from request
		origin := r.Header.Get("Origin")

		fmt.Println("Request Origin:", origin)

		// Check if the origin is allowed
		if m.isAllowedOrigin(origin, allowedOrigins) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
		}

		// Set standard CORS headers
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Accept, Origin, Cache-Control, X-Requested-With")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		// Proceed with the next handler
		next.ServeHTTP(w, r)
	})
}

// getAllowedOrigins returns the list of allowed origins
func (m *CorsMiddleware) getAllowedOrigins() []string {
	return m.config.CORS.AllowedOrigins
}

// isAllowedOrigin checks if the origin is in the allowed list or if all origins are allowed
func (m *CorsMiddleware) isAllowedOrigin(origin string, allowedOrigins []string) bool {
	if origin == "" {
		return false
	}

	// Check if "*" is in the allowed origins list
	if slices.Contains(allowedOrigins, "*") {
		return true
	}

	// Check if the specific origin is allowed
	return slices.Contains(allowedOrigins, origin)
}
