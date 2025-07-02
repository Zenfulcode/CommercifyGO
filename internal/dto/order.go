package dto

import (
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
)

// OrderDTO represents an order in the system
type OrderDTO struct {
	ID              uint               `json:"id"`
	UserID          uint               `json:"user_id"`
	OrderNumber     string             `json:"order_number"`
	Items           []OrderItemDTO     `json:"items"`
	Status          OrderStatus        `json:"status"`
	PaymentStatus   PaymentStatus      `json:"payment_status"`
	TotalAmount     float64            `json:"total_amount"`  // Subtotal (items only)
	ShippingCost    float64            `json:"shipping_cost"` // Shipping cost
	FinalAmount     float64            `json:"final_amount"`  // Total including shipping and discounts
	Currency        string             `json:"currency"`
	ShippingAddress AddressDTO         `json:"shipping_address"`
	BillingAddress  AddressDTO         `json:"billing_address"`
	PaymentDetails  PaymentDetails     `json:"payment_details"`
	ShippingDetails ShippingOptionDTO  `json:"shipping_details"`
	DiscountDetails AppliedDiscountDTO `json:"discount_details"`
	Customer        CustomerDetailsDTO `json:"customer"`
	CheckoutID      string             `json:"checkout_id"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

type OrderSummaryDTO struct {
	ID               uint               `json:"id"`
	OrderNumber      string             `json:"order_number"`
	CheckoutID       string             `json:"checkout_id"`
	UserID           uint               `json:"user_id"`
	Customer         CustomerDetailsDTO `json:"customer"`
	Status           OrderStatus        `json:"status"`
	PaymentStatus    PaymentStatus      `json:"payment_status"`
	TotalAmount      float64            `json:"total_amount"`  // Subtotal (items only)
	ShippingCost     float64            `json:"shipping_cost"` // Shipping cost
	FinalAmount      float64            `json:"final_amount"`  // Total including shipping and discounts
	OrderLinesAmount int                `json:"order_lines_amount"`
	Currency         string             `json:"currency"`
	CreatedAt        time.Time          `json:"created_at"`
	UpdatedAt        time.Time          `json:"updated_at"`
}

type PaymentDetails struct {
	PaymentID string          `json:"payment_id"`
	Provider  PaymentProvider `json:"provider"`
	Method    PaymentMethod   `json:"method"`
	Status    string          `json:"status"`
	Captured  bool            `json:"captured"`
	Refunded  bool            `json:"refunded"`
}

// OrderItemDTO represents an item in an order
type OrderItemDTO struct {
	ID          uint      `json:"id"`
	OrderID     uint      `json:"order_id"`
	ProductID   uint      `json:"product_id"`
	VariantID   uint      `json:"variant_id,omitempty"`
	SKU         string    `json:"sku"`
	ProductName string    `json:"product_name"`
	VariantName string    `json:"variant_name"`
	Quantity    int       `json:"quantity"`
	UnitPrice   float64   `json:"unit_price"`
	TotalPrice  float64   `json:"total_price"`
	ImageURL    string    `json:"image_url,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateOrderRequest represents the data needed to create a new order
type CreateOrderRequest struct {
	FirstName        string     `json:"first_name"`
	LastName         string     `json:"last_name"`
	Email            string     `json:"email"`
	PhoneNumber      string     `json:"phone_number,omitempty"`
	ShippingAddress  AddressDTO `json:"shipping_address"`
	BillingAddress   AddressDTO `json:"billing_address"`
	ShippingMethodID uint       `json:"shipping_method_id"`
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
	UserID        uint        `json:"user_id,omitempty"`
	Status        OrderStatus `json:"status,omitempty"`
	PaymentStatus string      `json:"payment_status,omitempty"`
	StartDate     *time.Time  `json:"start_date,omitempty"`
	EndDate       *time.Time  `json:"end_date,omitempty"`
	PaginationDTO `json:"pagination"`
}

// OrderStatus represents the status of an order
type OrderStatus string

const (
	OrderStatusPending   OrderStatus = "pending"
	OrderStatusPaid      OrderStatus = "paid"
	OrderStatusShipped   OrderStatus = "shipped"
	OrderStatusCancelled OrderStatus = "cancelled"
	OrderStatusCompleted OrderStatus = "completed"
)

// PaymentStatus represents the status of a payment
type PaymentStatus string

const (
	PaymentStatusPending    PaymentStatus = "pending"
	PaymentStatusAuthorized PaymentStatus = "authorized"
	PaymentStatusCaptured   PaymentStatus = "captured"
	PaymentStatusRefunded   PaymentStatus = "refunded"
	PaymentStatusCancelled  PaymentStatus = "cancelled"
	PaymentStatusFailed     PaymentStatus = "failed"
)

// PaymentMethod represents the payment method used for an order
type PaymentMethod string

const (
	PaymentMethodCard   PaymentMethod = "credit_card"
	PaymentMethodWallet PaymentMethod = "wallet"
)

// PaymentProvider represents the payment provider used for an order
type PaymentProvider string

const (
	PaymentProviderStripe    PaymentProvider = "stripe"
	PaymentProviderMobilePay PaymentProvider = "mobilepay"
)

func OrderUpdateStatusResponse(order *entity.Order) ResponseDTO[OrderSummaryDTO] {
	return SuccessResponseWithMessage(ToOrderSummaryDTO(order), "Order status updated successfully")
}

func OrderSummaryListResponse(orders []*entity.Order, page, pageSize, total int) ListResponseDTO[OrderSummaryDTO] {
	var orderSummaries []OrderSummaryDTO
	for _, order := range orders {
		orderSummaries = append(orderSummaries, ToOrderSummaryDTO(order))
	}

	return ListResponseDTO[OrderSummaryDTO]{
		Success: true,
		Message: "Order summaries retrieved successfully",
		Data:    orderSummaries,
		Pagination: PaginationDTO{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}
}

func OrderDetailResponse(order *entity.Order) ResponseDTO[OrderDTO] {
	return SuccessResponse(toOrderDTO(order))
}

// toOrderSummaryDTO converts an Order entity to OrderSummaryDTO
func ToOrderSummaryDTO(order *entity.Order) OrderSummaryDTO {
	return OrderSummaryDTO{
		ID:               order.ID,
		OrderNumber:      order.OrderNumber,
		CheckoutID:       order.CheckoutSessionID,
		UserID:           order.UserID,
		Status:           OrderStatus(order.Status),
		PaymentStatus:    PaymentStatus(order.PaymentStatus),
		TotalAmount:      money.FromCents(order.TotalAmount),
		ShippingCost:     money.FromCents(order.ShippingCost),
		FinalAmount:      money.FromCents(order.FinalAmount),
		OrderLinesAmount: len(order.Items),
		Currency:         order.Currency,
		Customer: CustomerDetailsDTO{
			Email:    order.CustomerDetails.Email,
			Phone:    order.CustomerDetails.Phone,
			FullName: order.CustomerDetails.FullName,
		},
		CreatedAt: order.CreatedAt,
		UpdatedAt: order.UpdatedAt,
	}
}

func toOrderDTO(order *entity.Order) OrderDTO {
	// Convert order items to DTOs
	var items []OrderItemDTO
	if len(order.Items) > 0 {
		items = make([]OrderItemDTO, len(order.Items))
		for i, item := range order.Items {
			items[i] = OrderItemDTO{
				ID:          item.ID,
				OrderID:     order.ID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				UnitPrice:   money.FromCents(item.Price),
				TotalPrice:  money.FromCents(item.Subtotal),
				ImageURL:    item.ImageURL,
				SKU:         item.SKU,
				ProductName: item.ProductName,
				VariantID:   item.ProductVariantID,
				CreatedAt:   order.CreatedAt,
				UpdatedAt:   order.UpdatedAt,
			}
		}
	}

	// Convert addresses to DTOs
	var shippingAddr *AddressDTO
	if order.ShippingAddr.Street != "" {
		shippingAddr = &AddressDTO{
			AddressLine1: order.ShippingAddr.Street,
			City:         order.ShippingAddr.City,
			State:        order.ShippingAddr.State,
			PostalCode:   order.ShippingAddr.PostalCode,
			Country:      order.ShippingAddr.Country,
		}
	}

	var billingAddr *AddressDTO
	if order.BillingAddr.Street != "" {
		billingAddr = &AddressDTO{
			AddressLine1: order.BillingAddr.Street,
			City:         order.BillingAddr.City,
			State:        order.BillingAddr.State,
			PostalCode:   order.BillingAddr.PostalCode,
			Country:      order.BillingAddr.Country,
		}
	}

	customerDetails := CustomerDetailsDTO{
		Email:    order.CustomerDetails.Email,
		Phone:    order.CustomerDetails.Phone,
		FullName: order.CustomerDetails.FullName,
	}

	paymentDetails := PaymentDetails{
		PaymentID: order.PaymentID,
		Provider:  PaymentProvider(order.PaymentProvider),
		Method:    PaymentMethod(order.PaymentMethod),
		Captured:  order.IsCaptured(),
		Refunded:  order.IsRefunded(),
	}

	var discountDetails AppliedDiscountDTO
	if order.AppliedDiscount != nil {
		discountDetails = AppliedDiscountDTO{
			ID:     order.AppliedDiscount.DiscountID,
			Code:   order.AppliedDiscount.DiscountCode,
			Amount: money.FromCents(order.AppliedDiscount.DiscountAmount),
			Type:   "",
			Method: "",
			Value:  0,
		}
	}

	var shippingDetails ShippingOptionDTO
	if order.ShippingOption != nil {
		shippingDetails = ConvertToShippingOptionDTO(order.ShippingOption)
	}

	return OrderDTO{
		ID:              order.ID,
		OrderNumber:     order.OrderNumber,
		UserID:          order.UserID,
		Status:          OrderStatus(order.Status),
		PaymentStatus:   PaymentStatus(order.PaymentStatus),
		TotalAmount:     money.FromCents(order.TotalAmount),
		ShippingCost:    money.FromCents(order.ShippingCost),
		FinalAmount:     money.FromCents(order.FinalAmount),
		Currency:        order.Currency,
		Items:           items,
		ShippingAddress: *shippingAddr,
		BillingAddress:  *billingAddr,
		PaymentDetails:  paymentDetails,
		ShippingDetails: shippingDetails,
		DiscountDetails: discountDetails,
		Customer:        customerDetails,
		CheckoutID:      order.CheckoutSessionID,
		CreatedAt:       order.CreatedAt,
		UpdatedAt:       order.UpdatedAt,
	}
}
