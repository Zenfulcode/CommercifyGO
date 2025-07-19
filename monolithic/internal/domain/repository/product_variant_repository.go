package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// ProductVariantRepository defines the interface for product variant data access
type ProductVariantRepository interface {
	Create(variant *entity.ProductVariant) error
	GetByID(variantID uint) (*entity.ProductVariant, error)
	GetBySKU(sku string) (*entity.ProductVariant, error)
	GetByProduct(productID uint) ([]*entity.ProductVariant, error)
	Update(variant *entity.ProductVariant) error
	Delete(variantID uint) error
	BatchCreate(variants []*entity.ProductVariant) error
}
