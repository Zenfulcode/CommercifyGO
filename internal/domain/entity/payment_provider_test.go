package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zenfulcode/commercify/internal/domain/common"
)

func TestPaymentProvider(t *testing.T) {
	t.Run("Validate success", func(t *testing.T) {
		provider := &PaymentProvider{
			Type:        common.PaymentProviderStripe,
			Name:        "Stripe",
			Description: "Stripe payment processor",
			Methods:     []common.PaymentMethod{common.PaymentMethodCreditCard},
		}

		err := provider.Validate()
		assert.NoError(t, err)
	})

	t.Run("Validate with multiple methods", func(t *testing.T) {
		provider := &PaymentProvider{
			Type:        common.PaymentProviderMobilePay,
			Name:        "MobilePay",
			Description: "MobilePay wallet",
			Methods:     []common.PaymentMethod{common.PaymentMethodWallet, common.PaymentMethodCreditCard},
		}

		err := provider.Validate()
		assert.NoError(t, err)
	})

	t.Run("Validate validation errors", func(t *testing.T) {
		tests := []struct {
			name        string
			provider    *PaymentProvider
			expectedErr string
		}{
			{
				name: "empty type",
				provider: &PaymentProvider{
					Type:    "",
					Name:    "Test Provider",
					Methods: []common.PaymentMethod{common.PaymentMethodCreditCard},
				},
				expectedErr: "payment provider type is required",
			},
			{
				name: "empty name",
				provider: &PaymentProvider{
					Type:    common.PaymentProviderStripe,
					Name:    "",
					Methods: []common.PaymentMethod{common.PaymentMethodCreditCard},
				},
				expectedErr: "payment provider name is required",
			},
			{
				name: "no methods",
				provider: &PaymentProvider{
					Type:    common.PaymentProviderStripe,
					Name:    "Test Provider",
					Methods: []common.PaymentMethod{},
				},
				expectedErr: "at least one payment method is required",
			},
			{
				name: "nil methods",
				provider: &PaymentProvider{
					Type:    common.PaymentProviderStripe,
					Name:    "Test Provider",
					Methods: nil,
				},
				expectedErr: "at least one payment method is required",
			},
			{
				name: "invalid method",
				provider: &PaymentProvider{
					Type:    common.PaymentProviderStripe,
					Name:    "Test Provider",
					Methods: []common.PaymentMethod{"invalid_method"},
				},
				expectedErr: "invalid payment method: invalid_method",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := tt.provider.Validate()
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			})
		}
	})

	t.Run("SetMethods and GetMethodsJSON", func(t *testing.T) {
		provider := &PaymentProvider{}
		methods := []common.PaymentMethod{common.PaymentMethodCreditCard, common.PaymentMethodWallet}

		provider.SetMethods(methods)
		assert.Equal(t, methods, provider.Methods)

		json, err := provider.GetMethodsJSON()
		assert.NoError(t, err)
		assert.Contains(t, json, "credit_card")
		assert.Contains(t, json, "wallet")
	})

	t.Run("SetMethodsFromJSON", func(t *testing.T) {
		provider := &PaymentProvider{}
		jsonData := []byte(`["credit_card", "wallet"]`)

		err := provider.SetMethodsFromJSON(jsonData)
		assert.NoError(t, err)
		assert.Len(t, provider.Methods, 2)
		assert.Contains(t, provider.Methods, common.PaymentMethodCreditCard)
		assert.Contains(t, provider.Methods, common.PaymentMethodWallet)

		// Test invalid JSON
		invalidJSON := []byte(`invalid json`)
		err = provider.SetMethodsFromJSON(invalidJSON)
		assert.Error(t, err)
	})

	t.Run("SetSupportedCurrencies and GetSupportedCurrenciesJSON", func(t *testing.T) {
		provider := &PaymentProvider{}
		currencies := []string{"USD", "EUR", "GBP"}

		provider.SetSupportedCurrencies(currencies)
		assert.Equal(t, currencies, provider.SupportedCurrencies)

		json, err := provider.GetSupportedCurrenciesJSON()
		assert.NoError(t, err)
		assert.Contains(t, json, "USD")
		assert.Contains(t, json, "EUR")
		assert.Contains(t, json, "GBP")
	})

	t.Run("SetSupportedCurrenciesFromJSON", func(t *testing.T) {
		provider := &PaymentProvider{}
		jsonData := []byte(`["USD", "EUR", "DKK"]`)

		err := provider.SetSupportedCurrenciesFromJSON(jsonData)
		assert.NoError(t, err)
		assert.Len(t, provider.SupportedCurrencies, 3)
		assert.Contains(t, provider.SupportedCurrencies, "USD")
		assert.Contains(t, provider.SupportedCurrencies, "EUR")
		assert.Contains(t, provider.SupportedCurrencies, "DKK")
	})

	t.Run("SetWebhookEvents and GetWebhookEventsJSON", func(t *testing.T) {
		provider := &PaymentProvider{}
		events := []string{"payment.succeeded", "payment.failed", "refund.created"}

		provider.SetWebhookEvents(events)
		assert.Equal(t, events, provider.WebhookEvents)

		json, err := provider.GetWebhookEventsJSON()
		assert.NoError(t, err)
		assert.Contains(t, json, "payment.succeeded")
		assert.Contains(t, json, "payment.failed")
		assert.Contains(t, json, "refund.created")
	})

	t.Run("SetWebhookEventsFromJSON", func(t *testing.T) {
		provider := &PaymentProvider{}
		jsonData := []byte(`["payment.succeeded", "payment.failed"]`)

		err := provider.SetWebhookEventsFromJSON(jsonData)
		assert.NoError(t, err)
		assert.Len(t, provider.WebhookEvents, 2)
		assert.Contains(t, provider.WebhookEvents, "payment.succeeded")
		assert.Contains(t, provider.WebhookEvents, "payment.failed")
	})

	t.Run("SetConfiguration and GetConfigurationJSON", func(t *testing.T) {
		provider := &PaymentProvider{}
		config := common.JSONB{
			"api_key":        "sk_test_123",
			"public_key":     "pk_test_456",
			"webhook_secret": "whsec_test_789",
		}

		provider.SetConfiguration(config)
		assert.Equal(t, config, provider.Configuration)

		json, err := provider.GetConfigurationJSON()
		assert.NoError(t, err)
		assert.Contains(t, json, "api_key")
		assert.Contains(t, json, "public_key")
		assert.Contains(t, json, "webhook_secret")

		// Test with nil config
		provider.SetConfiguration(nil)
		assert.NotNil(t, provider.Configuration)
		assert.Len(t, provider.Configuration, 0)
	})

	t.Run("SetConfigurationFromJSON", func(t *testing.T) {
		provider := &PaymentProvider{}
		jsonData := []byte(`{"api_key": "sk_test_123", "public_key": "pk_test_456"}`)

		err := provider.SetConfigurationFromJSON(jsonData)
		assert.NoError(t, err)
		assert.Equal(t, "sk_test_123", provider.Configuration["api_key"])
		assert.Equal(t, "pk_test_456", provider.Configuration["public_key"])
	})

	t.Run("SupportsCurrency", func(t *testing.T) {
		provider := &PaymentProvider{}

		// No currencies specified - should support all
		assert.True(t, provider.SupportsCurrency("USD"))
		assert.True(t, provider.SupportsCurrency("EUR"))

		// With specific currencies
		provider.SetSupportedCurrencies([]string{"USD", "EUR"})
		assert.True(t, provider.SupportsCurrency("USD"))
		assert.True(t, provider.SupportsCurrency("EUR"))
		assert.False(t, provider.SupportsCurrency("GBP"))
		assert.False(t, provider.SupportsCurrency("DKK"))
	})

	t.Run("SupportsMethod", func(t *testing.T) {
		provider := &PaymentProvider{}
		provider.SetMethods([]common.PaymentMethod{common.PaymentMethodCreditCard})

		assert.True(t, provider.SupportsMethod(common.PaymentMethodCreditCard))
		assert.False(t, provider.SupportsMethod(common.PaymentMethodWallet))

		// Add wallet method
		provider.SetMethods([]common.PaymentMethod{
			common.PaymentMethodCreditCard,
			common.PaymentMethodWallet,
		})
		assert.True(t, provider.SupportsMethod(common.PaymentMethodCreditCard))
		assert.True(t, provider.SupportsMethod(common.PaymentMethodWallet))
	})

	t.Run("ToPaymentProviderInfo", func(t *testing.T) {
		provider := &PaymentProvider{
			Type:                common.PaymentProviderStripe,
			Name:                "Stripe",
			Description:         "Stripe payment processor",
			IconURL:             "https://example.com/stripe-icon.png",
			Methods:             []common.PaymentMethod{common.PaymentMethodCreditCard},
			Enabled:             true,
			SupportedCurrencies: []string{"USD", "EUR"},
		}

		info := provider.ToPaymentProviderInfo()
		assert.Equal(t, common.PaymentProviderStripe, info.Type)
		assert.Equal(t, "Stripe", info.Name)
		assert.Equal(t, "Stripe payment processor", info.Description)
		assert.Equal(t, "https://example.com/stripe-icon.png", info.IconURL)
		assert.Equal(t, []common.PaymentMethod{common.PaymentMethodCreditCard}, info.Methods)
		assert.True(t, info.Enabled)
		assert.Equal(t, []string{"USD", "EUR"}, info.SupportedCurrencies)
	})
}

func TestPaymentProviderConstants(t *testing.T) {
	t.Run("PaymentProviderType constants", func(t *testing.T) {
		assert.Equal(t, common.PaymentProviderType("stripe"), common.PaymentProviderStripe)
		assert.Equal(t, common.PaymentProviderType("mobilepay"), common.PaymentProviderMobilePay)
		assert.Equal(t, common.PaymentProviderType("mock"), common.PaymentProviderMock)
	})

	t.Run("PaymentMethod constants", func(t *testing.T) {
		assert.Equal(t, common.PaymentMethod("credit_card"), common.PaymentMethodCreditCard)
		assert.Equal(t, common.PaymentMethod("wallet"), common.PaymentMethodWallet)
	})
}
