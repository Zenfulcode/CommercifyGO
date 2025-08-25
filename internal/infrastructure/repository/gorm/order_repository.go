package gorm

import (
	"errors"
	"fmt"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/dto"
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
	err := o.db.Table("order_items").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("order_items.product_id = ?", productID).
		Count(&count).Error
	if err != nil {
		return false, fmt.Errorf("failed to check orders with product %d: %w", productID, err)
	}
	return count > 0, nil
}

// IsDiscountIdUsed implements repository.OrderRepository.
func (o *OrderRepository) IsDiscountIdUsed(discountID uint) (bool, error) {
	var count int64
	err := o.db.Model(&entity.Order{}).
		Where("JSON_EXTRACT(applied_discount, '$.id') = ?", discountID).
		Count(&count).Error
	if err != nil {
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

// GetTotalRevenueByDateRange implements repository.OrderRepository.
func (o *OrderRepository) GetTotalRevenueByDateRange(startDate, endDate time.Time) (int64, error) {
	var totalRevenue int64
	err := o.db.Model(&entity.Order{}).
		Where("created_at >= ? AND created_at <= ? AND status IN (?)",
			startDate, endDate, []entity.OrderStatus{entity.OrderStatusPaid, entity.OrderStatusShipped, entity.OrderStatusCompleted}).
		Select("COALESCE(SUM(total_amount), 0)").
		Scan(&totalRevenue).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get total revenue: %w", err)
	}
	return totalRevenue, nil
}

// GetTotalOrdersByDateRange implements repository.OrderRepository.
func (o *OrderRepository) GetTotalOrdersByDateRange(startDate, endDate time.Time) (int64, error) {
	var count int64
	err := o.db.Model(&entity.Order{}).
		Where("created_at >= ? AND created_at <= ?", startDate, endDate).
		Count(&count).Error
	if err != nil {
		return 0, fmt.Errorf("failed to get total orders count: %w", err)
	}
	return count, nil
}

// GetRecentOrdersSummary implements repository.OrderRepository.
func (o *OrderRepository) GetRecentOrdersSummary(startDate, endDate time.Time, limit int) ([]dto.RecentOrderSummary, error) {
	var results []dto.RecentOrderSummary

	err := o.db.Table("orders").
		Select("orders.id, orders.order_number, orders.total_amount, orders.status, orders.created_at, "+
			"CASE WHEN orders.user_id IS NULL THEN 'Guest' ELSE (users.first_name || ' ' || users.last_name) END as customer_name, "+
			"CASE WHEN orders.user_id IS NULL THEN orders.customer_email ELSE users.email END as customer_email").
		Joins("LEFT JOIN users ON orders.user_id = users.id").
		Where("orders.created_at >= ? AND orders.created_at <= ?", startDate, endDate).
		Order("orders.created_at DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get recent orders summary: %w", err)
	}

	return results, nil
}

// GetTopProductsSummary implements repository.OrderRepository.
func (o *OrderRepository) GetTopProductsSummary(startDate, endDate time.Time, limit int) ([]dto.TopProductSummary, error) {
	var results []dto.TopProductSummary

	err := o.db.Table("order_items").
		Select("order_items.product_id, products.name as product_name, "+
			"order_items.product_variant_id, COALESCE(product_variants.sku, '') as variant_name, "+
			"SUM(order_items.quantity) as quantity_sold, "+
			"SUM(order_items.price * order_items.quantity) as revenue").
		Joins("JOIN products ON order_items.product_id = products.id").
		Joins("LEFT JOIN product_variants ON order_items.product_variant_id = product_variants.id").
		Joins("JOIN orders ON order_items.order_id = orders.id").
		Where("orders.created_at >= ? AND orders.created_at <= ? AND orders.status IN (?)",
			startDate, endDate, []entity.OrderStatus{entity.OrderStatusPaid, entity.OrderStatusShipped, entity.OrderStatusCompleted}).
		Group("order_items.product_id, products.name, order_items.product_variant_id, product_variants.sku").
		Order("quantity_sold DESC").
		Limit(limit).
		Scan(&results).Error

	if err != nil {
		return nil, fmt.Errorf("failed to get top products summary: %w", err)
	}

	return results, nil
}

// NewOrderRepository creates a new GORM-based OrderRepository
func NewOrderRepository(db *gorm.DB) repository.OrderRepository {
	return &OrderRepository{db: db}
}
