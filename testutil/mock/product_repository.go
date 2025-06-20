package mock

import (
	"errors"
	"strings"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// MockProductRepository is a mock implementation of product repository for testing
type MockProductRepository struct {
	products    map[uint]*entity.Product
	lastID      uint
	searchCount int
}

// NewMockProductRepository creates a new instance of MockProductRepository
func NewMockProductRepository() repository.ProductRepository {
	return &MockProductRepository{
		products:    make(map[uint]*entity.Product),
		lastID:      0,
		searchCount: 0,
	}
}

// Count returns the number of products in the repository
func (r *MockProductRepository) Count(searchQuery, currency string, categoryID uint, minPriceCents, maxPriceCents int64, active bool) (int, error) {
	return len(r.products), nil
}

// Create adds a product to the repository
func (r *MockProductRepository) Create(product *entity.Product) error {
	// Assign ID
	r.lastID++
	product.ID = r.lastID

	// Store product
	r.products[product.ID] = product

	return nil
}

// GetByID retrieves a product by ID
func (r *MockProductRepository) GetByID(id uint) (*entity.Product, error) {
	product, exists := r.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}
	return product, nil
}

// GetByIDWithVariants retrieves a product by ID including its variants
func (r *MockProductRepository) GetByIDWithVariants(id uint) (*entity.Product, error) {
	product, exists := r.products[id]
	if !exists {
		return nil, errors.New("product not found")
	}

	// Return a copy of the product to prevent unintended modifications
	productCopy := *product

	return &productCopy, nil
}

// Update updates a product
func (r *MockProductRepository) Update(product *entity.Product) error {
	if _, exists := r.products[product.ID]; !exists {
		return errors.New("product not found")
	}

	// Update product
	r.products[product.ID] = product

	return nil
}

// Delete removes a product
func (r *MockProductRepository) Delete(id uint) error {
	if _, exists := r.products[id]; !exists {
		return errors.New("product not found")
	}

	delete(r.products, id)
	return nil
}

// List retrieves products with pagination
func (r *MockProductRepository) List(query, currency string, categoryID, offset, limit uint, minPrice, maxPrice int64, active bool) ([]*entity.Product, error) {
	result := make([]*entity.Product, 0)
	count := uint(0)
	skip := offset

	for _, product := range r.products {
		// Apply search filters
		if query != "" && !strings.Contains(strings.ToLower(product.Name), strings.ToLower(query)) &&
			!strings.Contains(strings.ToLower(product.Description), strings.ToLower(query)) {
			continue
		}

		if categoryID > 0 && product.CategoryID != categoryID {
			continue
		}

		if minPrice > 0 && product.Price < minPrice {
			continue
		}

		if maxPrice > 0 && product.Price > maxPrice {
			continue
		}

		// Apply pagination
		if skip > 0 {
			skip--
			continue
		}

		result = append(result, product)
		count++

		if count >= limit {
			break
		}
	}

	return result, nil
}

// GetByProductNumber retrieves a product by product number (SKU)
func (r *MockProductRepository) GetByProductNumber(productNumber string) (*entity.Product, error) {
	for _, product := range r.products {
		if product.ProductNumber == productNumber {
			return product, nil
		}
	}
	return nil, errors.New("product not found")
}
