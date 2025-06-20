package dto

import (
	"time"

	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
)

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
	CategoryID  uint     `json:"category_id,omitempty"`
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

func (cp *CreateProductRequest) ToUseCaseInput() usecase.CreateProductInput {
	variants := make([]usecase.CreateVariantInput, len(cp.Variants))
	for i, v := range cp.Variants {
		variants[i] = v.ToUseCaseInput()
	}

	return usecase.CreateProductInput{
		Name:        cp.Name,
		Description: cp.Description,
		Currency:    cp.Currency,
		CategoryID:  cp.CategoryID,
		Images:      cp.Images,
		Active:      cp.Active,
		Variants:    variants,
	}
}

func (cv *CreateVariantRequest) ToUseCaseInput() usecase.CreateVariantInput {
	attributes := make([]entity.VariantAttribute, len(cv.Attributes))
	for i, attr := range cv.Attributes {
		attributes[i] = attr.ToEntity()
	}

	return usecase.CreateVariantInput{
		SKU:        cv.SKU,
		Price:      cv.Price,
		Stock:      cv.Stock,
		Attributes: attributes,
		Images:     cv.Images,
		IsDefault:  cv.IsDefault,
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

func (va *VariantAttributeDTO) ToEntity() entity.VariantAttribute {
	return entity.VariantAttribute{
		Name:  va.Name,
		Value: va.Value,
	}
}

func ToVariantDTO(variant *entity.ProductVariant) VariantDTO {
	if variant == nil {
		return VariantDTO{}
	}

	attributesDTO := make([]VariantAttributeDTO, len(variant.Attributes))
	for i, a := range variant.Attributes {
		attributesDTO[i] = VariantAttributeDTO{
			Name:  a.Name,
			Value: a.Value,
		}
	}

	return VariantDTO{
		ID:         variant.ID,
		ProductID:  variant.ProductID,
		SKU:        variant.SKU,
		Price:      money.FromCents(variant.Price),
		Currency:   variant.CurrencyCode,
		Stock:      variant.Stock,
		Attributes: attributesDTO,
		Images:     variant.Images,
		IsDefault:  variant.IsDefault,
		CreatedAt:  variant.CreatedAt,
		UpdatedAt:  variant.UpdatedAt,
	}
}

func ToProductDTO(product *entity.Product) ProductDTO {
	if product == nil {
		return ProductDTO{}
	}
	variantsDTO := make([]VariantDTO, len(product.Variants))
	for i, v := range product.Variants {
		variantsDTO[i] = ToVariantDTO(v)
	}

	return ProductDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		SKU:         product.ProductNumber,
		Price:       money.FromCents(product.Price),
		Currency:    product.CurrencyCode,
		Stock:       product.Stock,
		Weight:      product.Weight,
		CategoryID:  product.CategoryID,
		Images:      product.Images,
		HasVariants: product.HasVariants,
		Variants:    variantsDTO,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		Active:      product.Active,
	}
}
