package gorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/testutil"
)

func TestCheckoutRepository_Update_RemoveItems(t *testing.T) {
	db := testutil.SetupTestDB(t)
	checkoutRepo := NewCheckoutRepository(db)

	t.Run("Remove item should delete from database", func(t *testing.T) {
		// Create a checkout with items
		checkout, err := entity.NewCheckout("session123", "USD")
		require.NoError(t, err)

		// Add multiple items
		err = checkout.AddItem(1, 1, 2, 1000, 1.5, "Product 1", "Variant 1", "SKU-001")
		require.NoError(t, err)
		err = checkout.AddItem(2, 2, 1, 2000, 2.0, "Product 2", "Variant 2", "SKU-002")
		require.NoError(t, err)
		err = checkout.AddItem(3, 3, 3, 3000, 3.0, "Product 3", "Variant 3", "SKU-003")
		require.NoError(t, err)

		// Create checkout in database
		err = checkoutRepo.Create(checkout)
		require.NoError(t, err)
		require.NotZero(t, checkout.ID)

		// Verify all items were created
		retrievedCheckout, err := checkoutRepo.GetByID(checkout.ID)
		require.NoError(t, err)
		require.Len(t, retrievedCheckout.Items, 3)

		// Store item IDs for verification
		var itemIDs []uint
		for _, item := range retrievedCheckout.Items {
			itemIDs = append(itemIDs, item.ID)
		}

		// Remove one item from the checkout entity
		err = retrievedCheckout.RemoveItem(2, 2) // Remove Product 2
		require.NoError(t, err)
		require.Len(t, retrievedCheckout.Items, 2)

		// Update checkout in database
		err = checkoutRepo.Update(retrievedCheckout)
		require.NoError(t, err)

		// Verify the checkout now has only 2 items
		updatedCheckout, err := checkoutRepo.GetByID(checkout.ID)
		require.NoError(t, err)
		require.Len(t, updatedCheckout.Items, 2)

		// Verify the correct item was removed (Product 2 should be gone)
		productIDs := make(map[uint]bool)
		for _, item := range updatedCheckout.Items {
			productIDs[item.ProductID] = true
		}
		assert.True(t, productIDs[1], "Product 1 should still exist")
		assert.False(t, productIDs[2], "Product 2 should be removed")
		assert.True(t, productIDs[3], "Product 3 should still exist")

		// Verify the item was actually deleted from the database
		var itemCount int64
		err = db.Model(&entity.CheckoutItem{}).Where("checkout_id = ?", checkout.ID).Count(&itemCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(2), itemCount, "Should have exactly 2 items in database")

		// Verify the specific item with Product ID 2 is deleted
		var deletedItemCount int64
		err = db.Model(&entity.CheckoutItem{}).Where("checkout_id = ? AND product_id = ?", checkout.ID, 2).Count(&deletedItemCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), deletedItemCount, "Product 2 item should be deleted from database")
	})

	t.Run("Remove multiple items should delete all from database", func(t *testing.T) {
		// Create a checkout with items
		checkout, err := entity.NewCheckout("session456", "USD")
		require.NoError(t, err)

		// Add multiple items
		err = checkout.AddItem(1, 1, 2, 1000, 1.5, "Product 1", "Variant 1", "SKU-001")
		require.NoError(t, err)
		err = checkout.AddItem(2, 2, 1, 2000, 2.0, "Product 2", "Variant 2", "SKU-002")
		require.NoError(t, err)
		err = checkout.AddItem(3, 3, 3, 3000, 3.0, "Product 3", "Variant 3", "SKU-003")
		require.NoError(t, err)

		// Create checkout in database
		err = checkoutRepo.Create(checkout)
		require.NoError(t, err)
		require.NotZero(t, checkout.ID)

		// Get the checkout and remove multiple items
		retrievedCheckout, err := checkoutRepo.GetByID(checkout.ID)
		require.NoError(t, err)
		require.Len(t, retrievedCheckout.Items, 3)

		// Remove two items
		err = retrievedCheckout.RemoveItem(1, 1) // Remove Product 1
		require.NoError(t, err)
		err = retrievedCheckout.RemoveItem(3, 3) // Remove Product 3
		require.NoError(t, err)
		require.Len(t, retrievedCheckout.Items, 1)

		// Update checkout in database
		err = checkoutRepo.Update(retrievedCheckout)
		require.NoError(t, err)

		// Verify the checkout now has only 1 item
		updatedCheckout, err := checkoutRepo.GetByID(checkout.ID)
		require.NoError(t, err)
		require.Len(t, updatedCheckout.Items, 1)

		// Verify only Product 2 remains
		assert.Equal(t, uint(2), updatedCheckout.Items[0].ProductID)

		// Verify the correct count in database
		var itemCount int64
		err = db.Model(&entity.CheckoutItem{}).Where("checkout_id = ?", checkout.ID).Count(&itemCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(1), itemCount, "Should have exactly 1 item in database")
	})

	t.Run("Remove all items should clear database", func(t *testing.T) {
		// Create a checkout with items
		checkout, err := entity.NewCheckout("session789", "USD")
		require.NoError(t, err)

		// Add items
		err = checkout.AddItem(1, 1, 2, 1000, 1.5, "Product 1", "Variant 1", "SKU-001")
		require.NoError(t, err)
		err = checkout.AddItem(2, 2, 1, 2000, 2.0, "Product 2", "Variant 2", "SKU-002")
		require.NoError(t, err)

		// Create checkout in database
		err = checkoutRepo.Create(checkout)
		require.NoError(t, err)

		// Get the checkout and clear all items
		retrievedCheckout, err := checkoutRepo.GetByID(checkout.ID)
		require.NoError(t, err)

		// Clear all items using the Clear method
		retrievedCheckout.Clear()
		require.Len(t, retrievedCheckout.Items, 0)

		// Update checkout in database
		err = checkoutRepo.Update(retrievedCheckout)
		require.NoError(t, err)

		// Verify the checkout has no items
		updatedCheckout, err := checkoutRepo.GetByID(checkout.ID)
		require.NoError(t, err)
		require.Len(t, updatedCheckout.Items, 0)

		// Verify no items exist in database for this checkout
		var itemCount int64
		err = db.Model(&entity.CheckoutItem{}).Where("checkout_id = ?", checkout.ID).Count(&itemCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), itemCount, "Should have no items in database")
	})

	t.Run("Update existing items without removal should work", func(t *testing.T) {
		// Create a checkout with items
		checkout, err := entity.NewCheckout("session999", "USD")
		require.NoError(t, err)

		// Add items
		err = checkout.AddItem(1, 1, 2, 1000, 1.5, "Product 1", "Variant 1", "SKU-001")
		require.NoError(t, err)
		err = checkout.AddItem(2, 2, 1, 2000, 2.0, "Product 2", "Variant 2", "SKU-002")
		require.NoError(t, err)

		// Create checkout in database
		err = checkoutRepo.Create(checkout)
		require.NoError(t, err)

		// Get the checkout and update an item quantity
		retrievedCheckout, err := checkoutRepo.GetByID(checkout.ID)
		require.NoError(t, err)

		// Update item quantity
		err = retrievedCheckout.UpdateItem(1, 1, 5) // Change quantity from 2 to 5
		require.NoError(t, err)

		// Update checkout in database
		err = checkoutRepo.Update(retrievedCheckout)
		require.NoError(t, err)

		// Verify the checkout still has 2 items with updated quantity
		updatedCheckout, err := checkoutRepo.GetByID(checkout.ID)
		require.NoError(t, err)
		require.Len(t, updatedCheckout.Items, 2)

		// Find the updated item
		var updatedItem *entity.CheckoutItem
		for _, item := range updatedCheckout.Items {
			if item.ProductID == 1 {
				updatedItem = &item
				break
			}
		}
		require.NotNil(t, updatedItem)
		assert.Equal(t, 5, updatedItem.Quantity)

		// Verify count in database remains the same
		var itemCount int64
		err = db.Model(&entity.CheckoutItem{}).Where("checkout_id = ?", checkout.ID).Count(&itemCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(2), itemCount, "Should still have 2 items in database")
	})
}
