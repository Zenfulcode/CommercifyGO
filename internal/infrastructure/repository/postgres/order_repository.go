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

// OrderRepository implements the order repository interface using PostgreSQL
type OrderRepository struct {
	db *sql.DB
}

// NewOrderRepository creates a new OrderRepository
func NewOrderRepository(db *sql.DB) repository.OrderRepository {
	return &OrderRepository{db: db}
}

// Create creates a new order
func (r *OrderRepository) Create(order *entity.Order) error {
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

	// Marshal addresses to JSON
	shippingAddrJSON, err := json.Marshal(order.ShippingAddr)
	if err != nil {
		return err
	}

	billingAddrJSON, err := json.Marshal(order.BillingAddr)
	if err != nil {
		return err
	}

	// Insert order
	var query string
	var err2 error

	// For guest orders or orders with UserID=0, explicitly set user_id to NULL
	if order.IsGuestOrder || order.UserID == 0 {
		// Add guest order fields
		query = `
			INSERT INTO orders (
				user_id, total_amount, status, payment_status, shipping_address, billing_address,
				payment_id, payment_provider, tracking_code, created_at, updated_at, completed_at, final_amount,
				customer_email, customer_phone, customer_full_name, is_guest_order, shipping_method_id, shipping_cost,
				total_weight, currency
			)
			VALUES (NULL, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
			RETURNING id
		`

		err2 = tx.QueryRow(
			query,
			order.TotalAmount,
			order.Status,
			order.PaymentStatus,
			shippingAddrJSON,
			billingAddrJSON,
			order.PaymentID,
			order.PaymentProvider,
			order.TrackingCode,
			order.CreatedAt,
			order.UpdatedAt,
			order.CompletedAt,
			order.FinalAmount,
			order.CustomerDetails.Email,
			order.CustomerDetails.Phone,
			order.CustomerDetails.FullName,
			order.IsGuestOrder,
			order.ShippingMethodID,
			order.ShippingCost,
			order.TotalWeight,
			order.Currency,
		).Scan(&order.ID)
	} else {
		// Regular user order
		query = `
			INSERT INTO orders (
				user_id, total_amount, status, payment_status, shipping_address, billing_address,
				payment_id, payment_provider, tracking_code, created_at, updated_at, completed_at, final_amount,
				customer_email, customer_phone, customer_full_name, shipping_method_id, shipping_cost, total_weight,
				currency
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20)
			RETURNING id
		`

		err2 = tx.QueryRow(
			query,
			order.UserID,
			order.TotalAmount,
			order.Status,
			order.PaymentStatus,
			shippingAddrJSON,
			billingAddrJSON,
			order.PaymentID,
			order.PaymentProvider,
			order.TrackingCode,
			order.CreatedAt,
			order.UpdatedAt,
			order.CompletedAt,
			order.FinalAmount,
			order.CustomerDetails.Email,
			order.CustomerDetails.Phone,
			order.CustomerDetails.FullName,
			order.ShippingMethodID,
			order.ShippingCost,
			order.TotalWeight,
			order.Currency,
		).Scan(&order.ID)
	}

	if err2 != nil {
		return err2
	}

	// Generate and set the order number
	order.SetOrderNumber(order.ID)

	// Update the order with the generated order number
	_, err = tx.Exec(
		"UPDATE orders SET order_number = $1 WHERE id = $2",
		order.OrderNumber,
		order.ID,
	)
	if err != nil {
		return err
	}

	// Insert order items
	for i := range order.Items {
		order.Items[i].OrderID = order.ID
		query := `
			INSERT INTO order_items (order_id, product_id, quantity, price, subtotal, created_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id
		`
		err = tx.QueryRow(
			query,
			order.Items[i].OrderID,
			order.Items[i].ProductID,
			order.Items[i].Quantity,
			order.Items[i].Price,
			order.Items[i].Subtotal,
			order.CreatedAt,
		).Scan(&order.Items[i].ID)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetByID retrieves an order by ID
func (r *OrderRepository) GetByID(orderID uint) (*entity.Order, error) {
	// Get order
	query := `
		SELECT id, order_number, user_id, total_amount, status, payment_status, shipping_address, billing_address,
			payment_id, payment_provider, tracking_code, created_at, updated_at, completed_at,
			discount_amount, discount_id, discount_code, final_amount, action_url,
			customer_email, customer_phone, customer_full_name, is_guest_order, shipping_method_id, shipping_cost,
			total_weight, currency
		FROM orders
		WHERE id = $1
	`

	order := &entity.Order{}
	var shippingAddrJSON, billingAddrJSON []byte
	var completedAt sql.NullTime
	var paymentProvider sql.NullString
	var orderNumber sql.NullString
	var actionURL sql.NullString
	var userID sql.NullInt64 // Use NullInt64 to handle NULL user_id
	var customerEmail, customerPhone, customerFullName sql.NullString
	var isGuestOrder sql.NullBool
	var shippingMethodID sql.NullInt64
	var shippingCost sql.NullInt64
	var totalWeight sql.NullFloat64

	var discountID sql.NullInt64
	var discountCode sql.NullString

	err := r.db.QueryRow(query, orderID).Scan(
		&order.ID,
		&orderNumber,
		&userID,
		&order.TotalAmount,
		&order.Status,
		&order.PaymentStatus,
		&shippingAddrJSON,
		&billingAddrJSON,
		&order.PaymentID,
		&paymentProvider,
		&order.TrackingCode,
		&order.CreatedAt,
		&order.UpdatedAt,
		&completedAt,
		&order.DiscountAmount,
		&discountID,
		&discountCode,
		&order.FinalAmount,
		&actionURL,
		&customerEmail,
		&customerPhone,
		&customerFullName,
		&isGuestOrder,
		&shippingMethodID,
		&shippingCost,
		&totalWeight,
		&order.Currency,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("order not found")
	}

	if err != nil {
		return nil, err
	}

	// Handle user_id properly
	if userID.Valid {
		order.UserID = uint(userID.Int64)
	} else {
		order.UserID = 0 // Use 0 to represent NULL in our application
	}

	// Handle guest order fields
	if isGuestOrder.Valid && isGuestOrder.Bool {
		order.IsGuestOrder = true
		order.CustomerDetails = &entity.CustomerDetails{
			Email:    customerEmail.String,
			Phone:    customerPhone.String,
			FullName: customerFullName.String,
		}
	}

	order.AppliedDiscount = &entity.AppliedDiscount{
		DiscountID:     uint(discountID.Int64),
		DiscountCode:   discountCode.String,
		DiscountAmount: order.DiscountAmount,
	}

	if order.FinalAmount == 0 {
		order.FinalAmount = order.TotalAmount
	}

	// Set order number if valid
	if orderNumber.Valid {
		order.OrderNumber = orderNumber.String
	}

	// Set payment provider if valid
	if paymentProvider.Valid {
		order.PaymentProvider = paymentProvider.String
	}

	// Set action URL if valid
	if actionURL.Valid {
		order.ActionURL = actionURL.String
	}

	// Unmarshal addresses
	if err := json.Unmarshal(shippingAddrJSON, &order.ShippingAddr); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(billingAddrJSON, &order.BillingAddr); err != nil {
		return nil, err
	}

	// Set completed at if valid
	if completedAt.Valid {
		order.CompletedAt = &completedAt.Time
	}

	// Set shipping method ID if valid
	if shippingMethodID.Valid {
		order.ShippingMethodID = uint(shippingMethodID.Int64)
	}

	// Set shipping cost if valid
	if shippingCost.Valid {
		order.ShippingCost = shippingCost.Int64
	}

	// Set total weight if valid
	if totalWeight.Valid {
		order.TotalWeight = totalWeight.Float64
	}

	// Get order items
	query = `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price, oi.subtotal,
			p.name as product_name, p.product_number as sku
		FROM order_items oi
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id = $1
	`

	rows, err := r.db.Query(query, order.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	order.Items = []entity.OrderItem{}
	for rows.Next() {
		item := entity.OrderItem{}
		var productName, sku sql.NullString
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.Subtotal,
			&productName,
			&sku,
		)
		if err != nil {
			return nil, err
		}
		if productName.Valid {
			item.ProductName = productName.String
		}
		if sku.Valid {
			item.SKU = sku.String
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

// Update updates an order
func (r *OrderRepository) Update(order *entity.Order) error {
	// Marshal addresses to JSON
	shippingAddrJSON, err := json.Marshal(order.ShippingAddr)
	if err != nil {
		return err
	}

	billingAddrJSON, err := json.Marshal(order.BillingAddr)
	if err != nil {
		return err
	}

	// Update order
	query := `
		UPDATE orders
		SET status = $1, payment_status = $2, shipping_address = $3, billing_address = $4,
			payment_id = $5, payment_provider = $6, tracking_code = $7, updated_at = $8, completed_at = $9, order_number = $10,
			final_amount = $11,
			discount_id = $12,
			discount_amount = $13,
			discount_code = $14,
			action_url = $15,
			shipping_method_id = $16,
			shipping_cost = $17,
			total_weight = $18,
			customer_email = $19,
			customer_phone = $20,
			customer_full_name = $21
		WHERE id = $22
	`

	var discountID sql.NullInt64
	var discountCode sql.NullString
	var discountAmount int64 = 0

	if order.AppliedDiscount != nil && order.AppliedDiscount.DiscountID > 0 {
		discountID.Int64 = int64(order.AppliedDiscount.DiscountID)
		discountID.Valid = true
		discountAmount = order.AppliedDiscount.DiscountAmount
		discountCode.String = order.AppliedDiscount.DiscountCode
		discountCode.Valid = true
	}

	_, err = r.db.Exec(
		query,
		order.Status,
		order.PaymentStatus,
		shippingAddrJSON,
		billingAddrJSON,
		order.PaymentID,
		order.PaymentProvider,
		order.TrackingCode,
		time.Now(),
		order.CompletedAt,
		order.OrderNumber,
		order.FinalAmount,
		discountID,
		discountAmount,
		discountCode,
		order.ActionURL,
		order.ShippingMethodID,
		order.ShippingCost,
		order.TotalWeight,
		order.CustomerDetails.Email,
		order.CustomerDetails.Phone,
		order.CustomerDetails.FullName,
		order.ID,
	)

	return err
}

// GetByUser retrieves orders for a user
func (r *OrderRepository) GetByUser(userID uint, offset, limit int) ([]*entity.Order, error) {
	query := `
		SELECT id, order_number, user_id, total_amount, status, payment_status, shipping_address, billing_address,
			payment_id, payment_provider, tracking_code, created_at, updated_at, completed_at,
			customer_email, customer_phone, customer_full_name, is_guest_order, currency
		FROM orders
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*entity.Order{}
	for rows.Next() {
		order := &entity.Order{}
		var shippingAddrJSON, billingAddrJSON []byte
		var completedAt sql.NullTime
		var paymentProvider sql.NullString
		var orderNumber sql.NullString
		var userIDNull sql.NullInt64
		var customerEmail, customerPhone, customerFullName sql.NullString
		var isGuestOrder sql.NullBool

		err := rows.Scan(
			&order.ID,
			&orderNumber,
			&userIDNull,
			&order.TotalAmount,
			&order.Status,
			&order.PaymentStatus,
			&shippingAddrJSON,
			&billingAddrJSON,
			&order.PaymentID,
			&paymentProvider,
			&order.TrackingCode,
			&order.CreatedAt,
			&order.UpdatedAt,
			&completedAt,
			&customerEmail,
			&customerPhone,
			&customerFullName,
			&isGuestOrder,
			&order.Currency,
		)
		if err != nil {
			return nil, err
		}

		// Handle user_id properly
		if userIDNull.Valid {
			order.UserID = uint(userIDNull.Int64)
		} else {
			order.UserID = 0
		}

		// Handle guest order fields
		if isGuestOrder.Valid && isGuestOrder.Bool {
			order.IsGuestOrder = true
			order.CustomerDetails = &entity.CustomerDetails{
				Email:    customerEmail.String,
				Phone:    customerPhone.String,
				FullName: customerFullName.String,
			}
		}

		// Set order number if valid
		if orderNumber.Valid {
			order.OrderNumber = orderNumber.String
		}

		// Set payment provider if valid
		if paymentProvider.Valid {
			order.PaymentProvider = paymentProvider.String
		}

		// Unmarshal addresses
		if err := json.Unmarshal(shippingAddrJSON, &order.ShippingAddr); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(billingAddrJSON, &order.BillingAddr); err != nil {
			return nil, err
		}

		// Set completed at if valid
		if completedAt.Valid {
			order.CompletedAt = &completedAt.Time
		}

		// Get order items
		itemsQuery := `
			SELECT id, order_id, product_id, quantity, price, subtotal
			FROM order_items
			WHERE order_id = $1
		`

		itemRows, err := r.db.Query(itemsQuery, order.ID)
		if err != nil {
			return nil, err
		}

		order.Items = []entity.OrderItem{}
		for itemRows.Next() {
			item := entity.OrderItem{}
			err := itemRows.Scan(
				&item.ID,
				&item.OrderID,
				&item.ProductID,
				&item.Quantity,
				&item.Price,
				&item.Subtotal,
			)
			if err != nil {
				itemRows.Close()
				return nil, err
			}
			order.Items = append(order.Items, item)
		}
		itemRows.Close()

		orders = append(orders, order)
	}

	return orders, nil
}

// ListByStatus retrieves orders by status
func (r *OrderRepository) ListByStatus(status entity.OrderStatus, offset, limit int) ([]*entity.Order, error) {
	query := `
		SELECT id, order_number, user_id, total_amount, status, payment_status, created_at, updated_at, completed_at,
			customer_email, customer_full_name, is_guest_order, currency
		FROM orders
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, string(status), limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*entity.Order{}
	for rows.Next() {
		order := &entity.Order{}
		var completedAt sql.NullTime
		var orderNumber sql.NullString
		var userIDNull sql.NullInt64
		var customerEmail, customerFullName sql.NullString
		var isGuestOrder sql.NullBool

		err := rows.Scan(
			&order.ID,
			&orderNumber,
			&userIDNull,
			&order.TotalAmount,
			&order.Status,
			&order.PaymentStatus,
			&order.CreatedAt,
			&order.UpdatedAt,
			&completedAt,
			&customerEmail,
			&customerFullName,
			&isGuestOrder,
			&order.Currency,
		)
		if err != nil {
			return nil, err
		}

		// Handle user_id properly
		if userIDNull.Valid {
			order.UserID = uint(userIDNull.Int64)
		} else {
			order.UserID = 0
		}

		// Handle guest order fields
		if isGuestOrder.Valid && isGuestOrder.Bool {
			order.IsGuestOrder = true
			order.CustomerDetails = &entity.CustomerDetails{
				Email:    customerEmail.String,
				FullName: customerFullName.String,
			}
		}

		// Set order number if valid
		if orderNumber.Valid {
			order.OrderNumber = orderNumber.String
		}

		// Set completed at if valid
		if completedAt.Valid {
			order.CompletedAt = &completedAt.Time
		}

		// Note: This simplified query doesn't load all order details
		// For full order details, use GetByID method

		orders = append(orders, order)
	}

	return orders, nil
}

// HasOrdersWithProduct checks if a product has any associated orders
func (r *OrderRepository) HasOrdersWithProduct(productID uint) (bool, error) {
	if productID == 0 {
		return false, errors.New("product ID cannot be 0")
	}

	query := `
		SELECT EXISTS(
			SELECT 1 FROM order_items 
			WHERE product_id = $1
		)
	`

	var exists bool
	err := r.db.QueryRow(query, productID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("failed to check if product has orders: %w", err)
	}

	return exists, nil
}

func (r *OrderRepository) IsDiscountIdUsed(discountID uint) (bool, error) {
	query := `
		SELECT COUNT(*) > 0
		FROM orders
		WHERE discount_id = $1
	`

	var exists bool
	err := r.db.QueryRow(query, discountID).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// GetByPaymentID retrieves an order by payment ID
func (r *OrderRepository) GetByPaymentID(paymentID string) (*entity.Order, error) {
	if paymentID == "" {
		return nil, errors.New("payment ID cannot be empty")
	}

	// Get order by payment_id
	query := `
		SELECT id, order_number, user_id, total_amount, status, payment_status, shipping_address, billing_address,
			payment_id, payment_provider, tracking_code, created_at, updated_at, completed_at,
			discount_amount, discount_id, discount_code, final_amount, action_url,
			customer_email, customer_phone, customer_full_name, is_guest_order, shipping_method_id, shipping_cost,
			total_weight, currency
		FROM orders
		WHERE payment_id = $1
	`

	order := &entity.Order{}
	var shippingAddrJSON, billingAddrJSON []byte
	var completedAt sql.NullTime
	var paymentProvider sql.NullString
	var orderNumber sql.NullString
	var actionURL sql.NullString
	var userID sql.NullInt64 // Use NullInt64 to handle NULL user_id
	var customerEmail, customerPhone, customerFullName sql.NullString
	var isGuestOrder sql.NullBool
	var shippingMethodID sql.NullInt64
	var shippingCost sql.NullInt64
	var totalWeight sql.NullFloat64

	var discountID sql.NullInt64
	var discountCode sql.NullString

	err := r.db.QueryRow(query, paymentID).Scan(
		&order.ID,
		&orderNumber,
		&userID,
		&order.TotalAmount,
		&order.Status,
		&order.PaymentStatus,
		&shippingAddrJSON,
		&billingAddrJSON,
		&order.PaymentID,
		&paymentProvider,
		&order.TrackingCode,
		&order.CreatedAt,
		&order.UpdatedAt,
		&completedAt,
		&order.DiscountAmount,
		&discountID,
		&discountCode,
		&order.FinalAmount,
		&actionURL,
		&customerEmail,
		&customerPhone,
		&customerFullName,
		&isGuestOrder,
		&shippingMethodID,
		&shippingCost,
		&totalWeight,
		&order.Currency,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("order not found")
	}

	if err != nil {
		return nil, err
	}

	// Handle user_id properly
	if userID.Valid {
		order.UserID = uint(userID.Int64)
	} else {
		order.UserID = 0 // Use 0 to represent NULL in our application
	}

	// Handle guest order fields
	if isGuestOrder.Valid && isGuestOrder.Bool {
		order.IsGuestOrder = true

	}

	order.CustomerDetails = &entity.CustomerDetails{
		Email:    customerEmail.String,
		Phone:    customerPhone.String,
		FullName: customerFullName.String,
	}

	order.AppliedDiscount = &entity.AppliedDiscount{
		DiscountID:     uint(discountID.Int64),
		DiscountCode:   discountCode.String,
		DiscountAmount: order.DiscountAmount,
	}

	if order.FinalAmount == 0 {
		order.FinalAmount = order.TotalAmount
	}

	// Set order number if valid
	if orderNumber.Valid {
		order.OrderNumber = orderNumber.String
	}

	// Set payment provider if valid
	if paymentProvider.Valid {
		order.PaymentProvider = paymentProvider.String
	}

	// Set action URL if valid
	if actionURL.Valid {
		order.ActionURL = actionURL.String
	}

	// Unmarshal addresses
	if err := json.Unmarshal(shippingAddrJSON, &order.ShippingAddr); err != nil {
		return nil, err
	}

	if err := json.Unmarshal(billingAddrJSON, &order.BillingAddr); err != nil {
		return nil, err
	}

	// Set completed at if valid
	if completedAt.Valid {
		order.CompletedAt = &completedAt.Time
	}

	// Set shipping method ID if valid
	if shippingMethodID.Valid {
		order.ShippingMethodID = uint(shippingMethodID.Int64)
	}

	// Set shipping cost if valid
	if shippingCost.Valid {
		order.ShippingCost = shippingCost.Int64
	}

	// Set total weight if valid
	if totalWeight.Valid {
		order.TotalWeight = totalWeight.Float64
	}

	// Get order items
	query = `
		SELECT oi.id, oi.order_id, oi.product_id, oi.quantity, oi.price, oi.subtotal,
			p.name as product_name, p.product_number as sku
		FROM order_items oi
		LEFT JOIN products p ON p.id = oi.product_id
		WHERE oi.order_id = $1
	`

	rows, err := r.db.Query(query, order.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	order.Items = []entity.OrderItem{}
	for rows.Next() {
		item := entity.OrderItem{}
		var productName, sku sql.NullString
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.Subtotal,
			&productName,
			&sku,
		)
		if err != nil {
			return nil, err
		}
		if productName.Valid {
			item.ProductName = productName.String
		}
		if sku.Valid {
			item.SKU = sku.String
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}

// ListAll lists all orders
func (r *OrderRepository) ListAll(offset, limit int) ([]*entity.Order, error) {
	query := `
		SELECT id, order_number, user_id, total_amount, status, payment_status,
			payment_provider, created_at, updated_at, completed_at,
			final_amount, customer_email, customer_full_name, is_guest_order, currency
		FROM orders
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	orders := []*entity.Order{}
	for rows.Next() {
		order := &entity.Order{}
		var completedAt sql.NullTime
		var userID sql.NullInt64
		var guestEmail, guestFullName sql.NullString
		var isGuestOrder sql.NullBool

		err := rows.Scan(
			&order.ID,
			&order.OrderNumber,
			&userID,
			&order.TotalAmount,
			&order.Status,
			&order.PaymentStatus,
			&order.PaymentProvider,
			&order.CreatedAt,
			&order.UpdatedAt,
			&completedAt,
			&order.FinalAmount,
			&guestEmail,
			&guestFullName,
			&isGuestOrder,
			&order.Currency,
		)

		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, nil
}
