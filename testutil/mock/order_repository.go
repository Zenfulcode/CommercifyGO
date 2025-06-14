package mock

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// OrderRepository is a mock implementation of the order repository interface
type OrderRepository struct {
	orders           map[uint]*entity.Order
	paymentIDIndex   map[string]*entity.Order // Index to find orders by payment ID
	isDiscountIdUsed bool
}

// NewMockOrderRepository creates a new mock order repository
func NewMockOrderRepository(
	isDiscountIdUsed bool,
) repository.OrderRepository {
	return &OrderRepository{
		orders:           make(map[uint]*entity.Order),
		paymentIDIndex:   make(map[string]*entity.Order),
		isDiscountIdUsed: isDiscountIdUsed,
	}
}

// ListAll implements repository.OrderRepository.
func (r *OrderRepository) ListAll(offset int, limit int) ([]*entity.Order, error) {
	orders := make([]*entity.Order, 0, len(r.orders))
	for _, order := range r.orders {
		orders = append(orders, order)
	}
	return orders, nil
}

// Create adds a new order to the mock repository
func (r *OrderRepository) Create(order *entity.Order) error {
	// If no ID provided, generate one
	if order.ID == 0 {
		maxID := uint(0)
		for id := range r.orders {
			if id > maxID {
				maxID = id
			}
		}
		order.ID = maxID + 1
	}

	// Clone the order to prevent unintended modifications
	clone := *order
	r.orders[order.ID] = &clone

	// Index by payment ID if available
	if order.PaymentID != "" {
		r.paymentIDIndex[order.PaymentID] = &clone
	}

	return nil
}

// GetByID retrieves an order by ID from the mock repository
func (r *OrderRepository) GetByID(id uint) (*entity.Order, error) {
	order, exists := r.orders[id]
	if !exists {
		return nil, errors.New("order not found")
	}

	// Return a clone to prevent unintended modifications
	clone := *order
	return &clone, nil
}

// Update updates an existing order in the mock repository
func (r *OrderRepository) Update(order *entity.Order) error {
	if _, exists := r.orders[order.ID]; !exists {
		return errors.New("order not found")
	}

	// If payment ID has changed, update the index
	existingOrder := r.orders[order.ID]
	if existingOrder.PaymentID != order.PaymentID {
		if existingOrder.PaymentID != "" {
			delete(r.paymentIDIndex, existingOrder.PaymentID)
		}
		if order.PaymentID != "" {
			r.paymentIDIndex[order.PaymentID] = order
		}
	}

	// Clone the order to prevent unintended modifications
	clone := *order
	r.orders[order.ID] = &clone

	return nil
}

// GetByUser retrieves orders for a user from the mock repository
func (r *OrderRepository) GetByUser(userID uint, offset, limit int) ([]*entity.Order, error) {
	var orders []*entity.Order
	for _, order := range r.orders {
		if order.UserID == userID {
			clone := *order
			orders = append(orders, &clone)
		}
	}

	// Apply offset and limit
	if offset >= len(orders) {
		return []*entity.Order{}, nil
	}
	end := min(offset+limit, len(orders))

	return orders[offset:end], nil
}

// ListByStatus retrieves orders by status from the mock repository
func (r *OrderRepository) ListByStatus(status entity.OrderStatus, offset, limit int) ([]*entity.Order, error) {
	var orders []*entity.Order
	for _, order := range r.orders {
		if order.Status == status {
			clone := *order
			orders = append(orders, &clone)
		}
	}

	// Apply offset and limit
	if offset >= len(orders) {
		return []*entity.Order{}, nil
	}
	end := min(offset+limit, len(orders))

	return orders[offset:end], nil
}

func (r *OrderRepository) SetIsDiscountIdUsed(isDiscountIdUsed bool) {
	r.isDiscountIdUsed = isDiscountIdUsed
}

// IsDiscountIdUsed checks if a discount is used by any order in the mock repository
func (r *OrderRepository) IsDiscountIdUsed(discountID uint) (bool, error) {
	if r.isDiscountIdUsed {
		return true, nil
	}

	// Otherwise fall back to the default implementation
	for _, order := range r.orders {
		if order.AppliedDiscount != nil && order.AppliedDiscount.DiscountID == discountID {
			return true, nil
		}
	}
	return false, nil
}

// GetByPaymentID retrieves an order by payment ID from the mock repository
func (r *OrderRepository) GetByPaymentID(paymentID string) (*entity.Order, error) {
	order, exists := r.paymentIDIndex[paymentID]
	if !exists {
		return nil, errors.New("order not found for payment ID")
	}

	// Return a clone to prevent unintended modifications
	clone := *order
	return &clone, nil
}

// AddMockGetByPaymentID is a helper function to set up mock behavior for GetByPaymentID
func (r *OrderRepository) AddMockGetByPaymentID(order *entity.Order) {
	if order != nil && order.PaymentID != "" {
		r.paymentIDIndex[order.PaymentID] = order
	}
}

// HasOrdersWithProduct checks if any orders contain items with the specified product ID
func (r *OrderRepository) HasOrdersWithProduct(productID uint) (bool, error) {
	if productID == 0 {
		return false, errors.New("product ID cannot be 0")
	}

	for _, order := range r.orders {
		for _, item := range order.Items {
			if item.ProductID == productID {
				return true, nil
			}
		}
	}

	return false, nil
}
