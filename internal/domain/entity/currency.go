package entity

import (
	"errors"
	"strings"
	"time"
)

// Currency represents a currency in the system
type Currency struct {
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Symbol       string    `json:"symbol"`
	ExchangeRate float64   `json:"exchange_rate"`
	IsEnabled    bool      `json:"is_enabled"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ProductPrice represents a price for a product in a specific currency
type ProductPrice struct {
	ID           uint      `json:"id"`
	ProductID    uint      `json:"product_id"`
	CurrencyCode string    `json:"currency_code"`
	Price        int64     `json:"price"` // Price in cents
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ProductVariantPrice represents a price for a product variant in a specific currency
type ProductVariantPrice struct {
	ID           uint      `json:"id"`
	VariantID    uint      `json:"variant_id"`
	CurrencyCode string    `json:"currency_code"`
	Price        int64     `json:"price"` // Price in cents
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// NewCurrency creates a new Currency
func NewCurrency(code, name, symbol string, exchangeRate float64, isEnabled bool, isDefault bool) (*Currency, error) {
	// Validate required fields
	if strings.TrimSpace(code) == "" {
		return nil, errors.New("currency code is required")
	}

	if strings.TrimSpace(name) == "" {
		return nil, errors.New("currency name is required")
	}

	if strings.TrimSpace(symbol) == "" {
		return nil, errors.New("currency symbol is required")
	}

	if exchangeRate <= 0 {
		return nil, errors.New("exchange rate must be positive")
	}

	now := time.Now()
	return &Currency{
		Code:         strings.ToUpper(code),
		Name:         name,
		Symbol:       symbol,
		ExchangeRate: exchangeRate,
		IsEnabled:    isEnabled,
		IsDefault:    isDefault,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// SetExchangeRate sets the exchange rate for the currency
func (c *Currency) SetExchangeRate(rate float64) error {
	if rate <= 0 {
		return errors.New("exchange rate must be positive")
	}
	c.ExchangeRate = rate
	c.UpdatedAt = time.Now()
	return nil
}

// Enable enables the currency
func (c *Currency) Enable() {
	c.IsEnabled = true
	c.UpdatedAt = time.Now()
}

// Disable disables the currency
func (c *Currency) Disable() error {
	if c.IsDefault {
		return errors.New("cannot disable the default currency")
	}
	c.IsEnabled = false
	c.UpdatedAt = time.Now()
	return nil
}

// SetAsDefault sets this currency as the default currency
func (c *Currency) SetAsDefault() {
	c.IsDefault = true
	c.IsEnabled = true // Default currency must be enabled
	c.UpdatedAt = time.Now()
}

// UnsetAsDefault unsets this currency as the default currency
func (c *Currency) UnsetAsDefault() error {
	c.IsDefault = false
	c.UpdatedAt = time.Now()
	return nil
}

// ConvertAmount converts an amount from this currency to the target currency
func (c *Currency) ConvertAmount(amount int64, targetCurrency *Currency) int64 {
	if c.Code == targetCurrency.Code {
		return amount
	}

	// First convert to a base unit
	baseAmount := float64(amount) / c.ExchangeRate

	// Then convert to target currency
	targetAmount := baseAmount * targetCurrency.ExchangeRate

	// Round to nearest cent instead of truncating
	return int64(targetAmount)
}
