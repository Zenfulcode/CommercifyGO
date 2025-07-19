package container

import (
	"sync"

	"github.com/zenfulcode/commercify/internal/interfaces/api/middleware"
)

// MiddlewareProvider provides access to all middlewares
type MiddlewareProvider interface {
	AuthMiddleware() *middleware.AuthMiddleware
	CorsMiddleware() *middleware.CorsMiddleware
}

// middlewareProvider is the concrete implementation of MiddlewareProvider
type middlewareProvider struct {
	container Container
	mu        sync.Mutex

	authMiddleware *middleware.AuthMiddleware
	corsMiddleware *middleware.CorsMiddleware
}

// NewMiddlewareProvider creates a new middleware provider
func NewMiddlewareProvider(container Container) MiddlewareProvider {
	return &middlewareProvider{
		container: container,
	}
}

// AuthMiddleware returns the authentication middleware
func (p *middlewareProvider) AuthMiddleware() *middleware.AuthMiddleware {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.authMiddleware == nil {
		p.authMiddleware = middleware.NewAuthMiddleware(
			p.container.Services().JWTService(),
			p.container.Logger(),
		)
	}
	return p.authMiddleware
}

// CorsMiddleware returns the CORS middleware
func (p *middlewareProvider) CorsMiddleware() *middleware.CorsMiddleware {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.corsMiddleware == nil {
		p.corsMiddleware = middleware.NewCorsMiddleware(
			p.container.Config(),
		)
	}
	return p.corsMiddleware
}
