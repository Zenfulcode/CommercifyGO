package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/service"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"gorm.io/datatypes"
)

// PaymentProviderHandler handles payment provider management requests
type PaymentProviderHandler struct {
	paymentProviderService service.PaymentProviderService
	logger                 logger.Logger
}

// NewPaymentProviderHandler creates a new PaymentProviderHandler
func NewPaymentProviderHandler(paymentProviderService service.PaymentProviderService, logger logger.Logger) *PaymentProviderHandler {
	return &PaymentProviderHandler{
		paymentProviderService: paymentProviderService,
		logger:                 logger,
	}
}

// GetPaymentProviders handles getting all payment providers (admin only)
func (h *PaymentProviderHandler) GetPaymentProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := h.paymentProviderService.GetPaymentProviders()
	if err != nil {
		h.logger.Error("Failed to get payment providers: %v", err)
		http.Error(w, "Failed to get payment providers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providers)
}

// GetEnabledPaymentProviders handles getting only enabled payment providers
func (h *PaymentProviderHandler) GetEnabledPaymentProviders(w http.ResponseWriter, r *http.Request) {
	providers, err := h.paymentProviderService.GetEnabledPaymentProviders()
	if err != nil {
		h.logger.Error("Failed to get enabled payment providers: %v", err)
		http.Error(w, "Failed to get enabled payment providers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(providers)
}

// EnableProviderRequest represents a request to enable a payment provider
type EnableProviderRequest struct {
	Enabled bool `json:"enabled"`
}

// EnablePaymentProvider handles enabling a payment provider
func (h *PaymentProviderHandler) EnablePaymentProvider(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	providerType := common.PaymentProviderType(vars["providerType"])

	var req EnableProviderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to parse request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var err error
	if req.Enabled {
		err = h.paymentProviderService.EnableProvider(providerType)
	} else {
		err = h.paymentProviderService.DisableProvider(providerType)
	}

	if err != nil {
		h.logger.Error("Failed to update provider %s: %v", providerType, err)
		http.Error(w, "Failed to update payment provider", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Payment provider updated successfully",
		"enabled": req.Enabled,
	})
}

// UpdateConfigurationRequest represents a request to update provider configuration
type UpdateConfigurationRequest struct {
	Configuration datatypes.JSONMap `json:"configuration"`
}

// UpdateProviderConfiguration handles updating a payment provider's configuration
func (h *PaymentProviderHandler) UpdateProviderConfiguration(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	providerType := common.PaymentProviderType(vars["providerType"])

	var req UpdateConfigurationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to parse request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err := h.paymentProviderService.UpdateProviderConfiguration(providerType, req.Configuration)
	if err != nil {
		h.logger.Error("Failed to update configuration for provider %s: %v", providerType, err)
		http.Error(w, "Failed to update provider configuration", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Provider configuration updated successfully",
	})
}

// ProviderWebhookRequest represents a request to register a webhook
type ProviderWebhookRequest struct {
	URL    string   `json:"url"`
	Events []string `json:"events"`
}

// RegisterWebhook handles registering a webhook for a payment provider
func (h *PaymentProviderHandler) RegisterWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	providerType := common.PaymentProviderType(vars["providerType"])

	var req ProviderWebhookRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to parse request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.URL == "" {
		http.Error(w, "Webhook URL is required", http.StatusBadRequest)
		return
	}

	err := h.paymentProviderService.RegisterWebhook(providerType, req.URL, req.Events)
	if err != nil {
		h.logger.Error("Failed to register webhook for provider %s: %v", providerType, err)
		http.Error(w, "Failed to register webhook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Webhook registered successfully",
	})
}

// DeleteWebhook handles deleting a webhook for a payment provider
func (h *PaymentProviderHandler) DeleteWebhook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	providerType := common.PaymentProviderType(vars["providerType"])

	err := h.paymentProviderService.DeleteWebhook(providerType)
	if err != nil {
		h.logger.Error("Failed to delete webhook for provider %s: %v", providerType, err)
		http.Error(w, "Failed to delete webhook", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"message": "Webhook deleted successfully",
	})
}

// GetWebhookInfo handles getting webhook information for a payment provider
func (h *PaymentProviderHandler) GetWebhookInfo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	providerType := common.PaymentProviderType(vars["providerType"])

	provider, err := h.paymentProviderService.GetWebhookInfo(providerType)
	if err != nil {
		h.logger.Error("Failed to get webhook info for provider %s: %v", providerType, err)
		http.Error(w, "Failed to get webhook info", http.StatusInternalServerError)
		return
	}

	// Return only webhook-related information
	webhookInfo := map[string]any{
		"provider_type":       provider.Type,
		"webhook_url":         provider.WebhookURL,
		"webhook_secret":      provider.WebhookSecret,
		"webhook_events":      provider.WebhookEvents,
		"external_webhook_id": provider.ExternalWebhookID,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(webhookInfo)
}
