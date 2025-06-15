package dto

import "time"

// ProductDTO represents a product in the system
type ProductDTO struct {
	ID          uint         `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	SKU         string       `json:"sku"`
	Price       float64      `json:"price"`
	Currency    string       `json:"currency"`
	Stock       int          `json:"stock"`
	Weight      float64      `json:"weight"`
	CategoryID  uint         `json:"category_id"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	Images      []string     `json:"images"`
	HasVariants bool         `json:"has_variants"`
	Variants    []VariantDTO `json:"variants,omitempty"`
	Active      bool         `json:"active"`
}

// VariantDTO represents a product variant
type VariantDTO struct {
	ID         uint                  `json:"id"`
	ProductID  uint                  `json:"product_id"`
	SKU        string                `json:"sku"`
	Price      float64               `json:"price"`
	Currency   string                `json:"currency"`
	Stock      int                   `json:"stock"`
	Attributes []VariantAttributeDTO `json:"attributes"`
	Images     []string              `json:"images,omitempty"`
	IsDefault  bool                  `json:"is_default"`
	CreatedAt  time.Time             `json:"created_at"`
	UpdatedAt  time.Time             `json:"updated_at"`
}

type VariantAttributeDTO struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CreateProductRequest represents the data needed to create a new product
type CreateProductRequest struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Currency    string                 `json:"currency"`
	CategoryID  uint                   `json:"category_id"`
	Images      []string               `json:"images"`
	Active      bool                   `json:"active"`
	Variants    []CreateVariantRequest `json:"variants,omitempty"`
}

// CreateVariantRequest represents the data needed to create a new product variant
type CreateVariantRequest struct {
	SKU        string                `json:"sku"`
	Price      float64               `json:"price"`
	Stock      int                   `json:"stock"`
	Attributes []VariantAttributeDTO `json:"attributes"`
	Images     []string              `json:"images,omitempty"`
	IsDefault  bool                  `json:"is_default,omitempty"`
}

// UpdateProductRequest represents the data needed to update an existing product
type UpdateProductRequest struct {
	Name        string   `json:"name,omitempty"`
	Description string   `json:"description,omitempty"`
	Currency    string   `json:"currency,omitempty"`
	CategoryID  *uint    `json:"category_id,omitempty"`
	Images      []string `json:"images,omitempty"`
	Active      bool     `json:"active,omitempty"`
}

// UpdateVariantRequest represents the data needed to update an existing product variant
type UpdateVariantRequest struct {
	SKU        string                `json:"sku,omitempty"`
	Price      *float64              `json:"price,omitempty"`
	Stock      *int                  `json:"stock,omitempty"`
	Attributes []VariantAttributeDTO `json:"attributes,omitempty"`
	Images     []string              `json:"images,omitempty"`
	IsDefault  *bool                 `json:"is_default,omitempty"`
}

// ProductListResponse represents a paginated list of products
type ProductListResponse struct {
	ListResponseDTO[ProductDTO]
}
