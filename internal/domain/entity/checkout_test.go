package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckout(t *testing.T) {
	t.Run("NewCheckout success", func(t *testing.T) {
		checkout, err := NewCheckout("session123", "USD")

		require.NoError(t, err)
		assert.Equal(t, "session123", checkout.SessionID)
		assert.Equal(t, "USD", checkout.Currency)
		assert.Equal(t, CheckoutStatusActive, checkout.Status)
		assert.Equal(t, int64(0), checkout.TotalAmount)
		assert.Equal(t, int64(0), checkout.ShippingCost)
		assert.Equal(t, int64(0), checkout.DiscountAmount)
		assert.Equal(t, int64(0), checkout.FinalAmount)
		assert.NotNil(t, checkout.Items)
		assert.Len(t, checkout.Items, 0)
		assert.False(t, checkout.LastActivityAt.IsZero())
		assert.False(t, checkout.ExpiresAt.IsZero())
		assert.True(t, checkout.ExpiresAt.After(checkout.LastActivityAt))
	})

	t.Run("NewCheckout validation errors", func(t *testing.T) {
		tests := []struct {
			name        string
			sessionID   string
			currency    string
			expectedErr string
		}{
			{
				name:        "empty session ID",
				sessionID:   "",
				currency:    "USD",
				expectedErr: "session ID cannot be empty",
			},
			{
				name:        "empty currency",
				sessionID:   "session123",
				currency:    "",
				expectedErr: "currency cannot be empty",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				checkout, err := NewCheckout(tt.sessionID, tt.currency)
				assert.Nil(t, checkout)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			})
		}
	})

	t.Run("AddItem", func(t *testing.T) {
		checkout, err := NewCheckout("session123", "USD")
		require.NoError(t, err)

		// Add first item
		err = checkout.AddItem(1, 1, 2, 9999, 1.5, "Test Product", "Size M", "SKU-001")
		assert.NoError(t, err)
		assert.Len(t, checkout.Items, 1)
		assert.Equal(t, uint(1), checkout.Items[0].ProductID)
		assert.Equal(t, uint(1), checkout.Items[0].ProductVariantID)
		assert.Equal(t, 2, checkout.Items[0].Quantity)
		assert.Equal(t, int64(9999), checkout.Items[0].Price)
		assert.Equal(t, "Test Product", checkout.Items[0].ProductName)
		assert.Equal(t, "SKU-001", checkout.Items[0].SKU)

		// Add same item again (should update quantity)
		err = checkout.AddItem(1, 1, 1, 9999, 1.5, "Test Product", "Size M", "SKU-001")
		assert.NoError(t, err)
		assert.Len(t, checkout.Items, 1)               // Still only one item
		assert.Equal(t, 3, checkout.Items[0].Quantity) // Quantity should be updated

		// Add different item
		err = checkout.AddItem(2, 2, 1, 19999, 2.0, "Another Product", "Size L", "SKU-002")
		assert.NoError(t, err)
		assert.Len(t, checkout.Items, 2) // Now two items
	})

	t.Run("AddItem validation errors", func(t *testing.T) {
		checkout, err := NewCheckout("session123", "USD")
		require.NoError(t, err)

		tests := []struct {
			name        string
			productID   uint
			variantID   uint
			quantity    int
			price       int64
			expectedErr string
		}{
			{
				name:        "zero product ID",
				productID:   0,
				variantID:   1,
				quantity:    1,
				price:       9999,
				expectedErr: "product ID cannot be empty",
			},
			{
				name:        "zero quantity",
				productID:   1,
				variantID:   1,
				quantity:    0,
				price:       9999,
				expectedErr: "quantity must be greater than zero",
			},
			{
				name:        "negative quantity",
				productID:   1,
				variantID:   1,
				quantity:    -1,
				price:       9999,
				expectedErr: "quantity must be greater than zero",
			},
			{
				name:        "negative price",
				productID:   1,
				variantID:   1,
				quantity:    1,
				price:       -100,
				expectedErr: "price cannot be negative",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				err := checkout.AddItem(tt.productID, tt.variantID, tt.quantity, tt.price, 1.0, "Product", "Variant", "SKU")
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			})
		}
	})

	t.Run("UpdateItem", func(t *testing.T) {
		checkout, err := NewCheckout("session123", "USD")
		require.NoError(t, err)

		// Add an item first
		err = checkout.AddItem(1, 1, 2, 9999, 1.5, "Test Product", "Size M", "SKU-001")
		require.NoError(t, err)

		// Update the item quantity
		err = checkout.UpdateItem(1, 1, 5)
		assert.NoError(t, err)
		assert.Equal(t, 5, checkout.Items[0].Quantity)

		// Try to update non-existent item
		err = checkout.UpdateItem(999, 999, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product not found in checkout")
	})

	t.Run("RemoveItem", func(t *testing.T) {
		checkout, err := NewCheckout("session123", "USD")
		require.NoError(t, err)

		// Add items first
		err = checkout.AddItem(1, 1, 2, 9999, 1.5, "Product 1", "Variant 1", "SKU-001")
		require.NoError(t, err)
		err = checkout.AddItem(2, 2, 1, 19999, 2.0, "Product 2", "Variant 2", "SKU-002")
		require.NoError(t, err)

		assert.Len(t, checkout.Items, 2)

		// Remove one item
		err = checkout.RemoveItem(1, 1)
		assert.NoError(t, err)
		assert.Len(t, checkout.Items, 1)
		assert.Equal(t, uint(2), checkout.Items[0].ProductID)

		// Try to remove non-existent item
		err = checkout.RemoveItem(999, 999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "product not found in checkout")
	})

	t.Run("TotalItems", func(t *testing.T) {
		checkout, err := NewCheckout("session123", "USD")
		require.NoError(t, err)

		assert.Equal(t, 0, checkout.TotalItems())

		// Add items
		err = checkout.AddItem(1, 1, 2, 9999, 1.5, "Product 1", "Variant 1", "SKU-001")
		require.NoError(t, err)
		err = checkout.AddItem(2, 2, 3, 19999, 2.0, "Product 2", "Variant 2", "SKU-002")
		require.NoError(t, err)

		assert.Equal(t, 5, checkout.TotalItems()) // 2 + 3
	})

	t.Run("MarkAsCompleted", func(t *testing.T) {
		checkout, err := NewCheckout("session123", "USD")
		require.NoError(t, err)

		assert.Equal(t, CheckoutStatusActive, checkout.Status)
		assert.Nil(t, checkout.CompletedAt)

		checkout.MarkAsCompleted(123)
		assert.Equal(t, CheckoutStatusCompleted, checkout.Status)
		assert.NotNil(t, checkout.CompletedAt)
		assert.Equal(t, uint(123), *checkout.ConvertedOrderID)
	})

	t.Run("MarkAsAbandoned", func(t *testing.T) {
		checkout, err := NewCheckout("session123", "USD")
		require.NoError(t, err)

		assert.Equal(t, CheckoutStatusActive, checkout.Status)

		checkout.MarkAsAbandoned()
		assert.Equal(t, CheckoutStatusAbandoned, checkout.Status)
	})

	t.Run("MarkAsExpired", func(t *testing.T) {
		checkout, err := NewCheckout("session123", "USD")
		require.NoError(t, err)

		assert.Equal(t, CheckoutStatusActive, checkout.Status)

		checkout.MarkAsExpired()
		assert.Equal(t, CheckoutStatusExpired, checkout.Status)
	})

	t.Run("IsExpired", func(t *testing.T) {
		checkout, err := NewCheckout("session123", "USD")
		require.NoError(t, err)

		// Should not be expired initially
		assert.False(t, checkout.IsExpired())

		// Set expiry to past
		checkout.ExpiresAt = time.Now().Add(-1 * time.Hour)
		assert.True(t, checkout.IsExpired())
	})

	t.Run("ExtendExpiry", func(t *testing.T) {
		checkout, err := NewCheckout("session123", "USD")
		require.NoError(t, err)

		beforeExtend := time.Now()

		checkout.ExtendExpiry(2 * time.Hour)

		// The new expiry should be around 2 hours from now, not from the original expiry
		expectedMin := beforeExtend.Add(1*time.Hour + 50*time.Minute)
		expectedMax := beforeExtend.Add(2*time.Hour + 10*time.Minute)

		assert.True(t, checkout.ExpiresAt.After(expectedMin))
		assert.True(t, checkout.ExpiresAt.Before(expectedMax))
	})

	t.Run("HasCustomerInfo", func(t *testing.T) {
		checkout, err := NewCheckout("session123", "USD")
		require.NoError(t, err)

		// Initially no customer info
		assert.False(t, checkout.HasCustomerInfo())

		// Set customer details
		checkout.SetCustomerDetails(CustomerDetails{
			Email:    "test@example.com",
			FullName: "John Doe",
		})
		assert.True(t, checkout.HasCustomerInfo())
	})

	t.Run("IsEmpty", func(t *testing.T) {
		checkout, err := NewCheckout("session123", "USD")
		require.NoError(t, err)

		// Initially empty
		assert.True(t, checkout.IsEmpty())

		// Add an item
		err = checkout.AddItem(1, 1, 1, 9999, 1.5, "Product", "Variant", "SKU-001")
		require.NoError(t, err)
		assert.False(t, checkout.IsEmpty())

		// Remove the item but add customer info
		err = checkout.RemoveItem(1, 1)
		require.NoError(t, err)
		checkout.SetCustomerDetails(CustomerDetails{Email: "test@example.com"})
		assert.False(t, checkout.IsEmpty()) // Still not empty due to customer info
	})

	t.Run("Clear", func(t *testing.T) {
		checkout, err := NewCheckout("session123", "USD")
		require.NoError(t, err)

		// Add items and set details
		err = checkout.AddItem(1, 1, 2, 9999, 1.5, "Product", "Variant", "SKU-001")
		require.NoError(t, err)
		checkout.SetCustomerDetails(CustomerDetails{Email: "test@example.com"})

		assert.False(t, checkout.IsEmpty())

		// Clear the checkout
		checkout.Clear()

		assert.Len(t, checkout.Items, 0)
		assert.Equal(t, int64(0), checkout.TotalAmount)
		assert.Equal(t, int64(0), checkout.FinalAmount)
		assert.Equal(t, int64(0), checkout.DiscountAmount)
		assert.Equal(t, int64(0), checkout.ShippingCost)
		assert.Equal(t, "", checkout.DiscountCode)
	})
}

func TestCheckoutStatus(t *testing.T) {
	t.Run("CheckoutStatus constants", func(t *testing.T) {
		assert.Equal(t, CheckoutStatus("active"), CheckoutStatusActive)
		assert.Equal(t, CheckoutStatus("completed"), CheckoutStatusCompleted)
		assert.Equal(t, CheckoutStatus("abandoned"), CheckoutStatusAbandoned)
		assert.Equal(t, CheckoutStatus("expired"), CheckoutStatusExpired)
	})
}

func TestCheckoutDTOConversions(t *testing.T) {
	t.Run("ToCheckoutDTO", func(t *testing.T) {
		checkout, err := NewCheckout("session123", "USD")
		require.NoError(t, err)

		// Add some items to the checkout
		err = checkout.AddItem(1, 1, 2, 9999, 1.5, "Test Product", "Test Variant", "SKU-001")
		require.NoError(t, err)

		// Mock ID that would be set by GORM
		checkout.ID = 1

		dto := checkout.ToCheckoutDTO()
		assert.Equal(t, uint(1), dto.ID)
		assert.Equal(t, "session123", dto.SessionID)
		assert.Equal(t, "USD", dto.Currency)
		assert.Equal(t, string(CheckoutStatusActive), dto.Status)
		assert.Equal(t, float64(199.98), dto.TotalAmount) // 2 * 99.99 (converted from cents)
		assert.Equal(t, float64(0), dto.ShippingCost)
		assert.Equal(t, float64(0), dto.DiscountAmount)
		assert.Equal(t, float64(199.98), dto.FinalAmount)
		assert.Equal(t, float64(3.0), dto.TotalWeight) // 2 * 1.5
		assert.NotNil(t, dto.Items)
		assert.Len(t, dto.Items, 1)
		assert.Equal(t, uint(1), dto.Items[0].ProductID)
		assert.Equal(t, "Test Product", dto.Items[0].ProductName)
		assert.Equal(t, 2, dto.Items[0].Quantity)
		assert.Equal(t, float64(99.99), dto.Items[0].Price)
		assert.False(t, dto.LastActivityAt.IsZero())
		assert.False(t, dto.ExpiresAt.IsZero())
	})
}
