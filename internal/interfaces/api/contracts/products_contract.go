package contracts

import (
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/dto"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
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

func CreateProductListResponse(products []*entity.Product, totalCount, page, pageSize int) ListResponseDTO[dto.ProductDTO] {
	var productDTOs []dto.ProductDTO
	for _, product := range products {
		productDTOs = append(productDTOs, *product.ToProductDTO())
	}
	if len(productDTOs) == 0 {
		return ListResponseDTO[dto.ProductDTO]{
			Success:    true,
			Data:       []dto.ProductDTO{},
			Pagination: PaginationDTO{Page: page, PageSize: pageSize, Total: 0},
			Message:    "No products found",
		}
	}

	return ListResponseDTO[dto.ProductDTO]{
		Success: true,
		Data:    productDTOs,
		Pagination: PaginationDTO{
			Page:     page,
			PageSize: pageSize,
			Total:    totalCount,
		},
		Message: "Products retrieved successfully",
	}
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
			IsDefault:  cv.IsDefault,
		},
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

func (u UpdateVariantRequest) ToUseCaseInput() usecase.UpdateVariantInput {

	var variantInput usecase.VariantInput
	if u.SKU != nil {
		variantInput.SKU = *u.SKU
	}
	if u.Stock != nil {
		variantInput.Stock = *u.Stock
	}
	if u.Weight != nil {
		variantInput.Weight = *u.Weight
	}
	if u.Images != nil {
		variantInput.Images = *u.Images
	}
	if u.Price != nil {
		variantInput.Price = money.ToCents(*u.Price)
	}
	if u.IsDefault != nil {
		variantInput.IsDefault = *u.IsDefault
	}
	if u.Attributes != nil {
		variantInput.Attributes = *u.Attributes
	}

	return usecase.UpdateVariantInput{
		VariantInput: variantInput,
	}
}
