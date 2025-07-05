package contracts

import (
	"time"

	"github.com/zenfulcode/commercify/internal/domain/dto"
	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// CreateOrderRequest represents the data needed to create a new order
type CreateOrderRequest struct {
	FirstName        string         `json:"first_name"`
	LastName         string         `json:"last_name"`
	Email            string         `json:"email"`
	PhoneNumber      string         `json:"phone_number,omitempty"`
	ShippingAddress  dto.AddressDTO `json:"shipping_address"`
	BillingAddress   dto.AddressDTO `json:"billing_address"`
	ShippingMethodID uint           `json:"shipping_method_id"`
}

// CreateOrderItemRequest represents the data needed to create a new order item
type CreateOrderItemRequest struct {
	ProductID uint `json:"product_id"`
	VariantID uint `json:"variant_id,omitempty"`
	Quantity  int  `json:"quantity"`
}

// UpdateOrderRequest represents the data needed to update an existing order
type UpdateOrderRequest struct {
	Status            string     `json:"status,omitempty"`
	PaymentStatus     string     `json:"payment_status,omitempty"`
	TrackingNumber    string     `json:"tracking_number,omitempty"`
	EstimatedDelivery *time.Time `json:"estimated_delivery,omitempty"`
}

// OrderSearchRequest represents the parameters for searching orders
type OrderSearchRequest struct {
	UserID        uint            `json:"user_id,omitempty"`
	Status        dto.OrderStatus `json:"status,omitempty"`
	PaymentStatus string          `json:"payment_status,omitempty"`
	StartDate     *time.Time      `json:"start_date,omitempty"`
	EndDate       *time.Time      `json:"end_date,omitempty"`
	PaginationDTO `json:"pagination"`
}

func OrderUpdateStatusResponse(orderSummary dto.OrderSummaryDTO) ResponseDTO[dto.OrderSummaryDTO] {
	return SuccessResponseWithMessage(orderSummary, "Order status updated successfully")
}

func OrderSummaryListResponse(orderSummaries []*entity.Order, page, pageSize, total int) ListResponseDTO[dto.OrderSummaryDTO] {
	var orderSummaryDTOs []dto.OrderSummaryDTO
	for _, order := range orderSummaries {
		orderSummaryDTOs = append(orderSummaryDTOs, *order.ToOrderSummaryDTO())
	}

	if len(orderSummaryDTOs) == 0 {
		return ListResponseDTO[dto.OrderSummaryDTO]{
			Success:    true,
			Data:       []dto.OrderSummaryDTO{},
			Pagination: PaginationDTO{Page: page, PageSize: pageSize, Total: 0},
			Message:    "No orders found",
		}
	}

	return ListResponseDTO[dto.OrderSummaryDTO]{
		Success: true,
		Data:    orderSummaryDTOs,
		Pagination: PaginationDTO{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
		Message: "Order summaries retrieved successfully",
	}
}

func OrderDetailResponse(order *dto.OrderDTO) ResponseDTO[dto.OrderDTO] {
	return SuccessResponse(*order)
}
