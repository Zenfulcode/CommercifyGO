package gorm

import (
	"errors"
	"fmt"

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
	return t.db.Create(transaction).Error
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

// Update implements repository.PaymentTransactionRepository.
func (t *TransactionRepository) Update(transaction *entity.PaymentTransaction) error {
	return t.db.Save(transaction).Error
}

// NewTransactionRepository creates a new GORM-based TransactionRepository
func NewTransactionRepository(db *gorm.DB) repository.PaymentTransactionRepository {
	return &TransactionRepository{db: db}
}
