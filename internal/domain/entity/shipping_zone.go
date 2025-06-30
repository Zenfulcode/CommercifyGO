package entity

import (
	"errors"
	"slices"

	"github.com/zenfulcode/commercify/internal/dto"
	"gorm.io/gorm"
)

// ShippingZone represents a geographical shipping zone
type ShippingZone struct {
	gorm.Model
	Name        string   `gorm:"not null;size:255"`
	Description string   `gorm:"type:text"`
	Countries   []string `gorm:"type:jsonb;default:'[]'"`
	Active      bool     `gorm:"default:true"`
}

// NewShippingZone creates a new shipping zone
func NewShippingZone(name, description string, countries []string) (*ShippingZone, error) {
	if name == "" {
		return nil, errors.New("shipping zone name cannot be empty")
	}

	return &ShippingZone{
		Name:        name,
		Description: description,
		Countries:   countries,
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

	return slices.Contains(z.Countries, address.Country)
}

func (z *ShippingZone) ToShippingZoneDTO() *dto.ShippingZoneDTO {
	return &dto.ShippingZoneDTO{
		ID:          z.ID,
		Name:        z.Name,
		Description: z.Description,
		Countries:   z.Countries,
		Active:      z.Active,
		CreatedAt:   z.CreatedAt,
		UpdatedAt:   z.UpdatedAt,
	}
}
