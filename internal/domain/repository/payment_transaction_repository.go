package repository

import (
	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// PaymentTransactionRepository defines the interface for payment transaction persistence
type PaymentTransactionRepository interface {
	// Create creates a new payment transaction (always creates a new record)
	Create(transaction *entity.PaymentTransaction) error

	// CreateOrUpdate creates a new transaction or updates an existing one if a transaction
	// of the same type already exists for the order (upsert behavior)
	CreateOrUpdate(transaction *entity.PaymentTransaction) error

	// GetByID retrieves a payment transaction by ID
	GetByID(id uint) (*entity.PaymentTransaction, error)

	// GetByTransactionID retrieves a payment transaction by external transaction ID
	GetByTransactionID(transactionID string) (*entity.PaymentTransaction, error)

	// GetByOrderID retrieves all payment transactions for an order
	GetByOrderID(orderID uint) ([]*entity.PaymentTransaction, error)

	// Update updates a payment transaction
	Update(transaction *entity.PaymentTransaction) error

	// Delete deletes a payment transaction
	Delete(id uint) error

	// GetLatestByOrderIDAndType retrieves the latest transaction of a specific type for an order
	GetLatestByOrderIDAndType(orderID uint, transactionType entity.TransactionType) (*entity.PaymentTransaction, error)

	// CountSuccessfulByOrderIDAndType counts successful transactions of a specific type for an order
	CountSuccessfulByOrderIDAndType(orderID uint, transactionType entity.TransactionType) (int, error)

	// SumAmountByOrderIDAndType sums the amount of transactions of a specific type for an order
	SumAmountByOrderIDAndType(orderID uint, transactionType entity.TransactionType) (int64, error)

	// SumAuthorizedAmountByOrderID sums all authorized amounts for an order
	SumAuthorizedAmountByOrderID(orderID uint) (int64, error)

	// SumCapturedAmountByOrderID sums all captured amounts for an order
	SumCapturedAmountByOrderID(orderID uint) (int64, error)

	// SumRefundedAmountByOrderID sums all refunded amounts for an order
	SumRefundedAmountByOrderID(orderID uint) (int64, error)

	// GetByIdempotencyKey retrieves a payment transaction by idempotency key from metadata
	GetByIdempotencyKey(idempotencyKey string) (*entity.PaymentTransaction, error)
}
