package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/dto"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// OrderHandler handles order-related HTTP requests
type OrderHandler struct {
	orderUseCase *usecase.OrderUseCase
	logger       logger.Logger
}

// NewOrderHandler creates a new OrderHandler
func NewOrderHandler(orderUseCase *usecase.OrderUseCase, logger logger.Logger) *OrderHandler {
	return &OrderHandler{
		orderUseCase: orderUseCase,
		logger:       logger,
	}
}

// GetOrder handles getting an order by ID
func (h *OrderHandler) GetOrder(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get order ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["orderId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Get order
	order, err := h.orderUseCase.GetOrderByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get order: %v", err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Check if the user is authorized to view this order
	if order.UserID != userID {
		role, ok := r.Context().Value("role").(string)
		if !ok || role != "admin" {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
	}

	// Convert order to DTO
	orderDTO := convertToOrderDTO(order)

	// Return order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderDTO)
}

// ListOrders handles listing orders for a user
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse pagination parameters
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10 // Default limit
	}

	// Get orders
	orders, err := h.orderUseCase.GetUserOrders(userID, offset, limit)
	if err != nil {
		h.logger.Error("Failed to list orders: %v", err)
		http.Error(w, "Failed to list orders", http.StatusInternalServerError)
		return
	}

	// Convert orders to DTOs
	orderDTOs := make([]dto.OrderDTO, len(orders))
	for i, order := range orders {
		orderDTOs[i] = convertToOrderDTO(order)
	}

	// Create response
	response := dto.OrderListResponse{
		ListResponseDTO: dto.ListResponseDTO[dto.OrderDTO]{
			Success: true,
			Data:    orderDTOs,
			Pagination: dto.PaginationDTO{
				Page:     offset/limit + 1,
				PageSize: limit,
				Total:    len(orderDTOs),
			},
		},
	}

	// Return orders
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ProcessPayment handles payment processing for an order
func (h *OrderHandler) ProcessPayment(w http.ResponseWriter, r *http.Request) {
	// Get order ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["orderId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Get the order
	order, err := h.orderUseCase.GetOrderByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get order: %v", err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	// Parse request body
	var paymentInput dto.ProcessPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentInput); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// For registered users, verify authorization
	if order.UserID > 0 {
		// Get user ID from context
		userID, ok := r.Context().Value("user_id").(uint)
		if !ok || userID == 0 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Check if the user is authorized to process payment for this order
		if order.UserID != userID {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}
	} else {
		// For guest orders, check the session cookie
		if !order.IsGuestOrder {
			http.Error(w, "Invalid order type", http.StatusBadRequest)
			return
		}

		// Only allow payment processing for guest orders if they have a valid cookie
		cookie, cookieErr := r.Cookie(common.SessionCookieName)
		if cookieErr != nil || cookie.Value == "" {
			http.Error(w, "Invalid session", http.StatusUnauthorized)
			return
		}
	}

	// Set up payment method based on input
	var paymentMethod service.PaymentMethod
	switch paymentInput.PaymentMethod {
	case "credit_card":
		paymentMethod = service.PaymentMethodCreditCard
		if paymentInput.CardDetails == nil {
			http.Error(w, "Card details are required for credit card payment", http.StatusBadRequest)
			return
		}
	case "wallet":
		paymentMethod = service.PaymentMethodWallet
	default:
		http.Error(w, "Invalid payment method", http.StatusBadRequest)
		return
	}

	// Set up payment provider based on input
	var paymentProvider service.PaymentProviderType
	switch paymentInput.PaymentProvider {
	case "stripe":
		paymentProvider = service.PaymentProviderStripe
	case "mock":
		paymentProvider = service.PaymentProviderMock
	case "mobilepay":
		paymentProvider = service.PaymentProviderMobilePay
	default:
		http.Error(w, "Invalid payment provider", http.StatusBadRequest)
		return
	}

	// Get customer email based on order type
	customerEmail := ""
	if order.IsGuestOrder {
		customerEmail = order.CustomerDetails.Email
	} else {
		// For registered users, get email from user repository
		user, err := h.orderUseCase.GetUserByID(order.UserID)
		if err != nil {
			h.logger.Error("Failed to get user: %v", err)
			http.Error(w, "Failed to process payment", http.StatusInternalServerError)
			return
		}
		customerEmail = user.Email
	}

	// Process payment
	input := usecase.ProcessPaymentInput{
		OrderID:         uint(id),
		PaymentMethod:   paymentMethod,
		PaymentProvider: paymentProvider,
		CardDetails:     paymentInput.CardDetails,
		CustomerEmail:   customerEmail,
		PhoneNumber:     paymentInput.PhoneNumber,
	}

	updatedOrder, err := h.orderUseCase.ProcessPayment(input)
	if err != nil {
		h.logger.Error("Failed to process payment: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert order to DTO
	orderDTO := convertToOrderDTO(updatedOrder)

	// Return updated order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderDTO)
}

// ListAllOrders handles listing all orders (admin only)
func (h *OrderHandler) ListAllOrders(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	status := r.URL.Query().Get("status")

	if limit <= 0 {
		limit = 10 // Default limit
	}

	// Get orders by status if provided
	var orders []*entity.Order
	var err error

	if status != "" {
		orders, err = h.orderUseCase.ListOrdersByStatus(entity.OrderStatus(status), offset, limit)
	} else {
		orders, err = h.orderUseCase.ListAllOrders(offset, limit)
	}

	if err != nil {
		h.logger.Error("Failed to list orders: %v", err)
		http.Error(w, "Failed to list orders", http.StatusInternalServerError)
		return
	}

	// Convert orders to DTOs
	orderDTOs := make([]dto.OrderDTO, len(orders))
	for i, order := range orders {
		orderDTOs[i] = convertToOrderDTO(order)
	}

	// Create response
	response := dto.OrderListResponse{
		ListResponseDTO: dto.ListResponseDTO[dto.OrderDTO]{
			Success: true,
			Data:    orderDTOs,
			Pagination: dto.PaginationDTO{
				Page:     offset/limit + 1,
				PageSize: limit,
				Total:    len(orderDTOs),
			},
		},
	}

	// Return orders
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateOrderStatus handles updating an order's status (admin only)
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	// Get order ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["orderId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var statusInput struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&statusInput); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update order status
	input := usecase.UpdateOrderStatusInput{
		OrderID: uint(id),
		Status:  entity.OrderStatus(statusInput.Status),
	}

	updatedOrder, err := h.orderUseCase.UpdateOrderStatus(input)
	if err != nil {
		h.logger.Error("Failed to update order status: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert order to DTO
	orderDTO := convertToOrderDTO(updatedOrder)

	// Return updated order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderDTO)
}

// Helper functions to convert between entities and DTOs

func convertToOrderDTO(order *entity.Order) dto.OrderDTO {
	// Convert order items to DTOs
	var items []dto.OrderItemDTO
	if len(order.Items) > 0 {
		items = make([]dto.OrderItemDTO, len(order.Items))
		for i, item := range order.Items {
			items[i] = dto.OrderItemDTO{
				ID:         item.ID,
				OrderID:    order.ID,
				ProductID:  item.ProductID,
				Quantity:   item.Quantity,
				UnitPrice:  money.FromCents(item.Price),
				TotalPrice: money.FromCents(item.Subtotal),
				CreatedAt:  order.CreatedAt,
				UpdatedAt:  order.UpdatedAt,
			}
		}
	}

	// Convert addresses to DTOs
	var shippingAddr *dto.AddressDTO
	if order.ShippingAddr.Street != "" {
		shippingAddr = &dto.AddressDTO{
			AddressLine1: order.ShippingAddr.Street,
			City:         order.ShippingAddr.City,
			State:        order.ShippingAddr.State,
			PostalCode:   order.ShippingAddr.PostalCode,
			Country:      order.ShippingAddr.Country,
		}
	}

	var billingAddr *dto.AddressDTO
	if order.BillingAddr.Street != "" {
		billingAddr = &dto.AddressDTO{
			AddressLine1: order.BillingAddr.Street,
			City:         order.BillingAddr.City,
			State:        order.BillingAddr.State,
			PostalCode:   order.BillingAddr.PostalCode,
			Country:      order.BillingAddr.Country,
		}
	}

	customerDetails := dto.CustomerDetails{
		Email:    order.CustomerDetails.Email,
		Phone:    order.CustomerDetails.Phone,
		FullName: order.CustomerDetails.FullName,
	}

	paymentDetails := dto.PaymentDetails{
		PaymentID: order.PaymentID,
		Provider:  dto.PaymentProvider(order.PaymentProvider),
		Method:    dto.PaymentMethod(order.PaymentMethod),
		Captured:  order.IsCaptured(),
		Refunded:  order.IsRefunded(),
	}

	var discountDetails dto.DiscountDetails
	if order.AppliedDiscount != nil {
		discountDetails = dto.DiscountDetails{
			Code:   order.AppliedDiscount.DiscountCode,
			Amount: money.FromCents(order.AppliedDiscount.DiscountAmount),
		}
	}

	var shippingDetails dto.ShippingDetails
	if order.ShippingMethod != nil {
		shippingDetails = dto.ShippingDetails{
			MethodID: order.ShippingMethodID,
			Method:   order.ShippingMethod.Name,
			Cost:     money.FromCents(order.ShippingCost),
		}
	}

	return dto.OrderDTO{
		ID:              order.ID,
		OrderNumber:     order.OrderNumber,
		UserID:          order.UserID,
		Status:          dto.OrderStatus(order.Status),
		TotalAmount:     money.FromCents(order.TotalAmount),
		FinalAmount:     money.FromCents(order.FinalAmount),
		Currency:        "USD",
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
