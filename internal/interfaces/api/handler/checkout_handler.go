package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/dto"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// CheckoutHandler handles checkout-related HTTP requests
type CheckoutHandler struct {
	checkoutUseCase *usecase.CheckoutUseCase
	orderUseCase    *usecase.OrderUseCase
	logger          logger.Logger
}

// NewCheckoutHandler creates a new CheckoutHandler
func NewCheckoutHandler(checkoutUseCase *usecase.CheckoutUseCase, orderUseCase *usecase.OrderUseCase, logger logger.Logger) *CheckoutHandler {
	return &CheckoutHandler{
		checkoutUseCase: checkoutUseCase,
		orderUseCase:    orderUseCase,
		logger:          logger,
	}
}

// getCheckoutSessionID gets or creates a checkout session ID
func (h *CheckoutHandler) getCheckoutSessionID(w http.ResponseWriter, r *http.Request) string {
	// Check if checkout session cookie exists
	cookie, err := r.Cookie(common.CheckoutSessionCookie)
	if err == nil && cookie.Value != "" {
		return cookie.Value
	}

	// Create new checkout session ID if none exists
	sessionID := uuid.New().String()
	http.SetCookie(w, &http.Cookie{
		Name:     common.CheckoutSessionCookie,
		Value:    sessionID,
		Path:     "/",
		MaxAge:   common.CheckoutSessionMaxAge,
		HttpOnly: true,
		Secure:   r.TLS != nil, // Set secure flag if connection is HTTPS
		SameSite: http.SameSiteLaxMode,
	})

	return sessionID
}

// GetCheckout handles getting a user's checkout
func (h *CheckoutHandler) GetCheckout(w http.ResponseWriter, r *http.Request) {
	// Always get checkout session ID, needed for all checkouts
	checkoutSessionID := h.getCheckoutSessionID(w, r)

	// Check for optional currency parameter
	currency := r.URL.Query().Get("currency")

	var checkout *entity.Checkout
	var err error

	if currency != "" {
		checkout, err = h.checkoutUseCase.GetOrCreateCheckoutBySessionIDWithCurrency(checkoutSessionID, currency)
	} else {
		checkout, err = h.checkoutUseCase.GetOrCreateCheckoutBySessionID(checkoutSessionID)
	}

	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		response := dto.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.CreateCheckoutResponse(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// AddToCheckout handles adding an item to the checkout
func (h *CheckoutHandler) AddToCheckout(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.AddToCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Invalid request body: %v", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Always get checkout session ID, needed for all checkouts
	checkoutSessionID := h.getCheckoutSessionID(w, r)

	// print request and session ID for debugging
	h.logger.Debug("AddToCheckout request: %+v", request)
	fmt.Printf("Checkout session ID: %s\n", checkoutSessionID)

	// Try to find checkout by checkout session ID first
	var checkout *entity.Checkout
	var err error

	if request.Currency != "" {
		checkout, err = h.checkoutUseCase.GetOrCreateCheckoutBySessionIDWithCurrency(checkoutSessionID, request.Currency)
	} else {
		checkout, err = h.checkoutUseCase.GetOrCreateCheckoutBySessionID(checkoutSessionID)
	}
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert DTO to usecase input
	checkoutInput := usecase.CheckoutInput{
		SKU:      request.SKU,
		Quantity: request.Quantity,
	}

	// Add item to checkout
	checkout, err = h.checkoutUseCase.AddItemToCheckout(checkout.ID, checkoutInput)

	if err != nil {
		h.logger.Error("Failed to add to checkout: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.CreateCheckoutResponse(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateCheckoutItem handles updating an item in the checkout
func (h *CheckoutHandler) UpdateCheckoutItem(w http.ResponseWriter, r *http.Request) {
	// Get SKU from URL path
	vars := mux.Vars(r)
	sku := vars["sku"]
	if sku == "" {
		h.logger.Error("SKU is required in URL path")
		http.Error(w, "SKU is required in URL path", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request dto.UpdateCheckoutItemRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Invalid request body: %v", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	checkoutSessionID := h.getCheckoutSessionID(w, r)

	checkout, err := h.checkoutUseCase.GetOrCreateCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert path parameter and request body to usecase input
	updateInput := usecase.UpdateCheckoutItemInput{
		SKU:      sku,
		Quantity: request.Quantity,
	}

	// Update item in checkout using the new usecase method
	checkout, err = h.checkoutUseCase.UpdateCheckoutItemBySKU(checkout.ID, updateInput)

	if err != nil {
		h.logger.Error("Failed to update checkout item: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.CreateCheckoutResponse(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RemoveFromCheckout handles removing an item from the checkout
func (h *CheckoutHandler) RemoveFromCheckout(w http.ResponseWriter, r *http.Request) {
	// Get SKU from URL path
	vars := mux.Vars(r)
	sku := vars["sku"]
	if sku == "" {
		h.logger.Error("SKU is required in URL path")
		http.Error(w, "SKU is required in URL path", http.StatusBadRequest)
		return
	}

	checkoutSessionID := h.getCheckoutSessionID(w, r)

	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		response := dto.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert SKU to usecase input
	removeInput := usecase.RemoveItemInput{
		SKU: sku,
	}

	// Remove item from checkout using the new usecase method
	checkout, err = h.checkoutUseCase.RemoveItemBySKU(checkout.ID, removeInput)

	if err != nil {
		h.logger.Error("Failed to remove item from checkout: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.CreateCheckoutResponse(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ClearCheckout handles emptying the checkout
func (h *CheckoutHandler) ClearCheckout(w http.ResponseWriter, r *http.Request) {
	checkoutSessionID := h.getCheckoutSessionID(w, r)

	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		response := dto.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	checkout.Clear()
	checkout, err = h.checkoutUseCase.UpdateCheckout(checkout)

	if err != nil {
		h.logger.Error("Failed to clear checkout: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.CreateCheckoutResponse(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SetShippingAddress handles setting the shipping address for a checkout
func (h *CheckoutHandler) SetShippingAddress(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.SetShippingAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to parse shipping address request: %v", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	checkoutSessionID := h.getCheckoutSessionID(w, r)

	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	address := entity.Address{
		Street:     request.AddressLine1,
		City:       request.City,
		State:      request.State,
		PostalCode: request.PostalCode,
		Country:    request.Country,
	}

	checkout.SetShippingAddress(address)
	checkout, err = h.checkoutUseCase.UpdateCheckout(checkout)

	if err != nil {
		h.logger.Error("Failed to set shipping address: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.CreateCheckoutResponse(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SetBillingAddress handles setting the billing address for a checkout
func (h *CheckoutHandler) SetBillingAddress(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.SetBillingAddressRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to parse billing address request: %v", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	checkoutSessionID := h.getCheckoutSessionID(w, r)
	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		response := dto.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert DTO to address entity
	address := entity.Address{
		Street:     request.AddressLine1,
		City:       request.City,
		State:      request.State,
		PostalCode: request.PostalCode,
		Country:    request.Country,
	}

	checkout.SetBillingAddress(address)
	checkout, err = h.checkoutUseCase.UpdateCheckout(checkout)

	if err != nil {
		h.logger.Error("Failed to set billing address: %v", err)
		response := dto.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.CreateCheckoutResponse(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SetCustomerDetails handles setting the customer details for a checkout
func (h *CheckoutHandler) SetCustomerDetails(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.SetCustomerDetailsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to parse customer details request: %v", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	checkoutSessionID := h.getCheckoutSessionID(w, r)

	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert DTO to customer details entity
	customerDetails := entity.CustomerDetails{
		Email:    request.Email,
		Phone:    request.Phone,
		FullName: request.FullName,
	}

	checkout.SetCustomerDetails(customerDetails)
	checkout, err = h.checkoutUseCase.UpdateCheckout(checkout)

	if err != nil {
		h.logger.Error("Failed to set customer details: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.CreateCheckoutResponse(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SetShippingMethod handles setting the shipping method for a checkout
func (h *CheckoutHandler) SetShippingMethod(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.SetShippingMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to parse shipping method request: %v", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	checkoutSessionID := h.getCheckoutSessionID(w, r)
	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	checkout, err = h.checkoutUseCase.SetShippingMethod(checkout, request.ShippingMethodID)

	if err != nil {
		h.logger.Error("Failed to set shipping method: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.CreateCheckoutResponse(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ApplyDiscount handles applying a discount code to a checkout
func (h *CheckoutHandler) ApplyDiscount(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.ApplyDiscountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to parse discount code request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	checkoutSessionID := h.getCheckoutSessionID(w, r)
	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	checkout, err = h.checkoutUseCase.ApplyDiscountCode(checkout, request.DiscountCode)

	if err != nil {
		h.logger.Error("Failed to apply discount: %v", err)
		response := dto.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.CreateCheckoutResponse(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// RemoveDiscount handles removing a discount from a checkout
func (h *CheckoutHandler) RemoveDiscount(w http.ResponseWriter, r *http.Request) {
	checkoutSessionID := h.getCheckoutSessionID(w, r)
	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	checkout.ApplyDiscount(nil)

	checkout, err = h.checkoutUseCase.UpdateCheckout(checkout)

	if err != nil {
		h.logger.Error("Failed to remove discount: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.CreateCheckoutResponse(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SetCurrency handles changing the currency for a checkout
func (h *CheckoutHandler) SetCurrency(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.SetCurrencyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to parse currency change request: %v", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validate currency code
	if request.Currency == "" {
		h.logger.Error("Currency code is required")
		response := dto.ErrorResponse("Currency code is required")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	checkoutSessionID := h.getCheckoutSessionID(w, r)
	checkout, err := h.checkoutUseCase.ChangeCurrencyBySessionID(checkoutSessionID, request.Currency)
	if err != nil {
		h.logger.Error("Failed to change checkout currency: %v", err)
		response := dto.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.CreateCheckoutResponse(checkout)

	// Return updated checkout
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// CompleteOrder handles converting a checkout to an order
func (h *CheckoutHandler) CompleteOrder(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var paymentInput dto.CompleteCheckoutRequest
	if err := json.NewDecoder(r.Body).Decode(&paymentInput); err != nil {
		h.logger.Error("Failed to parse checkout completion request: %v", err)
		http.Error(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Get checkout session ID
	checkoutSessionID := h.getCheckoutSessionID(w, r)

	h.logger.Info("Converting checkout to order. CheckoutSessionID: %s", checkoutSessionID)

	var order *entity.Order
	var err error

	// Try to find checkout by checkout session ID first
	checkout, err := h.checkoutUseCase.GetCheckoutBySessionID(checkoutSessionID)
	if err != nil {
		h.logger.Error("Failed to get checkout with session ID %s: %v", checkoutSessionID, err)

		errResponse := dto.ErrorResponse("Checkout not found. Please create a checkout first.")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	// Check if checkout has items
	if checkout == nil || len(checkout.Items) == 0 {
		h.logger.Error("Checkout %s has no items", checkoutSessionID)

		errResponse := dto.ErrorResponse("Checkout is empty. Please add items to the checkout before completing.")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	// If checkout exists for this session, convert it to order
	order, err = h.checkoutUseCase.CreateOrderFromCheckout(checkout.ID)
	if err != nil {
		h.logger.Error("Failed to convert checkout to order: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate payment data
	if paymentInput.PaymentData.CardDetails == nil && paymentInput.PaymentData.PhoneNumber == "" {
		h.logger.Error("Missing payment data: both CardDetails and PhoneNumber are empty")
		response := dto.ErrorResponse("Payment data is required. Please provide either card details or a phone number for wallet payments.")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Validate that the payment provider is specified
	if paymentInput.PaymentProvider == "" {
		h.logger.Error("Missing payment provider")
		response := dto.ErrorResponse("Payment provider is required. Please specify a payment provider.")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Determine the payment method based on provided data
	paymentMethod := service.PaymentMethodWallet
	if paymentInput.PaymentData.CardDetails != nil {
		paymentMethod = service.PaymentMethodCreditCard
	}

	processInput := usecase.ProcessPaymentInput{
		PaymentProvider: service.PaymentProviderType(paymentInput.PaymentProvider),
		PaymentMethod:   paymentMethod,
		PhoneNumber:     paymentInput.PaymentData.PhoneNumber,
	}

	// Only add card details if they were provided
	if paymentInput.PaymentData.CardDetails != nil {
		processInput.CardDetails = &service.CardDetails{
			CardNumber:     paymentInput.PaymentData.CardDetails.CardNumber,
			ExpiryMonth:    paymentInput.PaymentData.CardDetails.ExpiryMonth,
			ExpiryYear:     paymentInput.PaymentData.CardDetails.ExpiryYear,
			CVV:            paymentInput.PaymentData.CardDetails.CVV,
			CardholderName: paymentInput.PaymentData.CardDetails.CardholderName,
			Token:          paymentInput.PaymentData.CardDetails.Token,
		}
	}

	// Process payment
	h.logger.Debug("Processing payment for order %d with provider %s and method %s",
		order.ID, processInput.PaymentProvider, processInput.PaymentMethod)

	processedOrder, err := h.checkoutUseCase.ProcessPayment(order, processInput)
	if err != nil {
		// print order
		h.logger.Debug("Order details: %+v", order)

		h.orderUseCase.FailOrder(order)

		h.logger.Error("Failed to process payment for order %d: %v", order.ID, err)

		// Return a more informative error to the client
		errResponse := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errResponse)
		return
	}

	// Create response
	response := dto.CreateCompleteCheckoutResponse(processedOrder)

	// Return created order
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// ListAdminCheckouts handles listing all checkouts (admin only)
func (h *CheckoutHandler) ListAdminCheckouts(w http.ResponseWriter, r *http.Request) {
	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("pageSize"))
	status := r.URL.Query().Get("status")

	if page <= 0 {
		page = 1 // Default to page 1
	}

	if pageSize <= 0 {
		pageSize = 10 // Default limit
	}

	// Get checkouts by status if provided
	var checkouts []*entity.Checkout
	var err error

	if status != "" {
		checkouts, err = h.checkoutUseCase.GetCheckoutsByStatus(entity.CheckoutStatus(status), page, pageSize)
	} else {
		checkouts, err = h.checkoutUseCase.GetAllCheckouts(page, pageSize)
	}

	if err != nil {
		h.logger.Error("Failed to list checkouts: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create response
	response := dto.CreateCheckoutsListResponse(checkouts, len(checkouts), page, pageSize)

	// Return checkouts
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetAdminCheckout handles retrieving a checkout by ID for admin
func (h *CheckoutHandler) GetAdminCheckout(w http.ResponseWriter, r *http.Request) {
	// Get checkout ID from URL
	vars := mux.Vars(r)
	checkoutID, err := strconv.ParseUint(vars["checkoutId"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid checkout ID: %v", err)
		http.Error(w, "Invalid checkout ID", http.StatusBadRequest)
		return
	}

	// Get checkout
	checkout, err := h.checkoutUseCase.GetCheckoutByID(uint(checkoutID))
	if err != nil {
		h.logger.Error("Failed to get checkout: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.CreateCheckoutResponse(checkout)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteAdminCheckout handles deleting a checkout by ID (admin only)
func (h *CheckoutHandler) DeleteAdminCheckout(w http.ResponseWriter, r *http.Request) {
	// Get checkout ID from URL
	vars := mux.Vars(r)
	checkoutID, err := strconv.ParseUint(vars["checkoutId"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid checkout ID: %v", err)
		http.Error(w, "Invalid checkout ID", http.StatusBadRequest)
		return
	}

	// Delete checkout
	err = h.checkoutUseCase.DeleteCheckout(uint(checkoutID))
	if err != nil {
		h.logger.Error("Failed to delete checkout: %v", err)
		response := dto.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Return success response
	response := dto.SuccessResponseMessage("Checkout deleted successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
