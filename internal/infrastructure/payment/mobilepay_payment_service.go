package payment

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/gkhaavik/vipps-mobilepay-sdk/pkg/client"
	"github.com/gkhaavik/vipps-mobilepay-sdk/pkg/models"
	"github.com/google/uuid"
	"github.com/zenfulcode/commercify/config"
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
			Type:                service.PaymentProviderMobilePay,
			Name:                "MobilePay",
			Description:         "Pay with MobilePay app",
			IconURL:             "/assets/images/mobilepay-logo.png",
			Methods:             []service.PaymentMethod{service.PaymentMethodWallet},
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
		for _, supportedCurrency := range provider.SupportedCurrencies {
			if supportedCurrency == currency {
				supportedProviders = append(supportedProviders, provider)
				break
			}
		}
	}

	return supportedProviders
}

// ProcessPayment processes a payment request using MobilePay
func (s *MobilePayPaymentService) ProcessPayment(request service.PaymentRequest) (*service.PaymentResult, error) {
	// Only wallet payment method is supported for MobilePay
	if request.PaymentMethod != service.PaymentMethodWallet {
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: "unsupported payment method for MobilePay, only wallet is supported",
			Provider:     service.PaymentProviderMobilePay,
		}, nil
	}

	if request.PhoneNumber == "" {
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: "phone number is required for MobilePay payments",
			Provider:     service.PaymentProviderMobilePay,
		}, nil
	}

	phoneNumber := request.PhoneNumber

	r := regexp.MustCompile(`^\+?[1-9]\d{1,14}$`)

	if !r.MatchString(phoneNumber) {
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: "invalid phone number format, must be in international format",
			Provider:     service.PaymentProviderMobilePay,
		}, nil
	}

	// Generate a unique reference for this payment
	reference := fmt.Sprintf("order-%d-%s", request.OrderID, uuid.New().String())

	// MobilePay only supports specific currencies (NOK, DKK, EUR)
	// Use the configured market currency instead of the request currency
	mobilePayCurrency := s.config.Market
	if mobilePayCurrency == "" {
		mobilePayCurrency = "NOK" // Default to NOK if not configured
	}

	// Construct the payment request
	paymentRequest := models.CreatePaymentRequest{
		Amount: models.Amount{
			Currency: mobilePayCurrency,
			Value:    int(request.Amount),
		},
		Customer: &models.Customer{
			PhoneNumber: &phoneNumber,
		},
		PaymentMethod: &models.PaymentMethod{
			Type: "WALLET",
		},
		Reference:          reference,
		ReturnURL:          s.config.ReturnURL + "?reference=" + reference,
		UserFlow:           models.UserFlowWebRedirect,
		PaymentDescription: s.config.PaymentDescription,
	}

	res, err := s.epayment.Create(paymentRequest)
	if err != nil {
		return &service.PaymentResult{
			Success:      false,
			ErrorMessage: fmt.Sprintf("failed to create payment: %v", err),
			Provider:     service.PaymentProviderMobilePay,
		}, nil
	}

	// MobilePay requires a redirect to complete the payment
	// Return a result with action URL
	return &service.PaymentResult{
		Success:        false,
		TransactionID:  res.Reference,
		ErrorMessage:   "payment requires user action",
		RequiresAction: true,
		ActionURL:      res.RedirectURL,
		Provider:       service.PaymentProviderMobilePay,
	}, nil
}

// VerifyPayment verifies a payment
func (s *MobilePayPaymentService) VerifyPayment(transactionID string, provider service.PaymentProviderType) (bool, error) {
	if provider != service.PaymentProviderMobilePay {
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
func (s *MobilePayPaymentService) RefundPayment(transactionID string, amount int64, provider service.PaymentProviderType) error {
	if provider != service.PaymentProviderMobilePay {
		return errors.New("invalid payment provider")
	}

	// MobilePay only supports specific currencies (NOK, DKK, EUR)
	// Use the configured market currency
	mobilePayCurrency := s.config.Market
	if mobilePayCurrency == "" {
		mobilePayCurrency = "NOK" // Default to NOK if not configured
	}

	// Prepare refund request
	refundRequest := models.ModificationRequest{
		ModificationAmount: models.Amount{
			Currency: mobilePayCurrency,
			Value:    int(amount),
		},
	}

	_, err := s.epayment.Refund(transactionID, refundRequest)

	if err != nil {
		return fmt.Errorf("failed to refund payment: %v", err)
	}

	return nil
}

// CapturePayment captures an authorized payment
func (s *MobilePayPaymentService) CapturePayment(transactionID string, amount int64, provider service.PaymentProviderType) error {
	if provider != service.PaymentProviderMobilePay {
		return errors.New("invalid payment provider")
	}

	if transactionID == "" {
		return errors.New("transaction ID is required")
	}

	// MobilePay only supports specific currencies (NOK, DKK, EUR)
	// Use the configured market currency
	mobilePayCurrency := s.config.Market
	if mobilePayCurrency == "" {
		mobilePayCurrency = "NOK" // Default to NOK if not configured
	}

	// Prepare capture request
	captureRequest := models.ModificationRequest{
		ModificationAmount: models.Amount{
			Currency: mobilePayCurrency,
			Value:    int(amount),
		},
	}

	_, err := s.epayment.Capture(transactionID, captureRequest)
	if err != nil {
		return fmt.Errorf("failed to capture payment: %v", err)
	}

	return nil
}

// CancelPayment cancels a payment
func (s *MobilePayPaymentService) CancelPayment(transactionID string, provider service.PaymentProviderType) error {
	if provider != service.PaymentProviderMobilePay {
		return errors.New("invalid payment provider")
	}

	if transactionID == "" {
		return errors.New("transaction ID is required")
	}

	_, err := s.epayment.Cancel(transactionID, &models.CancelModificationRequest{
		CancelTransactionOnly: false,
	})

	if err != nil {
		return fmt.Errorf("failed to cancel payment: %v", err)
	}

	return nil
}

func (s *MobilePayPaymentService) ForceApprovePayment(transactionID string, phoneNumber string, provider service.PaymentProviderType) error {
	if provider != service.PaymentProviderMobilePay {
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
