package payment

import (
	"fmt"
	"slices"

	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// MultiProviderPaymentService implements payment service with multiple providers
type MultiProviderPaymentService struct {
	providers map[service.PaymentProviderType]service.PaymentService
	config    *config.Config
	logger    logger.Logger
}

// ProviderWithService represents a provider type with its service implementation
type ProviderWithService struct {
	Type    service.PaymentProviderType
	Service service.PaymentService
}

// NewMultiProviderPaymentService creates a new MultiProviderPaymentService
func NewMultiProviderPaymentService(cfg *config.Config, logger logger.Logger) *MultiProviderPaymentService {
	providers := make(map[service.PaymentProviderType]service.PaymentService)

	// Initialize enabled providers
	for _, providerName := range cfg.Payment.EnabledProviders {
		switch providerName {
		case string(service.PaymentProviderStripe):
			if cfg.Stripe.Enabled {
				providers[service.PaymentProviderStripe] = NewStripePaymentService(cfg.Stripe, logger)
				logger.Info("Stripe payment provider initialized")
			}
		case string(service.PaymentProviderMobilePay):
			if cfg.MobilePay.Enabled {
				providers[service.PaymentProviderMobilePay] = NewMobilePayPaymentService(cfg.MobilePay, logger)
				logger.Info("MobilePay payment provider initialized")
			}
		case string(service.PaymentProviderMock):
			providers[service.PaymentProviderMock] = NewMockPaymentService()
			logger.Info("Mock payment provider initialized")
		}
	}

	return &MultiProviderPaymentService{
		providers: providers,
		config:    cfg,
		logger:    logger,
	}
}

// GetAvailableProviders returns a list of available payment providers
func (s *MultiProviderPaymentService) GetAvailableProviders() []service.PaymentProvider {
	var enabledProviders []service.PaymentProvider

	// Collect providers from all enabled payment services
	for _, providerService := range s.providers {
		providers := providerService.GetAvailableProviders()
		enabledProviders = append(enabledProviders, providers...)
	}

	return enabledProviders
}

// GetAvailableProvidersForCurrency returns a list of available payment providers that support the given currency
func (s *MultiProviderPaymentService) GetAvailableProvidersForCurrency(currency string) []service.PaymentProvider {
	var supportedProviders []service.PaymentProvider

	// Collect providers from all enabled payment services that support the currency
	for _, providerService := range s.providers {
		providers := providerService.GetAvailableProvidersForCurrency(currency)
		supportedProviders = append(supportedProviders, providers...)
	}

	return supportedProviders
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	return slices.Contains(slice, item)
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
func (s *MultiProviderPaymentService) VerifyPayment(transactionID string, provider service.PaymentProviderType) (bool, error) {
	paymentProvider, exists := s.providers[provider]
	if !exists {
		return false, fmt.Errorf("payment provider %s not available", provider)
	}

	return paymentProvider.VerifyPayment(transactionID, provider)
}

// RefundPayment refunds a payment
func (s *MultiProviderPaymentService) RefundPayment(transactionID, currency string, amount int64, provider service.PaymentProviderType) (*service.PaymentResult, error) {
	paymentProvider, exists := s.providers[provider]
	if !exists {
		return nil, fmt.Errorf("payment provider %s not available", provider)
	}

	return paymentProvider.RefundPayment(transactionID, currency, amount, provider)
}

// CapturePayment captures a payment
func (s *MultiProviderPaymentService) CapturePayment(transactionID, currency string, amount int64, provider service.PaymentProviderType) (*service.PaymentResult, error) {
	paymentProvider, exists := s.providers[provider]
	if !exists {
		return nil, fmt.Errorf("payment provider %s not available", provider)
	}

	return paymentProvider.CapturePayment(transactionID, currency, amount, provider)
}

// CancelPayment cancels a payment
func (s *MultiProviderPaymentService) CancelPayment(transactionID string, provider service.PaymentProviderType) (*service.PaymentResult, error) {
	paymentProvider, exists := s.providers[provider]
	if !exists {
		return nil, fmt.Errorf("payment provider %s not available", provider)
	}

	return paymentProvider.CancelPayment(transactionID, provider)
}

func (s *MultiProviderPaymentService) ForceApprovePayment(transactionID string, phoneNumber string, provider service.PaymentProviderType) error {
	paymentProvider, exists := s.providers[provider]
	if !exists {
		return fmt.Errorf("payment provider %s not available", provider)
	}

	return paymentProvider.ForceApprovePayment(transactionID, phoneNumber, provider)
}
