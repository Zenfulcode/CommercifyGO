package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
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
		h.logger.Error("Unauthorized access attempt")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
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

	// Get order
	order, err := h.orderUseCase.GetOrderByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get order: %v", err)
		response := dto.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if the user is authorized to view this order
	if order.UserID != userID {
		role, ok := r.Context().Value("role").(string)
		if !ok || role != "admin" {
			h.logger.Error("Unauthorized access to order %d by user %d", order.ID, userID)
			response := dto.ErrorResponse("You are not authorized to view this order")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(response)
			return
		}
	}

	orderDTO := dto.OrderDetailResponse(order)

	// Return order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderDTO)
}

// ListOrders handles listing orders for a user
func (h *OrderHandler) ListOrders(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		h.logger.Error("Unauthorized access attempt")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))

	if page <= 0 {
		page = 1 // Default to page 1
	}

	if pageSize <= 0 {
		page = 10 // Default limit
	}

	// Get orders
	orders, err := h.orderUseCase.GetUserOrders(userID, page, pageSize)
	if err != nil {
		h.logger.Error("Failed to list orders: %v", err)
		// TODO: Add proper error handling
		response := dto.ErrorResponse("Failed to list orders")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create response
	response := dto.OrderSummaryListResponse(orders, page, pageSize, len(orders))

	// Return orders
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListAllOrders handles listing all orders (admin only)
func (h *OrderHandler) ListAllOrders(w http.ResponseWriter, r *http.Request) {
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

	// Get orders by status if provided
	var orders []*entity.Order
	var err error

	if status != "" {
		orders, err = h.orderUseCase.ListOrdersByStatus(entity.OrderStatus(status), page, pageSize)
	} else {
		orders, err = h.orderUseCase.ListAllOrders(page, pageSize)
	}

	if err != nil {
		h.logger.Error("Failed to list orders: %v", err)
		// TODO: Add proper error handling
		response := dto.ErrorResponse("Failed to list orders")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create response
	// TODO: FIX total count logic
	response := dto.OrderSummaryListResponse(orders, page, pageSize, len(orders))

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
		response := dto.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert order to DTO
	orderDTO := dto.OrderUpdateStatusResponse(updatedOrder)

	// Return updated order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderDTO)
}
