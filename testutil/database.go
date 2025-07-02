package testutil

import (
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// SetupTestDB creates an in-memory SQLite database for testing
// It automatically migrates all entities and returns the database connection
func SetupTestDB(t *testing.T) *gorm.DB {
	// Create in-memory SQLite database
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.New(
			log.New(log.Writer(), "\r\n", log.LstdFlags), // Use Go's standard logger
			logger.Config{
				LogLevel: logger.Silent, // Set to Silent to reduce test output noise
			},
		),
	})
	require.NoError(t, err, "Failed to connect to test database")

	// Auto-migrate all entities
	err = autoMigrate(db)
	require.NoError(t, err, "Failed to migrate test database")

	return db
}

// SetupTestDBWithLogger creates an in-memory SQLite database with custom logging level
func SetupTestDBWithLogger(t *testing.T, logLevel logger.LogLevel) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.New(
			log.New(log.Writer(), "\r\n", log.LstdFlags),
			logger.Config{
				LogLevel: logLevel,
			},
		),
	})
	require.NoError(t, err, "Failed to connect to test database")

	err = autoMigrate(db)
	require.NoError(t, err, "Failed to migrate test database")

	return db
}

// autoMigrate performs automatic migration of all entities
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		// Core entities
		&entity.User{},
		&entity.Category{},

		// Product entities
		&entity.Product{},
		&entity.ProductVariant{},
		&entity.Currency{},

		// Order entities
		&entity.Order{},
		&entity.OrderItem{},

		// Checkout entities
		&entity.Checkout{},
		&entity.CheckoutItem{},

		// Discount entities
		&entity.Discount{},

		// Shipping entities
		&entity.ShippingMethod{},
		&entity.ShippingZone{},
		&entity.ShippingRate{},
		&entity.WeightBasedRate{},
		&entity.ValueBasedRate{},

		// Payment entities
		&entity.PaymentTransaction{},
		// Skip PaymentProvider for now due to slice field issues
		// &entity.PaymentProvider{},
	)
}

// CreateTestOrder creates a test order with the given ID
func CreateTestOrder(t *testing.T, db *gorm.DB, orderID uint) *entity.Order {
	order := &entity.Order{
		Model:         gorm.Model{ID: orderID},
		OrderNumber:   fmt.Sprintf("ORD-%d", orderID),
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

// CreateTestUser creates a test user with the given ID
func CreateTestUser(t *testing.T, db *gorm.DB, userID uint) *entity.User {
	user := &entity.User{
		Model:     gorm.Model{ID: userID},
		Email:     fmt.Sprintf("user%d@example.com", userID),
		Password:  "hashedpassword", // In real tests, you might want to hash this
		FirstName: fmt.Sprintf("User%d", userID),
		LastName:  "TestUser",
		Role:      "user",
	}
	err := db.Create(user).Error
	require.NoError(t, err)
	return user
}

// CreateTestProduct creates a test product with the given ID
func CreateTestProduct(t *testing.T, db *gorm.DB, productID uint) *entity.Product {
	// First create a test category
	category := CreateTestCategory(t, db, productID) // Use the same ID for simplicity

	product := &entity.Product{
		Model:       gorm.Model{ID: productID},
		Name:        fmt.Sprintf("Test Product %d", productID),
		Description: fmt.Sprintf("Test product %d description", productID),
		Currency:    "USD",
		CategoryID:  category.ID,
		Active:      true,
	}
	err := db.Create(product).Error
	require.NoError(t, err)
	return product
}

// CreateTestCategory creates a test category with the given ID
func CreateTestCategory(t *testing.T, db *gorm.DB, categoryID uint) *entity.Category {
	category := &entity.Category{
		Model:       gorm.Model{ID: categoryID},
		Name:        fmt.Sprintf("Test Category %d", categoryID),
		Description: fmt.Sprintf("Test category %d description", categoryID),
	}
	err := db.Create(category).Error
	require.NoError(t, err)
	return category
}

// CreateTestPaymentProvider creates a test payment provider
func CreateTestPaymentProvider(t *testing.T, db *gorm.DB, providerID uint, name string) *entity.PaymentProvider {
	provider := &entity.PaymentProvider{
		Model:   gorm.Model{ID: providerID},
		Type:    common.PaymentProviderMock, // Use mock type for testing
		Name:    name,
		Enabled: true,
	}
	err := db.Create(provider).Error
	require.NoError(t, err)
	return provider
}

// CleanupTestDB closes the database connection and cleans up resources
func CleanupTestDB(t *testing.T, db *gorm.DB) {
	sqlDB, err := db.DB()
	require.NoError(t, err)
	err = sqlDB.Close()
	require.NoError(t, err)
}

// TruncateAllTables removes all data from all tables (useful for test isolation)
func TruncateAllTables(t *testing.T, db *gorm.DB) {
	tables := []string{
		"payment_transactions",
		// "payment_providers", // Commented out since we don't migrate this entity
		"order_items",
		"orders",
		"checkout_items",
		"checkouts",
		"product_variants",
		"products",
		"categories",
		"users",
		"discounts",
		"shipping_methods",
		"shipping_zones",
		"shipping_rates",
		"weight_based_rates",
		"value_based_rates",
	}

	// Disable foreign key checks temporarily
	db.Exec("PRAGMA foreign_keys = OFF")

	for _, table := range tables {
		// Check if table exists before trying to truncate
		if db.Migrator().HasTable(table) {
			err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error
			require.NoError(t, err, fmt.Sprintf("Failed to truncate table %s", table))
		}
	}

	// Re-enable foreign key checks
	db.Exec("PRAGMA foreign_keys = ON")
}
