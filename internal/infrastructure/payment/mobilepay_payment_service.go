package payment

import (
	"errors"
	"fmt"
	"regexp"
	"slices"

	"github.com/gkhaavik/vipps-mobilepay-sdk/pkg/client"
	"github.com/gkhaavik/vipps-mobilepay-sdk/pkg/models"
	"github.com/google/uuid"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// MobilePayPaymentService implements a MobilePay payment service
type MobilePayPaymentService struct {
	vippsClient   *client.Client
	webhookClient *client.Webhook
	epayment      *client.Payment
	logger        logger.Logger
	config        config.MobilePayConfig
}

// NewMobilePayPaymentService creates a new MobilePayPaymentService
func NewMobilePayPaymentService(config config.MobilePayConfig, logger logger.Logger) *MobilePayPaymentService {
	vippsClient := client.NewClient(
		config.ClientID,
		config.ClientSecret,
		config.SubscriptionKey,
		config.MerchantSerialNumber,
		config.IsTestMode)

	paymentClient := client.NewPayment(vippsClient)
	webhookClient := client.NewWebhook(vippsClient)

	return &MobilePayPaymentService{
		vippsClient:   vippsClient,
		webhookClient: webhookClient,
		epayment:      paymentClient,
		logger:        logger,
		config:        config,
	}
}

// GetAvailableProviders returns a list of available payment providers
func (s *MobilePayPaymentService) GetAvailableProviders() []service.PaymentProvider {
	return []service.PaymentProvider{
		{
			Type:                common.PaymentProviderMobilePay,
			Name:                "MobilePay",
			Description:         "Pay with MobilePay app",
			IconURL:             "/assets/images/mobilepay-logo.png",
			Methods:             []common.PaymentMethod{common.PaymentMethodWallet},
			Enabled:             true,
			SupportedCurrencies: []string{"NOK", "DKK", "EUR"},
		},
	}
}

// GetAvailableProvidersForCurrency returns a list of available payment providers that support the given currency
func (s *MobilePayPaymentService) GetAvailableProvidersForCurrency(currency string) []service.PaymentProvider {
	providers := s.GetAvailableProviders()
	var supportedProviders []service.PaymentProvider

	for _, provider := range providers {
		if slices.Contains(provider.SupportedCurrencies, currency) {
			supportedProviders = append(supportedProviders, provider)
		}
	}

	return supportedProviders
}

// ProcessPayment processes a payment request using MobilePay
func (s *MobilePayPaymentService) ProcessPayment(request service.PaymentRequest) (*service.PaymentResult, error) {
	if !slices.Contains(s.GetAvailableProviders()[0].SupportedCurrencies, request.Currency) {
		return nil, fmt.Errorf("currency %s is not supported by MobilePay", request.Currency)
	}

	if request.PhoneNumber == "" {
		return nil, errors.New("phone number is required for MobilePay payments")
	}

	phoneNumber := request.PhoneNumber

	r := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

	if !r.MatchString(phoneNumber) {
		return nil, fmt.Errorf("invalid phone number format: %s", phoneNumber)
	}

	// Generate a unique reference for this payment
	reference := fmt.Sprintf("order-%d-%s", request.OrderID, uuid.New().String())

	// Construct the payment request
	paymentRequest := models.CreatePaymentRequest{
		Amount: models.Amount{
			Currency: request.Currency,
			Value:    int(request.Amount),
		},
		Customer: &models.Customer{
			PhoneNumber: &phoneNumber,
		},
		PaymentMethod: &models.PaymentMethod{
			Type: "WALLET",
		},
		Reference:          reference,
		ReturnURL:          s.config.ReturnURL + "?order=" + request.OrderNumber,
		UserFlow:           models.UserFlowWebRedirect,
		PaymentDescription: s.config.PaymentDescription,
	}

	res, err := s.epayment.Create(paymentRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to create MobilePay payment: %v", err)
	}

	// MobilePay requires a redirect to complete the payment
	// Return a result with action URL
	return &service.PaymentResult{
		Success:        false,
		TransactionID:  res.Reference,
		Message:        "payment requires user action",
		RequiresAction: true,
		ActionURL:      res.RedirectURL,
		Provider:       common.PaymentProviderMobilePay,
	}, nil
}

// VerifyPayment verifies a payment
func (s *MobilePayPaymentService) VerifyPayment(transactionID string, provider common.PaymentProviderType) (bool, error) {
	if provider != common.PaymentProviderMobilePay {
		return false, errors.New("invalid payment provider")
	}

	if transactionID == "" {
		return false, errors.New("transaction ID is required")
	}

	res, err := s.epayment.Get(transactionID)
	if err != nil {
		return false, fmt.Errorf("failed to get payment details: %v", err)
	}

	// Return true if payment is authorized
	return res.State == "AUTHORIZED", nil
}

// RefundPayment refunds a payment
func (s *MobilePayPaymentService) RefundPayment(transactionID, currency string, amount int64, provider common.PaymentProviderType) (*service.PaymentResult, error) {
	if provider != common.PaymentProviderMobilePay {
		return nil, errors.New("invalid payment provider")
	}

	if !slices.Contains(s.GetAvailableProviders()[0].SupportedCurrencies, currency) {
		return nil, fmt.Errorf("currency %s is not supported by MobilePay", currency)
	}

	// Prepare refund request
	refundRequest := models.ModificationRequest{
		ModificationAmount: models.Amount{
			Currency: currency,
			Value:    int(amount),
		},
	}

	result, err := s.epayment.Refund(transactionID, refundRequest)

	if err != nil {
		return nil, fmt.Errorf("failed to refund payment: %v", err)
	}

	return &service.PaymentResult{
		Success:        true,
		TransactionID:  result.Reference,
		Message:        "payment refunded successfully",
		RequiresAction: false,
		ActionURL:      "", // No action URL needed for refunds
		Provider:       common.PaymentProviderMobilePay,
	}, nil
}

// CapturePayment captures an authorized payment
func (s *MobilePayPaymentService) CapturePayment(transactionID, currency string, amount int64, provider common.PaymentProviderType) (*service.PaymentResult, error) {
	if provider != common.PaymentProviderMobilePay {
		return nil, errors.New("invalid payment provider")
	}

	if transactionID == "" {
		return nil, errors.New("transaction ID is required")
	}

	if !slices.Contains(s.GetAvailableProviders()[0].SupportedCurrencies, currency) {
		return nil, fmt.Errorf("currency %s is not supported by MobilePay", currency)
	}

	// Prepare capture request
	captureRequest := models.ModificationRequest{
		ModificationAmount: models.Amount{
			Currency: currency,
			Value:    int(amount),
		},
	}

	result, err := s.epayment.Capture(transactionID, captureRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to capture payment: %v", err)
	}

	return &service.PaymentResult{
		Success:        true,
		TransactionID:  result.Reference,
		Message:        "payment captured successfully",
		RequiresAction: false,
		ActionURL:      "", // No action URL needed for captures
		Provider:       common.PaymentProviderMobilePay,
	}, nil
}

// CancelPayment cancels a payment
func (s *MobilePayPaymentService) CancelPayment(transactionID string, provider common.PaymentProviderType) (*service.PaymentResult, error) {
	if provider != common.PaymentProviderMobilePay {
		return nil, errors.New("invalid payment provider")
	}

	if transactionID == "" {
		return nil, errors.New("transaction ID is required")
	}

	result, err := s.epayment.Cancel(transactionID, &models.CancelModificationRequest{
		CancelTransactionOnly: false,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to cancel payment: %v", err)
	}

	return &service.PaymentResult{
		Success:        true,
		TransactionID:  result.Reference,
		Message:        "payment cancelled successfully",
		RequiresAction: false,
		ActionURL:      "", // No action URL needed for cancellations
		Provider:       common.PaymentProviderMobilePay,
	}, nil
}

func (s *MobilePayPaymentService) ForceApprovePayment(transactionID string, phoneNumber string, provider common.PaymentProviderType) error {
	if provider != common.PaymentProviderMobilePay {
		return errors.New("invalid payment provider")
	}

	if transactionID == "" {
		return errors.New("transaction ID is required")
	}

	if phoneNumber == "" {
		return errors.New("phone number is required")
	}

	err := s.epayment.ForceApprove(transactionID, phoneNumber)

	if err != nil {
		return fmt.Errorf("failed to force approve payment: %v", err)
	}

	return nil
}

func (s *MobilePayPaymentService) GetAccessToken() error {
	err := s.vippsClient.EnsureValidToken()
	if err != nil {
		return s.vippsClient.GetAccessToken()
	}

	return nil
}
