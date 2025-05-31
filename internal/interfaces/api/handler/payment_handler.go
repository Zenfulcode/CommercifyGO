package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
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

	// Parse request body
	var input struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate amount
	if input.Amount <= 0 {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	// Capture payment
	err := h.orderUseCase.CapturePayment(paymentID, money.ToCents(input.Amount))
	if err != nil {
		h.logger.Error("Failed to capture payment: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Payment captured successfully",
	})
}

// CancelPayment handles cancelling a payment
func (h *PaymentHandler) CancelPayment(w http.ResponseWriter, r *http.Request) {
	// Get payment ID from URL
	vars := mux.Vars(r)
	paymentID := vars["paymentId"]
	if paymentID == "" {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	// Cancel payment
	err := h.orderUseCase.CancelPayment(paymentID)
	if err != nil {
		h.logger.Error("Failed to cancel payment: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Payment cancelled successfully",
	})
}

// RefundPayment handles refunding a payment
func (h *PaymentHandler) RefundPayment(w http.ResponseWriter, r *http.Request) {
	// Get payment ID from URL
	vars := mux.Vars(r)
	paymentID := vars["paymentId"]
	if paymentID == "" {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var input struct {
		Amount float64 `json:"amount"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate amount
	if input.Amount <= 0 {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	// Refund payment
	err := h.orderUseCase.RefundPayment(paymentID, money.ToCents(input.Amount))
	if err != nil {
		h.logger.Error("Failed to refund payment: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Payment refunded successfully",
	})
}

// ForceApproveMobilePayPayment handles force approving a MobilePay payment (admin only)
func (h *PaymentHandler) ForceApproveMobilePayPayment(w http.ResponseWriter, r *http.Request) {
	// Get payment ID from URL
	vars := mux.Vars(r)
	paymentID := vars["paymentId"]
	if paymentID == "" {
		http.Error(w, "Invalid payment ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var input struct {
		PhoneNumber string `json:"phone_number"`
	}
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate phone number
	if input.PhoneNumber == "" {
		http.Error(w, "Phone number is required", http.StatusBadRequest)
		return
	}

	// Force approve payment
	err := h.orderUseCase.ForceApproveMobilePayPayment(paymentID, input.PhoneNumber)
	if err != nil {
		h.logger.Error("Failed to force approve payment: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return success
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status":  "success",
		"message": "Payment force approved successfully",
	})
}
