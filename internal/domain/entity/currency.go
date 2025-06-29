package entity

import (
	"errors"
	"strings"
	"time"

	"gorm.io/gorm"
)

// Currency represents a currency in the system
type Currency struct {
	gorm.Model           // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	Code         string  `gorm:"primaryKey;size:3"`
	Name         string  `gorm:"size:100;not null"`
	Symbol       string  `gorm:"size:10;not null"`
	ExchangeRate float64 `gorm:"not null;default:1.0"`
	IsEnabled    bool    `gorm:"not null;default:true"`
	IsDefault    bool    ` gorm:"not null;default:false"`
}

// ProductPrice represents a price for a product variant in a specific currency
// This maps to the product_prices table which is linked to variants, not products directly
type ProductPrice struct {
	gorm.Model          // Includes ID, CreatedAt, UpdatedAt, DeletedAt
	VariantID    uint   `gorm:"uniqueIndex:idx_variantprice"`
	CurrencyCode string `gorm:"uniqueIndex:idx_variantprice;size:3;not null"`
	Price        int64  `gorm:"not null"` // Price in cents
} // Unique constraint: UNIQUE(variant_id, currency_code)

// TableName returns the table name for ProductPrice
func (ProductPrice) TableName() string {
	return "product_prices"
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

	return &Currency{
		Code:         strings.ToUpper(code),
		Name:         name,
		Symbol:       symbol,
		ExchangeRate: exchangeRate,
		IsEnabled:    isEnabled,
		IsDefault:    isDefault,
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
