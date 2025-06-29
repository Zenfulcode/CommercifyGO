package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/dto"
	"gorm.io/gorm"
)

// VariantAttributes represents JSONB attributes for a product variant
type VariantAttributes map[string]string

// Value implements the driver.Valuer interface for database storage
func (va VariantAttributes) Value() (driver.Value, error) {
	return json.Marshal(va)
}

// Scan implements the sql.Scanner interface for database retrieval
func (va *VariantAttributes) Scan(value any) error {
	if value == nil {
		*va = make(VariantAttributes)
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, va)
}

// ProductVariant represents a specific variant of a product
type ProductVariant struct {
	gorm.Model
	ProductID  uint              `gorm:"index"`
	SKU        string            `gorm:"uniqueIndex;size:100;not null"`
	Stock      int               `gorm:"default:0"`
	Attributes VariantAttributes `gorm:"type:json;not null"`
	Images     []string          `gorm:"type:json"`
	IsDefault  bool              `gorm:"default:false"`
	Weight     float64           `gorm:"default:0"`
	Prices     []ProductPrice    `gorm:"foreignKey:VariantID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE;"`
}

// NewProductVariant creates a new product variant
func NewProductVariant(sku string, stock int, weight float64, attributes VariantAttributes, prices []ProductPrice, images []string, isDefault bool) (*ProductVariant, error) {
	if sku == "" {
		return nil, errors.New("SKU cannot be empty")
	}
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}
	if weight < 0 {
		return nil, errors.New("weight cannot be negative")
	}

	if attributes == nil {
		attributes = make(VariantAttributes)
	}

	if len(prices) == 0 {
		return nil, errors.New("at least one price must be provided")
	}

	for _, price := range prices {
		if price.CurrencyCode == "" {
			return nil, errors.New("currency code cannot be empty")
		}
		if price.Price < 0 {
			return nil, errors.New("price cannot be negative for currency " + price.CurrencyCode)
		}
	}

	return &ProductVariant{
		SKU:        sku,
		Stock:      stock,
		Attributes: attributes,
		Images:     images,
		IsDefault:  isDefault,
		Weight:     weight,
		Prices:     prices,
	}, nil
}

// UpdateStock updates the variant's stock
func (v *ProductVariant) UpdateStock(quantity int) error {
	newStock := v.Stock + quantity
	if newStock < 0 {
		return errors.New("insufficient stock")
	}

	v.Stock = newStock
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
	return 0, false
}

// SetPriceInCurrency sets or updates the price for a specific currency
func (v *ProductVariant) SetPriceInCurrency(currencyCode string, priceInCents int64) error {
	if currencyCode == "" {
		return errors.New("currency code cannot be empty")
	}
	if priceInCents < 0 {
		return errors.New("price cannot be negative")
	}

	// Check if price already exists for this currency
	for i, existingPrice := range v.Prices {
		if existingPrice.CurrencyCode == currencyCode {
			v.Prices[i].Price = priceInCents
			return nil
		}
	}

	newPrice := ProductPrice{
		VariantID:    v.ID,
		CurrencyCode: currencyCode,
		Price:        priceInCents,
	}

	v.Prices = append(v.Prices, newPrice)
	return nil
}

// RemovePriceInCurrency removes the price for a specific currency
func (v *ProductVariant) RemovePriceInCurrency(currencyCode string) error {
	if currencyCode == "" {
		return errors.New("currency code cannot be empty")
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

// GetAllPrices returns all prices for this variant
func (v *ProductVariant) GetAllPrices() map[string]int64 {
	prices := make(map[string]int64)

	// Add all currency prices
	for _, price := range v.Prices {
		prices[price.CurrencyCode] = price.Price
	}

	return prices
}

// HasPriceInCurrency checks if the variant has a price set for the specified currency
func (v *ProductVariant) HasPriceInCurrency(currencyCode string) bool {
	// Check currency prices
	for _, price := range v.Prices {
		if price.CurrencyCode == currencyCode {
			return true
		}
	}

	return false
}

func (v *ProductVariant) Name() string {
	// Combine all attribute values to form a name
	name := ""
	for _, value := range v.Attributes {
		if name == "" {
			name = value
		} else {
			name += " / " + value
		}
	}
	return name
}

// Remove VariantAttributeDTO as we'll use map directly
func (variant *ProductVariant) ToVariantDTO() *dto.VariantDTO {
	if variant == nil {
		return nil
	}

	// Get all prices and convert from cents to float
	allPricesInCents := variant.GetAllPrices()
	allPrices := make(map[string]float64)
	for currency, priceInCents := range allPricesInCents {
		allPrices[currency] = money.FromCents(priceInCents)
	}

	return &dto.VariantDTO{
		ID:         variant.ID,
		ProductID:  variant.ProductID,
		SKU:        variant.SKU,
		Stock:      variant.Stock,
		Attributes: variant.Attributes,
		Images:     variant.Images,
		IsDefault:  variant.IsDefault,
		Weight:     variant.Weight,
		Prices:     allPrices,
		CreatedAt:  variant.CreatedAt,
		UpdatedAt:  variant.UpdatedAt,
	}
}
