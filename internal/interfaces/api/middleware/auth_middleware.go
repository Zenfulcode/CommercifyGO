package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/infrastructure/auth"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// AuthMiddleware handles authentication
type AuthMiddleware struct {
	jwtService *auth.JWTService
	logger     logger.Logger
}

type contextKey string

const (
	UserIDKey contextKey = "user_id"
	emailKey  contextKey = "email"
	RoleKey   contextKey = "role"
)

// NewAuthMiddleware creates a new AuthMiddleware
func NewAuthMiddleware(jwtService *auth.JWTService, logger logger.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		jwtService: jwtService,
		logger:     logger,
	}
}

// Authenticate authenticates a request
func (m *AuthMiddleware) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header is required", http.StatusUnauthorized)
			return
		}

		// Check if the header has the Bearer prefix
		if !strings.HasPrefix(authHeader, "Bearer ") {
			http.Error(w, "Invalid authorization format", http.StatusUnauthorized)
			return
		}

		// Extract token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Validate token
		claims, err := m.jwtService.ValidateToken(tokenString)
		if err != nil {
			m.logger.Error("Invalid token: %v", err)
			http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), UserIDKey, claims.UserID)
		ctx = context.WithValue(ctx, emailKey, claims.Email)
		ctx = context.WithValue(ctx, RoleKey, claims.Role)

		// Call the next handler with the updated context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AdminOnly middleware ensures the user has admin role
func AdminOnly(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get role from context
		role, ok := r.Context().Value(RoleKey).(string)
		if !ok || role != string(entity.RoleAdmin) {
			http.Error(w, "Admin access required", http.StatusForbidden)
			return
		}

		// Call the next handler
		next.ServeHTTP(w, r)
	})
}
