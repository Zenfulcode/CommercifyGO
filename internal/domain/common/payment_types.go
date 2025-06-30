package common

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
