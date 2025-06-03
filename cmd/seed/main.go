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
	"github.com/lib/pq"
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

	// Temporarily disable the variant deletion trigger
	if _, err := db.Exec("DROP TRIGGER IF EXISTS prevent_last_variant_deletion ON product_variants"); err != nil {
		return err
	}

	// Clear tables in reverse order of dependencies
	tables := []string{
		"checkout_items",
		"checkouts",
		"payment_transactions",
		"order_items",
		"orders",
		"product_variants",
		"products",
		"categories",
		"users",
		"discounts",
		"webhooks",
		"shipping_methods",
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

	// Re-create the variant deletion trigger
	triggerSQL := `
	CREATE TRIGGER prevent_last_variant_deletion
		BEFORE DELETE ON product_variants
		FOR EACH ROW
		EXECUTE FUNCTION check_product_has_variants();
	`
	if _, err := db.Exec(triggerSQL); err != nil {
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

	for _, product := range products {
		categoryID, ok := categoryIDs[product.categoryName]
		if !ok {
			continue
		}

		// Check if product with this name already exists
		var exists bool
		err := db.QueryRow(
			`SELECT EXISTS(SELECT 1 FROM products WHERE name = $1)`,
			product.name,
		).Scan(&exists)

		if err != nil {
			return err
		}

		// Only insert if product doesn't exist
		if !exists {
			_, err := db.Exec(
				`INSERT INTO products (name, description, price, stock, category_id, images, created_at, updated_at, active)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
				product.name, product.description, money.ToCents(product.price), product.stock, categoryID, product.images, now, now, product.active,
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
				variantName := product.name
				if len(variant.attributes) > 0 {
					var values []string
					for _, attr := range variant.attributes {
						if value, ok := attr["value"]; ok {
							values = append(values, value)
						}
					}
					if len(values) > 0 {
						variantName = fmt.Sprintf("%s - %s", product.name, strings.Join(values, ", "))
					}
				}

				_, err = db.Exec(
					`INSERT INTO product_variants (
						product_id, name, sku, price, stock, weight, attributes, created_at, updated_at
					)
					VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
					variant.productID,
					variantName,
					variant.sku,
					money.ToCents(variant.price),
					variant.stock,
					0.5, // Default weight in kg
					attributesJSON,
					now,
					now,
				)
				if err != nil {
					return err
				}
			}
		}

		// Notify that variants were created for this product
		fmt.Printf("Created %d variants for product: %s\n", len(variants), product.name)
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
		// Convert value based on method
		var valueInCents int64
		if discount.method == "percentage" {
			// Store percentage as value * 100 (e.g., 10% = 1000)
			valueInCents = int64(discount.value * 100)
		} else {
			// Store fixed amount as cents
			valueInCents = money.ToCents(discount.value)
		}

		// Convert uint slices to int slices for PostgreSQL
		productIDsInt := make([]int, len(discount.productIDs))
		for i, id := range discount.productIDs {
			productIDsInt[i] = int(id)
		}

		categoryIDsInt := make([]int, len(discount.categoryIDs))
		for i, id := range discount.categoryIDs {
			categoryIDsInt[i] = int(id)
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
			valueInCents,
			money.ToCents(discount.minOrderValue),
			money.ToCents(discount.maxDiscountValue),
			pq.Array(productIDsInt),
			pq.Array(categoryIDsInt),
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
		name            string
		description     string
		baseRate        float64
		ratePerKg       float64
		minDeliveryDays int
		maxDeliveryDays int
		active          bool
	}{
		{
			name:            "Standard Shipping",
			description:     "Standard delivery - 3-5 business days",
			baseRate:        5.99,
			ratePerKg:       2.00,
			minDeliveryDays: 3,
			maxDeliveryDays: 5,
			active:          true,
		},
		{
			name:            "Express Shipping",
			description:     "Express delivery - 1-2 business days",
			baseRate:        12.99,
			ratePerKg:       5.00,
			minDeliveryDays: 1,
			maxDeliveryDays: 2,
			active:          true,
		},
		{
			name:            "Next Day Delivery",
			description:     "Next business day delivery (order by 2pm)",
			baseRate:        19.99,
			ratePerKg:       8.00,
			minDeliveryDays: 1,
			maxDeliveryDays: 1,
			active:          true,
		},
		{
			name:            "Economy Shipping",
			description:     "Budget-friendly shipping - 5-8 business days",
			baseRate:        3.99,
			ratePerKg:       1.50,
			minDeliveryDays: 5,
			maxDeliveryDays: 8,
			active:          true,
		},
		{
			name:            "International Shipping",
			description:     "International delivery - 7-14 business days",
			baseRate:        25.99,
			ratePerKg:       10.00,
			minDeliveryDays: 7,
			maxDeliveryDays: 14,
			active:          true,
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
					name, description, base_rate, rate_per_kg, min_delivery_days, max_delivery_days, active, created_at, updated_at
				)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`,
				method.name,
				method.description,
				money.ToCents(method.baseRate),
				money.ToCents(method.ratePerKg),
				method.minDeliveryDays,
				method.maxDeliveryDays,
				method.active,
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
