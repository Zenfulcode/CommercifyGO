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
		}

		return nil
	})
}

// GetByID retrieves a product by ID without variants
func (r *ProductRepository) GetByID(productID uint) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.Preload("Variants").First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &product, nil
}

// GetBySKU implements repository.ProductRepository.
func (r *ProductRepository) GetBySKU(sku string) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.Preload("Variants").Joins("JOIN product_variants ON products.id = product_variants.product_id").
		Where("product_variants.sku = ?", sku).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	// Ensure product has variants loaded
	if len(product.Variants) == 0 {
		return nil, errors.New("product has no variants")
	}
	return &product, nil
}

func (r *ProductRepository) GetByIDAndCurrency(productID uint, currency string) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.Preload("Variants", func(db *gorm.DB) *gorm.DB {
		if currency != "" {
			return db.Where("currency = ?", currency)
		}
		return db
	}).First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	// Ensure product has variants loaded
	if len(product.Variants) == 0 {
		return nil, errors.New("product has no variants")
	}

	return &product, nil
}

// Update updates an existing product
func (r *ProductRepository) Update(product *entity.Product) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Update the product itself
		if err := tx.Save(product).Error; err != nil {
			return err
		}

		// Update variants
		for _, variant := range product.Variants {
			if err := tx.Save(variant).Error; err != nil {
				return err
			}
		}

		return nil
	})
}

// Delete deletes a product by ID
func (r *ProductRepository) Delete(productID uint) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
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
		tx = tx.Joins("JOIN product_variants ON products.id = product_variants.product_id")

		if currency != "" {
			tx = tx.Where("products.currency = ?", currency)
		}

		if minPriceCents > 0 {
			tx = tx.Where("product_variants.price >= ?", minPriceCents)
		}

		if maxPriceCents > 0 {
			tx = tx.Where("product_variants.price <= ?", maxPriceCents)
		}

		tx = tx.Distinct()
	}

	// Apply pagination
	if err := tx.Offset(int(offset)).Limit(int(limit)).Preload("Variants").Find(&products).Error; err != nil {
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
		tx = tx.Joins("JOIN product_variants ON products.id = product_variants.product_id")

		if currency != "" {
			tx = tx.Where("products.currency = ?", currency)
		}

		if minPriceCents > 0 {
			tx = tx.Where("product_variants.price >= ?", minPriceCents)
		}

		if maxPriceCents > 0 {
			tx = tx.Where("product_variants.price <= ?", maxPriceCents)
		}

		tx = tx.Distinct()
	}

	if err := tx.Count(&count).Error; err != nil {
		return 0, err
	}

	return int(count), nil
}
