package contracts

import (
	"time"

	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/interfaces/api/contracts/dto"
)

// CreateDiscountRequest represents the data needed to create a new discount
type CreateDiscountRequest struct {
	Code             string    `json:"code"`
	Type             string    `json:"type"`
	Method           string    `json:"method"`
	Value            float64   `json:"value"`
	MinOrderValue    float64   `json:"min_order_value,omitempty"`
	MaxDiscountValue float64   `json:"max_discount_value,omitempty"`
	ProductIDs       []uint    `json:"product_ids,omitempty"`
	CategoryIDs      []uint    `json:"category_ids,omitempty"`
	StartDate        time.Time `json:"start_date,omitempty"`
	EndDate          time.Time `json:"end_date,omitempty"`
	UsageLimit       int       `json:"usage_limit,omitempty"`
}

// UpdateDiscountRequest represents the data needed to update a discount
type UpdateDiscountRequest struct {
	Code             string    `json:"code,omitempty"`
	Type             string    `json:"type,omitempty"`
	Method           string    `json:"method,omitempty"`
	Value            float64   `json:"value,omitempty"`
	MinOrderValue    float64   `json:"min_order_value,omitempty"`
	MaxDiscountValue float64   `json:"max_discount_value,omitempty"`
	ProductIDs       []uint    `json:"product_ids,omitempty"`
	CategoryIDs      []uint    `json:"category_ids,omitempty"`
	StartDate        time.Time `json:"start_date"`
	EndDate          time.Time `json:"end_date"`
	UsageLimit       int       `json:"usage_limit,omitempty"`
	Active           bool      `json:"active"`
}

// ValidateDiscountRequest represents the data needed to validate a discount code
type ValidateDiscountRequest struct {
	DiscountCode string `json:"discount_code"`
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

func (r *CreateDiscountRequest) ToUseCaseInput() usecase.CreateDiscountInput {
	if r.MinOrderValue < 0 {
		r.MinOrderValue = 0
	}
	if r.MaxDiscountValue < 0 {
		r.MaxDiscountValue = 0
	}
	if r.UsageLimit < 0 {
		r.UsageLimit = 0
	}
	if r.StartDate.IsZero() {
		r.StartDate = time.Now().Local()
	}
	if r.EndDate.IsZero() {
		r.EndDate = time.Now().Local().AddDate(1, 0, 0) // Default to 1 year from now
	}
	if r.ProductIDs == nil {
		r.ProductIDs = []uint{}
	}
	if r.CategoryIDs == nil {
		r.CategoryIDs = []uint{}
	}

	return usecase.CreateDiscountInput{
		Code:             r.Code,
		Type:             r.Type,
		Method:           r.Method,
		Value:            r.Value,
		MinOrderValue:    r.MinOrderValue,
		MaxDiscountValue: r.MaxDiscountValue,
		ProductIDs:       r.ProductIDs,
		CategoryIDs:      r.CategoryIDs,
		StartDate:        r.StartDate,
		EndDate:          r.EndDate,
		UsageLimit:       r.UsageLimit,
	}
}

func (r *UpdateDiscountRequest) ToUseCaseInput() usecase.UpdateDiscountInput {
	return usecase.UpdateDiscountInput{
		Code:             r.Code,
		Type:             r.Type,
		Method:           r.Method,
		Value:            r.Value,
		MinOrderValue:    r.MinOrderValue,
		MaxDiscountValue: r.MaxDiscountValue,
		ProductIDs:       r.ProductIDs,
		CategoryIDs:      r.CategoryIDs,
		StartDate:        r.StartDate,
		EndDate:          r.EndDate,
		UsageLimit:       r.UsageLimit,
		Active:           r.Active,
	}
}

func DiscountCreateResponse(discount *dto.DiscountDTO) ResponseDTO[dto.DiscountDTO] {
	return SuccessResponseWithMessage(*discount, "Discount created successfully")
}

func DiscountRetrieveResponse(discount *dto.DiscountDTO) ResponseDTO[dto.DiscountDTO] {
	return SuccessResponse(*discount)
}

func DiscountUpdateResponse(discount *dto.DiscountDTO) ResponseDTO[dto.DiscountDTO] {
	return SuccessResponseWithMessage(*discount, "Discount updated successfully")
}

func DiscountDeleteResponse() ResponseDTO[any] {
	return SuccessResponseMessage("Discount deleted successfully")
}

func DiscountListResponse(discounts []*entity.Discount, totalCount, page, pageSize int) ListResponseDTO[dto.DiscountDTO] {
	var discountDTOs []dto.DiscountDTO
	for _, discount := range discounts {
		discountDTOs = append(discountDTOs, *discount.ToDiscountDTO())
	}

	if len(discountDTOs) == 0 {
		return ListResponseDTO[dto.DiscountDTO]{
			Success:    true,
			Data:       []dto.DiscountDTO{},
			Pagination: PaginationDTO{Page: page, PageSize: pageSize, Total: 0},
			Message:    "No discounts found",
		}
	}

	return ListResponseDTO[dto.DiscountDTO]{
		Success: true,
		Data:    discountDTOs,
		Pagination: PaginationDTO{
			Page:     page,
			PageSize: pageSize,
			Total:    totalCount,
		},
		Message: "Discounts retrieved successfully",
	}
}
