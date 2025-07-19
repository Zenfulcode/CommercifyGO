package handler

import (
	"fmt"
	"net/http"

	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/vipps-mobilepay-sdk/pkg/models"
	"github.com/zenfulcode/vipps-mobilepay-sdk/pkg/webhooks"
)

// MobilePayWebhookHandler handles MobilePay webhook callbacks
type MobilePayWebhookHandler struct {
	orderUseCase           *usecase.OrderUseCase
	paymentProviderService service.PaymentProviderService
	config                 *config.Config
	logger                 logger.Logger
	webhookHandler         *webhooks.Handler
	webhookRouter          *webhooks.Router
}

// NewMobilePayWebhookHandler creates a new MobilePayWebhookHandler
func NewMobilePayWebhookHandler(
	orderUseCase *usecase.OrderUseCase,
	paymentProviderService service.PaymentProviderService,
	cfg *config.Config,
	logger logger.Logger,
) *MobilePayWebhookHandler {
	// Get webhook secret from the MobilePay payment provider in the database
	secretKey := getWebhookSecretFromDatabase(paymentProviderService, logger)

	// Create webhook handler with secret key
	webhookHandler := webhooks.NewHandler(secretKey)

	// Create webhook router
	webhookRouter := webhooks.NewRouter()

	handler := &MobilePayWebhookHandler{
		orderUseCase:           orderUseCase,
		paymentProviderService: paymentProviderService,
		config:                 cfg,
		logger:                 logger,
		webhookHandler:         webhookHandler,
		webhookRouter:          webhookRouter,
	}

	// Register event handlers
	handler.setupEventHandlers()

	return handler
}

// HandleWebhook handles incoming MobilePay webhook events using the official SDK
func (h *MobilePayWebhookHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	// Use the SDK's HTTP handler with our router
	h.webhookHandler.HandleHTTP(h.webhookRouter.Process)(w, r)
}

// setupEventHandlers registers event handlers for different webhook event types
func (h *MobilePayWebhookHandler) setupEventHandlers() {
	// Register handlers for different event types using the SDK's event constants
	h.webhookRouter.HandleFunc(models.EventAuthorized, h.handleSDKPaymentAuthorized)
	h.webhookRouter.HandleFunc(models.EventCaptured, h.handleSDKPaymentCaptured)
	h.webhookRouter.HandleFunc(models.EventCancelled, h.handleSDKPaymentCancelled)
	h.webhookRouter.HandleFunc(models.EventExpired, h.handleSDKPaymentExpired)
	h.webhookRouter.HandleFunc(models.EventRefunded, h.handleSDKPaymentRefunded)
}

// handleSDKPaymentAuthorized handles payment authorized events from the SDK
func (h *MobilePayWebhookHandler) handleSDKPaymentAuthorized(event *models.WebhookEvent) error {
	orderID, err := h.getOrderIDFromSDKEvent(event)
	if err != nil {
		return err
	}

	h.logger.Info("Processing authorized MobilePay payment for order %d, transaction %s", orderID, event.Reference)

	// Get the order to access payment details
	order, err := h.orderUseCase.GetOrderByID(orderID)
	if err != nil {
		h.logger.Error("Failed to get order %d for payment transaction recording: %v", orderID, err)
		return err
	}

	// Check if payment is already authorized or in a later stage (idempotency check)
	if order.PaymentStatus == entity.PaymentStatusAuthorized ||
		order.PaymentStatus == entity.PaymentStatusCaptured ||
		order.PaymentStatus == entity.PaymentStatusRefunded {
		h.logger.Info("Payment for order %d is already authorized or beyond, skipping duplicate authorization webhook", orderID)
		return nil
	}

	// Check if we already processed this exact webhook event using idempotency key (prevents duplicate webhook processing)
	if event.IdempotencyKey != "" {
		existingTxn, err := h.orderUseCase.GetTransactionByIdempotencyKey(event.IdempotencyKey)
		if err == nil && existingTxn != nil {
			h.logger.Info("Transaction with idempotency key %s already exists for order %d, skipping duplicate authorization webhook", event.IdempotencyKey, orderID)
			return nil
		}
	}

	// Record/update the authorization transaction
	if err := h.recordPaymentTransaction(orderID, event.Reference, entity.TransactionTypeAuthorize, entity.TransactionStatusSuccessful, order.FinalAmount, order.Currency, "mobilepay", event); err != nil {
		h.logger.Error("Failed to record authorization transaction for order %d: %v", orderID, err)
		// Don't fail the webhook processing if transaction recording fails
	}

	_, err = h.orderUseCase.UpdatePaymentStatus(usecase.UpdatePaymentStatusInput{
		OrderID:       orderID,
		PaymentStatus: entity.PaymentStatusAuthorized,
		TransactionID: event.Reference,
	})

	return err
}

// handleSDKPaymentCaptured handles payment captured events from the SDK
func (h *MobilePayWebhookHandler) handleSDKPaymentCaptured(event *models.WebhookEvent) error {
	orderID, err := h.getOrderIDFromSDKEvent(event)
	if err != nil {
		return err
	}

	h.logger.Info("Processing captured MobilePay payment for order %d, transaction %s", orderID, event.Reference)

	// Get the order to access payment details
	order, err := h.orderUseCase.GetOrderByID(orderID)
	if err != nil {
		h.logger.Error("Failed to get order %d for payment transaction recording: %v", orderID, err)
		return err
	}

	// Check if payment is already captured or refunded (idempotency check)
	if order.PaymentStatus == entity.PaymentStatusCaptured ||
		order.PaymentStatus == entity.PaymentStatusRefunded {
		h.logger.Info("Payment for order %d is already captured or refunded, skipping duplicate capture webhook", orderID)
		return nil
	}

	// Check if we already processed this exact webhook event using idempotency key (prevents duplicate webhook processing)
	if event.IdempotencyKey != "" {
		existingTxn, err := h.orderUseCase.GetTransactionByIdempotencyKey(event.IdempotencyKey)
		if err == nil && existingTxn != nil {
			h.logger.Info("Transaction with idempotency key %s already exists for order %d, skipping duplicate capture webhook", event.IdempotencyKey, orderID)
			return nil
		}
	}

	// Record/update the capture transaction
	// Use the amount from the webhook event if available, otherwise use order amount
	captureAmount := order.FinalAmount
	if event.Amount.Value > 0 {
		captureAmount = int64(event.Amount.Value)
	}

	if err := h.recordPaymentTransaction(orderID, event.Reference, entity.TransactionTypeCapture, entity.TransactionStatusSuccessful, captureAmount, order.Currency, "mobilepay", event); err != nil {
		h.logger.Error("Failed to record capture transaction for order %d: %v", orderID, err)
		// Don't fail the webhook processing if transaction recording fails
	}

	// Update order payment status to captured
	_, err = h.orderUseCase.UpdatePaymentStatus(usecase.UpdatePaymentStatusInput{
		OrderID:       orderID,
		PaymentStatus: entity.PaymentStatusCaptured,
		TransactionID: event.Reference,
	})

	return err
}

// handleSDKPaymentCancelled handles payment cancelled events from the SDK
func (h *MobilePayWebhookHandler) handleSDKPaymentCancelled(event *models.WebhookEvent) error {
	orderID, err := h.getOrderIDFromSDKEvent(event)
	if err != nil {
		return err
	}

	h.logger.Info("Processing cancelled MobilePay payment for order %d, transaction %s", orderID, event.Reference)

	// Get the order to access payment details
	order, err := h.orderUseCase.GetOrderByID(orderID)
	if err != nil {
		h.logger.Error("Failed to get order %d for payment transaction recording: %v", orderID, err)
		return err
	}

	// Check if payment is already cancelled (idempotency check)
	if order.PaymentStatus == entity.PaymentStatusCancelled {
		h.logger.Info("Payment for order %d is already cancelled, skipping duplicate cancellation webhook", orderID)
		return nil
	}

	// Check if we already processed this exact webhook event using idempotency key (prevents duplicate webhook processing)
	if event.IdempotencyKey != "" {
		existingTxn, err := h.orderUseCase.GetTransactionByIdempotencyKey(event.IdempotencyKey)
		if err == nil && existingTxn != nil {
			h.logger.Info("Transaction with idempotency key %s already exists for order %d, skipping duplicate cancellation webhook", event.IdempotencyKey, orderID)
			return nil
		}
	}

	// Record/update the cancellation transaction
	if err := h.recordPaymentTransaction(orderID, event.Reference, entity.TransactionTypeCancel, entity.TransactionStatusSuccessful, 0, order.Currency, "mobilepay", event); err != nil {
		h.logger.Error("Failed to record cancellation transaction for order %d: %v", orderID, err)
		// Don't fail the webhook processing if transaction recording fails
	}

	// Update order payment status to cancelled
	_, err = h.orderUseCase.UpdatePaymentStatus(usecase.UpdatePaymentStatusInput{
		OrderID:       orderID,
		PaymentStatus: entity.PaymentStatusCancelled,
		TransactionID: event.Reference,
	})

	return err
}

// handleSDKPaymentExpired handles payment expired events from the SDK
func (h *MobilePayWebhookHandler) handleSDKPaymentExpired(event *models.WebhookEvent) error {
	orderID, err := h.getOrderIDFromSDKEvent(event)
	if err != nil {
		return err
	}

	h.logger.Info("Processing expired MobilePay payment for order %d, transaction %s", orderID, event.Reference)

	// Get the order to access payment details
	order, err := h.orderUseCase.GetOrderByID(orderID)
	if err != nil {
		h.logger.Error("Failed to get order %d for payment transaction recording: %v", orderID, err)
		return err
	}

	// Check if payment is already failed, cancelled, or in a successful state (idempotency check)
	if order.PaymentStatus == entity.PaymentStatusFailed ||
		order.PaymentStatus == entity.PaymentStatusCancelled ||
		order.PaymentStatus == entity.PaymentStatusAuthorized ||
		order.PaymentStatus == entity.PaymentStatusCaptured ||
		order.PaymentStatus == entity.PaymentStatusRefunded {
		h.logger.Info("Payment for order %d is already in a final state (%s), skipping duplicate expiration webhook", orderID, order.PaymentStatus)
		return nil
	}

	// Check if we already processed this exact webhook event using idempotency key (prevents duplicate webhook processing)
	if event.IdempotencyKey != "" {
		existingTxn, err := h.orderUseCase.GetTransactionByIdempotencyKey(event.IdempotencyKey)
		if err == nil && existingTxn != nil {
			h.logger.Info("Transaction with idempotency key %s already exists for order %d, skipping duplicate expiration webhook", event.IdempotencyKey, orderID)
			return nil
		}
	}

	// Record/update the expiration as a failed transaction
	if err := h.recordPaymentTransaction(orderID, event.Reference, entity.TransactionTypeAuthorize, entity.TransactionStatusFailed, 0, order.Currency, "mobilepay", event); err != nil {
		h.logger.Error("Failed to record expiration transaction for order %d: %v", orderID, err)
		// Don't fail the webhook processing if transaction recording fails
	}

	// Update order payment status to failed
	_, err = h.orderUseCase.UpdatePaymentStatus(usecase.UpdatePaymentStatusInput{
		OrderID:       orderID,
		PaymentStatus: entity.PaymentStatusFailed,
		TransactionID: event.Reference,
	})

	return err
}

// handleSDKPaymentRefunded handles payment refunded events from the SDK
func (h *MobilePayWebhookHandler) handleSDKPaymentRefunded(event *models.WebhookEvent) error {
	orderID, err := h.getOrderIDFromSDKEvent(event)
	if err != nil {
		return err
	}

	h.logger.Info("Processing refunded MobilePay payment for order %d, transaction %s", orderID, event.Reference)

	// Get the order to access payment details
	order, err := h.orderUseCase.GetOrderByID(orderID)
	if err != nil {
		h.logger.Error("Failed to get order %d for payment transaction recording: %v", orderID, err)
		return err
	}

	// Check if payment is in a refundable state (idempotency check)
	// Allow refunds for both captured and already partially refunded payments
	if order.PaymentStatus != entity.PaymentStatusCaptured && order.PaymentStatus != entity.PaymentStatusRefunded {
		h.logger.Info("Payment for order %d is not in a refundable state (%s), skipping refund webhook", orderID, order.PaymentStatus)
		return nil
	}

	// Check if we already processed this exact webhook event using idempotency key (prevents duplicate webhook processing)
	// For refunds, we check to prevent the exact same webhook event from being processed multiple times
	if event.IdempotencyKey != "" {
		existingTxn, err := h.orderUseCase.GetTransactionByIdempotencyKey(event.IdempotencyKey)
		if err == nil && existingTxn != nil {
			h.logger.Info("Transaction with idempotency key %s already exists for order %d, skipping duplicate refund webhook", event.IdempotencyKey, orderID)
			return nil
		}
	}

	// For refunds, always create a new transaction (don't update pending ones)
	// This allows multiple partial refunds to be tracked separately
	refundAmount := order.FinalAmount
	if event.Amount.Value > 0 {
		refundAmount = int64(event.Amount.Value)
	}

	h.logger.Info("Creating new refund transaction for order %d with amount %d", orderID, refundAmount)
	if err := h.createNewTransaction(orderID, event.Reference, entity.TransactionTypeRefund, entity.TransactionStatusSuccessful, refundAmount, order.Currency, "mobilepay", event); err != nil {
		h.logger.Error("Failed to record refund transaction for order %d: %v", orderID, err)
		// Don't fail the webhook processing if transaction recording fails
	}

	// Always mark order as refunded when any refund occurs
	// The system can track partial vs full refunds through transaction records
	// Business logic elsewhere can determine if it's a full or partial refund by comparing totals
	_, err = h.orderUseCase.UpdatePaymentStatus(usecase.UpdatePaymentStatusInput{
		OrderID:       orderID,
		PaymentStatus: entity.PaymentStatusRefunded,
		TransactionID: event.Reference,
	})

	return err
}

// getOrderIDFromSDKEvent gets the order ID associated with a MobilePay payment from SDK event
func (h *MobilePayWebhookHandler) getOrderIDFromSDKEvent(event *models.WebhookEvent) (uint, error) {
	// Try to find the order by PaymentID field using the event reference
	order, err := h.orderUseCase.GetOrderByExternalID(event.Reference)
	if err != nil {
		h.logger.Error("Could not find order for MobilePay payment %s", event.Reference)
		return 0, fmt.Errorf("order not found for MobilePay payment %s", event.Reference)
	}

	return order.ID, nil
}

// getWebhookSecretFromDatabase retrieves the webhook secret for MobilePay from the database
func getWebhookSecretFromDatabase(paymentProviderService service.PaymentProviderService, logger logger.Logger) string {
	// Get the MobilePay payment provider from the database
	provider, err := paymentProviderService.GetWebhookInfo("mobilepay")
	if err != nil {
		logger.Error("Failed to get MobilePay payment provider from database: %v", err)
		return ""
	}

	if provider == nil {
		logger.Error("MobilePay payment provider not found in database")
		return ""
	}

	// Get webhook secret from provider
	if provider.WebhookSecret != "" {
		logger.Info("Retrieved MobilePay webhook secret from database")
		return provider.WebhookSecret
	}

	logger.Warn("MobilePay webhook secret not found in provider configuration")
	return ""
}

// recordPaymentTransaction creates and saves a payment transaction record
func (h *MobilePayWebhookHandler) recordPaymentTransaction(orderID uint, transactionID string, txnType entity.TransactionType, status entity.TransactionStatus, amount int64, currency, provider string, event *models.WebhookEvent) error {
	// Try to update existing pending transaction first
	if err := h.updateOrCreateTransaction(orderID, transactionID, txnType, status, amount, currency, provider, event); err != nil {
		return fmt.Errorf("failed to update or create payment transaction: %w", err)
	}
	return nil
}

// updateOrCreateTransaction attempts to update an existing pending transaction or creates a new one
func (h *MobilePayWebhookHandler) updateOrCreateTransaction(orderID uint, transactionID string, txnType entity.TransactionType, status entity.TransactionStatus, amount int64, currency, provider string, event *models.WebhookEvent) error {
	// First, try to find an existing pending transaction of the same type
	existingTxn, err := h.orderUseCase.GetLatestPendingTransactionByType(orderID, txnType)
	if err == nil && existingTxn != nil {
		// Update the existing pending transaction
		h.logger.Info("Updating existing pending %s transaction for order %d from pending to %s", txnType, orderID, status)

		// Prepare metadata and raw response from webhook event
		metadata := make(map[string]string)
		var rawResponse string

		if event != nil {
			rawResponse = h.buildEventRawResponse(event)
			metadata = h.buildEventMetadata(event)
		}

		// Update the external ID to the webhook reference
		existingTxn.ExternalID = transactionID
		if err := h.orderUseCase.UpdatePaymentTransactionStatus(existingTxn, status, rawResponse, metadata); err != nil {
			return fmt.Errorf("failed to update existing transaction: %w", err)
		}

		return nil
	}

	// No pending transaction found, create a new one (fallback for edge cases)
	h.logger.Info("No pending %s transaction found for order %d, creating new transaction with status %s", txnType, orderID, status)
	return h.createNewTransaction(orderID, transactionID, txnType, status, amount, currency, provider, event)
}

// createNewTransaction creates a completely new transaction record
func (h *MobilePayWebhookHandler) createNewTransaction(orderID uint, transactionID string, txnType entity.TransactionType, status entity.TransactionStatus, amount int64, currency, provider string, event *models.WebhookEvent) error {
	// Get idempotency key from event if available
	idempotencyKey := ""
	if event != nil {
		idempotencyKey = event.IdempotencyKey
	}

	// Create payment transaction
	txn, err := entity.NewPaymentTransaction(
		orderID,
		transactionID,
		idempotencyKey,
		txnType,
		status,
		amount,
		currency,
		provider,
	)
	if err != nil {
		return fmt.Errorf("failed to create payment transaction: %w", err)
	}

	// Add webhook event data
	if event != nil {
		txn.SetRawResponse(h.buildEventRawResponse(event))

		// Add metadata
		for key, value := range h.buildEventMetadata(event) {
			txn.AddMetadata(key, value)
		}
	}

	// Save the transaction using the usecase
	return h.orderUseCase.RecordPaymentTransaction(txn)
}

// buildEventRawResponse builds the raw response string from webhook event
func (h *MobilePayWebhookHandler) buildEventRawResponse(event *models.WebhookEvent) string {
	eventData := map[string]interface{}{
		"event_name":    string(event.Name),
		"reference":     event.Reference,
		"psp_reference": event.PSPReference,
		"timestamp":     event.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
		"success":       event.Success,
		"msn":           event.MSN,
	}
	return fmt.Sprintf("%+v", eventData)
}

// buildEventMetadata builds metadata map from webhook event
func (h *MobilePayWebhookHandler) buildEventMetadata(event *models.WebhookEvent) map[string]string {
	metadata := make(map[string]string)
	metadata["webhook_event_name"] = string(event.Name)
	metadata["webhook_psp_reference"] = event.PSPReference
	metadata["webhook_timestamp"] = event.Timestamp.Format("2006-01-02T15:04:05Z07:00")
	metadata["webhook_success"] = fmt.Sprintf("%t", event.Success)
	if event.IdempotencyKey != "" {
		metadata["idempotency_key"] = event.IdempotencyKey
	}
	return metadata
}
