package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCurrency(t *testing.T) {
	t.Run("NewCurrency success", func(t *testing.T) {
		currency, err := NewCurrency(
			"usd",
			"US Dollar",
			"$",
			1.0,
			true,
			true,
		)

		require.NoError(t, err)
		assert.Equal(t, "USD", currency.Code) // Should be uppercase
		assert.Equal(t, "US Dollar", currency.Name)
		assert.Equal(t, "$", currency.Symbol)
		assert.Equal(t, 1.0, currency.ExchangeRate)
		assert.True(t, currency.IsEnabled)
		assert.True(t, currency.IsDefault)
	})

	t.Run("NewCurrency validation errors", func(t *testing.T) {
		tests := []struct {
			name         string
			code         string
			currencyName string
			symbol       string
			exchangeRate float64
			expectedErr  string
		}{
			{
				name:         "empty code",
				code:         "",
				currencyName: "US Dollar",
				symbol:       "$",
				exchangeRate: 1.0,
				expectedErr:  "currency code is required",
			},
			{
				name:         "whitespace code",
				code:         "   ",
				currencyName: "US Dollar",
				symbol:       "$",
				exchangeRate: 1.0,
				expectedErr:  "currency code is required",
			},
			{
				name:         "empty name",
				code:         "USD",
				currencyName: "",
				symbol:       "$",
				exchangeRate: 1.0,
				expectedErr:  "currency name is required",
			},
			{
				name:         "empty symbol",
				code:         "USD",
				currencyName: "US Dollar",
				symbol:       "",
				exchangeRate: 1.0,
				expectedErr:  "currency symbol is required",
			},
			{
				name:         "zero exchange rate",
				code:         "USD",
				currencyName: "US Dollar",
				symbol:       "$",
				exchangeRate: 0,
				expectedErr:  "exchange rate must be positive",
			},
			{
				name:         "negative exchange rate",
				code:         "USD",
				currencyName: "US Dollar",
				symbol:       "$",
				exchangeRate: -1.5,
				expectedErr:  "exchange rate must be positive",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				currency, err := NewCurrency(tt.code, tt.currencyName, tt.symbol, tt.exchangeRate, true, false)
				assert.Nil(t, currency)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			})
		}
	})

	t.Run("SetExchangeRate", func(t *testing.T) {
		currency, err := NewCurrency("EUR", "Euro", "€", 1.0, true, false)
		require.NoError(t, err)

		// Valid exchange rate
		err = currency.SetExchangeRate(0.85)
		assert.NoError(t, err)
		assert.Equal(t, 0.85, currency.ExchangeRate)

		// Invalid exchange rates
		err = currency.SetExchangeRate(0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exchange rate must be positive")

		err = currency.SetExchangeRate(-1.0)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "exchange rate must be positive")
	})

	t.Run("Enable/Disable", func(t *testing.T) {
		currency, err := NewCurrency("GBP", "British Pound", "£", 0.75, false, false)
		require.NoError(t, err)

		// Enable currency
		currency.Enable()
		assert.True(t, currency.IsEnabled)

		// Disable non-default currency
		err = currency.Disable()
		assert.NoError(t, err)
		assert.False(t, currency.IsEnabled)

		// Try to disable default currency
		currency.IsDefault = true
		err = currency.Disable()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot disable the default currency")
	})

	t.Run("SetAsDefault/UnsetAsDefault", func(t *testing.T) {
		currency, err := NewCurrency("CAD", "Canadian Dollar", "C$", 1.25, false, false)
		require.NoError(t, err)

		// Set as default
		currency.SetAsDefault()
		assert.True(t, currency.IsDefault)
		assert.True(t, currency.IsEnabled) // Should also be enabled

		// Unset as default
		err = currency.UnsetAsDefault()
		assert.NoError(t, err)
		assert.False(t, currency.IsDefault)
	})

	t.Run("ConvertAmount", func(t *testing.T) {
		usd, err := NewCurrency("USD", "US Dollar", "$", 1.0, true, true)
		require.NoError(t, err)

		eur, err := NewCurrency("EUR", "Euro", "€", 0.85, true, false)
		require.NoError(t, err)

		// Convert from USD to EUR
		amount := int64(10000) // $100.00
		converted := usd.ConvertAmount(amount, eur)
		assert.Equal(t, int64(8500), converted) // €85.00

		// Convert from EUR to USD
		amount = int64(8500) // €85.00
		converted = eur.ConvertAmount(amount, usd)
		assert.Equal(t, int64(10000), converted) // $100.00

		// Convert same currency
		amount = int64(10000)
		converted = usd.ConvertAmount(amount, usd)
		assert.Equal(t, int64(10000), converted)
	})

	t.Run("ToCurrencyDTO", func(t *testing.T) {
		currency, err := NewCurrency("JPY", "Japanese Yen", "¥", 110.0, true, false)
		require.NoError(t, err)

		dto := currency.ToCurrencyDTO()
		assert.Equal(t, "JPY", dto.Code)
		assert.Equal(t, "Japanese Yen", dto.Name)
		assert.Equal(t, "¥", dto.Symbol)
		assert.Equal(t, 110.0, dto.ExchangeRate)
		assert.True(t, dto.IsEnabled)
		assert.False(t, dto.IsDefault)
	})
}
