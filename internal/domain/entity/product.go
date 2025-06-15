package entity

import (
	"errors"
	"fmt"
	"time"
)

// Product represents a product in the system
// Note: All products must have at least one variant. ProductNumber is deprecated in favor of variant SKUs.
type Product struct {
	ID            uint              `json:"id"`
	ProductNumber string            `json:"product_number,omitempty"` // Deprecated: Use variant SKUs instead
	Name          string            `json:"name"`
	Description   string            `json:"description"`
	Price         int64             `json:"price"` // Stored as cents (default variant price)
	CurrencyCode  string            `json:"currency_code,omitempty"`
	Stock         int               `json:"stock"`  // Aggregate stock from variants
	Weight        float64           `json:"weight"` // Weight in kg
	CategoryID    uint              `json:"category_id"`
	Images        []string          `json:"images"`
	HasVariants   bool              `json:"has_variants"` // Always true, kept for backward compatibility
	Variants      []*ProductVariant `json:"variants,omitempty"`
	Prices        []ProductPrice    `json:"prices,omitempty"` // Prices in different currencies
	CreatedAt     time.Time         `json:"created_at"`
	UpdatedAt     time.Time         `json:"updated_at"`
	Active        bool              `json:"active"`
}

// NewProduct creates a new product with the given details (price in cents)
// Note: This creates a product structure, but at least one variant must be added before saving
func NewProduct(name, description string, currencyCode string, categoryID uint, images []string) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}

	now := time.Now()

	// Generate a temporary product number (deprecated, variants will have SKUs)
	productNumber := "PROD-TEMP"

	return &Product{
		Name:          name,
		ProductNumber: productNumber,
		Description:   description,
		Price:         0, // Already in cents
		CurrencyCode:  currencyCode,
		Stock:         0,
		Weight:        0.0,
		CategoryID:    categoryID,
		Images:        images,
		HasVariants:   false,
		Active:        false,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

func (p *Product) IsComplete() bool {
	// A product is complete if it has a name, description, and at least one variant
	if p.Name == "" || p.Description == "" || len(p.Variants) == 0 {
		return false
	}

	// Ensure at least one variant has a SKU and price
	for _, variant := range p.Variants {
		if variant.SKU == "" || variant.Price <= 0 {
			return false
		}
	}

	return true
}

// IsAvailable checks if the product is available in the requested quantity
func (p *Product) IsAvailable(quantity int) bool {
	if p.HasVariants {
		// For products with variants, availability depends on variants
		return true
	}
	return p.Stock >= quantity
}

// AddVariant adds a variant to the product
func (p *Product) AddVariant(variant *ProductVariant) error {
	if variant == nil {
		return errors.New("variant cannot be nil")
	}

	// Ensure variant belongs to this product
	if variant.ProductID != p.ID {
		return errors.New("variant does not belong to this product")
	}

	// If this is the first variant and it's the default, set product price to match
	if len(p.Variants) == 0 && variant.IsDefault {
		p.Price = variant.Price
		p.Stock = variant.Stock
	}

	variant.CurrencyCode = p.CurrencyCode

	// Add variant to product
	p.Variants = append(p.Variants, variant)

	// Only set has_variants=true if there are now multiple variants
	p.HasVariants = len(p.Variants) > 1

	p.UpdatedAt = time.Now()

	return nil
}

// GetDefaultVariant returns the default variant of the product
func (p *Product) GetDefaultVariant() *ProductVariant {
	if len(p.Variants) == 0 {
		return nil
	}

	for _, variant := range p.Variants {
		if variant.IsDefault {
			return variant
		}
	}

	// If no default is set, return the first variant
	return p.Variants[0]
}

// GetVariantByID returns a variant by its ID
func (p *Product) GetVariantByID(variantID uint) *ProductVariant {
	if len(p.Variants) == 0 {
		return nil
	}

	for _, variant := range p.Variants {
		if variant.ID == variantID {
			return variant
		}
	}

	return nil
}

// GetVariantBySKU returns a variant by its SKU
func (p *Product) GetVariantBySKU(sku string) *ProductVariant {
	if len(p.Variants) == 0 || sku == "" {
		return nil
	}

	for _, variant := range p.Variants {
		if variant.SKU == sku {
			return variant
		}
	}

	return nil
}

// SetProductNumber sets the product number
func (p *Product) SetProductNumber(id uint) {
	// Format: PROD-000001
	p.ProductNumber = fmt.Sprintf("PROD-%06d", id)
}

// GetTotalWeight calculates the total weight for a quantity of this product
func (p *Product) GetTotalWeight(quantity int) float64 {
	if quantity <= 0 {
		return 0
	}
	return p.Weight * float64(quantity)
}

// GetPriceInCurrency returns the price for a specific currency
func (p *Product) GetPriceInCurrency(currencyCode string) (int64, bool) {
	variant := p.GetDefaultVariant()
	if variant != nil {
		return variant.GetPriceInCurrency(currencyCode)
	}

	for _, productPrice := range p.Prices {
		if productPrice.CurrencyCode == currencyCode {
			return productPrice.Price, true
		}
	}

	return p.Price, false
}

// Category represents a product category
type Category struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    *uint     `json:"parent_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// NewCategory creates a new category
func NewCategory(name, description string, parentID *uint) (*Category, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	now := time.Now()
	return &Category{
		Name:        name,
		Description: description,
		ParentID:    parentID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}
