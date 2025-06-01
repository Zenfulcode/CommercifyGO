package payment

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/zenfulcode/commercify/internal/domain/service"
)

// MockPaymentService implements a mock payment service for testing and development
type MockPaymentService struct{}

// NewMockPaymentService creates a new MockPaymentService
func NewMockPaymentService() *MockPaymentService {
	return &MockPaymentService{}
}

// GetAvailableProviders returns a list of available payment providers
func (s *MockPaymentService) GetAvailableProviders() []service.PaymentProvider {
	return []service.PaymentProvider{
		{
			Type:                service.PaymentProviderMock,
			Name:                "Test Payment",
			Description:         "For testing purposes only",
			Methods:             []service.PaymentMethod{service.PaymentMethodCreditCard},
			Enabled:             true,
			SupportedCurrencies: []string{"USD", "EUR", "GBP", "NOK", "DKK"},
		},
	}
}

// GetAvailableProvidersForCurrency returns a list of available payment providers that support the given currency
func (s *MockPaymentService) GetAvailableProvidersForCurrency(currency string) []service.PaymentProvider {
	providers := s.GetAvailableProviders()
	var supportedProviders []service.PaymentProvider

	for _, provider := range providers {
		for _, supportedCurrency := range provider.SupportedCurrencies {
			if supportedCurrency == currency {
				supportedProviders = append(supportedProviders, provider)
				break
			}
		}
	}

	return supportedProviders
}

// ProcessPayment processes a payment request
func (s *MockPaymentService) ProcessPayment(request service.PaymentRequest) (*service.PaymentResult, error) {
	// Simulate payment processing
	time.Sleep(500 * time.Millisecond)

	// Generate a transaction ID
	transactionID := uuid.New().String()

	// Validate payment details based on method
	switch request.PaymentMethod {
	case service.PaymentMethodCreditCard:
		if request.CardDetails == nil {
			return &service.PaymentResult{
				Success:      false,
				ErrorMessage: "card details are required for credit card payment",
				Provider:     service.PaymentProviderMock,
			}, nil
		}
		// Validate card details
		if request.CardDetails.CardNumber == "" || request.CardDetails.CVV == "" {
			return &service.PaymentResult{
				Success:      false,
				ErrorMessage: "invalid card details",
				Provider:     service.PaymentProviderMock,
			}, nil
		}
	default:
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: "unsupported payment method",
			Provider:     service.PaymentProviderMock,
		}, nil
	}

	// Simulate successful payment
	return &service.PaymentResult{
		Success:       true,
		TransactionID: transactionID,
		Provider:      service.PaymentProviderMock,
	}, nil
}

// VerifyPayment verifies a payment
func (s *MockPaymentService) VerifyPayment(transactionID string, provider service.PaymentProviderType) (bool, error) {
	if transactionID == "" {
		return false, errors.New("transaction ID is required")
	}

	// Simulate verification
	time.Sleep(300 * time.Millisecond)

	// Always return true for mock service
	return true, nil
}

// RefundPayment refunds a payment
func (s *MockPaymentService) RefundPayment(transactionID string, amount int64, provider service.PaymentProviderType) error {
	if transactionID == "" {
		return errors.New("transaction ID is required")
	}
	if amount <= 0 {
		return errors.New("refund amount must be greater than zero")
	}

	// Simulate refund processing
	time.Sleep(500 * time.Millisecond)

	// Always succeed for mock service
	return nil
}

// CapturePayment captures a payment
func (s *MockPaymentService) CapturePayment(transactionID string, amount int64, provider service.PaymentProviderType) error {
	if transactionID == "" {
		return errors.New("transaction ID is required")
	}
	if amount <= 0 {
		return errors.New("capture amount must be greater than zero")
	}

	// Simulate capture processing
	time.Sleep(500 * time.Millisecond)

	// Always succeed for mock service
	return nil
}

// CancelPayment cancels a payment
func (s *MockPaymentService) CancelPayment(transactionID string, provider service.PaymentProviderType) error {
	if transactionID == "" {
		return errors.New("transaction ID is required")
	}

	// Simulate cancellation processing
	time.Sleep(500 * time.Millisecond)

	// Always succeed for mock service
	return nil
}

func (s *MockPaymentService) ForceApprovePayment(transactionID string, phoneNumber string, provider service.PaymentProviderType) error {
	return nil
}
