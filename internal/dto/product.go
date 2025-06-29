package dto

import (
	"time"
)

// ProductDTO represents a product in the system
type ProductDTO struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
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
	ID         uint               `json:"id"`
	ProductID  uint               `json:"product_id"`
	SKU        string             `json:"sku"`
	Stock      int                `json:"stock"`
	Attributes map[string]string  `json:"attributes"`
	Images     []string           `json:"images"`
	IsDefault  bool               `json:"is_default"`
	Weight     float64            `json:"weight"`
	Prices     map[string]float64 `json:"prices"` // All prices in different currencies
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}
