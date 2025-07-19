package gorm

import (
	"errors"
	"fmt"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/gorm"
)

// DiscountRepository implements repository.DiscountRepository using GORM
type DiscountRepository struct {
	db *gorm.DB
}

// Create implements repository.DiscountRepository.
func (d *DiscountRepository) Create(discount *entity.Discount) error {
	return d.db.Create(discount).Error
}

// Delete implements repository.DiscountRepository.
func (d *DiscountRepository) Delete(discountID uint) error {
	return d.db.Delete(&entity.Discount{}, discountID).Error
}

// GetByCode implements repository.DiscountRepository.
func (d *DiscountRepository) GetByCode(code string) (*entity.Discount, error) {
	var discount entity.Discount
	if err := d.db.Where("code = ?", code).First(&discount).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("discount with code %s not found", code)
		}
		return nil, fmt.Errorf("failed to fetch discount by code: %w", err)
	}
	return &discount, nil
}

// GetByID implements repository.DiscountRepository.
func (d *DiscountRepository) GetByID(discountID uint) (*entity.Discount, error) {
	var discount entity.Discount
	if err := d.db.First(&discount, discountID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("discount with ID %d not found", discountID)
		}
		return nil, fmt.Errorf("failed to fetch discount: %w", err)
	}
	return &discount, nil
}

// IncrementUsage implements repository.DiscountRepository.
func (d *DiscountRepository) IncrementUsage(discountID uint) error {
	return d.db.Model(&entity.Discount{}).Where("id = ?", discountID).
		UpdateColumn("current_usage", gorm.Expr("current_usage + ?", 1)).Error
}

// List implements repository.DiscountRepository.
func (d *DiscountRepository) List(offset int, limit int) ([]*entity.Discount, error) {
	var discounts []*entity.Discount
	if err := d.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&discounts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch discounts: %w", err)
	}
	return discounts, nil
}

// ListActive implements repository.DiscountRepository.
func (d *DiscountRepository) ListActive(offset int, limit int) ([]*entity.Discount, error) {
	var discounts []*entity.Discount
	now := time.Now()

	if err := d.db.Where("active = ? AND start_date <= ? AND end_date >= ? AND (usage_limit = 0 OR current_usage < usage_limit)",
		true, now, now).
		Offset(offset).Limit(limit).
		Order("created_at DESC").
		Find(&discounts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch active discounts: %w", err)
	}
	return discounts, nil
}

// Update implements repository.DiscountRepository.
func (d *DiscountRepository) Update(discount *entity.Discount) error {
	return d.db.Save(discount).Error
}

// NewDiscountRepository creates a new GORM-based DiscountRepository
func NewDiscountRepository(db *gorm.DB) repository.DiscountRepository {
	return &DiscountRepository{db: db}
}
