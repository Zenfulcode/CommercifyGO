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

	// Record the authorization transaction
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

	// Record the capture transaction
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

	// Record the cancellation transaction
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

	// Record the expiration as a failed transaction
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

	// Record the refund transaction
	// Use the amount from the webhook event if available, otherwise use order amount
	refundAmount := order.FinalAmount
	if event.Amount.Value > 0 {
		refundAmount = int64(event.Amount.Value)
	}

	if err := h.recordPaymentTransaction(orderID, event.Reference, entity.TransactionTypeRefund, entity.TransactionStatusSuccessful, refundAmount, order.Currency, "mobilepay", event); err != nil {
		h.logger.Error("Failed to record refund transaction for order %d: %v", orderID, err)
		// Don't fail the webhook processing if transaction recording fails
	}

	// Update order payment status to refunded
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
	// Create payment transaction
	txn, err := entity.NewPaymentTransaction(
		orderID,
		transactionID,
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
		eventData := map[string]interface{}{
			"event_name":    string(event.Name),
			"reference":     event.Reference,
			"psp_reference": event.PSPReference,
			"timestamp":     event.Timestamp.Format("2006-01-02T15:04:05Z07:00"),
			"success":       event.Success,
			"msn":           event.MSN,
		}

		// Convert to JSON string for storage
		eventJSON := fmt.Sprintf("%+v", eventData)
		txn.SetRawResponse(eventJSON)

		// Add metadata
		txn.AddMetadata("webhook_event_name", string(event.Name))
		txn.AddMetadata("webhook_psp_reference", event.PSPReference)
		txn.AddMetadata("webhook_timestamp", event.Timestamp.Format("2006-01-02T15:04:05Z07:00"))
		txn.AddMetadata("webhook_success", fmt.Sprintf("%t", event.Success))
		if event.IdempotencyKey != "" {
			txn.AddMetadata("idempotency_key", event.IdempotencyKey)
		}
	}

	// Save the transaction using the usecase
	return h.orderUseCase.RecordPaymentTransaction(txn)
}
