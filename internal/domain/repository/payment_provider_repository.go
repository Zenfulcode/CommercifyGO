package repository

import (
	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// PaymentProviderRepository defines the interface for payment provider operations
type PaymentProviderRepository interface {
	// Create creates a new payment provider
	Create(provider *entity.PaymentProvider) error

	// Update updates an existing payment provider
	Update(provider *entity.PaymentProvider) error

	// Delete deletes a payment provider
	Delete(id uint) error

	// GetByID returns a payment provider by ID
	GetByID(id uint) (*entity.PaymentProvider, error)

	// GetByType returns a payment provider by type
	GetByType(providerType common.PaymentProviderType) (*entity.PaymentProvider, error)

	// List returns all payment providers with pagination
	List(offset, limit int) ([]*entity.PaymentProvider, error)

	// GetEnabled returns all enabled payment providers
	GetEnabled() ([]*entity.PaymentProvider, error)

	// GetEnabledByMethod returns enabled payment providers that support a specific payment method
	GetEnabledByMethod(method common.PaymentMethod) ([]*entity.PaymentProvider, error)

	// GetEnabledByCurrency returns enabled payment providers that support a specific currency
	GetEnabledByCurrency(currency string) ([]*entity.PaymentProvider, error)

	// GetEnabledByMethodAndCurrency returns enabled payment providers that support both method and currency
	GetEnabledByMethodAndCurrency(method common.PaymentMethod, currency string) ([]*entity.PaymentProvider, error)

	// UpdateWebhookInfo updates webhook information for a payment provider
	UpdateWebhookInfo(providerType common.PaymentProviderType, webhookURL, webhookSecret, externalWebhookID string, events []string) error

	// GetWithWebhooks returns payment providers that have webhook configurations
	GetWithWebhooks() ([]*entity.PaymentProvider, error)
}
