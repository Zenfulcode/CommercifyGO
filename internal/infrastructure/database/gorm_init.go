package database

import (
	"fmt"
	"log"

	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB initializes the GORM database connection and auto-migrates tables
func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	// Get database connection details from environment
	host := cfg.Host
	port := cfg.Port
	user := cfg.User
	password := cfg.Password
	dbname := cfg.DBName
	sslmode := cfg.SSLMode

	// Build connection string
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, password, dbname, sslmode)

	// Open database connection
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		// Logger: gormLogger,
		// DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schema
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	log.Println("Database connected and migrated successfully")
	return db, nil
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
		&entity.ProductPrice{},

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

		// Webhook entities
		&entity.Webhook{},
	)
}

// CloseDB closes the database connection
func Close(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}
