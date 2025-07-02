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
			"idempotency-key-12345",
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
			"idempotency-key-12345",
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
			"idempotency-key-partial-123",
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
			"idempotency-key-partial-1",
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
			"idempotency-key-partial-2",
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
			"idempotency-key-webhook-123", // Unique idempotency key for the first delivery
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
			"idempotency-key-webhook-123",      // Same idempotency key
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

// TestAmountTracking tests the new amount tracking fields (authorized_amount, captured_amount, refunded_amount)
func TestAmountTracking(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := NewTransactionRepository(db)

	// Create a test order
	order := testutil.CreateTestOrder(t, db, 1)

	t.Run("Amount tracking fields are set correctly based on transaction type", func(t *testing.T) {
		// Test 1: Authorization transaction
		authTxn, err := entity.NewPaymentTransaction(
			order.ID,
			"mp_auth_123",
			"auth-idempotency-key",
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusSuccessful,
			10000, // $100.00
			"DKK",
			"mobilepay",
		)
		require.NoError(t, err)
		require.Equal(t, int64(10000), authTxn.AuthorizedAmount, "Authorized amount should be set for authorize transaction")
		require.Equal(t, int64(0), authTxn.CapturedAmount, "Captured amount should be 0 for authorize transaction")
		require.Equal(t, int64(0), authTxn.RefundedAmount, "Refunded amount should be 0 for authorize transaction")

		err = repo.Create(authTxn)
		require.NoError(t, err)

		// Test 2: Capture transaction
		captureTxn, err := entity.NewPaymentTransaction(
			order.ID,
			"mp_capture_123",
			"capture-idempotency-key",
			entity.TransactionTypeCapture,
			entity.TransactionStatusSuccessful,
			8000, // $80.00 (partial capture)
			"DKK",
			"mobilepay",
		)
		require.NoError(t, err)
		require.Equal(t, int64(0), captureTxn.AuthorizedAmount, "Authorized amount should be 0 for capture transaction")
		require.Equal(t, int64(8000), captureTxn.CapturedAmount, "Captured amount should be set for capture transaction")
		require.Equal(t, int64(0), captureTxn.RefundedAmount, "Refunded amount should be 0 for capture transaction")

		err = repo.Create(captureTxn)
		require.NoError(t, err)

		// Test 3: Refund transaction
		refundTxn, err := entity.NewPaymentTransaction(
			order.ID,
			"mp_refund_123",
			"refund-idempotency-key",
			entity.TransactionTypeRefund,
			entity.TransactionStatusSuccessful,
			3000, // $30.00 (partial refund)
			"DKK",
			"mobilepay",
		)
		require.NoError(t, err)
		require.Equal(t, int64(0), refundTxn.AuthorizedAmount, "Authorized amount should be 0 for refund transaction")
		require.Equal(t, int64(0), refundTxn.CapturedAmount, "Captured amount should be 0 for refund transaction")
		require.Equal(t, int64(3000), refundTxn.RefundedAmount, "Refunded amount should be set for refund transaction")

		err = repo.Create(refundTxn)
		require.NoError(t, err)

		// Test 4: Verify sum methods work correctly
		totalAuthorized, err := repo.SumAuthorizedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(10000), totalAuthorized, "Total authorized amount should be 10000")

		totalCaptured, err := repo.SumCapturedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(8000), totalCaptured, "Total captured amount should be 8000")

		totalRefunded, err := repo.SumRefundedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(3000), totalRefunded, "Total refunded amount should be 3000")

		// Test 5: Verify remaining refundable amount calculation
		remainingRefundable := totalCaptured - totalRefunded
		assert.Equal(t, int64(5000), remainingRefundable, "Remaining refundable amount should be 5000")
	})

	t.Run("Multiple transactions of same type are summed correctly", func(t *testing.T) {
		// Create another order for this test
		order2 := testutil.CreateTestOrder(t, db, 2)

		// Create two capture transactions
		capture1, err := entity.NewPaymentTransaction(
			order2.ID,
			"mp_capture_1",
			"capture-1-idempotency-key",
			entity.TransactionTypeCapture,
			entity.TransactionStatusSuccessful,
			6000, // $60.00
			"DKK",
			"mobilepay",
		)
		require.NoError(t, err)
		err = repo.Create(capture1)
		require.NoError(t, err)

		capture2, err := entity.NewPaymentTransaction(
			order2.ID,
			"mp_capture_2",
			"capture-2-idempotency-key",
			entity.TransactionTypeCapture,
			entity.TransactionStatusSuccessful,
			4000, // $40.00
			"DKK",
			"mobilepay",
		)
		require.NoError(t, err)
		err = repo.Create(capture2)
		require.NoError(t, err)

		// Verify sum is correct
		totalCaptured, err := repo.SumCapturedAmountByOrderID(order2.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(10000), totalCaptured, "Total captured from two transactions should be 10000")

		// Create two refund transactions
		refund1, err := entity.NewPaymentTransaction(
			order2.ID,
			"mp_refund_1",
			"refund-1-idempotency-key",
			entity.TransactionTypeRefund,
			entity.TransactionStatusSuccessful,
			2000, // $20.00
			"DKK",
			"mobilepay",
		)
		require.NoError(t, err)
		err = repo.Create(refund1)
		require.NoError(t, err)

		refund2, err := entity.NewPaymentTransaction(
			order2.ID,
			"mp_refund_2",
			"refund-2-idempotency-key",
			entity.TransactionTypeRefund,
			entity.TransactionStatusSuccessful,
			1500, // $15.00
			"DKK",
			"mobilepay",
		)
		require.NoError(t, err)
		err = repo.Create(refund2)
		require.NoError(t, err)

		// Verify refund sum is correct
		totalRefunded, err := repo.SumRefundedAmountByOrderID(order2.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(3500), totalRefunded, "Total refunded from two transactions should be 3500")

		// Verify remaining amount
		remainingRefundable := totalCaptured - totalRefunded
		assert.Equal(t, int64(6500), remainingRefundable, "Remaining refundable should be 6500")
	})
}

// TestPendingTransactionAmountHandling tests that pending transactions don't set amount fields
// until they transition to successful status
func TestPendingTransactionAmountHandling(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := NewTransactionRepository(db)

	// Create a test order
	order := testutil.CreateTestOrder(t, db, 1)

	t.Run("Pending authorization should not set authorized amount", func(t *testing.T) {
		// Create a pending authorization
		pendingAuth, err := entity.NewPaymentTransaction(
			order.ID,
			"pi_pending_123",
			"idempotency-key-pending-123",
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusPending,
			10000, // $100.00
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		// Verify amount fields are 0 for pending transaction
		assert.Equal(t, int64(0), pendingAuth.AuthorizedAmount)
		assert.Equal(t, int64(0), pendingAuth.CapturedAmount)
		assert.Equal(t, int64(0), pendingAuth.RefundedAmount)
		assert.Equal(t, int64(10000), pendingAuth.Amount) // Original amount is still stored

		err = repo.Create(pendingAuth)
		require.NoError(t, err)

		// Verify sums are still 0 for pending transactions
		authSum, err := repo.SumAuthorizedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), authSum)

		captureSum, err := repo.SumCapturedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), captureSum)

		refundSum, err := repo.SumRefundedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), refundSum)
	})

	t.Run("Updating pending to successful should set authorized amount", func(t *testing.T) {
		// Get the pending transaction
		transactions, err := repo.GetByOrderID(order.ID)
		require.NoError(t, err)
		require.Len(t, transactions, 1)

		pendingAuth := transactions[0]
		require.Equal(t, entity.TransactionStatusPending, pendingAuth.Status)
		require.Equal(t, int64(0), pendingAuth.AuthorizedAmount)

		// Update status to successful
		pendingAuth.UpdateStatus(entity.TransactionStatusSuccessful)

		// Verify authorized amount is now set
		assert.Equal(t, int64(10000), pendingAuth.AuthorizedAmount)
		assert.Equal(t, int64(0), pendingAuth.CapturedAmount)
		assert.Equal(t, int64(0), pendingAuth.RefundedAmount)

		// Save the updated transaction
		err = repo.Update(pendingAuth)
		require.NoError(t, err)

		// Verify sums now reflect the successful authorization
		authSum, err := repo.SumAuthorizedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(10000), authSum)

		captureSum, err := repo.SumCapturedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), captureSum)
	})

	t.Run("Failed transactions should not contribute to amounts", func(t *testing.T) {
		// Create a failed capture attempt
		failedCapture, err := entity.NewPaymentTransaction(
			order.ID,
			"ch_failed_123",
			"idempotency-key-failed-123",
			entity.TransactionTypeCapture,
			entity.TransactionStatusFailed,
			10000, // $100.00
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		// Verify amount fields are 0 for failed transaction
		assert.Equal(t, int64(0), failedCapture.AuthorizedAmount)
		assert.Equal(t, int64(0), failedCapture.CapturedAmount)
		assert.Equal(t, int64(0), failedCapture.RefundedAmount)

		err = repo.Create(failedCapture)
		require.NoError(t, err)

		// Verify sums are not affected by failed transactions
		authSum, err := repo.SumAuthorizedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(10000), authSum) // Still the same from successful auth

		captureSum, err := repo.SumCapturedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), captureSum) // Failed capture doesn't count
	})

	t.Run("Transitioning successful to failed should clear amount field", func(t *testing.T) {
		// Create a successful transaction first
		successfulCapture, err := entity.NewPaymentTransaction(
			order.ID,
			"ch_success_then_fail",
			"idempotency-key-success-fail",
			entity.TransactionTypeCapture,
			entity.TransactionStatusSuccessful,
			5000, // $50.00
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		// Verify captured amount is set
		assert.Equal(t, int64(5000), successfulCapture.CapturedAmount)

		err = repo.Create(successfulCapture)
		require.NoError(t, err)

		// Verify capture sum includes this transaction
		captureSum, err := repo.SumCapturedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(5000), captureSum)

		// Now update to failed status
		successfulCapture.UpdateStatus(entity.TransactionStatusFailed)

		// Verify captured amount is cleared
		assert.Equal(t, int64(0), successfulCapture.CapturedAmount)

		// Save the updated transaction
		err = repo.Update(successfulCapture)
		require.NoError(t, err)

		// Verify capture sum no longer includes this transaction
		captureSum, err = repo.SumCapturedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), captureSum)
	})
}

// TestTransactionStatusUpdateWithAmountFields tests that updating transaction status
// properly sets the amount fields when transitioning from pending to successful
func TestTransactionStatusUpdateWithAmountFields(t *testing.T) {
	db := testutil.SetupTestDB(t)
	repo := NewTransactionRepository(db)

	// Create a test order
	order := testutil.CreateTestOrder(t, db, 1)

	t.Run("Pending authorization updated to successful should set authorized amount", func(t *testing.T) {
		// Create a pending authorization transaction
		pendingAuth, err := entity.NewPaymentTransaction(
			order.ID,
			"pi_pending_update_test",
			"idempotency-key-update-test",
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusPending,
			15000, // $150.00
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		// Verify it starts with no amount fields set
		assert.Equal(t, int64(0), pendingAuth.AuthorizedAmount)
		assert.Equal(t, int64(0), pendingAuth.CapturedAmount)
		assert.Equal(t, int64(0), pendingAuth.RefundedAmount)
		assert.Equal(t, int64(15000), pendingAuth.Amount)

		// Save the pending transaction
		err = repo.Create(pendingAuth)
		require.NoError(t, err)

		// Verify sums are 0 for pending
		authSum, err := repo.SumAuthorizedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), authSum)

		// Update status to successful (this should set the authorized amount)
		pendingAuth.UpdateStatus(entity.TransactionStatusSuccessful)

		// Verify authorized amount is now set
		assert.Equal(t, int64(15000), pendingAuth.AuthorizedAmount)
		assert.Equal(t, int64(0), pendingAuth.CapturedAmount)
		assert.Equal(t, int64(0), pendingAuth.RefundedAmount)

		// Save the updated transaction
		err = repo.Update(pendingAuth)
		require.NoError(t, err)

		// Verify sums now reflect the successful authorization
		authSum, err = repo.SumAuthorizedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(15000), authSum)
	})

	t.Run("Pending capture updated to successful should set captured amount", func(t *testing.T) {
		// Create a pending capture transaction
		pendingCapture, err := entity.NewPaymentTransaction(
			order.ID,
			"ch_pending_capture_test",
			"idempotency-key-capture-test",
			entity.TransactionTypeCapture,
			entity.TransactionStatusPending,
			15000, // $150.00
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		// Save the pending transaction
		err = repo.Create(pendingCapture)
		require.NoError(t, err)

		// Verify capture sum is still 0 for pending
		captureSum, err := repo.SumCapturedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), captureSum)

		// Update status to successful
		pendingCapture.UpdateStatus(entity.TransactionStatusSuccessful)

		// Verify captured amount is now set (but not authorized amount)
		assert.Equal(t, int64(0), pendingCapture.AuthorizedAmount) // Should remain 0 for capture transaction
		assert.Equal(t, int64(15000), pendingCapture.CapturedAmount)
		assert.Equal(t, int64(0), pendingCapture.RefundedAmount)

		// Save the updated transaction
		err = repo.Update(pendingCapture)
		require.NoError(t, err)

		// Verify capture sum now includes this transaction
		captureSum, err = repo.SumCapturedAmountByOrderID(order.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(15000), captureSum)
	})

	t.Run("Multiple pending transactions updated individually", func(t *testing.T) {
		// Create a new order for this test
		order2 := testutil.CreateTestOrder(t, db, 2)

		// Create multiple pending transactions
		pending1, err := entity.NewPaymentTransaction(
			order2.ID,
			"txn_multi_1",
			"idempotency-key-multi-1",
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusPending,
			8000, // $80.00
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		pending2, err := entity.NewPaymentTransaction(
			order2.ID,
			"txn_multi_2",
			"idempotency-key-multi-2",
			entity.TransactionTypeAuthorize,
			entity.TransactionStatusPending,
			2000, // $20.00
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		// Save both pending transactions
		err = repo.Create(pending1)
		require.NoError(t, err)
		err = repo.Create(pending2)
		require.NoError(t, err)

		// Verify no authorized amounts yet
		authSum, err := repo.SumAuthorizedAmountByOrderID(order2.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(0), authSum)

		// Update first transaction to successful
		pending1.UpdateStatus(entity.TransactionStatusSuccessful)
		err = repo.Update(pending1)
		require.NoError(t, err)

		// Verify only first transaction contributes to sum
		authSum, err = repo.SumAuthorizedAmountByOrderID(order2.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(8000), authSum)

		// Update second transaction to successful
		pending2.UpdateStatus(entity.TransactionStatusSuccessful)
		err = repo.Update(pending2)
		require.NoError(t, err)

		// Verify both transactions contribute to sum
		authSum, err = repo.SumAuthorizedAmountByOrderID(order2.ID)
		require.NoError(t, err)
		assert.Equal(t, int64(10000), authSum) // $80 + $20 = $100
	})
}
