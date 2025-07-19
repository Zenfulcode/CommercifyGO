package gorm

import (
	"errors"
	"fmt"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/gorm"
)

// CurrencyRepository implements repository.CurrencyRepository using GORM
type CurrencyRepository struct {
	db *gorm.DB
}

// Create implements repository.CurrencyRepository.
func (c *CurrencyRepository) Create(currency *entity.Currency) error {
	return c.db.Create(currency).Error
}

// Delete implements repository.CurrencyRepository.
func (c *CurrencyRepository) Delete(code string) error {
	return c.db.Where("code = ?", code).Delete(&entity.Currency{}).Error
}

// GetByCode implements repository.CurrencyRepository.
func (c *CurrencyRepository) GetByCode(code string) (*entity.Currency, error) {
	var currency entity.Currency
	if err := c.db.Where("code = ?", code).First(&currency).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("currency with code %s not found", code)
		}
		return nil, fmt.Errorf("failed to fetch currency by code: %w", err)
	}
	return &currency, nil
}

// GetDefault implements repository.CurrencyRepository.
func (c *CurrencyRepository) GetDefault() (*entity.Currency, error) {
	var currency entity.Currency
	if err := c.db.Where("is_default = ?", true).First(&currency).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("no default currency found")
		}
		return nil, fmt.Errorf("failed to fetch default currency: %w", err)
	}
	return &currency, nil
}

// List implements repository.CurrencyRepository.
func (c *CurrencyRepository) List() ([]*entity.Currency, error) {
	var currencies []*entity.Currency
	if err := c.db.Order("code ASC").Find(&currencies).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch currencies: %w", err)
	}
	return currencies, nil
}

// ListEnabled implements repository.CurrencyRepository.
func (c *CurrencyRepository) ListEnabled() ([]*entity.Currency, error) {
	var currencies []*entity.Currency
	if err := c.db.Where("is_enabled = ?", true).Order("code ASC").Find(&currencies).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch enabled currencies: %w", err)
	}
	return currencies, nil
}

// SetDefault implements repository.CurrencyRepository.
func (c *CurrencyRepository) SetDefault(code string) error {
	return c.db.Transaction(func(tx *gorm.DB) error {
		// First, unset all currencies as default
		if err := tx.Model(&entity.Currency{}).Where("is_default = ?", true).
			Update("is_default", false).Error; err != nil {
			return fmt.Errorf("failed to unset existing default currency: %w", err)
		}

		// Then set the specified currency as default and ensure it's enabled
		if err := tx.Model(&entity.Currency{}).Where("code = ?", code).
			Updates(map[string]any{
				"is_default": true,
				"is_enabled": true,
			}).Error; err != nil {
			return fmt.Errorf("failed to set currency %s as default: %w", code, err)
		}

		return nil
	})
}

// Update implements repository.CurrencyRepository.
func (c *CurrencyRepository) Update(currency *entity.Currency) error {
	return c.db.Save(currency).Error
}

// NewCurrencyRepository creates a new GORM-based CurrencyRepository
func NewCurrencyRepository(db *gorm.DB) repository.CurrencyRepository {
	return &CurrencyRepository{db: db}
}
