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

	// Update order payment status to captured
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
