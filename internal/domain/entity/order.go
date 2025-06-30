package entity

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/dto"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"gorm.io/gorm"
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
	gorm.Model
	OrderNumber         string         `gorm:"uniqueIndex;not null;size:100"`
	Currency            string         `gorm:"not null;size:3"` // e.g., "USD", "EUR"
	UserID              uint           `gorm:"index"`           // 0 for guest orders
	User                *User          `gorm:"foreignKey:UserID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Items               []OrderItem    `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	TotalAmount         int64          `gorm:"not null"` // stored in cents
	Status              OrderStatus    `gorm:"not null;size:50;default:'pending'"`
	PaymentStatus       PaymentStatus  `gorm:"not null;size:50;default:'pending'"` // New field for payment status
	ShippingAddressJSON *string        `gorm:"column:shipping_address_json;type:text"`
	BillingAddressJSON  *string        `gorm:"column:billing_address_json;type:text"`
	ShippingOptionJSON  *string        `gorm:"column:shipping_option_json;type:text"`
	AppliedDiscountJSON *string        `gorm:"column:applied_discount_json;type:text"`
	PaymentID           string         `gorm:"size:255"`
	PaymentProvider     string         `gorm:"size:100"`
	PaymentMethod       string         `gorm:"size:100"`
	TrackingCode        sql.NullString `gorm:"size:255"`
	ActionURL           sql.NullString `gorm:"size:500"` // URL for redirect to payment provider
	CompletedAt         *time.Time
	CheckoutSessionID   string `gorm:"size:255"` // Tracks which checkout session created this order

	// Guest information (only used for guest orders where UserID is 0)
	CustomerDetails *CustomerDetails `gorm:"embedded;embeddedPrefix:customer_"`
	IsGuestOrder    bool             `gorm:"default:false"`

	// Shipping information stored as JSON
	ShippingCost int64
	TotalWeight  float64

	// Discount-related fields
	DiscountAmount int64
	FinalAmount    int64 `gorm:"not null"` // stored in cents

	// Payment transactions
	PaymentTransactions []PaymentTransaction `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
}

// OrderItem represents an item in an order
type OrderItem struct {
	gorm.Model
	OrderID          uint           `gorm:"index;not null"`
	Order            Order          `gorm:"foreignKey:OrderID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	ProductID        uint           `gorm:"index;not null"`
	Product          Product        `gorm:"foreignKey:ProductID;constraint:OnDelete:RESTRICT,OnUpdate:CASCADE"`
	ProductVariantID uint           `gorm:"index;not null"`
	ProductVariant   ProductVariant `gorm:"foreignKey:ProductVariantID;constraint:OnDelete:RESTRICT,OnUpdate:CASCADE"`
	Quantity         int            `gorm:"not null"`
	Price            int64          `gorm:"not null"` // Price at time of order
	Subtotal         int64          `gorm:"not null"`
	Weight           float64        `gorm:"default:0"`

	// Snapshot data at time of order
	ProductName string `gorm:"not null;size:255"`
	SKU         string `gorm:"not null;size:100"`
	ImageURL    string `gorm:"size:500"`
}

// Address represents a shipping or billing address
type Address struct {
	Street1    string `gorm:"size:255"`
	Street2    string `gorm:"size:255"`
	City       string `gorm:"size:100"`
	State      string `gorm:"size:100"` // Nullable for international addresses
	PostalCode string `gorm:"size:20"`
	Country    string `gorm:"size:100"`
}

type CustomerDetails struct {
	Email    string `gorm:"size:255"`
	Phone    string `gorm:"size:50"`
	FullName string `gorm:"size:200"`
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

	order := &Order{
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
		CustomerDetails: &customerDetails,
		IsGuestOrder:    false,
	}

	// Set addresses using JSON helper methods
	order.SetShippingAddressJSON(&shippingAddr)
	order.SetBillingAddressJSON(&billingAddr)

	return order, nil
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

	order := &Order{
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

		// Guest-specific information
		CustomerDetails: &customerDetails,
		IsGuestOrder:    true,
	}

	// Set addresses using JSON helper methods
	order.SetShippingAddressJSON(&shippingAddr)
	order.SetBillingAddressJSON(&billingAddr)

	return order, nil
}

// UpdateStatus updates the order status
func (o *Order) UpdateStatus(status OrderStatus) error {
	if !isValidStatusTransition(OrderStatus(o.Status), status) {
		return errors.New("invalid status transition: " + string(OrderStatus(o.Status)) + " -> " + string(status))
	}

	o.Status = status

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
	return nil
}

// SetPaymentProvider sets the payment provider for the order
func (o *Order) SetPaymentProvider(provider string) error {
	if provider == "" {
		return errors.New("payment provider cannot be empty")
	}

	o.PaymentProvider = provider

	return nil
}

// SetPaymentMethod sets the payment method for the order
func (o *Order) SetPaymentMethod(method string) error {
	if method == "" {
		return errors.New("payment method cannot be empty")
	}

	o.PaymentMethod = method

	return nil
}

// SetTrackingCode sets the tracking code for the order
func (o *Order) SetTrackingCode(trackingCode string) error {
	if trackingCode == "" {
		return errors.New("tracking code cannot be empty")
	}

	o.TrackingCode = sql.NullString{
		String: trackingCode,
		Valid:  true,
	}

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

	// Record the applied discount using JSON storage
	appliedDiscount := &AppliedDiscount{
		DiscountID:     discount.ID,
		DiscountCode:   discount.Code,
		DiscountAmount: discountAmount,
	}
	o.SetAppliedDiscountJSON(appliedDiscount)

	return nil
}

// RemoveDiscount removes any applied discount from the order
func (o *Order) RemoveDiscount() {
	o.DiscountAmount = 0
	o.FinalAmount = o.TotalAmount + o.ShippingCost
	o.SetAppliedDiscountJSON(nil)
}

// SetActionURL sets the action URL for the order
func (o *Order) SetActionURL(actionURL string) error {
	if actionURL == "" {
		return errors.New("action URL cannot be empty")
	}

	o.ActionURL = sql.NullString{
		String: actionURL,
		Valid:  true,
	}

	return nil
}

// SetShippingMethod sets the shipping method for the order and updates shipping cost
func (o *Order) SetShippingMethod(option *ShippingOption) error {
	if option == nil {
		return errors.New("shipping method cannot be nil")
	}

	o.SetShippingOptionJSON(option)
	o.ShippingCost = option.Cost

	// Update final amount with new shipping cost
	o.FinalAmount = o.TotalAmount + o.ShippingCost - o.DiscountAmount

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

func (o *Order) ToOrderSummaryDTO() *dto.OrderSummaryDTO {
	return &dto.OrderSummaryDTO{
		ID:               o.ID,
		OrderNumber:      o.OrderNumber,
		CheckoutID:       o.CheckoutSessionID,
		UserID:           o.UserID,
		Customer:         *o.CustomerDetails.ToCustomerDetailsDTO(),
		Status:           dto.OrderStatus(o.Status),
		PaymentStatus:    dto.PaymentStatus(o.PaymentStatus),
		TotalAmount:      money.FromCents(o.TotalAmount),
		OrderLinesAmount: len(o.Items),
		Currency:         o.Currency,
		CreatedAt:        o.CreatedAt,
		UpdatedAt:        o.UpdatedAt,
	}
}
func (o *Order) ToOrderDetailsDTO() *dto.OrderDTO {
	var discountDetails *dto.AppliedDiscountDTO
	if appliedDiscount := o.GetAppliedDiscount(); appliedDiscount != nil {
		discountDetails = appliedDiscount.ToAppliedDiscountDTO()
	}

	var shippingDetails *dto.ShippingOptionDTO
	if shippingOption := o.GetShippingOption(); shippingOption != nil {
		shippingDetails = shippingOption.ToShippingOptionDTO()
	}

	shippingAddr := o.GetShippingAddress()
	billingAddr := o.GetBillingAddress()

	return &dto.OrderDTO{
		ID:              o.ID,
		OrderNumber:     o.OrderNumber,
		UserID:          o.UserID,
		CheckoutID:      o.CheckoutSessionID,
		CustomerDetails: o.CustomerDetails.ToCustomerDetailsDTO(),
		ShippingDetails: shippingDetails,
		DiscountDetails: discountDetails,
		Status:          dto.OrderStatus(o.Status),
		PaymentStatus:   dto.PaymentStatus(o.PaymentStatus),
		Currency:        o.Currency,
		TotalAmount:     money.FromCents(o.TotalAmount),
		ShippingCost:    money.FromCents(o.ShippingCost),
		DiscountAmount:  money.FromCents(o.DiscountAmount),
		FinalAmount:     money.FromCents(o.FinalAmount),
		Items:           o.ToOrderItemsDTO(),
		ShippingAddress: shippingAddr.ToAddressDTO(),
		BillingAddress:  billingAddr.ToAddressDTO(),
		ActionRequired:  o.ActionURL.Valid && o.ActionURL.String != "",
		ActionURL:       o.ActionURL.String,
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}
}

func (o *Order) ToOrderItemsDTO() []dto.OrderItemDTO {
	itemsDTO := make([]dto.OrderItemDTO, len(o.Items))
	for i, item := range o.Items {
		itemsDTO[i] = dto.OrderItemDTO{
			ID:          item.ID,
			OrderID:     item.OrderID,
			ProductID:   item.ProductID,
			VariantID:   item.ProductVariantID,
			SKU:         item.SKU,
			ProductName: item.ProductName,
			VariantName: item.ProductVariant.Name(),
			ImageURL:    item.ImageURL,
			Quantity:    item.Quantity,
			UnitPrice:   money.FromCents(item.Price),
			TotalPrice:  money.FromCents(item.Subtotal),
		}
	}
	return itemsDTO
}

func (a *Address) ToAddressDTO() *dto.AddressDTO {
	return &dto.AddressDTO{
		AddressLine1: a.Street1,
		AddressLine2: a.Street2,
		City:         a.City,
		State:        a.State,
		PostalCode:   a.PostalCode,
		Country:      a.Country,
	}
}

func (c *CustomerDetails) ToCustomerDetailsDTO() *dto.CustomerDetailsDTO {
	return &dto.CustomerDetailsDTO{
		Email:    c.Email,
		Phone:    c.Phone,
		FullName: c.FullName,
	}
}

// GetShippingAddress returns the shipping address from JSON
func (o *Order) GetShippingAddress() Address {
	if o.ShippingAddressJSON == nil || *o.ShippingAddressJSON == "" {
		return Address{}
	}

	var address Address
	if err := json.Unmarshal([]byte(*o.ShippingAddressJSON), &address); err != nil {
		return Address{}
	}
	return address
}

// SetShippingAddressJSON sets the shipping address as JSON
func (o *Order) SetShippingAddressJSON(address *Address) {
	if address == nil {
		o.ShippingAddressJSON = nil
		return
	}

	jsonData, err := json.Marshal(address)
	if err != nil {
		o.ShippingAddressJSON = nil
		return
	}

	jsonStr := string(jsonData)
	o.ShippingAddressJSON = &jsonStr
}

// GetBillingAddress returns the billing address from JSON
func (o *Order) GetBillingAddress() Address {
	if o.BillingAddressJSON == nil || *o.BillingAddressJSON == "" {
		return Address{}
	}

	var address Address
	if err := json.Unmarshal([]byte(*o.BillingAddressJSON), &address); err != nil {
		return Address{}
	}
	return address
}

// SetBillingAddressJSON sets the billing address as JSON
func (o *Order) SetBillingAddressJSON(address *Address) {
	if address == nil {
		o.BillingAddressJSON = nil
		return
	}

	jsonData, err := json.Marshal(address)
	if err != nil {
		o.BillingAddressJSON = nil
		return
	}

	jsonStr := string(jsonData)
	o.BillingAddressJSON = &jsonStr
}

// GetAppliedDiscount returns the applied discount from JSON
func (o *Order) GetAppliedDiscount() *AppliedDiscount {
	if o.AppliedDiscountJSON == nil || *o.AppliedDiscountJSON == "" {
		return nil
	}

	var discount AppliedDiscount
	if err := json.Unmarshal([]byte(*o.AppliedDiscountJSON), &discount); err != nil {
		return nil
	}
	return &discount
}

// SetAppliedDiscountJSON sets the applied discount as JSON
func (o *Order) SetAppliedDiscountJSON(discount *AppliedDiscount) {
	if discount == nil {
		o.AppliedDiscountJSON = nil
		return
	}

	jsonData, err := json.Marshal(discount)
	if err != nil {
		o.AppliedDiscountJSON = nil
		return
	}

	jsonStr := string(jsonData)
	o.AppliedDiscountJSON = &jsonStr
}

// GetShippingOption returns the shipping option from JSON
func (o *Order) GetShippingOption() *ShippingOption {
	if o.ShippingOptionJSON == nil || *o.ShippingOptionJSON == "" {
		return nil
	}

	var option ShippingOption
	if err := json.Unmarshal([]byte(*o.ShippingOptionJSON), &option); err != nil {
		return nil
	}
	return &option
}

// SetShippingOptionJSON sets the shipping option as JSON
func (o *Order) SetShippingOptionJSON(option *ShippingOption) {
	if option == nil {
		o.ShippingOptionJSON = nil
		return
	}

	jsonData, err := json.Marshal(option)
	if err != nil {
		o.ShippingOptionJSON = nil
		return
	}

	jsonStr := string(jsonData)
	o.ShippingOptionJSON = &jsonStr
}
