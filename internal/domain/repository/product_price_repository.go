package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// ProductPriceRepository defines the contract for product price operations
type ProductPriceRepository interface {
	Create(price *entity.ProductPrice) error
	Update(price *entity.ProductPrice) error
	Delete(id uint) error
	GetByVariantID(variantID uint) ([]entity.ProductPrice, error)
	GetByVariantIDAndCurrency(variantID uint, currencyCode string) (*entity.ProductPrice, error)
	DeleteByVariantIDAndCurrency(variantID uint, currencyCode string) error
}
