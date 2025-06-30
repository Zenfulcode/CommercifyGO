package contracts

import (
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/dto"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
)

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

// CreateShippingRateRequest represents the data needed to create a new shipping rate
type CreateShippingRateRequest struct {
	ShippingMethodID      uint     `json:"shipping_method_id"`
	ShippingZoneID        uint     `json:"shipping_zone_id"`
	BaseRate              float64  `json:"base_rate"`
	MinOrderValue         float64  `json:"min_order_value"`
	FreeShippingThreshold *float64 `json:"free_shipping_threshold"`
	Active                bool     `json:"active"`
}

// CreateValueBasedRateRequest represents the data needed to create a value-based rate
type CreateValueBasedRateRequest struct {
	ShippingRateID uint    `json:"shipping_rate_id"`
	MinOrderValue  float64 `json:"min_order_value"`
	MaxOrderValue  float64 `json:"max_order_value"`
	Rate           float64 `json:"rate"`
}

// UpdateShippingRateRequest represents the data needed to update a shipping rate
type UpdateShippingRateRequest struct {
	BaseRate              float64  `json:"base_rate,omitempty"`
	MinOrderValue         float64  `json:"min_order_value,omitempty"`
	FreeShippingThreshold *float64 `json:"free_shipping_threshold"`
	Active                bool     `json:"active"`
}

// CreateWeightBasedRateRequest represents the data needed to create a weight-based rate
type CreateWeightBasedRateRequest struct {
	ShippingRateID uint    `json:"shipping_rate_id"`
	MinWeight      float64 `json:"min_weight"`
	MaxWeight      float64 `json:"max_weight"`
	Rate           float64 `json:"rate"`
}

// CalculateShippingOptionsRequest represents the request to calculate shipping options
type CalculateShippingOptionsRequest struct {
	Address     dto.AddressDTO `json:"address"`
	OrderValue  float64        `json:"order_value"`
	OrderWeight float64        `json:"order_weight"`
}

func (c CalculateShippingOptionsRequest) ToUseCaseInput() usecase.CalculateShippingOptionsInput {
	return usecase.CalculateShippingOptionsInput{
		Address: entity.Address{
			Street1:    c.Address.AddressLine1,
			Street2:    c.Address.AddressLine2,
			City:       c.Address.City,
			State:      c.Address.State,
			Country:    c.Address.Country,
			PostalCode: c.Address.PostalCode,
		},
		OrderValue:  money.ToCents(c.OrderValue),
		OrderWeight: c.OrderWeight,
	}
}

// CalculateShippingOptionsResponse represents the response with available shipping options
type CalculateShippingOptionsResponse struct {
	Options []dto.ShippingOptionDTO `json:"options"`
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

func CreateShippingOptionsListResponse(options []*entity.ShippingOption, totalCount, page, pageSize int) ListResponseDTO[dto.ShippingOptionDTO] {
	var response []dto.ShippingOptionDTO
	for _, option := range options {
		response = append(response, *option.ToShippingOptionDTO())
	}
	if len(response) == 0 {
		return ListResponseDTO[dto.ShippingOptionDTO]{
			Success:    true,
			Data:       []dto.ShippingOptionDTO{},
			Pagination: PaginationDTO{Page: page, PageSize: pageSize, Total: 0},
			Message:    "No shipping options found",
		}
	}
	return ListResponseDTO[dto.ShippingOptionDTO]{
		Success: true,
		Data:    response,
		Pagination: PaginationDTO{
			Page:     page,
			PageSize: pageSize,
			Total:    totalCount,
		},
		Message: "Shipping options retrieved successfully",
	}
}
