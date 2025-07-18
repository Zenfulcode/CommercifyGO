package contracts

import (
	"github.com/zenfulcode/commercify/internal/domain/dto"
	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// AddToCheckoutRequest represents the data needed to add an item to a checkout
type AddToCheckoutRequest struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
	Currency string `json:"currency,omitempty"` // Optional currency for checkout creation/updates
}

// UpdateCheckoutItemRequest represents the data needed to update a checkout item
type UpdateCheckoutItemRequest struct {
	Quantity int `json:"quantity"`
}

// SetShippingAddressRequest represents the data needed to set a shipping address
type SetShippingAddressRequest struct {
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
}

// SetBillingAddressRequest represents the data needed to set a billing address
type SetBillingAddressRequest struct {
	AddressLine1 string `json:"address_line1"`
	AddressLine2 string `json:"address_line2"`
	City         string `json:"city"`
	State        string `json:"state"`
	PostalCode   string `json:"postal_code"`
	Country      string `json:"country"`
}

// SetCustomerDetailsRequest represents the data needed to set customer details
type SetCustomerDetailsRequest struct {
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	FullName string `json:"full_name"`
}

// SetShippingMethodRequest represents the data needed to set a shipping method
type SetShippingMethodRequest struct {
	ShippingMethodID uint `json:"shipping_method_id"`
}

// SetCurrencyRequest represents the data needed to change checkout currency
type SetCurrencyRequest struct {
	Currency string `json:"currency"`
}

// ApplyDiscountRequest represents the data needed to apply a discount
type ApplyDiscountRequest struct {
	DiscountCode string `json:"discount_code"`
}

// CheckoutListResponse represents a paginated list of checkouts
type CheckoutListResponse struct {
	ListResponseDTO[dto.CheckoutDTO]
}

// CheckoutSearchRequest represents the parameters for searching checkouts
type CheckoutSearchRequest struct {
	UserID uint   `json:"user_id,omitempty"`
	Status string `json:"status,omitempty"`
	PaginationDTO
}

type CheckoutCompleteResponse struct {
	Order          dto.OrderSummaryDTO `json:"order"`
	ActionRequired bool                `json:"action_required,omitempty"`
	ActionURL      string              `json:"redirect_url,omitempty"`
}

// CompleteCheckoutRequest represents the data needed to convert a checkout to an order
type CompleteCheckoutRequest struct {
	PaymentProvider string      `json:"payment_provider"`
	PaymentData     PaymentData `json:"payment_data"`
}

type PaymentData struct {
	CardDetails *dto.CardDetailsDTO `json:"card_details,omitempty"`
	PhoneNumber string              `json:"phone_number,omitempty"`
}

func CreateCheckoutsListResponse(checkouts []*entity.Checkout, totalCount, page, pageSize int) ListResponseDTO[dto.CheckoutDTO] {
	var checkoutDTOs []dto.CheckoutDTO
	for _, checkout := range checkouts {
		checkoutDTOs = append(checkoutDTOs, *checkout.ToCheckoutDTO())
	}
	if len(checkoutDTOs) == 0 {
		return ListResponseDTO[dto.CheckoutDTO]{
			Success:    true,
			Data:       []dto.CheckoutDTO{},
			Pagination: PaginationDTO{Page: page, PageSize: pageSize, Total: 0},
			Message:    "No checkouts found",
		}
	}

	return ListResponseDTO[dto.CheckoutDTO]{
		Success: true,
		Data:    checkoutDTOs,
		Pagination: PaginationDTO{
			Page:     page,
			PageSize: pageSize,
			Total:    totalCount,
		},
		Message: "Checkouts retrieved successfully",
	}
}

func CreateCheckoutResponse(checkout *dto.CheckoutDTO) ResponseDTO[dto.CheckoutDTO] {
	return SuccessResponse(*checkout)
}

func CreateCompleteCheckoutResponse(order *entity.Order) ResponseDTO[CheckoutCompleteResponse] {
	response := CheckoutCompleteResponse{
		Order:          *order.ToOrderSummaryDTO(),
		ActionRequired: order.ActionRequired(),
		ActionURL:      order.ActionURL.String,
	}
	return SuccessResponse(response)
}
