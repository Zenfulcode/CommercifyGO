package gorm

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/gorm"
)

// ProductRepository implements repository.ProductRepository using GORM
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository creates a new GORM-based ProductRepository
func NewProductRepository(db *gorm.DB) repository.ProductRepository {
	return &ProductRepository{db: db}
}

// Create creates a new product with its variants and prices
func (r *ProductRepository) Create(product *entity.Product) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create the product
		if err := tx.Create(product).Error; err != nil {
			return err
		}

		// Create variants if any
		for _, variant := range product.Variants {
			variant.ProductID = product.ID
			if err := tx.Create(variant).Error; err != nil {
				return err
			}

			// Create prices for this variant
			for i := range variant.Prices {
				variant.Prices[i].VariantID = variant.ID
				if err := tx.Create(&variant.Prices[i]).Error; err != nil {
					return err
				}
			}
		}

		return nil
	})
}

// GetByID retrieves a product by ID without variants
func (r *ProductRepository) GetByID(productID uint) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.Preload("Variants.Prices").First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(product *entity.Product) error {
	return r.db.Save(product).Error
}

// Delete deletes a product by ID
func (r *ProductRepository) Delete(productID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Delete all associated prices first (cascade should handle this, but being explicit)
		if err := tx.Where("variant_id IN (SELECT id FROM product_variants WHERE product_id = ?)", productID).Delete(&entity.ProductPrice{}).Error; err != nil {
			return err
		}

		// Delete all variants (cascade should handle this)
		if err := tx.Where("product_id = ?", productID).Delete(&entity.ProductVariant{}).Error; err != nil {
			return err
		}

		// Delete the product
		return tx.Delete(&entity.Product{}, productID).Error
	})
}

// List retrieves products with filtering and pagination
func (r *ProductRepository) List(query, currency string, categoryID, offset, limit uint, minPriceCents, maxPriceCents int64, active bool) ([]*entity.Product, error) {
	var products []*entity.Product

	tx := r.db.Model(&entity.Product{})

	// Apply filters
	if query != "" {
		tx = tx.Where("name ILIKE ? OR description ILIKE ?", "%"+query+"%", "%"+query+"%")
	}

	if categoryID > 0 {
		tx = tx.Where("category_id = ?", categoryID)
	}

	tx = tx.Where("active = ?", active)

	// Price filtering requires joining with variants and prices
	if minPriceCents > 0 || maxPriceCents > 0 || currency != "" {
		tx = tx.Joins("JOIN product_variants ON products.id = product_variants.product_id").
			Joins("JOIN product_prices ON product_variants.id = product_prices.variant_id")

		if currency != "" {
			tx = tx.Where("product_prices.currency_code = ?", currency)
		}

		if minPriceCents > 0 {
			tx = tx.Where("product_prices.price >= ?", minPriceCents)
		}

		if maxPriceCents > 0 {
			tx = tx.Where("product_prices.price <= ?", maxPriceCents)
		}

		tx = tx.Distinct()
	}

	// Apply pagination
	if err := tx.Offset(int(offset)).Limit(int(limit)).Preload("Variants.Prices").Find(&products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

// Count returns the total count of products matching the filter criteria
func (r *ProductRepository) Count(searchQuery, currency string, categoryID uint, minPriceCents, maxPriceCents int64, active bool) (int, error) {
	var count int64

	tx := r.db.Model(&entity.Product{})

	// Apply same filters as List method
	if searchQuery != "" {
		tx = tx.Where("name ILIKE ? OR description ILIKE ?", "%"+searchQuery+"%", "%"+searchQuery+"%")
	}

	if categoryID > 0 {
		tx = tx.Where("category_id = ?", categoryID)
	}

	tx = tx.Where("active = ?", active)

	// Price filtering requires joining with variants and prices
	if minPriceCents > 0 || maxPriceCents > 0 || currency != "" {
		tx = tx.Joins("JOIN product_variants ON products.id = product_variants.product_id").
			Joins("JOIN product_prices ON product_variants.id = product_prices.variant_id")

		if currency != "" {
			tx = tx.Where("product_prices.currency_code = ?", currency)
		}

		if minPriceCents > 0 {
			tx = tx.Where("product_prices.price >= ?", minPriceCents)
		}

		if maxPriceCents > 0 {
			tx = tx.Where("product_prices.price <= ?", maxPriceCents)
		}

		tx = tx.Distinct()
	}

	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
