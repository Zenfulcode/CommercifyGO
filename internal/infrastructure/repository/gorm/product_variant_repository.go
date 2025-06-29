package gorm

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/gorm"
)

// ProductVariantRepository implements repository.ProductVariantRepository using GORM
type ProductVariantRepository struct {
	db *gorm.DB
}

// NewProductVariantRepository creates a new GORM-based ProductVariantRepository
func NewProductVariantRepository(db *gorm.DB) repository.ProductVariantRepository {
	return &ProductVariantRepository{db: db}
}

// Create creates a new product variant with its prices
func (r *ProductVariantRepository) Create(variant *entity.ProductVariant) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create the variant
		if err := tx.Create(variant).Error; err != nil {
			return err
		}

		return nil
	})
}

// BatchCreate creates multiple variants at once
func (r *ProductVariantRepository) BatchCreate(variants []*entity.ProductVariant) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		for _, variant := range variants {
			if err := tx.Create(variant).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// GetByID retrieves a variant by ID with prices
func (r *ProductVariantRepository) GetByID(variantID uint) (*entity.ProductVariant, error) {
	var variant entity.ProductVariant
	if err := r.db.Preload("Prices").First(&variant, variantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("variant not found")
		}
		return nil, err
	}
	return &variant, nil
}

// GetBySKU retrieves a variant by SKU with prices
func (r *ProductVariantRepository) GetBySKU(sku string) (*entity.ProductVariant, error) {
	var variant entity.ProductVariant
	if err := r.db.Preload("Prices").Where("sku = ?", sku).First(&variant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("variant not found")
		}
		return nil, err
	}
	return &variant, nil
}

// GetByProduct retrieves all variants for a product with prices
func (r *ProductVariantRepository) GetByProduct(productID uint) ([]*entity.ProductVariant, error) {
	var variants []*entity.ProductVariant
	if err := r.db.Preload("Prices").Where("product_id = ?", productID).Find(&variants).Error; err != nil {
		return nil, err
	}
	return variants, nil
}

// Update updates an existing variant
func (r *ProductVariantRepository) Update(variant *entity.ProductVariant) error {
	return r.db.Save(variant).Error
}

// Delete deletes a variant by ID
func (r *ProductVariantRepository) Delete(variantID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete associated prices first (cascade should handle this)
		if err := tx.Where("variant_id = ?", variantID).Delete(&entity.ProductPrice{}).Error; err != nil {
			return err
		}

		// Delete the variant
		return tx.Delete(&entity.ProductVariant{}, variantID).Error
	})
}

// UpdateStock updates the stock for a variant
func (r *ProductVariantRepository) UpdateStock(variantID uint, quantity int) error {
	return r.db.Model(&entity.ProductVariant{}).Where("id = ?", variantID).Update("stock", quantity).Error
}

// List retrieves variants with filtering and pagination
func (r *ProductVariantRepository) List(offset, limit uint) ([]*entity.ProductVariant, error) {
	var variants []*entity.ProductVariant
	if err := r.db.Preload("Prices").Offset(int(offset)).Limit(int(limit)).Find(&variants).Error; err != nil {
		return nil, err
	}
	return variants, nil
}
