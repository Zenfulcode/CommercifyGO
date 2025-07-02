package service

// PaymentProviderType represents a payment provider type
type PaymentProviderType string

const (
	PaymentProviderStripe    PaymentProviderType = "stripe"
	PaymentProviderMobilePay PaymentProviderType = "mobilepay"
	PaymentProviderMock      PaymentProviderType = "mock"
)

// PaymentMethod represents a payment method type
type PaymentMethod string

const (
	PaymentMethodCreditCard PaymentMethod = "credit_card"
	PaymentMethodWallet     PaymentMethod = "wallet"
)

// PaymentProvider represents information about a payment provider
type PaymentProvider struct {
	Type                PaymentProviderType `json:"type"`
	Name                string              `json:"name"`
	Description         string              `json:"description"`
	IconURL             string              `json:"icon_url,omitempty"`
	Methods             []PaymentMethod     `json:"methods"`
	Enabled             bool                `json:"enabled"`
	SupportedCurrencies []string            `json:"supported_currencies,omitempty"`
}

// PaymentRequest represents a request to process a payment
type PaymentRequest struct {
	OrderID         uint
	OrderNumber     string
	Amount          int64
	Currency        string
	PaymentMethod   PaymentMethod
	PaymentProvider PaymentProviderType
	CardDetails     *CardDetails
	PhoneNumber     string
	CustomerEmail   string
}

// CardDetails represents credit card payment details
type CardDetails struct {
	CardNumber     string `json:"card_number"`
	ExpiryMonth    int    `json:"expiry_month"`
	ExpiryYear     int    `json:"expiry_year"`
	CVV            string `json:"cvv"`
	CardholderName string `json:"cardholder_name"`
	Token          string `json:"token,omitempty"`
}

// PayPalDetails represents PayPal payment details
type PayPalDetails struct {
	Email string
	Token string
}

// BankDetails represents bank transfer details
type BankDetails struct {
	AccountNumber string `json:"account_number"`
	BankCode      string `json:"bank_code"`
	AccountName   string `json:"account_name"`
}

// PaymentResult represents the result of a payment processing
type PaymentResult struct {
	Success        bool
	TransactionID  string
	Message        string
	RequiresAction bool
	ActionURL      string
	Provider       PaymentProviderType
}

// PaymentService defines the interface for payment processing
type PaymentService interface {
	// GetAvailableProviders returns a list of available payment providers
	GetAvailableProviders() []PaymentProvider

	// GetAvailableProvidersForCurrency returns a list of available payment providers that support the given currency
	GetAvailableProvidersForCurrency(currency string) []PaymentProvider

	// ProcessPayment processes a payment request
	ProcessPayment(request PaymentRequest) (*PaymentResult, error)

	// VerifyPayment verifies a payment
	VerifyPayment(transactionID string, provider PaymentProviderType) (bool, error)

	// RefundPayment refunds a payment
	RefundPayment(transactionID, currency string, amount int64, provider PaymentProviderType) (*PaymentResult, error)

	// CapturePayment captures a payment
	CapturePayment(transactionID, currency string, amount int64, provider PaymentProviderType) (*PaymentResult, error)

	// CancelPayment cancels a payment
	CancelPayment(transactionID string, provider PaymentProviderType) (*PaymentResult, error)

	// ForceApprovePayment force approves a payment
	ForceApprovePayment(transactionID string, phoneNumber string, provider PaymentProviderType) error
}
