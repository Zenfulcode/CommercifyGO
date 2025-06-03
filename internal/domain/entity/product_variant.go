package entity

import (
	"errors"
	"time"
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
func NewProductVariant(productID uint, sku string, price int64, currencyCode string, stock int, attributes []VariantAttribute, images []string, isDefault bool) (*ProductVariant, error) {
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

	now := time.Now()
	return &ProductVariant{
		ProductID:    productID,
		SKU:          sku,
		Price:        price, // Already in cents
		CurrencyCode: currencyCode,
		Stock:        stock,
		Attributes:   attributes,
		Images:       images,
		IsDefault:    isDefault,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

// NewDefaultProductVariant creates a default product variant using the product's basic information
// This is used when a product needs at least one variant but no specific variants are provided
func NewDefaultProductVariant(productID uint, sku string, price int64, currencyCode string, stock int) (*ProductVariant, error) {
	return NewProductVariant(productID, sku, price, currencyCode, stock, []VariantAttribute{}, []string{}, true)
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
