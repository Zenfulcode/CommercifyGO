package mock

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// PaymentProviderRepository is a mock implementation of repository.PaymentProviderRepository
type PaymentProviderRepository struct {
	mu        sync.RWMutex
	providers map[uint]*entity.PaymentProvider
	nextID    uint
	// Map provider type to entity for quick lookup
	typeMap map[common.PaymentProviderType]*entity.PaymentProvider
	// For testing error scenarios
	CreateError                        error
	UpdateError                        error
	DeleteError                        error
	GetByIDError                       error
	GetByTypeError                     error
	ListError                          error
	GetEnabledError                    error
	GetEnabledByMethodError            error
	GetEnabledByCurrencyError          error
	GetEnabledByMethodAndCurrencyError error
	UpdateWebhookInfoError             error
	GetWithWebhooksError               error
}

// NewPaymentProviderRepository creates a new mock payment provider repository
func NewPaymentProviderRepository() repository.PaymentProviderRepository {
	return &PaymentProviderRepository{
		providers: make(map[uint]*entity.PaymentProvider),
		typeMap:   make(map[common.PaymentProviderType]*entity.PaymentProvider),
		nextID:    1,
	}
}

// Create creates a new payment provider
func (m *PaymentProviderRepository) Create(provider *entity.PaymentProvider) error {
	if m.CreateError != nil {
		return m.CreateError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if provider type already exists
	if _, exists := m.typeMap[provider.Type]; exists {
		return errors.New("payment provider with this type already exists")
	}

	// Assign ID if not set
	if provider.ID == 0 {
		provider.ID = m.nextID
		m.nextID++
	}

	// Create a copy to avoid external mutations
	providerCopy := *provider
	m.providers[providerCopy.ID] = &providerCopy
	m.typeMap[providerCopy.Type] = &providerCopy

	return nil
}

// Update updates an existing payment provider
func (m *PaymentProviderRepository) Update(provider *entity.PaymentProvider) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	existing, exists := m.providers[provider.ID]
	if !exists {
		return fmt.Errorf("payment provider with ID %d not found", provider.ID)
	}

	// Update type mapping if type changed
	if existing.Type != provider.Type {
		delete(m.typeMap, existing.Type)
		m.typeMap[provider.Type] = provider
	}

	// Update both maps
	providerCopy := *provider
	m.providers[providerCopy.ID] = &providerCopy
	m.typeMap[providerCopy.Type] = &providerCopy

	return nil
}

// Delete deletes a payment provider
func (m *PaymentProviderRepository) Delete(id uint) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	provider, exists := m.providers[id]
	if !exists {
		return fmt.Errorf("payment provider with ID %d not found", id)
	}

	delete(m.providers, id)
	delete(m.typeMap, provider.Type)

	return nil
}

// GetByID returns a payment provider by ID
func (m *PaymentProviderRepository) GetByID(id uint) (*entity.PaymentProvider, error) {
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	provider, exists := m.providers[id]
	if !exists {
		return nil, fmt.Errorf("payment provider with ID %d not found", id)
	}

	// Return a copy
	providerCopy := *provider
	return &providerCopy, nil
}

// GetByType returns a payment provider by type
func (m *PaymentProviderRepository) GetByType(providerType common.PaymentProviderType) (*entity.PaymentProvider, error) {
	if m.GetByTypeError != nil {
		return nil, m.GetByTypeError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	provider, exists := m.typeMap[providerType]
	if !exists {
		return nil, fmt.Errorf("payment provider with type %s not found", providerType)
	}

	// Return a copy
	providerCopy := *provider
	return &providerCopy, nil
}

// List returns all payment providers with pagination
func (m *PaymentProviderRepository) List(offset, limit int) ([]*entity.PaymentProvider, error) {
	if m.ListError != nil {
		return nil, m.ListError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*entity.PaymentProvider
	var count int
	for _, provider := range m.providers {
		if count >= offset {
			if limit > 0 && len(results) >= limit {
				break
			}
			providerCopy := *provider
			results = append(results, &providerCopy)
		}
		count++
	}

	return results, nil
}

// GetEnabled returns all enabled payment providers
func (m *PaymentProviderRepository) GetEnabled() ([]*entity.PaymentProvider, error) {
	if m.GetEnabledError != nil {
		return nil, m.GetEnabledError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*entity.PaymentProvider
	for _, provider := range m.providers {
		if provider.Enabled {
			providerCopy := *provider
			results = append(results, &providerCopy)
		}
	}

	return results, nil
}

// GetEnabledByMethod returns enabled payment providers that support a specific payment method
func (m *PaymentProviderRepository) GetEnabledByMethod(method common.PaymentMethod) ([]*entity.PaymentProvider, error) {
	if m.GetEnabledByMethodError != nil {
		return nil, m.GetEnabledByMethodError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*entity.PaymentProvider
	for _, provider := range m.providers {
		if provider.Enabled && m.supportsMethod(provider, method) {
			providerCopy := *provider
			results = append(results, &providerCopy)
		}
	}

	return results, nil
}

// GetEnabledByCurrency returns enabled payment providers that support a specific currency
func (m *PaymentProviderRepository) GetEnabledByCurrency(currency string) ([]*entity.PaymentProvider, error) {
	if m.GetEnabledByCurrencyError != nil {
		return nil, m.GetEnabledByCurrencyError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*entity.PaymentProvider
	for _, provider := range m.providers {
		if provider.Enabled && m.supportsCurrency(provider, currency) {
			providerCopy := *provider
			results = append(results, &providerCopy)
		}
	}

	return results, nil
}

// GetEnabledByMethodAndCurrency returns enabled payment providers that support both method and currency
func (m *PaymentProviderRepository) GetEnabledByMethodAndCurrency(method common.PaymentMethod, currency string) ([]*entity.PaymentProvider, error) {
	if m.GetEnabledByMethodAndCurrencyError != nil {
		return nil, m.GetEnabledByMethodAndCurrencyError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*entity.PaymentProvider
	for _, provider := range m.providers {
		if provider.Enabled && m.supportsMethod(provider, method) && m.supportsCurrency(provider, currency) {
			providerCopy := *provider
			results = append(results, &providerCopy)
		}
	}

	return results, nil
}

// UpdateWebhookInfo updates webhook information for a payment provider
func (m *PaymentProviderRepository) UpdateWebhookInfo(providerType common.PaymentProviderType, webhookURL, webhookSecret, externalWebhookID string, events []string) error {
	if m.UpdateWebhookInfoError != nil {
		return m.UpdateWebhookInfoError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	provider, exists := m.typeMap[providerType]
	if !exists {
		return fmt.Errorf("payment provider with type %s not found", providerType)
	}

	provider.WebhookURL = webhookURL
	provider.WebhookSecret = webhookSecret
	provider.ExternalWebhookID = externalWebhookID
	provider.WebhookEvents = events

	return nil
}

// GetWithWebhooks returns payment providers that have webhook configurations
func (m *PaymentProviderRepository) GetWithWebhooks() ([]*entity.PaymentProvider, error) {
	if m.GetWithWebhooksError != nil {
		return nil, m.GetWithWebhooksError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*entity.PaymentProvider
	for _, provider := range m.providers {
		if provider.WebhookURL != "" {
			providerCopy := *provider
			results = append(results, &providerCopy)
		}
	}

	return results, nil
}

// Helper methods

// supportsMethod checks if a provider supports a specific payment method
func (m *PaymentProviderRepository) supportsMethod(provider *entity.PaymentProvider, method common.PaymentMethod) bool {
	for _, supportedMethod := range provider.Methods {
		if supportedMethod == method {
			return true
		}
	}
	return false
}

// supportsCurrency checks if a provider supports a specific currency
func (m *PaymentProviderRepository) supportsCurrency(provider *entity.PaymentProvider, currency string) bool {
	// If no currencies specified, assume all are supported
	if len(provider.SupportedCurrencies) == 0 {
		return true
	}

	currency = strings.ToUpper(currency)
	for _, supportedCurrency := range provider.SupportedCurrencies {
		if strings.ToUpper(supportedCurrency) == currency {
			return true
		}
	}
	return false
}

// Helper methods for testing

// SetCreateError sets an error to be returned by Create
func (m *PaymentProviderRepository) SetCreateError(err error) {
	m.CreateError = err
}

// SetUpdateError sets an error to be returned by Update
func (m *PaymentProviderRepository) SetUpdateError(err error) {
	m.UpdateError = err
}

// SetDeleteError sets an error to be returned by Delete
func (m *PaymentProviderRepository) SetDeleteError(err error) {
	m.DeleteError = err
}

// SetGetByIDError sets an error to be returned by GetByID
func (m *PaymentProviderRepository) SetGetByIDError(err error) {
	m.GetByIDError = err
}

// SetGetByTypeError sets an error to be returned by GetByType
func (m *PaymentProviderRepository) SetGetByTypeError(err error) {
	m.GetByTypeError = err
}

// SetListError sets an error to be returned by List
func (m *PaymentProviderRepository) SetListError(err error) {
	m.ListError = err
}

// SetGetEnabledError sets an error to be returned by GetEnabled
func (m *PaymentProviderRepository) SetGetEnabledError(err error) {
	m.GetEnabledError = err
}

// SetGetEnabledByMethodError sets an error to be returned by GetEnabledByMethod
func (m *PaymentProviderRepository) SetGetEnabledByMethodError(err error) {
	m.GetEnabledByMethodError = err
}

// SetGetEnabledByCurrencyError sets an error to be returned by GetEnabledByCurrency
func (m *PaymentProviderRepository) SetGetEnabledByCurrencyError(err error) {
	m.GetEnabledByCurrencyError = err
}

// SetGetEnabledByMethodAndCurrencyError sets an error to be returned by GetEnabledByMethodAndCurrency
func (m *PaymentProviderRepository) SetGetEnabledByMethodAndCurrencyError(err error) {
	m.GetEnabledByMethodAndCurrencyError = err
}

// SetUpdateWebhookInfoError sets an error to be returned by UpdateWebhookInfo
func (m *PaymentProviderRepository) SetUpdateWebhookInfoError(err error) {
	m.UpdateWebhookInfoError = err
}

// SetGetWithWebhooksError sets an error to be returned by GetWithWebhooks
func (m *PaymentProviderRepository) SetGetWithWebhooksError(err error) {
	m.GetWithWebhooksError = err
}

// Reset clears all data and errors
func (m *PaymentProviderRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.providers = make(map[uint]*entity.PaymentProvider)
	m.typeMap = make(map[common.PaymentProviderType]*entity.PaymentProvider)
	m.nextID = 1
	m.CreateError = nil
	m.UpdateError = nil
	m.DeleteError = nil
	m.GetByIDError = nil
	m.GetByTypeError = nil
	m.ListError = nil
	m.GetEnabledError = nil
	m.GetEnabledByMethodError = nil
	m.GetEnabledByCurrencyError = nil
	m.GetEnabledByMethodAndCurrencyError = nil
	m.UpdateWebhookInfoError = nil
	m.GetWithWebhooksError = nil
}

// GetAllProviders returns all providers (for testing purposes)
func (m *PaymentProviderRepository) GetAllProviders() []*entity.PaymentProvider {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*entity.PaymentProvider
	for _, provider := range m.providers {
		providerCopy := *provider
		results = append(results, &providerCopy)
	}

	return results
}

// AddTestProvider adds a test provider (for testing purposes)
func (m *PaymentProviderRepository) AddTestProvider(provider *entity.PaymentProvider) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if provider.ID == 0 {
		provider.ID = m.nextID
		m.nextID++
	}

	providerCopy := *provider
	m.providers[providerCopy.ID] = &providerCopy
	m.typeMap[providerCopy.Type] = &providerCopy
}
