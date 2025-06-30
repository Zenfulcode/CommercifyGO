package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOrder(t *testing.T) {
	t.Run("NewOrder success", func(t *testing.T) {
		// Create test items
		items := []OrderItem{
			{
				ProductID:   1,
				ProductName: "Test Product 1",
				SKU:         "SKU-001",
				Quantity:    2,
				Price:       9999, // $99.99
				Weight:      1.5,
			},
			{
				ProductID:   2,
				ProductName: "Test Product 2",
				SKU:         "SKU-002",
				Quantity:    1,
				Price:       4999, // $49.99
				Weight:      0.8,
			},
		}

		shippingAddr := Address{
			Street:     "123 Main St",
			City:       "Anytown",
			State:      "CA",
			PostalCode: "12345",
			Country:    "US",
		}

		billingAddr := Address{
			Street:     "456 Oak Ave",
			City:       "Another City",
			State:      "NY",
			PostalCode: "67890",
			Country:    "US",
		}

		customerDetails := CustomerDetails{
			Email:    "test@example.com",
			Phone:    "555-1234",
			FullName: "John Doe",
		}

		order, err := NewOrder(1, items, "USD", shippingAddr, billingAddr, customerDetails)

		require.NoError(t, err)
		assert.Contains(t, order.OrderNumber, "ORD-")
		assert.Equal(t, "USD", order.Currency)
		assert.Equal(t, uint(1), order.UserID)
		assert.Equal(t, OrderStatusPending, order.Status)
		assert.Equal(t, PaymentStatusPending, order.PaymentStatus)
		assert.Equal(t, int64(24997), order.TotalAmount) // (2*9999) + (1*4999)
		assert.Equal(t, int64(24997), order.FinalAmount)
		assert.Equal(t, 3.8, order.TotalWeight) // (2*1.5) + (1*0.8)
		assert.Len(t, order.Items, 2)
		assert.Equal(t, shippingAddr, order.ShippingAddr)
		assert.Equal(t, billingAddr, order.BillingAddr)
		assert.Equal(t, customerDetails, *order.CustomerDetails)
		assert.False(t, order.IsGuestOrder)
	})

	t.Run("NewOrder validation errors", func(t *testing.T) {
		validItems := []OrderItem{
			{ProductID: 1, ProductName: "Test", SKU: "SKU-001", Quantity: 1, Price: 9999, Weight: 1.0},
		}
		validAddr := Address{Street: "123 Main St", City: "City", Country: "US"}
		validCustomer := CustomerDetails{Email: "test@example.com", FullName: "John Doe"}

		tests := []struct {
			name          string
			userID        uint
			items         []OrderItem
			currency      string
			expectedError string
		}{
			{
				name:          "zero user ID",
				userID:        0,
				items:         validItems,
				currency:      "USD",
				expectedError: "user ID cannot be empty",
			},
			{
				name:          "empty items",
				userID:        1,
				items:         []OrderItem{},
				currency:      "USD",
				expectedError: "order must have at least one item",
			},
			{
				name:          "empty currency",
				userID:        1,
				items:         validItems,
				currency:      "",
				expectedError: "currency cannot be empty",
			},
			{
				name:   "invalid item quantity",
				userID: 1,
				items: []OrderItem{
					{ProductID: 1, ProductName: "Test", SKU: "SKU-001", Quantity: 0, Price: 9999, Weight: 1.0},
				},
				currency:      "USD",
				expectedError: "item quantity must be greater than zero",
			},
			{
				name:   "invalid item price",
				userID: 1,
				items: []OrderItem{
					{ProductID: 1, ProductName: "Test", SKU: "SKU-001", Quantity: 1, Price: 0, Weight: 1.0},
				},
				currency:      "USD",
				expectedError: "item price must be greater than zero",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				order, err := NewOrder(tt.userID, tt.items, tt.currency, validAddr, validAddr, validCustomer)
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, order)
			})
		}
	})

	t.Run("NewGuestOrder success", func(t *testing.T) {
		items := []OrderItem{
			{
				ProductID:   1,
				ProductName: "Test Product",
				SKU:         "SKU-001",
				Quantity:    1,
				Price:       9999,
				Weight:      1.5,
			},
		}

		shippingAddr := Address{Street: "123 Main St", City: "City", Country: "US"}
		billingAddr := Address{Street: "456 Oak Ave", City: "City", Country: "US"}
		customerDetails := CustomerDetails{Email: "guest@example.com", FullName: "Guest User"}

		order, err := NewGuestOrder(items, shippingAddr, billingAddr, customerDetails)

		require.NoError(t, err)
		assert.Contains(t, order.OrderNumber, "GS-")
		assert.Equal(t, uint(0), order.UserID)
		assert.True(t, order.IsGuestOrder)
		assert.Equal(t, int64(9999), order.TotalAmount)
		assert.Equal(t, 1.5, order.TotalWeight)
	})
}

func TestOrderStatusConstants(t *testing.T) {
	assert.Equal(t, OrderStatus("pending"), OrderStatusPending)
	assert.Equal(t, OrderStatus("paid"), OrderStatusPaid)
	assert.Equal(t, OrderStatus("shipped"), OrderStatusShipped)
	assert.Equal(t, OrderStatus("cancelled"), OrderStatusCancelled)
	assert.Equal(t, OrderStatus("completed"), OrderStatusCompleted)
}

func TestPaymentStatusConstants(t *testing.T) {
	assert.Equal(t, PaymentStatus("pending"), PaymentStatusPending)
	assert.Equal(t, PaymentStatus("authorized"), PaymentStatusAuthorized)
	assert.Equal(t, PaymentStatus("captured"), PaymentStatusCaptured)
	assert.Equal(t, PaymentStatus("refunded"), PaymentStatusRefunded)
	assert.Equal(t, PaymentStatus("cancelled"), PaymentStatusCancelled)
	assert.Equal(t, PaymentStatus("failed"), PaymentStatusFailed)
}
