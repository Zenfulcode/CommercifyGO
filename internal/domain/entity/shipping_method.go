package entity

import (
	"errors"
	"time"
)

// ShippingMethod represents a shipping method option (e.g., standard, express)
type ShippingMethod struct {
	ID                    uint      `json:"id"`
	Name                  string    `json:"name"`
	Description           string    `json:"description"`
	EstimatedDeliveryDays int       `json:"estimated_delivery_days"`
	Active                bool      `json:"active"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// NewShippingMethod creates a new shipping method
func NewShippingMethod(name string, description string, estimatedDeliveryDays int) (*ShippingMethod, error) {
	if name == "" {
		return nil, errors.New("shipping method name cannot be empty")
	}

	if estimatedDeliveryDays < 0 {
		return nil, errors.New("estimated delivery days must be a non-negative number")
	}

	now := time.Now()
	return &ShippingMethod{
		Name:                  name,
		Description:           description,
		EstimatedDeliveryDays: estimatedDeliveryDays,
		Active:                true,
		CreatedAt:             now,
		UpdatedAt:             now,
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
	s.UpdatedAt = time.Now()
	return nil
}

// Activate activates a shipping method
func (s *ShippingMethod) Activate() {
	if !s.Active {
		s.Active = true
		s.UpdatedAt = time.Now()
	}
}

// Deactivate deactivates a shipping method
func (s *ShippingMethod) Deactivate() {
	if s.Active {
		s.Active = false
		s.UpdatedAt = time.Now()
	}
}
