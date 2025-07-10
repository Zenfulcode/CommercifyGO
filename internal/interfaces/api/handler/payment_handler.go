package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/interfaces/api/contracts"
)

// PaymentHandler handles payment-related HTTP requests
type PaymentHandler struct {
	orderUseCase *usecase.OrderUseCase
	logger       logger.Logger
}

// NewPaymentHandler creates a new PaymentHandler
func NewPaymentHandler(orderUseCase *usecase.OrderUseCase, logger logger.Logger) *PaymentHandler {
	return &PaymentHandler{
		orderUseCase: orderUseCase,
		logger:       logger,
	}
}

// handleValidationError handles request validation errors
func (h *PaymentHandler) handleValidationError(w http.ResponseWriter, err error, context string) {
	h.logger.Error("Validation error in %s: %v", context, err)
	h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
}

// writeErrorResponse is a helper to write error responses consistently
func (h *PaymentHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := contracts.ErrorResponse(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// GetAvailablePaymentProviders returns a list of available payment providers
func (h *PaymentHandler) GetAvailablePaymentProviders(w http.ResponseWriter, r *http.Request) {
	// Check for currency parameter
	currency := r.URL.Query().Get("currency")

	var providers []service.PaymentProvider
	if currency != "" {
		// Get providers that support the specific currency
		providers = h.orderUseCase.GetAvailablePaymentProvidersForCurrency(currency)
	} else {
		// Get all available payment providers
		providers = h.orderUseCase.GetAvailablePaymentProviders()
	}

	// Return providers
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providers)
}

// CapturePayment handles capturing an authorized payment
func (h *PaymentHandler) CapturePayment(w http.ResponseWriter, r *http.Request) {
	// Get payment ID from URL
	vars := mux.Vars(r)
	paymentID := vars["paymentId"]
	if paymentID == "" {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	var request contracts.CapturePaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.handleValidationError(w, err, "CapturePayment")
		return
	}

	// Validate input - either amount or is_full must be specified
	if !request.IsFull && request.Amount <= 0 {
		h.writeErrorResponse(w, http.StatusBadRequest, "Amount must be greater than zero when is_full is false")
		return
	}

	// If both amount and is_full are specified, prioritize is_full
	if request.IsFull && request.Amount > 0 {
		h.logger.Info("Both amount and is_full specified for payment %s, using is_full=true", paymentID)
	}

	// Capture payment
	var err error
	if request.IsFull {
		// For full capture, we need to get the order first to determine the full amount
		order, orderErr := h.orderUseCase.GetOrderByPaymentID(paymentID)
		if orderErr != nil {
			h.logger.Error("Failed to get order for payment %s: %v", paymentID, orderErr)
			http.Error(w, "Order not found for payment ID", http.StatusNotFound)
			return
		}
		err = h.orderUseCase.CapturePayment(paymentID, order.FinalAmount)
	} else {
		err = h.orderUseCase.CapturePayment(paymentID, money.ToCents(request.Amount))
	}
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to capture payment: "+err.Error())
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.SuccessResponseMessage("Payment captured successfully"))
}

// CancelPayment handles cancelling a payment
func (h *PaymentHandler) CancelPayment(w http.ResponseWriter, r *http.Request) {
	// Get payment ID from URL
	vars := mux.Vars(r)
	paymentID := vars["paymentId"]
	if paymentID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	// Cancel payment
	err := h.orderUseCase.CancelPayment(paymentID)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to cancel payment: "+err.Error())
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.SuccessResponseMessage("Payment cancelled successfully"))
}

// RefundPayment handles refunding a payment
func (h *PaymentHandler) RefundPayment(w http.ResponseWriter, r *http.Request) {
	// Get payment ID from URL
	vars := mux.Vars(r)
	paymentID := vars["paymentId"]
	if paymentID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	var request contracts.RefundPaymentRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.handleValidationError(w, err, "RefundPayment")
		return
	}

	// Validate input - either amount or is_full must be specified
	if !request.IsFull && request.Amount <= 0 {
		h.writeErrorResponse(w, http.StatusBadRequest, "Amount must be greater than zero when is_full is false")
		return
	}

	// If both amount and is_full are specified, prioritize is_full
	if request.IsFull && request.Amount > 0 {
		h.logger.Info("Both amount and is_full specified for payment %s, using is_full=true", paymentID)
	}

	// Refund payment
	var err error
	if request.IsFull {
		// For full refund, we need to get the order first to determine the full amount
		order, orderErr := h.orderUseCase.GetOrderByPaymentID(paymentID)
		if orderErr != nil {
			h.logger.Error("Failed to get order for payment %s: %v", paymentID, orderErr)
			http.Error(w, "Order not found for payment ID", http.StatusNotFound)
			return
		}
		err = h.orderUseCase.RefundPayment(paymentID, order.FinalAmount)
	} else {
		err = h.orderUseCase.RefundPayment(paymentID, money.ToCents(request.Amount))
	}
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to refund payment: "+err.Error())
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.SuccessResponseMessage("Payment refunded successfully"))
}

// ForceApproveMobilePayPayment handles force approving a MobilePay payment (admin only)
func (h *PaymentHandler) ForceApproveMobilePayPayment(w http.ResponseWriter, r *http.Request) {
	// Get payment ID from URL
	vars := mux.Vars(r)
	paymentID := vars["paymentId"]
	if paymentID == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Invalid payment ID")
		return
	}

	// Parse request body
	var input struct {
		PhoneNumber string `json:"phone_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.handleValidationError(w, err, "ForceApproveMobilePayPayment")
		return
	}

	// Validate phone number
	if input.PhoneNumber == "" {
		h.writeErrorResponse(w, http.StatusBadRequest, "Phone number is required")
		return
	}

	// Force approve payment
	err := h.orderUseCase.ForceApproveMobilePayPayment(paymentID, input.PhoneNumber)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "Failed to force approve payment: "+err.Error())
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(contracts.SuccessResponseMessage("Payment force approved successfully"))
}
