package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// ProductRepository is the PostgreSQL implementation of the ProductRepository interface
type ProductRepository struct {
	db                *sql.DB
	variantRepository repository.ProductVariantRepository
}

// NewProductRepository creates a new ProductRepository
func NewProductRepository(db *sql.DB, variantRepository repository.ProductVariantRepository) repository.ProductRepository {
	return &ProductRepository{
		db:                db,
		variantRepository: variantRepository,
	}
}

// Create creates a new product
func (r *ProductRepository) Create(product *entity.Product) error {
	query := `
	INSERT INTO products (name, description, price, currency_code, stock, weight, category_id, images, has_variants, active, created_at, updated_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	RETURNING id
	`

	imagesJSON, err := json.Marshal(product.Images)
	if err != nil {
		return err
	}

	err = r.db.QueryRow(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.CurrencyCode,
		product.Stock,
		product.Weight,
		product.CategoryID,
		imagesJSON,
		product.HasVariants,
		product.Active,
		product.CreatedAt,
		product.UpdatedAt,
	).Scan(&product.ID)
	if err != nil {
		return err
	}

	// Generate and set the product number
	product.SetProductNumber(product.ID)

	// Update the product number in the database
	updateQuery := "UPDATE products SET product_number = $1 WHERE id = $2"
	_, err = r.db.Exec(updateQuery, product.ProductNumber, product.ID)
	if err != nil {
		return err
	}

	// If the product has currency-specific prices, save them
	if len(product.Prices) > 0 {
		for i := range product.Prices {
			product.Prices[i].ProductID = product.ID
			if err = r.createProductPrice(&product.Prices[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

// createProductPrice creates a product price entry for a specific currency
func (r *ProductRepository) createProductPrice(price *entity.ProductPrice) error {
	query := `
		INSERT INTO product_prices (product_id, currency_code, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (product_id, currency_code) DO UPDATE SET
		price = EXCLUDED.price,
		updated_at = EXCLUDED.updated_at
		RETURNING id
		`

	now := time.Now()

	return r.db.QueryRow(
		query,
		price.ProductID,
		price.CurrencyCode,
		price.Price,
		now,
		now,
	).Scan(&price.ID)
}

// GetByID gets a product by ID
func (r *ProductRepository) GetByID(productID uint) (*entity.Product, error) {
	query := `
			SELECT id, product_number, name, description, price, currency_code, stock, weight, category_id, images, has_variants, active, created_at, updated_at
			FROM products
			WHERE id = $1
			`

	var imagesJSON []byte
	product := &entity.Product{}
	var productNumber sql.NullString

	err := r.db.QueryRow(query, productID).Scan(
		&product.ID,
		&productNumber,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CurrencyCode,
		&product.Stock,
		&product.Weight,
		&product.CategoryID,
		&imagesJSON,
		&product.HasVariants,
		&product.Active,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Set product number if valid
	if productNumber.Valid {
		product.ProductNumber = productNumber.String
	}

	// Unmarshal images JSON
	if err := json.Unmarshal(imagesJSON, &product.Images); err != nil {
		return nil, err
	}

	// Load currency-specific prices
	prices, err := r.getProductPrices(product.ID)
	if err != nil {
		return nil, err
	}
	product.Prices = prices

	return product, nil
}

// GetByProductNumber gets a product by product number
func (r *ProductRepository) GetByProductNumber(productNumber string) (*entity.Product, error) {
	query := `
			SELECT id, product_number, name, description, price, currency_code, stock, weight, category_id, images, has_variants, active, created_at, updated_at
			FROM products
			WHERE product_number = $1
			`

	var imagesJSON []byte
	product := &entity.Product{}
	var productNumberResult sql.NullString

	err := r.db.QueryRow(query, productNumber).Scan(
		&product.ID,
		&productNumberResult,
		&product.Name,
		&product.Description,
		&product.Price,
		&product.CurrencyCode,
		&product.Stock,
		&product.Weight,
		&product.CategoryID,
		&imagesJSON,
		&product.HasVariants,
		&product.Active,
		&product.CreatedAt,
		&product.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	// Set product number if valid
	if productNumberResult.Valid {
		product.ProductNumber = productNumberResult.String
	}

	// Unmarshal images JSON
	if err := json.Unmarshal(imagesJSON, &product.Images); err != nil {
		return nil, err
	}

	// Load currency-specific prices
	prices, err := r.getProductPrices(product.ID)
	if err != nil {
		return nil, err
	}
	product.Prices = prices

	return product, nil
}

// getProductPrices retrieves all prices for a product in different currencies
func (r *ProductRepository) getProductPrices(productID uint) ([]entity.ProductPrice, error) {
	query := `
			SELECT id, product_id, currency_code, price, created_at, updated_at
			FROM product_prices
			WHERE product_id = $1
			`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []entity.ProductPrice
	for rows.Next() {
		var price entity.ProductPrice

		err := rows.Scan(
			&price.ID,
			&price.ProductID,
			&price.CurrencyCode,
			&price.Price,
			&price.CreatedAt,
			&price.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		prices = append(prices, price)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return prices, nil
}

// GetByIDWithVariants gets a product by ID with variants
func (r *ProductRepository) GetByIDWithVariants(productID uint) (*entity.Product, error) {
	// Get the base product
	product, err := r.GetByID(productID)
	if err != nil {
		return nil, err
	}

	// If product has variants, get them
	if product.HasVariants {
		variants, err := r.variantRepository.GetByProduct(productID)
		if err != nil {
			return nil, err
		}

		product.Variants = variants
	}

	return product, nil
}

// Update updates a product
func (r *ProductRepository) Update(product *entity.Product) error {
	query := `
			UPDATE products
			SET name = $1, description = $2, price = $3, currency_code = $4, stock = $5, weight = $6, category_id = $7, 
		    images = $8, has_variants = $9, updated_at = $10
			WHERE id = $11
			`

	imagesJSON, err := json.Marshal(product.Images)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		query,
		product.Name,
		product.Description,
		product.Price,
		product.CurrencyCode,
		product.Stock,
		product.Weight,
		product.CategoryID,
		imagesJSON,
		product.HasVariants,
		time.Now(),
		product.ID,
	)
	if err != nil {
		return err
	}

	// Update currency-specific prices if they exist
	if len(product.Prices) > 0 {
		// Use an upsert query to update or insert prices
		query := `
			INSERT INTO product_prices (product_id, currency_code, price)
			VALUES ($1, $2, $3)
			ON CONFLICT (product_id, currency_code)
			DO UPDATE SET price = EXCLUDED.price
		`
		for _, price := range product.Prices {
			_, err := r.db.Exec(query, product.ID, price.CurrencyCode, price.Price)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// Delete deletes a product and all its related data
// This operation cascades to delete variants, variant prices, and product prices
func (r *ProductRepository) Delete(productID uint) error {
	if productID == 0 {
		return fmt.Errorf("invalid product ID: %d", productID)
	}

	result, err := r.db.Exec("DELETE FROM products WHERE id = $1", productID)
	if err != nil {
		return fmt.Errorf("failed to delete product: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check deletion result: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("product with ID %d not found", productID)
	}

	return nil
}

// List lists products with pagination
func (r *ProductRepository) List(query, currency string, categoryID, offset, limit uint, minPriceCents, maxPriceCents int64, active bool) ([]*entity.Product, error) {
	// Build dynamic query parts
	searchQuery := `
		SELECT 
			p.id, p.product_number, p.name, p.description, 
			COALESCE(pv.price, p.price) as price,
			p.currency_code, p.stock, p.weight, p.category_id, p.images, p.has_variants, p.active, p.created_at, p.updated_at
		FROM products p
		LEFT JOIN product_variants pv ON p.id = pv.product_id AND pv.is_default = true
	`
	queryParams := []interface{}{}
	paramCounter := 1

	var whereAdded bool
	if active {
		searchQuery += " WHERE p.active = true"
		whereAdded = true
	}

	if query != "" {
		if whereAdded {
			searchQuery += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.description ILIKE $%d)", paramCounter, paramCounter)
		} else {
			searchQuery += fmt.Sprintf(" WHERE (p.name ILIKE $%d OR p.description ILIKE $%d)", paramCounter, paramCounter)
			whereAdded = true
		}
		queryParams = append(queryParams, "%"+query+"%")
		paramCounter++
	}

	if currency != "" {
		if whereAdded {
			searchQuery += fmt.Sprintf(" AND p.currency_code = $%d", paramCounter)
		} else {
			searchQuery += fmt.Sprintf(" WHERE p.currency_code = $%d", paramCounter)
			whereAdded = true
		}
		queryParams = append(queryParams, currency)
		paramCounter++
	}

	if categoryID > 0 {
		if whereAdded {
			searchQuery += fmt.Sprintf(" AND p.category_id = $%d", paramCounter)
		} else {
			searchQuery += fmt.Sprintf(" WHERE p.category_id = $%d", paramCounter)
			whereAdded = true
		}
		queryParams = append(queryParams, categoryID)
		paramCounter++
	}

	if minPriceCents > 0 {
		if whereAdded {
			searchQuery += fmt.Sprintf(" AND COALESCE(pv.price, p.price) >= $%d", paramCounter)
		} else {
			searchQuery += fmt.Sprintf(" WHERE COALESCE(pv.price, p.price) >= $%d", paramCounter)
			whereAdded = true
		}
		queryParams = append(queryParams, minPriceCents) // Use cents
		paramCounter++
	}

	if maxPriceCents > 0 {
		if whereAdded {
			searchQuery += fmt.Sprintf(" AND COALESCE(pv.price, p.price) <= $%d", paramCounter)
		} else {
			searchQuery += fmt.Sprintf(" WHERE COALESCE(pv.price, p.price) <= $%d", paramCounter)
		}
		queryParams = append(queryParams, maxPriceCents) // Use cents
		paramCounter++
	}

	// Add pagination
	searchQuery += " ORDER BY p.created_at DESC LIMIT $" + strconv.Itoa(paramCounter) + " OFFSET $" + strconv.Itoa(paramCounter+1)
	queryParams = append(queryParams, limit, offset)

	// Execute query
	rows, err := r.db.Query(searchQuery, queryParams...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	products := []*entity.Product{}
	for rows.Next() {
		var imagesJSON []byte
		product := &entity.Product{}
		var productNumber sql.NullString

		err := rows.Scan(
			&product.ID,
			&productNumber,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.CurrencyCode,
			&product.Stock,
			&product.Weight,
			&product.CategoryID,
			&imagesJSON,
			&product.HasVariants,
			&product.Active,
			&product.CreatedAt,
			&product.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Set product number if valid
		if productNumber.Valid {
			product.ProductNumber = productNumber.String
		}

		// Unmarshal images JSON
		if err := json.Unmarshal(imagesJSON, &product.Images); err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return products, nil
}

func (r *ProductRepository) Count(searchQuery, currency string, categoryID uint, minPriceCents, maxPriceCents int64, active bool) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM products p
		LEFT JOIN product_variants pv ON p.id = pv.product_id AND pv.is_default = true
	`

	queryParams := []any{}
	paramCounter := 1
	var whereAdded bool

	if active {
		query += " WHERE p.active = true"
		whereAdded = true
	}

	if searchQuery != "" {
		if whereAdded {
			query += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.description ILIKE $%d)", paramCounter, paramCounter)
		} else {
			query += fmt.Sprintf(" WHERE (p.name ILIKE $%d OR p.description ILIKE $%d)", paramCounter, paramCounter)
			whereAdded = true
		}
		queryParams = append(queryParams, "%"+searchQuery+"%")
		paramCounter++
	}

	if categoryID > 0 {
		if whereAdded {
			query += fmt.Sprintf(" AND p.category_id = $%d", paramCounter)
		} else {
			query += fmt.Sprintf(" WHERE p.category_id = $%d", paramCounter)
			whereAdded = true
		}
		queryParams = append(queryParams, categoryID)
		paramCounter++
	}

	if minPriceCents > 0 {
		if whereAdded {
			query += fmt.Sprintf(" AND COALESCE(pv.price, p.price) >= $%d", paramCounter)
		} else {
			query += fmt.Sprintf(" WHERE COALESCE(pv.price, p.price) >= $%d", paramCounter)
			whereAdded = true
		}
		queryParams = append(queryParams, minPriceCents)
		paramCounter++
	}

	if maxPriceCents > 0 {
		if whereAdded {
			query += fmt.Sprintf(" AND COALESCE(pv.price, p.price) <= $%d", paramCounter)
		} else {
			query += fmt.Sprintf(" WHERE COALESCE(pv.price, p.price) <= $%d", paramCounter)
		}
		queryParams = append(queryParams, maxPriceCents)
		paramCounter++
	}

	var count int
	err := r.db.QueryRow(query, queryParams...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *ProductRepository) CountSearch(searchQuery string, categoryID uint, minPriceCents, maxPriceCents int64) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM products p
		LEFT JOIN product_variants pv ON p.id = pv.product_id AND pv.is_default = true
		WHERE p.active = true
	`

	queryParams := []any{}
	paramCounter := 1

	if searchQuery != "" {
		query += fmt.Sprintf(" AND (p.name ILIKE $%d OR p.description ILIKE $%d)", paramCounter, paramCounter)
		queryParams = append(queryParams, "%"+searchQuery+"%")
		paramCounter++
	}

	if categoryID > 0 {
		query += fmt.Sprintf(" AND p.category_id = $%d", paramCounter)
		queryParams = append(queryParams, categoryID)
		paramCounter++
	}

	if minPriceCents > 0 {
		query += fmt.Sprintf(" AND COALESCE(pv.price, p.price) >= $%d", paramCounter)
		queryParams = append(queryParams, minPriceCents)
		paramCounter++
	}

	if maxPriceCents > 0 {
		query += fmt.Sprintf(" AND COALESCE(pv.price, p.price) <= $%d", paramCounter)
		queryParams = append(queryParams, maxPriceCents)
		paramCounter++
	}

	var count int
	err := r.db.QueryRow(query, queryParams...).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
