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
	Variants    []CreateVariantRequest `json:"variants"`
}

// AttributeKeyValue represents a key-value pair for product attributes
type AttributeKeyValue struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

// CreateVariantRequest represents the data needed to create a new product variant
type CreateVariantRequest struct {
	SKU        string              `json:"sku"`
	Stock      int                 `json:"stock"`
	Attributes []AttributeKeyValue `json:"attributes"`
	Images     []string            `json:"images"`
	IsDefault  bool                `json:"is_default"`
	Weight     float64             `json:"weight"`
	Price      float64             `json:"price"`
}

// UpdateProductRequest represents the data needed to update an existing product
type UpdateProductRequest struct {
	Name        *string                 `json:"name,omitempty"`
	Description *string                 `json:"description,omitempty"`
	Currency    *string                 `json:"currency,omitempty"`
	CategoryID  *uint                   `json:"category_id,omitempty"`
	Images      *[]string               `json:"images,omitempty"`
	Active      *bool                   `json:"active,omitempty"`
	Variants    *[]UpdateVariantRequest `json:"variants,omitempty"` // Optional, can be nil if no variants are updated
}

// UpdateVariantRequest represents the data needed to update an existing product variant
type UpdateVariantRequest struct {
	SKU        *string              `json:"sku,omitempty"`
	Stock      *int                 `json:"stock,omitempty"`
	Attributes *[]AttributeKeyValue `json:"attributes,omitempty"`
	Images     *[]string            `json:"images,omitempty"`
	IsDefault  *bool                `json:"is_default,omitempty"`
	Weight     *float64             `json:"weight,omitempty"`
	Price      *float64             `json:"price,omitempty"`
}

func CreateProductListResponse(products []*entity.Product, totalCount, page, pageSize int) ListResponseDTO[dto.ProductDTO] {
	var productDTOs []dto.ProductDTO
	for _, product := range products {
		productDTOs = append(productDTOs, *product.ToProductSummaryDTO())
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
		Currency:    cp.Currency,
		CategoryID:  cp.CategoryID,
		Images:      cp.Images,
		Active:      cp.Active,
		Variants:    variants,
	}
}

func (cv *CreateVariantRequest) ToUseCaseInput() usecase.CreateVariantInput {
	// Convert attributes array to map
	attributesMap := make(entity.VariantAttributes)
	for _, attr := range cv.Attributes {
		attributesMap[attr.Name] = attr.Value
	}

	return usecase.CreateVariantInput{
		VariantInput: usecase.VariantInput{
			SKU:        cv.SKU,
			Stock:      cv.Stock,
			Weight:     cv.Weight,
			Images:     cv.Images,
			Attributes: attributesMap,
			Price:      money.ToCents(cv.Price),
			IsDefault:  cv.IsDefault,
		},
	}
}

func (up *UpdateProductRequest) ToUseCaseInput() usecase.UpdateProductInput {
	input := usecase.UpdateProductInput{
		Name:        up.Name,
		Description: up.Description,
		CategoryID:  up.CategoryID,
		Images:      up.Images,
		Active:      up.Active,
	}

	// Convert variants if provided
	if up.Variants != nil {
		variants := make([]usecase.UpdateVariantInput, len(*up.Variants))
		for i, v := range *up.Variants {
			variants[i] = v.ToUseCaseInput()
		}
		input.Variants = &variants
	}

	return input
}

func (u UpdateVariantRequest) ToUseCaseInput() usecase.UpdateVariantInput {
	var variantInput usecase.VariantInput

	// Set defaults for required fields
	variantInput.SKU = ""
	variantInput.Stock = 0
	variantInput.Price = 0
	variantInput.Weight = 0
	variantInput.IsDefault = false
	variantInput.Images = []string{}
	variantInput.Attributes = make(map[string]string)

	// Update with provided values
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
		// Convert attributes array to map
		attributesMap := make(map[string]string)
		for _, attr := range *u.Attributes {
			attributesMap[attr.Name] = attr.Value
		}
		variantInput.Attributes = attributesMap
	}

	return usecase.UpdateVariantInput{
		VariantInput: variantInput,
	}
}
