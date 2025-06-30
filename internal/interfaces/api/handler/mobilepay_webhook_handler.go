package handler

import (
	"fmt"
	"net/http"

	"github.com/gkhaavik/vipps-mobilepay-sdk/pkg/models"
	"github.com/gkhaavik/vipps-mobilepay-sdk/pkg/webhooks"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
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

	// Update order payment status to authorized (reserved but not captured)
	return h.orderUseCase.UpdateOrderPaymentStatus(orderID, entity.PaymentStatusAuthorized, event.Reference)
}

// handleSDKPaymentCaptured handles payment captured events from the SDK
func (h *MobilePayWebhookHandler) handleSDKPaymentCaptured(event *models.WebhookEvent) error {
	orderID, err := h.getOrderIDFromSDKEvent(event)
	if err != nil {
		return err
	}

	h.logger.Info("Processing captured MobilePay payment for order %d, transaction %s", orderID, event.Reference)

	// Update order payment status to captured
	return h.orderUseCase.UpdateOrderPaymentStatus(orderID, entity.PaymentStatusCaptured, event.Reference)
}

// handleSDKPaymentCancelled handles payment cancelled events from the SDK
func (h *MobilePayWebhookHandler) handleSDKPaymentCancelled(event *models.WebhookEvent) error {
	orderID, err := h.getOrderIDFromSDKEvent(event)
	if err != nil {
		return err
	}

	h.logger.Info("Processing cancelled MobilePay payment for order %d, transaction %s", orderID, event.Reference)

	// Update order payment status to cancelled
	return h.orderUseCase.UpdateOrderPaymentStatus(orderID, entity.PaymentStatusCancelled, event.Reference)
}

// handleSDKPaymentExpired handles payment expired events from the SDK
func (h *MobilePayWebhookHandler) handleSDKPaymentExpired(event *models.WebhookEvent) error {
	orderID, err := h.getOrderIDFromSDKEvent(event)
	if err != nil {
		return err
	}

	h.logger.Info("Processing expired MobilePay payment for order %d, transaction %s", orderID, event.Reference)

	// Update order payment status to failed
	return h.orderUseCase.UpdateOrderPaymentStatus(orderID, entity.PaymentStatusFailed, event.Reference)
}

// handleSDKPaymentRefunded handles payment refunded events from the SDK
func (h *MobilePayWebhookHandler) handleSDKPaymentRefunded(event *models.WebhookEvent) error {
	orderID, err := h.getOrderIDFromSDKEvent(event)
	if err != nil {
		return err
	}

	h.logger.Info("Processing refunded MobilePay payment for order %d, transaction %s", orderID, event.Reference)

	// Update order payment status to refunded
	return h.orderUseCase.UpdateOrderPaymentStatus(orderID, entity.PaymentStatusRefunded, event.Reference)
}

// getOrderIDFromSDKEvent gets the order ID associated with a MobilePay payment from SDK event
func (h *MobilePayWebhookHandler) getOrderIDFromSDKEvent(event *models.WebhookEvent) (uint, error) {
	// Try to find the order by PaymentID field using the event reference
	order, err := h.orderUseCase.GetOrderByPaymentID(event.Reference)
	if err == nil && order != nil {
		return order.ID, nil
	}

	h.logger.Error("Could not find order for MobilePay payment %s", event.Reference)
	return 0, fmt.Errorf("order not found for MobilePay payment %s", event.Reference)
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
