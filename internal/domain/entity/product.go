package entity

import (
	"errors"
	"fmt"
	"slices"

	"github.com/zenfulcode/commercify/internal/domain/dto"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// Product represents a product in the system
// All products must have at least one variant as per the database schema
type Product struct {
	gorm.Model
	Name        string                      `gorm:"not null;size:255"`
	Description string                      `gorm:"type:text"`
	Currency    string                      `gorm:"not null;size:3"`
	CategoryID  uint                        `gorm:"not null;index"`
	Category    Category                    `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT,OnUpdate:CASCADE"`
	Images      datatypes.JSONSlice[string] `gorm:"type:text[];default:'[]'"`
	Active      bool                        `gorm:"default:true"`
	Variants    []*ProductVariant           `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

// NewProduct creates a new product with the given details
// Note: At least one variant must be added before the product can be considered complete
func NewProduct(name, description, currency string, categoryID uint, images []string, variants []*ProductVariant, isActive bool) (*Product, error) {
	if name == "" {
		return nil, errors.New("product name cannot be empty")
	}

	if categoryID == 0 {
		return nil, errors.New("category ID cannot be zero")
	}

	if len(variants) == 0 {
		return nil, errors.New("at least one variant must be provided")
	}

	// Copy variants to ensure product has its own slice
	productVariants := make([]*ProductVariant, len(variants))
	copy(productVariants, variants)

	return &Product{
		Name:        name,
		Description: description,
		Currency:    currency,
		CategoryID:  categoryID,
		Images:      images,
		Variants:    productVariants,
		Active:      isActive,
	}, nil
}

// IsAvailable checks if the product is available in the requested quantity
// For products with variants, this checks if any variant has sufficient stock
func (p *Product) IsAvailable(quantity int) bool {
	if !p.HasVariants() {
		return false // Product must have variants
	}

	// Check if any variant has sufficient stock
	for _, variant := range p.Variants {
		if variant.IsAvailable(quantity) {
			return true
		}
	}
	return false
}

// AddVariant adds a variant to the product
func (p *Product) AddVariant(variant *ProductVariant) error {
	if variant == nil {
		return errors.New("variant cannot be nil")
	}

	variant.ProductID = p.ID

	// Add variant to product
	p.Variants = append(p.Variants, variant)
	return nil
}

// RemoveVariant removes a variant from the product by its ID
func (p *Product) RemoveVariant(variantID uint) error {
	if len(p.Variants) == 0 {
		return errors.New("no variants available to remove")
	}

	for i, variant := range p.Variants {
		if variant.ID == variantID {
			// Remove the variant from the slice
			p.Variants = append(p.Variants[:i], p.Variants[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("variant with ID %d not found", variantID)
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

// GetTotalWeight calculates the total weight for a quantity of the default variant
func (p *Product) GetTotalWeight(quantity int) float64 {
	if quantity <= 0 {
		return 0
	}

	defaultVariant := p.GetDefaultVariant()
	if defaultVariant == nil {
		return 0
	}

	return defaultVariant.Weight * float64(quantity)
}

func (p *Product) GetStockForVariant(variantID uint) (int, error) {
	if len(p.Variants) == 0 {
		return 0, errors.New("no variants available for this product")
	}

	for _, variant := range p.Variants {
		if variant.ID == variantID {
			return variant.Stock, nil
		}
	}

	return 0, fmt.Errorf("variant with ID %d not found", variantID)
}

// GetTotalStock calculates the total stock across all variants
func (p *Product) GetTotalStock() int {
	totalStock := 0
	for _, variant := range p.Variants {
		totalStock += variant.Stock
	}
	return totalStock
}

func (p *Product) HasVariants() bool {
	return len(p.Variants) > 0
}

func (p *Product) Update(name *string, description *string, currency *string, images *[]string, active *bool, categoryID *uint) bool {
	updated := false
	if name != nil && *name != "" && p.Name != *name {
		p.Name = *name
		updated = true
	}
	if description != nil && *description != "" && p.Description != *description {
		p.Description = *description
		updated = true
	}
	if currency != nil && *currency != "" && p.Currency != *currency {
		p.Currency = *currency
		updated = true
	}
	if images != nil && len(*images) > 0 && !slices.Equal(p.Images, *images) {
		p.Images = *images
		updated = true
	}
	if active != nil && p.Active != *active {
		p.Active = *active
		updated = true
	}
	if categoryID != nil && p.CategoryID != *categoryID {
		p.CategoryID = *categoryID
		updated = true
	}

	return updated
}

func (p *Product) GetProdNumber() string {
	if p == nil || len(p.Variants) == 0 {
		return ""
	}

	defaultVariant := p.GetDefaultVariant()
	if defaultVariant != nil {
		return defaultVariant.SKU
	}

	return ""
}

func (p *Product) GetPrice() int64 {
	if p == nil || len(p.Variants) == 0 {
		return 0
	}

	defaultVariant := p.GetDefaultVariant()
	if defaultVariant != nil {
		return defaultVariant.Price
	}

	return 0
}

func (p *Product) ToProductDTO() *dto.ProductDTO {
	if p == nil {
		return nil
	}

	variantsDTO := make([]dto.VariantDTO, len(p.Variants))
	for i, v := range p.Variants {
		variantsDTO[i] = *v.ToVariantDTO()
	}

	defaultVariant := p.GetDefaultVariant()
	if defaultVariant == nil {

	}

	return &dto.ProductDTO{
		ID:          p.ID,
		Name:        p.Name,
		SKU:         p.GetProdNumber(),
		Description: p.Description,
		Currency:    p.Currency,
		TotalStock:  p.GetTotalStock(),
		Price:       money.FromCents(p.GetPrice()),
		Category:    p.Category.Name,
		CategoryID:  p.CategoryID,
		Images:      p.Images,
		HasVariants: p.HasVariants(),
		Active:      p.Active,
		Variants:    variantsDTO,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}

func (p *Product) ToProductSummaryDTO() *dto.ProductDTO {
	return &dto.ProductDTO{
		ID:          p.ID,
		Name:        p.Name,
		SKU:         p.GetProdNumber(),
		Description: p.Description,
		Currency:    p.Currency,
		TotalStock:  p.GetTotalStock(),
		Price:       money.FromCents(p.GetPrice()),
		Category:    p.Category.Name,
		Images:      p.Images,
		HasVariants: p.HasVariants(),
		Active:      p.Active,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
