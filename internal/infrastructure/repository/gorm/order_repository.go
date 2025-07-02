package gorm

import (
	"errors"
	"fmt"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/gorm"
)

// OrderRepository implements repository.OrderRepository using GORM
type OrderRepository struct {
	db *gorm.DB
}

// Create implements repository.OrderRepository.
func (o *OrderRepository) Create(order *entity.Order) error {
	return o.db.Create(order).Error
}

// GetByCheckoutSessionID implements repository.OrderRepository.
func (o *OrderRepository) GetByCheckoutSessionID(checkoutSessionID string) (*entity.Order, error) {
	var order entity.Order
	if err := o.db.Preload("Items").Preload("Items.Product").Preload("Items.ProductVariant").
		Preload("User").Preload("PaymentTransactions").
		Where("checkout_session_id = ?", checkoutSessionID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order with checkout session ID %s not found", checkoutSessionID)
		}
		return nil, fmt.Errorf("failed to fetch order by checkout session ID: %w", err)
	}
	return &order, nil
}

// GetByID implements repository.OrderRepository.
func (o *OrderRepository) GetByID(orderID uint) (*entity.Order, error) {
	var order entity.Order
	if err := o.db.Preload("Items").Preload("Items.Product").Preload("Items.ProductVariant").
		Preload("User").Preload("PaymentTransactions").
		First(&order, orderID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order with ID %d not found", orderID)
		}
		return nil, fmt.Errorf("failed to fetch order: %w", err)
	}
	return &order, nil
}

// GetByPaymentID implements repository.OrderRepository.
func (o *OrderRepository) GetByPaymentID(paymentID string) (*entity.Order, error) {
	var order entity.Order
	if err := o.db.Preload("Items").Preload("Items.Product").Preload("Items.ProductVariant").
		Preload("User").Preload("PaymentTransactions").
		Where("payment_id = ?", paymentID).First(&order).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("order with payment ID %s not found", paymentID)
		}
		return nil, fmt.Errorf("failed to fetch order by payment ID: %w", err)
	}
	return &order, nil
}

// GetByUser implements repository.OrderRepository.
func (o *OrderRepository) GetByUser(userID uint, offset int, limit int) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := o.db.Preload("Items").Preload("Items.Product").Preload("Items.ProductVariant").
		Preload("User").
		Where("user_id = ?", userID).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch orders for user %d: %w", userID, err)
	}
	return orders, nil
}

// HasOrdersWithProduct implements repository.OrderRepository.
func (o *OrderRepository) HasOrdersWithProduct(productID uint) (bool, error) {
	var count int64
	if err := o.db.Model(&entity.Order{}).
		Joins("JOIN order_items ON orders.id = order_items.order_id").
		Where("order_items.product_id = ?", productID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check orders with product %d: %w", productID, err)
	}
	return count > 0, nil
}

// IsDiscountIdUsed implements repository.OrderRepository.
func (o *OrderRepository) IsDiscountIdUsed(discountID uint) (bool, error) {
	var count int64
	if err := o.db.Model(&entity.Order{}).
		Where("discount_discount_id = ?", discountID).
		Count(&count).Error; err != nil {
		return false, fmt.Errorf("failed to check if discount %d is used: %w", discountID, err)
	}
	return count > 0, nil
}

// ListAll implements repository.OrderRepository.
func (o *OrderRepository) ListAll(offset int, limit int) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := o.db.Preload("Items").Preload("Items.Product").Preload("Items.ProductVariant").
		Preload("User").
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch all orders: %w", err)
	}
	return orders, nil
}

// ListByStatus implements repository.OrderRepository.
func (o *OrderRepository) ListByStatus(status entity.OrderStatus, offset int, limit int) ([]*entity.Order, error) {
	var orders []*entity.Order
	if err := o.db.Preload("Items").Preload("Items.Product").Preload("Items.ProductVariant").
		Preload("User").
		Where("status = ?", status).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&orders).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch orders by status %s: %w", status, err)
	}
	return orders, nil
}

// Update implements repository.OrderRepository.
func (o *OrderRepository) Update(order *entity.Order) error {
	return o.db.Session(&gorm.Session{FullSaveAssociations: true}).Save(order).Error
}

// NewOrderRepository creates a new GORM-based OrderRepository
func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &OrderRepository{db: db}
}
