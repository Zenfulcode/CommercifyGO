package entity

import (
	"errors"
	"fmt"
	"slices"
	"time"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusCompleted OrderStatus = "completed" // Set automatically when payment is captured
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusAuthorized PaymentStatus = "authorized"
	PaymentStatusCaptured   PaymentStatus = "captured"
	PaymentStatusRefunded   PaymentStatus = "refunded"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
	PaymentStatusFailed     PaymentStatus = "failed"
)

// Order represents an order entity
type Order struct {
	ID                uint
	OrderNumber       string
	Currency          string // e.g., "USD", "EUR"
	UserID            uint   // 0 for guest orders
	Items             []OrderItem
	TotalAmount       int64 // stored in cents
	Status            OrderStatus
	PaymentStatus     PaymentStatus // New field for payment status
	ShippingAddr      Address
	BillingAddr       Address
	PaymentID         string
	PaymentProvider   string
	PaymentMethod     string
	TrackingCode      string
	ActionURL         string // URL for redirect to payment provider
	CreatedAt         time.Time
	UpdatedAt         time.Time
	CompletedAt       *time.Time
	CheckoutSessionID string // Tracks which checkout session created this order

	// Guest information (only used for guest orders where UserID is 0)
	CustomerDetails *CustomerDetails `json:"customer_details"`
	IsGuestOrder    bool             `json:"is_guest_order"`

	// Shipping information
	ShippingMethodID uint            `json:"shipping_method_id,omitempty"`
	ShippingOption   *ShippingOption `json:"shipping_option,omitempty"`
	ShippingCost     int64           `json:"shipping_cost"` // stored in cents
	TotalWeight      float64         `json:"total_weight"`

	// Discount-related fields
	DiscountAmount  int64 // stored in cents
	FinalAmount     int64 // stored in cents
	AppliedDiscount *AppliedDiscount
}

// OrderItem represents an item in an order
type OrderItem struct {
	ID               uint    `json:"id"`
	OrderID          uint    `json:"order_id"`
	ProductID        uint    `json:"product_id"`
	ProductVariantID uint    `json:"product_variant_id,omitempty"`
	Quantity         int     `json:"quantity"`
	Price            int64   `json:"price"`    // stored in cents
	Subtotal         int64   `json:"subtotal"` // stored in cents
	Weight           float64 `json:"weight"`   // Weight per item

	ProductName string `json:"product_name"`
	SKU         string `json:"sku"`
	ImageURL    string `json:"image_url,omitempty"`
}

// Address represents a shipping or billing address
type Address struct {
	Street     string `json:"street"`
	City       string `json:"city"`
	State      string `json:"state"`
	PostalCode string `json:"postal_code"`
	Country    string `json:"country"`
}

type CustomerDetails struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
}

// NewOrder creates a new order
func NewOrder(userID uint, items []OrderItem, currency string, shippingAddr, billingAddr Address, customerDetails CustomerDetails) (*Order, error) {
	if userID == 0 {
		return nil, errors.New("user ID cannot be empty")
	}
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}
	if currency == "" {
		return nil, errors.New("currency cannot be empty")
	}

	var totalAmount int64
	totalWeight := 0.0
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, errors.New("item quantity must be greater than zero")
		}
		if item.Price <= 0 {
			return nil, errors.New("item price must be greater than zero")
		}
		item.Subtotal = int64(item.Quantity) * item.Price
		totalAmount += item.Subtotal
		totalWeight += item.Weight * float64(item.Quantity)
	}

	now := time.Now()

	// Generate a friendly order number (will be replaced with actual ID after creation)
	// Format: ORD-YYYYMMDD-TEMP
	orderNumber := fmt.Sprintf("ORD-%s-TEMP", now.Format("20060102"))

	return &Order{
		UserID:          userID,
		OrderNumber:     orderNumber,
		Currency:        currency,
		Items:           items,
		TotalAmount:     totalAmount,
		TotalWeight:     totalWeight,
		ShippingCost:    0, // Default to 0, will be set later
		DiscountAmount:  0,
		FinalAmount:     totalAmount, // Initially same as total amount
		Status:          OrderStatusPending,
		PaymentStatus:   PaymentStatusPending, // Initialize payment status
		ShippingAddr:    shippingAddr,
		BillingAddr:     billingAddr,
		CreatedAt:       now,
		UpdatedAt:       now,
		CustomerDetails: &customerDetails,
		IsGuestOrder:    false,
	}, nil
}

// NewGuestOrder creates a new order for a guest user
func NewGuestOrder(items []OrderItem, shippingAddr, billingAddr Address, customerDetails CustomerDetails) (*Order, error) {
	if len(items) == 0 {
		return nil, errors.New("order must have at least one item")
	}

	totalAmount := int64(0)
	totalWeight := 0.0
	for _, item := range items {
		if item.Quantity <= 0 {
			return nil, errors.New("item quantity must be greater than zero")
		}
		if item.Price <= 0 {
			return nil, errors.New("item price must be greater than zero")
		}
		item.Subtotal = int64(item.Quantity) * item.Price
		totalAmount += item.Subtotal
		totalWeight += item.Weight * float64(item.Quantity)
	}

	now := time.Now()

	// Format: GS-YYYYMMDD-TEMP (GS prefix for guest orders)
	orderNumber := fmt.Sprintf("GS-%s-TEMP", now.Format("20060102"))

	return &Order{
		UserID:         0, // Using 0 to indicate it should be NULL in the database
		OrderNumber:    orderNumber,
		Items:          items,
		TotalAmount:    totalAmount,
		TotalWeight:    totalWeight,
		ShippingCost:   0, // Default to 0, will be set later
		DiscountAmount: 0,
		FinalAmount:    totalAmount, // Initially same as total amount
		Status:         OrderStatusPending,
		PaymentStatus:  PaymentStatusPending, // Initialize payment status
		ShippingAddr:   shippingAddr,
		BillingAddr:    billingAddr,
		CreatedAt:      now,
		UpdatedAt:      now,

		// Guest-specific information
		CustomerDetails: &customerDetails,
		IsGuestOrder:    true,
	}, nil
}

// UpdateStatus updates the order status
func (o *Order) UpdateStatus(status OrderStatus) error {
	if !isValidStatusTransition(OrderStatus(o.Status), status) {
		return errors.New("invalid status transition: " + string(OrderStatus(o.Status)) + " -> " + string(status))
	}

	o.Status = status
	o.UpdatedAt = time.Now()

	// If the status is cancelled or completed, set the completed_at timestamp
	if status == OrderStatusCancelled || status == OrderStatusCompleted {
		now := time.Now()
		o.CompletedAt = &now
	}

	return nil
}

// isValidStatusTransition checks if a status transition is valid
func isValidStatusTransition(from, to OrderStatus) bool {
	validTransitions := map[OrderStatus][]OrderStatus{
		OrderStatusPending:   {OrderStatusPaid, OrderStatusCancelled},
		OrderStatusPaid:      {OrderStatusShipped, OrderStatusCancelled},
		OrderStatusShipped:   {OrderStatusCompleted, OrderStatusCancelled},
		OrderStatusCancelled: {},
		OrderStatusCompleted: {},
	}

	return slices.Contains(validTransitions[from], to)
}

// SetPaymentID sets the payment ID for the order
func (o *Order) SetPaymentID(paymentID string) error {
	if paymentID == "" {
		return errors.New("payment ID cannot be empty")
	}

	o.PaymentID = paymentID
	o.UpdatedAt = time.Now()
	return nil
}

// SetPaymentProvider sets the payment provider for the order
func (o *Order) SetPaymentProvider(provider string) error {
	if provider == "" {
		return errors.New("payment provider cannot be empty")
	}

	o.PaymentProvider = provider
	o.UpdatedAt = time.Now()
	return nil
}

// SetPaymentMethod sets the payment method for the order
func (o *Order) SetPaymentMethod(method string) error {
	if method == "" {
		return errors.New("payment method cannot be empty")
	}

	o.PaymentMethod = method
	o.UpdatedAt = time.Now()
	return nil
}

// SetTrackingCode sets the tracking code for the order
func (o *Order) SetTrackingCode(trackingCode string) error {
	if trackingCode == "" {
		return errors.New("tracking code cannot be empty")
	}

	o.TrackingCode = trackingCode
	o.UpdatedAt = time.Now()
	return nil
}

// SetOrderNumber sets the order number
func (o *Order) SetOrderNumber(id uint) {
	// Format: ORD-YYYYMMDD-000001
	o.OrderNumber = fmt.Sprintf("ORD-%s-%06d", o.CreatedAt.Format("20060102"), id)
}

// ApplyDiscount applies a discount to the order
func (o *Order) ApplyDiscount(discount *Discount) error {
	if discount == nil {
		return errors.New("discount cannot be nil")
	}

	// Validate discount
	if !discount.IsValid() || !discount.Active {
		return errors.New("discount is invalid or inactive")
	}

	// Use the Discount entity's CalculateDiscount method to calculate the discount amount
	discountAmount := discount.CalculateDiscount(o)
	if discountAmount <= 0 {
		return errors.New("discount is not applicable to this order")
	}

	// Apply the calculated discount
	o.DiscountAmount = discountAmount
	o.FinalAmount = o.TotalAmount + o.ShippingCost - discountAmount

	// Record the applied discount
	o.AppliedDiscount = &AppliedDiscount{
		DiscountID:     discount.ID,
		DiscountCode:   discount.Code,
		DiscountAmount: discountAmount,
	}

	o.UpdatedAt = time.Now()
	return nil
}

// RemoveDiscount removes any applied discount from the order
func (o *Order) RemoveDiscount() {
	o.DiscountAmount = 0
	o.FinalAmount = o.TotalAmount + o.ShippingCost
	o.AppliedDiscount = nil
	o.UpdatedAt = time.Now()
}

// SetActionURL sets the action URL for the order
func (o *Order) SetActionURL(actionURL string) error {
	if actionURL == "" {
		return errors.New("action URL cannot be empty")
	}

	o.ActionURL = actionURL
	o.UpdatedAt = time.Now()
	return nil
}

// SetShippingMethod sets the shipping method for the order and updates shipping cost
func (o *Order) SetShippingMethod(option *ShippingOption) error {
	if option == nil {
		return errors.New("shipping method cannot be nil")
	}

	o.ShippingMethodID = option.ShippingMethodID
	o.ShippingOption = option
	o.ShippingCost = option.Cost

	// Update final amount with new shipping cost
	o.FinalAmount = o.TotalAmount + o.ShippingCost - o.DiscountAmount

	o.UpdatedAt = time.Now()
	return nil
}

// CalculateTotalWeight calculates the total weight of all items in the order
func (o *Order) CalculateTotalWeight() float64 {
	totalWeight := 0.0
	for _, item := range o.Items {
		totalWeight += item.Weight * float64(item.Quantity)
	}
	o.TotalWeight = totalWeight
	return totalWeight
}

// IsCaptured returns true if the payment is captured
func (o *Order) IsCaptured() bool {
	return o.PaymentStatus == PaymentStatusCaptured
}

// IsRefunded returns true if the payment is refunded
func (o *Order) IsRefunded() bool {
	return o.PaymentStatus == PaymentStatusRefunded
}

// UpdatePaymentStatus updates the payment status and handles order status transitions
func (o *Order) UpdatePaymentStatus(status PaymentStatus) error {
	if !isValidPaymentStatusTransition(o.PaymentStatus, status) {
		return errors.New("invalid payment status transition: " + string(o.PaymentStatus) + " -> " + string(status))
	}

	o.PaymentStatus = status
	o.UpdatedAt = time.Now()

	// Handle automatic order status transitions based on payment status
	switch status {
	case PaymentStatusAuthorized:
		// When payment is authorized, order becomes "paid"
		if o.Status == OrderStatusPending {
			o.Status = OrderStatusPaid
		}
	case PaymentStatusFailed:
		// When payment fails, order is cancelled
		if o.Status == OrderStatusPending {
			o.Status = OrderStatusCancelled
			now := time.Now()
			o.CompletedAt = &now
		}
	case PaymentStatusCaptured:
		// When payment is captured and order is shipped, order is completed
		if o.Status == OrderStatusShipped {
			o.Status = OrderStatusCompleted
			now := time.Now()
			o.CompletedAt = &now
		}
	case PaymentStatusCancelled:
		// When payment is cancelled, order is cancelled
		if o.Status == OrderStatusPending || o.Status == OrderStatusPaid {
			o.Status = OrderStatusCancelled
			now := time.Now()
			o.CompletedAt = &now
		}
	}

	return nil
}

// isValidPaymentStatusTransition checks if a payment status transition is valid
func isValidPaymentStatusTransition(from, to PaymentStatus) bool {
	validTransitions := map[PaymentStatus][]PaymentStatus{
		PaymentStatusPending:    {PaymentStatusAuthorized, PaymentStatusFailed},
		PaymentStatusAuthorized: {PaymentStatusCaptured, PaymentStatusRefunded, PaymentStatusCancelled},
		PaymentStatusCaptured:   {PaymentStatusRefunded},
		PaymentStatusRefunded:   {},
		PaymentStatusCancelled:  {},
		PaymentStatusFailed:     {},
	}

	return slices.Contains(validTransitions[from], to)
}
