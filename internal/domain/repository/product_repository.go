package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(product *entity.Product) error
	GetByID(productID uint) (*entity.Product, error)
	GetByIDAndCurrency(productID uint, currency string) (*entity.Product, error)
	GetBySKU(sku string) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(productID uint) error
	List(query, currency string, categoryID, offset, limit uint, minPriceCents, maxPriceCents int64, active bool) ([]*entity.Product, error)
	Count(searchQuery, currency string, categoryID uint, minPriceCents, maxPriceCents int64, active bool) (int, error)
}

// CategoryRepository defines the interface for category data access
type CategoryRepository interface {
	Create(category *entity.Category) error
	GetByID(categoryID uint) (*entity.Category, error)
	Update(category *entity.Category) error
	Delete(categoryID uint) error
	List() ([]*entity.Category, error)
	GetChildren(parentID uint) ([]*entity.Category, error)
}
