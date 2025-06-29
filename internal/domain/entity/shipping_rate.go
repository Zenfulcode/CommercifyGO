package entity

import (
	"errors"

	"gorm.io/gorm"
)

// ShippingRate connects shipping methods to zones with pricing
type ShippingRate struct {
	gorm.Model
	ShippingMethodID      uint              `gorm:"index;not null"`
	ShippingMethod        *ShippingMethod   `gorm:"foreignKey:ShippingMethodID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	ShippingZoneID        uint              `gorm:"index;not null"`
	ShippingZone          *ShippingZone     `gorm:"foreignKey:ShippingZoneID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	BaseRate              int64             `gorm:"not null;default:0"`
	MinOrderValue         int64             `gorm:"default:0"`
	FreeShippingThreshold *int64            `gorm:"default:null"`
	WeightBasedRates      []WeightBasedRate `gorm:"foreignKey:ShippingRateID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	ValueBasedRates       []ValueBasedRate  `gorm:"foreignKey:ShippingRateID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	Active                bool              `gorm:"default:true"`
}

// WeightBasedRate represents additional costs based on order weight
type WeightBasedRate struct {
	gorm.Model
	ShippingRateID uint         `gorm:"index;not null"`
	ShippingRate   ShippingRate `gorm:"foreignKey:ShippingRateID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	MinWeight      float64      `gorm:"not null"`
	MaxWeight      float64      `gorm:"not null"`
	Rate           int64        `gorm:"not null"`
}

// ValueBasedRate represents additional costs/discounts based on order value
type ValueBasedRate struct {
	gorm.Model
	ShippingRateID uint         `gorm:"index;not null"`
	ShippingRate   ShippingRate `gorm:"foreignKey:ShippingRateID;constraint:OnDelete:CASCADE,OnUpdate:CASCADE"`
	MinOrderValue  int64        `gorm:"not null"`
	MaxOrderValue  int64        `gorm:"not null"`
	Rate           int64        `gorm:"not null"`
}

// ShippingOption represents a single shipping option with its cost
type ShippingOption struct {
	ShippingRateID        uint
	ShippingMethodID      uint
	Name                  string
	Description           string
	EstimatedDeliveryDays int
	Cost                  int64
	FreeShipping          bool
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

	return &ShippingRate{
		ShippingMethodID: shippingMethodID,
		ShippingZoneID:   shippingZoneID,
		BaseRate:         baseRate,
		MinOrderValue:    minOrderValue,
		Active:           true,
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

	return nil
}

// SetFreeShippingThreshold sets the free shipping threshold
func (r *ShippingRate) SetFreeShippingThreshold(threshold *int64) {
	// Validate that threshold is either nil or positive
	if threshold != nil && *threshold < 0 {
		return
	}

	r.FreeShippingThreshold = threshold

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

	}
}

// Deactivate deactivates a shipping rate
func (r *ShippingRate) Deactivate() {
	if r.Active {
		r.Active = false

	}
}
