package gorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/testutil"
)

// TestPaymentCaptureWorkflow tests the complete payment workflow
// to ensure that each action (authorize, capture) creates a separate transaction record
func TestPaymentCaptureWorkflow(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := NewTransactionRepository(db)

	// Create a test order
	order := testutil.CreateTestOrder(t, db, 1)

	t.Run("Complete payment workflow with separate transactions", func(t *testing.T) {
		// Step 1: Create authorization transaction (initial payment processing)
		authTxn, err := entity.NewPaymentTransaction(
			order.ID,
			"pi_1234567890", // Payment Intent ID from Stripe
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusSuccessful,
			10000, // $100.00
			"USD",
			"stripe",
		)
		require.NoError(t, err)
		authTxn.RawResponse = `{"id": "pi_1234567890", "status": "requires_capture"}`
		authTxn.AddMetadata("payment_intent_id", "pi_1234567890")

		err = repo.Create(authTxn)
		require.NoError(t, err)

		// Step 2: Create capture transaction (when payment is captured via webhook)
		// In real scenarios, this might have the same transaction ID or a different one
		captureTxn, err := entity.NewPaymentTransaction(
			order.ID,
			"ch_1234567890", // Charge ID from Stripe (different from payment intent)
			entity.TransactionTypeCapture,
			entity.TransactionStatusSuccessful,
			10000, // Same amount
			"USD",
			"stripe",
		)
		require.NoError(t, err)
		captureTxn.RawResponse = `{"id": "ch_1234567890", "status": "succeeded", "captured": true}`
		captureTxn.AddMetadata("charge_id", "ch_1234567890")
		captureTxn.AddMetadata("webhook_id", "we_1234567890")

		err = repo.Create(captureTxn)
		require.NoError(t, err)

		// Verify both transactions were created as separate records
		assert.NotEqual(t, authTxn.ID, captureTxn.ID)

		// Verify we can retrieve all transactions for the order
		transactions, err := repo.GetByOrderID(order.ID)
		require.NoError(t, err)
		assert.Len(t, transactions, 2)

		// Verify we can get the latest transaction of each type
		latestAuth, err := repo.GetLatestByOrderIDAndType(order.ID, entity.TransactionTypeAuthorize)
		require.NoError(t, err)
		assert.Equal(t, "pi_1234567890", latestAuth.ExternalID) // Check ExternalID instead of TransactionID
		assert.Equal(t, "pi_1234567890", latestAuth.Metadata["payment_intent_id"])

		latestCapture, err := repo.GetLatestByOrderIDAndType(order.ID, entity.TransactionTypeCapture)
		require.NoError(t, err)
		assert.Equal(t, "ch_1234567890", latestCapture.ExternalID) // Check ExternalID instead of TransactionID
		assert.Equal(t, "ch_1234567890", latestCapture.Metadata["charge_id"])
		assert.Equal(t, "we_1234567890", latestCapture.Metadata["webhook_id"])

		// Verify transaction counts and sums
		authCount, err := repo.CountSuccessfulByOrderIDAndType(order.ID, entity.TransactionTypeAuthorize)
		require.NoError(t, err)
		assert.Equal(t, 1, authCount)

		captureCount, err := repo.CountSuccessfulByOrderIDAndType(order.ID, entity.TransactionTypeCapture)
		require.NoError(t, err)
		assert.Equal(t, 1, captureCount)

		totalCaptured, err := repo.SumAmountByOrderIDAndType(order.ID, entity.TransactionTypeCapture)
		require.NoError(t, err)
		assert.Equal(t, int64(10000), totalCaptured)
	})

	t.Run("Partial capture workflow", func(t *testing.T) {
		// Create another order for partial capture testing
		order2 := testutil.CreateTestOrder(t, db, 2)

		// Step 1: Authorization for $100
		authTxn, err := entity.NewPaymentTransaction(
			order2.ID,
			"pi_partial_123",
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusSuccessful,
			10000, // $100.00
			"USD",
			"stripe",
		)
		require.NoError(t, err)
		err = repo.Create(authTxn)
		require.NoError(t, err)

		// Step 2: First partial capture for $60
		capture1, err := entity.NewPaymentTransaction(
			order2.ID,
			"ch_partial_1",
			entity.TransactionTypeCapture,
			entity.TransactionStatusSuccessful,
			6000, // $60.00
			"USD",
			"stripe",
		)
		require.NoError(t, err)
		capture1.AddMetadata("partial_capture", "1")
		err = repo.Create(capture1)
		require.NoError(t, err)

		// Step 3: Second partial capture for $40
		capture2, err := entity.NewPaymentTransaction(
			order2.ID,
			"ch_partial_2",
			entity.TransactionTypeCapture,
			entity.TransactionStatusSuccessful,
			4000, // $40.00
			"USD",
			"stripe",
		)
		require.NoError(t, err)
		capture2.AddMetadata("partial_capture", "2")
		err = repo.Create(capture2)
		require.NoError(t, err)

		// Verify all transactions are separate
		transactions, err := repo.GetByOrderID(order2.ID)
		require.NoError(t, err)
		assert.Len(t, transactions, 3) // 1 auth + 2 captures

		// Verify capture count and total
		captureCount, err := repo.CountSuccessfulByOrderIDAndType(order2.ID, entity.TransactionTypeCapture)
		require.NoError(t, err)
		assert.Equal(t, 2, captureCount)

		totalCaptured, err := repo.SumAmountByOrderIDAndType(order2.ID, entity.TransactionTypeCapture)
		require.NoError(t, err)
		assert.Equal(t, int64(10000), totalCaptured) // $60 + $40 = $100

		// Verify each capture transaction has correct metadata
		captures, err := repo.GetByOrderID(order2.ID)
		require.NoError(t, err)

		var partialCapture1, partialCapture2 *entity.PaymentTransaction
		for _, txn := range captures {
			if txn.Type == entity.TransactionTypeCapture {
				switch txn.Metadata["partial_capture"] {
				case "1":
					partialCapture1 = txn
				case "2":
					partialCapture2 = txn
				}
			}
		}

		require.NotNil(t, partialCapture1)
		require.NotNil(t, partialCapture2)
		assert.Equal(t, int64(6000), partialCapture1.Amount)
		assert.Equal(t, int64(4000), partialCapture2.Amount)
	})
}

// TestWebhookDuplicationHandling tests scenarios where the same webhook might be received multiple times
func TestWebhookDuplicationHandling(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := NewTransactionRepository(db)

	// Create a test order
	order := testutil.CreateTestOrder(t, db, 1)

	t.Run("Same webhook received multiple times creates multiple records", func(t *testing.T) {
		// This demonstrates the current behavior - each create call will create a new record
		// In practice, you might want to implement webhook idempotency at the application level
		// using webhook IDs or other identifiers

		webhookID := "we_duplicate_test"

		// First webhook delivery
		txn1, err := entity.NewPaymentTransaction(
			order.ID,
			"ch_webhook_test",
			entity.TransactionTypeCapture,
			entity.TransactionStatusSuccessful,
			5000,
			"USD",
			"stripe",
		)
		require.NoError(t, err)
		txn1.AddMetadata("webhook_id", webhookID)
		err = repo.Create(txn1)
		require.NoError(t, err)

		// Second webhook delivery (duplicate)
		txn2, err := entity.NewPaymentTransaction(
			order.ID,
			"ch_webhook_test",                  // Same transaction ID
			entity.TransactionTypeCapture,      // Same type
			entity.TransactionStatusSuccessful, // Same status
			5000,                               // Same amount
			"USD",
			"stripe",
		)
		require.NoError(t, err)
		txn2.AddMetadata("webhook_id", webhookID) // Same webhook ID
		err = repo.Create(txn2)
		require.NoError(t, err)

		// Both transactions are created as separate records
		assert.NotEqual(t, txn1.ID, txn2.ID)

		// Count shows 2 transactions
		count, err := repo.CountSuccessfulByOrderIDAndType(order.ID, entity.TransactionTypeCapture)
		require.NoError(t, err)
		assert.Equal(t, 2, count)

		// Note: In a real application, you would typically implement webhook idempotency
		// at the application layer by checking for existing transactions with the same
		// webhook_id before creating new ones.
	})
}
