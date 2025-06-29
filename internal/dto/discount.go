package dto

import (
	"time"
)

// DiscountDTO represents a discount in the system
type DiscountDTO struct {
	ID               uint      `json:"id"`
	Code             string    `json:"code"`
	Type             string    `json:"type"`
	Method           string    `json:"method"`
	Value            float64   `json:"value"`
	MinOrderValue    float64   `json:"min_order_value"`
	MaxDiscountValue float64   `json:"max_discount_value"`
	ProductIDs       []uint    `json:"product_ids,omitempty"`
	CategoryIDs      []uint    `json:"category_ids,omitempty"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	UsageLimit       int       `json:"usage_limit"`
	CurrentUsage     int       `json:"current_usage"`
	Active           bool      `json:"active"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// AppliedDiscountDTO represents an applied discount in a checkout
type AppliedDiscountDTO struct {
	ID     uint    `json:"id"`
	Code   string  `json:"code"`
	Type   string  `json:"type"`
	Method string  `json:"method"`
	Value  float64 `json:"value"`
	Amount float64 `json:"amount"`
}
