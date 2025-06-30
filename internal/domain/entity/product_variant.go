package entity

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"maps"
	"slices"

	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/dto"
	"github.com/zenfulcode/commercify/internal/domain/money"
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
	ProductID  uint               `gorm:"index;not null"`
	Product    Product            `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	SKU        string             `gorm:"uniqueIndex;size:100;not null"`
	Stock      int                `gorm:"default:0"`
	Attributes VariantAttributes  `gorm:"type:jsonb;not null"`
	Images     common.StringSlice `gorm:"type:json;default:'[]'"`
	IsDefault  bool               `gorm:"default:false"`
	Weight     float64            `gorm:"default:0"`
	Price      int64              `gorm:"not null"`
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
		Attributes: attributes,
		Images:     common.StringSlice(images),
		IsDefault:  isDefault,
		Weight:     weight,
		Price:      price,
	}, nil
}

func (v *ProductVariant) Update(SKU string, stock int, price int64, weight float64, images []string, attributes VariantAttributes) (bool, error) {
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
		v.Images = common.StringSlice(images)
		updated = true
	}
	if len(attributes) > 0 {
		// Convert slice of maps to VariantAttributes map
		newAttributes := make(VariantAttributes)
		maps.Copy(newAttributes, attributes)
		if !maps.Equal(v.Attributes, newAttributes) {
			v.Attributes = newAttributes
			updated = true
		}
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

	return &dto.VariantDTO{
		ID:          variant.ID,
		ProductID:   variant.ProductID,
		VariantName: variant.Name(),
		SKU:         variant.SKU,
		Stock:       variant.Stock,
		Attributes:  variant.Attributes,
		Images:      variant.Images,
		IsDefault:   variant.IsDefault,
		Weight:      variant.Weight,
		Price:       money.FromCents(variant.Price),
		Currency:    variant.Product.Currency,
		CreatedAt:   variant.CreatedAt,
		UpdatedAt:   variant.UpdatedAt,
	}
}
