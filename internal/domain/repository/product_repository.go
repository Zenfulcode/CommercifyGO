package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// ProductRepository defines the interface for product data access
type ProductRepository interface {
	Create(product *entity.Product) error
	GetByID(productID uint) (*entity.Product, error)
	GetByIDWithVariants(productID uint) (*entity.Product, error)
	GetByProductNumber(productNumber string) (*entity.Product, error)
	Update(product *entity.Product) error
	Delete(productID uint) error
	List(offset, limit int) ([]*entity.Product, error)
	// Search expects minPriceCents and maxPriceCents as int64 (cents)
	Search(query string, categoryID uint, minPriceCents, maxPriceCents int64, offset, limit int) ([]*entity.Product, error)
	Count() (int, error)
	CountSearch(searchQuery string, categoryID uint, minPriceCents, maxPriceCents int64) (int, error)
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
