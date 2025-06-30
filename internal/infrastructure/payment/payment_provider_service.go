package payment

import (
	"fmt"

	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// PaymentProviderServiceImpl implements service.PaymentProviderService
type PaymentProviderServiceImpl struct {
	repo   repository.PaymentProviderRepository
	logger logger.Logger
}

// NewPaymentProviderService creates a new PaymentProviderServiceImpl
func NewPaymentProviderService(repo repository.PaymentProviderRepository, logger logger.Logger) service.PaymentProviderService {
	return &PaymentProviderServiceImpl{
		repo:   repo,
		logger: logger,
	}
}

// convertToServiceProvider converts entity.PaymentProvider to service.PaymentProvider
func (s *PaymentProviderServiceImpl) convertToServiceProvider(provider *entity.PaymentProvider) service.PaymentProvider {
	return service.PaymentProvider{
		Type:                provider.Type,
		Name:                provider.Name,
		Description:         provider.Description,
		IconURL:             provider.IconURL,
		Methods:             provider.Methods,
		Enabled:             provider.Enabled,
		SupportedCurrencies: provider.SupportedCurrencies,
	}
}

// convertToServiceProviders converts a slice of entity.PaymentProvider to service.PaymentProvider
func (s *PaymentProviderServiceImpl) convertToServiceProviders(providers []*entity.PaymentProvider) []service.PaymentProvider {
	result := make([]service.PaymentProvider, len(providers))
	for i, provider := range providers {
		result[i] = s.convertToServiceProvider(provider)
	}
	return result
}

// GetPaymentProviders implements service.PaymentProviderService.
func (s *PaymentProviderServiceImpl) GetPaymentProviders() ([]service.PaymentProvider, error) {
	providers, err := s.repo.List(0, 0) // Get all providers
	if err != nil {
		return nil, fmt.Errorf("failed to list payment providers: %w", err)
	}

	return s.convertToServiceProviders(providers), nil
}

// GetEnabledPaymentProviders implements service.PaymentProviderService.
func (s *PaymentProviderServiceImpl) GetEnabledPaymentProviders() ([]service.PaymentProvider, error) {
	providers, err := s.repo.GetEnabled()
	if err != nil {
		return nil, fmt.Errorf("failed to get enabled payment providers: %w", err)
	}

	return s.convertToServiceProviders(providers), nil
}

// GetPaymentProvidersForCurrency implements service.PaymentProviderService.
func (s *PaymentProviderServiceImpl) GetPaymentProvidersForCurrency(currency string) ([]service.PaymentProvider, error) {
	providers, err := s.repo.GetEnabledByCurrency(currency)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment providers for currency %s: %w", currency, err)
	}

	return s.convertToServiceProviders(providers), nil
}

// GetPaymentProvidersForMethod implements service.PaymentProviderService.
func (s *PaymentProviderServiceImpl) GetPaymentProvidersForMethod(method common.PaymentMethod) ([]service.PaymentProvider, error) {
	providers, err := s.repo.GetEnabledByMethod(method)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment providers for method %s: %w", method, err)
	}

	return s.convertToServiceProviders(providers), nil
}

// RegisterWebhook implements service.PaymentProviderService.
func (s *PaymentProviderServiceImpl) RegisterWebhook(providerType common.PaymentProviderType, webhookURL string, events []string) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL cannot be empty")
	}

	// Generate a mock external webhook ID for now
	// In a real implementation, this would come from the payment provider's API
	externalWebhookID := fmt.Sprintf("webhook_%s_%d", providerType, len(events))

	err := s.repo.UpdateWebhookInfo(providerType, webhookURL, "", externalWebhookID, events)
	if err != nil {
		return fmt.Errorf("failed to register webhook for provider %s: %w", providerType, err)
	}

	s.logger.Info("Successfully registered webhook for provider %s: %s", providerType, webhookURL)
	return nil
}

// DeleteWebhook implements service.PaymentProviderService.
func (s *PaymentProviderServiceImpl) DeleteWebhook(providerType common.PaymentProviderType) error {
	err := s.repo.UpdateWebhookInfo(providerType, "", "", "", nil)
	if err != nil {
		return fmt.Errorf("failed to delete webhook for provider %s: %w", providerType, err)
	}

	s.logger.Info("Successfully deleted webhook for provider %s", providerType)
	return nil
}

// GetWebhookInfo implements service.PaymentProviderService.
func (s *PaymentProviderServiceImpl) GetWebhookInfo(providerType common.PaymentProviderType) (*entity.PaymentProvider, error) {
	provider, err := s.repo.GetByType(providerType)
	if err != nil {
		return nil, fmt.Errorf("failed to get webhook info for provider %s: %w", providerType, err)
	}

	return provider, nil
}

// UpdateProviderConfiguration implements service.PaymentProviderService.
func (s *PaymentProviderServiceImpl) UpdateProviderConfiguration(providerType common.PaymentProviderType, config map[string]any) error {
	if config == nil {
		return fmt.Errorf("configuration cannot be nil")
	}

	provider, err := s.repo.GetByType(providerType)
	if err != nil {
		return fmt.Errorf("failed to get provider %s: %w", providerType, err)
	}

	provider.SetConfiguration(config)

	err = s.repo.Update(provider)
	if err != nil {
		return fmt.Errorf("failed to update configuration for provider %s: %w", providerType, err)
	}

	s.logger.Info("Successfully updated configuration for provider %s", providerType)
	return nil
}

// EnableProvider implements service.PaymentProviderService.
func (s *PaymentProviderServiceImpl) EnableProvider(providerType common.PaymentProviderType) error {
	provider, err := s.repo.GetByType(providerType)
	if err != nil {
		return fmt.Errorf("failed to get provider %s: %w", providerType, err)
	}

	if provider.Enabled {
		s.logger.Info("Provider %s is already enabled", providerType)
		return nil
	}

	provider.Enabled = true

	err = s.repo.Update(provider)
	if err != nil {
		return fmt.Errorf("failed to enable provider %s: %w", providerType, err)
	}

	s.logger.Info("Successfully enabled provider %s", providerType)
	return nil
}

// DisableProvider implements service.PaymentProviderService.
func (s *PaymentProviderServiceImpl) DisableProvider(providerType common.PaymentProviderType) error {
	provider, err := s.repo.GetByType(providerType)
	if err != nil {
		return fmt.Errorf("failed to get provider %s: %w", providerType, err)
	}

	if !provider.Enabled {
		s.logger.Info("Provider %s is already disabled", providerType)
		return nil
	}

	provider.Enabled = false

	err = s.repo.Update(provider)
	if err != nil {
		return fmt.Errorf("failed to disable provider %s: %w", providerType, err)
	}

	s.logger.Info("Successfully disabled provider %s", providerType)
	return nil
}

// InitializeDefaultProviders implements service.PaymentProviderService.
func (s *PaymentProviderServiceImpl) InitializeDefaultProviders() error {
	s.logger.Info("Initializing default payment providers...")

	// Define default providers
	defaultProviders := []*entity.PaymentProvider{
		{
			Type:        common.PaymentProviderStripe,
			Name:        "Stripe",
			Description: "Pay with credit or debit card",
			Methods:     []common.PaymentMethod{common.PaymentMethodCreditCard},
			Enabled:     false, // Disabled by default until configured
			SupportedCurrencies: []string{
				"USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "SEK", "NOK", "DKK",
				"PLN", "CZK", "HUF", "BGN", "RON", "HRK", "ISK", "MXN", "BRL", "SGD",
				"HKD", "INR", "MYR", "PHP", "THB", "TWD", "KRW", "NZD", "ILS", "ZAR",
			},
			Priority: 100,
		},
		{
			Type:                common.PaymentProviderMobilePay,
			Name:                "MobilePay",
			Description:         "Pay with MobilePay app",
			Methods:             []common.PaymentMethod{common.PaymentMethodWallet},
			Enabled:             false, // Disabled by default until configured
			SupportedCurrencies: []string{"NOK", "DKK", "EUR"},
			Priority:            90,
		},
		{
			Type:                common.PaymentProviderMock,
			Name:                "Test Payment",
			Description:         "For testing purposes only",
			Methods:             []common.PaymentMethod{common.PaymentMethodCreditCard},
			Enabled:             true, // Enabled by default for testing
			SupportedCurrencies: []string{"USD", "EUR", "GBP", "NOK", "DKK"},
			Priority:            10,
		},
	}

	createdCount := 0
	existingCount := 0

	// Create providers if they don't exist
	for _, provider := range defaultProviders {
		_, err := s.repo.GetByType(provider.Type)
		if err != nil {
			// Provider doesn't exist, create it
			if err := s.repo.Create(provider); err != nil {
				s.logger.Error("Failed to create default provider %s: %v", provider.Type, err)
				return fmt.Errorf("failed to create default provider %s: %w", provider.Type, err)
			}
			s.logger.Info("Created default provider: %s", provider.Type)
			createdCount++
		} else {
			s.logger.Debug("Provider %s already exists, skipping creation", provider.Type)
			existingCount++
		}
	}

	s.logger.Info("Default provider initialization complete. Created: %d, Existing: %d, Total: %d",
		createdCount, existingCount, len(defaultProviders))

	return nil
}
