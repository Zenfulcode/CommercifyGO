package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// DiscountHandler handles discount-related HTTP requests
type DiscountHandler struct {
	discountUseCase *usecase.DiscountUseCase
	orderUseCase    *usecase.OrderUseCase
	logger          logger.Logger
}

// NewDiscountHandler creates a new DiscountHandler
func NewDiscountHandler(discountUseCase *usecase.DiscountUseCase, orderUseCase *usecase.OrderUseCase, logger logger.Logger) *DiscountHandler {
	return &DiscountHandler{
		discountUseCase: discountUseCase,
		orderUseCase:    orderUseCase,
		logger:          logger,
	}
}

// CreateDiscount handles creating a new discount (admin only)
func (h *DiscountHandler) CreateDiscount(w http.ResponseWriter, r *http.Request) {
	var input usecase.CreateDiscountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	discount, err := h.discountUseCase.CreateDiscount(input)
	if err != nil {
		h.logger.Error("Failed to create discount: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(discount)
}

// GetDiscount handles getting a discount by ID (admin only)
func (h *DiscountHandler) GetDiscount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["discountId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid discount ID", http.StatusBadRequest)
		return
	}

	discount, err := h.discountUseCase.GetDiscountByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get discount: %v", err)
		http.Error(w, "Discount not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discount)
}

// UpdateDiscount handles updating a discount (admin only)
func (h *DiscountHandler) UpdateDiscount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["discountId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid discount ID", http.StatusBadRequest)
		return
	}

	var input usecase.UpdateDiscountInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	discount, err := h.discountUseCase.UpdateDiscount(uint(id), input)
	if err != nil {
		h.logger.Error("Failed to update discount: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discount)
}

// DeleteDiscount handles deleting a discount (admin only)
func (h *DiscountHandler) DeleteDiscount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["discountId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid discount ID", http.StatusBadRequest)
		return
	}

	if err := h.discountUseCase.DeleteDiscount(uint(id)); err != nil {
		h.logger.Error("Failed to delete discount: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// ListDiscounts handles listing all discounts (admin only)
func (h *DiscountHandler) ListDiscounts(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10 // Default limit
	}

	discounts, err := h.discountUseCase.ListDiscounts(offset, limit)
	if err != nil {
		h.logger.Error("Failed to list discounts: %v", err)
		http.Error(w, "Failed to list discounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discounts)
}

// ListActiveDiscounts handles listing active discounts (public)
func (h *DiscountHandler) ListActiveDiscounts(w http.ResponseWriter, r *http.Request) {
	offset, _ := strconv.Atoi(r.URL.Query().Get("offset"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10 // Default limit
	}

	discounts, err := h.discountUseCase.ListActiveDiscounts(offset, limit)
	if err != nil {
		h.logger.Error("Failed to list active discounts: %v", err)
		http.Error(w, "Failed to list discounts", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(discounts)
}

// ApplyDiscountToOrder handles applying a discount to an order
func (h *DiscountHandler) ApplyDiscountToOrder(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get order ID from URL
	vars := mux.Vars(r)
	orderID, err := strconv.ParseUint(vars["orderId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var input struct {
		DiscountCode string `json:"discount_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the order to verify ownership
	order, err := h.orderUseCase.GetOrderByID(uint(orderID))
	if err != nil {
		h.logger.Error("Failed to get order: %v", err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	role, _ := r.Context().Value("role").(string)

	// Check if the user is authorized to apply discount to this order
	if order.UserID != userID && role != string(entity.RoleAdmin) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Check if order is in a state where discounts can be applied
	if order.Status != entity.OrderStatusPending {
		http.Error(w, "Discount can only be applied to pending orders", http.StatusBadRequest)
		return
	}

	// Apply discount to order
	discountInput := usecase.ApplyDiscountToOrderInput{
		OrderID:      uint(orderID),
		DiscountCode: input.DiscountCode,
	}

	updatedOrder, err := h.discountUseCase.ApplyDiscountToOrder(discountInput, order)
	if err != nil {
		h.logger.Error("Failed to apply discount: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return updated order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrder)
}

// RemoveDiscountFromOrder handles removing a discount from an order
func (h *DiscountHandler) RemoveDiscountFromOrder(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get order ID from URL
	vars := mux.Vars(r)
	orderID, err := strconv.ParseUint(vars["orderId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Get the order to verify ownership
	order, err := h.orderUseCase.GetOrderByID(uint(orderID))
	if err != nil {
		h.logger.Error("Failed to get order: %v", err)
		http.Error(w, "Order not found", http.StatusNotFound)
		return
	}

	role, _ := r.Context().Value("role").(string)

	// Check if the user is authorized to remove discount from this order
	if order.UserID != userID && role != string(entity.RoleAdmin) {
		http.Error(w, "Unauthorized", http.StatusForbidden)
		return
	}

	// Check if order is in a state where discounts can be removed
	if order.Status != entity.OrderStatusPending {
		http.Error(w, "Discount can only be removed from pending orders", http.StatusBadRequest)
		return
	}

	// Check if order has a discount applied
	if order.AppliedDiscount == nil {
		http.Error(w, "No discount applied to this order", http.StatusBadRequest)
		return
	}

	// Remove discount from order
	h.discountUseCase.RemoveDiscountFromOrder(order)

	// Return updated order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// ValidateDiscountCode handles validating a discount code without applying it
func (h *DiscountHandler) ValidateDiscountCode(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var input struct {
		DiscountCode string `json:"discount_code"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get discount by code
	discount, err := h.discountUseCase.GetDiscountByCode(input.DiscountCode)
	if err != nil {
		http.Error(w, "Invalid discount code", http.StatusBadRequest)
		return
	}

	// Check if discount is valid
	if !discount.IsValid() {
		response := map[string]interface{}{
			"valid":  false,
			"reason": "Discount is not valid (expired, inactive, or usage limit reached)",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Return discount details
	response := map[string]interface{}{
		"valid": true,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
