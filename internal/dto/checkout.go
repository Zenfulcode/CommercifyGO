package dto

import (
	"time"
)

// CheckoutDTO represents a checkout session in the system
type CheckoutDTO struct {
	ID               uint                `json:"id"`
	UserID           uint                `json:"user_id,omitempty"`
	SessionID        string              `json:"session_id,omitempty"`
	Items            []CheckoutItemDTO   `json:"items"`
	Status           string              `json:"status"`
	ShippingAddress  AddressDTO          `json:"shipping_address"`
	BillingAddress   AddressDTO          `json:"billing_address"`
	ShippingMethodID uint                `json:"shipping_method_id,omitempty"`
	ShippingOption   *ShippingOptionDTO  `json:"shipping_option,omitempty"`
	PaymentProvider  string              `json:"payment_provider,omitempty"`
	TotalAmount      float64             `json:"total_amount"`
	ShippingCost     float64             `json:"shipping_cost"`
	TotalWeight      float64             `json:"total_weight"`
	CustomerDetails  CustomerDetailsDTO  `json:"customer_details"`
	Currency         string              `json:"currency"`
	DiscountCode     string              `json:"discount_code,omitempty"`
	DiscountAmount   float64             `json:"discount_amount"`
	FinalAmount      float64             `json:"final_amount"`
	AppliedDiscount  *AppliedDiscountDTO `json:"applied_discount,omitempty"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	LastActivityAt   time.Time           `json:"last_activity_at"`
	ExpiresAt        time.Time           `json:"expires_at"`
	CompletedAt      *time.Time          `json:"completed_at,omitempty"`
	ConvertedOrderID uint                `json:"converted_order_id,omitempty"`
}

// CheckoutItemDTO represents an item in a checkout
type CheckoutItemDTO struct {
	ID          uint      `json:"id"`
	ProductID   uint      `json:"product_id"`
	VariantID   uint      `json:"variant_id,omitempty"`
	ProductName string    `json:"product_name"`
	VariantName string    `json:"variant_name,omitempty"`
	ImageURL    string    `json:"image_url"`
	SKU         string    `json:"sku"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	Weight      float64   `json:"weight"`
	Subtotal    float64   `json:"subtotal"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CardDetailsDTO represents card details for payment processing
type CardDetailsDTO struct {
	CardNumber     string `json:"card_number"`
	ExpiryMonth    int    `json:"expiry_month"`
	ExpiryYear     int    `json:"expiry_year"`
	CVV            string `json:"cvv"`
	CardholderName string `json:"cardholder_name"`
	Token          string `json:"token,omitempty"` // Optional token for saved cards
}
