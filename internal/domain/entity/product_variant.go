package entity

import (
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/money"
)

// VariantAttribute represents a single attribute of a product variant
type VariantAttribute struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// ProductVariant represents a specific variant of a product
type ProductVariant struct {
	ID           uint                  `json:"id"`
	ProductID    uint                  `json:"product_id"`
	SKU          string                `json:"sku"`
	Price        int64                 `json:"price"` // Stored as cents (in default currency)
	CurrencyCode string                `json:"currency"`
	Stock        int                   `json:"stock"`
	Attributes   []VariantAttribute    `json:"attributes"`
	Images       []string              `json:"images"`
	IsDefault    bool                  `json:"is_default"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
	Prices       []ProductVariantPrice `json:"prices,omitempty"` // Prices in different currencies
}

// NewProductVariant creates a new product variant
func NewProductVariant(productID uint, sku string, price float64, currencyCode string, stock int, attributes []VariantAttribute, images []string, isDefault bool) (*ProductVariant, error) {
	if productID == 0 {
		return nil, errors.New("product ID cannot be empty")
	}
	if sku == "" {
		return nil, errors.New("SKU cannot be empty")
	}
	if price <= 0 { // Check cents
		return nil, errors.New("price must be greater than zero")
	}
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}
	// Note: attributes can be empty for default variants

	// Convert price to cents
	priceInCents := money.ToCents(price)

	now := time.Now()
	return &ProductVariant{
		ProductID:    productID,
		SKU:          sku,
		Price:        priceInCents, // Already in cents
		CurrencyCode: currencyCode,
		Stock:        stock,
		Attributes:   attributes,
		Images:       images,
		IsDefault:    isDefault,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// UpdateStock updates the variant's stock
func (v *ProductVariant) UpdateStock(quantity int) error {
	newStock := v.Stock + quantity
	if newStock < 0 {
		return errors.New("insufficient stock")
	}

	v.Stock = newStock
	v.UpdatedAt = time.Now()
	return nil
}

// IsAvailable checks if the variant is available in the requested quantity
func (v *ProductVariant) IsAvailable(quantity int) bool {
	return v.Stock >= quantity
}

// GetPriceInCurrency returns the price in the specified currency
func (v *ProductVariant) GetPriceInCurrency(currencyCode string) (int64, bool) {
	for _, price := range v.Prices {
		if price.CurrencyCode == currencyCode {
			return price.Price, true
		}
	}

	return v.Price, false
}

// SetPriceInCurrency sets or updates the price for a specific currency
func (v *ProductVariant) SetPriceInCurrency(currencyCode string, price float64) error {
	if currencyCode == "" {
		return errors.New("currency code cannot be empty")
	}
	if price <= 0 {
		return errors.New("price must be greater than zero")
	}

	priceInCents := money.ToCents(price)

	// Check if price already exists for this currency
	for i, existingPrice := range v.Prices {
		if existingPrice.CurrencyCode == currencyCode {
			// Update existing price
			v.Prices[i].Price = priceInCents
			v.Prices[i].UpdatedAt = time.Now()
			v.UpdatedAt = time.Now()
			return nil
		}
	}

	// Add new price
	now := time.Now()
	newPrice := ProductVariantPrice{
		VariantID:    v.ID,
		CurrencyCode: currencyCode,
		Price:        priceInCents,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	v.Prices = append(v.Prices, newPrice)
	v.UpdatedAt = time.Now()
	return nil
}

// RemovePriceInCurrency removes the price for a specific currency
func (v *ProductVariant) RemovePriceInCurrency(currencyCode string) error {
	if currencyCode == "" {
		return errors.New("currency code cannot be empty")
	}

	// Don't allow removing the default currency price
	if currencyCode == v.CurrencyCode {
		return errors.New("cannot remove default currency price")
	}

	for i, price := range v.Prices {
		if price.CurrencyCode == currencyCode {
			// Remove the price by slicing
			v.Prices = append(v.Prices[:i], v.Prices[i+1:]...)
			v.UpdatedAt = time.Now()
			return nil
		}
	}

	return errors.New("price not found for the specified currency")
}

// GetAllPrices returns all prices including the default price
func (v *ProductVariant) GetAllPrices() map[string]int64 {
	prices := make(map[string]int64)

	// Add default price
	prices[v.CurrencyCode] = v.Price

	// Add additional currency prices
	for _, price := range v.Prices {
		prices[price.CurrencyCode] = price.Price
	}

	return prices
}

// HasPriceInCurrency checks if the variant has a price set for the specified currency
func (v *ProductVariant) HasPriceInCurrency(currencyCode string) bool {
	// Check if it's the default currency
	if currencyCode == v.CurrencyCode {
		return true
	}

	// Check additional currency prices
	for _, price := range v.Prices {
		if price.CurrencyCode == currencyCode {
			return true
		}
	}

	return false
}
