package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/infrastructure/database"
	"gorm.io/gorm"
)

func main() {
	// Define command line flags
	allFlag := flag.Bool("all", false, "Seed all data")
	usersFlag := flag.Bool("users", false, "Seed users data")
	categoriesFlag := flag.Bool("categories", false, "Seed categories data")
	productsFlag := flag.Bool("products", false, "Seed products data")
	productVariantsFlag := flag.Bool("product-variants", false, "Seed product variants data")
	discountsFlag := flag.Bool("discounts", false, "Seed discounts data")
	ordersFlag := flag.Bool("orders", false, "Seed orders data")
	checkoutsFlag := flag.Bool("checkouts", false, "Seed checkouts data")
	paymentTransactionsFlag := flag.Bool("payment-transactions", false, "Seed payment transactions data")
	shippingFlag := flag.Bool("shipping", false, "Seed shipping data (methods, zones, rates)")
	clearFlag := flag.Bool("clear", false, "Clear all data before seeding")
	flag.Parse()

	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using environment variables")
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Connect to database
	db, err := database.InitDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close(db)

	// Clear data if requested
	if *clearFlag {
		if err := clearData(db); err != nil {
			log.Fatalf("Failed to clear data: %v", err)
		}
		fmt.Println("All data cleared")
	}

	// Seed data based on flags
	if *allFlag || *usersFlag {
		if err := seedUsers(db); err != nil {
			log.Fatalf("Failed to seed users: %v", err)
		}
		fmt.Println("Users seeded successfully")
	}

	if *allFlag || *categoriesFlag {
		if err := seedCategories(db); err != nil {
			log.Fatalf("Failed to seed categories: %v", err)
		}
		fmt.Println("Categories seeded successfully")
	}

	if *allFlag || *productsFlag {
		if err := seedProducts(db); err != nil {
			log.Fatalf("Failed to seed products: %v", err)
		}
		fmt.Println("Products seeded successfully")
	}

	if *allFlag || *productVariantsFlag {
		if err := seedProductVariants(db); err != nil {
			log.Fatalf("Failed to seed product variants: %v", err)
		}
		fmt.Println("Product variants seeded successfully")
	}

	if *allFlag || *discountsFlag {
		if err := seedDiscounts(db); err != nil {
			log.Fatalf("Failed to seed discounts: %v", err)
		}
		fmt.Println("Discounts seeded successfully")
	}

	if *allFlag || *shippingFlag {
		if err := seedShippingMethods(db); err != nil {
			log.Fatalf("Failed to seed shipping methods: %v", err)
		}
		fmt.Println("Shipping methods seeded successfully")

		if err := seedShippingZones(db); err != nil {
			log.Fatalf("Failed to seed shipping zones: %v", err)
		}
		fmt.Println("Shipping zones seeded successfully")

		if err := seedShippingRates(db); err != nil {
			log.Fatalf("Failed to seed shipping rates: %v", err)
		}
		fmt.Println("Shipping rates seeded successfully")
	}

	if *allFlag || *ordersFlag {
		if err := seedOrders(db); err != nil {
			log.Fatalf("Failed to seed orders: %v", err)
		}
		fmt.Println("Orders seeded successfully")
	}

	if *allFlag || *checkoutsFlag {
		if err := seedCheckouts(db); err != nil {
			log.Fatalf("Failed to seed checkouts: %v", err)
		}
		fmt.Println("Checkouts seeded successfully")
	}

	// if *allFlag || *paymentTransactionsFlag {
	// 	if err := seedPaymentTransactions(db); err != nil {
	// 		log.Fatalf("Failed to seed payment transactions: %v", err)
	// 	}
	// 	fmt.Println("Payment transactions seeded successfully")
	// }

	if !*allFlag && !*usersFlag && !*categoriesFlag && !*productsFlag && !*productVariantsFlag &&
		!*ordersFlag && !*checkoutsFlag && !*clearFlag && !*discountsFlag &&
		!*paymentTransactionsFlag && !*shippingFlag {
		fmt.Println("No action specified")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
	}
}

// clearData clears all data from the database
func clearData(db *gorm.DB) error {
	tables := []string{
		"payment_transactions",
		"shipping_rates",
		"shipping_zones",
		"shipping_methods",
		"orders",
		"checkouts",
		"discounts",
		"product_variants",
		"products",
		"categories",
		"users",
	}

	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table)).Error; err != nil {
			return fmt.Errorf("failed to clear table %s: %w", table, err)
		}
	}

	return nil
}

// seedUsers seeds user data
func seedUsers(db *gorm.DB) error {
	return nil
}

// seedCategories seeds category data
func seedCategories(db *gorm.DB) error {
	return nil
}

// seedProducts seeds product data
func seedProducts(db *gorm.DB) error {
	return nil
}

// seedProductVariants seeds product variant data
func seedProductVariants(db *gorm.DB) error {
	return nil
}

// seedOrders seeds order data
func seedOrders(db *gorm.DB) error {
	return nil
}

// seedDiscounts seeds discount data
func seedDiscounts(db *gorm.DB) error {
	return nil
}

// seedShippingMethods seeds shipping method data
func seedShippingMethods(db *gorm.DB) error {
	return nil
}

// seedShippingZones seeds shipping zone data
func seedShippingZones(db *gorm.DB) error {
	return nil
}

// seedShippingRates seeds shipping rate data
func seedShippingRates(db *gorm.DB) error {
	return nil
}

// seedPaymentTransactions seeds payment transaction data
func seedPaymentTransactions(db *gorm.DB) error {
	return nil
}

// seedCheckouts seeds checkout data for testing expiry and cleanup logic
func seedCheckouts(db *gorm.DB) error {
	return nil
}
