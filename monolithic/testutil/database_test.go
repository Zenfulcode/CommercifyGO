package testutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

func TestSetupTestDB(t *testing.T) {
	t.Run("SetupTestDB creates working database", func(t *testing.T) {
		db := SetupTestDB(t)
		defer CleanupTestDB(t, db)

		// Test that we can create and retrieve an order
		order := CreateTestOrder(t, db, 1)
		assert.Equal(t, uint(1), order.ID)
		assert.Equal(t, "ORD-1", order.OrderNumber)

		// Test that we can retrieve the order
		var retrievedOrder entity.Order
		err := db.First(&retrievedOrder, 1).Error
		require.NoError(t, err)
		assert.Equal(t, order.ID, retrievedOrder.ID)
		assert.Equal(t, order.OrderNumber, retrievedOrder.OrderNumber)
	})

	t.Run("CreateTestUser creates valid user", func(t *testing.T) {
		db := SetupTestDB(t)
		defer CleanupTestDB(t, db)

		user := CreateTestUser(t, db, 1)
		assert.Equal(t, uint(1), user.ID)
		assert.Equal(t, "user1@example.com", user.Email)
		assert.Equal(t, "User1", user.FirstName)
		assert.Equal(t, "TestUser", user.LastName)
		assert.Equal(t, "user", user.Role)
	})

	t.Run("CreateTestCategory creates valid category", func(t *testing.T) {
		db := SetupTestDB(t)
		defer CleanupTestDB(t, db)

		category := CreateTestCategory(t, db, 1)
		assert.Equal(t, uint(1), category.ID)
		assert.Equal(t, "Test Category 1", category.Name)
		assert.Contains(t, category.Description, "Test category 1")
	})

	t.Run("CreateTestProduct creates valid product with category", func(t *testing.T) {
		db := SetupTestDB(t)
		defer CleanupTestDB(t, db)

		product := CreateTestProduct(t, db, 1)
		assert.Equal(t, uint(1), product.ID)
		assert.Equal(t, "Test Product 1", product.Name)
		assert.Equal(t, "USD", product.Currency)
		assert.True(t, product.Active)
		assert.Equal(t, uint(1), product.CategoryID) // Should reference the category created in the function
	})

	t.Run("TruncateAllTables cleans database", func(t *testing.T) {
		db := SetupTestDB(t)
		defer CleanupTestDB(t, db)

		// Create some test data
		CreateTestOrder(t, db, 1)
		CreateTestUser(t, db, 1)

		// Verify data exists
		var orderCount int64
		var userCount int64
		db.Model(&entity.Order{}).Count(&orderCount)
		db.Model(&entity.User{}).Count(&userCount)
		assert.Equal(t, int64(1), orderCount)
		assert.Equal(t, int64(1), userCount)

		// Truncate all tables
		TruncateAllTables(t, db)

		// Verify data is gone
		db.Model(&entity.Order{}).Count(&orderCount)
		db.Model(&entity.User{}).Count(&userCount)
		assert.Equal(t, int64(0), orderCount)
		assert.Equal(t, int64(0), userCount)
	})
}
