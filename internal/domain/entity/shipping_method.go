package entity

import (
	"errors"

	"gorm.io/gorm"
)

// ShippingMethod represents a shipping method option (e.g., standard, express)
type ShippingMethod struct {
	gorm.Model
	Name                  string `gorm:"not null;size:255"`
	Description           string `gorm:"type:text"`
	EstimatedDeliveryDays int    `gorm:"not null;default:0"`
	Active                bool   `gorm:"default:true"`
}

// NewShippingMethod creates a new shipping method
func NewShippingMethod(name string, description string, estimatedDeliveryDays int) (*ShippingMethod, error) {
	if name == "" {
		return nil, errors.New("shipping method name cannot be empty")
	}

	if estimatedDeliveryDays < 0 {
		return nil, errors.New("estimated delivery days must be a non-negative number")
	}

	return &ShippingMethod{
		Name:                  name,
		Description:           description,
		EstimatedDeliveryDays: estimatedDeliveryDays,
		Active:                true,
	}, nil
}

// Update updates a shipping method's details
func (s *ShippingMethod) Update(name string, description string, estimatedDeliveryDays int) error {
	if name == "" {
		return errors.New("shipping method name cannot be empty")
	}

	if estimatedDeliveryDays < 0 {
		return errors.New("estimated delivery days must be a non-negative number")
	}

	s.Name = name
	s.Description = description
	s.EstimatedDeliveryDays = estimatedDeliveryDays

	return nil
}

// Activate activates a shipping method
func (s *ShippingMethod) Activate() {
	if !s.Active {
		s.Active = true

	}
}

// Deactivate deactivates a shipping method
func (s *ShippingMethod) Deactivate() {
	if s.Active {
		s.Active = false

	}
}
