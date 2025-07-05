// Package container provides a dependency injection container for the application
package container

import (
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"gorm.io/gorm"
)

// Container defines the interface for dependency injection container
type Container interface {
	// Config returns the application configuration
	Config() *config.Config

	// DB returns the database connection
	DB() *gorm.DB

	// Logger returns the application logger
	Logger() logger.Logger

	// Repositories provides access to all repositories
	Repositories() RepositoryProvider

	// Services provides access to all services
	Services() ServiceProvider

	// UseCases provides access to all use cases
	UseCases() UseCaseProvider

	// Handlers provides access to all handlers
	Handlers() HandlerProvider

	// Middlewares provides access to all middlewares
	Middlewares() MiddlewareProvider
}

// DIContainer is the concrete implementation of the Container interface
type DIContainer struct {
	config *config.Config
	db     *gorm.DB
	logger logger.Logger

	// Providers
	repositories RepositoryProvider
	services     ServiceProvider
	useCases     UseCaseProvider
	handlers     HandlerProvider
	middlewares  MiddlewareProvider
}

// NewContainer creates a new dependency injection container
func NewContainer(config *config.Config, db *gorm.DB, logger logger.Logger) Container {
	container := &DIContainer{
		config: config,
		db:     db,
		logger: logger,
	}

	// Initialize providers
	container.repositories = NewRepositoryProvider(container)
	container.services = NewServiceProvider(container)
	container.useCases = NewUseCaseProvider(container)
	container.handlers = NewHandlerProvider(container)
	container.middlewares = NewMiddlewareProvider(container)

	return container
}

// Config returns the application configuration
func (c *DIContainer) Config() *config.Config {
	return c.config
}

// DB returns the database connection
func (c *DIContainer) DB() *gorm.DB {
	return c.db
}

// Logger returns the application logger
func (c *DIContainer) Logger() logger.Logger {
	return c.logger
}

// Repositories provides access to all repositories
func (c *DIContainer) Repositories() RepositoryProvider {
	return c.repositories
}

// Services provides access to all services
func (c *DIContainer) Services() ServiceProvider {
	return c.services
}

// UseCases provides access to all use cases
func (c *DIContainer) UseCases() UseCaseProvider {
	return c.useCases
}

// Handlers provides access to all handlers
func (c *DIContainer) Handlers() HandlerProvider {
	return c.handlers
}

// Middlewares provides access to all middlewares
func (c *DIContainer) Middlewares() MiddlewareProvider {
	return c.middlewares
}
