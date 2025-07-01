package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
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
	currenciesFlag := flag.Bool("currencies", false, "Seed currencies data")
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
	if *allFlag || *currenciesFlag {
		if err := seedCurrencies(db); err != nil {
			log.Fatalf("Failed to seed currencies: %v", err)
		}
		fmt.Println("Currencies seeded successfully")
	}

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

	if *allFlag || *paymentTransactionsFlag {
		if err := seedPaymentTransactions(db); err != nil {
			log.Fatalf("Failed to seed payment transactions: %v", err)
		}
		fmt.Println("Payment transactions seeded successfully")
	}

	if !*allFlag && !*usersFlag && !*categoriesFlag && !*productsFlag && !*productVariantsFlag &&
		!*ordersFlag && !*checkoutsFlag && !*clearFlag && !*discountsFlag &&
		!*paymentTransactionsFlag && !*shippingFlag && !*currenciesFlag {
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
		"currencies",
	}

	// For SQLite, use DELETE instead of TRUNCATE
	for _, table := range tables {
		if err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)).Error; err != nil {
			return fmt.Errorf("failed to clear table %s: %w", table, err)
		}
	}

	return nil
}

// seedUsers seeds user data
func seedUsers(db *gorm.DB) error {
	users := []struct {
		email     string
		password  string
		firstName string
		lastName  string
		role      string
	}{
		{"admin@example.com", "password123", "Admin", "User", "admin"},
		{"john.doe@example.com", "password123", "John", "Doe", "user"},
		{"jane.smith@example.com", "password123", "Jane", "Smith", "user"},
		{"customer@example.com", "password123", "Test", "Customer", "user"},
	}

	for _, userData := range users {
		// Check if user already exists
		var existingUser struct{ ID uint }
		if err := db.Table("users").Select("id").Where("email = ?", userData.email).First(&existingUser).Error; err == nil {
			continue // User already exists, skip
		}

		// Create new user using entity constructor
		user, err := entity.NewUser(userData.email, userData.password, userData.firstName, userData.lastName, entity.UserRole(userData.role))
		if err != nil {
			return fmt.Errorf("failed to create user %s: %w", userData.email, err)
		}

		if err := db.Create(user).Error; err != nil {
			return fmt.Errorf("failed to save user %s: %w", userData.email, err)
		}
	}

	return nil
}

// seedCurrencies seeds currency data
func seedCurrencies(db *gorm.DB) error {
	currencies := []struct {
		code         string
		name         string
		symbol       string
		exchangeRate float64
		isEnabled    bool
		isDefault    bool
	}{
		{"USD", "US Dollar", "$", 1.0, true, true},       // USD as default base currency
		{"EUR", "Euro", "€", 0.85, true, false},          // Approximate exchange rate
		{"DKK", "Danish Krone", "kr", 6.80, true, false}, // Approximate exchange rate
	}

	for _, currData := range currencies {
		// Check if currency already exists
		var existingCurrency entity.Currency
		if err := db.Where("code = ?", currData.code).First(&existingCurrency).Error; err == nil {
			// Currency exists, update it with our seed data
			existingCurrency.Name = currData.name
			existingCurrency.Symbol = currData.symbol
			existingCurrency.ExchangeRate = currData.exchangeRate
			existingCurrency.IsEnabled = currData.isEnabled
			existingCurrency.IsDefault = currData.isDefault

			if err := db.Save(&existingCurrency).Error; err != nil {
				return fmt.Errorf("failed to update currency %s: %w", currData.code, err)
			}
			continue
		}

		// Create new currency using entity constructor
		currency, err := entity.NewCurrency(
			currData.code,
			currData.name,
			currData.symbol,
			currData.exchangeRate,
			currData.isEnabled,
			currData.isDefault,
		)
		if err != nil {
			return fmt.Errorf("failed to create currency %s: %w", currData.code, err)
		}

		if err := db.Create(currency).Error; err != nil {
			return fmt.Errorf("failed to save currency %s: %w", currData.code, err)
		}
	}

	return nil
}

// seedCategories seeds category data
func seedCategories(db *gorm.DB) error {
	categories := []struct {
		name        string
		description string
		parentID    *uint
	}{
		{"Clothing", "All clothing items", nil},
		{"Electronics", "Electronic devices and accessories", nil},
		{"Home & Garden", "Home and garden products", nil},
		{"Men's Clothing", "Clothing for men", nil}, // Will be updated with parentID after creation
		{"Women's Clothing", "Clothing for women", nil},
		{"Accessories", "Fashion accessories", nil},
	}

	var clothingID uint

	for i, catData := range categories {
		// Check if category already exists
		var existingCat struct{ ID uint }
		if err := db.Table("categories").Select("id").Where("name = ?", catData.name).First(&existingCat).Error; err == nil {
			if catData.name == "Clothing" {
				clothingID = existingCat.ID
			}
			continue // Category already exists, skip
		}

		// Create new category using entity constructor
		category, err := entity.NewCategory(catData.name, catData.description, catData.parentID)
		if err != nil {
			return fmt.Errorf("failed to create category %s: %w", catData.name, err)
		}

		if err := db.Create(category).Error; err != nil {
			return fmt.Errorf("failed to save category %s: %w", catData.name, err)
		}

		// Store clothing ID for subcategories
		if catData.name == "Clothing" {
			clothingID = category.ID
		}

		// Update men's and women's clothing to be children of clothing
		if i == len(categories)-3 && clothingID > 0 { // After creating all main categories
			db.Model(&entity.Category{}).Where("name IN ?", []string{"Men's Clothing", "Women's Clothing"}).Update("parent_id", clothingID)
		}
	}

	return nil
}

// seedProducts seeds product data
func seedProducts(db *gorm.DB) error {
	// First, get category IDs
	var categories []struct {
		ID   uint
		Name string
	}
	if err := db.Table("categories").Select("id, name").Find(&categories).Error; err != nil {
		return fmt.Errorf("failed to fetch categories: %w", err)
	}

	categoryMap := make(map[string]uint)
	for _, cat := range categories {
		categoryMap[cat.Name] = cat.ID
	}

	clothingCategoryID := categoryMap["Clothing"]
	electronicsCategoryID := categoryMap["Electronics"]

	if clothingCategoryID == 0 {
		return fmt.Errorf("clothing category not found")
	}

	products := []struct {
		name        string
		description string
		currency    string
		categoryID  uint
		images      []string
		active      bool
		variants    []struct {
			sku        string
			stock      int
			price      int64 // Price in cents
			weight     float64
			attributes map[string]string
			images     []string
			isDefault  bool
		}
	}{
		{
			name:        "Classic T-Shirt",
			description: "Comfortable cotton t-shirt perfect for everyday wear",
			currency:    "USD",
			categoryID:  clothingCategoryID,
			images:      []string{"tshirt1.jpg", "tshirt2.jpg"},
			active:      true,
			variants: []struct {
				sku        string
				stock      int
				price      int64
				weight     float64
				attributes map[string]string
				images     []string
				isDefault  bool
			}{
				{"Men-B-M", 50, 1999, 0.2, map[string]string{"Color": "Black", "Size": "M", "Gender": "Men"}, []string{}, true},
				{"Men-B-L", 30, 1999, 0.2, map[string]string{"Color": "Black", "Size": "L", "Gender": "Men"}, []string{}, false},
				{"Women-R-M", 40, 1999, 0.18, map[string]string{"Color": "Red", "Size": "M", "Gender": "Women"}, []string{}, false},
				{"Women-R-L", 25, 1999, 0.18, map[string]string{"Color": "Red", "Size": "L", "Gender": "Women"}, []string{}, false},
			},
		},
		{
			name:        "Premium Jeans",
			description: "High-quality denim jeans with perfect fit",
			currency:    "USD",
			categoryID:  clothingCategoryID,
			images:      []string{"jeans1.jpg", "jeans2.jpg"},
			active:      true,
			variants: []struct {
				sku        string
				stock      int
				price      int64
				weight     float64
				attributes map[string]string
				images     []string
				isDefault  bool
			}{
				{"JEANS-30-32", 20, 7999, 0.8, map[string]string{"Waist": "30", "Length": "32", "Color": "Blue"}, []string{}, true},
				{"JEANS-32-32", 25, 7999, 0.8, map[string]string{"Waist": "32", "Length": "32", "Color": "Blue"}, []string{}, false},
				{"JEANS-34-32", 15, 7999, 0.8, map[string]string{"Waist": "34", "Length": "32", "Color": "Blue"}, []string{}, false},
			},
		},
	}

	// Add electronics product if category exists
	if electronicsCategoryID > 0 {
		products = append(products, struct {
			name        string
			description string
			currency    string
			categoryID  uint
			images      []string
			active      bool
			variants    []struct {
				sku        string
				stock      int
				price      int64
				weight     float64
				attributes map[string]string
				images     []string
				isDefault  bool
			}
		}{
			name:        "Wireless Headphones",
			description: "Premium wireless headphones with noise cancellation",
			currency:    "USD",
			categoryID:  electronicsCategoryID,
			images:      []string{"headphones1.jpg", "headphones2.jpg"},
			active:      true,
			variants: []struct {
				sku        string
				stock      int
				price      int64
				weight     float64
				attributes map[string]string
				images     []string
				isDefault  bool
			}{
				{"WH-BLACK", 15, 19999, 0.35, map[string]string{"Color": "Black", "Type": "Wireless"}, []string{}, true},
				{"WH-WHITE", 10, 19999, 0.35, map[string]string{"Color": "White", "Type": "Wireless"}, []string{}, false},
			},
		})
	}

	// Add DKK products for MobilePay testing
	dkkProducts := []struct {
		name        string
		description string
		currency    string
		categoryID  uint
		images      []string
		active      bool
		variants    []struct {
			sku        string
			stock      int
			price      int64 // Price in øre (DKK cents)
			weight     float64
			attributes map[string]string
			images     []string
			isDefault  bool
		}
	}{
		{
			name:        "Danish Design T-Shirt",
			description: "Stylish Danish design t-shirt made from organic cotton",
			currency:    "DKK",
			categoryID:  clothingCategoryID,
			images:      []string{"danish_tshirt1.jpg", "danish_tshirt2.jpg"},
			active:      true,
			variants: []struct {
				sku        string
				stock      int
				price      int64
				weight     float64
				attributes map[string]string
				images     []string
				isDefault  bool
			}{
				{"DK-TSHIRT-M", 30, 14900, 0.2, map[string]string{"Color": "Navy", "Size": "M", "Origin": "Denmark"}, []string{}, true},    // 149 DKK
				{"DK-TSHIRT-L", 25, 14900, 0.2, map[string]string{"Color": "Navy", "Size": "L", "Origin": "Denmark"}, []string{}, false},   // 149 DKK
				{"DK-TSHIRT-XL", 20, 14900, 0.2, map[string]string{"Color": "Navy", "Size": "XL", "Origin": "Denmark"}, []string{}, false}, // 149 DKK
			},
		},
		{
			name:        "Copenhagen Hoodie",
			description: "Premium hoodie with Copenhagen city design",
			currency:    "DKK",
			categoryID:  clothingCategoryID,
			images:      []string{"cph_hoodie1.jpg", "cph_hoodie2.jpg"},
			active:      true,
			variants: []struct {
				sku        string
				stock      int
				price      int64
				weight     float64
				attributes map[string]string
				images     []string
				isDefault  bool
			}{
				{"CPH-HOODIE-M", 15, 39900, 0.6, map[string]string{"Color": "Gray", "Size": "M", "Design": "Copenhagen"}, []string{}, true},  // 399 DKK
				{"CPH-HOODIE-L", 12, 39900, 0.6, map[string]string{"Color": "Gray", "Size": "L", "Design": "Copenhagen"}, []string{}, false}, // 399 DKK
			},
		},
	}

	// Add DKK electronics if category exists
	if electronicsCategoryID > 0 {
		dkkProducts = append(dkkProducts, struct {
			name        string
			description string
			currency    string
			categoryID  uint
			images      []string
			active      bool
			variants    []struct {
				sku        string
				stock      int
				price      int64
				weight     float64
				attributes map[string]string
				images     []string
				isDefault  bool
			}
		}{
			name:        "Danish Audio Speakers",
			description: "High-quality Danish audio speakers with premium sound",
			currency:    "DKK",
			categoryID:  electronicsCategoryID,
			images:      []string{"dk_speakers1.jpg", "dk_speakers2.jpg"},
			active:      true,
			variants: []struct {
				sku        string
				stock      int
				price      int64
				weight     float64
				attributes map[string]string
				images     []string
				isDefault  bool
			}{
				{"DK-SPEAKERS-BLK", 8, 149900, 2.5, map[string]string{"Color": "Black", "Brand": "Danish Audio", "Type": "Bluetooth"}, []string{}, true},  // 1499 DKK
				{"DK-SPEAKERS-WHT", 5, 149900, 2.5, map[string]string{"Color": "White", "Brand": "Danish Audio", "Type": "Bluetooth"}, []string{}, false}, // 1499 DKK
			},
		})
	}

	// Combine USD and DKK products
	allProducts := append(products, dkkProducts...)

	for _, prodData := range allProducts {
		// Check if product already exists
		var existingProduct struct{ ID uint }
		if err := db.Table("products").Select("id").Where("name = ?", prodData.name).First(&existingProduct).Error; err == nil {
			continue // Product already exists, skip
		}

		// Create product without variants first
		product := &entity.Product{
			Name:        prodData.name,
			Description: prodData.description,
			Currency:    prodData.currency,
			CategoryID:  prodData.categoryID,
			Active:      prodData.active,
			Images:      common.StringSlice(prodData.images),
		}

		if err := db.Create(product).Error; err != nil {
			return fmt.Errorf("failed to save product %s: %w", prodData.name, err)
		}

		// Now create variants for this product
		for _, varData := range prodData.variants {
			// Check if variant already exists
			var existingVariant struct{ ID uint }
			if err := db.Table("product_variants").Select("id").Where("sku = ?", varData.sku).First(&existingVariant).Error; err == nil {
				continue // Variant already exists, skip
			}

			variant, err := entity.NewProductVariant(
				varData.sku,
				varData.stock,
				varData.price,
				varData.weight,
				varData.attributes,
				varData.images,
				varData.isDefault,
			)
			if err != nil {
				return fmt.Errorf("failed to create variant %s: %w", varData.sku, err)
			}

			variant.ProductID = product.ID

			if err := db.Create(variant).Error; err != nil {
				return fmt.Errorf("failed to save variant %s: %w", varData.sku, err)
			}
		}
	}

	return nil
}

// seedProductVariants seeds product variant data
// Note: This is typically called automatically when seeding products
// but can be used independently to add more variants to existing products
func seedProductVariants(db *gorm.DB) error {
	// Get existing products
	var products []struct {
		ID   uint
		Name string
	}
	if err := db.Table("products").Select("id, name").Find(&products).Error; err != nil {
		return fmt.Errorf("failed to fetch products: %w", err)
	}

	// Add additional variants to existing products if any
	for _, product := range products {
		if product.Name == "Classic T-Shirt" {
			// Add more color/size combinations
			additionalVariants := []struct {
				sku        string
				stock      int
				price      int64
				weight     float64
				attributes map[string]string
				images     []string
				isDefault  bool
			}{
				{"Men-W-M", 45, 1999, 0.2, map[string]string{"Color": "White", "Size": "M", "Gender": "Men"}, []string{}, false},
				{"Men-W-L", 35, 1999, 0.2, map[string]string{"Color": "White", "Size": "L", "Gender": "Men"}, []string{}, false},
				{"Women-B-M", 40, 1999, 0.18, map[string]string{"Color": "Black", "Size": "M", "Gender": "Women"}, []string{}, false},
				{"Women-B-L", 30, 1999, 0.18, map[string]string{"Color": "Black", "Size": "L", "Gender": "Women"}, []string{}, false},
			}

			for _, varData := range additionalVariants {
				// Check if variant already exists
				var existingVariant struct{ ID uint }
				if err := db.Table("product_variants").Select("id").Where("sku = ?", varData.sku).First(&existingVariant).Error; err == nil {
					continue // Variant already exists, skip
				}

				variant, err := entity.NewProductVariant(
					varData.sku,
					varData.stock,
					varData.price,
					varData.weight,
					varData.attributes,
					[]string{}, // Use empty slice for now
					varData.isDefault,
				)
				if err != nil {
					return fmt.Errorf("failed to create variant %s: %w", varData.sku, err)
				}

				variant.ProductID = product.ID

				if err := db.Create(variant).Error; err != nil {
					return fmt.Errorf("failed to save variant %s: %w", varData.sku, err)
				}
			}
		}
	}

	return nil
}

// seedOrders seeds order data
func seedOrders(db *gorm.DB) error {
	// Get users for order assignments
	var users []struct {
		ID    uint
		Email string
	}
	if err := db.Table("users").Select("id, email").Find(&users).Error; err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("No users found - skipping order seeding")
		return nil
	}

	// Get some products and variants for order items
	var variants []struct {
		ID        uint
		ProductID uint
		SKU       string
		Price     int64
		Weight    float64
	}
	if err := db.Table("product_variants").Select("id, product_id, sku, price, weight").Limit(5).Find(&variants).Error; err != nil {
		return fmt.Errorf("failed to fetch product variants: %w", err)
	}

	if len(variants) == 0 {
		fmt.Println("No product variants found - skipping order seeding")
		return nil
	}

	now := time.Now()
	orders := []struct {
		orderNumber   string
		userID        uint
		currency      string
		totalAmount   int64
		status        entity.OrderStatus
		paymentStatus entity.PaymentStatus
		items         []struct {
			productVariantID uint
			productID        uint
			sku              string
			quantity         int
			price            int64
			weight           float64
			productName      string
		}
		shippingAddr entity.Address
		billingAddr  entity.Address
		createdAt    time.Time
	}{
		{
			orderNumber:   "ORD-001",
			userID:        users[0].ID,
			currency:      "USD",
			totalAmount:   5997, // $59.97
			status:        entity.OrderStatusPaid,
			paymentStatus: entity.PaymentStatusCaptured,
			items: []struct {
				productVariantID uint
				productID        uint
				sku              string
				quantity         int
				price            int64
				weight           float64
				productName      string
			}{
				{
					productVariantID: variants[0].ID,
					productID:        variants[0].ProductID,
					sku:              variants[0].SKU,
					quantity:         3,
					price:            variants[0].Price,
					weight:           variants[0].Weight,
					productName:      "Classic T-Shirt",
				},
			},
			shippingAddr: entity.Address{
				Street1:    "123 Main St, Apt 1",
				City:       "New York",
				State:      "NY",
				PostalCode: "10001",
				Country:    "USA",
			},
			billingAddr: entity.Address{
				Street1:    "123 Main St, Apt 1",
				City:       "New York",
				State:      "NY",
				PostalCode: "10001",
				Country:    "USA",
			},
			createdAt: now.AddDate(0, 0, -7), // 7 days ago
		},
		{
			orderNumber:   "ORD-002",
			userID:        users[1].ID,
			currency:      "USD",
			totalAmount:   8999, // $89.99
			status:        entity.OrderStatusPending,
			paymentStatus: entity.PaymentStatusAuthorized,
			items: []struct {
				productVariantID uint
				productID        uint
				sku              string
				quantity         int
				price            int64
				weight           float64
				productName      string
			}{
				{
					productVariantID: variants[1].ID,
					productID:        variants[1].ProductID,
					sku:              variants[1].SKU,
					quantity:         1,
					price:            7999,
					weight:           0.8,
					productName:      "Premium Jeans",
				},
			},
			shippingAddr: entity.Address{
				Street1:    "456 Oak Ave",
				City:       "Los Angeles",
				State:      "CA",
				PostalCode: "90210",
				Country:    "USA",
			},
			billingAddr: entity.Address{
				Street1:    "456 Oak Ave",
				City:       "Los Angeles",
				State:      "CA",
				PostalCode: "90210",
				Country:    "USA",
			},
			createdAt: now.AddDate(0, 0, -3), // 3 days ago
		},
		{
			orderNumber:   "ORD-003",
			userID:        users[0].ID, // Use an existing user instead of guest
			currency:      "USD",
			totalAmount:   21999, // $219.99
			status:        entity.OrderStatusShipped,
			paymentStatus: entity.PaymentStatusCaptured,
			items: []struct {
				productVariantID uint
				productID        uint
				sku              string
				quantity         int
				price            int64
				weight           float64
				productName      string
			}{
				{
					productVariantID: variants[2].ID,
					productID:        variants[2].ProductID,
					sku:              variants[2].SKU,
					quantity:         1,
					price:            19999,
					weight:           0.35,
					productName:      "Wireless Headphones",
				},
			},
			shippingAddr: entity.Address{
				Street1:    "789 Pine St",
				City:       "Chicago",
				State:      "IL",
				PostalCode: "60601",
				Country:    "USA",
			},
			billingAddr: entity.Address{
				Street1:    "789 Pine St",
				City:       "Chicago",
				State:      "IL",
				PostalCode: "60601",
				Country:    "USA",
			},
			createdAt: now.AddDate(0, 0, -1), // 1 day ago
		},
	}

	for _, orderData := range orders {
		// Check if order already exists
		var existingOrder struct{ ID uint }
		if err := db.Table("orders").Select("id").Where("order_number = ?", orderData.orderNumber).First(&existingOrder).Error; err == nil {
			continue // Order already exists, skip
		}

		// Create order directly using entity struct (since NewOrder constructor might be complex)
		order := &entity.Order{
			OrderNumber:   orderData.orderNumber,
			Currency:      orderData.currency,
			TotalAmount:   orderData.totalAmount,
			FinalAmount:   orderData.totalAmount, // Same as total for simplicity
			Status:        orderData.status,
			PaymentStatus: orderData.paymentStatus,
			IsGuestOrder:  orderData.userID == 0,
		}

		if orderData.userID > 0 {
			order.UserID = orderData.userID
		}
		// For guest orders (userID = 0), UserID will remain 0

		// Set addresses using JSON helper methods
		order.SetShippingAddressJSON(&orderData.shippingAddr)
		order.SetBillingAddressJSON(&orderData.billingAddr)

		// Set the creation time
		order.CreatedAt = orderData.createdAt
		order.UpdatedAt = orderData.createdAt

		if err := db.Create(order).Error; err != nil {
			return fmt.Errorf("failed to save order %s: %w", orderData.orderNumber, err)
		}

		// Create order items
		for _, itemData := range orderData.items {
			subtotal := int64(itemData.quantity) * itemData.price
			orderItem := &entity.OrderItem{
				OrderID:          order.ID,
				ProductID:        itemData.productID,
				ProductVariantID: itemData.productVariantID,
				SKU:              itemData.sku,
				Quantity:         itemData.quantity,
				Price:            itemData.price,
				Subtotal:         subtotal,
				Weight:           itemData.weight,
				ProductName:      itemData.productName,
			}

			if err := db.Create(orderItem).Error; err != nil {
				return fmt.Errorf("failed to save order item for order %s: %w", orderData.orderNumber, err)
			}
		}
	}

	return nil
}

// seedDiscounts seeds discount data
func seedDiscounts(db *gorm.DB) error {
	now := time.Now()
	nextMonth := now.AddDate(0, 1, 0)
	nextYear := now.AddDate(1, 0, 0)

	discounts := []struct {
		code             string
		discountType     entity.DiscountType
		method           entity.DiscountMethod
		value            float64
		minOrderValue    int64 // in cents
		maxDiscountValue int64 // in cents
		startDate        time.Time
		endDate          time.Time
		usageLimit       int
	}{
		{
			code:             "WELCOME10",
			discountType:     entity.DiscountTypeBasket,
			method:           entity.DiscountMethodPercentage,
			value:            10.0,
			minOrderValue:    0,
			maxDiscountValue: 5000,                  // $50 max
			startDate:        now.AddDate(0, 0, -7), // Started a week ago
			endDate:          nextYear,
			usageLimit:       1000,
		},
		{
			code:             "SAVE20",
			discountType:     entity.DiscountTypeBasket,
			method:           entity.DiscountMethodFixed,
			value:            20.0,  // $20 off
			minOrderValue:    10000, // $100 minimum
			maxDiscountValue: 0,     // No max limit
			startDate:        now,
			endDate:          nextMonth,
			usageLimit:       500,
		},
		{
			code:             "SUMMER25",
			discountType:     entity.DiscountTypeBasket,
			method:           entity.DiscountMethodPercentage,
			value:            25.0,
			minOrderValue:    5000,  // $50 minimum
			maxDiscountValue: 10000, // $100 max
			startDate:        now,
			endDate:          nextMonth,
			usageLimit:       200,
		},
		{
			code:             "FREESHIP",
			discountType:     entity.DiscountTypeBasket,
			method:           entity.DiscountMethodFixed,
			value:            10.0, // Typical shipping cost
			minOrderValue:    7500, // $75 minimum
			maxDiscountValue: 0,
			startDate:        now,
			endDate:          nextYear,
			usageLimit:       0, // No limit
		},
	}

	for _, discData := range discounts {
		// Check if discount already exists
		var existingDiscount struct{ ID uint }
		if err := db.Table("discounts").Select("id").Where("code = ?", discData.code).First(&existingDiscount).Error; err == nil {
			continue // Discount already exists, skip
		}

		// Create discount using entity constructor
		discount, err := entity.NewDiscount(
			discData.code,
			discData.discountType,
			discData.method,
			discData.value,
			discData.minOrderValue,
			discData.maxDiscountValue,
			[]uint{}, // No specific products
			[]uint{}, // No specific categories
			discData.startDate,
			discData.endDate,
			discData.usageLimit,
		)
		if err != nil {
			return fmt.Errorf("failed to create discount %s: %w", discData.code, err)
		}

		if err := db.Create(discount).Error; err != nil {
			return fmt.Errorf("failed to save discount %s: %w", discData.code, err)
		}
	}

	return nil
}

// seedShippingMethods seeds shipping method data
func seedShippingMethods(db *gorm.DB) error {
	methods := []struct {
		name                  string
		description           string
		estimatedDeliveryDays int
	}{
		{"Standard Shipping", "Regular shipping with standard delivery time", 5},
		{"Express Shipping", "Fast shipping for urgent orders", 2},
		{"Next Day Delivery", "Guaranteed next business day delivery", 1},
		{"Economy Shipping", "Budget-friendly shipping option", 7},
	}

	for _, methodData := range methods {
		// Check if shipping method already exists
		var existingMethod struct{ ID uint }
		if err := db.Table("shipping_methods").Select("id").Where("name = ?", methodData.name).First(&existingMethod).Error; err == nil {
			continue // Method already exists, skip
		}

		// Create shipping method using entity constructor
		method, err := entity.NewShippingMethod(
			methodData.name,
			methodData.description,
			methodData.estimatedDeliveryDays,
		)
		if err != nil {
			return fmt.Errorf("failed to create shipping method %s: %w", methodData.name, err)
		}

		if err := db.Create(method).Error; err != nil {
			return fmt.Errorf("failed to save shipping method %s: %w", methodData.name, err)
		}
	}

	return nil
}

// seedShippingZones seeds shipping zone data
func seedShippingZones(db *gorm.DB) error {
	zones := []struct {
		name        string
		description string
		countries   []string
	}{
		{
			name:        "Domestic US",
			description: "United States domestic shipping zone",
			countries:   []string{"US", "USA"},
		},
		{
			name:        "Canada",
			description: "Canadian shipping zone",
			countries:   []string{"CA", "CAN"},
		},
		{
			name:        "Europe",
			description: "European Union and nearby countries",
			countries:   []string{"DE", "FR", "GB", "IT", "ES", "NL", "BE", "AT", "CH", "SE", "NO", "DK", "FI"},
		},
		{
			name:        "Rest of World",
			description: "All other countries worldwide",
			countries:   []string{}, // Empty means it covers all other countries
		},
	}

	for _, zoneData := range zones {
		// Check if shipping zone already exists
		var existingZone struct{ ID uint }
		if err := db.Table("shipping_zones").Select("id").Where("name = ?", zoneData.name).First(&existingZone).Error; err == nil {
			continue // Zone already exists, skip
		}

		// Create shipping zone using entity constructor
		zone, err := entity.NewShippingZone(
			zoneData.name,
			zoneData.description,
			[]string{}, // Use empty slice for now to avoid JSONB issues
		)
		if err != nil {
			return fmt.Errorf("failed to create shipping zone %s: %w", zoneData.name, err)
		}

		if err := db.Create(zone).Error; err != nil {
			return fmt.Errorf("failed to save shipping zone %s: %w", zoneData.name, err)
		}
	}

	return nil
}

// seedShippingRates seeds shipping rate data
func seedShippingRates(db *gorm.DB) error {
	// Get shipping methods and zones
	var methods []struct {
		ID   uint
		Name string
	}
	if err := db.Table("shipping_methods").Select("id, name").Find(&methods).Error; err != nil {
		return fmt.Errorf("failed to fetch shipping methods: %w", err)
	}

	var zones []struct {
		ID   uint
		Name string
	}
	if err := db.Table("shipping_zones").Select("id, name").Find(&zones).Error; err != nil {
		return fmt.Errorf("failed to fetch shipping zones: %w", err)
	}

	// Create maps for easy lookup
	methodMap := make(map[string]uint)
	for _, method := range methods {
		methodMap[method.Name] = method.ID
	}

	zoneMap := make(map[string]uint)
	for _, zone := range zones {
		zoneMap[zone.Name] = zone.ID
	}

	// Define rates for different method-zone combinations
	rates := []struct {
		methodName    string
		zoneName      string
		baseRate      int64 // in cents
		minOrderValue int64 // in cents
	}{
		// Standard Shipping rates
		{"Standard Shipping", "Domestic US", 599, 0},    // $5.99
		{"Standard Shipping", "Canada", 1299, 0},        // $12.99
		{"Standard Shipping", "Europe", 1999, 0},        // $19.99
		{"Standard Shipping", "Rest of World", 2999, 0}, // $29.99

		// Express Shipping rates
		{"Express Shipping", "Domestic US", 1299, 0},   // $12.99
		{"Express Shipping", "Canada", 2499, 0},        // $24.99
		{"Express Shipping", "Europe", 3999, 0},        // $39.99
		{"Express Shipping", "Rest of World", 5999, 0}, // $59.99

		// Next Day Delivery (only domestic)
		{"Next Day Delivery", "Domestic US", 2499, 0}, // $24.99

		// Economy Shipping rates
		{"Economy Shipping", "Domestic US", 399, 2500},     // $3.99, min $25 order
		{"Economy Shipping", "Canada", 899, 5000},          // $8.99, min $50 order
		{"Economy Shipping", "Europe", 1499, 7500},         // $14.99, min $75 order
		{"Economy Shipping", "Rest of World", 1999, 10000}, // $19.99, min $100 order
	}

	for _, rateData := range rates {
		methodID, methodExists := methodMap[rateData.methodName]
		zoneID, zoneExists := zoneMap[rateData.zoneName]

		if !methodExists {
			continue // Skip if method doesn't exist
		}
		if !zoneExists {
			continue // Skip if zone doesn't exist
		}

		// Check if rate already exists
		var existingRate struct{ ID uint }
		if err := db.Table("shipping_rates").Select("id").
			Where("shipping_method_id = ? AND shipping_zone_id = ?", methodID, zoneID).
			First(&existingRate).Error; err == nil {
			continue // Rate already exists, skip
		}

		// Create shipping rate using entity constructor
		rate, err := entity.NewShippingRate(
			methodID,
			zoneID,
			rateData.baseRate,
			rateData.minOrderValue,
		)
		if err != nil {
			return fmt.Errorf("failed to create shipping rate for %s in %s: %w", rateData.methodName, rateData.zoneName, err)
		}

		if err := db.Create(rate).Error; err != nil {
			return fmt.Errorf("failed to save shipping rate for %s in %s: %w", rateData.methodName, rateData.zoneName, err)
		}
	}

	return nil
}

// seedPaymentTransactions seeds payment transaction data for completed orders
func seedPaymentTransactions(db *gorm.DB) error {
	// Get orders that need payment transactions
	var orders []struct {
		ID            uint
		TotalAmount   int64
		Currency      string
		Status        string
		PaymentStatus string
	}
	if err := db.Table("orders").
		Select("id, total_amount, currency, status, payment_status").
		Where("payment_status IN ?", []string{"captured", "authorized"}).
		Find(&orders).Error; err != nil {
		return fmt.Errorf("failed to fetch completed orders: %w", err)
	}

	if len(orders) == 0 {
		fmt.Println("No completed orders found - skipping payment transaction seeding")
		return nil
	}

	now := time.Now()
	providers := []string{"stripe", "paypal", "square"}

	for i, order := range orders {
		provider := providers[i%len(providers)]

		// Check if payment transactions already exist for this order
		var existingCount int64
		if err := db.Table("payment_transactions").Where("order_id = ?", order.ID).Count(&existingCount).Error; err != nil {
			return fmt.Errorf("failed to check existing payment transactions for order %d: %w", order.ID, err)
		}

		if existingCount > 0 {
			continue // Payment transactions already exist for this order
		}

		// Generate unique transaction IDs
		authTxnID := fmt.Sprintf("%s_auth_%d_%d", provider, order.ID, now.Unix())
		captureTxnID := fmt.Sprintf("%s_capture_%d_%d", provider, order.ID, now.Unix()+1)

		// For SQLite compatibility, create payment transactions via SQL to avoid metadata issues
		authSQL := `INSERT INTO payment_transactions (created_at, updated_at, order_id, transaction_id, type, status, amount, currency, provider, raw_response) 
		            VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		authRawResponse := fmt.Sprintf(`{"id":"%s","amount":%d,"currency":"%s","status":"succeeded","created":%d}`,
			authTxnID, order.TotalAmount, order.Currency, now.Add(-time.Duration(i+1)*time.Hour).Unix())

		if err := db.Exec(authSQL,
			now.Add(-time.Duration(i+1)*time.Hour),
			now.Add(-time.Duration(i+1)*time.Hour),
			order.ID, authTxnID, "authorize", "successful",
			order.TotalAmount, order.Currency, provider, authRawResponse).Error; err != nil {
			return fmt.Errorf("failed to save auth transaction for order %d: %w", order.ID, err)
		}

		// Create capture transaction (usually follows authorization)
		captureSQL := `INSERT INTO payment_transactions (created_at, updated_at, order_id, transaction_id, type, status, amount, currency, provider, raw_response) 
		               VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

		captureTime := now.Add(-time.Duration(i+1) * time.Hour).Add(5 * time.Minute)
		captureRawResponse := fmt.Sprintf(`{"id":"%s","amount":%d,"currency":"%s","status":"succeeded","captured":true,"created":%d}`,
			captureTxnID, order.TotalAmount, order.Currency, captureTime.Unix())

		if err := db.Exec(captureSQL,
			captureTime, captureTime,
			order.ID, captureTxnID, "capture", "successful",
			order.TotalAmount, order.Currency, provider, captureRawResponse).Error; err != nil {
			return fmt.Errorf("failed to save capture transaction for order %d: %w", order.ID, err)
		}

		// For some orders, add a partial refund transaction (about 20% of orders)
		if i%5 == 0 && order.TotalAmount > 1000 { // Only for orders > $10.00
			refundAmount := order.TotalAmount / 4 // Refund 25%
			refundTxnID := fmt.Sprintf("%s_refund_%d_%d", provider, order.ID, now.Unix()+2)

			// Create refund transaction via SQL
			refundSQL := `INSERT INTO payment_transactions (created_at, updated_at, order_id, transaction_id, type, status, amount, currency, provider, raw_response) 
			              VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

			refundTime := captureTime.Add(24 * time.Hour)
			refundRawResponse := fmt.Sprintf(`{"id":"%s","amount":%d,"currency":"%s","status":"succeeded","refunded":true,"created":%d}`,
				refundTxnID, refundAmount, order.Currency, refundTime.Unix())

			if err := db.Exec(refundSQL,
				refundTime, refundTime,
				order.ID, refundTxnID, "refund", "successful",
				refundAmount, order.Currency, provider, refundRawResponse).Error; err != nil {
				return fmt.Errorf("failed to save refund transaction for order %d: %w", order.ID, err)
			}
		}
	}

	return nil
}

// seedCheckouts seeds checkout data for testing expiry and cleanup logic
func seedCheckouts(db *gorm.DB) error {
	// Get users for checkout assignments
	var users []struct {
		ID    uint
		Email string
	}
	if err := db.Table("users").Select("id, email").Find(&users).Error; err != nil {
		return fmt.Errorf("failed to fetch users: %w", err)
	}

	if len(users) == 0 {
		fmt.Println("No users found - skipping checkout seeding")
		return nil
	}

	// Get some product variants for checkout items
	var variants []struct {
		ID          uint
		ProductID   uint
		SKU         string
		Price       int64
		Weight      float64
		ProductName string
	}
	if err := db.Raw(`
		SELECT pv.id, pv.product_id, pv.sku, pv.price, pv.weight, p.name as product_name
		FROM product_variants pv 
		JOIN products p ON p.id = pv.product_id 
		LIMIT 5
	`).Scan(&variants).Error; err != nil {
		return fmt.Errorf("failed to fetch product variants: %w", err)
	}

	if len(variants) == 0 {
		fmt.Println("No product variants found - skipping checkout seeding")
		return nil
	}

	// Get shipping methods
	var shippingMethods []struct {
		ID   uint
		Name string
	}
	if err := db.Table("shipping_methods").Select("id, name").Find(&shippingMethods).Error; err != nil {
		return fmt.Errorf("failed to fetch shipping methods: %w", err)
	}

	// Get shipping rates
	var shippingRates []struct {
		ID               uint
		ShippingMethodID uint
		BaseRate         int64
	}
	if err := db.Table("shipping_rates").Select("id, shipping_method_id, base_rate").Find(&shippingRates).Error; err != nil {
		return fmt.Errorf("failed to fetch shipping rates: %w", err)
	}

	// Get discounts
	var discounts []struct {
		ID     uint
		Code   string
		Value  float64
		Method string
	}
	if err := db.Table("discounts").Select("id, code, value, method").Where("active = ?", true).Find(&discounts).Error; err != nil {
		return fmt.Errorf("failed to fetch discounts: %w", err)
	}

	now := time.Now()

	// Create comprehensive checkout with all features
	comprehensiveCheckout := struct {
		sessionID       string
		userID          uint
		currency        string
		status          entity.CheckoutStatus
		shippingAddress entity.Address
		billingAddress  entity.Address
		customerDetails entity.CustomerDetails
		items           []struct {
			productID   uint
			variantID   uint
			quantity    int
			price       int64
			weight      float64
			productName string
			variantName string
			sku         string
		}
		shippingOption  *entity.ShippingOption
		appliedDiscount *entity.AppliedDiscount
		createdAt       time.Time
		lastActivityAt  time.Time
		expiresAt       time.Time
	}{
		sessionID: "sess_comprehensive_001",
		userID:    users[0].ID,
		currency:  "USD",
		status:    entity.CheckoutStatusActive,
		shippingAddress: entity.Address{
			Street1:    "123 Commerce Street",
			Street2:    "Suite 456",
			City:       "Copenhagen",
			State:      "Capital Region",
			PostalCode: "2100",
			Country:    "Denmark",
		},
		billingAddress: entity.Address{
			Street1:    "789 Business Avenue",
			Street2:    "Floor 3",
			City:       "Aarhus",
			State:      "Central Denmark",
			PostalCode: "8000",
			Country:    "Denmark",
		},
		customerDetails: entity.CustomerDetails{
			Email:    "customer@example.com",
			Phone:    "+45 12 34 56 78",
			FullName: "John Doe Nielsen",
		},
		items: []struct {
			productID   uint
			variantID   uint
			quantity    int
			price       int64
			weight      float64
			productName string
			variantName string
			sku         string
		}{
			{
				productID:   variants[0].ProductID,
				variantID:   variants[0].ID,
				quantity:    2,
				price:       variants[0].Price,
				weight:      variants[0].Weight,
				productName: variants[0].ProductName,
				variantName: "Medium",
				sku:         variants[0].SKU,
			},
			{
				productID:   variants[1].ProductID,
				variantID:   variants[1].ID,
				quantity:    1,
				price:       variants[1].Price,
				weight:      variants[1].Weight,
				productName: variants[1].ProductName,
				variantName: "Large",
				sku:         variants[1].SKU,
			},
			{
				productID:   variants[2].ProductID,
				variantID:   variants[2].ID,
				quantity:    3,
				price:       variants[2].Price,
				weight:      variants[2].Weight,
				productName: variants[2].ProductName,
				variantName: "One Size",
				sku:         variants[2].SKU,
			},
		},
		createdAt:      now.Add(-2 * time.Hour),
		lastActivityAt: now.Add(-30 * time.Minute),
		expiresAt:      now.Add(22 * time.Hour),
	}

	// Add shipping option if shipping methods exist
	if len(shippingMethods) > 0 && len(shippingRates) > 0 {
		comprehensiveCheckout.shippingOption = &entity.ShippingOption{
			ShippingRateID:        shippingRates[0].ID,
			ShippingMethodID:      shippingMethods[0].ID,
			Name:                  shippingMethods[0].Name,
			Description:           "Standard shipping with tracking",
			EstimatedDeliveryDays: 3,
			Cost:                  shippingRates[0].BaseRate,
			FreeShipping:          false,
		}
	}

	// Add discount if discounts exist
	if len(discounts) > 0 {
		comprehensiveCheckout.appliedDiscount = &entity.AppliedDiscount{
			DiscountID:     discounts[0].ID,
			DiscountCode:   discounts[0].Code,
			DiscountAmount: 500, // $5.00 discount
		}
	}

	// Check if comprehensive checkout already exists
	var existingCheckout struct{ ID uint }
	if err := db.Table("checkouts").Select("id").Where("session_id = ?", comprehensiveCheckout.sessionID).First(&existingCheckout).Error; err != nil {
		// Create comprehensive checkout
		checkout, err := entity.NewCheckout(comprehensiveCheckout.sessionID, comprehensiveCheckout.currency)
		if err != nil {
			return fmt.Errorf("failed to create comprehensive checkout: %w", err)
		}

		// Set user and basic fields
		checkout.UserID = &comprehensiveCheckout.userID
		checkout.Status = comprehensiveCheckout.status
		checkout.CustomerDetails = comprehensiveCheckout.customerDetails
		checkout.CreatedAt = comprehensiveCheckout.createdAt
		checkout.UpdatedAt = comprehensiveCheckout.createdAt
		checkout.LastActivityAt = comprehensiveCheckout.lastActivityAt
		checkout.ExpiresAt = comprehensiveCheckout.expiresAt

		// Set addresses using JSON methods
		checkout.SetShippingAddressJSON(&comprehensiveCheckout.shippingAddress)
		checkout.SetBillingAddressJSON(&comprehensiveCheckout.billingAddress)

		// Set shipping option if available
		if comprehensiveCheckout.shippingOption != nil {
			checkout.SetShippingOptionJSON(comprehensiveCheckout.shippingOption)
			checkout.ShippingCost = comprehensiveCheckout.shippingOption.Cost
		}

		// Set applied discount if available
		if comprehensiveCheckout.appliedDiscount != nil {
			checkout.SetAppliedDiscountJSON(comprehensiveCheckout.appliedDiscount)
			checkout.DiscountCode = comprehensiveCheckout.appliedDiscount.DiscountCode
			checkout.DiscountAmount = comprehensiveCheckout.appliedDiscount.DiscountAmount
		}

		// Calculate totals
		var totalAmount int64
		var totalWeight float64
		for _, item := range comprehensiveCheckout.items {
			totalAmount += int64(item.quantity) * item.price
			totalWeight += float64(item.quantity) * item.weight
		}

		checkout.TotalAmount = totalAmount
		checkout.TotalWeight = totalWeight

		// Calculate final amount (total + shipping - discount)
		finalAmount := totalAmount
		if comprehensiveCheckout.shippingOption != nil {
			finalAmount += comprehensiveCheckout.shippingOption.Cost
		}
		if comprehensiveCheckout.appliedDiscount != nil {
			finalAmount -= comprehensiveCheckout.appliedDiscount.DiscountAmount
		}
		if finalAmount < 0 {
			finalAmount = 0
		}
		checkout.FinalAmount = finalAmount

		// Save checkout
		if err := db.Create(checkout).Error; err != nil {
			return fmt.Errorf("failed to save comprehensive checkout: %w", err)
		}

		// Create checkout items
		for _, itemData := range comprehensiveCheckout.items {
			checkoutItem := &entity.CheckoutItem{
				CheckoutID:       checkout.ID,
				ProductID:        itemData.productID,
				ProductVariantID: itemData.variantID,
				Quantity:         itemData.quantity,
				Price:            itemData.price,
				Weight:           itemData.weight,
				ProductName:      itemData.productName,
				VariantName:      itemData.variantName,
				SKU:              itemData.sku,
			}

			if err := db.Create(checkoutItem).Error; err != nil {
				return fmt.Errorf("failed to save checkout item for comprehensive checkout: %w", err)
			}
		}

		fmt.Printf("Created comprehensive checkout with ID %d\n", checkout.ID)
	} else {
		fmt.Println("Comprehensive checkout already exists")
	}

	// Create DKK checkout for MobilePay testing
	// Get DKK product variants
	var dkkVariants []struct {
		ID          uint
		ProductID   uint
		SKU         string
		Price       int64
		Weight      float64
		ProductName string
	}
	if err := db.Raw(`
		SELECT pv.id, pv.product_id, pv.sku, pv.price, pv.weight, p.name as product_name
		FROM product_variants pv 
		JOIN products p ON p.id = pv.product_id 
		WHERE p.currency = 'DKK'
		LIMIT 3
	`).Scan(&dkkVariants).Error; err == nil && len(dkkVariants) > 0 {

		dkkCheckout := struct {
			sessionID       string
			userID          uint
			currency        string
			status          entity.CheckoutStatus
			shippingAddress entity.Address
			billingAddress  entity.Address
			customerDetails entity.CustomerDetails
			items           []struct {
				productID   uint
				variantID   uint
				quantity    int
				price       int64
				weight      float64
				productName string
				variantName string
				sku         string
			}
			createdAt      time.Time
			lastActivityAt time.Time
			expiresAt      time.Time
		}{
			sessionID: "sess_dkk_mobilepay_001",
			userID:    users[0].ID,
			currency:  "DKK",
			status:    entity.CheckoutStatusActive,
			shippingAddress: entity.Address{
				Street1:    "Strøget 15",
				Street2:    "",
				City:       "København K",
				State:      "Capital Region",
				PostalCode: "1001",
				Country:    "Denmark",
			},
			billingAddress: entity.Address{
				Street1:    "Strøget 15",
				Street2:    "",
				City:       "København K",
				State:      "Capital Region",
				PostalCode: "1001",
				Country:    "Denmark",
			},
			customerDetails: entity.CustomerDetails{
				Email:    "mobilepay.test@example.dk",
				Phone:    "+45 12 34 56 78",
				FullName: "Lars Nielsen",
			},
			items: []struct {
				productID   uint
				variantID   uint
				quantity    int
				price       int64
				weight      float64
				productName string
				variantName string
				sku         string
			}{
				{
					productID:   dkkVariants[0].ProductID,
					variantID:   dkkVariants[0].ID,
					quantity:    1,
					price:       dkkVariants[0].Price,
					weight:      dkkVariants[0].Weight,
					productName: dkkVariants[0].ProductName,
					variantName: "M",
					sku:         dkkVariants[0].SKU,
				},
			},
			createdAt:      now.Add(-1 * time.Hour),
			lastActivityAt: now.Add(-5 * time.Minute),
			expiresAt:      now.Add(23 * time.Hour),
		}

		// Add second item if available
		if len(dkkVariants) > 1 {
			dkkCheckout.items = append(dkkCheckout.items, struct {
				productID   uint
				variantID   uint
				quantity    int
				price       int64
				weight      float64
				productName string
				variantName string
				sku         string
			}{
				productID:   dkkVariants[1].ProductID,
				variantID:   dkkVariants[1].ID,
				quantity:    1,
				price:       dkkVariants[1].Price,
				weight:      dkkVariants[1].Weight,
				productName: dkkVariants[1].ProductName,
				variantName: "M",
				sku:         dkkVariants[1].SKU,
			})
		}

		// Check if DKK checkout already exists
		var existingDKKCheckout struct{ ID uint }
		if err := db.Table("checkouts").Select("id").Where("session_id = ?", dkkCheckout.sessionID).First(&existingDKKCheckout).Error; err != nil {
			// Create DKK checkout
			checkout, err := entity.NewCheckout(dkkCheckout.sessionID, dkkCheckout.currency)
			if err != nil {
				return fmt.Errorf("failed to create DKK checkout: %w", err)
			}

			// Set user and basic fields
			checkout.UserID = &dkkCheckout.userID
			checkout.Status = dkkCheckout.status
			checkout.CustomerDetails = dkkCheckout.customerDetails
			checkout.CreatedAt = dkkCheckout.createdAt
			checkout.UpdatedAt = dkkCheckout.createdAt
			checkout.LastActivityAt = dkkCheckout.lastActivityAt
			checkout.ExpiresAt = dkkCheckout.expiresAt

			// Set addresses using JSON methods
			checkout.SetShippingAddressJSON(&dkkCheckout.shippingAddress)
			checkout.SetBillingAddressJSON(&dkkCheckout.billingAddress)

			// Calculate totals
			var totalAmount int64
			var totalWeight float64
			for _, item := range dkkCheckout.items {
				totalAmount += int64(item.quantity) * item.price
				totalWeight += float64(item.quantity) * item.weight
			}

			checkout.TotalAmount = totalAmount
			checkout.TotalWeight = totalWeight
			checkout.FinalAmount = totalAmount

			// Save checkout
			if err := db.Create(checkout).Error; err != nil {
				return fmt.Errorf("failed to save DKK checkout: %w", err)
			}

			// Create checkout items
			for _, itemData := range dkkCheckout.items {
				checkoutItem := &entity.CheckoutItem{
					CheckoutID:       checkout.ID,
					ProductID:        itemData.productID,
					ProductVariantID: itemData.variantID,
					Quantity:         itemData.quantity,
					Price:            itemData.price,
					Weight:           itemData.weight,
					ProductName:      itemData.productName,
					VariantName:      itemData.variantName,
					SKU:              itemData.sku,
				}

				if err := db.Create(checkoutItem).Error; err != nil {
					return fmt.Errorf("failed to save checkout item for DKK checkout: %w", err)
				}
			}

			fmt.Printf("Created DKK checkout for MobilePay testing with ID %d\n", checkout.ID)
		} else {
			fmt.Println("DKK checkout already exists")
		}
	}

	// Create additional simpler checkouts
	simpleCheckouts := []struct {
		sessionID string
		userID    uint
		currency  string
		status    entity.CheckoutStatus
		items     []struct {
			productID   uint
			variantID   uint
			quantity    int
			price       int64
			weight      float64
			productName string
			sku         string
		}
		createdAt      time.Time
		lastActivityAt time.Time
		expiresAt      time.Time
	}{
		{
			sessionID: "sess_simple_002",
			userID:    0, // Guest checkout
			currency:  "USD",
			status:    entity.CheckoutStatusActive,
			items: []struct {
				productID   uint
				variantID   uint
				quantity    int
				price       int64
				weight      float64
				productName string
				sku         string
			}{
				{
					productID:   variants[1].ProductID,
					variantID:   variants[1].ID,
					quantity:    1,
					price:       variants[1].Price,
					weight:      variants[1].Weight,
					productName: variants[1].ProductName,
					sku:         variants[1].SKU,
				},
			},
			createdAt:      now.Add(-1 * time.Hour),
			lastActivityAt: now.Add(-10 * time.Minute),
			expiresAt:      now.Add(23 * time.Hour),
		},
		{
			sessionID: "sess_abandoned_003",
			userID:    users[1%len(users)].ID,
			currency:  "USD",
			status:    entity.CheckoutStatusAbandoned,
			items: []struct {
				productID   uint
				variantID   uint
				quantity    int
				price       int64
				weight      float64
				productName string
				sku         string
			}{
				{
					productID:   variants[2].ProductID,
					variantID:   variants[2].ID,
					quantity:    2,
					price:       variants[2].Price,
					weight:      variants[2].Weight,
					productName: variants[2].ProductName,
					sku:         variants[2].SKU,
				},
			},
			createdAt:      now.Add(-48 * time.Hour),
			lastActivityAt: now.Add(-24 * time.Hour),
			expiresAt:      now.Add(-23 * time.Hour),
		},
	}

	// Create simple checkouts
	for _, checkoutData := range simpleCheckouts {
		// Check if checkout already exists
		var existingCheckout struct{ ID uint }
		if err := db.Table("checkouts").Select("id").Where("session_id = ?", checkoutData.sessionID).First(&existingCheckout).Error; err == nil {
			continue // Checkout already exists, skip
		}

		// Create checkout using entity constructor
		checkout, err := entity.NewCheckout(checkoutData.sessionID, checkoutData.currency)
		if err != nil {
			return fmt.Errorf("failed to create checkout %s: %w", checkoutData.sessionID, err)
		}

		// Set additional fields
		if checkoutData.userID > 0 {
			checkout.UserID = &checkoutData.userID
		}
		checkout.Status = checkoutData.status
		checkout.CreatedAt = checkoutData.createdAt
		checkout.UpdatedAt = checkoutData.createdAt
		checkout.LastActivityAt = checkoutData.lastActivityAt
		checkout.ExpiresAt = checkoutData.expiresAt

		// Calculate totals
		var totalAmount int64
		var totalWeight float64
		for _, item := range checkoutData.items {
			totalAmount += int64(item.quantity) * item.price
			totalWeight += float64(item.quantity) * item.weight
		}

		checkout.TotalAmount = totalAmount
		checkout.TotalWeight = totalWeight
		checkout.FinalAmount = totalAmount

		if err := db.Create(checkout).Error; err != nil {
			return fmt.Errorf("failed to save checkout %s: %w", checkoutData.sessionID, err)
		}

		// Create checkout items
		for _, itemData := range checkoutData.items {
			checkoutItem := &entity.CheckoutItem{
				CheckoutID:       checkout.ID,
				ProductID:        itemData.productID,
				ProductVariantID: itemData.variantID,
				Quantity:         itemData.quantity,
				Price:            itemData.price,
				Weight:           itemData.weight,
				ProductName:      itemData.productName,
				SKU:              itemData.sku,
			}

			if err := db.Create(checkoutItem).Error; err != nil {
				return fmt.Errorf("failed to save checkout item for checkout %s: %w", checkoutData.sessionID, err)
			}
		}
	}

	return nil
}
