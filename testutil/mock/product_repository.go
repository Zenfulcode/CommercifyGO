package mock

import (
	"fmt"
	"strings"
	"sync"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// ProductRepository is a mock implementation of repository.ProductRepository
type ProductRepository struct {
	mu       sync.RWMutex
	products map[uint]*entity.Product
	nextID   uint
	// Map SKU to entity for quick lookup
	skuMap map[string]*entity.Product
	// For testing error scenarios
	CreateError             error
	GetByIDError            error
	GetByIDAndCurrencyError error
	GetBySKUError           error
	UpdateError             error
	DeleteError             error
	ListError               error
	CountError              error
}

// NewProductRepository creates a new mock product repository
func NewProductRepository() repository.ProductRepository {
	return &ProductRepository{
		products: make(map[uint]*entity.Product),
		skuMap:   make(map[string]*entity.Product),
		nextID:   1,
	}
}

// Create creates a new product
func (m *ProductRepository) Create(product *entity.Product) error {
	if m.CreateError != nil {
		return m.CreateError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if any variant SKU already exists
	for _, variant := range product.Variants {
		if _, exists := m.skuMap[variant.SKU]; exists {
			return fmt.Errorf("product variant with SKU %s already exists", variant.SKU)
		}
	}

	// Assign ID if not set
	if product.ID == 0 {
		product.ID = m.nextID
		m.nextID++
	}

	// Assign IDs to variants if not set
	for _, variant := range product.Variants {
		if variant.ID == 0 {
			variant.ID = m.nextID
			m.nextID++
		}
		variant.ProductID = product.ID
	}

	// Create a copy to avoid external mutations
	productCopy := *product
	// Deep copy variants
	productCopy.Variants = make([]*entity.ProductVariant, len(product.Variants))
	for i, variant := range product.Variants {
		variantCopy := *variant
		productCopy.Variants[i] = &variantCopy
		// Update SKU mapping
		m.skuMap[variantCopy.SKU] = &productCopy
	}

	m.products[productCopy.ID] = &productCopy

	return nil
}

// GetByID retrieves a product by ID
func (m *ProductRepository) GetByID(productID uint) (*entity.Product, error) {
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	product, exists := m.products[productID]
	if !exists {
		return nil, fmt.Errorf("product with ID %d not found", productID)
	}

	// Return a deep copy
	return m.copyProduct(product), nil
}

// GetByIDAndCurrency retrieves a product by ID and currency
func (m *ProductRepository) GetByIDAndCurrency(productID uint, currency string) (*entity.Product, error) {
	if m.GetByIDAndCurrencyError != nil {
		return nil, m.GetByIDAndCurrencyError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	product, exists := m.products[productID]
	if !exists {
		return nil, fmt.Errorf("product with ID %d not found", productID)
	}

	if product.Currency != currency {
		return nil, fmt.Errorf("product with ID %d not found for currency %s", productID, currency)
	}

	// Return a deep copy
	return m.copyProduct(product), nil
}

// GetBySKU retrieves a product by variant SKU
func (m *ProductRepository) GetBySKU(sku string) (*entity.Product, error) {
	if m.GetBySKUError != nil {
		return nil, m.GetBySKUError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	product, exists := m.skuMap[sku]
	if !exists {
		return nil, fmt.Errorf("product with SKU %s not found", sku)
	}

	// Return a deep copy
	return m.copyProduct(product), nil
}

// Update updates a product
func (m *ProductRepository) Update(product *entity.Product) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	existing, exists := m.products[product.ID]
	if !exists {
		return fmt.Errorf("product with ID %d not found", product.ID)
	}

	// Update SKU mappings
	// First remove old SKU mappings
	for _, variant := range existing.Variants {
		delete(m.skuMap, variant.SKU)
	}

	// Check for SKU conflicts with other products
	for _, variant := range product.Variants {
		if existingProduct, exists := m.skuMap[variant.SKU]; exists && existingProduct.ID != product.ID {
			return fmt.Errorf("product variant with SKU %s already exists in another product", variant.SKU)
		}
	}

	// Update the product
	productCopy := *product
	// Deep copy variants
	productCopy.Variants = make([]*entity.ProductVariant, len(product.Variants))
	for i, variant := range product.Variants {
		variantCopy := *variant
		variantCopy.ProductID = product.ID
		productCopy.Variants[i] = &variantCopy
		// Update SKU mapping
		m.skuMap[variantCopy.SKU] = &productCopy
	}

	m.products[productCopy.ID] = &productCopy

	return nil
}

// Delete deletes a product
func (m *ProductRepository) Delete(productID uint) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	product, exists := m.products[productID]
	if !exists {
		return fmt.Errorf("product with ID %d not found", productID)
	}

	// Remove SKU mappings
	for _, variant := range product.Variants {
		delete(m.skuMap, variant.SKU)
	}

	delete(m.products, productID)

	return nil
}

// List retrieves products with filtering and pagination
func (m *ProductRepository) List(query, currency string, categoryID, offset, limit uint, minPriceCents, maxPriceCents int64, active bool) ([]*entity.Product, error) {
	if m.ListError != nil {
		return nil, m.ListError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*entity.Product
	var count uint

	for _, product := range m.products {
		// Apply filters
		if !m.matchesFilters(product, query, currency, categoryID, minPriceCents, maxPriceCents, active) {
			continue
		}

		// Apply pagination
		if count < offset {
			count++
			continue
		}

		if limit > 0 && uint(len(results)) >= limit {
			break
		}

		results = append(results, m.copyProduct(product))
		count++
	}

	return results, nil
}

// Count counts products matching the filters
func (m *ProductRepository) Count(searchQuery, currency string, categoryID uint, minPriceCents, maxPriceCents int64, active bool) (int, error) {
	if m.CountError != nil {
		return 0, m.CountError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, product := range m.products {
		if m.matchesFilters(product, searchQuery, currency, categoryID, minPriceCents, maxPriceCents, active) {
			count++
		}
	}

	return count, nil
}

// Helper methods

// copyProduct creates a deep copy of a product
func (m *ProductRepository) copyProduct(product *entity.Product) *entity.Product {
	productCopy := *product

	// Deep copy variants
	productCopy.Variants = make([]*entity.ProductVariant, len(product.Variants))
	for i, variant := range product.Variants {
		variantCopy := *variant
		productCopy.Variants[i] = &variantCopy
	}

	// Deep copy images slice
	if product.Images != nil {
		productCopy.Images = make([]string, len(product.Images))
		copy(productCopy.Images, product.Images)
	}

	return &productCopy
}

// matchesFilters checks if a product matches the given filters
func (m *ProductRepository) matchesFilters(product *entity.Product, query, currency string, categoryID uint, minPriceCents, maxPriceCents int64, active bool) bool {
	// Active filter
	if product.Active != active {
		return false
	}

	// Currency filter
	if currency != "" && product.Currency != currency {
		return false
	}

	// Category filter
	if categoryID > 0 && product.CategoryID != categoryID {
		return false
	}

	// Query filter (search in name and description)
	if query != "" {
		query = strings.ToLower(query)
		nameMatch := strings.Contains(strings.ToLower(product.Name), query)
		descMatch := strings.Contains(strings.ToLower(product.Description), query)
		if !nameMatch && !descMatch {
			return false
		}
	}

	// Price range filter (check variants)
	if minPriceCents > 0 || maxPriceCents > 0 {
		hasVariantInRange := false
		for _, variant := range product.Variants {
			if minPriceCents > 0 && variant.Price < minPriceCents {
				continue
			}
			if maxPriceCents > 0 && variant.Price > maxPriceCents {
				continue
			}
			hasVariantInRange = true
			break
		}
		if !hasVariantInRange {
			return false
		}
	}

	return true
}

// Helper methods for testing

// SetCreateError sets an error to be returned by Create
func (m *ProductRepository) SetCreateError(err error) {
	m.CreateError = err
}

// SetGetByIDError sets an error to be returned by GetByID
func (m *ProductRepository) SetGetByIDError(err error) {
	m.GetByIDError = err
}

// SetGetByIDAndCurrencyError sets an error to be returned by GetByIDAndCurrency
func (m *ProductRepository) SetGetByIDAndCurrencyError(err error) {
	m.GetByIDAndCurrencyError = err
}

// SetGetBySKUError sets an error to be returned by GetBySKU
func (m *ProductRepository) SetGetBySKUError(err error) {
	m.GetBySKUError = err
}

// SetUpdateError sets an error to be returned by Update
func (m *ProductRepository) SetUpdateError(err error) {
	m.UpdateError = err
}

// SetDeleteError sets an error to be returned by Delete
func (m *ProductRepository) SetDeleteError(err error) {
	m.DeleteError = err
}

// SetListError sets an error to be returned by List
func (m *ProductRepository) SetListError(err error) {
	m.ListError = err
}

// SetCountError sets an error to be returned by Count
func (m *ProductRepository) SetCountError(err error) {
	m.CountError = err
}

// Reset clears all data and errors
func (m *ProductRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.products = make(map[uint]*entity.Product)
	m.skuMap = make(map[string]*entity.Product)
	m.nextID = 1
	m.CreateError = nil
	m.GetByIDError = nil
	m.GetByIDAndCurrencyError = nil
	m.GetBySKUError = nil
	m.UpdateError = nil
	m.DeleteError = nil
	m.ListError = nil
	m.CountError = nil
}

// GetAllProducts returns all products (for testing purposes)
func (m *ProductRepository) GetAllProducts() []*entity.Product {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*entity.Product
	for _, product := range m.products {
		results = append(results, m.copyProduct(product))
	}

	return results
}

// AddTestProduct adds a test product (for testing purposes)
func (m *ProductRepository) AddTestProduct(product *entity.Product) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if product.ID == 0 {
		product.ID = m.nextID
		m.nextID++
	}

	// Assign IDs to variants if not set
	for _, variant := range product.Variants {
		if variant.ID == 0 {
			variant.ID = m.nextID
			m.nextID++
		}
		variant.ProductID = product.ID
		m.skuMap[variant.SKU] = product
	}

	productCopy := m.copyProduct(product)
	m.products[productCopy.ID] = productCopy
}

// GetProductBySKU returns the product containing the variant with the given SKU (for testing purposes)
func (m *ProductRepository) GetProductBySKU(sku string) *entity.Product {
	m.mu.RLock()
	defer m.mu.RUnlock()

	product, exists := m.skuMap[sku]
	if !exists {
		return nil
	}

	return m.copyProduct(product)
}
