package gorm

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/gorm"
)

// ProductPriceRepository handles product price operations using GORM
type ProductPriceRepository struct {
	db *gorm.DB
}

// NewProductPriceRepository creates a new GORM-based ProductPriceRepository
func NewProductPriceRepository(db *gorm.DB) repository.ProductPriceRepository {
	return &ProductPriceRepository{
		db: db,
	}
}

// Create creates a new price entry
func (r ProductPriceRepository) Create(price *entity.ProductPrice) error {
	return r.db.Create(price).Error
}

// GetByVariantAndCurrency retrieves a price by variant ID and currency code
func (r ProductPriceRepository) GetByVariantIDAndCurrency(variantID uint, currencyCode string) (*entity.ProductPrice, error) {
	var price entity.ProductPrice
	if err := r.db.Where("variant_id = ? AND currency_code = ?", variantID, currencyCode).First(&price).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("price not found")
		}
		return nil, err
	}
	return &price, nil
}

// GetByVariant retrieves all prices for a variant
func (r ProductPriceRepository) GetByVariantID(variantID uint) ([]entity.ProductPrice, error) {
	var prices []entity.ProductPrice
	if err := r.db.Where("variant_id = ?", variantID).Find(&prices).Error; err != nil {
		return nil, err
	}
	return prices, nil
}

// Update updates an existing price
func (r ProductPriceRepository) Update(price *entity.ProductPrice) error {
	return r.db.Save(price).Error
}

// Upsert creates or updates a price for a variant in a specific currency
func (r ProductPriceRepository) Upsert(variantID uint, currencyCode string, priceInCents int64) error {
	price := entity.ProductPrice{
		VariantID:    variantID,
		CurrencyCode: currencyCode,
		Price:        priceInCents,
	}

	// Use ON CONFLICT to update if exists, insert if not
	return r.db.Where("variant_id = ? AND currency_code = ?", variantID, currencyCode).
		Assign(entity.ProductPrice{Price: priceInCents}).
		FirstOrCreate(&price).Error
}

// Delete deletes a price by variant ID and currency code
func (r ProductPriceRepository) Delete(id uint) error {
	return r.db.Where("id = ?", id).
		Delete(&entity.ProductPrice{}).Error
}

func (r ProductPriceRepository) DeleteByVariantIDAndCurrency(variantID uint, currencyCode string) error {
	return r.db.Where("variant_id = ? AND currency_code = ?", variantID, currencyCode).
		Delete(&entity.ProductPrice{}).Error
}

// BatchUpsert creates or updates multiple prices at once
func (r ProductPriceRepository) BatchUpsert(prices []entity.ProductPrice) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, price := range prices {
			if err := tx.Where("variant_id = ? AND currency_code = ?", price.VariantID, price.CurrencyCode).
				Assign(entity.ProductPrice{Price: price.Price}).
				FirstOrCreate(&price).Error; err != nil {
				return err
			}
		}
		return nil
	})
}
