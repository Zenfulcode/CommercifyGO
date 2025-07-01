package dto

import (
	"time"

	"github.com/zenfulcode/commercify/internal/domain/common"
)

// ProductDTO represents a product in the system
type ProductDTO struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Currency    string       `json:"currency"`
	Price       float64      `json:"price"`       // Default variant price in given currency
	SKU         string       `json:"sku"`         // Default variant SKU
	TotalStock  int          `json:"total_stock"` // Total stock across all variants
	Category    string       `json:"category"`
	CategoryID  uint         `json:"category_id"`
	Images      []string     `json:"images"`
	HasVariants bool         `json:"has_variants"`
	Active      bool         `json:"active"`
	Variants    []VariantDTO `json:"variants,omitempty"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
}

// VariantDTO represents a product variant
type VariantDTO struct {
	ID          uint             `json:"id"`
	ProductID   uint             `json:"product_id"`
	VariantName string           `json:"variant_name"`
	SKU         string           `json:"sku"`
	Stock       int              `json:"stock"`
	Attributes  common.StringMap `json:"attributes"`
	Images      []string         `json:"images"`
	IsDefault   bool             `json:"is_default"`
	Weight      float64          `json:"weight"`
	Price       float64          `json:"price"`
	Currency    string           `json:"currency"`
	CreatedAt   time.Time        `json:"created_at"`
	UpdatedAt   time.Time        `json:"updated_at"`
}
