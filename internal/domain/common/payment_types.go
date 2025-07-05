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

// IsValidPaymentMethod checks if the payment method is valid
func IsValidPaymentMethod(method string) bool {
	switch PaymentMethod(method) {
	case PaymentMethodCreditCard, PaymentMethodWallet:
		return true
	default:
		return false
	}
}
