package gorm

import (
	"errors"
	"fmt"

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

// Create creates a new product with its variants
func (r *ProductRepository) Create(product *entity.Product) error {
	// GORM will automatically create associated variants due to the relationship definition
	return r.db.Create(product).Error
}

// GetByID retrieves a product by ID with all related data
func (r *ProductRepository) GetByID(productID uint) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.Preload("Variants").Preload("Category").First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product with ID %d not found", productID)
		}
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}
	return &product, nil
}

// GetBySKU retrieves a product by variant SKU
func (r *ProductRepository) GetBySKU(sku string) (*entity.Product, error) {
	var product entity.Product
	if err := r.db.Preload("Variants").Preload("Category").
		Joins("JOIN product_variants ON products.id = product_variants.product_id").
		Where("product_variants.sku = ?", sku).First(&product).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product with SKU %s not found", sku)
		}
		return nil, fmt.Errorf("failed to fetch product by SKU: %w", err)
	}

	// Ensure product has variants loaded
	if len(product.Variants) == 0 {
		return nil, fmt.Errorf("product with SKU %s has no variants", sku)
	}
	return &product, nil
}

// GetByIDAndCurrency retrieves a product by ID, filtering for the specified currency
func (r *ProductRepository) GetByIDAndCurrency(productID uint, currency string) (*entity.Product, error) {
	var product entity.Product

	// Build the query
	query := r.db.Preload("Category")

	// Filter variants by currency if specified, otherwise load all variants
	if currency != "" {
		query = query.Preload("Variants", "price IS NOT NULL") // Basic validation that variant has a price
		query = query.Where("currency = ?", currency)
	} else {
		query = query.Preload("Variants")
	}

	if err := query.First(&product, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("product with ID %d not found", productID)
		}
		return nil, fmt.Errorf("failed to fetch product: %w", err)
	}

	// Ensure product has variants loaded
	if len(product.Variants) == 0 {
		return nil, fmt.Errorf("product with ID %d has no variants for currency %s", productID, currency)
	}

	return &product, nil
}

// Update updates an existing product and its variants
func (r *ProductRepository) Update(product *entity.Product) error {
	// Use FullSaveAssociations to handle variant updates properly
	return r.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(product).Error
}

// Delete deletes a product by ID and its associated variants (hard deletion)
func (r *ProductRepository) Delete(productID uint) error {
	// Use a transaction to ensure data consistency
	return r.db.Transaction(func(tx *gorm.DB) error {
		// First, hard delete all variants for this product
		if err := tx.Unscoped().Where("product_id = ?", productID).Delete(&entity.ProductVariant{}).Error; err != nil {
			return fmt.Errorf("failed to delete product variants: %w", err)
		}

		// Then hard delete the product itself
		if err := tx.Unscoped().Delete(&entity.Product{}, productID).Error; err != nil {
			return fmt.Errorf("failed to delete product: %w", err)
		}

		return nil
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

	if currency != "" {
		tx = tx.Where("currency = ?", currency)
	}

	// Active filter: if active=true, only show active products; if active=false, show all products
	if active {
		tx = tx.Where("active = ?", true)
	}
	// If active is false, don't add any filter to show all products (active and inactive)

	// Price filtering requires joining with variants
	if minPriceCents > 0 || maxPriceCents > 0 {
		tx = tx.Joins("JOIN product_variants ON products.id = product_variants.product_id")

		if minPriceCents > 0 {
			tx = tx.Where("product_variants.price >= ?", minPriceCents)
		}

		if maxPriceCents > 0 {
			tx = tx.Where("product_variants.price <= ?", maxPriceCents)
		}

		tx = tx.Distinct()
	}

	// Apply pagination and load relationships
	if err := tx.Offset(int(offset)).Limit(int(limit)).
		Preload("Variants").Preload("Category").
		Find(&products).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch products: %w", err)
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

	if currency != "" {
		tx = tx.Where("currency = ?", currency)
	}

	// Only filter by active status if active=true
	// If active=false, return all products (active and inactive)
	if active {
		tx = tx.Where("active = ?", true)
	}

	// Price filtering requires joining with variants
	if minPriceCents > 0 || maxPriceCents > 0 {
		tx = tx.Joins("JOIN product_variants ON products.id = product_variants.product_id")

		if minPriceCents > 0 {
			tx = tx.Where("product_variants.price >= ?", minPriceCents)
		}

		if maxPriceCents > 0 {
			tx = tx.Where("product_variants.price <= ?", maxPriceCents)
		}

		tx = tx.Distinct()
	}

	if err := tx.Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count products: %w", err)
	}

	return int(count), nil
}
