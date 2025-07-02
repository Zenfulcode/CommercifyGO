package gorm

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/gorm"
)

// TransactionRepository implements repository.TransactionRepository using GORM
type TransactionRepository struct {
	db *gorm.DB
}

// CountSuccessfulByOrderIDAndType implements repository.PaymentTransactionRepository.
func (t *TransactionRepository) CountSuccessfulByOrderIDAndType(orderID uint, transactionType entity.TransactionType) (int, error) {
	var count int64
	if err := t.db.Model(&entity.PaymentTransaction{}).
		Where("order_id = ? AND type = ? AND status = ?", orderID, transactionType, entity.TransactionStatusSuccessful).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count successful payment transactions: %w", err)
	}
	return int(count), nil
}

// Create implements repository.PaymentTransactionRepository.
func (t *TransactionRepository) Create(transaction *entity.PaymentTransaction) error {
	// Always create a new transaction record (no upsert behavior)
	// This allows multiple transactions of the same type for the same order
	// which is useful for scenarios like partial captures, webhook retries, etc.
	if transaction.TransactionID == "" {
		sequence, err := t.getNextSequenceNumber(transaction.Type)
		if err != nil {
			return fmt.Errorf("failed to generate sequence number: %w", err)
		}
		transaction.SetTransactionID(sequence)
	}
	return t.db.Create(transaction).Error
}

// CreateOrUpdate creates a new transaction or updates an existing one if a transaction
// of the same type already exists for the order. This method implements upsert behavior
// for cases where you want to ensure only one transaction per type per order.
func (t *TransactionRepository) CreateOrUpdate(transaction *entity.PaymentTransaction) error {
	// Check if a transaction of this type already exists for this order
	var existingTransaction entity.PaymentTransaction
	err := t.db.Where("order_id = ? AND type = ?", transaction.OrderID, transaction.Type).
		First(&existingTransaction).Error

	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check for existing transaction: %w", err)
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// No existing transaction, create a new one
		return t.Create(transaction)
	} else {
		// Transaction exists, update it with new information
		existingTransaction.Status = transaction.Status
		existingTransaction.Amount = transaction.Amount
		existingTransaction.ExternalID = transaction.ExternalID
		existingTransaction.RawResponse = transaction.RawResponse
		existingTransaction.Metadata = transaction.Metadata

		// Update the transaction in the database
		err = t.db.Save(&existingTransaction).Error
		if err != nil {
			return fmt.Errorf("failed to update existing transaction: %w", err)
		}

		// Copy the updated values back to the input transaction for consistency
		*transaction = existingTransaction
		return nil
	}
}

// Delete implements repository.PaymentTransactionRepository.
func (t *TransactionRepository) Delete(id uint) error {
	return t.db.Delete(&entity.PaymentTransaction{}, id).Error
}

// GetByID implements repository.PaymentTransactionRepository.
func (t *TransactionRepository) GetByID(id uint) (*entity.PaymentTransaction, error) {
	var transaction entity.PaymentTransaction
	if err := t.db.Preload("Order").First(&transaction, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("payment transaction with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch payment transaction: %w", err)
	}
	return &transaction, nil
}

// GetByOrderID implements repository.PaymentTransactionRepository.
func (t *TransactionRepository) GetByOrderID(orderID uint) ([]*entity.PaymentTransaction, error) {
	var transactions []*entity.PaymentTransaction
	if err := t.db.Preload("Order").Where("order_id = ?", orderID).
		Order("created_at DESC").Find(&transactions).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch payment transactions for order %d: %w", orderID, err)
	}
	return transactions, nil
}

// GetByTransactionID implements repository.PaymentTransactionRepository.
func (t *TransactionRepository) GetByTransactionID(transactionID string) (*entity.PaymentTransaction, error) {
	var transaction entity.PaymentTransaction
	if err := t.db.Preload("Order").Where("transaction_id = ?", transactionID).First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("payment transaction with transaction ID %s not found", transactionID)
		}
		return nil, fmt.Errorf("failed to fetch payment transaction by transaction ID: %w", err)
	}
	return &transaction, nil
}

// GetLatestByOrderIDAndType implements repository.PaymentTransactionRepository.
func (t *TransactionRepository) GetLatestByOrderIDAndType(orderID uint, transactionType entity.TransactionType) (*entity.PaymentTransaction, error) {
	var transaction entity.PaymentTransaction
	if err := t.db.Preload("Order").
		Where("order_id = ? AND type = ?", orderID, transactionType).
		Order("created_at DESC").
		First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no payment transaction of type %s found for order %d", transactionType, orderID)
		}
		return nil, fmt.Errorf("failed to fetch latest payment transaction: %w", err)
	}
	return &transaction, nil
}

// SumAmountByOrderIDAndType implements repository.PaymentTransactionRepository.
func (t *TransactionRepository) SumAmountByOrderIDAndType(orderID uint, transactionType entity.TransactionType) (int64, error) {
	var result struct {
		TotalAmount int64
	}
	if err := t.db.Model(&entity.PaymentTransaction{}).
		Select("COALESCE(SUM(amount), 0) as total_amount").
		Where("order_id = ? AND type = ? AND status = ?", orderID, transactionType, entity.TransactionStatusSuccessful).
		Scan(&result).Error; err != nil {
		return 0, fmt.Errorf("failed to sum payment transaction amounts: %w", err)
	}
	return result.TotalAmount, nil
}

// SumAuthorizedAmountByOrderID implements repository.PaymentTransactionRepository.
func (t *TransactionRepository) SumAuthorizedAmountByOrderID(orderID uint) (int64, error) {
	var result struct {
		TotalAmount int64
	}
	if err := t.db.Model(&entity.PaymentTransaction{}).
		Select("COALESCE(SUM(authorized_amount), 0) as total_amount").
		Where("order_id = ? AND status = ?", orderID, entity.TransactionStatusSuccessful).
		Scan(&result).Error; err != nil {
		return 0, fmt.Errorf("failed to sum authorized amounts: %w", err)
	}
	return result.TotalAmount, nil
}

// SumCapturedAmountByOrderID implements repository.PaymentTransactionRepository.
func (t *TransactionRepository) SumCapturedAmountByOrderID(orderID uint) (int64, error) {
	var result struct {
		TotalAmount int64
	}
	if err := t.db.Model(&entity.PaymentTransaction{}).
		Select("COALESCE(SUM(captured_amount), 0) as total_amount").
		Where("order_id = ? AND status = ?", orderID, entity.TransactionStatusSuccessful).
		Scan(&result).Error; err != nil {
		return 0, fmt.Errorf("failed to sum captured amounts: %w", err)
	}
	return result.TotalAmount, nil
}

// SumRefundedAmountByOrderID implements repository.PaymentTransactionRepository.
func (t *TransactionRepository) SumRefundedAmountByOrderID(orderID uint) (int64, error) {
	var result struct {
		TotalAmount int64
	}
	if err := t.db.Model(&entity.PaymentTransaction{}).
		Select("COALESCE(SUM(refunded_amount), 0) as total_amount").
		Where("order_id = ? AND status = ?", orderID, entity.TransactionStatusSuccessful).
		Scan(&result).Error; err != nil {
		return 0, fmt.Errorf("failed to sum refunded amounts: %w", err)
	}
	return result.TotalAmount, nil
}

// Update implements repository.PaymentTransactionRepository.
func (t *TransactionRepository) Update(transaction *entity.PaymentTransaction) error {
	return t.db.Save(transaction).Error
}

// NewTransactionRepository creates a new GORM-based TransactionRepository
func NewTransactionRepository(db *gorm.DB) repository.PaymentTransactionRepository {
	return &TransactionRepository{db: db}
}

// getNextSequenceNumber generates the next sequence number for a given transaction type and year
func (t *TransactionRepository) getNextSequenceNumber(transactionType entity.TransactionType) (int, error) {
	var count int64

	// Count existing transactions of this type for the current year
	// This creates a sequence like: TXN-AUTH-2025-001, TXN-AUTH-2025-002, etc.
	year := time.Now().Year()

	// Count transactions with IDs matching the pattern for this type and year
	if err := t.db.Model(&entity.PaymentTransaction{}).
		Where("transaction_id LIKE ?", fmt.Sprintf("TXN-%s-%d-%%", getTypeCode(transactionType), year)).
		Count(&count).Error; err != nil {
		return 0, fmt.Errorf("failed to count existing transactions: %w", err)
	}

	return int(count) + 1, nil
}

// getTypeCode returns the short type code for transaction types
func getTypeCode(transactionType entity.TransactionType) string {
	switch transactionType {
	case entity.TransactionTypeAuthorize:
		return "AUTH"
	case entity.TransactionTypeCapture:
		return "CAPT"
	case entity.TransactionTypeRefund:
		return "REFUND"
	case entity.TransactionTypeCancel:
		return "CANCEL"
	default:
		return strings.ToUpper(string(transactionType))
	}
}

// GetByIdempotencyKey implements repository.PaymentTransactionRepository.
func (t *TransactionRepository) GetByIdempotencyKey(idempotencyKey string) (*entity.PaymentTransaction, error) {
	var transaction entity.PaymentTransaction

	// Search for transactions by the dedicated idempotency_key field
	if err := t.db.Preload("Order").
		Where("idempotency_key = ?", idempotencyKey).
		First(&transaction).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("payment transaction with idempotency key %s not found", idempotencyKey)
		}
		return nil, fmt.Errorf("failed to fetch payment transaction by idempotency key: %w", err)
	}
	return &transaction, nil
}
