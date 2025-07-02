package entity

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPaymentTransaction(t *testing.T) {
	t.Run("NewPaymentTransaction success", func(t *testing.T) {
		txn, err := NewPaymentTransaction(
			1,
			"txn_123",
			TransactionTypeAuthorize,
			TransactionStatusSuccessful,
			10000,
			"USD",
			"stripe",
		)

		fmt.Printf("orderID: %d, transactionID: %s, type: %s, status: %s, amount: %d, currency: %s, provider: %s\n",
			txn.OrderID, txn.ExternalID, txn.Type, txn.Status, txn.Amount,
			txn.Currency, txn.Provider)

		require.NoError(t, err)
		assert.Equal(t, uint(1), txn.OrderID)
		assert.Equal(t, "txn_123", txn.ExternalID)
		assert.Equal(t, TransactionTypeAuthorize, txn.Type)
		assert.Equal(t, TransactionStatusSuccessful, txn.Status)
		assert.Equal(t, int64(10000), txn.Amount)
		assert.Equal(t, "USD", txn.Currency)
		assert.Equal(t, "stripe", txn.Provider)
		assert.NotNil(t, txn.Metadata)
		assert.Empty(t, txn.Metadata)
	})

	t.Run("NewPaymentTransaction validation errors", func(t *testing.T) {
		tests := []struct {
			name          string
			orderID       uint
			externalID    string
			txnType       TransactionType
			status        TransactionStatus
			amount        int64
			currency      string
			provider      string
			expectedError string
		}{
			{
				name:          "zero orderID",
				orderID:       0,
				externalID:    "txn_123",
				txnType:       TransactionTypeAuthorize,
				status:        TransactionStatusSuccessful,
				amount:        10000,
				currency:      "USD",
				provider:      "stripe",
				expectedError: "orderID cannot be zero",
			},
			{
				name:          "empty transactionType",
				orderID:       1,
				externalID:    "txn_123",
				txnType:       "",
				status:        TransactionStatusSuccessful,
				amount:        10000,
				currency:      "USD",
				provider:      "stripe",
				expectedError: "transactionType cannot be empty",
			},
			{
				name:          "empty status",
				orderID:       1,
				externalID:    "txn_123",
				txnType:       TransactionTypeAuthorize,
				status:        "",
				amount:        10000,
				currency:      "USD",
				provider:      "stripe",
				expectedError: "status cannot be empty",
			},
			{
				name:          "empty currency",
				orderID:       1,
				externalID:    "txn_123",
				txnType:       TransactionTypeAuthorize,
				status:        TransactionStatusSuccessful,
				amount:        10000,
				currency:      "",
				provider:      "stripe",
				expectedError: "currency cannot be empty",
			},
			{
				name:          "empty provider",
				orderID:       1,
				externalID:    "txn_123",
				txnType:       TransactionTypeAuthorize,
				status:        TransactionStatusSuccessful,
				amount:        10000,
				currency:      "USD",
				provider:      "",
				expectedError: "provider cannot be empty",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				txn, err := NewPaymentTransaction(
					tt.orderID,
					tt.externalID,
					tt.txnType,
					tt.status,
					tt.amount,
					tt.currency,
					tt.provider,
				)

				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, txn)
			})
		}
	})

	t.Run("AddMetadata", func(t *testing.T) {
		txn, err := NewPaymentTransaction(
			1,
			"txn_123",
			TransactionTypeAuthorize,
			TransactionStatusSuccessful,
			10000,
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		// Add metadata
		txn.AddMetadata("key1", "value1")
		txn.AddMetadata("key2", "value2")

		assert.Equal(t, "value1", txn.Metadata["key1"])
		assert.Equal(t, "value2", txn.Metadata["key2"])
		assert.Len(t, txn.Metadata, 2)
	})

	t.Run("AddMetadata with nil map", func(t *testing.T) {
		txn := &PaymentTransaction{}
		txn.AddMetadata("key1", "value1")

		assert.Equal(t, "value1", txn.Metadata["key1"])
		assert.Len(t, txn.Metadata, 1)
	})

	t.Run("SetRawResponse", func(t *testing.T) {
		txn, err := NewPaymentTransaction(
			1,
			"txn_123",
			TransactionTypeAuthorize,
			TransactionStatusSuccessful,
			10000,
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		response := `{"id": "ch_123", "status": "succeeded"}`
		txn.SetRawResponse(response)

		assert.Equal(t, response, txn.RawResponse)
	})

	t.Run("UpdateStatus", func(t *testing.T) {
		txn, err := NewPaymentTransaction(
			1,
			"txn_123",
			TransactionTypeAuthorize,
			TransactionStatusPending,
			10000,
			"USD",
			"stripe",
		)
		require.NoError(t, err)

		txn.UpdateStatus(TransactionStatusSuccessful)
		assert.Equal(t, TransactionStatusSuccessful, txn.Status)
	})
}

func TestTransactionTypeConstants(t *testing.T) {
	assert.Equal(t, TransactionType("authorize"), TransactionTypeAuthorize)
	assert.Equal(t, TransactionType("capture"), TransactionTypeCapture)
	assert.Equal(t, TransactionType("refund"), TransactionTypeRefund)
	assert.Equal(t, TransactionType("cancel"), TransactionTypeCancel)
}

func TestTransactionStatusConstants(t *testing.T) {
	assert.Equal(t, TransactionStatus("successful"), TransactionStatusSuccessful)
	assert.Equal(t, TransactionStatus("failed"), TransactionStatusFailed)
	assert.Equal(t, TransactionStatus("pending"), TransactionStatusPending)
}
