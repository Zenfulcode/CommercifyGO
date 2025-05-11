package postgres

import (
	"database/sql"
	"encoding/json"
	"errors"
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

	// For guest orders, explicitly set user_id to NULL
	if order.IsGuestOrder {
		// Add guest order fields
		query = `
			INSERT INTO orders (
				user_id, total_amount, status, shipping_address, billing_address,
				payment_id, payment_provider, tracking_code, created_at, updated_at, completed_at, final_amount,
				guest_email, guest_phone, guest_full_name, is_guest_order, shipping_method_id, shipping_cost,
				total_weight
			)
			VALUES (NULL, $1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18)
			RETURNING id
		`

		err2 = tx.QueryRow(
			query,
			order.TotalAmount,
			order.Status,
			shippingAddrJSON,
			billingAddrJSON,
			order.PaymentID,
			order.PaymentProvider,
			order.TrackingCode,
			order.CreatedAt,
			order.UpdatedAt,
			order.CompletedAt,
			order.FinalAmount,
			order.GuestEmail,
			order.GuestPhone,
			order.GuestFullName,
			order.IsGuestOrder,
			order.ShippingMethodID,
			order.ShippingCost,
			order.TotalWeight,
		).Scan(&order.ID)
	} else {
		// Regular user order
		query = `
			INSERT INTO orders (
				user_id, total_amount, status, shipping_address, billing_address,
				payment_id, payment_provider, tracking_code, created_at, updated_at, completed_at, final_amount,
				shipping_method_id, shipping_cost, total_weight
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15)
			RETURNING id
		`

		err2 = tx.QueryRow(
			query,
			order.UserID,
			order.TotalAmount,
			order.Status,
			shippingAddrJSON,
			billingAddrJSON,
			order.PaymentID,
			order.PaymentProvider,
			order.TrackingCode,
			order.CreatedAt,
			order.UpdatedAt,
			order.CompletedAt,
			order.FinalAmount,
			order.ShippingMethodID,
			order.ShippingCost,
			order.TotalWeight,
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
func (r *OrderRepository) GetByID(id uint) (*entity.Order, error) {
	// Get order
	query := `
		SELECT id, order_number, user_id, total_amount, status, shipping_address, billing_address,
			payment_id, payment_provider, tracking_code, created_at, updated_at, completed_at,
			discount_amount, discount_id, discount_code, final_amount, action_url,
			guest_email, guest_phone, guest_full_name, is_guest_order, shipping_method_id, shipping_cost,
			total_weight
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
	var guestEmail, guestPhone, guestFullName sql.NullString
	var isGuestOrder sql.NullBool
	var shippingMethodID sql.NullInt64
	var shippingCost sql.NullInt64
	var totalWeight sql.NullFloat64

	var discountID sql.NullInt64
	var discountCode sql.NullString

	err := r.db.QueryRow(query, id).Scan(
		&order.ID,
		&orderNumber,
		&userID,
		&order.TotalAmount,
		&order.Status,
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
		&guestEmail,
		&guestPhone,
		&guestFullName,
		&isGuestOrder,
		&shippingMethodID,
		&shippingCost,
		&totalWeight,
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
		if guestEmail.Valid {
			order.GuestEmail = guestEmail.String
		}
		if guestPhone.Valid {
			order.GuestPhone = guestPhone.String
		}
		if guestFullName.Valid {
			order.GuestFullName = guestFullName.String
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
		SELECT id, order_id, product_id, quantity, price, subtotal
		FROM order_items
		WHERE order_id = $1
	`

	rows, err := r.db.Query(query, order.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	order.Items = []entity.OrderItem{}
	for rows.Next() {
		item := entity.OrderItem{}
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.Subtotal,
		)
		if err != nil {
			return nil, err
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
		SET status = $1, shipping_address = $2, billing_address = $3,
			payment_id = $4, payment_provider = $5, tracking_code = $6, updated_at = $7, completed_at = $8, order_number = $9,
			final_amount = $10,
			discount_id = $11,
			discount_amount = $12,
			discount_code = $13,
			action_url = $14,
			shipping_method_id = $15,
			shipping_cost = $16,
			total_weight = $17
		WHERE id = $18
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
		order.ID,
	)

	return err
}

// GetByUser retrieves orders for a user
func (r *OrderRepository) GetByUser(userID uint, offset, limit int) ([]*entity.Order, error) {
	query := `
		SELECT id, order_number, user_id, total_amount, status, shipping_address, billing_address,
			payment_id, payment_provider, tracking_code, created_at, updated_at, completed_at,
			guest_email, guest_phone, guest_full_name, is_guest_order
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
		var guestEmail, guestPhone, guestFullName sql.NullString
		var isGuestOrder sql.NullBool

		err := rows.Scan(
			&order.ID,
			&orderNumber,
			&userIDNull,
			&order.TotalAmount,
			&order.Status,
			&shippingAddrJSON,
			&billingAddrJSON,
			&order.PaymentID,
			&paymentProvider,
			&order.TrackingCode,
			&order.CreatedAt,
			&order.UpdatedAt,
			&completedAt,
			&guestEmail,
			&guestPhone,
			&guestFullName,
			&isGuestOrder,
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
			if guestEmail.Valid {
				order.GuestEmail = guestEmail.String
			}
			if guestPhone.Valid {
				order.GuestPhone = guestPhone.String
			}
			if guestFullName.Valid {
				order.GuestFullName = guestFullName.String
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
		SELECT id, order_number, user_id, total_amount, status, shipping_address, billing_address,
			payment_id, payment_provider, tracking_code, created_at, updated_at, completed_at,
			guest_email, guest_phone, guest_full_name, is_guest_order
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
		var shippingAddrJSON, billingAddrJSON []byte
		var completedAt sql.NullTime
		var paymentProvider sql.NullString
		var orderNumber sql.NullString
		var userIDNull sql.NullInt64
		var guestEmail, guestPhone, guestFullName sql.NullString
		var isGuestOrder sql.NullBool

		err := rows.Scan(
			&order.ID,
			&orderNumber,
			&userIDNull,
			&order.TotalAmount,
			&order.Status,
			&shippingAddrJSON,
			&billingAddrJSON,
			&order.PaymentID,
			&paymentProvider,
			&order.TrackingCode,
			&order.CreatedAt,
			&order.UpdatedAt,
			&completedAt,
			&guestEmail,
			&guestPhone,
			&guestFullName,
			&isGuestOrder,
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
			if guestEmail.Valid {
				order.GuestEmail = guestEmail.String
			}
			if guestPhone.Valid {
				order.GuestPhone = guestPhone.String
			}
			if guestFullName.Valid {
				order.GuestFullName = guestFullName.String
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

		// Get order items (simplified to avoid N+1 query issue in production)
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
		SELECT id, order_number, user_id, total_amount, status, shipping_address, billing_address,
			payment_id, payment_provider, tracking_code, created_at, updated_at, completed_at,
			discount_amount, discount_id, discount_code, final_amount, action_url,
			guest_email, guest_phone, guest_full_name, is_guest_order
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
	var guestEmail, guestPhone, guestFullName sql.NullString
	var isGuestOrder sql.NullBool

	var discountID sql.NullInt64
	var discountCode sql.NullString

	err := r.db.QueryRow(query, paymentID).Scan(
		&order.ID,
		&orderNumber,
		&userID,
		&order.TotalAmount,
		&order.Status,
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
		&guestEmail,
		&guestPhone,
		&guestFullName,
		&isGuestOrder,
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
		if guestEmail.Valid {
			order.GuestEmail = guestEmail.String
		}
		if guestPhone.Valid {
			order.GuestPhone = guestPhone.String
		}
		if guestFullName.Valid {
			order.GuestFullName = guestFullName.String
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

	// Get order items
	query = `
		SELECT id, order_id, product_id, quantity, price, subtotal
		FROM order_items
		WHERE order_id = $1
	`

	rows, err := r.db.Query(query, order.ID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	order.Items = []entity.OrderItem{}
	for rows.Next() {
		item := entity.OrderItem{}
		err := rows.Scan(
			&item.ID,
			&item.OrderID,
			&item.ProductID,
			&item.Quantity,
			&item.Price,
			&item.Subtotal,
		)
		if err != nil {
			return nil, err
		}
		order.Items = append(order.Items, item)
	}

	return order, nil
}
