package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/infrastructure/database"
	"golang.org/x/crypto/bcrypt"
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
	db, err := database.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

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

	// if *allFlag || *paymentTransactionsFlag {
	// 	if err := seedPaymentTransactions(db); err != nil {
	// 		log.Fatalf("Failed to seed payment transactions: %v", err)
	// 	}
	// 	fmt.Println("Payment transactions seeded successfully")
	// }

	if !*allFlag && !*usersFlag && !*categoriesFlag && !*productsFlag && !*productVariantsFlag &&
		!*ordersFlag && !*clearFlag && !*discountsFlag &&
		!*paymentTransactionsFlag && !*shippingFlag {
		fmt.Println("No action specified")
		fmt.Println("\nUsage:")
		flag.PrintDefaults()
	}
}

// clearData clears all data from the database
func clearData(db *sql.DB) error {
	// Disable foreign key checks temporarily
	if _, err := db.Exec("SET CONSTRAINTS ALL DEFERRED"); err != nil {
		return err
	}

	// Clear tables in reverse order of dependencies
	tables := []string{
		"order_items",
		"orders",
		"products",
		"categories",
		"users",
	}

	for _, table := range tables {
		if _, err := db.Exec(fmt.Sprintf("DELETE FROM %s", table)); err != nil {
			return err
		}
		// Reset sequence
		if _, err := db.Exec(fmt.Sprintf("ALTER SEQUENCE %s_id_seq RESTART WITH 1", table)); err != nil {
			return err
		}
	}

	// Re-enable foreign key checks
	if _, err := db.Exec("SET CONSTRAINTS ALL IMMEDIATE"); err != nil {
		return err
	}

	return nil
}

// seedUsers seeds user data
func seedUsers(db *sql.DB) error {
	// Hash passwords
	adminPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	userPassword, err := bcrypt.GenerateFromPassword([]byte("user123"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	now := time.Now()

	// Insert users
	users := []struct {
		email     string
		password  []byte
		firstName string
		lastName  string
		role      string
	}{
		{"admin@example.com", adminPassword, "Admin", "User", "admin"},
		{"user@example.com", userPassword, "Regular", "User", "user"},
	}

	for _, user := range users {
		_, err := db.Exec(
			`INSERT INTO users (email, password, first_name, last_name, role, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (email) DO NOTHING`,
			user.email, user.password, user.firstName, user.lastName, user.role, now, now,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedCategories seeds category data
func seedCategories(db *sql.DB) error {
	now := time.Now()

	// Insert parent categories
	parentCategories := []struct {
		name        string
		description string
	}{
		{"Electronics", "Electronic devices and accessories"},
		{"Clothing", "Apparel and fashion items"},
		{"Home & Kitchen", "Home goods and kitchen appliances"},
		{"Books", "Books and publications"},
		{"Sports & Outdoors", "Sports equipment and outdoor gear"},
	}

	for _, category := range parentCategories {
		_, err := db.Exec(
			`INSERT INTO categories (name, description, parent_id, created_at, updated_at)
			VALUES ($1, $2, NULL, $3, $4)`,
			category.name, category.description, now, now,
		)
		if err != nil {
			return err
		}
	}

	// Get parent category IDs
	rows, err := db.Query("SELECT id, name FROM categories WHERE parent_id IS NULL")
	if err != nil {
		return err
	}
	defer rows.Close()

	parentCategoryIDs := make(map[string]int)
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err
		}
		parentCategoryIDs[name] = id
	}

	// Insert subcategories
	subcategories := []struct {
		name        string
		description string
		parentName  string
	}{
		{"Smartphones", "Mobile phones and accessories", "Electronics"},
		{"Laptops", "Notebook computers", "Electronics"},
		{"Audio", "Headphones, speakers, and audio equipment", "Electronics"},
		{"Men's Clothing", "Clothing for men", "Clothing"},
		{"Women's Clothing", "Clothing for women", "Clothing"},
		{"Footwear", "Shoes and boots", "Clothing"},
		{"Kitchen Appliances", "Appliances for the kitchen", "Home & Kitchen"},
		{"Furniture", "Home furniture", "Home & Kitchen"},
		{"Fiction", "Fiction books", "Books"},
		{"Non-Fiction", "Non-fiction books", "Books"},
		{"Fitness Equipment", "Equipment for exercise and fitness", "Sports & Outdoors"},
		{"Outdoor Gear", "Gear for outdoor activities", "Sports & Outdoors"},
	}

	for _, subcategory := range subcategories {
		parentID, ok := parentCategoryIDs[subcategory.parentName]
		if !ok {
			continue
		}

		_, err := db.Exec(
			`INSERT INTO categories (name, description, parent_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5)`,
			subcategory.name, subcategory.description, parentID, now, now,
		)
		if err != nil {
			return err
		}
	}

	return nil
}

// seedProducts seeds product data
func seedProducts(db *sql.DB) error {
	// Get category IDs
	rows, err := db.Query("SELECT id, name FROM categories")
	if err != nil {
		return err
	}
	defer rows.Close()

	categoryIDs := make(map[string]int)
	for rows.Next() {
		var id int
		var name string
		if err := rows.Scan(&id, &name); err != nil {
			return err
		}
		categoryIDs[name] = id
	}

	now := time.Now()

	// Insert products
	products := []struct {
		name         string
		description  string
		price        float64
		currencyCode string
		stock        int
		categoryName string
		images       string
		active       bool
	}{
		{
			"iPhone 13",
			"Apple iPhone 13 with A15 Bionic chip",
			999.99,
			"USD",
			50,
			"Smartphones",
			`["/images/iphone13.jpg"]`,
			true,
		},
		{
			"Samsung Galaxy S21",
			"Samsung Galaxy S21 with 5G capability",
			899.99,
			"USD",
			75,
			"Smartphones",
			`["/images/galaxys21.jpg"]`,
			true,
		},
		{
			"MacBook Pro",
			"Apple MacBook Pro with M1 chip",
			1299.99,
			"USD",
			30,
			"Laptops",
			`["/images/macbookpro.jpg"]`,
			true,
		},
		{
			"Dell XPS 13",
			"Dell XPS 13 with Intel Core i7",
			1199.99,
			"USD",
			25,
			"Laptops",
			`["/images/dellxps13.jpg"]`,
			true,
		},
		{
			"Sony WH-1000XM4",
			"Sony noise-cancelling headphones",
			349.99,
			"USD",
			100,
			"Audio",
			`["/images/sonywh1000xm4.jpg"]`,
			true,
		},
		{
			"Men's Casual Shirt",
			"Comfortable casual shirt for men",
			39.99,
			"USD",
			200,
			"Men's Clothing",
			`["/images/mencasualshirt.jpg"]`,
			true,
		},
		{
			"Women's Summer Dress",
			"Lightweight summer dress for women",
			49.99,
			"USD",
			150,
			"Women's Clothing",
			`["/images/womendress.jpg"]`,
			true,
		},
		{
			"Running Shoes",
			"Comfortable shoes for running",
			89.99,
			"USD",
			120,
			"Footwear",
			`["/images/runningshoes.jpg"]`,
			true,
		},
		{
			"Coffee Maker",
			"Automatic coffee maker for home use",
			79.99,
			"USD",
			80,
			"Kitchen Appliances",
			`["/images/coffeemaker.jpg"]`,
			true,
		},
		{
			"Sofa Set",
			"3-piece sofa set for living room",
			599.99,
			"USD",
			15,
			"Furniture",
			`["/images/sofaset.jpg"]`,
			false,
		},
		{
			"The Great Gatsby",
			"Classic novel by F. Scott Fitzgerald",
			12.99,
			"USD",
			300,
			"Fiction",
			`["/images/greatgatsby.jpg"]`,
			true,
		},
		{
			"Atomic Habits",
			"Self-improvement book by James Clear",
			14.99,
			"USD",
			250,
			"Non-Fiction",
			`["/images/atomichabits.jpg"]`,
			false,
		},
		{
			"Yoga Mat",
			"Non-slip yoga mat for exercise",
			24.99,
			"USD",
			180,
			"Fitness Equipment",
			`["/images/yogamat.jpg"]`,
			false,
		},
		{
			"Camping Tent",
			"4-person camping tent for outdoor adventures",
			129.99,
			"USD",
			60,
			"Outdoor Gear",
			`["/images/campingtent.jpg"]`,
			false,
		},
	}

	for i, product := range products {
		categoryID, ok := categoryIDs[product.categoryName]
		if !ok {
			continue
		}
		// Generate product number
		productNumber := fmt.Sprintf("PROD-%06d", i+1)

		// Check if product with this product_number already exists
		var exists bool
		err := db.QueryRow(
			`SELECT EXISTS(SELECT 1 FROM products WHERE product_number = $1)`,
			productNumber,
		).Scan(&exists)

		if err != nil {
			return err
		}

		// Only insert if product doesn't exist
		if !exists {
			_, err := db.Exec(
				`INSERT INTO products (name, description, price, currency_code, stock, category_id, images, created_at, updated_at, product_number, active)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
				product.name, product.description, money.ToCents(product.price), product.currencyCode, product.stock, categoryID, product.images, now, now, productNumber, product.active,
			)
			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("Seeded products successfully\n")
	return nil
}

// seedProductVariants seeds product variant data
func seedProductVariants(db *sql.DB) error {
	// Get product IDs
	rows, err := db.Query("SELECT id, name FROM products LIMIT 8")
	if err != nil {
		return err
	}
	defer rows.Close()

	type productInfo struct {
		id   int
		name string
	}

	var products []productInfo
	for rows.Next() {
		var p productInfo
		if err := rows.Scan(&p.id, &p.name); err != nil {
			return err
		}
		products = append(products, p)
	}

	if len(products) == 0 {
		return fmt.Errorf("no products found to create variants for")
	}

	now := time.Now()

	// Sample attributes for different product types
	colorOptions := []string{"Black", "White", "Red", "Blue", "Green"}
	sizeOptions := []string{"XS", "S", "M", "L", "XL", "XXL"}
	capacityOptions := []string{"64GB", "128GB", "256GB", "512GB", "1TB"}
	materialOptions := []string{"Cotton", "Polyester", "Leather", "Wool", "Silk"}

	for _, product := range products {
		var variants []struct {
			sku        string
			price      float64
			stock      int
			attributes []map[string]string
			isDefault  bool
			productID  int
			images     string
		}

		// Create different variants based on product type
		if product.name == "iPhone 13" || product.name == "Samsung Galaxy S21" {
			// Phone variants with different colors and capacities
			for i, color := range colorOptions[:3] {
				for j, capacity := range capacityOptions[:3] {
					isDefault := (i == 0 && j == 0)
					priceAdjustment := float64(j) * 100.0 // Higher capacity costs more
					basePrice := 999.99 + priceAdjustment

					variants = append(variants, struct {
						sku        string
						price      float64
						stock      int
						attributes []map[string]string
						isDefault  bool
						productID  int
						images     string
					}{
						sku:   fmt.Sprintf("%s-%s-%s", product.name[:3], color[:1], capacity[:3]),
						price: basePrice,
						stock: 50 - (i * 10) - (j * 5),
						attributes: []map[string]string{
							{"name": "Color", "value": color},
							{"name": "Capacity", "value": capacity},
						},
						isDefault: isDefault,
						productID: product.id,
						images:    fmt.Sprintf(`["/images/%s_%s.jpg"]`, strings.ToLower(strings.ReplaceAll(product.name, " ", "")), strings.ToLower(color)),
					})
				}
			}
		} else if product.name == "Men's Casual Shirt" || product.name == "Women's Summer Dress" {
			// Clothing variants with different colors and sizes
			for i, color := range colorOptions {
				for j, size := range sizeOptions {
					// Skip some combinations to avoid too many variants
					if i > 3 || j > 4 {
						continue
					}

					isDefault := (i == 0 && j == 2) // M size in first color is default
					basePrice := 39.99

					variants = append(variants, struct {
						sku        string
						price      float64
						stock      int
						attributes []map[string]string
						isDefault  bool
						productID  int
						images     string
					}{
						sku:   fmt.Sprintf("%s-%s-%s", strings.ReplaceAll(strings.TrimSpace(product.name), "'s", ""), color[:1], size),
						price: basePrice,
						stock: 20 - (i * 2) - (j * 1),
						attributes: []map[string]string{
							{"name": "Color", "value": strings.TrimSpace(color)},
							{"name": "Size", "value": strings.TrimSpace(size)},
							{"name": "Material", "value": strings.TrimSpace(materialOptions[i%len(materialOptions)])},
						},
						isDefault: isDefault,
						productID: product.id,
						images:    fmt.Sprintf(`["/images/%s_%s.jpg"]`, strings.ToLower(strings.ReplaceAll(strings.TrimSpace(product.name), " ", "")), strings.ToLower(strings.TrimSpace(color))),
					})
				}
			}
		} else if product.name == "MacBook Pro" || product.name == "Dell XPS 13" {
			// Laptop variants with different specs
			ramOptions := []string{"8GB", "16GB", "32GB"}
			storageOptions := []string{"256GB", "512GB", "1TB"}

			for i, ram := range ramOptions {
				for j, storage := range storageOptions {
					isDefault := (i == 1 && j == 1)                        // 16GB RAM, 512GB storage is default
					priceAdjustment := float64(i)*200.0 + float64(j)*150.0 // Higher specs cost more
					basePrice := 1299.99 + priceAdjustment

					variants = append(variants, struct {
						sku        string
						price      float64
						stock      int
						attributes []map[string]string
						isDefault  bool
						productID  int
						images     string
					}{
						sku:   fmt.Sprintf("%s-%s-%s", strings.ReplaceAll(product.name, " ", "")[:3], ram[:2], storage[:3]),
						price: basePrice,
						stock: 15 - (i * 3) - (j * 2),
						attributes: []map[string]string{
							{"name": "RAM", "value": ram},
							{"name": "Storage", "value": storage},
						},
						isDefault: isDefault,
						productID: product.id,
						images:    fmt.Sprintf(`["/images/%s.jpg"]`, strings.ToLower(strings.ReplaceAll(product.name, " ", ""))),
					})
				}
			}
		}

		// Insert variants for this product
		variantsCreated := 0
		for _, variant := range variants {
			// Check if variant with this SKU already exists
			var exists bool
			err := db.QueryRow(
				`SELECT EXISTS(SELECT 1 FROM product_variants WHERE sku = $1)`,
				variant.sku,
			).Scan(&exists)

			if err != nil {
				return err
			}

			// Only insert if variant doesn't exist
			if !exists {
				attributesJSON, err := json.Marshal(variant.attributes)
				if err != nil {
					return err
				}

				// Insert product variant
				_, err = db.Exec(
					`INSERT INTO product_variants (
						sku, price, stock, attributes, is_default, product_id, currency_code,
						images, created_at, updated_at
					)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`,
					variant.sku,
					money.ToCents(variant.price),
					variant.stock,
					attributesJSON,
					variant.isDefault,
					variant.productID,
					"USD",
					variant.images,
					now,
					now,
				)
				if err != nil {
					return err
				}
				variantsCreated++
			}
		}

		// Only set has_variants=true if there are multiple variants for this product
		if len(variants) > 1 {
			_, err = db.Exec(
				`UPDATE products SET has_variants = true WHERE id = $1`,
				product.id,
			)
			if err != nil {
				return err
			}
		}

		// Notify that variants were created for this product
		fmt.Printf("Created %d variants for product: %s\n", variantsCreated, product.name)
	}

	return nil
}

// seedOrders seeds order data
func seedOrders(db *sql.DB) error {
	// Get user IDs
	rows, err := db.Query("SELECT id FROM users WHERE role = 'user' OR role = 'admin'")
	if err != nil {
		return err
	}
	defer rows.Close()

	var userIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return err
		}
		userIDs = append(userIDs, id)
	}

	if len(userIDs) == 0 {
		return fmt.Errorf("no users found to create orders for")
	}

	// Get product data
	productRows, err := db.Query("SELECT id, price FROM products")
	if err != nil {
		return err
	}
	defer productRows.Close()

	type productInfo struct {
		id    int
		price float64
	}

	var products []productInfo
	for productRows.Next() {
		var p productInfo
		if err := productRows.Scan(&p.id, &p.price); err != nil {
			return err
		}
		products = append(products, p)
	}

	if len(products) == 0 {
		return fmt.Errorf("no products found to create orders with")
	}

	// Sample addresses
	addresses := []map[string]string{
		{
			"street":      "123 Main St",
			"city":        "New York",
			"state":       "NY",
			"postal_code": "10001",
			"country":     "USA",
		},
		{
			"street":      "456 Oak Ave",
			"city":        "Los Angeles",
			"state":       "CA",
			"postal_code": "90001",
			"country":     "USA",
		},
		{
			"street":      "789 Pine Rd",
			"city":        "Chicago",
			"state":       "IL",
			"postal_code": "60601",
			"country":     "USA",
		},
		{
			"street":      "101 Maple Dr",
			"city":        "Seattle",
			"state":       "WA",
			"postal_code": "98101",
			"country":     "USA",
		},
	}

	// Order statuses
	statuses := []string{"pending", "paid", "shipped", "completed", "cancelled"}

	// Payment providers
	paymentProviders := []string{"stripe", "paypal", "mock"}

	// Create orders
	for i := 0; i < 10; i++ {
		// Select random user
		userID := userIDs[i%len(userIDs)]

		// Select random address
		addrIndex := i % len(addresses)
		shippingAddr := addresses[addrIndex]
		billingAddr := addresses[addrIndex] // Use same address for billing

		// Convert addresses to JSON
		shippingAddrJSON, err := json.Marshal(shippingAddr)
		if err != nil {
			return err
		}

		billingAddrJSON, err := json.Marshal(billingAddr)
		if err != nil {
			return err
		}

		// Select random status
		status := statuses[i%len(statuses)]

		// Create timestamps
		now := time.Now()
		createdAt := now.Add(time.Duration(-i*24) * time.Hour) // Each order created a day apart
		updatedAt := createdAt

		// Set completed_at for completed orders
		var completedAt *time.Time
		if status == "completed" {
			completedTime := updatedAt.Add(3 * 24 * time.Hour) // 3 days after creation
			completedAt = &completedTime
		}

		// Set payment details for paid, shipped, or completed orders
		var paymentID string
		var paymentProvider string
		var trackingCode string

		if status == "paid" || status == "shipped" || status == "completed" {
			paymentID = fmt.Sprintf("payment_%d_%s", i, time.Now().Format("20060102"))
			paymentProvider = paymentProviders[i%len(paymentProviders)]
		}

		if status == "shipped" || status == "completed" {
			trackingCode = fmt.Sprintf("TRACK%d%s", i, time.Now().Format("20060102"))
		}

		// Generate order number
		orderNumber := fmt.Sprintf("ORD-%s-%06d", createdAt.Format("20060102"), i+1)

		// Start a transaction
		tx, err := db.Begin()
		if err != nil {
			return err
		}

		// Insert order
		var orderID int
		// Set payment status based on order status
		var paymentStatus string
		switch status {
		case "pending":
			paymentStatus = "pending"
		case "paid", "shipped", "completed":
			paymentStatus = "captured"
		case "cancelled":
			paymentStatus = "cancelled"
		default:
			paymentStatus = "pending"
		}

		err = tx.QueryRow(`
			INSERT INTO orders (
				user_id, total_amount, status, payment_status, shipping_address, billing_address,
				payment_id, payment_provider, tracking_code, created_at, updated_at, completed_at, order_number
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
			RETURNING id
		`,
			userID,
			0, // Total amount will be updated after adding items
			status,
			paymentStatus,
			shippingAddrJSON,
			billingAddrJSON,
			paymentID,
			paymentProvider,
			trackingCode,
			createdAt,
			updatedAt,
			completedAt,
			orderNumber,
		).Scan(&orderID)

		if err != nil {
			tx.Rollback()
			return err
		}

		// Add 1-3 random products as order items
		numItems := (i % 3) + 1
		totalAmount := 0.0

		// Ensure we don't try to add more items than we have products
		if numItems > len(products) {
			numItems = len(products)
		}

		for j := 0; j < numItems; j++ {
			// Select product
			product := products[(i+j)%len(products)]

			// Random quantity between 1 and 3
			quantity := (j % 3) + 1

			// Calculate subtotal
			subtotal := float64(quantity) * product.price
			totalAmount += subtotal

			// Insert order item
			_, err = tx.Exec(`
				INSERT INTO order_items (
					order_id, product_id, quantity, price, subtotal, created_at
				)
				VALUES ($1, $2, $3, $4, $5, $6)
			`,
				orderID,
				product.id,
				quantity,
				int64(product.price),
				int64(subtotal),
				createdAt,
			)

			if err != nil {
				tx.Rollback()
				return err
			}
		}

		// Update order with total amount
		_, err = tx.Exec(`
			UPDATE orders
			SET total_amount = $1
			WHERE id = $2
		`,
			int64(totalAmount),
			orderID,
		)

		if err != nil {
			tx.Rollback()
			return err
		}

		// Commit transaction
		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

// seedDiscounts seeds discount data
func seedDiscounts(db *sql.DB) error {
	now := time.Now()
	startDate := now.Add(-24 * time.Hour)   // Start date is yesterday
	endDate := now.Add(30 * 24 * time.Hour) // End date is 30 days from now

	// Sample discounts
	discounts := []struct {
		code             string
		discountType     string
		method           string
		value            float64
		minOrderValue    float64
		maxDiscountValue float64
		productIDs       []uint
		categoryIDs      []uint
		startDate        time.Time
		endDate          time.Time
		usageLimit       int
		currentUsage     int
		active           bool
	}{
		{
			code:             "WELCOME10",
			discountType:     "basket",
			method:           "percentage",
			value:            10.0,
			minOrderValue:    0,
			maxDiscountValue: 0,
			productIDs:       []uint{},
			categoryIDs:      []uint{},
			startDate:        startDate,
			endDate:          endDate,
			usageLimit:       0,
			currentUsage:     0,
			active:           true,
		},
		{
			code:             "SAVE20",
			discountType:     "basket",
			method:           "percentage",
			value:            20.0,
			minOrderValue:    100.0,
			maxDiscountValue: 50.0,
			productIDs:       []uint{},
			categoryIDs:      []uint{},
			startDate:        startDate,
			endDate:          endDate,
			usageLimit:       100,
			currentUsage:     0,
			active:           true,
		},
		{
			code:             "FLAT25",
			discountType:     "basket",
			method:           "fixed",
			value:            25.0,
			minOrderValue:    150.0,
			maxDiscountValue: 0,
			productIDs:       []uint{},
			categoryIDs:      []uint{},
			startDate:        startDate,
			endDate:          endDate,
			usageLimit:       50,
			currentUsage:     0,
			active:           true,
		},
	}

	// Get product IDs for product-specific discounts
	productRows, err := db.Query("SELECT id FROM products LIMIT 5")
	if err != nil {
		return err
	}
	defer productRows.Close()

	var productIDs []uint
	for productRows.Next() {
		var id uint
		if err := productRows.Scan(&id); err != nil {
			return err
		}
		productIDs = append(productIDs, id)
	}

	// Get category IDs for category-specific discounts
	categoryRows, err := db.Query("SELECT id FROM categories WHERE parent_id IS NOT NULL LIMIT 3")
	if err != nil {
		return err
	}
	defer categoryRows.Close()

	var categoryIDs []uint
	for categoryRows.Next() {
		var id uint
		if err := categoryRows.Scan(&id); err != nil {
			return err
		}
		categoryIDs = append(categoryIDs, id)
	}

	// Add product-specific discounts if we have products
	if len(productIDs) > 0 {
		// Product-specific percentage discount
		productDiscount := struct {
			code             string
			discountType     string
			method           string
			value            float64
			minOrderValue    float64
			maxDiscountValue float64
			productIDs       []uint
			categoryIDs      []uint
			startDate        time.Time
			endDate          time.Time
			usageLimit       int
			currentUsage     int
			active           bool
		}{
			code:             "PRODUCT15",
			discountType:     "product",
			method:           "percentage",
			value:            15.0,
			minOrderValue:    0,
			maxDiscountValue: 0,
			productIDs:       productIDs[:2], // Use first 2 products
			categoryIDs:      []uint{},
			startDate:        startDate,
			endDate:          endDate,
			usageLimit:       0,
			currentUsage:     0,
			active:           true,
		}
		discounts = append(discounts, productDiscount)

		// Product-specific fixed discount
		productFixedDiscount := struct {
			code             string
			discountType     string
			method           string
			value            float64
			minOrderValue    float64
			maxDiscountValue float64
			productIDs       []uint
			categoryIDs      []uint
			startDate        time.Time
			endDate          time.Time
			usageLimit       int
			currentUsage     int
			active           bool
		}{
			code:             "PRODUCT10OFF",
			discountType:     "product",
			method:           "fixed",
			value:            100.0,
			minOrderValue:    0,
			maxDiscountValue: 0,
			productIDs:       productIDs[2:], // Use remaining products
			categoryIDs:      []uint{},
			startDate:        startDate,
			endDate:          endDate,
			usageLimit:       0,
			currentUsage:     0,
			active:           true,
		}
		discounts = append(discounts, productFixedDiscount)
	}

	// Add category-specific discounts if we have categories
	if len(categoryIDs) > 0 {
		categoryDiscount := struct {
			code             string
			discountType     string
			method           string
			value            float64
			minOrderValue    float64
			maxDiscountValue float64
			productIDs       []uint
			categoryIDs      []uint
			startDate        time.Time
			endDate          time.Time
			usageLimit       int
			currentUsage     int
			active           bool
		}{
			code:             "CATEGORY25",
			discountType:     "product",
			method:           "percentage",
			value:            25.0,
			minOrderValue:    0,
			maxDiscountValue: 0,
			productIDs:       []uint{},
			categoryIDs:      categoryIDs,
			startDate:        startDate,
			endDate:          endDate,
			usageLimit:       0,
			currentUsage:     0,
			active:           true,
		}
		discounts = append(discounts, categoryDiscount)
	}

	// Insert discounts
	for _, discount := range discounts {
		productIDsJSON, err := json.Marshal(discount.productIDs)
		if err != nil {
			return err
		}

		categoryIDsJSON, err := json.Marshal(discount.categoryIDs)
		if err != nil {
			return err
		}

		_, err = db.Exec(
			`INSERT INTO discounts (
				code, type, method, value, min_order_value, max_discount_value,
				product_ids, category_ids, start_date, end_date,
				usage_limit, current_usage, active, created_at, updated_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			ON CONFLICT (code) DO NOTHING`,
			discount.code,
			discount.discountType,
			discount.method,
			discount.value,
			money.ToCents(discount.minOrderValue),
			money.ToCents(discount.maxDiscountValue),
			productIDsJSON,
			categoryIDsJSON,
			discount.startDate,
			discount.endDate,
			discount.usageLimit,
			discount.currentUsage,
			discount.active,
			now,
			now,
		)
		if err != nil {
			return err
		}
	}

	fmt.Printf("Seeded %d discounts\n", len(discounts))
	return nil
}

// seedShippingMethods seeds shipping method data
func seedShippingMethods(db *sql.DB) error {
	now := time.Now()

	// Insert shipping methods
	methods := []struct {
		name                  string
		description           string
		active                bool
		estimatedDeliveryDays int
	}{
		{
			name:                  "Standard Shipping",
			description:           "Standard delivery - 3-5 business days",
			active:                true,
			estimatedDeliveryDays: 4, // average of 3-5 days
		},
		{
			name:                  "Express Shipping",
			description:           "Express delivery - 1-2 business days",
			active:                true,
			estimatedDeliveryDays: 1, // minimum delivery time
		},
		{
			name:                  "Next Day Delivery",
			description:           "Next business day delivery (order by 2pm)",
			active:                true,
			estimatedDeliveryDays: 1,
		},
		{
			name:                  "Economy Shipping",
			description:           "Budget-friendly shipping - 5-8 business days",
			active:                true,
			estimatedDeliveryDays: 7, // average of 5-8 days
		},
		{
			name:                  "International Shipping",
			description:           "International delivery - 7-14 business days",
			active:                true,
			estimatedDeliveryDays: 10, // average of 7-14 days
		},
	}

	for _, method := range methods {
		// Check if the shipping method already exists
		var exists bool
		err := db.QueryRow(
			`SELECT EXISTS(SELECT 1 FROM shipping_methods WHERE name = $1)`,
			method.name,
		).Scan(&exists)

		if err != nil {
			return err
		}

		// Only insert if the shipping method doesn't exist
		if !exists {
			_, err := db.Exec(
				`INSERT INTO shipping_methods (
					name, description, active, estimated_delivery_days, created_at, updated_at
				)
				VALUES ($1, $2, $3, $4, $5, $6)`,
				method.name,
				method.description,
				method.active,
				method.estimatedDeliveryDays,
				now,
				now,
			)
			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("Seeded %d shipping methods\n", len(methods))
	return nil
}

// seedShippingZones seeds shipping zone data
func seedShippingZones(db *sql.DB) error {
	now := time.Now()

	// Insert shipping zones
	zones := []struct {
		name        string
		description string
		countries   []string
		active      bool
	}{
		{
			name:        "Domestic",
			description: "Shipping within the United States",
			countries:   []string{"USA"},
			active:      true,
		},
		{
			name:        "North America",
			description: "Shipping to North American countries",
			countries:   []string{"USA", "CAN", "MEX"},
			active:      true,
		},
		{
			name:        "Europe",
			description: "Shipping to European countries",
			countries:   []string{"GBR", "DEU", "FRA", "ESP", "ITA", "NLD", "SWE", "NOR", "DNK", "FIN"},
			active:      true,
		},
		{
			name:        "Asia Pacific",
			description: "Shipping to Asia-Pacific countries",
			countries:   []string{"JPN", "CHN", "KOR", "AUS", "NZL", "SGP", "THA", "IDN"},
			active:      true,
		},
		{
			name:        "Rest of World",
			description: "Shipping to all other countries",
			countries:   []string{"*"},
			active:      true,
		},
	}

	for _, zone := range zones {
		// Check if the shipping zone already exists
		var exists bool
		err := db.QueryRow(
			`SELECT EXISTS(SELECT 1 FROM shipping_zones WHERE name = $1)`,
			zone.name,
		).Scan(&exists)

		if err != nil {
			return err
		}

		// Only insert if the shipping zone doesn't exist
		if !exists {
			countriesJSON, err := json.Marshal(zone.countries)
			if err != nil {
				return err
			}

			_, err = db.Exec(
				`INSERT INTO shipping_zones (
					name, description, countries, active, created_at, updated_at
				)
				VALUES ($1, $2, $3, $4, $5, $6)`,
				zone.name,
				zone.description,
				countriesJSON,
				zone.active,
				now,
				now,
			)
			if err != nil {
				return err
			}
		}
	}

	fmt.Printf("Seeded %d shipping zones\n", len(zones))
	return nil
}

// seedShippingRates seeds shipping rate data
func seedShippingRates(db *sql.DB) error {
	// Get shipping method IDs
	methodRows, err := db.Query("SELECT id, name FROM shipping_methods")
	if err != nil {
		return err
	}
	defer methodRows.Close()

	methodIDs := make(map[string]int)
	for methodRows.Next() {
		var id int
		var name string
		if err := methodRows.Scan(&id, &name); err != nil {
			return err
		}
		methodIDs[name] = id
	}

	// Get shipping zone IDs
	zoneRows, err := db.Query("SELECT id, name FROM shipping_zones")
	if err != nil {
		return err
	}
	defer zoneRows.Close()

	zoneIDs := make(map[string]int)
	for zoneRows.Next() {
		var id int
		var name string
		if err := zoneRows.Scan(&id, &name); err != nil {
			return err
		}
		zoneIDs[name] = id
	}

	now := time.Now()

	// Insert base shipping rates
	baseRates := []struct {
		displayName           string // For logging only, not stored in DB
		methodName            string
		zoneName              string
		baseRate              float64
		minOrderValue         float64
		freeShippingThreshold *float64
		active                bool
		rateType              string
	}{
		{
			displayName:           "Domestic Standard",
			methodName:            "Standard Shipping",
			zoneName:              "Domestic",
			baseRate:              5.99,
			minOrderValue:         0,
			freeShippingThreshold: nil,
			active:                true,
			rateType:              "flat",
		},
		{
			displayName:           "Domestic Express",
			methodName:            "Express Shipping",
			zoneName:              "Domestic",
			baseRate:              12.99,
			minOrderValue:         0,
			freeShippingThreshold: &[]float64{75.0}[0], // Free shipping over $75
			active:                true,
			rateType:              "flat",
		},
		{
			displayName:           "North America Standard",
			methodName:            "Standard Shipping",
			zoneName:              "North America",
			baseRate:              15.99,
			minOrderValue:         0,
			freeShippingThreshold: &[]float64{100.0}[0], // Free shipping over $100
			active:                true,
			rateType:              "flat",
		},
		{
			displayName:           "Europe Standard",
			methodName:            "Standard Shipping",
			zoneName:              "Europe",
			baseRate:              24.99,
			minOrderValue:         0,
			freeShippingThreshold: nil,
			active:                true,
			rateType:              "weight_based",
		},
		{
			displayName:           "Europe Express",
			methodName:            "Express Shipping",
			zoneName:              "Europe",
			baseRate:              34.99,
			minOrderValue:         0,
			freeShippingThreshold: nil,
			active:                true,
			rateType:              "weight_based",
		},
		{
			displayName:           "Asia Pacific Standard",
			methodName:            "Standard Shipping",
			zoneName:              "Asia Pacific",
			baseRate:              29.99,
			minOrderValue:         0,
			freeShippingThreshold: nil,
			active:                true,
			rateType:              "value_based",
		},
		{
			displayName:           "Worldwide Economy",
			methodName:            "Economy Shipping",
			zoneName:              "Rest of World",
			baseRate:              39.99,
			minOrderValue:         0,
			freeShippingThreshold: nil,
			active:                true,
			rateType:              "value_based",
		},
	}

	// Start a transaction for inserting rates
	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, rate := range baseRates {
		methodID, ok := methodIDs[rate.methodName]
		if !ok {
			tx.Rollback()
			return fmt.Errorf("shipping method not found: %s", rate.methodName)
		}

		zoneID, ok := zoneIDs[rate.zoneName]
		if !ok {
			tx.Rollback()
			return fmt.Errorf("shipping zone not found: %s", rate.zoneName)
		}

		// Insert basic shipping rate
		var rateID int
		var freeShippingThresholdCents *int64
		if rate.freeShippingThreshold != nil {
			thresholdCents := money.ToCents(*rate.freeShippingThreshold)
			freeShippingThresholdCents = &thresholdCents
		}

		err := tx.QueryRow(
			`INSERT INTO shipping_rates (
				shipping_method_id, shipping_zone_id, base_rate, min_order_value, 
				free_shipping_threshold, active, created_at, updated_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
			RETURNING id`,
			methodID,
			zoneID,
			money.ToCents(rate.baseRate),
			money.ToCents(rate.minOrderValue),
			freeShippingThresholdCents,
			rate.active,
			now,
			now,
		).Scan(&rateID)

		if err != nil {
			tx.Rollback()
			return err
		}

		// Add weight-based rules for weight-based rates
		if rate.rateType == "weight_based" {
			weightRules := []struct {
				minWeight float64
				maxWeight float64
				rate      float64
			}{
				{0.0, 1.0, rate.baseRate},
				{1.01, 2.0, rate.baseRate * 1.5},
				{2.01, 5.0, rate.baseRate * 2.0},
				{5.01, 10.0, rate.baseRate * 3.0},
				{10.01, 20.0, rate.baseRate * 4.0},
			}

			for _, rule := range weightRules {
				_, err := tx.Exec(
					`INSERT INTO weight_based_rates (
						shipping_rate_id, min_weight, max_weight, rate
					)
					VALUES ($1, $2, $3, $4)`,
					rateID,
					rule.minWeight,
					rule.maxWeight,
					money.ToCents(rule.rate),
				)

				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		// Add value-based rules for value-based rates
		if rate.rateType == "value_based" {
			valueRules := []struct {
				minValue float64
				maxValue float64
				rate     float64
			}{
				{0.0, 50.0, rate.baseRate},
				{50.01, 100.0, rate.baseRate * 1.25},
				{100.01, 250.0, rate.baseRate * 1.5},
				{250.01, 500.0, rate.baseRate * 1.75},
				{500.01, 1000.0, rate.baseRate * 2.0},
				{1000.01, 9999999.0, rate.baseRate * 2.5},
			}

			for _, rule := range valueRules {
				_, err := tx.Exec(
					`INSERT INTO value_based_rates (
						shipping_rate_id, min_order_value, max_order_value, rate
					)
					VALUES ($1, $2, $3, $4)`,
					rateID,
					money.ToCents(rule.minValue),
					money.ToCents(rule.maxValue),
					money.ToCents(rule.rate),
				)

				if err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		fmt.Printf("Created shipping rate: %s (%s to %s)\n", rate.displayName, rate.methodName, rate.zoneName)
	}

	// Commit the transaction
	if err := tx.Commit(); err != nil {
		return err
	}

	fmt.Printf("Seeded %d shipping rates with associated rules\n", len(baseRates))
	return nil
}

// seedPaymentTransactions seeds payment transaction data
func seedPaymentTransactions(db *sql.DB) error {
	// Get order IDs with payment providers set
	orderRows, err := db.Query(`
		SELECT id, payment_id, payment_provider, total_amount, order_number 
		FROM orders 
		WHERE payment_provider IS NOT NULL 
		AND status IN ('paid', 'shipped', 'completed')
	`)
	if err != nil {
		return err
	}
	defer orderRows.Close()

	type orderInfo struct {
		id              int
		paymentID       string
		paymentProvider string
		totalAmount     int64
		orderNumber     string
	}

	var orders []orderInfo
	for orderRows.Next() {
		var o orderInfo
		if err := orderRows.Scan(&o.id, &o.paymentID, &o.paymentProvider, &o.totalAmount, &o.orderNumber); err != nil {
			return err
		}
		orders = append(orders, o)
	}

	if len(orders) == 0 {
		return fmt.Errorf("no paid orders found to create payment transactions for")
	}

	now := time.Now()

	// Transaction statuses by provider
	statuses := map[string][]string{
		"stripe":    {"successful", "pending", "failed"},
		"paypal":    {"successful", "pending", "failed"},
		"mobilepay": {"successful", "pending", "failed"},
		"mock":      {"successful", "pending", "failed"},
	}

	// Transaction types
	transactionTypes := []string{"authorize", "capture", "refund"}

	// Create payment transactions
	for i, order := range orders {
		// Set transaction status (mostly successful, with a few failures for testing)
		statusList := statuses[order.paymentProvider]
		if statusList == nil {
			statusList = statuses["mock"] // Fallback to mock statuses
		}

		var status string
		if i < len(orders)-2 {
			status = statusList[0] // Success status (first in each list)
		} else {
			status = statusList[i%len(statusList)] // Mix of statuses for the last few
		}

		// Determine transaction type based on index
		transactionType := transactionTypes[i%len(transactionTypes)]

		// Generate metadata
		metadata := map[string]interface{}{
			"order_number": order.orderNumber,
			"customer_ip":  fmt.Sprintf("192.168.1.%d", 100+i%100),
			"user_agent":   "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36",
		}

		metadataJSON, err := json.Marshal(metadata)
		if err != nil {
			return err
		}

		// Insert payment transaction using the correct column names from the schema
		_, err = db.Exec(`
			INSERT INTO payment_transactions (
				order_id, transaction_id, type, status, amount, currency, provider,
				metadata, created_at, updated_at
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		`,
			order.id,
			order.paymentID,
			transactionType,
			status,
			order.totalAmount,
			"USD", // Default currency
			order.paymentProvider,
			metadataJSON,
			now,
			now,
		)

		if err != nil {
			return err
		}
	}

	fmt.Printf("Seeded %d payment transactions\n", len(orders))
	return nil
}
