package handler

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
)

// StripeWebhookHandler handles Stripe webhook callbacks
type StripeWebhookHandler struct {
	orderUseCase *usecase.OrderUseCase
	config       *config.Config
	logger       logger.Logger
}

// NewStripeWebhookHandler creates a new StripeWebhookHandler
func NewStripeWebhookHandler(orderUseCase *usecase.OrderUseCase, cfg *config.Config, logger logger.Logger) *StripeWebhookHandler {
	return &StripeWebhookHandler{
		orderUseCase: orderUseCase,
		config:       cfg,
		logger:       logger,
	}
}

// HandleWebhook handles incoming Stripe webhook events
func (h *StripeWebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

	payload, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Failed to read Stripe webhook body: %v", err)
		http.Error(w, "Error reading request body", http.StatusServiceUnavailable)
		return
	}

	// Verify webhook signature
	if !h.verifySignature(payload, r.Header.Get("Stripe-Signature")) {
		h.logger.Error("Invalid Stripe webhook signature")
		http.Error(w, "Invalid signature", http.StatusUnauthorized)
		return
	}

	// Parse the webhook event
	var event StripeWebhookEvent
	if err := json.Unmarshal(payload, &event); err != nil {
		h.logger.Error("Failed to parse Stripe webhook event: %v", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	h.logger.Info("Received Stripe webhook event: %s", event.Type)

	// Process the event
	if err := h.processEvent(&event); err != nil {
		h.logger.Error("Failed to process Stripe webhook event: %v", err)
		http.Error(w, "Error processing event", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// StripeWebhookEvent represents a Stripe webhook event
type StripeWebhookEvent struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Created int64  `json:"created"`
	Data    struct {
		Object any `json:"object"`
	} `json:"data"`
	Request struct {
		ID             string `json:"id"`
		IdempotencyKey string `json:"idempotency_key"`
	} `json:"request"`
}

// recordPaymentTransaction creates and saves a payment transaction record for Stripe events
func (h *StripeWebhookHandler) recordPaymentTransaction(orderID uint, transactionID string, txnType entity.TransactionType, status entity.TransactionStatus, amount int64, currency, provider string, event *StripeWebhookEvent) error {
	// Create payment transaction
	txn, err := entity.NewPaymentTransaction(
		orderID,
		transactionID,
		"", // No idempotency key for Stripe events currently
		txnType,
		status,
		amount,
		currency,
		provider,
	)
	if err != nil {
		return fmt.Errorf("failed to create payment transaction: %w", err)
	}

	// Add webhook event data as raw response
	if event != nil {
		// Convert the entire event to JSON string for storage
		if eventJSON, err := json.Marshal(event); err == nil {
			txn.SetRawResponse(string(eventJSON))
		}

		// Add metadata
		txn.AddMetadata("webhook_event_type", event.Type)
		txn.AddMetadata("webhook_event_id", event.ID)
		if event.Created > 0 {
			txn.AddMetadata("webhook_created", strconv.FormatInt(event.Created, 10))
		}
		// Add request metadata if available
		if event.Request.ID != "" {
			txn.AddMetadata("webhook_request_id", event.Request.ID)
		}
		if event.Request.IdempotencyKey != "" {
			txn.AddMetadata("idempotency_key", event.Request.IdempotencyKey)
		}
	}

	// Save the transaction using the usecase
	return h.orderUseCase.RecordPaymentTransaction(txn)
}

// verifySignature verifies the Stripe webhook signature
func (h *StripeWebhookHandler) verifySignature(payload []byte, signature string) bool {
	if h.config.Stripe.WebhookSecret == "" {
		h.logger.Warn("Stripe webhook secret not configured, skipping signature verification")
		return true // In development, allow unsigned webhooks
	}

	// Parse the signature header
	signatureParts := strings.Split(signature, ",")
	var timestamp, signature256 string

	for _, part := range signatureParts {
		if after, ok := strings.CutPrefix(part, "t="); ok {
			timestamp = after
		} else if after0, ok0 := strings.CutPrefix(part, "v1="); ok0 {
			signature256 = after0
		}
	}

	if timestamp == "" || signature256 == "" {
		return false
	}

	// Compute the expected signature
	expectedPayload := timestamp + "." + string(payload)
	mac := hmac.New(sha256.New, []byte(h.config.Stripe.WebhookSecret))
	mac.Write([]byte(expectedPayload))
	expectedSignature := hex.EncodeToString(mac.Sum(nil))

	return hmac.Equal([]byte(signature256), []byte(expectedSignature))
}

// processEvent processes a Stripe webhook event
func (h *StripeWebhookHandler) processEvent(event *StripeWebhookEvent) error {
	switch event.Type {
	case "payment_intent.succeeded":
		return h.handlePaymentSucceeded(event)
	case "payment_intent.payment_failed":
		return h.handlePaymentFailed(event)
	case "payment_intent.canceled":
		return h.handlePaymentCanceled(event)
	case "payment_intent.requires_action":
		return h.handlePaymentRequiresAction(event)
	case "payment_intent.amount_capturable_updated":
		return h.handleAmountCapturableUpdated(event)
	case "payment_intent.partially_funded":
		return h.handlePartiallyFunded(event)
	case "charge.captured":
		return h.handleChargeCaptured(event)
	case "charge.dispute.created":
		return h.handleChargeDispute(event)
	case "invoice.payment_succeeded":
		return h.handleInvoicePaymentSucceeded(event)
	case "invoice.payment_failed":
		return h.handleInvoicePaymentFailed(event)
	default:
		h.logger.Info("Unhandled Stripe webhook event type: %s", event.Type)
		return nil
	}
}

// handlePaymentSucceeded handles successful Stripe payments
func (h *StripeWebhookHandler) handlePaymentSucceeded(event *StripeWebhookEvent) error {
	paymentIntent, ok := event.Data.Object.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid payment intent data")
	}

	transactionID, _ := paymentIntent["id"].(string)
	orderIDStr, _ := paymentIntent["metadata"].(map[string]any)["order_id"].(string)

	if orderIDStr == "" {
		h.logger.Warn("No order_id in Stripe payment intent metadata")
		return nil
	}

	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid order ID: %w", err)
	}

	h.logger.Info("Processing successful Stripe payment for order %d, transaction %s", orderID, transactionID)

	// Get order to access payment details for transaction recording
	order, err := h.orderUseCase.GetOrderByID(uint(orderID))
	if err != nil {
		h.logger.Error("Failed to get order %d for payment transaction recording: %v", orderID, err)
		return err
	}

	// Get the amount from the payment intent
	amount := int64(0)
	if amountFloat, ok := paymentIntent["amount"].(float64); ok {
		amount = int64(amountFloat)
	}

	// Get the currency from the payment intent
	currency := order.Currency
	if currencyStr, ok := paymentIntent["currency"].(string); ok && currencyStr != "" {
		currency = currencyStr
	}

	// Record the successful authorization transaction
	if recordErr := h.recordPaymentTransaction(uint(orderID), transactionID, entity.TransactionTypeAuthorize, entity.TransactionStatusSuccessful, amount, currency, "stripe", event); recordErr != nil {
		h.logger.Error("Failed to record authorization transaction for order %d: %v", orderID, recordErr)
		// Don't fail the webhook processing if transaction recording fails
	}

	// Update order payment status to authorized
	_, err = h.orderUseCase.UpdatePaymentStatus(usecase.UpdatePaymentStatusInput{
		OrderID:       uint(orderID),
		PaymentStatus: entity.PaymentStatusAuthorized,
		TransactionID: transactionID,
	})

	return err
}

// handlePaymentFailed handles failed Stripe payments
func (h *StripeWebhookHandler) handlePaymentFailed(event *StripeWebhookEvent) error {
	paymentIntent, ok := event.Data.Object.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid payment intent data")
	}

	transactionID, _ := paymentIntent["id"].(string)
	orderIDStr, _ := paymentIntent["metadata"].(map[string]any)["order_id"].(string)

	if orderIDStr == "" {
		h.logger.Warn("No order_id in Stripe payment intent metadata")
		return nil
	}

	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid order ID: %w", err)
	}

	h.logger.Info("Processing failed Stripe payment for order %d, transaction %s", orderID, transactionID)

	// Get order to access payment details for transaction recording
	order, err := h.orderUseCase.GetOrderByID(uint(orderID))
	if err != nil {
		h.logger.Error("Failed to get order %d for payment transaction recording: %v", orderID, err)
		return err
	}

	// Record the failed authorization transaction
	if recordErr := h.recordPaymentTransaction(uint(orderID), transactionID, entity.TransactionTypeAuthorize, entity.TransactionStatusFailed, 0, order.Currency, "stripe", event); recordErr != nil {
		h.logger.Error("Failed to record failed transaction for order %d: %v", orderID, recordErr)
		// Don't fail the webhook processing if transaction recording fails
	}

	// Update order payment status to failed
	_, err = h.orderUseCase.UpdatePaymentStatus(usecase.UpdatePaymentStatusInput{
		OrderID:       uint(orderID),
		PaymentStatus: entity.PaymentStatusFailed,
		TransactionID: transactionID,
	})

	return err

}

// handlePaymentCanceled handles canceled Stripe payments
func (h *StripeWebhookHandler) handlePaymentCanceled(event *StripeWebhookEvent) error {
	paymentIntent, ok := event.Data.Object.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid payment intent data")
	}

	transactionID, _ := paymentIntent["id"].(string)
	orderIDStr, _ := paymentIntent["metadata"].(map[string]any)["order_id"].(string)

	if orderIDStr == "" {
		h.logger.Warn("No order_id in Stripe payment intent metadata")
		return nil
	}

	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid order ID: %w", err)
	}

	h.logger.Info("Processing canceled Stripe payment for order %d, transaction %s", orderID, transactionID)

	// Get order to access payment details for transaction recording
	order, err := h.orderUseCase.GetOrderByID(uint(orderID))
	if err != nil {
		h.logger.Error("Failed to get order %d for payment transaction recording: %v", orderID, err)
		return err
	}

	// Record the cancellation transaction
	if recordErr := h.recordPaymentTransaction(uint(orderID), transactionID, entity.TransactionTypeCancel, entity.TransactionStatusSuccessful, 0, order.Currency, "stripe", event); recordErr != nil {
		h.logger.Error("Failed to record cancellation transaction for order %d: %v", orderID, recordErr)
		// Don't fail the webhook processing if transaction recording fails
	}

	// Update order payment status to cancelled
	_, err = h.orderUseCase.UpdatePaymentStatus(usecase.UpdatePaymentStatusInput{
		OrderID:       uint(orderID),
		PaymentStatus: entity.PaymentStatusCancelled,
		TransactionID: transactionID,
	})

	return err

}

// handlePaymentRequiresAction handles Stripe payments that require action
func (h *StripeWebhookHandler) handlePaymentRequiresAction(event *StripeWebhookEvent) error {
	paymentIntent, ok := event.Data.Object.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid payment intent data")
	}

	transactionID, _ := paymentIntent["id"].(string)
	orderIDStr, _ := paymentIntent["metadata"].(map[string]any)["order_id"].(string)

	if orderIDStr == "" {
		h.logger.Warn("No order_id in Stripe payment intent metadata")
		return nil
	}

	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid order ID: %w", err)
	}

	h.logger.Info("Processing Stripe payment requiring action for order %d, transaction %s", orderID, transactionID)

	// Update order payment status to pending (awaiting action)
	_, err = h.orderUseCase.UpdatePaymentStatus(usecase.UpdatePaymentStatusInput{
		OrderID:       uint(orderID),
		PaymentStatus: entity.PaymentStatusPending,
		TransactionID: transactionID,
	})
	return err

}

// handleAmountCapturableUpdated handles when the capturable amount is updated
func (h *StripeWebhookHandler) handleAmountCapturableUpdated(event *StripeWebhookEvent) error {
	paymentIntent, ok := event.Data.Object.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid payment intent data")
	}

	transactionID, _ := paymentIntent["id"].(string)
	orderIDStr, _ := paymentIntent["metadata"].(map[string]any)["order_id"].(string)

	if orderIDStr == "" {
		h.logger.Warn("No order_id in Stripe payment intent metadata")
		return nil
	}

	orderID, err := strconv.ParseUint(orderIDStr, 10, 32)
	if err != nil {
		return fmt.Errorf("invalid order ID: %w", err)
	}

	h.logger.Info("Processing Stripe capturable amount update for order %d, transaction %s", orderID, transactionID)

	_, err = h.orderUseCase.UpdatePaymentStatus(usecase.UpdatePaymentStatusInput{
		OrderID:       uint(orderID),
		PaymentStatus: entity.PaymentStatusCaptured,
		TransactionID: transactionID,
	})

	return err
}

// handlePartiallyFunded handles partially funded payments
func (h *StripeWebhookHandler) handlePartiallyFunded(event *StripeWebhookEvent) error {
	paymentIntent, ok := event.Data.Object.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid payment intent data")
	}

	transactionID, _ := paymentIntent["id"].(string)
	h.logger.Info("Processing partially funded Stripe payment: %s", transactionID)

	// For partially funded payments, we might want to handle them differently
	// For now, just log the event
	return nil
}

// handleChargeCaptured handles charge captured events from Stripe
func (h *StripeWebhookHandler) handleChargeCaptured(event *StripeWebhookEvent) error {
	charge, ok := event.Data.Object.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid charge data")
	}

	transactionID, _ := charge["id"].(string)

	// Get payment intent from charge to access metadata
	paymentIntentID, _ := charge["payment_intent"].(string)
	if paymentIntentID == "" {
		h.logger.Warn("No payment_intent in Stripe charge")
		return nil
	}

	// For now, we'll use the payment_intent ID to find the order
	// In a full implementation, you might need to make an API call to Stripe to get the payment intent metadata
	// For this implementation, we'll extract order_id from description or other fields if available
	description, _ := charge["description"].(string)

	// Extract order ID from description if it follows a pattern like "Order #123"
	orderID := extractOrderIDFromDescription(description)
	if orderID == 0 {
		h.logger.Warn("Could not extract order ID from charge description: %s", description)
		return nil
	}

	h.logger.Info("Processing captured Stripe charge for order %d, transaction %s", orderID, transactionID)

	// Get order to access payment details for transaction recording
	order, err := h.orderUseCase.GetOrderByID(orderID)
	if err != nil {
		h.logger.Error("Failed to get order %d for charge capture transaction recording: %v", orderID, err)
		return err
	}

	// Get the amount from the charge
	amount := int64(0)
	if amountFloat, ok := charge["amount"].(float64); ok {
		amount = int64(amountFloat)
	}

	// Get the currency from the charge
	currency := order.Currency
	if currencyStr, ok := charge["currency"].(string); ok && currencyStr != "" {
		currency = currencyStr
	}

	// Record the capture transaction
	if recordErr := h.recordPaymentTransaction(orderID, transactionID, entity.TransactionTypeCapture, entity.TransactionStatusSuccessful, amount, currency, "stripe", event); recordErr != nil {
		h.logger.Error("Failed to record capture transaction for order %d: %v", orderID, recordErr)
		// Don't fail the webhook processing if transaction recording fails
	}

	// Update order payment status to captured
	_, err = h.orderUseCase.UpdatePaymentStatus(usecase.UpdatePaymentStatusInput{
		OrderID:       orderID,
		PaymentStatus: entity.PaymentStatusCaptured,
		TransactionID: transactionID,
	})

	return err
}

// handleChargeDispute handles charge dispute events from Stripe
func (h *StripeWebhookHandler) handleChargeDispute(event *StripeWebhookEvent) error {
	dispute, ok := event.Data.Object.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid dispute data")
	}

	// Get charge information from dispute
	chargeID, _ := dispute["charge"].(string)
	if chargeID == "" {
		h.logger.Warn("No charge ID in Stripe dispute")
		return nil
	}

	h.logger.Info("Processing Stripe charge dispute for charge %s", chargeID)

	// For now, just log the dispute - in a full implementation you might want to:
	// 1. Find the order associated with this charge
	// 2. Update order status to indicate dispute
	// 3. Send notifications to admin
	// 4. Record a transaction for the dispute

	h.logger.Warn("Charge dispute received for charge %s - manual review required", chargeID)

	return nil
}

// handleInvoicePaymentSucceeded handles successful invoice payments from Stripe
func (h *StripeWebhookHandler) handleInvoicePaymentSucceeded(event *StripeWebhookEvent) error {
	invoice, ok := event.Data.Object.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid invoice data")
	}

	invoiceID, _ := invoice["id"].(string)
	h.logger.Info("Processing successful Stripe invoice payment for invoice %s", invoiceID)

	// For subscription-based orders, you might handle these differently
	// For now, just log the successful payment
	h.logger.Info("Invoice payment succeeded for invoice %s", invoiceID)

	return nil
}

// handleInvoicePaymentFailed handles failed invoice payments from Stripe
func (h *StripeWebhookHandler) handleInvoicePaymentFailed(event *StripeWebhookEvent) error {
	invoice, ok := event.Data.Object.(map[string]any)
	if !ok {
		return fmt.Errorf("invalid invoice data")
	}

	invoiceID, _ := invoice["id"].(string)
	h.logger.Info("Processing failed Stripe invoice payment for invoice %s", invoiceID)

	// For subscription-based orders, you might handle these differently
	// For now, just log the failed payment
	h.logger.Warn("Invoice payment failed for invoice %s", invoiceID)

	return nil
}

// extractOrderIDFromDescription extracts order ID from description string
// This is a helper function that tries to parse order ID from charge description
func extractOrderIDFromDescription(description string) uint {
	// This is a simple implementation that looks for "Order #123" pattern
	// You might need to adjust this based on your actual description format
	if description == "" {
		return 0
	}

	// Try to extract order ID from description
	// Implementation depends on your description format
	// For now, return 0 to indicate that order ID extraction is not implemented
	// You would implement pattern matching here based on your description format
	return 0
}
