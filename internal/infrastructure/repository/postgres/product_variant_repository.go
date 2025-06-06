package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// ProductVariantRepository is the PostgreSQL implementation of the ProductVariantRepository interface
type ProductVariantRepository struct {
	db *sql.DB
}

// NewProductVariantRepository creates a new ProductVariantRepository
func NewProductVariantRepository(db *sql.DB) repository.ProductVariantRepository {
	return &ProductVariantRepository{
		db: db,
	}
}

// Create creates a new product variant
func (r *ProductVariantRepository) Create(variant *entity.ProductVariant) error {
	query := `
		INSERT INTO product_variants (product_id, sku, price, currency_code, stock, attributes, images, is_default, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`

	// Marshal attributes directly
	attributesJSON, err := json.Marshal(variant.Attributes)
	if err != nil {
		return err
	}

	// Convert images to JSON
	imagesJSON, err := json.Marshal(variant.Images)
	if err != nil {
		return err
	}

	err = r.db.QueryRow(
		query,
		variant.ProductID,
		variant.SKU,
		variant.Price,
		variant.CurrencyCode,
		variant.Stock,
		attributesJSON,
		imagesJSON,
		variant.IsDefault,
		variant.CreatedAt,
		variant.UpdatedAt,
	).Scan(&variant.ID)

	if err != nil {
		// Check for duplicate SKU error
		if strings.Contains(err.Error(), "product_variants_sku_key") {
			return errors.New("a variant with this SKU already exists")
		}
		return err
	}

	// If the variant has currency-specific prices, save them
	if len(variant.Prices) > 0 {
		for i := range variant.Prices {
			variant.Prices[i].VariantID = variant.ID
			if err = r.createVariantPrice(&variant.Prices[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

// createVariantPrice creates a variant price entry for a specific currency
func (r *ProductVariantRepository) createVariantPrice(price *entity.ProductVariantPrice) error {
	query := `
		INSERT INTO product_variant_prices (variant_id, currency_code, price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (variant_id, currency_code) DO UPDATE SET
			price = EXCLUDED.price,
			updated_at = EXCLUDED.updated_at
		RETURNING id
	`

	now := time.Now()

	return r.db.QueryRow(
		query,
		price.VariantID,
		price.CurrencyCode,
		price.Price,
		now,
		now,
	).Scan(&price.ID)
}

// GetByID gets a variant by ID
func (r *ProductVariantRepository) GetByID(variantID uint) (*entity.ProductVariant, error) {
	query := `
		SELECT id, product_id, sku, price, currency_code, stock, attributes, images, is_default, created_at, updated_at
		FROM product_variants
		WHERE id = $1
	`

	var attributesJSON, imagesJSON []byte
	variant := &entity.ProductVariant{}

	err := r.db.QueryRow(query, variantID).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.SKU,
		&variant.Price,
		&variant.CurrencyCode,
		&variant.Stock,
		&attributesJSON,
		&imagesJSON,
		&variant.IsDefault,
		&variant.CreatedAt,
		&variant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("variant not found")
		}
		return nil, err
	}

	// Unmarshal attributes JSON directly into VariantAttribute slice
	if err := json.Unmarshal(attributesJSON, &variant.Attributes); err != nil {
		return nil, err
	}

	// Unmarshal images JSON
	if err := json.Unmarshal(imagesJSON, &variant.Images); err != nil {
		return nil, err
	}

	// Load currency-specific prices
	prices, err := r.getVariantPrices(variant.ID)
	if err != nil {
		return nil, err
	}
	variant.Prices = prices

	return variant, nil
}

// getVariantPrices retrieves all prices for a variant in different currencies
func (r *ProductVariantRepository) getVariantPrices(variantID uint) ([]entity.ProductVariantPrice, error) {
	query := `
		SELECT id, variant_id, currency_code, price, created_at, updated_at
		FROM product_variant_prices
		WHERE variant_id = $1
	`

	rows, err := r.db.Query(query, variantID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []entity.ProductVariantPrice
	for rows.Next() {
		var price entity.ProductVariantPrice

		err := rows.Scan(
			&price.ID,
			&price.VariantID,
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

// Update updates a product variant
func (r *ProductVariantRepository) Update(variant *entity.ProductVariant) error {
	query := `
		UPDATE product_variants
		SET sku = $1, price = $2, currency_code = $3, stock = $4, 
		    attributes = $5, images = $6, is_default = $7, updated_at = $8
		WHERE id = $9
	`

	// Marshal attributes directly
	attributesJSON, err := json.Marshal(variant.Attributes)
	if err != nil {
		return err
	}

	// Convert images to JSON
	imagesJSON, err := json.Marshal(variant.Images)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(
		query,
		variant.SKU,
		variant.Price,
		variant.CurrencyCode,
		variant.Stock,
		attributesJSON,
		imagesJSON,
		variant.IsDefault,
		time.Now(),
		variant.ID,
	)

	if err != nil {
		return err
	}

	// Update currency-specific prices
	if len(variant.Prices) > 0 {
		// First, delete existing prices (to handle removes)
		if _, err := r.db.Exec("DELETE FROM product_variant_prices WHERE variant_id = $1", variant.ID); err != nil {
			return err
		}

		// Then add all current prices
		for i := range variant.Prices {
			variant.Prices[i].VariantID = variant.ID
			if err := r.createVariantPrice(&variant.Prices[i]); err != nil {
				return err
			}
		}
	}

	return nil
}

// Delete deletes a product variant
func (r *ProductVariantRepository) Delete(variantID uint) error {
	// Check if this is the only variant or if it's the default variant
	var isDefault bool
	var productID uint
	var variantCount int

	// First verify the variant exists and get its product ID
	err := r.db.QueryRow(
		"SELECT is_default, product_id FROM product_variants WHERE id = $1",
		variantID,
	).Scan(&isDefault, &productID)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("variant not found")
		}
		return err
	}

	// Count variants for this product
	err = r.db.QueryRow(
		"SELECT COUNT(*) FROM product_variants WHERE product_id = $1",
		productID,
	).Scan(&variantCount)
	if err != nil {
		return err
	}

	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Delete the variant
	result, err := tx.Exec("DELETE FROM product_variants WHERE id = $1", variantID)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("variant not found or already deleted")
	}

	// If this deletion will leave only one variant, update product to not have multiple variants
	if variantCount == 2 {
		_, err = tx.Exec(
			"UPDATE products SET has_variants = false WHERE id = $1",
			productID,
		)
		if err != nil {
			return err
		}
	} else if isDefault {
		// If this was the default variant, set another variant as default
		_, err = tx.Exec(`
			UPDATE product_variants 
			SET is_default = true 
			WHERE id = (
				SELECT id 
				FROM product_variants 
				WHERE product_id = $1 
				AND id != $2 
				ORDER BY id ASC 
				LIMIT 1
			)
		`, productID, variantID)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// GetByProduct gets all variants for a product
func (r *ProductVariantRepository) GetByProduct(productID uint) ([]*entity.ProductVariant, error) {
	query := `
		SELECT id, product_id, sku, price, currency_code, stock, attributes, images, is_default, created_at, updated_at
		FROM product_variants
		WHERE product_id = $1
		ORDER BY is_default DESC, id ASC
	`

	rows, err := r.db.Query(query, productID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	variants := []*entity.ProductVariant{}
	for rows.Next() {
		var attributesJSON, imagesJSON []byte
		variant := &entity.ProductVariant{}

		err := rows.Scan(
			&variant.ID,
			&variant.ProductID,
			&variant.SKU,
			&variant.Price,
			&variant.CurrencyCode,
			&variant.Stock,
			&attributesJSON,
			&imagesJSON,
			&variant.IsDefault,
			&variant.CreatedAt,
			&variant.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}

		// Unmarshal attributes JSON directly into VariantAttribute slice
		if err := json.Unmarshal(attributesJSON, &variant.Attributes); err != nil {
			return nil, err
		}

		// Unmarshal images JSON
		if err := json.Unmarshal(imagesJSON, &variant.Images); err != nil {
			return nil, err
		}

		// Load currency-specific prices
		prices, err := r.getVariantPrices(variant.ID)
		if err != nil {
			return nil, err
		}

		variant.Prices = prices

		variants = append(variants, variant)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return variants, nil
}

// GetBySKU gets a variant by SKU
func (r *ProductVariantRepository) GetBySKU(sku string) (*entity.ProductVariant, error) {
	query := `
		SELECT id, product_id, sku, price, currency_code, stock, attributes, images, is_default, created_at, updated_at
		FROM product_variants
		WHERE sku = $1
	`

	var attributesJSON, imagesJSON []byte
	variant := &entity.ProductVariant{}

	err := r.db.QueryRow(query, sku).Scan(
		&variant.ID,
		&variant.ProductID,
		&variant.SKU,
		&variant.Price,
		&variant.CurrencyCode,
		&variant.Stock,
		&attributesJSON,
		&imagesJSON,
		&variant.IsDefault,
		&variant.CreatedAt,
		&variant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("variant not found")
		}
		return nil, err
	}

	// Unmarshal attributes JSON directly into VariantAttribute slice
	if err := json.Unmarshal(attributesJSON, &variant.Attributes); err != nil {
		return nil, err
	}

	// Unmarshal images JSON
	if err := json.Unmarshal(imagesJSON, &variant.Images); err != nil {
		return nil, err
	}

	// Load currency-specific prices
	prices, err := r.getVariantPrices(variant.ID)
	if err != nil {
		return nil, err
	}
	variant.Prices = prices

	return variant, nil
}

func (r *ProductVariantRepository) BatchCreate(variants []*entity.ProductVariant) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	for _, variant := range variants {
		err = r.Create(variant)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}
