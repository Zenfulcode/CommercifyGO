package mock

import (
	"fmt"
	"sync"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// OrderRepository is a mock implementation of repository.OrderRepository
type OrderRepository struct {
	mu     sync.RWMutex
	orders map[uint]*entity.Order
	nextID uint
	// Maps for quick lookup
	checkoutSessionMap map[string]*entity.Order
	paymentIDMap       map[string]*entity.Order
	// For testing error scenarios
	CreateError                 error
	GetByIDError                error
	GetByCheckoutSessionIDError error
	UpdateError                 error
	GetByUserError              error
	ListByStatusError           error
	IsDiscountIdUsedError       error
	GetByPaymentIDError         error
	ListAllError                error
	HasOrdersWithProductError   error
}

// NewOrderRepository creates a new mock order repository
func NewOrderRepository() repository.OrderRepository {
	return &OrderRepository{
		orders:             make(map[uint]*entity.Order),
		checkoutSessionMap: make(map[string]*entity.Order),
		paymentIDMap:       make(map[string]*entity.Order),
		nextID:             1,
	}
}

// Create creates a new order
func (m *OrderRepository) Create(order *entity.Order) error {
	if m.CreateError != nil {
		return m.CreateError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Assign ID if not set
	if order.ID == 0 {
		order.ID = m.nextID
		m.nextID++
	}

	// Create a copy to avoid external mutations
	orderCopy := *order
	m.orders[orderCopy.ID] = &orderCopy

	// Update lookup maps
	if orderCopy.CheckoutSessionID != "" {
		m.checkoutSessionMap[orderCopy.CheckoutSessionID] = &orderCopy
	}
	if orderCopy.PaymentID != "" {
		m.paymentIDMap[orderCopy.PaymentID] = &orderCopy
	}

	return nil
}

// GetByID retrieves an order by ID
func (m *OrderRepository) GetByID(orderID uint) (*entity.Order, error) {
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	order, exists := m.orders[orderID]
	if !exists {
		return nil, fmt.Errorf("order with ID %d not found", orderID)
	}

	// Return a copy
	orderCopy := *order
	return &orderCopy, nil
}

// GetByCheckoutSessionID retrieves an order by checkout session ID
func (m *OrderRepository) GetByCheckoutSessionID(checkoutSessionID string) (*entity.Order, error) {
	if m.GetByCheckoutSessionIDError != nil {
		return nil, m.GetByCheckoutSessionIDError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	order, exists := m.checkoutSessionMap[checkoutSessionID]
	if !exists {
		return nil, fmt.Errorf("order with checkout session ID %s not found", checkoutSessionID)
	}

	// Return a copy
	orderCopy := *order
	return &orderCopy, nil
}

// Update updates an order
func (m *OrderRepository) Update(order *entity.Order) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	existing, exists := m.orders[order.ID]
	if !exists {
		return fmt.Errorf("order with ID %d not found", order.ID)
	}

	// Update lookup maps if these fields changed
	if existing.CheckoutSessionID != order.CheckoutSessionID {
		if existing.CheckoutSessionID != "" {
			delete(m.checkoutSessionMap, existing.CheckoutSessionID)
		}
		if order.CheckoutSessionID != "" {
			m.checkoutSessionMap[order.CheckoutSessionID] = order
		}
	}

	if existing.PaymentID != order.PaymentID {
		if existing.PaymentID != "" {
			delete(m.paymentIDMap, existing.PaymentID)
		}
		if order.PaymentID != "" {
			m.paymentIDMap[order.PaymentID] = order
		}
	}

	// Update the order
	orderCopy := *order
	m.orders[orderCopy.ID] = &orderCopy

	return nil
}

// GetByUser retrieves orders for a specific user with pagination
func (m *OrderRepository) GetByUser(userID uint, offset, limit int) ([]*entity.Order, error) {
	if m.GetByUserError != nil {
		return nil, m.GetByUserError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*entity.Order
	var count int
	for _, order := range m.orders {
		if order.UserID == userID {
			if count >= offset {
				if limit > 0 && len(results) >= limit {
					break
				}
				orderCopy := *order
				results = append(results, &orderCopy)
			}
			count++
		}
	}

	return results, nil
}

// ListByStatus retrieves orders by status with pagination
func (m *OrderRepository) ListByStatus(status entity.OrderStatus, offset, limit int) ([]*entity.Order, error) {
	if m.ListByStatusError != nil {
		return nil, m.ListByStatusError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*entity.Order
	var count int
	for _, order := range m.orders {
		if order.Status == status {
			if count >= offset {
				if limit > 0 && len(results) >= limit {
					break
				}
				orderCopy := *order
				results = append(results, &orderCopy)
			}
			count++
		}
	}

	return results, nil
}

// IsDiscountIdUsed checks if a discount ID is used in any order
func (m *OrderRepository) IsDiscountIdUsed(discountID uint) (bool, error) {
	if m.IsDiscountIdUsedError != nil {
		return false, m.IsDiscountIdUsedError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, order := range m.orders {
		appliedDiscount := order.GetAppliedDiscount()
		if appliedDiscount != nil && appliedDiscount.DiscountID == discountID {
			return true, nil
		}
	}

	return false, nil
}

// GetByPaymentID retrieves an order by payment ID
func (m *OrderRepository) GetByPaymentID(paymentID string) (*entity.Order, error) {
	if m.GetByPaymentIDError != nil {
		return nil, m.GetByPaymentIDError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	order, exists := m.paymentIDMap[paymentID]
	if !exists {
		return nil, fmt.Errorf("order with payment ID %s not found", paymentID)
	}

	// Return a copy
	orderCopy := *order
	return &orderCopy, nil
}

// ListAll retrieves all orders with pagination
func (m *OrderRepository) ListAll(offset, limit int) ([]*entity.Order, error) {
	if m.ListAllError != nil {
		return nil, m.ListAllError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*entity.Order
	var count int
	for _, order := range m.orders {
		if count >= offset {
			if limit > 0 && len(results) >= limit {
				break
			}
			orderCopy := *order
			results = append(results, &orderCopy)
		}
		count++
	}

	return results, nil
}

// HasOrdersWithProduct checks if there are any orders with the specified product
func (m *OrderRepository) HasOrdersWithProduct(productID uint) (bool, error) {
	if m.HasOrdersWithProductError != nil {
		return false, m.HasOrdersWithProductError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, order := range m.orders {
		for _, item := range order.Items {
			if item.ProductID == productID {
				return true, nil
			}
		}
	}

	return false, nil
}

// Helper methods for testing

// SetCreateError sets an error to be returned by Create
func (m *OrderRepository) SetCreateError(err error) {
	m.CreateError = err
}

// SetGetByIDError sets an error to be returned by GetByID
func (m *OrderRepository) SetGetByIDError(err error) {
	m.GetByIDError = err
}

// SetGetByCheckoutSessionIDError sets an error to be returned by GetByCheckoutSessionID
func (m *OrderRepository) SetGetByCheckoutSessionIDError(err error) {
	m.GetByCheckoutSessionIDError = err
}

// SetUpdateError sets an error to be returned by Update
func (m *OrderRepository) SetUpdateError(err error) {
	m.UpdateError = err
}

// SetGetByUserError sets an error to be returned by GetByUser
func (m *OrderRepository) SetGetByUserError(err error) {
	m.GetByUserError = err
}

// SetListByStatusError sets an error to be returned by ListByStatus
func (m *OrderRepository) SetListByStatusError(err error) {
	m.ListByStatusError = err
}

// SetIsDiscountIdUsedError sets an error to be returned by IsDiscountIdUsed
func (m *OrderRepository) SetIsDiscountIdUsedError(err error) {
	m.IsDiscountIdUsedError = err
}

// SetGetByPaymentIDError sets an error to be returned by GetByPaymentID
func (m *OrderRepository) SetGetByPaymentIDError(err error) {
	m.GetByPaymentIDError = err
}

// SetListAllError sets an error to be returned by ListAll
func (m *OrderRepository) SetListAllError(err error) {
	m.ListAllError = err
}

// SetHasOrdersWithProductError sets an error to be returned by HasOrdersWithProduct
func (m *OrderRepository) SetHasOrdersWithProductError(err error) {
	m.HasOrdersWithProductError = err
}

// Reset clears all data and errors
func (m *OrderRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.orders = make(map[uint]*entity.Order)
	m.checkoutSessionMap = make(map[string]*entity.Order)
	m.paymentIDMap = make(map[string]*entity.Order)
	m.nextID = 1
	m.CreateError = nil
	m.GetByIDError = nil
	m.GetByCheckoutSessionIDError = nil
	m.UpdateError = nil
	m.GetByUserError = nil
	m.ListByStatusError = nil
	m.IsDiscountIdUsedError = nil
	m.GetByPaymentIDError = nil
	m.ListAllError = nil
	m.HasOrdersWithProductError = nil
}

// GetAllOrders returns all orders (for testing purposes)
func (m *OrderRepository) GetAllOrders() []*entity.Order {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*entity.Order
	for _, order := range m.orders {
		orderCopy := *order
		results = append(results, &orderCopy)
	}

	return results
}

// AddTestOrder adds a test order (for testing purposes)
func (m *OrderRepository) AddTestOrder(order *entity.Order) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if order.ID == 0 {
		order.ID = m.nextID
		m.nextID++
	}

	orderCopy := *order
	m.orders[orderCopy.ID] = &orderCopy

	// Update lookup maps
	if orderCopy.CheckoutSessionID != "" {
		m.checkoutSessionMap[orderCopy.CheckoutSessionID] = &orderCopy
	}
	if orderCopy.PaymentID != "" {
		m.paymentIDMap[orderCopy.PaymentID] = &orderCopy
	}
}
