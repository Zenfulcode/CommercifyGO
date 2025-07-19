package entity

import (
	"errors"
	"slices"

	"github.com/zenfulcode/commercify/internal/domain/dto"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type VariantAttributes = map[string]string

// ProductVariant represents a specific variant of a product
type ProductVariant struct {
	gorm.Model
	ProductID  uint                                  `gorm:"index;not null"`
	Product    Product                               `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	SKU        string                                `gorm:"uniqueIndex;size:100;not null"`
	Stock      int                                   `gorm:"default:0"`
	Attributes datatypes.JSONType[VariantAttributes] `gorm:"not null"`
	IsDefault  bool                                  `gorm:"default:false"`
	Weight     float64                               `gorm:"default:0"`
	Price      int64                                 `gorm:"not null"`
	Images     datatypes.JSONSlice[string]
}

// NewProductVariant creates a new product variant
func NewProductVariant(sku string, stock int, price int64, weight float64, attributes VariantAttributes, images []string, isDefault bool) (*ProductVariant, error) {
	if sku == "" {
		return nil, errors.New("SKU cannot be empty")
	}
	if stock < 0 {
		return nil, errors.New("stock cannot be negative")
	}
	if price < 0 {
		return nil, errors.New("price cannot be negative")
	}
	if weight < 0 {
		return nil, errors.New("weight cannot be negative")
	}

	if attributes == nil {
		attributes = make(VariantAttributes)
	}

	return &ProductVariant{
		SKU:        sku,
		Stock:      stock,
		Attributes: datatypes.NewJSONType(attributes),
		Images:     images,
		IsDefault:  isDefault,
		Weight:     weight,
		Price:      price,
	}, nil
}

func (v *ProductVariant) Update(SKU string, stock int, price int64, weight float64, images []string, attributes VariantAttributes, isDefault *bool) (bool, error) {
	updated := false
	if SKU != "" && v.SKU != SKU {
		v.SKU = SKU
		updated = true
	}
	if stock >= 0 && v.Stock != stock {
		v.Stock = stock
		updated = true
	}
	if price >= 0 && v.Price != price {
		v.Price = price
		updated = true
	}
	if weight >= 0 && v.Weight != weight {
		v.Weight = weight
		updated = true
	}

	if len(images) > 0 && !slices.Equal([]string(v.Images), images) {
		v.Images = images
		updated = true
	}
	if len(attributes) > 0 {
		v.Attributes = datatypes.NewJSONType(attributes)
		updated = true
	}
	if isDefault != nil && v.IsDefault != *isDefault {
		v.IsDefault = *isDefault
		updated = true
	}

	return updated, nil
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

func (v *ProductVariant) Name() string {
	// Combine all attribute values to form a name
	name := ""
	for _, value := range v.Attributes.Data() {
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

	return &dto.VariantDTO{
		ID:          variant.ID,
		ProductID:   variant.ProductID,
		VariantName: variant.Name(),
		SKU:         variant.SKU,
		Stock:       variant.Stock,
		Attributes:  variant.Attributes.Data(),
		Images:      variant.Images,
		IsDefault:   variant.IsDefault,
		Weight:      variant.Weight,
		Price:       money.FromCents(variant.Price),
		Currency:    variant.Product.Currency,
		CreatedAt:   variant.CreatedAt,
		UpdatedAt:   variant.UpdatedAt,
	}
}
