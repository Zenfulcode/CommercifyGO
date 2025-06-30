package mock

import (
	"errors"
	"fmt"
	"sync"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// PaymentTransactionRepository is a mock implementation of repository.PaymentTransactionRepository
type PaymentTransactionRepository struct {
	mu           sync.RWMutex
	transactions map[uint]*entity.PaymentTransaction
	nextID       uint
	// Map transaction ID to entity for quick lookup
	transactionIDMap map[string]*entity.PaymentTransaction
	// For testing error scenarios
	CreateError                          error
	GetByIDError                         error
	GetByTransactionIDError              error
	GetByOrderIDError                    error
	UpdateError                          error
	DeleteError                          error
	GetLatestByOrderIDAndTypeError       error
	CountSuccessfulByOrderIDAndTypeError error
	SumAmountByOrderIDAndTypeError       error
}

// NewPaymentTransactionRepository creates a new mock payment transaction repository
func NewPaymentTransactionRepository() repository.PaymentTransactionRepository {
	return &PaymentTransactionRepository{
		transactions:     make(map[uint]*entity.PaymentTransaction),
		transactionIDMap: make(map[string]*entity.PaymentTransaction),
		nextID:           1,
	}
}

// Create creates a new payment transaction
func (m *PaymentTransactionRepository) Create(transaction *entity.PaymentTransaction) error {
	if m.CreateError != nil {
		return m.CreateError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	// Check if transaction ID already exists
	if _, exists := m.transactionIDMap[transaction.TransactionID]; exists {
		return errors.New("transaction with this ID already exists")
	}

	// Assign ID if not set
	if transaction.ID == 0 {
		transaction.ID = m.nextID
		m.nextID++
	}

	// Create a copy to avoid external mutations
	txnCopy := *transaction
	m.transactions[txnCopy.ID] = &txnCopy
	m.transactionIDMap[txnCopy.TransactionID] = &txnCopy

	return nil
}

// GetByID retrieves a payment transaction by ID
func (m *PaymentTransactionRepository) GetByID(id uint) (*entity.PaymentTransaction, error) {
	if m.GetByIDError != nil {
		return nil, m.GetByIDError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	transaction, exists := m.transactions[id]
	if !exists {
		return nil, fmt.Errorf("payment transaction with ID %d not found", id)
	}

	// Return a copy
	txnCopy := *transaction
	return &txnCopy, nil
}

// GetByTransactionID retrieves a payment transaction by external transaction ID
func (m *PaymentTransactionRepository) GetByTransactionID(transactionID string) (*entity.PaymentTransaction, error) {
	if m.GetByTransactionIDError != nil {
		return nil, m.GetByTransactionIDError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	transaction, exists := m.transactionIDMap[transactionID]
	if !exists {
		return nil, fmt.Errorf("payment transaction with transaction ID %s not found", transactionID)
	}

	// Return a copy
	txnCopy := *transaction
	return &txnCopy, nil
}

// GetByOrderID retrieves all payment transactions for an order
func (m *PaymentTransactionRepository) GetByOrderID(orderID uint) ([]*entity.PaymentTransaction, error) {
	if m.GetByOrderIDError != nil {
		return nil, m.GetByOrderIDError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*entity.PaymentTransaction
	for _, transaction := range m.transactions {
		if transaction.OrderID == orderID {
			txnCopy := *transaction
			results = append(results, &txnCopy)
		}
	}

	return results, nil
}

// Update updates a payment transaction
func (m *PaymentTransactionRepository) Update(transaction *entity.PaymentTransaction) error {
	if m.UpdateError != nil {
		return m.UpdateError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.transactions[transaction.ID]; !exists {
		return fmt.Errorf("payment transaction with ID %d not found", transaction.ID)
	}

	// Update both maps
	txnCopy := *transaction
	m.transactions[txnCopy.ID] = &txnCopy
	m.transactionIDMap[txnCopy.TransactionID] = &txnCopy

	return nil
}

// Delete deletes a payment transaction
func (m *PaymentTransactionRepository) Delete(id uint) error {
	if m.DeleteError != nil {
		return m.DeleteError
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	transaction, exists := m.transactions[id]
	if !exists {
		return fmt.Errorf("payment transaction with ID %d not found", id)
	}

	delete(m.transactions, id)
	delete(m.transactionIDMap, transaction.TransactionID)

	return nil
}

// GetLatestByOrderIDAndType retrieves the latest transaction of a specific type for an order
func (m *PaymentTransactionRepository) GetLatestByOrderIDAndType(orderID uint, transactionType entity.TransactionType) (*entity.PaymentTransaction, error) {
	if m.GetLatestByOrderIDAndTypeError != nil {
		return nil, m.GetLatestByOrderIDAndTypeError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var latest *entity.PaymentTransaction
	for _, transaction := range m.transactions {
		if transaction.OrderID == orderID && transaction.Type == transactionType {
			if latest == nil || transaction.ID > latest.ID {
				latest = transaction
			}
		}
	}

	if latest == nil {
		return nil, fmt.Errorf("no payment transaction of type %s found for order %d", transactionType, orderID)
	}

	// Return a copy
	txnCopy := *latest
	return &txnCopy, nil
}

// CountSuccessfulByOrderIDAndType counts successful transactions of a specific type for an order
func (m *PaymentTransactionRepository) CountSuccessfulByOrderIDAndType(orderID uint, transactionType entity.TransactionType) (int, error) {
	if m.CountSuccessfulByOrderIDAndTypeError != nil {
		return 0, m.CountSuccessfulByOrderIDAndTypeError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	count := 0
	for _, transaction := range m.transactions {
		if transaction.OrderID == orderID &&
			transaction.Type == transactionType &&
			transaction.Status == entity.TransactionStatusSuccessful {
			count++
		}
	}

	return count, nil
}

// SumAmountByOrderIDAndType sums the amount of transactions of a specific type for an order
func (m *PaymentTransactionRepository) SumAmountByOrderIDAndType(orderID uint, transactionType entity.TransactionType) (int64, error) {
	if m.SumAmountByOrderIDAndTypeError != nil {
		return 0, m.SumAmountByOrderIDAndTypeError
	}

	m.mu.RLock()
	defer m.mu.RUnlock()

	var total int64
	for _, transaction := range m.transactions {
		if transaction.OrderID == orderID &&
			transaction.Type == transactionType &&
			transaction.Status == entity.TransactionStatusSuccessful {
			total += transaction.Amount
		}
	}

	return total, nil
}

// Helper methods for testing

// SetCreateError sets an error to be returned by Create
func (m *PaymentTransactionRepository) SetCreateError(err error) {
	m.CreateError = err
}

// SetGetByIDError sets an error to be returned by GetByID
func (m *PaymentTransactionRepository) SetGetByIDError(err error) {
	m.GetByIDError = err
}

// SetGetByTransactionIDError sets an error to be returned by GetByTransactionID
func (m *PaymentTransactionRepository) SetGetByTransactionIDError(err error) {
	m.GetByTransactionIDError = err
}

// SetGetByOrderIDError sets an error to be returned by GetByOrderID
func (m *PaymentTransactionRepository) SetGetByOrderIDError(err error) {
	m.GetByOrderIDError = err
}

// SetUpdateError sets an error to be returned by Update
func (m *PaymentTransactionRepository) SetUpdateError(err error) {
	m.UpdateError = err
}

// SetDeleteError sets an error to be returned by Delete
func (m *PaymentTransactionRepository) SetDeleteError(err error) {
	m.DeleteError = err
}

// Reset clears all data and errors
func (m *PaymentTransactionRepository) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.transactions = make(map[uint]*entity.PaymentTransaction)
	m.transactionIDMap = make(map[string]*entity.PaymentTransaction)
	m.nextID = 1
	m.CreateError = nil
	m.GetByIDError = nil
	m.GetByTransactionIDError = nil
	m.GetByOrderIDError = nil
	m.UpdateError = nil
	m.DeleteError = nil
	m.GetLatestByOrderIDAndTypeError = nil
	m.CountSuccessfulByOrderIDAndTypeError = nil
	m.SumAmountByOrderIDAndTypeError = nil
}

// GetAllTransactions returns all transactions (for testing purposes)
func (m *PaymentTransactionRepository) GetAllTransactions() []*entity.PaymentTransaction {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var results []*entity.PaymentTransaction
	for _, transaction := range m.transactions {
		txnCopy := *transaction
		results = append(results, &txnCopy)
	}

	return results
}
