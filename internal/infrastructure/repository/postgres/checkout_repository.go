package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// CheckoutRepository implements the checkout repository interface using PostgreSQL
type CheckoutRepository struct {
	db *sql.DB
}

// NewCheckoutRepository creates a new CheckoutRepository
func NewCheckoutRepository(db *sql.DB) repository.CheckoutRepository {
	return &CheckoutRepository{db: db}
}

// Create creates a new checkout
func (r *CheckoutRepository) Create(checkout *entity.Checkout) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Marshal addresses and customer details to JSON
	shippingAddrJSON, err := json.Marshal(checkout.ShippingAddr)
	if err != nil {
		return err
	}

	billingAddrJSON, err := json.Marshal(checkout.BillingAddr)
	if err != nil {
		return err
	}

	customerDetailsJSON, err := json.Marshal(checkout.CustomerDetails)
	if err != nil {
		return err
	}

	// Marshal applied discount to JSON if it exists
	var appliedDiscountJSON []byte = []byte("null") // Default to JSON null
	if checkout.AppliedDiscount != nil {
		appliedDiscountJSON, err = json.Marshal(checkout.AppliedDiscount)
		if err != nil {
			return err
		}
	}

	// Insert checkout
	query := `
		INSERT INTO checkouts (
			user_id, session_id, status, shipping_address, billing_address,
			shipping_method_id, payment_provider, total_amount, shipping_cost,
			total_weight, customer_details, currency, discount_code,
			discount_amount, final_amount, applied_discount, created_at,
			updated_at, last_activity_at, expires_at, completed_at,
			converted_order_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14,
			$15, $16, $17, $18, $19, $20, $21, $22
		) RETURNING id`

	var userID sql.NullInt64
	if checkout.UserID != 0 {
		userID.Int64 = int64(checkout.UserID)
		userID.Valid = true
	}

	var shippingMethodID sql.NullInt64
	if checkout.ShippingMethodID != 0 {
		shippingMethodID.Int64 = int64(checkout.ShippingMethodID)
		shippingMethodID.Valid = true
	}

	var completedAt sql.NullTime
	if checkout.CompletedAt != nil {
		completedAt.Time = *checkout.CompletedAt
		completedAt.Valid = true
	}

	var convertedOrderID sql.NullInt64
	if checkout.ConvertedOrderID != 0 {
		convertedOrderID.Int64 = int64(checkout.ConvertedOrderID)
		convertedOrderID.Valid = true
	}

	var paymentProviderNull sql.NullString
	if checkout.PaymentProvider != "" {
		paymentProviderNull.String = checkout.PaymentProvider
		paymentProviderNull.Valid = true
	}

	var discountCodeNull sql.NullString
	if checkout.DiscountCode != "" {
		discountCodeNull.String = checkout.DiscountCode
		discountCodeNull.Valid = true
	}

	// Execute query
	row := tx.QueryRow(
		query,
		userID, checkout.SessionID, checkout.Status, shippingAddrJSON, billingAddrJSON,
		shippingMethodID, paymentProviderNull, checkout.TotalAmount, checkout.ShippingCost,
		checkout.TotalWeight, customerDetailsJSON, checkout.Currency, discountCodeNull,
		checkout.DiscountAmount, checkout.FinalAmount, appliedDiscountJSON, checkout.CreatedAt,
		checkout.UpdatedAt, checkout.LastActivityAt, checkout.ExpiresAt, completedAt,
		convertedOrderID,
	)

	var id uint
	if err := row.Scan(&id); err != nil {
		return err
	}
	checkout.ID = id

	// Insert checkout items
	if len(checkout.Items) > 0 {
		for i := range checkout.Items {
			item := &checkout.Items[i]
			item.CheckoutID = checkout.ID

			var productVariantIDNull sql.NullInt64
			if item.ProductVariantID != 0 {
				productVariantIDNull.Int64 = int64(item.ProductVariantID)
				productVariantIDNull.Valid = true
			}

			var variantNameNull sql.NullString
			if item.VariantName != "" {
				variantNameNull.String = item.VariantName
				variantNameNull.Valid = true
			}

			var skuNull sql.NullString
			if item.SKU != "" {
				skuNull.String = item.SKU
				skuNull.Valid = true
			}

			itemQuery := `
				INSERT INTO checkout_items (
					checkout_id, product_id, product_variant_id, quantity,
					price, weight, product_name, variant_name, sku,
					created_at, updated_at
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
				) RETURNING id`

			var itemID uint
			err = tx.QueryRow(
				itemQuery,
				checkout.ID, item.ProductID, productVariantIDNull, item.Quantity,
				item.Price, item.Weight, item.ProductName, variantNameNull, skuNull,
				item.CreatedAt, item.UpdatedAt,
			).Scan(&itemID)

			if err != nil {
				return err
			}

			item.ID = itemID
		}
	}

	return nil
}

// GetByID retrieves a checkout by ID
func (r *CheckoutRepository) GetByID(checkoutID uint) (*entity.Checkout, error) {
	query := `
		SELECT 
			id, user_id, session_id, status, shipping_address, 
			billing_address, shipping_method_id, payment_provider, 
			total_amount, shipping_cost, total_weight, customer_details, 
			currency, discount_code, discount_amount, final_amount, 
			applied_discount, created_at, updated_at, last_activity_at, 
			expires_at, completed_at, converted_order_id
		FROM checkouts 
		WHERE id = $1`

	checkout, err := r.scanCheckout(r.db.QueryRow(query, checkoutID))
	if err != nil {
		return nil, err
	}

	// Get checkout items
	itemsQuery := `
		SELECT 
			id, checkout_id, product_id, product_variant_id, quantity, 
			price, weight, product_name, variant_name, sku, 
			created_at, updated_at
		FROM checkout_items 
		WHERE checkout_id = $1
		ORDER BY id ASC`

	rows, err := r.db.Query(itemsQuery, checkoutID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []entity.CheckoutItem{}
	for rows.Next() {
		item, err := r.scanCheckoutItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	checkout.Items = items
	return checkout, nil
}

// GetByUserID retrieves an active checkout by user ID
func (r *CheckoutRepository) GetByUserID(userID uint) (*entity.Checkout, error) {
	query := `
		SELECT 
			id, user_id, session_id, status, shipping_address, 
			billing_address, shipping_method_id, payment_provider, 
			total_amount, shipping_cost, total_weight, customer_details, 
			currency, discount_code, discount_amount, final_amount, 
			applied_discount, created_at, updated_at, last_activity_at, 
			expires_at, completed_at, converted_order_id
		FROM checkouts 
		WHERE user_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT 1`

	checkout, err := r.scanCheckout(r.db.QueryRow(query, userID, entity.CheckoutStatusActive))
	if err != nil {
		return nil, err
	}

	// Get checkout items
	itemsQuery := `
		SELECT 
			id, checkout_id, product_id, product_variant_id, quantity, 
			price, weight, product_name, variant_name, sku, 
			created_at, updated_at
		FROM checkout_items 
		WHERE checkout_id = $1
		ORDER BY id ASC`

	rows, err := r.db.Query(itemsQuery, checkout.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []entity.CheckoutItem{}
	for rows.Next() {
		item, err := r.scanCheckoutItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	checkout.Items = items
	return checkout, nil
}

// GetBySessionID retrieves an active checkout by session ID
func (r *CheckoutRepository) GetBySessionID(sessionID string) (*entity.Checkout, error) {
	query := `
		SELECT 
			id, user_id, session_id, status, shipping_address, 
			billing_address, shipping_method_id, payment_provider, 
			total_amount, shipping_cost, total_weight, customer_details, 
			currency, discount_code, discount_amount, final_amount, 
			applied_discount, created_at, updated_at, last_activity_at, 
			expires_at, completed_at, converted_order_id
		FROM checkouts 
		WHERE session_id = $1 AND status = $2
		ORDER BY created_at DESC
		LIMIT 1`

	checkout, err := r.scanCheckout(r.db.QueryRow(query, sessionID, entity.CheckoutStatusActive))
	if err != nil {
		return nil, err
	}

	// Get checkout items
	itemsQuery := `
		SELECT 
			id, checkout_id, product_id, product_variant_id, quantity, 
			price, weight, product_name, variant_name, sku, 
			created_at, updated_at
		FROM checkout_items 
		WHERE checkout_id = $1
		ORDER BY id ASC`

	rows, err := r.db.Query(itemsQuery, checkout.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := []entity.CheckoutItem{}
	for rows.Next() {
		item, err := r.scanCheckoutItem(rows)
		if err != nil {
			return nil, err
		}
		items = append(items, *item)
	}

	checkout.Items = items
	return checkout, nil
}

// Update updates a checkout
func (r *CheckoutRepository) Update(checkout *entity.Checkout) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Marshal addresses and customer details to JSON
	shippingAddrJSON, err := json.Marshal(checkout.ShippingAddr)
	if err != nil {
		return err
	}

	billingAddrJSON, err := json.Marshal(checkout.BillingAddr)
	if err != nil {
		return err
	}

	customerDetailsJSON, err := json.Marshal(checkout.CustomerDetails)
	if err != nil {
		return err
	}

	// Marshal applied discount to JSON if it exists
	var appliedDiscountJSON []byte = []byte("null") // Default to JSON null
	if checkout.AppliedDiscount != nil {
		appliedDiscountJSON, err = json.Marshal(checkout.AppliedDiscount)
		if err != nil {
			return err
		}
	}

	// Update checkout
	query := `
		UPDATE checkouts
		SET 
			user_id = $1,
			session_id = $2,
			status = $3,
			shipping_address = $4,
			billing_address = $5,
			shipping_method_id = $6,
			payment_provider = $7,
			total_amount = $8,
			shipping_cost = $9,
			total_weight = $10,
			customer_details = $11,
			currency = $12,
			discount_code = $13,
			discount_amount = $14,
			final_amount = $15,
			applied_discount = $16,
			updated_at = $17,
			last_activity_at = $18,
			expires_at = $19,
			completed_at = $20,
			converted_order_id = $21
		WHERE id = $22`

	var userID sql.NullInt64
	if checkout.UserID != 0 {
		userID.Int64 = int64(checkout.UserID)
		userID.Valid = true
	}

	var shippingMethodID sql.NullInt64
	if checkout.ShippingMethodID != 0 {
		shippingMethodID.Int64 = int64(checkout.ShippingMethodID)
		shippingMethodID.Valid = true
	}

	var completedAt sql.NullTime
	if checkout.CompletedAt != nil {
		completedAt.Time = *checkout.CompletedAt
		completedAt.Valid = true
	}

	var convertedOrderID sql.NullInt64
	if checkout.ConvertedOrderID != 0 {
		convertedOrderID.Int64 = int64(checkout.ConvertedOrderID)
		convertedOrderID.Valid = true
	}

	var paymentProviderNull sql.NullString
	if checkout.PaymentProvider != "" {
		paymentProviderNull.String = checkout.PaymentProvider
		paymentProviderNull.Valid = true
	}

	var discountCodeNull sql.NullString
	if checkout.DiscountCode != "" {
		discountCodeNull.String = checkout.DiscountCode
		discountCodeNull.Valid = true
	}

	// Execute update query
	_, err = tx.Exec(
		query,
		userID, checkout.SessionID, checkout.Status, shippingAddrJSON, billingAddrJSON,
		shippingMethodID, paymentProviderNull, checkout.TotalAmount, checkout.ShippingCost,
		checkout.TotalWeight, customerDetailsJSON, checkout.Currency, discountCodeNull,
		checkout.DiscountAmount, checkout.FinalAmount, appliedDiscountJSON, checkout.UpdatedAt,
		checkout.LastActivityAt, checkout.ExpiresAt, completedAt, convertedOrderID, checkout.ID,
	)
	if err != nil {
		return err
	}

	// Delete existing checkout items
	_, err = tx.Exec("DELETE FROM checkout_items WHERE checkout_id = $1", checkout.ID)
	if err != nil {
		return err
	}

	// Insert updated checkout items
	if len(checkout.Items) > 0 {
		for i := range checkout.Items {
			item := &checkout.Items[i]
			item.CheckoutID = checkout.ID

			var productVariantIDNull sql.NullInt64
			if item.ProductVariantID != 0 {
				productVariantIDNull.Int64 = int64(item.ProductVariantID)
				productVariantIDNull.Valid = true
			}

			var variantNameNull sql.NullString
			if item.VariantName != "" {
				variantNameNull.String = item.VariantName
				variantNameNull.Valid = true
			}

			var skuNull sql.NullString
			if item.SKU != "" {
				skuNull.String = item.SKU
				skuNull.Valid = true
			}

			itemQuery := `
				INSERT INTO checkout_items (
					checkout_id, product_id, product_variant_id, quantity,
					price, weight, product_name, variant_name, sku,
					created_at, updated_at
				) VALUES (
					$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11
				) RETURNING id`

			var itemID uint
			err = tx.QueryRow(
				itemQuery,
				checkout.ID, item.ProductID, productVariantIDNull, item.Quantity,
				item.Price, item.Weight, item.ProductName, variantNameNull, skuNull,
				item.CreatedAt, item.UpdatedAt,
			).Scan(&itemID)

			if err != nil {
				return err
			}

			item.ID = itemID
		}
	}

	return nil
}

// Delete deletes a checkout
func (r *CheckoutRepository) Delete(checkoutID uint) error {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// First delete checkout items
	_, err = tx.Exec("DELETE FROM checkout_items WHERE checkout_id = $1", checkoutID)
	if err != nil {
		return err
	}

	// Then delete checkout
	result, err := tx.Exec("DELETE FROM checkouts WHERE id = $1", checkoutID)
	if err != nil {
		return err
	}

	// Check if any rows were affected
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("checkout with ID %d not found", checkoutID)
	}

	return nil
}

// ConvertGuestCheckoutToUserCheckout converts a guest checkout to a user checkout
func (r *CheckoutRepository) ConvertGuestCheckoutToUserCheckout(sessionID string, userID uint) (*entity.Checkout, error) {
	// Start a transaction
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	// Find the guest checkout
	query := `
		UPDATE checkouts
		SET user_id = $1, updated_at = $2, last_activity_at = $3
		WHERE session_id = $4 AND status = $5
		RETURNING id`

	now := time.Now()
	var checkoutID uint
	err = tx.QueryRow(
		query, userID, now, now, sessionID, entity.CheckoutStatusActive,
	).Scan(&checkoutID)
	if err != nil {
		return nil, err
	}

	// Get updated checkout
	checkout, err := r.GetByID(checkoutID)
	if err != nil {
		return nil, err
	}

	return checkout, nil
}

// GetExpiredCheckouts retrieves all checkouts that have expired
func (r *CheckoutRepository) GetExpiredCheckouts() ([]*entity.Checkout, error) {
	query := `
		SELECT 
			id, user_id, session_id, status, shipping_address, 
			billing_address, shipping_method_id, payment_provider, 
			total_amount, shipping_cost, total_weight, customer_details, 
			currency, discount_code, discount_amount, final_amount, 
			applied_discount, created_at, updated_at, last_activity_at, 
			expires_at, completed_at, converted_order_id
		FROM checkouts 
		WHERE status = $1 AND expires_at < $2`

	rows, err := r.db.Query(query, entity.CheckoutStatusActive, time.Now())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checkouts := []*entity.Checkout{}
	for rows.Next() {
		checkout, err := r.scanCheckout(rows)
		if err != nil {
			return nil, err
		}

		// Get checkout items
		itemsQuery := `
			SELECT 
				id, checkout_id, product_id, product_variant_id, quantity, 
				price, weight, product_name, variant_name, sku, 
				created_at, updated_at
			FROM checkout_items 
			WHERE checkout_id = $1
			ORDER BY id ASC`

		itemRows, err := r.db.Query(itemsQuery, checkout.ID)
		if err != nil {
			return nil, err
		}

		items := []entity.CheckoutItem{}
		for itemRows.Next() {
			item, err := r.scanCheckoutItem(itemRows)
			if err != nil {
				itemRows.Close()
				return nil, err
			}
			items = append(items, *item)
		}
		itemRows.Close()

		checkout.Items = items
		checkouts = append(checkouts, checkout)
	}

	return checkouts, nil
}

// GetCheckoutsByStatus retrieves checkouts by status
func (r *CheckoutRepository) GetCheckoutsByStatus(status entity.CheckoutStatus, offset, limit int) ([]*entity.Checkout, error) {
	var query string
	var args []interface{}

	if status == "" {
		query = `
			SELECT 
				id, user_id, session_id, status, shipping_address, 
				billing_address, shipping_method_id, payment_provider, 
				total_amount, shipping_cost, total_weight, customer_details, 
				currency, discount_code, discount_amount, final_amount, 
				applied_discount, created_at, updated_at, last_activity_at, 
				expires_at, completed_at, converted_order_id
			FROM checkouts 
			ORDER BY created_at DESC
			OFFSET $1 LIMIT $2`

		args = []interface{}{offset, limit}
	} else {
		query = `
			SELECT 
				id, user_id, session_id, status, shipping_address, 
				billing_address, shipping_method_id, payment_provider, 
				total_amount, shipping_cost, total_weight, customer_details, 
				currency, discount_code, discount_amount, final_amount, 
				applied_discount, created_at, updated_at, last_activity_at, 
				expires_at, completed_at, converted_order_id
			FROM checkouts 
			WHERE status = $1
			ORDER BY created_at DESC
			OFFSET $2 LIMIT $3`

		args = []interface{}{status, offset, limit}
	}

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checkouts := []*entity.Checkout{}
	for rows.Next() {
		checkout, err := r.scanCheckout(rows)
		if err != nil {
			return nil, err
		}

		// Get checkout items
		itemsQuery := `
			SELECT 
				id, checkout_id, product_id, product_variant_id, quantity, 
				price, weight, product_name, variant_name, sku, 
				created_at, updated_at
			FROM checkout_items 
			WHERE checkout_id = $1
			ORDER BY id ASC`

		itemRows, err := r.db.Query(itemsQuery, checkout.ID)
		if err != nil {
			return nil, err
		}

		items := []entity.CheckoutItem{}
		for itemRows.Next() {
			item, err := r.scanCheckoutItem(itemRows)
			if err != nil {
				itemRows.Close()
				return nil, err
			}
			items = append(items, *item)
		}
		itemRows.Close()

		checkout.Items = items
		checkouts = append(checkouts, checkout)
	}

	return checkouts, nil
}

// GetActiveCheckoutsByUserID retrieves all active checkouts for a user
func (r *CheckoutRepository) GetActiveCheckoutsByUserID(userID uint) ([]*entity.Checkout, error) {
	query := `
		SELECT 
			id, user_id, session_id, status, shipping_address, 
			billing_address, shipping_method_id, payment_provider, 
			total_amount, shipping_cost, total_weight, customer_details, 
			currency, discount_code, discount_amount, final_amount, 
			applied_discount, created_at, updated_at, last_activity_at, 
			expires_at, completed_at, converted_order_id
		FROM checkouts 
		WHERE user_id = $1 AND status = $2
		ORDER BY created_at DESC`

	rows, err := r.db.Query(query, userID, entity.CheckoutStatusActive)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checkouts := []*entity.Checkout{}
	for rows.Next() {
		checkout, err := r.scanCheckout(rows)
		if err != nil {
			return nil, err
		}

		// Get checkout items
		itemsQuery := `
			SELECT 
				id, checkout_id, product_id, product_variant_id, quantity, 
				price, weight, product_name, variant_name, sku, 
				created_at, updated_at
			FROM checkout_items 
			WHERE checkout_id = $1
			ORDER BY id ASC`

		itemRows, err := r.db.Query(itemsQuery, checkout.ID)
		if err != nil {
			return nil, err
		}

		items := []entity.CheckoutItem{}
		for itemRows.Next() {
			item, err := r.scanCheckoutItem(itemRows)
			if err != nil {
				itemRows.Close()
				return nil, err
			}
			items = append(items, *item)
		}
		itemRows.Close()

		checkout.Items = items
		checkouts = append(checkouts, checkout)
	}

	return checkouts, nil
}

// GetCompletedCheckoutsByUserID retrieves all completed checkouts for a user
func (r *CheckoutRepository) GetCompletedCheckoutsByUserID(userID uint, offset, limit int) ([]*entity.Checkout, error) {
	query := `
		SELECT 
			id, user_id, session_id, status, shipping_address, 
			billing_address, shipping_method_id, payment_provider, 
			total_amount, shipping_cost, total_weight, customer_details, 
			currency, discount_code, discount_amount, final_amount, 
			applied_discount, created_at, updated_at, last_activity_at, 
			expires_at, completed_at, converted_order_id
		FROM checkouts 
		WHERE user_id = $1 AND status = $2
		ORDER BY created_at DESC
		OFFSET $3 LIMIT $4`

	rows, err := r.db.Query(query, userID, entity.CheckoutStatusCompleted, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	checkouts := []*entity.Checkout{}
	for rows.Next() {
		checkout, err := r.scanCheckout(rows)
		if err != nil {
			return nil, err
		}

		// Get checkout items
		itemsQuery := `
			SELECT 
				id, checkout_id, product_id, product_variant_id, quantity, 
				price, weight, product_name, variant_name, sku, 
				created_at, updated_at
			FROM checkout_items 
			WHERE checkout_id = $1
			ORDER BY id ASC`

		itemRows, err := r.db.Query(itemsQuery, checkout.ID)
		if err != nil {
			return nil, err
		}

		items := []entity.CheckoutItem{}
		for itemRows.Next() {
			item, err := r.scanCheckoutItem(itemRows)
			if err != nil {
				itemRows.Close()
				return nil, err
			}
			items = append(items, *item)
		}
		itemRows.Close()

		checkout.Items = items
		checkouts = append(checkouts, checkout)
	}

	return checkouts, nil
}

// HasActiveCheckoutsWithProduct checks if a product has any associated active checkouts
func (r *CheckoutRepository) HasActiveCheckoutsWithProduct(productID uint) (bool, error) {
	if productID == 0 {
		return false, errors.New("product ID cannot be 0")
	}

	query := `
		SELECT EXISTS(
			SELECT 1 FROM checkout_items ci
			JOIN checkouts c ON ci.checkout_id = c.id
			WHERE ci.product_id = $1 
			AND c.status = 'active'
		)
	`

	var exists bool
	err := r.db.QueryRow(query, productID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if product has active checkouts: %w", err)
	}

	return exists, nil
}

// Helper to scan checkout rows
func (r *CheckoutRepository) scanCheckout(row interface{}) (*entity.Checkout, error) {
	var checkout entity.Checkout
	var userID sql.NullInt64
	var shippingMethodID sql.NullInt64
	var completedAt sql.NullTime
	var convertedOrderID sql.NullInt64
	var shippingAddrJSON, billingAddrJSON, customerDetailsJSON []byte
	var appliedDiscountJSON sql.NullString
	var paymentProvider, discountCode sql.NullString
	var sessionID sql.NullString

	var scanner interface {
		Scan(...interface{}) error
	}

	switch v := row.(type) {
	case *sql.Row:
		scanner = v
	case *sql.Rows:
		scanner = v
	default:
		return nil, errors.New("invalid row type")
	}

	err := scanner.Scan(
		&checkout.ID,
		&userID,
		&sessionID,
		&checkout.Status,
		&shippingAddrJSON,
		&billingAddrJSON,
		&shippingMethodID,
		&paymentProvider,
		&checkout.TotalAmount,
		&checkout.ShippingCost,
		&checkout.TotalWeight,
		&customerDetailsJSON,
		&checkout.Currency,
		&discountCode,
		&checkout.DiscountAmount,
		&checkout.FinalAmount,
		&appliedDiscountJSON,
		&checkout.CreatedAt,
		&checkout.UpdatedAt,
		&checkout.LastActivityAt,
		&checkout.ExpiresAt,
		&completedAt,
		&convertedOrderID,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("checkout not found")
		}
		return nil, err
	}

	// Set values from nullable fields
	if userID.Valid {
		checkout.UserID = uint(userID.Int64)
	}

	if sessionID.Valid {
		checkout.SessionID = sessionID.String
	}

	if shippingMethodID.Valid {
		checkout.ShippingMethodID = uint(shippingMethodID.Int64)
	}

	if completedAt.Valid {
		checkout.CompletedAt = &completedAt.Time
	}

	if convertedOrderID.Valid {
		checkout.ConvertedOrderID = uint(convertedOrderID.Int64)
	}

	if paymentProvider.Valid {
		checkout.PaymentProvider = paymentProvider.String
	}

	if discountCode.Valid {
		checkout.DiscountCode = discountCode.String
	}

	// Unmarshal addresses
	if err := json.Unmarshal(shippingAddrJSON, &checkout.ShippingAddr); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(billingAddrJSON, &checkout.BillingAddr); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(customerDetailsJSON, &checkout.CustomerDetails); err != nil {
		return nil, err
	}

	// Unmarshal applied discount if it exists
	if appliedDiscountJSON.Valid && appliedDiscountJSON.String != "" {
		checkout.AppliedDiscount = &entity.AppliedDiscount{}
		if err := json.Unmarshal([]byte(appliedDiscountJSON.String), checkout.AppliedDiscount); err != nil {
			return nil, err
		}
	}

	return &checkout, nil
}

// Helper to scan checkout item rows
func (r *CheckoutRepository) scanCheckoutItem(rows *sql.Rows) (*entity.CheckoutItem, error) {
	var item entity.CheckoutItem
	var productVariantID sql.NullInt64
	var variantName, sku sql.NullString

	err := rows.Scan(
		&item.ID,
		&item.CheckoutID,
		&item.ProductID,
		&productVariantID,
		&item.Quantity,
		&item.Price,
		&item.Weight,
		&item.ProductName,
		&variantName,
		&sku,
		&item.CreatedAt,
		&item.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	if productVariantID.Valid {
		item.ProductVariantID = uint(productVariantID.Int64)
	}

	if variantName.Valid {
		item.VariantName = variantName.String
	}

	if sku.Valid {
		item.SKU = sku.String
	}

	return &item, nil
}

// GetCheckoutsToAbandon retrieves active checkouts with customer/shipping info that should be marked as abandoned
func (r *CheckoutRepository) GetCheckoutsToAbandon() ([]*entity.Checkout, error) {
	// Find active checkouts with customer or shipping info that haven't been active for 15 minutes
	abandonThreshold := time.Now().Add(-15 * time.Minute)

	query := `
		SELECT 
			id, user_id, session_id, status, shipping_address, 
			billing_address, shipping_method_id, payment_provider, 
			total_amount, shipping_cost, total_weight, customer_details, 
			currency, discount_code, discount_amount, final_amount, 
			applied_discount, created_at, updated_at, last_activity_at, 
			expires_at, completed_at, converted_order_id
		FROM checkouts 
		WHERE status = $1 
		AND last_activity_at < $2
		AND (
			(customer_details->>'email' != '' AND customer_details->>'email' IS NOT NULL)
			OR (customer_details->>'phone' != '' AND customer_details->>'phone' IS NOT NULL)
			OR (customer_details->>'full_name' != '' AND customer_details->>'full_name' IS NOT NULL)
			OR (shipping_address->>'street' != '' AND shipping_address->>'street' IS NOT NULL)
			OR (shipping_address->>'city' != '' AND shipping_address->>'city' IS NOT NULL)
			OR (shipping_address->>'state' != '' AND shipping_address->>'state' IS NOT NULL)
			OR (shipping_address->>'postal_code' != '' AND shipping_address->>'postal_code' IS NOT NULL)
			OR (shipping_address->>'country' != '' AND shipping_address->>'country' IS NOT NULL)
		)`

	rows, err := r.db.Query(query, entity.CheckoutStatusActive, abandonThreshold)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanCheckoutsWithItems(rows)
}

// GetCheckoutsToDelete retrieves checkouts that should be deleted
func (r *CheckoutRepository) GetCheckoutsToDelete() ([]*entity.Checkout, error) {
	now := time.Now()
	emptyDeleteThreshold := now.Add(-24 * time.Hour)
	abandonedDeleteThreshold := now.Add(-7 * 24 * time.Hour)

	query := `
		SELECT 
			id, user_id, session_id, status, shipping_address, 
			billing_address, shipping_method_id, payment_provider, 
			total_amount, shipping_cost, total_weight, customer_details, 
			currency, discount_code, discount_amount, final_amount, 
			applied_discount, created_at, updated_at, last_activity_at, 
			expires_at, completed_at, converted_order_id
		FROM checkouts 
		WHERE 
		(
			-- Delete empty checkouts after 24 hours
			(
				status = $1 
				AND last_activity_at < $2
				AND (customer_details->>'email' = '' OR customer_details->>'email' IS NULL)
				AND (customer_details->>'phone' = '' OR customer_details->>'phone' IS NULL)
				AND (customer_details->>'full_name' = '' OR customer_details->>'full_name' IS NULL)
				AND (shipping_address->>'street' = '' OR shipping_address->>'street' IS NULL)
				AND (shipping_address->>'city' = '' OR shipping_address->>'city' IS NULL)
				AND (shipping_address->>'state' = '' OR shipping_address->>'state' IS NULL)
				AND (shipping_address->>'postal_code' = '' OR shipping_address->>'postal_code' IS NULL)
				AND (shipping_address->>'country' = '' OR shipping_address->>'country' IS NULL)
			)
			OR
			-- Delete abandoned checkouts after 7 days
			(
				status = $3 
				AND updated_at < $4
			)
			OR
			-- Delete all expired checkouts
			(
				status = $5
			)
		)`

	rows, err := r.db.Query(query,
		entity.CheckoutStatusActive, emptyDeleteThreshold,
		entity.CheckoutStatusAbandoned, abandonedDeleteThreshold,
		entity.CheckoutStatusExpired)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return r.scanCheckoutsWithItems(rows)
}

// scanCheckoutsWithItems is a helper method to scan checkouts and their items
func (r *CheckoutRepository) scanCheckoutsWithItems(rows *sql.Rows) ([]*entity.Checkout, error) {
	checkouts := []*entity.Checkout{}
	for rows.Next() {
		checkout, err := r.scanCheckout(rows)
		if err != nil {
			return nil, err
		}

		// Get checkout items
		itemsQuery := `
			SELECT 
				id, checkout_id, product_id, product_variant_id, quantity, 
				price, weight, product_name, variant_name, sku, 
				created_at, updated_at
			FROM checkout_items 
			WHERE checkout_id = $1
			ORDER BY id ASC`

		itemRows, err := r.db.Query(itemsQuery, checkout.ID)
		if err != nil {
			return nil, err
		}

		items := []entity.CheckoutItem{}
		for itemRows.Next() {
			item, err := r.scanCheckoutItem(itemRows)
			if err != nil {
				itemRows.Close()
				return nil, err
			}
			items = append(items, *item)
		}
		itemRows.Close()

		checkout.Items = items
		checkouts = append(checkouts, checkout)
	}

	return checkouts, nil
}
