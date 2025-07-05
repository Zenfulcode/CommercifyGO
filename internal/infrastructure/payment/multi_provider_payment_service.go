package payment

import (
	"fmt"

	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// MultiProviderPaymentService implements payment service with multiple providers
type MultiProviderPaymentService struct {
	providers           map[common.PaymentProviderType]service.PaymentService
	paymentProviderRepo repository.PaymentProviderRepository
	config              *config.Config
	logger              logger.Logger
}

// ProviderWithService represents a provider type with its service implementation
type ProviderWithService struct {
	Type    common.PaymentProviderType
	Service service.PaymentService
}

// NewMultiProviderPaymentService creates a new MultiProviderPaymentService
func NewMultiProviderPaymentService(cfg *config.Config, paymentProviderRepo repository.PaymentProviderRepository, logger logger.Logger) *MultiProviderPaymentService {
	providers := make(map[common.PaymentProviderType]service.PaymentService)

	// Initialize enabled providers
	for _, providerName := range cfg.Payment.EnabledProviders {
		switch providerName {
		case string(common.PaymentProviderStripe):
			if cfg.Stripe.Enabled {
				providers[common.PaymentProviderStripe] = NewStripePaymentService(cfg.Stripe, logger)
				logger.Info("Stripe payment provider initialized")
			}
		case string(common.PaymentProviderMobilePay):
			if cfg.MobilePay.Enabled {
				providers[common.PaymentProviderMobilePay] = NewMobilePayPaymentService(cfg.MobilePay, logger)
				logger.Info("MobilePay payment provider initialized")
			}
		case string(common.PaymentProviderMock):
			providers[common.PaymentProviderMock] = NewMockPaymentService()
			logger.Info("Mock payment provider initialized")
		}
	}

	return &MultiProviderPaymentService{
		providers:           providers,
		paymentProviderRepo: paymentProviderRepo,
		config:              cfg,
		logger:              logger,
	}
}

// GetAvailableProviders returns a list of available payment providers
func (s *MultiProviderPaymentService) GetAvailableProviders() []service.PaymentProvider {
	// Get enabled providers from repository
	providers, err := s.paymentProviderRepo.GetEnabled()
	if err != nil {
		s.logger.Error("Failed to get enabled payment providers: %v", err)
		return []service.PaymentProvider{}
	}

	// Convert entity providers to service providers
	result := make([]service.PaymentProvider, len(providers))
	for i, provider := range providers {
		result[i] = service.PaymentProvider{
			Type:                provider.Type,
			Name:                provider.Name,
			Description:         provider.Description,
			IconURL:             provider.IconURL,
			Methods:             provider.GetMethods(),
			Enabled:             provider.Enabled,
			SupportedCurrencies: provider.SupportedCurrencies,
		}
	}

	return result
}

// GetAvailableProvidersForCurrency returns a list of available payment providers that support the given currency
func (s *MultiProviderPaymentService) GetAvailableProvidersForCurrency(currency string) []service.PaymentProvider {
	// Get enabled providers that support the currency from repository
	providers, err := s.paymentProviderRepo.GetEnabledByCurrency(currency)
	if err != nil {
		s.logger.Error("Failed to get payment providers for currency %s: %v", currency, err)
		return []service.PaymentProvider{}
	}

	// Convert entity providers to service providers
	result := make([]service.PaymentProvider, len(providers))
	for i, provider := range providers {
		result[i] = service.PaymentProvider{
			Type:                provider.Type,
			Name:                provider.Name,
			Description:         provider.Description,
			IconURL:             provider.IconURL,
			Methods:             provider.GetMethods(),
			Enabled:             provider.Enabled,
			SupportedCurrencies: provider.SupportedCurrencies,
		}
	}

	return result
}

// GetProviders returns all configured payment providers
func (s *MultiProviderPaymentService) GetProviders() []ProviderWithService {
	result := make([]ProviderWithService, 0, len(s.providers))
	for providerType, providerService := range s.providers {
		result = append(result, ProviderWithService{
			Type:    providerType,
			Service: providerService,
		})
	}
	return result
}

// ProcessPayment processes a payment request
func (s *MultiProviderPaymentService) ProcessPayment(request service.PaymentRequest) (*service.PaymentResult, error) {
	provider, exists := s.providers[request.PaymentProvider]
	if !exists {
		return nil, fmt.Errorf("payment provider %s not available", request.PaymentProvider)
	}

	result, err := provider.ProcessPayment(request)
	if err != nil {
		s.logger.Error("Error processing payment with provider %s: %v", request.PaymentProvider, err)
		return nil, err
	}

	// Set the provider in the result
	result.Provider = request.PaymentProvider
	return result, nil
}

// VerifyPayment verifies a payment
func (s *MultiProviderPaymentService) VerifyPayment(transactionID string, provider common.PaymentProviderType) (bool, error) {
	paymentProvider, exists := s.providers[provider]
	if !exists {
		return false, fmt.Errorf("payment provider %s not available", provider)
	}

	return paymentProvider.VerifyPayment(transactionID, provider)
}

// RefundPayment refunds a payment
func (s *MultiProviderPaymentService) RefundPayment(transactionID, currency string, amount int64, provider common.PaymentProviderType) (*service.PaymentResult, error) {
	paymentProvider, exists := s.providers[provider]
	if !exists {
		return nil, fmt.Errorf("payment provider %s not available", provider)
	}

	return paymentProvider.RefundPayment(transactionID, currency, amount, provider)
}

// CapturePayment captures a payment
func (s *MultiProviderPaymentService) CapturePayment(transactionID, currency string, amount int64, provider common.PaymentProviderType) (*service.PaymentResult, error) {
	paymentProvider, exists := s.providers[provider]
	if !exists {
		return nil, fmt.Errorf("payment provider %s not available", provider)
	}

	return paymentProvider.CapturePayment(transactionID, currency, amount, provider)
}

// CancelPayment cancels a payment
func (s *MultiProviderPaymentService) CancelPayment(transactionID string, provider common.PaymentProviderType) (*service.PaymentResult, error) {
	paymentProvider, exists := s.providers[provider]
	if !exists {
		return nil, fmt.Errorf("payment provider %s not available", provider)
	}

	return paymentProvider.CancelPayment(transactionID, provider)
}

func (s *MultiProviderPaymentService) ForceApprovePayment(transactionID string, phoneNumber string, provider common.PaymentProviderType) error {
	paymentProvider, exists := s.providers[provider]
	if !exists {
		return fmt.Errorf("payment provider %s not available", provider)
	}

	return paymentProvider.ForceApprovePayment(transactionID, phoneNumber, provider)
}
