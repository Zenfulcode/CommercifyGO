package dto

import (
	"time"
)

// ShippingMethodDetailDTO represents a shipping method in the system with full details
type ShippingMethodDetailDTO struct {
	ID                    uint      `json:"id"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	EstimatedDeliveryDays int       `json:"estimated_delivery_days"`
	Active                bool      `json:"active"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// ShippingZoneDTO represents a shipping zone in the system
type ShippingZoneDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Countries   []string  `json:"countries"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ShippingRateDTO represents a shipping rate in the system
type ShippingRateDTO struct {
	ID                    uint                     `json:"id"`
	ShippingMethodID      uint                     `json:"shipping_method_id"`
	ShippingMethod        *ShippingMethodDetailDTO `json:"shipping_method,omitempty"`
	ShippingZoneID        uint                     `json:"shipping_zone_id"`
	ShippingZone          *ShippingZoneDTO         `json:"shipping_zone,omitempty"`
	BaseRate              float64                  `json:"base_rate"`
	MinOrderValue         float64                  `json:"min_order_value"`
	FreeShippingThreshold float64                  `json:"free_shipping_threshold"`
	WeightBasedRates      []WeightBasedRateDTO     `json:"weight_based_rates,omitempty"`
	ValueBasedRates       []ValueBasedRateDTO      `json:"value_based_rates,omitempty"`
	Active                bool                     `json:"active"`
	CreatedAt             time.Time                `json:"created_at"`
	UpdatedAt             time.Time                `json:"updated_at"`
}

// WeightBasedRateDTO represents a weight-based rate in the system
type WeightBasedRateDTO struct {
	ID             uint      `json:"id"`
	ShippingRateID uint      `json:"shipping_rate_id"`
	MinWeight      float64   `json:"min_weight"`
	MaxWeight      float64   `json:"max_weight"`
	Rate           float64   `json:"rate"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ValueBasedRateDTO represents a value-based rate in the system
type ValueBasedRateDTO struct {
	ID             uint      `json:"id"`
	ShippingRateID uint      `json:"shipping_rate_id"`
	MinOrderValue  float64   `json:"min_order_value"`
	MaxOrderValue  float64   `json:"max_order_value"`
	Rate           float64   `json:"rate"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ShippingOptionDTO represents a shipping option with calculated cost
type ShippingOptionDTO struct {
	ShippingRateID        uint    `json:"shipping_rate_id,omitempty"`
	ShippingMethodID      uint    `json:"shipping_method_id,omitempty"`
	Name                  string  `json:"name"`
	Description           string  `json:"description"`
	EstimatedDeliveryDays int     `json:"estimated_delivery_days"`
	Cost                  float64 `json:"cost"`
	FreeShipping          bool    `json:"free_shipping"`
}
