package dto

import (
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
)

// CheckoutDTO represents a checkout session in the system
type CheckoutDTO struct {
	ID               uint                `json:"id"`
	UserID           uint                `json:"user_id,omitempty"`
	SessionID        string              `json:"session_id,omitempty"`
	Items            []CheckoutItemDTO   `json:"items"`
	Status           string              `json:"status"`
	ShippingAddress  AddressDTO          `json:"shipping_address"`
	BillingAddress   AddressDTO          `json:"billing_address"`
	ShippingMethodID uint                `json:"shipping_method_id,omitempty"`
	ShippingOption   *ShippingOptionDTO  `json:"shipping_option,omitempty"`
	PaymentProvider  string              `json:"payment_provider,omitempty"`
	TotalAmount      float64             `json:"total_amount"`
	ShippingCost     float64             `json:"shipping_cost"`
	TotalWeight      float64             `json:"total_weight"`
	CustomerDetails  CustomerDetailsDTO  `json:"customer_details"`
	Currency         string              `json:"currency"`
	DiscountCode     string              `json:"discount_code,omitempty"`
	DiscountAmount   float64             `json:"discount_amount"`
	FinalAmount      float64             `json:"final_amount"`
	AppliedDiscount  *AppliedDiscountDTO `json:"applied_discount,omitempty"`
	CreatedAt        time.Time           `json:"created_at"`
	UpdatedAt        time.Time           `json:"updated_at"`
	LastActivityAt   time.Time           `json:"last_activity_at"`
	ExpiresAt        time.Time           `json:"expires_at"`
	CompletedAt      *time.Time          `json:"completed_at,omitempty"`
	ConvertedOrderID uint                `json:"converted_order_id,omitempty"`
}

// CheckoutItemDTO represents an item in a checkout
type CheckoutItemDTO struct {
	ID          uint      `json:"id"`
	ProductID   uint      `json:"product_id"`
	VariantID   uint      `json:"variant_id,omitempty"`
	ProductName string    `json:"product_name"`
	VariantName string    `json:"variant_name,omitempty"`
	ImageURL    string    `json:"image_url"`
	SKU         string    `json:"sku"`
	Price       float64   `json:"price"`
	Quantity    int       `json:"quantity"`
	Weight      float64   `json:"weight"`
	Subtotal    float64   `json:"subtotal"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// AddToCheckoutRequest represents the data needed to add an item to a checkout
type AddToCheckoutRequest struct {
	SKU      string `json:"sku"`
	Quantity int    `json:"quantity"`
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
	ListResponseDTO[CheckoutDTO]
}

// CheckoutSearchRequest represents the parameters for searching checkouts
type CheckoutSearchRequest struct {
	UserID uint   `json:"user_id,omitempty"`
	Status string `json:"status,omitempty"`
	PaginationDTO
}

type CheckoutCompleteResponse struct {
	Order          OrderSummaryDTO `json:"order"`
	ActionRequired bool            `json:"action_required,omitempty"`
	ActionURL      string          `json:"redirect_url,omitempty"`
}

// CompleteCheckoutRequest represents the data needed to convert a checkout to an order
type CompleteCheckoutRequest struct {
	PaymentProvider string      `json:"payment_provider"`
	PaymentData     PaymentData `json:"payment_data"`
	// RedirectURL     string      `json:"redirect_url"`
}

type PaymentData struct {
	CardDetails *CardDetailsDTO `json:"card_details,omitempty"`
	PhoneNumber string          `json:"phone_number,omitempty"`
}

// CardDetailsDTO represents card details for payment processing
type CardDetailsDTO struct {
	CardNumber     string `json:"card_number"`
	ExpiryMonth    int    `json:"expiry_month"`
	ExpiryYear     int    `json:"expiry_year"`
	CVV            string `json:"cvv"`
	CardholderName string `json:"cardholder_name"`
	Token          string `json:"token,omitempty"` // Optional token for saved cards
}

func CreateCheckoutsListResponse(checkouts []*entity.Checkout, totalCount, page, pageSize int) ListResponseDTO[CheckoutDTO] {
	var checkoutDTOs []CheckoutDTO
	for _, checkout := range checkouts {
		checkoutDTOs = append(checkoutDTOs, toCheckoutDTO(checkout))
	}

	return ListResponseDTO[CheckoutDTO]{
		Success: true,
		Data:    checkoutDTOs,
		Pagination: PaginationDTO{
			Page:     page,
			PageSize: pageSize,
			Total:    totalCount,
		},
	}
}

func CreateCheckoutResponse(checkout *entity.Checkout) ResponseDTO[CheckoutDTO] {
	return SuccessResponse(toCheckoutDTO(checkout))
}

func CreateCompleteCheckoutResponse(order *entity.Order) ResponseDTO[CheckoutCompleteResponse] {
	response := CheckoutCompleteResponse{
		Order:          ToOrderSummaryDTO(order),
		ActionRequired: order.Status == entity.OrderStatusPendingAction,
		ActionURL:      order.ActionURL,
	}
	return SuccessResponse(response)
}

// ConvertToCheckoutDTO converts a checkout entity to a DTO
func toCheckoutDTO(checkout *entity.Checkout) CheckoutDTO {
	dto := CheckoutDTO{
		ID:               checkout.ID,
		UserID:           checkout.UserID,
		SessionID:        checkout.SessionID,
		Status:           string(checkout.Status),
		ShippingMethodID: checkout.ShippingMethodID,
		PaymentProvider:  checkout.PaymentProvider,
		TotalAmount:      float64(checkout.TotalAmount) / 100,  // Convert cents to currency units
		ShippingCost:     float64(checkout.ShippingCost) / 100, // Convert cents to currency units
		TotalWeight:      checkout.TotalWeight,
		Currency:         checkout.Currency,
		DiscountCode:     checkout.DiscountCode,
		DiscountAmount:   float64(checkout.DiscountAmount) / 100, // Convert cents to currency units
		FinalAmount:      float64(checkout.FinalAmount) / 100,    // Convert cents to currency units
		CreatedAt:        checkout.CreatedAt,
		UpdatedAt:        checkout.UpdatedAt,
		LastActivityAt:   checkout.LastActivityAt,
		ExpiresAt:        checkout.ExpiresAt,
		CompletedAt:      checkout.CompletedAt,
		ConvertedOrderID: checkout.ConvertedOrderID,
	}

	// Convert items
	items := make([]CheckoutItemDTO, len(checkout.Items))
	for i, item := range checkout.Items {
		items[i] = CheckoutItemDTO{
			ID:          item.ID,
			ProductID:   item.ProductID,
			VariantID:   item.ProductVariantID,
			ProductName: item.ProductName,
			VariantName: item.VariantName,
			ImageURL:    item.ImageURL,
			SKU:         item.SKU,
			Price:       float64(item.Price) / 100, // Convert cents to currency units
			Quantity:    item.Quantity,
			Weight:      item.Weight,
			Subtotal:    float64(item.Price*int64(item.Quantity)) / 100, // Convert cents to currency units
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		}
	}
	dto.Items = items

	// Convert shipping method if present
	if checkout.ShippingOption != nil {
		option := ConvertToShippingOptionDTO(checkout.ShippingOption)
		dto.ShippingOption = &ShippingOptionDTO{
			ShippingMethodID:      option.ShippingMethodID,
			ShippingRateID:        option.ShippingRateID,
			Name:                  option.Name,
			Description:           option.Description,
			Cost:                  money.FromCents(int64(option.Cost)), // Convert cents to currency units
			EstimatedDeliveryDays: option.EstimatedDeliveryDays,
			FreeShipping:          option.FreeShipping,
		}
	}

	// Convert shipping address
	dto.ShippingAddress = AddressDTO{
		AddressLine1: checkout.ShippingAddr.Street,
		City:         checkout.ShippingAddr.City,
		State:        checkout.ShippingAddr.State,
		PostalCode:   checkout.ShippingAddr.PostalCode,
		Country:      checkout.ShippingAddr.Country,
	}

	// Convert billing address
	dto.BillingAddress = AddressDTO{
		AddressLine1: checkout.BillingAddr.Street,
		City:         checkout.BillingAddr.City,
		State:        checkout.BillingAddr.State,
		PostalCode:   checkout.BillingAddr.PostalCode,
		Country:      checkout.BillingAddr.Country,
	}

	// Convert customer details
	dto.CustomerDetails = CustomerDetailsDTO{
		Email:    checkout.CustomerDetails.Email,
		Phone:    checkout.CustomerDetails.Phone,
		FullName: checkout.CustomerDetails.FullName,
	}

	// Convert applied discount if present
	if checkout.AppliedDiscount != nil {
		dto.AppliedDiscount = &AppliedDiscountDTO{
			ID:     checkout.AppliedDiscount.DiscountID,
			Code:   checkout.AppliedDiscount.DiscountCode,
			Type:   "", // We don't have this info in the AppliedDiscount
			Method: "", // We don't have this info in the AppliedDiscount
			Value:  0,  // We don't have this info in the AppliedDiscount
			Amount: float64(checkout.AppliedDiscount.DiscountAmount) / 100,
		}
	}

	return dto
}
