package dto

import (
	"time"
)

// OrderDTO represents an order in the system
type OrderDTO struct {
	ID              uint                `json:"id"`
	OrderNumber     string              `json:"order_number"`
	UserID          uint                `json:"user_id"`
	CheckoutID      string              `json:"checkout_id"`
	Items           []OrderItemDTO      `json:"items"`
	Status          OrderStatus         `json:"status"`
	PaymentStatus   PaymentStatus       `json:"payment_status"`
	TotalAmount     float64             `json:"total_amount"`    // Subtotal (items only)
	ShippingCost    float64             `json:"shipping_cost"`   // Shipping cost
	DiscountAmount  float64             `json:"discount_amount"` // Discount applied amount
	FinalAmount     float64             `json:"final_amount"`    // Total including shipping and discounts
	Currency        string              `json:"currency"`
	ShippingAddress AddressDTO          `json:"shipping_address"`
	BillingAddress  AddressDTO          `json:"billing_address"`
	PaymentDetails  *PaymentDetails     `json:"payment_details"`
	ShippingDetails ShippingOptionDTO   `json:"shipping_details"`
	DiscountDetails *AppliedDiscountDTO `json:"discount_details"`
	CustomerDetails CustomerDetailsDTO  `json:"customer"`
	ActionRequired  bool                `json:"action_required"`      // Indicates if action is needed (e.g., payment)
	ActionURL       string              `json:"action_url,omitempty"` // URL for payment or order actions
	CreatedAt       time.Time           `json:"created_at"`
	UpdatedAt       time.Time           `json:"updated_at"`
}

type OrderSummaryDTO struct {
	ID               uint               `json:"id"`
	OrderNumber      string             `json:"order_number"`
	CheckoutID       string             `json:"checkout_id"`
	UserID           uint               `json:"user_id"`
	Customer         CustomerDetailsDTO `json:"customer"`
	Status           OrderStatus        `json:"status"`
	PaymentStatus    PaymentStatus      `json:"payment_status"`
	TotalAmount      float64            `json:"total_amount"`  // Subtotal (items only)
	ShippingCost     float64            `json:"shipping_cost"` // Shipping cost
	FinalAmount      float64            `json:"final_amount"`  // Total including shipping and discounts
	OrderLinesAmount int                `json:"order_lines_amount"`
	Currency         string             `json:"currency"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

type PaymentDetails struct {
	PaymentID string          `json:"payment_id"`
	Provider  PaymentProvider `json:"provider"`
	Method    PaymentMethod   `json:"method"`
	Status    string          `json:"status"`
	Captured  bool            `json:"captured"`
	Refunded  bool            `json:"refunded"`
}

// OrderItemDTO represents an item in an order
type OrderItemDTO struct {
	ID          uint      `json:"id"`
	OrderID     uint      `json:"order_id"`
	ProductID   uint      `json:"product_id"`
	VariantID   uint      `json:"variant_id,omitempty"`
	SKU         string    `json:"sku"`
	ProductName string    `json:"product_name"`
	VariantName string    `json:"variant_name"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
	ImageURL    string    `json:"image_url"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// PaymentMethod represents the payment method used for an order
type PaymentMethod string

const (
	PaymentMethodCard   PaymentMethod = "credit_card"
	PaymentMethodWallet PaymentMethod = "wallet"
)

// PaymentProvider represents the payment provider used for an order
type PaymentProvider string

const (
	PaymentProviderStripe    PaymentProvider = "stripe"
	PaymentProviderMobilePay PaymentProvider = "mobilepay"
)

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusCompleted OrderStatus = "completed"
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
