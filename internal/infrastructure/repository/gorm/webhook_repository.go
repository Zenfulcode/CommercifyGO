package gorm

import (
	"errors"
	"fmt"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/gorm"
)

// WebhookRepository implements repository.WebhookRepository using GORM
// DEPRECATED: This is kept for backward compatibility. Use PaymentProviderRepository instead.
type WebhookRepository struct {
	db *gorm.DB
}

// Create implements repository.WebhookRepository.
func (r *WebhookRepository) Create(webhook *entity.Webhook) error {
	if err := webhook.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return r.db.Create(webhook).Error
}

// Update implements repository.WebhookRepository.
func (r *WebhookRepository) Update(webhook *entity.Webhook) error {
	if err := webhook.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return r.db.Save(webhook).Error
}

// Delete implements repository.WebhookRepository.
func (r *WebhookRepository) Delete(id uint) error {
	return r.db.Delete(&entity.Webhook{}, id).Error
}

// GetByID implements repository.WebhookRepository.
func (r *WebhookRepository) GetByID(id uint) (*entity.Webhook, error) {
	var webhook entity.Webhook
	if err := r.db.First(&webhook, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("webhook with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch webhook: %w", err)
	}
	return &webhook, nil
}

// GetByProvider implements repository.WebhookRepository.
func (r *WebhookRepository) GetByProvider(provider string) ([]*entity.Webhook, error) {
	var webhooks []*entity.Webhook
	if err := r.db.Where("provider = ?", provider).Find(&webhooks).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch webhooks by provider: %w", err)
	}
	return webhooks, nil
}

// GetActive implements repository.WebhookRepository.
func (r *WebhookRepository) GetActive() ([]*entity.Webhook, error) {
	var webhooks []*entity.Webhook
	if err := r.db.Where("is_active = ?", true).Find(&webhooks).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch active webhooks: %w", err)
	}
	return webhooks, nil
}

// GetByExternalID implements repository.WebhookRepository.
func (r *WebhookRepository) GetByExternalID(provider string, externalID string) (*entity.Webhook, error) {
	var webhook entity.Webhook
	if err := r.db.Where("provider = ? AND external_id = ?", provider, externalID).First(&webhook).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("webhook with provider %s and external ID %s not found", provider, externalID)
		}
		return nil, fmt.Errorf("failed to fetch webhook by external ID: %w", err)
	}
	return &webhook, nil
}

// NewWebhookRepository creates a new GORM-based WebhookRepository
func NewWebhookRepository(db *gorm.DB) repository.WebhookRepository {
	return &WebhookRepository{db: db}
}
