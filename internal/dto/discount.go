package dto

import (
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
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

// CreateDiscountRequest represents the data needed to create a new discount
type CreateDiscountRequest struct {
	Code             string    `json:"code" validate:"required,min=3,max=50"`
	Type             string    `json:"type" validate:"required,oneof=basket product"`
	Method           string    `json:"method" validate:"required,oneof=fixed percentage"`
	Value            float64   `json:"value" validate:"required,gt=0"`
	MinOrderValue    float64   `json:"min_order_value" validate:"gte=0"`
	MaxDiscountValue float64   `json:"max_discount_value" validate:"gte=0"`
	ProductIDs       []uint    `json:"product_ids,omitempty"`
	CategoryIDs      []uint    `json:"category_ids,omitempty"`
	StartDate        time.Time `json:"start_date" validate:"required"`
	EndDate          time.Time `json:"end_date" validate:"required"`
	UsageLimit       int       `json:"usage_limit" validate:"gte=0"`
}

// UpdateDiscountRequest represents the data needed to update a discount
type UpdateDiscountRequest struct {
	Code             string    `json:"code,omitempty" validate:"omitempty,min=3,max=50"`
	Type             string    `json:"type,omitempty" validate:"omitempty,oneof=basket product"`
	Method           string    `json:"method,omitempty" validate:"omitempty,oneof=fixed percentage"`
	Value            float64   `json:"value,omitempty" validate:"omitempty,gt=0"`
	MinOrderValue    float64   `json:"min_order_value,omitempty" validate:"omitempty,gte=0"`
	MaxDiscountValue float64   `json:"max_discount_value,omitempty" validate:"omitempty,gte=0"`
	ProductIDs       []uint    `json:"product_ids,omitempty"`
	CategoryIDs      []uint    `json:"category_ids,omitempty"`
	StartDate        time.Time `json:"start_date,omitempty"`
	EndDate          time.Time `json:"end_date,omitempty"`
	UsageLimit       int       `json:"usage_limit,omitempty" validate:"omitempty,gte=0"`
	Active           bool      `json:"active"`
}

// ValidateDiscountRequest represents the data needed to validate a discount code
type ValidateDiscountRequest struct {
	DiscountCode string `json:"discount_code" validate:"required"`
}

// ValidateDiscountResponse represents the response for discount validation
type ValidateDiscountResponse struct {
	Valid            bool    `json:"valid"`
	Reason           string  `json:"reason,omitempty"`
	DiscountID       uint    `json:"discount_id,omitempty"`
	Code             string  `json:"code,omitempty"`
	Type             string  `json:"type,omitempty"`
	Method           string  `json:"method,omitempty"`
	Value            float64 `json:"value,omitempty"`
	MinOrderValue    float64 `json:"min_order_value,omitempty"`
	MaxDiscountValue float64 `json:"max_discount_value,omitempty"`
}

// ConvertToDiscountDTO converts a domain discount entity to a DTO
func ConvertToDiscountDTO(discount *entity.Discount) DiscountDTO {
	if discount == nil {
		return DiscountDTO{}
	}

	return DiscountDTO{
		ID:               discount.ID,
		Code:             discount.Code,
		Type:             string(discount.Type),
		Method:           string(discount.Method),
		Value:            discount.Value,
		MinOrderValue:    money.FromCents(discount.MinOrderValue),
		MaxDiscountValue: money.FromCents(discount.MaxDiscountValue),
		ProductIDs:       discount.ProductIDs,
		CategoryIDs:      discount.CategoryIDs,
		StartDate:        discount.StartDate,
		EndDate:          discount.EndDate,
		UsageLimit:       discount.UsageLimit,
		CurrentUsage:     discount.CurrentUsage,
		Active:           discount.Active,
		CreatedAt:        discount.CreatedAt,
		UpdatedAt:        discount.UpdatedAt,
	}
}

// ConvertToAppliedDiscountDTO converts a domain applied discount entity to a DTO
func ConvertToAppliedDiscountDTO(appliedDiscount *entity.AppliedDiscount) AppliedDiscountDTO {
	if appliedDiscount == nil {
		return AppliedDiscountDTO{}
	}

	return AppliedDiscountDTO{
		ID:     appliedDiscount.DiscountID,
		Code:   appliedDiscount.DiscountCode,
		Type:   "", // We don't have this info in the AppliedDiscount
		Method: "", // We don't have this info in the AppliedDiscount
		Value:  0,  // We don't have this info in the AppliedDiscount
		Amount: money.FromCents(appliedDiscount.DiscountAmount),
	}
}

// ConvertDiscountListToDTO converts a slice of domain discount entities to DTOs
func ConvertDiscountListToDTO(discounts []*entity.Discount) []DiscountDTO {
	dtos := make([]DiscountDTO, len(discounts))
	for i, discount := range discounts {
		dtos[i] = ConvertToDiscountDTO(discount)
	}
	return dtos
}
