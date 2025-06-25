package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/gkhaavik/vipps-mobilepay-sdk/pkg/models"
	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// WebhookHandler handles webhook requests from payment providers
type WebhookHandler struct {
	cfg            *config.Config
	orderUseCase   *usecase.OrderUseCase
	webhookUseCase *usecase.WebhookUseCase
	logger         logger.Logger
}

// NewWebhookHandler creates a new WebhookHandler
func NewWebhookHandler(
	cfg *config.Config,
	orderUseCase *usecase.OrderUseCase,
	webhookUseCase *usecase.WebhookUseCase,
	logger logger.Logger,
) *WebhookHandler {
	return &WebhookHandler{
		cfg:            cfg,
		orderUseCase:   orderUseCase,
		webhookUseCase: webhookUseCase,
		logger:         logger,
	}
}

// RegisterWebhookRequest represents a request to register a webhook
type RegisterWebhookRequest struct {
	Provider string   `json:"provider"`
	URL      string   `json:"url"`
	Events   []string `json:"events"`
}

// RegisterWebhook handles registering a new webhook
func (h *WebhookHandler) RegisterMobilePayWebhook(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var req RegisterWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to parse request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request
	if req.URL == "" || len(req.Events) == 0 {
		h.logger.Error("Invalid request: missing required fields")
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Register webhook
	input := usecase.RegisterWebhookInput{
		URL:    req.URL,
		Events: req.Events,
	}

	webhook, err := h.webhookUseCase.RegisterMobilePayWebhook(input)

	if err != nil {
		h.logger.Error("Failed to register webhook: %v", err)
		http.Error(w, "Failed to register webhook", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(webhook)
}

// GetMobilePayWebhooks handles getting all webhooks for MobilePay
func (h *WebhookHandler) GetMobilePayWebhooks(w http.ResponseWriter, r *http.Request) {
	webhooks, err := h.webhookUseCase.GetMobilePayWebhooks()
	if err != nil {
		h.logger.Error("Failed to get webhooks: %v", err)
		http.Error(w, "Failed to get webhooks", http.StatusInternalServerError)
		return
	}

	// Return webhooks
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhooks)
}

func (h *WebhookHandler) DeleteMobilePayWebhook(w http.ResponseWriter, r *http.Request) {
	// Get webhook ID from URL
	vars := mux.Vars(r)
	mpWebhookID, ok := vars["externalId"]
	if !ok {
		http.Error(w, "External ID is required", http.StatusBadRequest)
		return
	}

	// Delete webhook
	if err := h.webhookUseCase.DeleteMobilePayWebhook(mpWebhookID); err != nil {
		h.logger.Error("Failed to delete webhook: %v", err)
		http.Error(w, "Failed to delete webhook", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}

// ListWebhooks handles listing all webhooks
func (h *WebhookHandler) ListWebhooks(w http.ResponseWriter, r *http.Request) {
	webhooks, err := h.webhookUseCase.GetAllWebhooks()
	if err != nil {
		h.logger.Error("Failed to list webhooks: %v", err)
		http.Error(w, "Failed to list webhooks", http.StatusInternalServerError)
		return
	}

	// Return webhooks
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhooks)
}

// GetWebhook handles getting a webhook by ID
func (h *WebhookHandler) GetWebhook(w http.ResponseWriter, r *http.Request) {
	// Get webhook ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["webhookId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid webhook ID", http.StatusBadRequest)
		return
	}

	// Get webhook
	webhook, err := h.webhookUseCase.GetWebhookByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get webhook: %v", err)
		http.Error(w, "Webhook not found", http.StatusNotFound)
		return
	}

	// Return webhook
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhook)
}

// DeleteWebhook handles deleting a webhook
func (h *WebhookHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	// Get webhook ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["webhookId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid webhook ID", http.StatusBadRequest)
		return
	}

	// Delete webhook
	if err := h.webhookUseCase.DeleteWebhook(uint(id)); err != nil {
		h.logger.Error("Failed to delete webhook: %v", err)
		http.Error(w, "Failed to delete webhook", http.StatusInternalServerError)
		return
	}

	// Return success
	w.WriteHeader(http.StatusNoContent)
}

// HandleMobilePayAuthorized handles the AUTHORIZED event
func (h *WebhookHandler) HandleMobilePayAuthorized(event *models.WebhookEvent) error {
	orderID, err := extractOrderIDFromReference(event.Reference)
	if err != nil {
		h.logger.Error("Failed to extract order ID from reference: %v", err)
		return err
	}

	// Update payment status to authorized
	input := usecase.UpdatePaymentStatusInput{
		OrderID:       orderID,
		PaymentStatus: entity.PaymentStatusAuthorized,
	}

	order, err := h.orderUseCase.UpdatePaymentStatus(input)
	if err != nil {
		h.logger.Error("Failed to update payment status for MobilePay payment: %v", err)
		return err
	}

	err = h.recordMobilePayPaymentTransaction(order, entity.TransactionStatusSuccessful, event)
	if err != nil {
		h.logger.Error("Failed to record payment transaction: %v", err)
	}

	h.logger.Info("MobilePay payment authorized for order %d", orderID)
	return nil
}

func (h *WebhookHandler) recordMobilePayPaymentTransaction(order *entity.Order, transactionStatus entity.TransactionStatus, event *models.WebhookEvent) error {
	return h.orderUseCase.UpdatePaymentTransaction(order.PaymentID, transactionStatus, map[string]string{
		"pspReference":   event.Reference,
		"amount":         strconv.FormatInt(int64(event.Amount.Value), 10),
		"name":           string(event.Name),
		"idempotencyKey": event.IdempotencyKey,
	})
}

// HandleMobilePayCaptured handles the CAPTURED event
func (h *WebhookHandler) HandleMobilePayCaptured(event *models.WebhookEvent) error {
	orderID, err := extractOrderIDFromReference(event.Reference)
	if err != nil {
		h.logger.Error("Failed to extract order ID from reference: %v", err)
		return err
	}

	h.logger.Info("MobilePay payment captured for order %d", orderID)

	// Update payment status to captured
	input := usecase.UpdatePaymentStatusInput{
		OrderID:       orderID,
		PaymentStatus: entity.PaymentStatusCaptured,
	}

	order, err := h.orderUseCase.UpdatePaymentStatus(input)
	if err != nil {
		h.logger.Error("Failed to update payment status for MobilePay payment: %v", err)
		return err
	}

	err = h.recordMobilePayPaymentTransaction(order, entity.TransactionStatusSuccessful, event)
	if err != nil {
		h.logger.Error("Failed to record payment transaction: %v", err)
		return err
	}

	h.logger.Info("MobilePay payment captured for order %d", orderID)
	return nil
}

// HandleMobilePayCancelled handles the CANCELLED event
func (h *WebhookHandler) HandleMobilePayCancelled(event *models.WebhookEvent) error {
	orderID, err := extractOrderIDFromReference(event.Reference)
	if err != nil {
		h.logger.Error("Failed to extract order ID from reference: %v", err)
		return err
	}

	// Update order status to cancelled
	input := usecase.UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  entity.OrderStatusCancelled,
	}

	order, err2 := h.orderUseCase.UpdateOrderStatus(input)
	if err2 != nil {
		h.logger.Error("Failed to cancel order for MobilePay payment: %v", err2)
		return err2
	}

	err = h.recordMobilePayPaymentTransaction(order, entity.TransactionStatusFailed, event)
	if err != nil {
		h.logger.Error("Failed to record payment transaction: %v", err)
		return err
	}

	h.logger.Info("MobilePay payment cancelled for order %d", orderID)
	return nil
}

// HandleMobilePayRefunded handles the REFUNDED event
func (h *WebhookHandler) HandleMobilePayRefunded(event *models.WebhookEvent) error {
	orderID, err := extractOrderIDFromReference(event.Reference)
	if err != nil {
		h.logger.Error("Failed to extract order ID from reference: %v", err)
		return err
	}

	// Update payment status to refunded
	input := usecase.UpdatePaymentStatusInput{
		OrderID:       orderID,
		PaymentStatus: entity.PaymentStatusRefunded,
	}

	order, err2 := h.orderUseCase.UpdatePaymentStatus(input)
	if err2 != nil {
		h.logger.Error("Failed to update payment status to refunded for MobilePay payment: %v", err2)
		return err2
	}

	err = h.recordMobilePayPaymentTransaction(order, entity.TransactionStatusSuccessful, event)
	if err != nil {
		h.logger.Error("Failed to record payment transaction: %v", err)
		return err
	}

	h.logger.Info("MobilePay payment refunded for order %d", orderID)
	return nil
}

// HandleMobilePayAborted handles the ABORTED event
func (h *WebhookHandler) HandleMobilePayAborted(event *models.WebhookEvent) error {
	orderID, err := extractOrderIDFromReference(event.Reference)
	if err != nil {
		h.logger.Error("Failed to extract order ID from reference: %v", err)
		return err
	}

	// Update order status to cancelled
	input := usecase.UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  entity.OrderStatusCancelled,
	}

	order, err2 := h.orderUseCase.UpdateOrderStatus(input)
	if err2 != nil {
		h.logger.Error("Failed to cancel order for MobilePay aborted payment: %v", err2)
		return err2
	}

	err = h.recordMobilePayPaymentTransaction(order, entity.TransactionStatusFailed, event)
	if err != nil {
		h.logger.Error("Failed to record payment transaction: %v", err)
		return err
	}

	h.logger.Info("MobilePay payment aborted for order %d", orderID)
	return nil
}

// HandleMobilePayExpired handles the EXPIRED event
func (h *WebhookHandler) HandleMobilePayExpired(event *models.WebhookEvent) error {
	orderID, err := extractOrderIDFromReference(event.Reference)
	if err != nil {
		h.logger.Error("Failed to extract order ID from reference: %v", err)
		return err
	}

	// Update order status to cancelled
	input := usecase.UpdateOrderStatusInput{
		OrderID: orderID,
		Status:  entity.OrderStatusCancelled,
	}

	order, err2 := h.orderUseCase.UpdateOrderStatus(input)
	if err2 != nil {
		h.logger.Error("Failed to cancel order for MobilePay expired payment: %v", err2)
		return err2
	}

	err = h.recordMobilePayPaymentTransaction(order, entity.TransactionStatusFailed, event)
	if err != nil {
		h.logger.Error("Failed to record payment transaction: %v", err)
		return err
	}

	h.logger.Info("MobilePay payment expired for order %d", orderID)
	return nil
}

// extractOrderIDFromReference extracts the order ID from the reference
// Reference format: "order-{orderID}-{uuid}"
func extractOrderIDFromReference(reference string) (uint, error) {
	var orderID uint
	_, err := fmt.Sscanf(reference, "order-%d-", &orderID)
	if err != nil {
		return 0, fmt.Errorf("invalid reference format: %v", err)
	}
	return orderID, nil
}

// HandleStripeWebhook handles webhook events from Stripe
func (h *WebhookHandler) HandleStripeWebhook(w http.ResponseWriter, r *http.Request) {
	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Failed to read webhook body: %v", err)
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	// Verify the webhook signature
	webhookSecret := h.cfg.Stripe.WebhookSecret
	event, err := webhook.ConstructEvent(body, r.Header.Get("Stripe-Signature"), webhookSecret)
	if err != nil {
		h.logger.Error("Failed to verify webhook signature: %v", err)
		http.Error(w, "Failed to verify webhook signature", http.StatusBadRequest)
		return
	}

	// Handle different event types
	switch event.Type {
	case "payment_intent.succeeded":
		h.handlePaymentSucceeded(event)
	case "payment_intent.payment_failed":
		h.handlePaymentFailed(event)
	case "payment_intent.canceled":
		h.handlePaymentCanceled(event)
	case "payment_intent.requires_action":
		h.handlePaymentRequiresAction(event)
	case "payment_intent.processing":
		h.handlePaymentProcessing(event)
	case "payment_intent.amount_capturable_updated":
		h.handlePaymentCapturableUpdated(event)
	case "charge.succeeded":
		h.handleChargeSucceeded(event)
	case "charge.failed":
		h.handleChargeFailed(event)
	case "charge.refunded":
		h.handleRefund(event)
	case "charge.dispute.created":
		h.handleDisputeCreated(event)
	case "charge.dispute.closed":
		h.handleDisputeClosed(event)
	default:
		h.logger.Info("Received unhandled webhook event: %s", event.Type)
	}

	// Return a successful response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

// handlePaymentSucceeded handles the payment_intent.succeeded event
func (h *WebhookHandler) handlePaymentSucceeded(event stripe.Event) {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		h.logger.Error("Failed to parse payment intent: %v", err)
		return
	}

	// Get the order ID from metadata
	orderIDStr, ok := paymentIntent.Metadata["order_id"]
	if !ok {
		h.logger.Error("Order ID not found in payment intent metadata")
		return
	}

	// Convert order ID to uint
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid order ID in metadata: %v", err)
		return
	}

	// Record the successful payment transaction
	txn, err := entity.NewPaymentTransaction(
		uint(orderID),
		paymentIntent.ID,
		entity.TransactionTypeAuthorize,
		entity.TransactionStatusSuccessful,
		paymentIntent.Amount,
		string(paymentIntent.Currency),
		"stripe",
	)

	if err == nil {
		// Add raw response for debugging
		txn.SetRawResponse(string(event.Data.Raw))

		// Add metadata
		if method, exists := paymentIntent.Metadata["method"]; exists {
			txn.AddMetadata("payment_method", method)
		}

		// Record the transaction
		err = h.orderUseCase.RecordPaymentTransaction(txn)
		if err != nil {
			h.logger.Error("Failed to record payment transaction: %v", err)
		}
	}

	// Update the order status to paid
	input := usecase.UpdateOrderStatusInput{
		OrderID: uint(orderID),
		Status:  entity.OrderStatusPaid,
	}

	_, err = h.orderUseCase.UpdateOrderStatus(input)
	if err != nil {
		h.logger.Error("Failed to update order status: %v", err)
		return
	}

	h.logger.Info("Payment succeeded for order %d", orderID)
}

// handlePaymentFailed handles the payment_intent.payment_failed event
func (h *WebhookHandler) handlePaymentFailed(event stripe.Event) {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		h.logger.Error("Failed to parse payment intent: %v", err)
		return
	}

	// Get the order ID from metadata
	orderIDStr, ok := paymentIntent.Metadata["order_id"]
	if !ok {
		h.logger.Error("Order ID not found in payment intent metadata")
		return
	}

	// Convert order ID to uint
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid order ID in metadata: %v", err)
		return
	}

	// Record the failed payment transaction
	txn, err := entity.NewPaymentTransaction(
		uint(orderID),
		paymentIntent.ID,
		entity.TransactionTypeAuthorize,
		entity.TransactionStatusFailed,
		paymentIntent.Amount,
		string(paymentIntent.Currency),
		"stripe",
	)

	if err == nil {
		txn.SetRawResponse(string(event.Data.Raw))

		// Add metadata including error message
		if paymentIntent.LastPaymentError != nil {
			txn.AddMetadata("error_message", paymentIntent.LastPaymentError.Msg)
			txn.AddMetadata("error_code", string(paymentIntent.LastPaymentError.Code))
		}

		// Record the transaction
		err = h.orderUseCase.RecordPaymentTransaction(txn)
		if err != nil {
			h.logger.Error("Failed to record payment transaction: %v", err)
		}
	}

	// Update order status to payment_failed
	input := usecase.UpdateOrderStatusInput{
		OrderID: uint(orderID),
		Status:  entity.OrderStatusCancelled,
	}

	_, err = h.orderUseCase.UpdateOrderStatus(input)
	if err != nil {
		h.logger.Error("Failed to update order status: %v", err)
		return
	}

	// Log the payment failure
	errorMsg := "Unknown error"
	if paymentIntent.LastPaymentError != nil {
		errorMsg = paymentIntent.LastPaymentError.Msg
	}
	h.logger.Info("Payment failed for order %d: %s", orderID, errorMsg)
}

// handlePaymentCanceled handles the payment_intent.canceled event
func (h *WebhookHandler) handlePaymentCanceled(event stripe.Event) {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		h.logger.Error("Failed to parse payment intent: %v", err)
		return
	}

	// Get the order ID from metadata
	orderIDStr, ok := paymentIntent.Metadata["order_id"]
	if !ok {
		h.logger.Error("Order ID not found in payment intent metadata")
		return
	}

	// Convert order ID to uint
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid order ID in metadata: %v", err)
		return
	}

	// Record the cancel transaction
	txn, err := entity.NewPaymentTransaction(
		uint(orderID),
		paymentIntent.ID,
		entity.TransactionTypeCancel,
		entity.TransactionStatusSuccessful,
		0, // No amount for cancellation
		string(paymentIntent.Currency),
		"stripe",
	)

	if err == nil {
		txn.SetRawResponse(string(event.Data.Raw))

		// Record the transaction
		err = h.orderUseCase.RecordPaymentTransaction(txn)
		if err != nil {
			h.logger.Error("Failed to record payment transaction: %v", err)
		}
	}

	// Update order status to cancelled
	input := usecase.UpdateOrderStatusInput{
		OrderID: uint(orderID),
		Status:  entity.OrderStatusCancelled,
	}

	_, err = h.orderUseCase.UpdateOrderStatus(input)
	if err != nil {
		h.logger.Error("Failed to update order status: %v", err)
		return
	}

	h.logger.Info("Payment canceled for order %d", orderID)
}

// handlePaymentRequiresAction handles the payment_intent.requires_action event
func (h *WebhookHandler) handlePaymentRequiresAction(event stripe.Event) {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		h.logger.Error("Failed to parse payment intent: %v", err)
		return
	}

	// Get the order ID from metadata
	orderIDStr, ok := paymentIntent.Metadata["order_id"]
	if !ok {
		h.logger.Error("Order ID not found in payment intent metadata")
		return
	}

	// Convert order ID to uint
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid order ID in metadata: %v", err)
		return
	}

	h.logger.Info("Payment requires action for order %d", orderID)
}

// handlePaymentProcessing handles the payment_intent.processing event
func (h *WebhookHandler) handlePaymentProcessing(event stripe.Event) {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		h.logger.Error("Failed to parse payment intent: %v", err)
		return
	}

	// Get the order ID from metadata
	orderIDStr, ok := paymentIntent.Metadata["order_id"]
	if !ok {
		h.logger.Error("Order ID not found in payment intent metadata")
		return
	}

	// Convert order ID to uint
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid order ID in metadata: %v", err)
		return
	}

	h.logger.Info("Payment is processing for order %d", orderID)

	// Update order status to processing_payment if needed
	input := usecase.UpdateOrderStatusInput{
		OrderID: uint(orderID),
		Status:  entity.OrderStatusPending,
	}

	_, err = h.orderUseCase.UpdateOrderStatus(input)
	if err != nil {
		h.logger.Error("Failed to update order status: %v", err)
		return
	}
}

// handlePaymentCapturableUpdated handles the payment_intent.amount_capturable_updated event
func (h *WebhookHandler) handlePaymentCapturableUpdated(event stripe.Event) {
	var paymentIntent stripe.PaymentIntent
	err := json.Unmarshal(event.Data.Raw, &paymentIntent)
	if err != nil {
		h.logger.Error("Failed to parse payment intent: %v", err)
		return
	}

	// Get the order ID from metadata
	orderIDStr, ok := paymentIntent.Metadata["order_id"]
	if !ok {
		h.logger.Error("Order ID not found in payment intent metadata")
		return
	}

	// Convert order ID to uint
	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		h.logger.Error("Invalid order ID in metadata: %v", err)
		return
	}

	h.logger.Info("Payment is now capturable for order %d", orderID)
}

// handleChargeSucceeded handles the charge.succeeded event
func (h *WebhookHandler) handleChargeSucceeded(event stripe.Event) {
	var charge stripe.Charge
	err := json.Unmarshal(event.Data.Raw, &charge)
	if err != nil {
		h.logger.Error("Failed to parse charge: %v", err)
		return
	}

	// If there's no payment intent attached, we can't process further
	if charge.PaymentIntent == nil || charge.PaymentIntent.ID == "" {
		h.logger.Warn("Charge without payment intent ID received")
		return
	}

	h.logger.Info("Charge succeeded for payment intent %s", charge.PaymentIntent.ID)
}

// handleChargeFailed handles the charge.failed event
func (h *WebhookHandler) handleChargeFailed(event stripe.Event) {
	var charge stripe.Charge
	err := json.Unmarshal(event.Data.Raw, &charge)
	if err != nil {
		h.logger.Error("Failed to parse charge: %v", err)
		return
	}

	// If there's no payment intent attached, we can't process further
	if charge.PaymentIntent == nil || charge.PaymentIntent.ID == "" {
		h.logger.Warn("Charge without payment intent ID received")
		return
	}

	h.logger.Info("Charge failed for payment intent %s: %s",
		charge.PaymentIntent.ID,
		charge.FailureMessage)
}

// handleRefund handles the charge.refunded event
func (h *WebhookHandler) handleRefund(event stripe.Event) {
	var charge stripe.Charge
	err := json.Unmarshal(event.Data.Raw, &charge)
	if err != nil {
		h.logger.Error("Failed to parse charge: %v", err)
		return
	}

	// If there's no payment intent attached, we can't process further
	if charge.PaymentIntent == nil || charge.PaymentIntent.ID == "" {
		h.logger.Warn("Charge without payment intent ID received")
		return
	}

	// Find order by payment ID
	order, err := h.orderUseCase.GetOrderByPaymentID(charge.PaymentIntent.ID)
	if err != nil {
		h.logger.Error("Failed to find order for payment intent %s: %v", charge.PaymentIntent.ID, err)
		return
	}

	// Record the refund transaction
	txn, err := entity.NewPaymentTransaction(
		order.ID,
		charge.PaymentIntent.ID,
		entity.TransactionTypeRefund,
		entity.TransactionStatusSuccessful,
		charge.AmountRefunded,
		string(charge.Currency),
		"stripe",
	)

	if err == nil {
		txn.SetRawResponse(string(event.Data.Raw))

		// Record the transaction
		err = h.orderUseCase.RecordPaymentTransaction(txn)
		if err != nil {
			h.logger.Error("Failed to record refund transaction: %v", err)
		}
	}

	// If the charge was fully refunded, update the payment status
	if charge.Refunded {
		input := usecase.UpdatePaymentStatusInput{
			OrderID:       order.ID,
			PaymentStatus: entity.PaymentStatusRefunded,
		}

		_, err = h.orderUseCase.UpdatePaymentStatus(input)
		if err != nil {
			h.logger.Error("Failed to update payment status to refunded: %v", err)
			return
		}
	}

	h.logger.Info("Refund processed for order %d, payment %s, amount: %v",
		order.ID,
		charge.PaymentIntent.ID,
		money.FromCents(charge.AmountRefunded))
}

// handleDisputeCreated handles the charge.dispute.created event
func (h *WebhookHandler) handleDisputeCreated(event stripe.Event) {
	var dispute stripe.Dispute
	err := json.Unmarshal(event.Data.Raw, &dispute)
	if err != nil {
		h.logger.Error("Failed to parse dispute: %v", err)
		return
	}

	// If there's no payment intent attached, we can't process further
	if dispute.PaymentIntent == nil || dispute.PaymentIntent.ID == "" {
		h.logger.Warn("Dispute without payment intent ID received")
		return
	}

	h.logger.Warn("Dispute created for payment intent %s, reason: %s",
		dispute.PaymentIntent.ID,
		dispute.Reason)
}

// handleDisputeClosed handles the charge.dispute.closed event
func (h *WebhookHandler) handleDisputeClosed(event stripe.Event) {
	var dispute stripe.Dispute
	err := json.Unmarshal(event.Data.Raw, &dispute)
	if err != nil {
		h.logger.Error("Failed to parse dispute: %v", err)
		return
	}

	// If there's no payment intent attached, we can't process further
	if dispute.PaymentIntent == nil || dispute.PaymentIntent.ID == "" {
		h.logger.Warn("Dispute without payment intent ID received")
		return
	}

	h.logger.Info("Dispute closed for payment intent %s with status: %s",
		dispute.PaymentIntent.ID,
		dispute.Status)
}
