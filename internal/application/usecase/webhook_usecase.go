package usecase

import (
	"errors"

	"github.com/gkhaavik/vipps-mobilepay-sdk/pkg/models"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"github.com/zenfulcode/commercify/internal/infrastructure/payment"
)

// WebhookUseCase handles webhook-related operations
type WebhookUseCase struct {
	webhookRepo    repository.WebhookRepository
	webhookService *payment.WebhookService
}

// RegisterWebhookInput represents the input for registering a webhook
type RegisterWebhookInput struct {
	Provider string   `json:"provider"`
	URL      string   `json:"url"`
	Events   []string `json:"events"`
}

// NewWebhookUseCase creates a new WebhookUseCase
func NewWebhookUseCase(webhookRepo repository.WebhookRepository, webhookService *payment.WebhookService) *WebhookUseCase {
	return &WebhookUseCase{
		webhookRepo:    webhookRepo,
		webhookService: webhookService,
	}
}

// RegisterMobilePayWebhook registers a webhook with MobilePay
func (u *WebhookUseCase) RegisterMobilePayWebhook(input RegisterWebhookInput) (*entity.Webhook, error) {
	// Validate input
	if input.URL == "" {
		return nil, errors.New("URL is required")
	}
	if len(input.Events) == 0 {
		return nil, errors.New("at least one event is required")
	}

	// Register webhook with MobilePay
	return u.webhookService.RegisterMobilePayWebhook(input.URL, input.Events)
}

func (u *WebhookUseCase) DeleteMobilePayWebhook(externalID string) error {
	// Validate ID
	if externalID == "" {
		return errors.New("external ID is required")
	}

	// Delete webhook from MobilePay
	return u.webhookService.ForceDeleteMobilePayWebhook(externalID)
}

// DeleteWebhook deletes a webhook
func (u *WebhookUseCase) DeleteWebhook(id uint) error {
	webhook, err := u.webhookRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Delete from provider if supported
	if webhook.Provider == "mobilepay" {
		return u.webhookService.DeleteMobilePayWebhook(webhook.ExternalID)
	}

	// Otherwise just delete from our database
	return u.webhookRepo.Delete(id)
}

// GetWebhookByID returns a webhook by ID
func (u *WebhookUseCase) GetWebhookByID(id uint) (*entity.Webhook, error) {
	return u.webhookRepo.GetByID(id)
}

// GetAllWebhooks returns all webhooks
func (u *WebhookUseCase) GetAllWebhooks() ([]*entity.Webhook, error) {
	return u.webhookRepo.GetActive()
}

// GetMobilePayWebhooks returns all MobilePay webhooks
func (u *WebhookUseCase) GetMobilePayWebhooks() ([]models.WebhookRegistration, error) {
	return u.webhookService.GetMobilePayWebhooks()
}
