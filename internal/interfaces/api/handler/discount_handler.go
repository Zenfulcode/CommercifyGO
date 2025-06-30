package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/interfaces/api/contracts"
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
	var req contracts.CreateDiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert DTO to usecase input
	input := req.ToUseCaseInput()

	discount, err := h.discountUseCase.CreateDiscount(input)
	if err != nil {
		h.logger.Error("Failed to create discount: %v", err)
		response := contracts.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := contracts.DiscountCreateResponse(discount.ToDiscountDTO())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetDiscount handles getting a discount by ID (admin only)
func (h *DiscountHandler) GetDiscount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["discountId"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid discount ID: %v", err)
		http.Error(w, "Invalid discount ID", http.StatusBadRequest)
		return
	}

	discount, err := h.discountUseCase.GetDiscountByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get discount: %v", err)
		response := contracts.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := contracts.DiscountRetrieveResponse(discount.ToDiscountDTO())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateDiscount handles updating a discount (admin only)
func (h *DiscountHandler) UpdateDiscount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["discountId"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid discount ID: %v", err)
		http.Error(w, "Invalid discount ID", http.StatusBadRequest)
		return
	}

	var req contracts.UpdateDiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert DTO to usecase input
	input := req.ToUseCaseInput()

	discount, err := h.discountUseCase.UpdateDiscount(uint(id), input)
	if err != nil {
		h.logger.Error("Failed to update discount: %v", err)
		response := contracts.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := contracts.DiscountUpdateResponse(discount.ToDiscountDTO())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteDiscount handles deleting a discount (admin only)
func (h *DiscountHandler) DeleteDiscount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["discountId"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid discount ID: %v", err)
		http.Error(w, "Invalid discount ID", http.StatusBadRequest)
		return
	}

	if err := h.discountUseCase.DeleteDiscount(uint(id)); err != nil {
		h.logger.Error("Failed to delete discount: %v", err)
		response := contracts.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := contracts.DiscountDeleteResponse()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
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
		response := contracts.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Calculate page information
	page := (offset / limit) + 1
	if limit == 0 {
		page = 1
	}

	// TODO: Get total count from the use case if available
	response := contracts.DiscountListResponse(discounts, len(discounts), page, limit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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
		response := contracts.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Calculate page information
	page := (offset / limit) + 1
	if limit == 0 {
		page = 1
	}

	response := contracts.DiscountListResponse(discounts, len(discounts), page, limit)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ApplyDiscountToOrder handles applying a discount to an order
func (h *DiscountHandler) ApplyDiscountToOrder(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		h.logger.Error("Unauthorized access: user ID not found in context")
		response := contracts.ErrorResponse("Unauthorized access")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get order ID from URL
	vars := mux.Vars(r)
	orderID, err := strconv.ParseUint(vars["orderId"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid order ID: %v", err)
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var req contracts.ApplyDiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get the order to verify ownership
	order, err := h.orderUseCase.GetOrderByID(uint(orderID))
	if err != nil {
		h.logger.Error("Failed to get order: %v", err)
		response := contracts.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	role, _ := r.Context().Value("role").(string)

	// Check if the user is authorized to apply discount to this order
	if order.UserID != userID && role != string(entity.RoleAdmin) {
		h.logger.Error("Unauthorized access: user does not own the order")
		response := contracts.ErrorResponse("Unauthorized access: user does not own the order")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if order is in a state where discounts can be applied
	if order.Status != entity.OrderStatusPending {
		h.logger.Error("Discount can only be applied to pending orders")
		response := contracts.ErrorResponse("Discount can only be applied to pending orders")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Apply discount to order
	discountInput := usecase.ApplyDiscountToOrderInput{
		OrderID:      uint(orderID),
		DiscountCode: req.DiscountCode,
	}

	updatedOrder, err := h.discountUseCase.ApplyDiscountToOrder(discountInput, order)
	if err != nil {
		h.logger.Error("Failed to apply discount: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// TODO: Order Response DTO

	// Return updated order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedOrder)
}

// RemoveDiscountFromOrder handles removing a discount from an order
func (h *DiscountHandler) RemoveDiscountFromOrder(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		h.logger.Error("Unauthorized access: user ID not found in context")
		response := contracts.ErrorResponse("Unauthorized access")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get order ID from URL
	vars := mux.Vars(r)
	orderID, err := strconv.ParseUint(vars["orderId"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid order ID: %v", err)
		http.Error(w, "Invalid order ID", http.StatusBadRequest)
		return
	}

	// Get the order to verify ownership
	order, err := h.orderUseCase.GetOrderByID(uint(orderID))
	if err != nil {
		h.logger.Error("Failed to get order: %v", err)
		response := contracts.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	role, _ := r.Context().Value("role").(string)

	// Check if the user is authorized to remove discount from this order
	if order.UserID != userID && role != string(entity.RoleAdmin) {
		response := contracts.ErrorResponse("Unauthorized access: user does not own the order")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if order is in a state where discounts can be removed
	if order.Status != entity.OrderStatusPending {
		response := contracts.ErrorResponse("Discount can only be removed from pending orders")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if order has a discount applied
	if order.AppliedDiscount == nil {
		h.logger.Error("No discount applied to this order")
		response := contracts.ErrorResponse("No discount applied to this order")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Remove discount from order
	h.discountUseCase.RemoveDiscountFromOrder(order)

	// TODO: Order Response DTO

	// Return updated order
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(order)
}

// ValidateDiscountCode handles validating a discount code without applying it
func (h *DiscountHandler) ValidateDiscountCode(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req contracts.ValidateDiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Get discount by code
	discount, err := h.discountUseCase.GetDiscountByCode(req.DiscountCode)
	if err != nil {
		response := contracts.ValidateDiscountResponse{
			Valid:  false,
			Reason: "Invalid discount code",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Check if discount is valid
	if !discount.IsValid() {
		response := contracts.ValidateDiscountResponse{
			Valid:  false,
			Reason: "Discount is not valid (expired, inactive, or usage limit reached)",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
		return
	}

	// Return discount details
	response := contracts.ValidateDiscountResponse{
		Valid:            true,
		DiscountID:       discount.ID,
		Code:             discount.Code,
		Type:             string(discount.Type),
		Method:           string(discount.Method),
		Value:            discount.Value,
		MinOrderValue:    float64(discount.MinOrderValue) / 100,    // Convert from cents to dollars
		MaxDiscountValue: float64(discount.MaxDiscountValue) / 100, // Convert from cents to dollars
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
