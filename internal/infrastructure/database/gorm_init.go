package database

import (
	"fmt"
	"log"
	"strings"

	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// InitDB initializes the GORM database connection and auto-migrates tables
func InitDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	var db *gorm.DB
	var err error

	// Choose database driver based on configuration
	switch strings.ToLower(cfg.Driver) {
	case "sqlite":
		db, err = initSQLiteDB(cfg)
	case "postgres":
		db, err = initPostgresDB(cfg)
	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Auto-migrate the schema
	if err := autoMigrate(db); err != nil {
		return nil, fmt.Errorf("failed to auto-migrate: %w", err)
	}

	log.Printf("Database connected (%s) and migrated successfully", cfg.Driver)
	return db, nil
}

// initSQLiteDB initializes SQLite database connection
func initSQLiteDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
	dbPath := cfg.DBName
	if dbPath == "" {
		dbPath = "commercify.db"
	}

	// Open SQLite database connection
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to SQLite database: %w", err)
	}

	// Enable foreign key constraints for SQLite
	if err := db.Exec("PRAGMA foreign_keys = ON").Error; err != nil {
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	return db, nil
}

// initPostgresDB initializes PostgreSQL database connection
func initPostgresDB(cfg config.DatabaseConfig) (*gorm.DB, error) {
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
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL database: %w", err)
	}

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
		&entity.PaymentProvider{},
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
