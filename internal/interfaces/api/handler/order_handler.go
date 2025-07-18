package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/interfaces/api/contracts"
	"github.com/zenfulcode/commercify/internal/interfaces/api/middleware"
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
	// Get user ID from context (optional for checkout session access)
	userID, isAuthenticated := r.Context().Value(middleware.UserIDKey).(uint)
	// Get order ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["orderId"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid order ID: %v", err)
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Parse query parameters for includes
	includePaymentTransactions := r.URL.Query().Get("includePaymentTransactions") == "true"
	includeItems := r.URL.Query().Get("includeItems") != "false" // Default to true for backward compatibility

	// Get order
	order, err := h.orderUseCase.GetOrderByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get order: %v", err)
		response := contracts.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check authorization: user owns the order, admin, or checkout session matches
	authorized := false

	// Check if authenticated user owns the order or is admin
	if isAuthenticated {
		if order.UserID != nil && *order.UserID == userID {
			authorized = true
		} else {
			// Check if user is admin
			role, ok := r.Context().Value(middleware.RoleKey).(string)
			if ok && role == string(entity.RoleAdmin) {
				authorized = true
			}
		}
	}

	// If not authorized by user auth, check checkout session cookie
	if !authorized {
		cookie, err := r.Cookie(common.CheckoutSessionCookie)
		if err == nil && cookie.Value != "" && cookie.Value == order.CheckoutSessionID {
			authorized = true
			h.logger.Info("Order %d accessed via checkout session: %s", order.ID, cookie.Value)
		}
	}

	if !authorized {
		h.logger.Error("Unauthorized access to order %d", order.ID)
		response := contracts.ErrorResponse("You are not authorized to view this order")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create order DTO with conditional includes
	options := entity.OrderDetailOptions{
		IncludePaymentTransactions: includePaymentTransactions,
		IncludeItems:               includeItems,
	}
	orderDTO := contracts.OrderDetailResponse(order.ToOrderDetailsDTOWithOptions(options))

	// Return order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderDTO)
}

// ListOrders handles listing orders for a user
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		h.logger.Error("Unauthorized access attempt")
		response := contracts.ErrorResponse("Unauthorized")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	if page <= 0 {
		page = 1 // Default to page 1
	}

	if pageSize <= 0 {
		pageSize = 10 // Default limit
	}

	// Calculate offset for pagination
	offset := (page - 1) * pageSize

	// Get orders
	orders, err := h.orderUseCase.GetUserOrders(userID, offset, pageSize)
	if err != nil {
		h.logger.Error("Failed to list orders: %v", err)
		// TODO: Add proper error handling
		response := contracts.ErrorResponse("Failed to list orders")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create response
	response := contracts.OrderSummaryListResponse(orders, page, pageSize, len(orders))

	// Return orders
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListAllOrders handles listing all orders (admin only)
func (h *OrderHandler) ListAllOrders(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	_, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		h.logger.Error("Unauthorized access attempt")
		response := contracts.ErrorResponse("Unauthorized")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	status := r.URL.Query().Get("status")

	if page <= 0 {
		page = 1 // Default to page 1
	}
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}

	// Calculate offset for pagination
	offset := (page - 1) * pageSize

	// Get orders by status if provided
	var orders []*entity.Order
	var err error

	if status != "" {
		orders, err = h.orderUseCase.ListOrdersByStatus(entity.OrderStatus(status), offset, pageSize)
	} else {
		orders, err = h.orderUseCase.ListAllOrders(offset, pageSize)
	}

	if err != nil {
		h.logger.Error("Failed to list orders: %v", err)
		// TODO: Add proper error handling
		response := contracts.ErrorResponse("Failed to list orders")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create response
	// TODO: FIX total count logic
	response := contracts.OrderSummaryListResponse(orders, page, pageSize, len(orders))

	// Return orders
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateOrderStatus handles updating an order's status (admin only)
func (h *OrderHandler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		h.logger.Error("Unauthorized access attempt")
		response := contracts.ErrorResponse("Unauthorized")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get order ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["orderId"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid order ID: %v", err)
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var statusInput struct {
		Status string `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&statusInput); err != nil {
		h.logger.Error("Failed to decode request body: %v", err)
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
		response := contracts.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert order to DTO
	response := contracts.OrderUpdateStatusResponse(*updatedOrder.ToOrderSummaryDTO())

	// Return updated order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateOrderStatusWithTracking handles updating an order's status with tracking information (admin only)
func (h *OrderHandler) UpdateOrderStatusWithTracking(w http.ResponseWriter, r *http.Request) {
	_, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok {
		h.logger.Error("Unauthorized access attempt")
		response := contracts.ErrorResponse("Unauthorized")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get order ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["orderId"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid order ID: %v", err)
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var statusInput struct {
		Status         string `json:"status"`
		TrackingNumber string `json:"tracking_number,omitempty"`
		TrackingURL    string `json:"tracking_url,omitempty"`
	}
	if err := json.NewDecoder(r.Body).Decode(&statusInput); err != nil {
		h.logger.Error("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Update order status with tracking
	input := usecase.UpdateOrderStatusWithTrackingInput{
		OrderID:        uint(id),
		Status:         entity.OrderStatus(statusInput.Status),
		TrackingNumber: statusInput.TrackingNumber,
		TrackingURL:    statusInput.TrackingURL,
	}

	updatedOrder, err := h.orderUseCase.UpdateOrderStatusWithTracking(input)
	if err != nil {
		h.logger.Error("Failed to update order status: %v", err)
		response := contracts.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert order to DTO
	response := contracts.OrderUpdateStatusResponse(*updatedOrder.ToOrderSummaryDTO())

	// Return updated order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
