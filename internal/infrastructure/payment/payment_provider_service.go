package payment

import (
	"fmt"

	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// PaymentProviderServiceImpl implements service.PaymentProviderService
type PaymentProviderServiceImpl struct {
	repo             repository.PaymentProviderRepository
	config           *config.Config
	logger           logger.Logger
	mobilePayService *MobilePayPaymentService
}

// NewPaymentProviderService creates a new PaymentProviderServiceImpl
func NewPaymentProviderService(repo repository.PaymentProviderRepository, cfg *config.Config, logger logger.Logger) service.PaymentProviderService {
	// Create MobilePay service if enabled
	var mobilePayService *MobilePayPaymentService
	if cfg.MobilePay.Enabled {
		mobilePayService = NewMobilePayPaymentService(cfg.MobilePay, logger)
	}

	return &PaymentProviderServiceImpl{
		repo:             repo,
		config:           cfg,
		logger:           logger,
		mobilePayService: mobilePayService,
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

	// Get the provider to check its configuration
	provider, err := s.repo.GetByType(providerType)
	if err != nil {
		return fmt.Errorf("failed to get provider %s: %w", providerType, err)
	}

	// Handle MobilePay webhook registration using the dedicated MobilePay service
	if providerType == common.PaymentProviderMobilePay {
		if s.mobilePayService == nil {
			return fmt.Errorf("MobilePay service not initialized")
		}

		// Register webhook using MobilePay service
		if err := s.mobilePayService.RegisterWebhook(provider, webhookURL); err != nil {
			return fmt.Errorf("failed to register MobilePay webhook via API: %w", err)
		}

		// Update the provider in the database
		if err := s.repo.Update(provider); err != nil {
			return fmt.Errorf("failed to update provider in database: %w", err)
		}

		s.logger.Info("Successfully registered MobilePay webhook via API: %s", webhookURL)
		return nil
	}

	// For other providers, use mock implementation
	err = s.repo.UpdateWebhookInfo(providerType, provider.WebhookURL, provider.WebhookSecret, provider.ExternalWebhookID, provider.WebhookEvents)
	if err != nil {
		return fmt.Errorf("failed to register webhook for provider %s: %w", providerType, err)
	}

	s.logger.Info("Successfully registered webhook for provider %s: %s", providerType, webhookURL)
	return nil
}

// DeleteWebhook implements service.PaymentProviderService.
func (s *PaymentProviderServiceImpl) DeleteWebhook(providerType common.PaymentProviderType) error {
	// Handle MobilePay webhook deletion using the dedicated MobilePay service
	if providerType == common.PaymentProviderMobilePay {
		if s.mobilePayService == nil {
			return fmt.Errorf("MobilePay service not initialized")
		}

		provider, err := s.repo.GetByType(providerType)
		if err != nil {
			return fmt.Errorf("failed to get MobilePay provider: %w", err)
		}

		// If there's an external webhook ID, delete it via API
		if provider.ExternalWebhookID != "" {
			if err := s.mobilePayService.DeleteWebhook(provider); err != nil {
				s.logger.Error("Failed to delete MobilePay webhook via API: %v", err)
				// Continue with database cleanup even if API call fails
			}
		}
	}

	// Update database to remove webhook info
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
func (s *PaymentProviderServiceImpl) UpdateProviderConfiguration(providerType common.PaymentProviderType, config common.JSONB) error {
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
			Enabled:     s.config.Stripe.Enabled,
			SupportedCurrencies: []string{
				"USD", "EUR", "GBP", "JPY", "CAD", "AUD", "CHF", "SEK", "NOK", "DKK",
				"PLN", "CZK", "HUF", "BGN", "RON", "HRK", "ISK", "MXN", "BRL", "SGD",
				"HKD", "INR", "MYR", "PHP", "THB", "TWD", "KRW", "NZD", "ILS", "ZAR",
			},
			Configuration: common.JSONB{
				"SecretKey":          s.config.Stripe.SecretKey,
				"PublicKey":          s.config.Stripe.PublicKey,
				"WebhookSecret":      s.config.Stripe.WebhookSecret,
				"PaymentDescription": s.config.Stripe.PaymentDescription,
				"ReturnURL":          s.config.Stripe.ReturnURL,
				"Enabled":            s.config.Stripe.Enabled,
			},
			Priority: 100,
		},
		{
			Type:                common.PaymentProviderMobilePay,
			Name:                "MobilePay",
			Description:         "Pay with MobilePay app",
			Methods:             []common.PaymentMethod{common.PaymentMethodWallet},
			Enabled:             s.config.MobilePay.Enabled,
			SupportedCurrencies: []string{"NOK", "DKK", "EUR"},
			Configuration: common.JSONB{
				"MerchantSerialNumber": s.config.MobilePay.MerchantSerialNumber,
				"SubscriptionKey":      s.config.MobilePay.SubscriptionKey,
				"ClientID":             s.config.MobilePay.ClientID,
				"ClientSecret":         s.config.MobilePay.ClientSecret,
				"ReturnURL":            s.config.MobilePay.ReturnURL,
				"WebhookURL":           s.config.MobilePay.WebhookURL,
				"PaymentDescription":   s.config.MobilePay.PaymentDescription,
				"Market":               s.config.MobilePay.Market,
				"Enabled":              s.config.MobilePay.Enabled,
				"IsTestMode":           s.config.MobilePay.IsTestMode,
			},
			Priority: 90,
		},
		{
			Type:                common.PaymentProviderMock,
			Name:                "Test Payment",
			Description:         "For testing purposes only",
			Methods:             []common.PaymentMethod{common.PaymentMethodCreditCard},
			Enabled:             true, // Always enabled for testing
			SupportedCurrencies: []string{"USD", "EUR", "GBP", "NOK", "DKK"},
			Configuration: common.JSONB{
				"Enabled":            true,
				"IsTestMode":         true,
				"PaymentDescription": "Test payment for development",
				"AutoConfirm":        true,
			},
			Priority: 10,
		},
	}

	createdCount := 0
	existingCount := 0

	// Create providers if they don't exist
	for _, provider := range defaultProviders {
		existingProvider, err := s.repo.GetByType(provider.Type)
		if err != nil {
			// Provider doesn't exist, create it
			if err := s.repo.Create(provider); err != nil {
				s.logger.Error("Failed to create default provider %s: %v", provider.Type, err)
				return fmt.Errorf("failed to create default provider %s: %w", provider.Type, err)
			}
			s.logger.Info("Created default provider: %s", provider.Type)
			createdCount++

			// Register webhook for MobilePay if enabled
			if provider.Type == common.PaymentProviderMobilePay && provider.Enabled && s.mobilePayService != nil {
				webhookURL, _ := provider.Configuration["WebhookURL"].(string)
				if err := s.mobilePayService.RegisterWebhook(provider, webhookURL); err != nil {
					s.logger.Error("Failed to register MobilePay webhook during initialization: %v", err)
					// Don't fail the entire initialization if webhook registration fails
				} else {
					// Update the provider in the database
					if err := s.repo.Update(provider); err != nil {
						s.logger.Error("Failed to update provider with webhook info: %v", err)
					}
				}
			}
		} else {
			s.logger.Debug("Provider %s already exists, skipping creation", provider.Type)
			existingCount++

			// Check if MobilePay webhook needs to be registered for existing provider
			if existingProvider.Type == common.PaymentProviderMobilePay && existingProvider.Enabled &&
				(existingProvider.WebhookSecret == "" || existingProvider.ExternalWebhookID == "") && s.mobilePayService != nil {
				s.logger.Info("Registering webhook for existing MobilePay provider (missing webhook data)")
				webhookURL, _ := existingProvider.Configuration["WebhookURL"].(string)
				if err := s.mobilePayService.RegisterWebhook(existingProvider, webhookURL); err != nil {
					s.logger.Error("Failed to register MobilePay webhook for existing provider: %v", err)
					// Don't fail the entire initialization if webhook registration fails
				} else {
					// Update the provider in the database
					if err := s.repo.Update(existingProvider); err != nil {
						s.logger.Error("Failed to update existing provider with webhook info: %v", err)
					}
				}
			}
		}
	}

	s.logger.Info("Default provider initialization complete. Created: %d, Existing: %d, Total: %d",
		createdCount, existingCount, len(defaultProviders))

	return nil
}
