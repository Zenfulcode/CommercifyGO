package dto

import (
	"time"

	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
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

// CreateShippingMethodRequest represents the data needed to create a new shipping method
type CreateShippingMethodRequest struct {
	Name                  string `json:"name"`
	Description           string `json:"description"`
	EstimatedDeliveryDays int    `json:"estimated_delivery_days"`
}

// UpdateShippingMethodRequest represents the data needed to update a shipping method
type UpdateShippingMethodRequest struct {
	Name                  string `json:"name,omitempty"`
	Description           string `json:"description,omitempty"`
	EstimatedDeliveryDays int    `json:"estimated_delivery_days,omitempty"`
	Active                bool   `json:"active"`
}

// ShippingZoneDTO represents a shipping zone in the system
type ShippingZoneDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Countries   []string  `json:"countries"`
	States      []string  `json:"states"`
	ZipCodes    []string  `json:"zip_codes"`
	Active      bool      `json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateShippingZoneRequest represents the data needed to create a new shipping zone
type CreateShippingZoneRequest struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Countries   []string `json:"countries"`
	States      []string `json:"states"`
	ZipCodes    []string `json:"zip_codes"`
}

// UpdateShippingZoneRequest represents the data needed to update a shipping zone
type UpdateShippingZoneRequest struct {
	Name        string   `json:"name,omitempty" `
	Description string   `json:"description,omitempty"`
	Countries   []string `json:"countries,omitempty"`
	States      []string `json:"states,omitempty"`
	ZipCodes    []string `json:"zip_codes,omitempty"`
	Active      bool     `json:"active"`
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
	FreeShippingThreshold *float64                 `json:"free_shipping_threshold"`
	WeightBasedRates      []WeightBasedRateDTO     `json:"weight_based_rates,omitempty"`
	ValueBasedRates       []ValueBasedRateDTO      `json:"value_based_rates,omitempty"`
	Active                bool                     `json:"active"`
	CreatedAt             time.Time                `json:"created_at"`
	UpdatedAt             time.Time                `json:"updated_at"`
}

// CreateShippingRateRequest represents the data needed to create a new shipping rate
type CreateShippingRateRequest struct {
	ShippingMethodID      uint     `json:"shipping_method_id"`
	ShippingZoneID        uint     `json:"shipping_zone_id"`
	BaseRate              float64  `json:"base_rate"`
	MinOrderValue         float64  `json:"min_order_value"`
	FreeShippingThreshold *float64 `json:"free_shipping_threshold"`
	Active                bool     `json:"active"`
}

// UpdateShippingRateRequest represents the data needed to update a shipping rate
type UpdateShippingRateRequest struct {
	BaseRate              float64  `json:"base_rate,omitempty"`
	MinOrderValue         float64  `json:"min_order_value,omitempty"`
	FreeShippingThreshold *float64 `json:"free_shipping_threshold"`
	Active                bool     `json:"active"`
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

// CreateWeightBasedRateRequest represents the data needed to create a weight-based rate
type CreateWeightBasedRateRequest struct {
	ShippingRateID uint    `json:"shipping_rate_id"`
	MinWeight      float64 `json:"min_weight"`
	MaxWeight      float64 `json:"max_weight"`
	Rate           float64 `json:"rate"`
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

// CreateValueBasedRateRequest represents the data needed to create a value-based rate
type CreateValueBasedRateRequest struct {
	ShippingRateID uint    `json:"shipping_rate_id"`
	MinOrderValue  float64 `json:"min_order_value"`
	MaxOrderValue  float64 `json:"max_order_value"`
	Rate           float64 `json:"rate"`
}

// ShippingOptionDTO represents a shipping option with calculated cost
type ShippingOptionDTO struct {
	ShippingRateID        uint    `json:"shipping_rate_id"`
	ShippingMethodID      uint    `json:"shipping_method_id"`
	Name                  string  `json:"name"`
	Description           string  `json:"description"`
	EstimatedDeliveryDays int     `json:"estimated_delivery_days"`
	Cost                  float64 `json:"cost"`
	FreeShipping          bool    `json:"free_shipping"`
}

// CalculateShippingOptionsRequest represents the request to calculate shipping options
type CalculateShippingOptionsRequest struct {
	Address     AddressDTO `json:"address"`
	OrderValue  float64    `json:"order_value"`
	OrderWeight float64    `json:"order_weight"`
}

// CalculateShippingOptionsResponse represents the response with available shipping options
type CalculateShippingOptionsResponse struct {
	Options []ShippingOptionDTO `json:"options"`
}

// CalculateShippingCostRequest represents the request to calculate shipping cost for a specific rate
type CalculateShippingCostRequest struct {
	OrderValue  float64 `json:"order_value"`
	OrderWeight float64 `json:"order_weight"`
}

// CalculateShippingCostResponse represents the response with calculated shipping cost
type CalculateShippingCostResponse struct {
	Cost float64 `json:"cost"`
}

// ConvertToShippingMethodDetailDTO converts a domain shipping method entity to a DTO
func ConvertToShippingMethodDetailDTO(method *entity.ShippingMethod) ShippingMethodDetailDTO {
	if method == nil {
		return ShippingMethodDetailDTO{}
	}

	return ShippingMethodDetailDTO{
		ID:                    method.ID,
		Name:                  method.Name,
		Description:           method.Description,
		EstimatedDeliveryDays: method.EstimatedDeliveryDays,
		Active:                method.Active,
		CreatedAt:             method.CreatedAt,
		UpdatedAt:             method.UpdatedAt,
	}
}

// ConvertToShippingZoneDTO converts a domain shipping zone entity to a DTO
func ConvertToShippingZoneDTO(zone *entity.ShippingZone) ShippingZoneDTO {
	if zone == nil {
		return ShippingZoneDTO{}
	}

	return ShippingZoneDTO{
		ID:          zone.ID,
		Name:        zone.Name,
		Description: zone.Description,
		Countries:   zone.Countries,
		States:      zone.States,
		ZipCodes:    zone.ZipCodes,
		Active:      zone.Active,
		CreatedAt:   zone.CreatedAt,
		UpdatedAt:   zone.UpdatedAt,
	}
}

// ConvertToShippingRateDTO converts a domain shipping rate entity to a DTO
func ConvertToShippingRateDTO(rate *entity.ShippingRate) ShippingRateDTO {
	if rate == nil {
		return ShippingRateDTO{}
	}

	dto := ShippingRateDTO{
		ID:               rate.ID,
		ShippingMethodID: rate.ShippingMethodID,
		ShippingZoneID:   rate.ShippingZoneID,
		BaseRate:         money.FromCents(rate.BaseRate),
		MinOrderValue:    money.FromCents(rate.MinOrderValue),
		Active:           rate.Active,
		CreatedAt:        rate.CreatedAt,
		UpdatedAt:        rate.UpdatedAt,
	}

	// Convert free shipping threshold
	if rate.FreeShippingThreshold != nil {
		threshold := money.FromCents(*rate.FreeShippingThreshold)
		dto.FreeShippingThreshold = &threshold
	}

	// Convert shipping method if available
	if rate.ShippingMethod != nil {
		method := ConvertToShippingMethodDetailDTO(rate.ShippingMethod)
		dto.ShippingMethod = &method
	}

	// Convert shipping zone if available
	if rate.ShippingZone != nil {
		zone := ConvertToShippingZoneDTO(rate.ShippingZone)
		dto.ShippingZone = &zone
	}

	// Convert weight-based rates
	if len(rate.WeightBasedRates) > 0 {
		dto.WeightBasedRates = make([]WeightBasedRateDTO, len(rate.WeightBasedRates))
		for i, wbr := range rate.WeightBasedRates {
			dto.WeightBasedRates[i] = ConvertToWeightBasedRateDTO(&wbr)
		}
	}

	// Convert value-based rates
	if len(rate.ValueBasedRates) > 0 {
		dto.ValueBasedRates = make([]ValueBasedRateDTO, len(rate.ValueBasedRates))
		for i, vbr := range rate.ValueBasedRates {
			dto.ValueBasedRates[i] = ConvertToValueBasedRateDTO(&vbr)
		}
	}

	return dto
}

// ConvertToWeightBasedRateDTO converts a domain weight-based rate entity to a DTO
func ConvertToWeightBasedRateDTO(rate *entity.WeightBasedRate) WeightBasedRateDTO {
	if rate == nil {
		return WeightBasedRateDTO{}
	}

	return WeightBasedRateDTO{
		ID:             rate.ID,
		ShippingRateID: rate.ShippingRateID,
		MinWeight:      rate.MinWeight,
		MaxWeight:      rate.MaxWeight,
		Rate:           money.FromCents(rate.Rate),
		CreatedAt:      rate.CreatedAt,
		UpdatedAt:      rate.UpdatedAt,
	}
}

// ConvertToValueBasedRateDTO converts a domain value-based rate entity to a DTO
func ConvertToValueBasedRateDTO(rate *entity.ValueBasedRate) ValueBasedRateDTO {
	if rate == nil {
		return ValueBasedRateDTO{}
	}

	return ValueBasedRateDTO{
		ID:             rate.ID,
		ShippingRateID: rate.ShippingRateID,
		MinOrderValue:  money.FromCents(rate.MinOrderValue),
		MaxOrderValue:  money.FromCents(rate.MaxOrderValue),
		Rate:           money.FromCents(rate.Rate),
		CreatedAt:      rate.CreatedAt,
		UpdatedAt:      rate.UpdatedAt,
	}
}

// ConvertToShippingOptionDTO converts a domain shipping option entity to a DTO
func ConvertToShippingOptionDTO(option *entity.ShippingOption) ShippingOptionDTO {
	if option == nil {
		return ShippingOptionDTO{}
	}

	return ShippingOptionDTO{
		ShippingRateID:        option.ShippingRateID,
		ShippingMethodID:      option.ShippingMethodID,
		Name:                  option.Name,
		Description:           option.Description,
		EstimatedDeliveryDays: option.EstimatedDeliveryDays,
		Cost:                  money.FromCents(option.Cost),
		FreeShipping:          option.FreeShipping,
	}
}

// ConvertShippingMethodListToDTO converts a slice of domain shipping method entities to DTOs
func ConvertShippingMethodListToDTO(methods []*entity.ShippingMethod) []ShippingMethodDetailDTO {
	dtos := make([]ShippingMethodDetailDTO, len(methods))
	for i, method := range methods {
		dtos[i] = ConvertToShippingMethodDetailDTO(method)
	}
	return dtos
}

// ConvertShippingZoneListToDTO converts a slice of domain shipping zone entities to DTOs
func ConvertShippingZoneListToDTO(zones []*entity.ShippingZone) []ShippingZoneDTO {
	dtos := make([]ShippingZoneDTO, len(zones))
	for i, zone := range zones {
		dtos[i] = ConvertToShippingZoneDTO(zone)
	}
	return dtos
}

// ConvertShippingRateListToDTO converts a slice of domain shipping rate entities to DTOs
func ConvertShippingRateListToDTO(rates []*entity.ShippingRate) []ShippingRateDTO {
	dtos := make([]ShippingRateDTO, len(rates))
	for i, rate := range rates {
		dtos[i] = ConvertToShippingRateDTO(rate)
	}
	return dtos
}

// ConvertShippingOptionListToDTO converts a slice of domain shipping option entities to DTOs
func ConvertShippingOptionListToDTO(options []*entity.ShippingOption) []ShippingOptionDTO {
	dtos := make([]ShippingOptionDTO, len(options))
	for i, option := range options {
		dtos[i] = ConvertToShippingOptionDTO(option)
	}
	return dtos
}

// Conversion functions from DTOs to use case inputs

// ToCreateShippingMethodInput converts a CreateShippingMethodRequest DTO to use case input
func (req CreateShippingMethodRequest) ToCreateShippingMethodInput() usecase.CreateShippingMethodInput {
	return usecase.CreateShippingMethodInput{
		Name:                  req.Name,
		Description:           req.Description,
		EstimatedDeliveryDays: req.EstimatedDeliveryDays,
	}
}

// ToUpdateShippingMethodInput converts an UpdateShippingMethodRequest DTO to use case input
func (req UpdateShippingMethodRequest) ToUpdateShippingMethodInput(id uint) usecase.UpdateShippingMethodInput {
	return usecase.UpdateShippingMethodInput{
		ID:                    id,
		Name:                  req.Name,
		Description:           req.Description,
		EstimatedDeliveryDays: req.EstimatedDeliveryDays,
		Active:                req.Active,
	}
}

// ToCreateShippingZoneInput converts a CreateShippingZoneRequest DTO to use case input
func (req CreateShippingZoneRequest) ToCreateShippingZoneInput() usecase.CreateShippingZoneInput {
	return usecase.CreateShippingZoneInput{
		Name:        req.Name,
		Description: req.Description,
		Countries:   req.Countries,
		States:      req.States,
		ZipCodes:    req.ZipCodes,
	}
}

// ToUpdateShippingZoneInput converts an UpdateShippingZoneRequest DTO to use case input
func (req UpdateShippingZoneRequest) ToUpdateShippingZoneInput(id uint) usecase.UpdateShippingZoneInput {
	return usecase.UpdateShippingZoneInput{
		ID:          id,
		Name:        req.Name,
		Description: req.Description,
		Countries:   req.Countries,
		States:      req.States,
		ZipCodes:    req.ZipCodes,
		Active:      req.Active,
	}
}

// ToCreateShippingRateInput converts a CreateShippingRateRequest DTO to use case input
func (req CreateShippingRateRequest) ToCreateShippingRateInput() usecase.CreateShippingRateInput {
	return usecase.CreateShippingRateInput{
		ShippingMethodID:      req.ShippingMethodID,
		ShippingZoneID:        req.ShippingZoneID,
		BaseRate:              req.BaseRate,
		MinOrderValue:         req.MinOrderValue,
		FreeShippingThreshold: req.FreeShippingThreshold,
		Active:                req.Active,
	}
}

// ToUpdateShippingRateInput converts an UpdateShippingRateRequest DTO to use case input
func (req UpdateShippingRateRequest) ToUpdateShippingRateInput(id uint) usecase.UpdateShippingRateInput {
	return usecase.UpdateShippingRateInput{
		ID:                    id,
		BaseRate:              req.BaseRate,
		MinOrderValue:         req.MinOrderValue,
		FreeShippingThreshold: req.FreeShippingThreshold,
		Active:                req.Active,
	}
}

// ToCreateWeightBasedRateInput converts a CreateWeightBasedRateRequest DTO to use case input
func (req CreateWeightBasedRateRequest) ToCreateWeightBasedRateInput() usecase.CreateWeightBasedRateInput {
	return usecase.CreateWeightBasedRateInput{
		ShippingRateID: req.ShippingRateID,
		MinWeight:      req.MinWeight,
		MaxWeight:      req.MaxWeight,
		Rate:           req.Rate,
	}
}

// ToCreateValueBasedRateInput converts a CreateValueBasedRateRequest DTO to use case input
func (req CreateValueBasedRateRequest) ToCreateValueBasedRateInput() usecase.CreateValueBasedRateInput {
	return usecase.CreateValueBasedRateInput{
		ShippingRateID: req.ShippingRateID,
		MinOrderValue:  req.MinOrderValue,
		MaxOrderValue:  req.MaxOrderValue,
		Rate:           req.Rate,
	}
}

// ToEntityAddress converts an AddressDTO to an entity.Address for use case operations
func (addr AddressDTO) ToEntityAddress() entity.Address {
	return entity.Address{
		Street:     addr.AddressLine1,
		City:       addr.City,
		State:      addr.State,
		PostalCode: addr.PostalCode,
		Country:    addr.Country,
	}
}

// ToDomainAddress is an alias for ToEntityAddress for consistency
func (addr AddressDTO) ToDomainAddress() entity.Address {
	return addr.ToEntityAddress()
}
