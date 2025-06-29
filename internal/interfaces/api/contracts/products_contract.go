package contracts

import (
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/dto"
)

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
	SKU        string            `json:"sku"`
	Stock      int               `json:"stock"`
	Attributes map[string]string `json:"attributes"`
	Images     []string          `json:"images"`
	IsDefault  bool              `json:"is_default"`
	Weight     float64           `json:"weight"`
	Price      float64           `json:"price"`
}

// UpdateProductRequest represents the data needed to update an existing product
type UpdateProductRequest struct {
	Name        *string   `json:"name,omitempty"`
	Description *string   `json:"description,omitempty"`
	Currency    *string   `json:"currency,omitempty"`
	CategoryID  *uint     `json:"category_id,omitempty"`
	Images      *[]string `json:"images,omitempty"`
	Active      *bool     `json:"active,omitempty"`
}

// UpdateVariantRequest represents the data needed to update an existing product variant
type UpdateVariantRequest struct {
	SKU        *string            `json:"sku,omitempty"`
	Stock      *int               `json:"stock,omitempty"`
	Attributes *map[string]string `json:"attributes,omitempty"`
	Images     *[]string          `json:"images,omitempty"`
	IsDefault  *bool              `json:"is_default,omitempty"`
	Weight     *float64           `json:"weight,omitempty"`
	Price      *float64           `json:"price,omitempty"`
}

// ProductListResponse represents a paginated list of products
type ProductListResponse struct {
	ListResponseDTO[dto.ProductDTO]
}

func (cp *CreateProductRequest) ToUseCaseInput() usecase.CreateProductInput {
	variants := make([]usecase.CreateVariantInput, len(cp.Variants))
	for i, v := range cp.Variants {
		variants[i] = v.ToUseCaseInput()
	}

	return usecase.CreateProductInput{
		Name:        cp.Name,
		Description: cp.Description,
		CategoryID:  cp.CategoryID,
		Images:      cp.Images,
		Active:      cp.Active,
		Variants:    variants,
	}
}

func (cv *CreateVariantRequest) ToUseCaseInput() usecase.CreateVariantInput {
	return usecase.CreateVariantInput{
		VariantInput: usecase.VariantInput{
			SKU:        cv.SKU,
			Stock:      cv.Stock,
			Weight:     cv.Weight,
			Images:     cv.Images,
			Attributes: cv.Attributes,
			Price:      money.ToCents(cv.Price),
		},
		IsDefault: cv.IsDefault,
	}
}

func (up *UpdateProductRequest) ToUseCaseInput() usecase.UpdateProductInput {
	return usecase.UpdateProductInput{
		Name:        up.Name,
		Description: up.Description,
		CategoryID:  up.CategoryID,
		Images:      up.Images,
		Active:      up.Active,
	}
}
