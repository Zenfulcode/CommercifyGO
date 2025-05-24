package dto

import (
	"time"

	"github.com/zenfulcode/commercify/internal/domain/service"
)

// OrderDTO represents an order in the system
type OrderDTO struct {
	ID              uint            `json:"id"`
	UserID          uint            `json:"user_id"`
	OrderNumber     string          `json:"order_number"`
	Items           []OrderItemDTO  `json:"items"`
	Status          OrderStatus     `json:"status"`
	TotalAmount     float64         `json:"total_amount"`
	FinalAmount     float64         `json:"final_amount"`
	Currency        string          `json:"currency"`
	ShippingAddress AddressDTO      `json:"shipping_address"`
	BillingAddress  AddressDTO      `json:"billing_address"`
	PaymentDetails  PaymentDetails  `json:"payment_details"`
	ShippingDetails ShippingDetails `json:"shipping_details"`
	DiscountDetails DiscountDetails `json:"discount_details"`
	Customer        CustomerDetails `json:"customer"`
	CheckoutID      string          `json:"checkout_id,omitempty"`
	CreatedAt       time.Time       `json:"created_at"`
	UpdatedAt       time.Time       `json:"updated_at"`
}

type PaymentDetails struct {
	PaymentID string          `json:"payment_id"`
	Provider  PaymentProvider `json:"provider"`
	Method    PaymentMethod   `json:"method"`
	Status    string          `json:"status"`
	Captured  bool            `json:"captured"`
	Refunded  bool            `json:"refunded"`
}

type ShippingDetails struct {
	MethodID uint    `json:"method_id"`
	Method   string  `json:"method"`
	Cost     float64 `json:"cost"`
}

type CustomerDetails struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
}

type DiscountDetails struct {
	Code   string  `json:"code"`
	Amount float64 `json:"amount"`
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
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateOrderRequest represents the data needed to create a new order
type CreateOrderRequest struct {
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Email            string     `json:"email"`
	PhoneNumber      string     `json:"phone_number,omitempty"`
	ShippingAddress  AddressDTO `json:"shipping_address"`
	BillingAddress   AddressDTO `json:"billing_address"`
	ShippingMethodID uint       `json:"shipping_method_id"`
}

// CreateOrderItemRequest represents the data needed to create a new order item
type CreateOrderItemRequest struct {
	ProductID uint `json:"product_id"`
	VariantID uint `json:"variant_id,omitempty"`
	Quantity  int  `json:"quantity"`
}

// UpdateOrderRequest represents the data needed to update an existing order
type UpdateOrderRequest struct {
	Status            string     `json:"status,omitempty"`
	PaymentStatus     string     `json:"payment_status,omitempty"`
	TrackingNumber    string     `json:"tracking_number,omitempty"`
	EstimatedDelivery *time.Time `json:"estimated_delivery,omitempty"`
}

// OrderListResponse represents a paginated list of orders
type OrderListResponse struct {
	ListResponseDTO[OrderDTO]
}

// OrderSearchRequest represents the parameters for searching orders
type OrderSearchRequest struct {
	UserID        uint        `json:"user_id,omitempty"`
	Status        OrderStatus `json:"status,omitempty"`
	PaymentStatus string      `json:"payment_status,omitempty"`
	StartDate     *time.Time  `json:"start_date,omitempty"`
	EndDate       *time.Time  `json:"end_date,omitempty"`
	PaginationDTO
}

// ProcessPaymentRequest represents the data needed to process a payment
type ProcessPaymentRequest struct {
	PaymentMethod   PaymentMethod        `json:"payment_method"`
	PaymentProvider PaymentProvider      `json:"payment_provider"`
	CardDetails     *service.CardDetails `json:"card_details,omitempty"`
	PhoneNumber     string               `json:"phone_number,omitempty"`
}

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending       OrderStatus = "pending"
	OrderStatusPendingAction OrderStatus = "pending_action" // Requires user action (e.g., redirect to payment provider)
	OrderStatusPaid          OrderStatus = "paid"
	OrderStatusCaptured      OrderStatus = "captured" // Payment captured
	OrderStatusShipped       OrderStatus = "shipped"
	OrderStatusDelivered     OrderStatus = "delivered"
	OrderStatusCancelled     OrderStatus = "cancelled"
	OrderStatusRefunded      OrderStatus = "refunded"
)

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
