package gorm

import (
	"errors"
	"fmt"

	"github.com/zenfulcode/commercify/internal/domain/common"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/gorm"
)

// PaymentProviderRepository implements repository.PaymentProviderRepository using GORM
type PaymentProviderRepository struct {
	db *gorm.DB
}

// Create implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) Create(provider *entity.PaymentProvider) error {
	if err := provider.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return r.db.Create(provider).Error
}

// Update implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) Update(provider *entity.PaymentProvider) error {
	if err := provider.Validate(); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}
	return r.db.Save(provider).Error
}

// Delete implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) Delete(id uint) error {
	return r.db.Delete(&entity.PaymentProvider{}, id).Error
}

// GetByID implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) GetByID(id uint) (*entity.PaymentProvider, error) {
	var provider entity.PaymentProvider
	if err := r.db.First(&provider, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("payment provider with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch payment provider: %w", err)
	}
	return &provider, nil
}

// GetByType implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) GetByType(providerType common.PaymentProviderType) (*entity.PaymentProvider, error) {
	var provider entity.PaymentProvider
	if err := r.db.Where("type = ?", providerType).First(&provider).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("payment provider with type %s not found", providerType)
		}
		return nil, fmt.Errorf("failed to fetch payment provider by type: %w", err)
	}
	return &provider, nil
}

// List implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) List(offset, limit int) ([]*entity.PaymentProvider, error) {
	var providers []*entity.PaymentProvider
	query := r.db.Order("priority DESC, created_at ASC")

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	if err := query.Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch payment providers: %w", err)
	}
	return providers, nil
}

// GetEnabled implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) GetEnabled() ([]*entity.PaymentProvider, error) {
	var providers []*entity.PaymentProvider
	if err := r.db.Where("enabled = ?", true).Order("priority DESC, created_at ASC").Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch enabled payment providers: %w", err)
	}
	return providers, nil
}

// GetEnabledByMethod implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) GetEnabledByMethod(method common.PaymentMethod) ([]*entity.PaymentProvider, error) {
	var providers []*entity.PaymentProvider
	// Use SQLite-compatible JSON query to check if the method exists in the methods array
	if err := r.db.Where("enabled = ? AND json_extract(methods, '$') LIKE ?", true, "%\""+string(method)+"\"%").
		Order("priority DESC, created_at ASC").Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch payment providers by method: %w", err)
	}
	return providers, nil
}

// GetEnabledByCurrency implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) GetEnabledByCurrency(currency string) ([]*entity.PaymentProvider, error) {
	var providers []*entity.PaymentProvider
	// Use SQLite-compatible JSON query to check if the currency exists in the supported_currencies array
	// If supported_currencies is empty/null, include the provider (supports all currencies)
	if err := r.db.Where("enabled = ? AND (supported_currencies IS NULL OR supported_currencies = '[]' OR json_extract(supported_currencies, '$') LIKE ?)",
		true, "%\""+currency+"\"%").
		Order("priority DESC, created_at ASC").Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch payment providers by currency: %w", err)
	}
	return providers, nil
}

// GetEnabledByMethodAndCurrency implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) GetEnabledByMethodAndCurrency(method common.PaymentMethod, currency string) ([]*entity.PaymentProvider, error) {
	var providers []*entity.PaymentProvider
	// Combine both method and currency filters using SQLite-compatible JSON queries
	if err := r.db.Where("enabled = ? AND json_extract(methods, '$') LIKE ? AND (supported_currencies IS NULL OR supported_currencies = '[]' OR json_extract(supported_currencies, '$') LIKE ?)",
		true, "%\""+string(method)+"\"%", "%\""+currency+"\"%").
		Order("priority DESC, created_at ASC").Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch payment providers by method and currency: %w", err)
	}
	return providers, nil
}

// UpdateWebhookInfo implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) UpdateWebhookInfo(providerType common.PaymentProviderType, webhookURL, webhookSecret, externalWebhookID string, events []string) error {
	updates := map[string]any{
		"webhook_url":         webhookURL,
		"webhook_secret":      webhookSecret,
		"external_webhook_id": externalWebhookID,
		"webhook_events":      events,
	}

	result := r.db.Model(&entity.PaymentProvider{}).Where("type = ?", providerType).Updates(updates)
	if result.Error != nil {
		return fmt.Errorf("failed to update webhook info: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("payment provider with type %s not found", providerType)
	}

	return nil
}

// GetWithWebhooks implements repository.PaymentProviderRepository.
func (r *PaymentProviderRepository) GetWithWebhooks() ([]*entity.PaymentProvider, error) {
	var providers []*entity.PaymentProvider
	if err := r.db.Where("webhook_url IS NOT NULL AND webhook_url != ''").
		Order("priority DESC, created_at ASC").Find(&providers).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch payment providers with webhooks: %w", err)
	}
	return providers, nil
}

// NewPaymentProviderRepository creates a new GORM-based PaymentProviderRepository
func NewPaymentProviderRepository(db *gorm.DB) repository.PaymentProviderRepository {
	return &PaymentProviderRepository{db: db}
}
