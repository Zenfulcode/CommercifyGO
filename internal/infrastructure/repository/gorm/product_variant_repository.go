package gorm

import (
	"errors"
	"fmt"

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

// Create creates a new product variant
func (r *ProductVariantRepository) Create(variant *entity.ProductVariant) error {
	return r.db.Create(variant).Error
}

// BatchCreate creates multiple variants at once
func (r *ProductVariantRepository) BatchCreate(variants []*entity.ProductVariant) error {
	if len(variants) == 0 {
		return nil
	}
	// Use GORM's CreateInBatches for better performance
	return r.db.CreateInBatches(variants, 100).Error
}

// GetByID retrieves a variant by ID with product relationship
func (r *ProductVariantRepository) GetByID(variantID uint) (*entity.ProductVariant, error) {
	var variant entity.ProductVariant
	if err := r.db.Preload("Product").First(&variant, variantID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("variant with ID %d not found", variantID)
		}
		return nil, fmt.Errorf("failed to fetch variant: %w", err)
	}
	return &variant, nil
}

// GetBySKU retrieves a variant by SKU with product relationship
func (r *ProductVariantRepository) GetBySKU(sku string) (*entity.ProductVariant, error) {
	var variant entity.ProductVariant
	if err := r.db.Preload("Product").Where("sku = ?", sku).First(&variant).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("variant with SKU %s not found", sku)
		}
		return nil, fmt.Errorf("failed to fetch variant by SKU: %w", err)
	}
	return &variant, nil
}

// GetByProduct retrieves all variants for a product with product relationship
func (r *ProductVariantRepository) GetByProduct(productID uint) ([]*entity.ProductVariant, error) {
	var variants []*entity.ProductVariant
	if err := r.db.Preload("Product").Where("product_id = ?", productID).Find(&variants).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch variants for product %d: %w", productID, err)
	}
	return variants, nil
}

// Update updates an existing variant
func (r *ProductVariantRepository) Update(variant *entity.ProductVariant) error {
	return r.db.Save(variant).Error
}

// Delete deletes a variant by ID
func (r *ProductVariantRepository) Delete(variantID uint) error {
	return r.db.Delete(&entity.ProductVariant{}, variantID).Error
}
