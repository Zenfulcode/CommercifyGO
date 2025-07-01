package service

import (
	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// PaymentProviderService defines the interface for payment provider management
type PaymentProviderService interface {
	// GetPaymentProviders returns all payment providers
	GetPaymentProviders() ([]PaymentProvider, error)

	// GetEnabledPaymentProviders returns all enabled payment providers
	GetEnabledPaymentProviders() ([]PaymentProvider, error)

	// GetPaymentProvidersForCurrency returns payment providers that support a specific currency
	GetPaymentProvidersForCurrency(currency string) ([]PaymentProvider, error)

	// GetPaymentProvidersForMethod returns payment providers that support a specific payment method
	GetPaymentProvidersForMethod(method common.PaymentMethod) ([]PaymentProvider, error)

	// RegisterWebhook registers a webhook for a payment provider
	RegisterWebhook(providerType common.PaymentProviderType, webhookURL string, events []string) error

	// DeleteWebhook deletes a webhook for a payment provider
	DeleteWebhook(providerType common.PaymentProviderType) error

	// GetWebhookInfo returns webhook information for a payment provider
	GetWebhookInfo(providerType common.PaymentProviderType) (*entity.PaymentProvider, error)

	// UpdateProviderConfiguration updates the configuration for a payment provider
	UpdateProviderConfiguration(providerType common.PaymentProviderType, config common.JSONB) error

	// EnableProvider enables a payment provider
	EnableProvider(providerType common.PaymentProviderType) error

	// DisableProvider disables a payment provider
	DisableProvider(providerType common.PaymentProviderType) error

	// InitializeDefaultProviders creates default payment provider entries if they don't exist
	InitializeDefaultProviders() error
}
