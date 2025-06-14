package entity

import (
	"errors"
	"time"
)

// ShippingRate connects shipping methods to zones with pricing
type ShippingRate struct {
	ID                    uint              `json:"id"`
	ShippingMethodID      uint              `json:"shipping_method_id"`
	ShippingMethod        *ShippingMethod   `json:"shipping_method,omitempty"`
	ShippingZoneID        uint              `json:"shipping_zone_id"`
	ShippingZone          *ShippingZone     `json:"shipping_zone,omitempty"`
	BaseRate              int64             `json:"base_rate"`
	MinOrderValue         int64             `json:"min_order_value"`
	FreeShippingThreshold *int64            `json:"free_shipping_threshold"`
	WeightBasedRates      []WeightBasedRate `json:"weight_based_rates,omitempty"`
	ValueBasedRates       []ValueBasedRate  `json:"value_based_rates,omitempty"`
	Active                bool              `json:"active"`
	CreatedAt             time.Time         `json:"created_at"`
	UpdatedAt             time.Time         `json:"updated_at"`
}

// WeightBasedRate represents additional costs based on order weight
type WeightBasedRate struct {
	ID             uint      `json:"id"`
	ShippingRateID uint      `json:"shipping_rate_id"`
	MinWeight      float64   `json:"min_weight"`
	MaxWeight      float64   `json:"max_weight"`
	Rate           int64     `json:"rate"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ValueBasedRate represents additional costs/discounts based on order value
type ValueBasedRate struct {
	ID             uint      `json:"id"`
	ShippingRateID uint      `json:"shipping_rate_id"`
	MinOrderValue  int64     `json:"min_order_value"`
	MaxOrderValue  int64     `json:"max_order_value"`
	Rate           int64     `json:"rate"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

// ShippingOption represents a single shipping option with its cost
type ShippingOption struct {
	ShippingRateID        uint   `json:"shipping_rate_id"`
	ShippingMethodID      uint   `json:"shipping_method_id"`
	Name                  string `json:"name"`
	Description           string `json:"description"`
	EstimatedDeliveryDays int    `json:"estimated_delivery_days"`
	Cost                  int64  `json:"cost"`
	FreeShipping          bool   `json:"free_shipping"`
}

// NewShippingRate creates a new shipping rate
func NewShippingRate(
	shippingMethodID uint,
	shippingZoneID uint,
	baseRate,
	minOrderValue int64,
) (*ShippingRate, error) {
	if shippingMethodID == 0 {
		return nil, errors.New("shipping method ID cannot be empty")
	}

	if shippingZoneID == 0 {
		return nil, errors.New("shipping zone ID cannot be empty")
	}

	if baseRate < 0 {
		return nil, errors.New("base rate cannot be negative")
	}

	if minOrderValue < 0 {
		return nil, errors.New("minimum order value cannot be negative")
	}

	now := time.Now()
	return &ShippingRate{
		ShippingMethodID: shippingMethodID,
		ShippingZoneID:   shippingZoneID,
		BaseRate:         baseRate,
		MinOrderValue:    minOrderValue,
		Active:           true,
		CreatedAt:        now,
		UpdatedAt:        now,
	}, nil
}

// Update updates a shipping rate's details
func (r *ShippingRate) Update(baseRate, minOrderValue int64) error {
	if baseRate < 0 {
		return errors.New("base rate cannot be negative")
	}

	if minOrderValue < 0 {
		return errors.New("minimum order value cannot be negative")
	}

	r.BaseRate = baseRate
	r.MinOrderValue = minOrderValue
	r.UpdatedAt = time.Now()
	return nil
}

// SetFreeShippingThreshold sets the free shipping threshold
func (r *ShippingRate) SetFreeShippingThreshold(threshold *int64) {
	// Validate that threshold is either nil or positive
	if threshold != nil && *threshold < 0 {
		return
	}

	r.FreeShippingThreshold = threshold
	r.UpdatedAt = time.Now()
}

// CalculateShippingCost calculates the shipping cost for an order
func (r *ShippingRate) CalculateShippingCost(orderValue int64, weight float64) (int64, error) {
	// Check if order qualifies for free shipping
	if r.FreeShippingThreshold != nil && orderValue >= *r.FreeShippingThreshold {
		return 0, nil // Free shipping applies
	}

	// Check if order meets minimum value
	if orderValue < r.MinOrderValue {
		return 0, errors.New("order value does not meet minimum requirement")
	}

	// Start with the base rate
	cost := r.BaseRate

	// Apply weight-based rates
	for _, wbr := range r.WeightBasedRates {
		if weight >= wbr.MinWeight && weight <= wbr.MaxWeight {
			cost += wbr.Rate
			break // Only apply the first matching weight rate
		}
	}

	// Apply value-based rates
	for _, vbr := range r.ValueBasedRates {
		if orderValue >= vbr.MinOrderValue && orderValue <= vbr.MaxOrderValue {
			cost += vbr.Rate
			break // Only apply the first matching value rate
		}
	}

	return cost, nil
}

// Activate activates a shipping rate
func (r *ShippingRate) Activate() {
	if !r.Active {
		r.Active = true
		r.UpdatedAt = time.Now()
	}
}

// Deactivate deactivates a shipping rate
func (r *ShippingRate) Deactivate() {
	if r.Active {
		r.Active = false
		r.UpdatedAt = time.Now()
	}
}
