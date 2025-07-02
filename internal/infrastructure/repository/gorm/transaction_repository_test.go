package gorm

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// setupTestDB creates an in-memory SQLite database for testing
func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	require.NoError(t, err)

	// Auto-migrate the payment_transactions table
	err = db.AutoMigrate(&entity.PaymentTransaction{}, &entity.Order{})
	require.NoError(t, err)

	return db
}

// createTestOrder creates a test order in the database
func createTestOrder(t *testing.T, db *gorm.DB, orderID uint) *entity.Order {
	order := &entity.Order{
		Model:         gorm.Model{ID: orderID},
		OrderNumber:   fmt.Sprintf("ORD-%d", orderID), // Make order number unique
		TotalAmount:   10000,
		Currency:      "USD",
		Status:        entity.OrderStatusPending,
		PaymentStatus: entity.PaymentStatusPending,
		IsGuestOrder:  true,
	}
	err := db.Create(order).Error
	require.NoError(t, err)
	return order
}

func TestTransactionRepository_Create(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)

	// Create a test order
	createTestOrder(t, db, 1)

	t.Run("Create new transaction successfully", func(t *testing.T) {
		txn, err := entity.NewPaymentTransaction(
			1,
			"txn_123",
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusSuccessful,
			10000,
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		err = repo.Create(txn)
		assert.NoError(t, err)
		assert.NotZero(t, txn.ID)
	})

	t.Run("Create identical transaction should update existing record", func(t *testing.T) {
		// Create first transaction
		txn1, err := entity.NewPaymentTransaction(
			1,
			"external_id_duplicate",
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusPending,
			5000,
			"USD",
			"stripe",
		)
		require.NoError(t, err)
		txn1.RawResponse = "original response"

		err = repo.CreateOrUpdate(txn1) // Use CreateOrUpdate for upsert behavior
		require.NoError(t, err)
		originalID := txn1.ID
		originalTransactionID := txn1.TransactionID

		// Create "identical" transaction (same order + type) with different status and external ID
		txn2, err := entity.NewPaymentTransaction(
			1,
			"external_id_updated",              // Different external ID
			entity.TransactionTypeAuthorize,    // Same type (this will trigger update)
			entity.TransactionStatusSuccessful, // Different status
			5000,                               // Same amount
			"USD",
			"stripe",
		)
		require.NoError(t, err)
		txn2.RawResponse = "updated response"
		txn2.AddMetadata("webhook_id", "wh_123")

		err = repo.CreateOrUpdate(txn2) // Use CreateOrUpdate for upsert behavior
		assert.NoError(t, err)

		// Verify that the existing transaction was updated, not a new one created
		assert.Equal(t, originalID, txn2.ID)
		assert.Equal(t, originalTransactionID, txn2.TransactionID)
		assert.Equal(t, entity.TransactionStatusSuccessful, txn2.Status)
		assert.Equal(t, "external_id_updated", txn2.ExternalID)

		// Verify only one transaction exists for this order + type
		var count int64
		err = db.Model(&entity.PaymentTransaction{}).Where("order_id = ? AND type = ?", 1, entity.TransactionTypeAuthorize).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Create transaction with different amount should update existing record", func(t *testing.T) {
		// Create first transaction
		txn1, err := entity.NewPaymentTransaction(
			1,
			"external_id_amount_test",
			entity.TransactionTypeCapture,
			entity.TransactionStatusSuccessful,
			5000,
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		err = repo.CreateOrUpdate(txn1) // Use CreateOrUpdate for upsert behavior
		require.NoError(t, err)
		originalID := txn1.ID
		originalTransactionID := txn1.TransactionID

		// Create transaction with same order + type but different amount (should update)
		txn2, err := entity.NewPaymentTransaction(
			1,
			"external_id_amount_updated",
			entity.TransactionTypeCapture, // Same type, so will update
			entity.TransactionStatusSuccessful,
			3000, // Different amount
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		err = repo.CreateOrUpdate(txn2) // Use CreateOrUpdate for upsert behavior
		assert.NoError(t, err)

		// Verify the existing transaction was updated
		assert.Equal(t, originalID, txn2.ID)
		assert.Equal(t, originalTransactionID, txn2.TransactionID)
		assert.Equal(t, int64(3000), txn2.Amount)

		// Verify only one transaction exists for this order + type
		var count int64
		err = db.Model(&entity.PaymentTransaction{}).Where("order_id = ? AND type = ?", 1, entity.TransactionTypeCapture).Count(&count).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), count)
	})

	t.Run("Create multiple transactions with different types should create separate records", func(t *testing.T) {
		// Create a new test order specifically for this test to avoid conflicts with previous tests
		createTestOrder(t, db, 99)

		// Create authorization transaction
		txn1, err := entity.NewPaymentTransaction(
			99, // Use order 99 to avoid conflicts
			"external_id_auth",
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusSuccessful,
			10000,
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		err = repo.Create(txn1)
		require.NoError(t, err)

		// Create capture transaction (different type, so should create new record)
		txn2, err := entity.NewPaymentTransaction(
			99, // Use order 99 to avoid conflicts
			"external_id_capture",
			entity.TransactionTypeCapture, // Different type
			entity.TransactionStatusSuccessful,
			10000,
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		err = repo.Create(txn2)
		assert.NoError(t, err)

		// Verify both transactions exist (different types)
		var authCount int64
		err = db.Model(&entity.PaymentTransaction{}).Where("order_id = ? AND type = ?", 99, entity.TransactionTypeAuthorize).Count(&authCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), authCount)

		var captureCount int64
		err = db.Model(&entity.PaymentTransaction{}).Where("order_id = ? AND type = ?", 99, entity.TransactionTypeCapture).Count(&captureCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), captureCount)

		// Verify they have different IDs and transaction IDs
		assert.NotEqual(t, txn1.ID, txn2.ID)
		assert.NotEqual(t, txn1.TransactionID, txn2.TransactionID)
	})
}

func TestTransactionRepository_GetByID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)

	// Create a test order
	createTestOrder(t, db, 1)

	t.Run("Get existing transaction", func(t *testing.T) {
		txn, err := entity.NewPaymentTransaction(
			1,
			"txn_get_by_id",
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusSuccessful,
			10000,
			"USD",
			"stripe",
		)
		require.NoError(t, err)
		txn.AddMetadata("test_key", "test_value")

		err = repo.Create(txn)
		require.NoError(t, err)

		retrieved, err := repo.GetByID(txn.ID)
		assert.NoError(t, err)
		assert.Equal(t, txn.TransactionID, retrieved.TransactionID)
		assert.Equal(t, txn.Type, retrieved.Type)
		assert.Equal(t, txn.Status, retrieved.Status)
		assert.Equal(t, txn.Amount, retrieved.Amount)
		assert.Equal(t, "test_value", retrieved.Metadata["test_key"])
	})

	t.Run("Get non-existent transaction", func(t *testing.T) {
		_, err := repo.GetByID(99999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestTransactionRepository_GetByTransactionID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)

	// Create a test order
	createTestOrder(t, db, 1)

	t.Run("Get existing transaction by transaction ID", func(t *testing.T) {
		txn, err := entity.NewPaymentTransaction(
			1,
			"external_id_123", // This is the external ID
			entity.TransactionTypeCapture,
			entity.TransactionStatusSuccessful,
			5000,
			"EUR",
			"paypal",
		)
		require.NoError(t, err)

		err = repo.Create(txn)
		require.NoError(t, err)

		// Use the generated friendly transaction ID (like TXN-CAPT-2025-001)
		retrieved, err := repo.GetByTransactionID(txn.TransactionID)
		assert.NoError(t, err)
		assert.Equal(t, txn.OrderID, retrieved.OrderID)
		assert.Equal(t, txn.Type, retrieved.Type)
		assert.Equal(t, "EUR", retrieved.Currency)
		assert.Equal(t, "paypal", retrieved.Provider)
		assert.Equal(t, "external_id_123", retrieved.ExternalID)
	})

	t.Run("Get non-existent transaction by transaction ID", func(t *testing.T) {
		_, err := repo.GetByTransactionID("non_existent_txn")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestTransactionRepository_GetByOrderID(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)

	// Create test orders
	createTestOrder(t, db, 1)
	createTestOrder(t, db, 2)

	t.Run("Get transactions for order with multiple transactions", func(t *testing.T) {
		// Create multiple transactions for order 1
		txn1, err := entity.NewPaymentTransaction(1, "txn_order_1_auth", entity.TransactionTypeAuthorize, entity.TransactionStatusSuccessful, 10000, "USD", "stripe")
		require.NoError(t, err)
		err = repo.Create(txn1)
		require.NoError(t, err)

		txn2, err := entity.NewPaymentTransaction(1, "txn_order_1_capture", entity.TransactionTypeCapture, entity.TransactionStatusSuccessful, 10000, "USD", "stripe")
		require.NoError(t, err)
		err = repo.Create(txn2)
		require.NoError(t, err)

		// Create transaction for order 2
		txn3, err := entity.NewPaymentTransaction(2, "txn_order_2_auth", entity.TransactionTypeAuthorize, entity.TransactionStatusSuccessful, 5000, "USD", "stripe")
		require.NoError(t, err)
		err = repo.Create(txn3)
		require.NoError(t, err)

		// Get transactions for order 1
		transactions, err := repo.GetByOrderID(1)
		assert.NoError(t, err)
		assert.Len(t, transactions, 2)

		// Verify transactions are ordered by created_at DESC
		assert.True(t, transactions[0].CreatedAt.After(transactions[1].CreatedAt) || transactions[0].CreatedAt.Equal(transactions[1].CreatedAt))
	})

	t.Run("Get transactions for order with no transactions", func(t *testing.T) {
		createTestOrder(t, db, 3)
		transactions, err := repo.GetByOrderID(3)
		assert.NoError(t, err)
		assert.Empty(t, transactions)
	})
}

func TestTransactionRepository_GetLatestByOrderIDAndType(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)

	// Create a test order
	createTestOrder(t, db, 1)

	t.Run("Get latest transaction of specific type", func(t *testing.T) {
		// Create authorization transaction
		txn1, err := entity.NewPaymentTransaction(1, "external_auth_1", entity.TransactionTypeAuthorize, entity.TransactionStatusSuccessful, 10000, "USD", "stripe")
		require.NoError(t, err)
		err = repo.Create(txn1)
		require.NoError(t, err)

		// Create a capture transaction (different type)
		txn2, err := entity.NewPaymentTransaction(1, "external_capture_1", entity.TransactionTypeCapture, entity.TransactionStatusSuccessful, 5000, "USD", "stripe")
		require.NoError(t, err)
		err = repo.Create(txn2)
		require.NoError(t, err)

		// Get latest authorization transaction (should be the only one)
		latest, err := repo.GetLatestByOrderIDAndType(1, entity.TransactionTypeAuthorize)
		assert.NoError(t, err)
		assert.Equal(t, txn1.TransactionID, latest.TransactionID)
		assert.Equal(t, txn1.ID, latest.ID)
	})

	t.Run("Get latest transaction when none exist of that type", func(t *testing.T) {
		createTestOrder(t, db, 4)
		_, err := repo.GetLatestByOrderIDAndType(4, entity.TransactionTypeRefund)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "no payment transaction of type")
	})
}

func TestTransactionRepository_CountSuccessfulByOrderIDAndType(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)

	// Create a test order
	createTestOrder(t, db, 1)

	t.Run("Count successful transactions", func(t *testing.T) {
		// Create test orders
		createTestOrder(t, db, 10)
		createTestOrder(t, db, 11)

		// Create successful capture transactions for different orders
		txn1, err := entity.NewPaymentTransaction(10, "external_success_1", entity.TransactionTypeCapture, entity.TransactionStatusSuccessful, 5000, "USD", "stripe")
		require.NoError(t, err)
		err = repo.Create(txn1)
		require.NoError(t, err)

		txn2, err := entity.NewPaymentTransaction(11, "external_success_2", entity.TransactionTypeCapture, entity.TransactionStatusSuccessful, 3000, "USD", "stripe")
		require.NoError(t, err)
		err = repo.Create(txn2)
		require.NoError(t, err)

		// Create failed capture transaction for order 1
		txn3, err := entity.NewPaymentTransaction(1, "external_failed_1", entity.TransactionTypeCapture, entity.TransactionStatusFailed, 2000, "USD", "stripe")
		require.NoError(t, err)
		err = repo.Create(txn3)
		require.NoError(t, err)

		// Count successful capture transactions (should find 2 across all orders)
		count, err := repo.CountSuccessfulByOrderIDAndType(10, entity.TransactionTypeCapture)
		assert.NoError(t, err)
		assert.Equal(t, 1, count) // Only the one for order 10

		count, err = repo.CountSuccessfulByOrderIDAndType(11, entity.TransactionTypeCapture)
		assert.NoError(t, err)
		assert.Equal(t, 1, count) // Only the one for order 11
	})

	t.Run("Count when no successful transactions exist", func(t *testing.T) {
		createTestOrder(t, db, 5)
		count, err := repo.CountSuccessfulByOrderIDAndType(5, entity.TransactionTypeRefund)
		assert.NoError(t, err)
		assert.Equal(t, 0, count)
	})
}

func TestTransactionRepository_SumAmountByOrderIDAndType(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)

	// Create a test order
	createTestOrder(t, db, 1)

	t.Run("Sum amounts for successful transactions", func(t *testing.T) {
		// Create a successful capture transaction for order 1
		txn1, err := entity.NewPaymentTransaction(1, "external_sum_1", entity.TransactionTypeCapture, entity.TransactionStatusSuccessful, 5000, "USD", "stripe")
		require.NoError(t, err)
		err = repo.Create(txn1)
		require.NoError(t, err)

		// Test the sum of that one transaction
		total, err := repo.SumAmountByOrderIDAndType(1, entity.TransactionTypeCapture)
		assert.NoError(t, err)
		assert.Equal(t, int64(5000), total)

		// Now update the transaction with a failed status using CreateOrUpdate - should not be included in sum
		txn_update, err := entity.NewPaymentTransaction(1, "external_sum_updated", entity.TransactionTypeCapture, entity.TransactionStatusFailed, 3000, "USD", "stripe")
		require.NoError(t, err)
		err = repo.CreateOrUpdate(txn_update) // Use CreateOrUpdate to update the existing capture transaction
		require.NoError(t, err)

		// Sum should now be 0 since the transaction is failed
		total, err = repo.SumAmountByOrderIDAndType(1, entity.TransactionTypeCapture)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
	})

	t.Run("Sum when no successful transactions exist", func(t *testing.T) {
		createTestOrder(t, db, 6)
		total, err := repo.SumAmountByOrderIDAndType(6, entity.TransactionTypeRefund)
		assert.NoError(t, err)
		assert.Equal(t, int64(0), total)
	})
}

func TestTransactionRepository_Update(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)

	// Create a test order
	createTestOrder(t, db, 1)

	t.Run("Update transaction successfully", func(t *testing.T) {
		txn, err := entity.NewPaymentTransaction(
			1,
			"txn_update",
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusPending,
			10000,
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		err = repo.Create(txn)
		require.NoError(t, err)

		// Update the transaction
		txn.UpdateStatus(entity.TransactionStatusSuccessful)
		txn.RawResponse = "updated response"

		err = repo.Update(txn)
		assert.NoError(t, err)

		// Verify the update
		retrieved, err := repo.GetByID(txn.ID)
		require.NoError(t, err)
		assert.Equal(t, entity.TransactionStatusSuccessful, retrieved.Status)
		assert.Equal(t, "updated response", retrieved.RawResponse)
	})
}

func TestTransactionRepository_Delete(t *testing.T) {
	db := setupTestDB(t)
	repo := NewTransactionRepository(db)

	// Create a test order
	createTestOrder(t, db, 1)

	t.Run("Delete transaction successfully", func(t *testing.T) {
		txn, err := entity.NewPaymentTransaction(
			1,
			"txn_delete",
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusSuccessful,
			10000,
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		err = repo.Create(txn)
		require.NoError(t, err)
		txnID := txn.ID

		// Delete the transaction
		err = repo.Delete(txnID)
		assert.NoError(t, err)

		// Verify it's deleted
		_, err = repo.GetByID(txnID)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("Delete non-existent transaction", func(t *testing.T) {
		err := repo.Delete(99999)
		// GORM doesn't return an error when deleting non-existent records
		assert.NoError(t, err)
	})
}
