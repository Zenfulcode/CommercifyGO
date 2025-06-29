package entity

import (
	"errors"

	"gorm.io/gorm"
)

// ShippingZone represents a geographical shipping zone
type ShippingZone struct {
	gorm.Model
	Name        string   `gorm:"not null;size:255"`
	Description string   `gorm:"type:text"`
	Countries   []string `gorm:"type:jsonb;default:'[]'"`
	States      []string `gorm:"type:jsonb;default:'[]'"`
	ZipCodes    []string `gorm:"type:jsonb;default:'[]'"`
	Active      bool     `gorm:"default:true"`
}

// NewShippingZone creates a new shipping zone
func NewShippingZone(name string, description string) (*ShippingZone, error) {
	if name == "" {
		return nil, errors.New("shipping zone name cannot be empty")
	}

	return &ShippingZone{
		Name:        name,
		Description: description,
		Countries:   []string{},
		States:      []string{},
		ZipCodes:    []string{},
		Active:      true,
	}, nil
}

// Update updates a shipping zone's details
func (z *ShippingZone) Update(name string, description string) error {
	if name == "" {
		return errors.New("shipping zone name cannot be empty")
	}

	z.Name = name
	z.Description = description

	return nil
}

// SetCountries sets the countries for this shipping zone
func (z *ShippingZone) SetCountries(countries []string) {
	z.Countries = countries

}

// SetStates sets the states/provinces for this shipping zone
func (z *ShippingZone) SetStates(states []string) {
	z.States = states

}

// SetZipCodes sets the zip/postal codes for this shipping zone
func (z *ShippingZone) SetZipCodes(zipCodes []string) {
	z.ZipCodes = zipCodes

}

// Activate activates a shipping zone
func (z *ShippingZone) Activate() {
	if !z.Active {
		z.Active = true

	}
}

// Deactivate deactivates a shipping zone
func (z *ShippingZone) Deactivate() {
	if z.Active {
		z.Active = false

	}
}

// IsAddressInZone checks if an address is within this zone
func (z *ShippingZone) IsAddressInZone(address Address) bool {
	// If no countries are specified, all countries match
	if len(z.Countries) == 0 {
		return true
	}

	// Check country match
	countryMatch := false
	for _, country := range z.Countries {
		if country == address.Country {
			countryMatch = true
			break
		}
	}

	if !countryMatch {
		return false
	}

	// If we matched country and no states are specified, it's a match
	if len(z.States) == 0 {
		return true
	}

	// Check state match
	stateMatch := false
	for _, state := range z.States {
		if state == address.State {
			stateMatch = true
			break
		}
	}

	if !stateMatch {
		return false
	}

	// If we matched country and state, and no zip codes are specified, it's a match
	if len(z.ZipCodes) == 0 {
		return true
	}

	// Check zip code match (exact match only - could be extended for patterns/ranges)
	for _, zipCode := range z.ZipCodes {
		if zipCode == address.PostalCode {
			return true
		}
	}

	return false
}
