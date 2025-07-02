package mock

import (
	"fmt"
	"sync"

	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/service"
)

// PaymentProviderService is a mock implementation of service.PaymentProviderService
type PaymentProviderService struct {
	mu        sync.RWMutex
	providers map[common.PaymentProviderType]*entity.PaymentProvider
	// For testing error scenarios
	GetPaymentProvidersError            error
	GetEnabledPaymentProvidersError     error
	GetPaymentProvidersForCurrencyError error
	GetPaymentProvidersForMethodError   error
	RegisterWebhookError                error
	DeleteWebhookError                  error
	GetWebhookInfoError                 error
	UpdateProviderConfigurationError    error
	EnableProviderError                 error
	DisableProviderError                error
	InitializeDefaultProvidersError     error
}

// NewPaymentProviderService creates a new mock payment provider service
func NewPaymentProviderService() service.PaymentProviderService {
	return &PaymentProviderService{
		providers: make(map[common.PaymentProviderType]*entity.PaymentProvider),
	}
}

// GetPaymentProviders returns all payment providers
func (m *PaymentProviderService) GetPaymentProviders() ([]service.PaymentProvider, error) {
	if m.GetPaymentProvidersError != nil {
		return nil, m.GetPaymentProvidersError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []service.PaymentProvider
	for _, provider := range m.providers {
		result = append(result, service.PaymentProvider{
			Type:                provider.Type,
			Name:                provider.Name,
			Description:         provider.Description,
			IconURL:             provider.IconURL,
			Methods:             provider.Methods,
			Enabled:             provider.Enabled,
			SupportedCurrencies: provider.SupportedCurrencies,
		})
	}

	return result, nil
}

// GetEnabledPaymentProviders returns all enabled payment providers
func (m *PaymentProviderService) GetEnabledPaymentProviders() ([]service.PaymentProvider, error) {
	if m.GetEnabledPaymentProvidersError != nil {
		return nil, m.GetEnabledPaymentProvidersError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []service.PaymentProvider
	for _, provider := range m.providers {
		if provider.Enabled {
			result = append(result, service.PaymentProvider{
				Type:                provider.Type,
				Name:                provider.Name,
				Description:         provider.Description,
				IconURL:             provider.IconURL,
				Methods:             provider.Methods,
				Enabled:             provider.Enabled,
				SupportedCurrencies: provider.SupportedCurrencies,
			})
		}
	}

	return result, nil
}

// GetPaymentProvidersForCurrency returns payment providers that support a specific currency
func (m *PaymentProviderService) GetPaymentProvidersForCurrency(currency string) ([]service.PaymentProvider, error) {
	if m.GetPaymentProvidersForCurrencyError != nil {
		return nil, m.GetPaymentProvidersForCurrencyError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []service.PaymentProvider
	for _, provider := range m.providers {
		if provider.Enabled && m.supportsCurrency(provider, currency) {
			result = append(result, service.PaymentProvider{
				Type:                provider.Type,
				Name:                provider.Name,
				Description:         provider.Description,
				IconURL:             provider.IconURL,
				Methods:             provider.Methods,
				Enabled:             provider.Enabled,
				SupportedCurrencies: provider.SupportedCurrencies,
			})
		}
	}

	return result, nil
}

// GetPaymentProvidersForMethod returns payment providers that support a specific payment method
func (m *PaymentProviderService) GetPaymentProvidersForMethod(method common.PaymentMethod) ([]service.PaymentProvider, error) {
	if m.GetPaymentProvidersForMethodError != nil {
		return nil, m.GetPaymentProvidersForMethodError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var result []service.PaymentProvider
	for _, provider := range m.providers {
		if provider.Enabled && m.supportsMethod(provider, method) {
			result = append(result, service.PaymentProvider{
				Type:                provider.Type,
				Name:                provider.Name,
				Description:         provider.Description,
				IconURL:             provider.IconURL,
				Methods:             provider.Methods,
				Enabled:             provider.Enabled,
				SupportedCurrencies: provider.SupportedCurrencies,
			})
		}
	}

	return result, nil
}

// RegisterWebhook registers a webhook for a payment provider
func (m *PaymentProviderService) RegisterWebhook(providerType common.PaymentProviderType, webhookURL string, events []string) error {
	if m.RegisterWebhookError != nil {
		return m.RegisterWebhookError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	provider, exists := m.providers[providerType]
	if !exists {
		return fmt.Errorf("payment provider %s not found", providerType)
	}

	provider.WebhookURL = webhookURL
	provider.WebhookEvents = events

	return nil
}

// DeleteWebhook deletes a webhook for a payment provider
func (m *PaymentProviderService) DeleteWebhook(providerType common.PaymentProviderType) error {
	if m.DeleteWebhookError != nil {
		return m.DeleteWebhookError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	provider, exists := m.providers[providerType]
	if !exists {
		return fmt.Errorf("payment provider %s not found", providerType)
	}

	provider.WebhookURL = ""
	provider.WebhookSecret = ""
	provider.ExternalWebhookID = ""
	provider.WebhookEvents = nil

	return nil
}

// GetWebhookInfo returns webhook information for a payment provider
func (m *PaymentProviderService) GetWebhookInfo(providerType common.PaymentProviderType) (*entity.PaymentProvider, error) {
	if m.GetWebhookInfoError != nil {
		return nil, m.GetWebhookInfoError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	provider, exists := m.providers[providerType]
	if !exists {
		return nil, fmt.Errorf("payment provider %s not found", providerType)
	}

	// Return a copy
	providerCopy := *provider
	return &providerCopy, nil
}

// UpdateProviderConfiguration updates the configuration for a payment provider
func (m *PaymentProviderService) UpdateProviderConfiguration(providerType common.PaymentProviderType, config map[string]interface{}) error {
	if m.UpdateProviderConfigurationError != nil {
		return m.UpdateProviderConfigurationError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	provider, exists := m.providers[providerType]
	if !exists {
		return fmt.Errorf("payment provider %s not found", providerType)
	}

	provider.Configuration = config

	return nil
}

// EnableProvider enables a payment provider
func (m *PaymentProviderService) EnableProvider(providerType common.PaymentProviderType) error {
	if m.EnableProviderError != nil {
		return m.EnableProviderError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	provider, exists := m.providers[providerType]
	if !exists {
		return fmt.Errorf("payment provider %s not found", providerType)
	}

	provider.Enabled = true

	return nil
}

// DisableProvider disables a payment provider
func (m *PaymentProviderService) DisableProvider(providerType common.PaymentProviderType) error {
	if m.DisableProviderError != nil {
		return m.DisableProviderError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	provider, exists := m.providers[providerType]
	if !exists {
		return fmt.Errorf("payment provider %s not found", providerType)
	}

	provider.Enabled = false

	return nil
}

// InitializeDefaultProviders creates default payment provider entries if they don't exist
func (m *PaymentProviderService) InitializeDefaultProviders() error {
	if m.InitializeDefaultProvidersError != nil {
		return m.InitializeDefaultProvidersError
	}

	// Mock implementation - just return success
	return nil
}

// Helper methods

// supportsCurrency checks if a provider supports a specific currency
func (m *PaymentProviderService) supportsCurrency(provider *entity.PaymentProvider, currency string) bool {
	if len(provider.SupportedCurrencies) == 0 {
		return true // If no currencies specified, assume all are supported
	}

	for _, supportedCurrency := range provider.SupportedCurrencies {
		if supportedCurrency == currency {
			return true
		}
	}
	return false
}

// supportsMethod checks if a provider supports a specific payment method
func (m *PaymentProviderService) supportsMethod(provider *entity.PaymentProvider, method common.PaymentMethod) bool {
	for _, supportedMethod := range provider.Methods {
		if supportedMethod == method {
			return true
		}
	}
	return false
}

// Helper methods for testing

// SetGetPaymentProvidersError sets an error to be returned by GetPaymentProviders
func (m *PaymentProviderService) SetGetPaymentProvidersError(err error) {
	m.GetPaymentProvidersError = err
}

// SetGetEnabledPaymentProvidersError sets an error to be returned by GetEnabledPaymentProviders
func (m *PaymentProviderService) SetGetEnabledPaymentProvidersError(err error) {
	m.GetEnabledPaymentProvidersError = err
}

// SetGetPaymentProvidersForCurrencyError sets an error to be returned by GetPaymentProvidersForCurrency
func (m *PaymentProviderService) SetGetPaymentProvidersForCurrencyError(err error) {
	m.GetPaymentProvidersForCurrencyError = err
}

// SetGetPaymentProvidersForMethodError sets an error to be returned by GetPaymentProvidersForMethod
func (m *PaymentProviderService) SetGetPaymentProvidersForMethodError(err error) {
	m.GetPaymentProvidersForMethodError = err
}

// SetRegisterWebhookError sets an error to be returned by RegisterWebhook
func (m *PaymentProviderService) SetRegisterWebhookError(err error) {
	m.RegisterWebhookError = err
}

// SetDeleteWebhookError sets an error to be returned by DeleteWebhook
func (m *PaymentProviderService) SetDeleteWebhookError(err error) {
	m.DeleteWebhookError = err
}

// SetGetWebhookInfoError sets an error to be returned by GetWebhookInfo
func (m *PaymentProviderService) SetGetWebhookInfoError(err error) {
	m.GetWebhookInfoError = err
}

// SetUpdateProviderConfigurationError sets an error to be returned by UpdateProviderConfiguration
func (m *PaymentProviderService) SetUpdateProviderConfigurationError(err error) {
	m.UpdateProviderConfigurationError = err
}

// SetEnableProviderError sets an error to be returned by EnableProvider
func (m *PaymentProviderService) SetEnableProviderError(err error) {
	m.EnableProviderError = err
}

// SetDisableProviderError sets an error to be returned by DisableProvider
func (m *PaymentProviderService) SetDisableProviderError(err error) {
	m.DisableProviderError = err
}

// SetInitializeDefaultProvidersError sets an error to be returned by InitializeDefaultProviders
func (m *PaymentProviderService) SetInitializeDefaultProvidersError(err error) {
	m.InitializeDefaultProvidersError = err
}

// Reset clears all data and errors
func (m *PaymentProviderService) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.providers = make(map[common.PaymentProviderType]*entity.PaymentProvider)
	m.GetPaymentProvidersError = nil
	m.GetEnabledPaymentProvidersError = nil
	m.GetPaymentProvidersForCurrencyError = nil
	m.GetPaymentProvidersForMethodError = nil
	m.RegisterWebhookError = nil
	m.DeleteWebhookError = nil
	m.GetWebhookInfoError = nil
	m.UpdateProviderConfigurationError = nil
	m.EnableProviderError = nil
	m.DisableProviderError = nil
	m.InitializeDefaultProvidersError = nil
}

// AddTestProvider adds a test provider (for testing purposes)
func (m *PaymentProviderService) AddTestProvider(provider *entity.PaymentProvider) {
	m.mu.Lock()
	defer m.mu.Unlock()

	providerCopy := *provider
	m.providers[providerCopy.Type] = &providerCopy
}

// GetTestProvider gets a test provider by type (for testing purposes)
func (m *PaymentProviderService) GetTestProvider(providerType common.PaymentProviderType) *entity.PaymentProvider {
	m.mu.RLock()
	defer m.mu.RUnlock()

	provider, exists := m.providers[providerType]
	if !exists {
		return nil
	}

	providerCopy := *provider
	return &providerCopy
}
