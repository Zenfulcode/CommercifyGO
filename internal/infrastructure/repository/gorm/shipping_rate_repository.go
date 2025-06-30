package gorm

import (
	"fmt"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/gorm"
)

// ShippingRateRepository implements repository.ShippingRateRepository using GORM
type ShippingRateRepository struct {
	db *gorm.DB
}

// Create implements repository.ShippingRateRepository.
func (r *ShippingRateRepository) Create(rate *entity.ShippingRate) error {
	return r.db.Create(rate).Error
}

// GetByID implements repository.ShippingRateRepository.
func (r *ShippingRateRepository) GetByID(rateID uint) (*entity.ShippingRate, error) {
	var rate entity.ShippingRate
	if err := r.db.First(&rate, rateID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch shipping rate: %w", err)
	}
	return &rate, nil
}

// GetByMethodID implements repository.ShippingRateRepository.
func (r *ShippingRateRepository) GetByMethodID(methodID uint) ([]*entity.ShippingRate, error) {
	var rates []*entity.ShippingRate
	if err := r.db.Where("method_id = ?", methodID).Find(&rates).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch shipping rates by method: %w", err)
	}
	return rates, nil
}

// GetByZoneID implements repository.ShippingRateRepository.
func (r *ShippingRateRepository) GetByZoneID(zoneID uint) ([]*entity.ShippingRate, error) {
	var rates []*entity.ShippingRate
	if err := r.db.Where("zone_id = ?", zoneID).Find(&rates).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch shipping rates by zone: %w", err)
	}
	return rates, nil
}

// GetAvailableRatesForAddress implements repository.ShippingRateRepository.
func (r *ShippingRateRepository) GetAvailableRatesForAddress(address entity.Address, orderValue int64) ([]*entity.ShippingRate, error) {
	// This is a placeholder implementation
	var rates []*entity.ShippingRate
	if err := r.db.Find(&rates).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch available shipping rates: %w", err)
	}
	return rates, nil
}

// CreateWeightBasedRate implements repository.ShippingRateRepository.
func (r *ShippingRateRepository) CreateWeightBasedRate(weightRate *entity.WeightBasedRate) error {
	return r.db.Create(weightRate).Error
}

// CreateValueBasedRate implements repository.ShippingRateRepository.
func (r *ShippingRateRepository) CreateValueBasedRate(valueRate *entity.ValueBasedRate) error {
	return r.db.Create(valueRate).Error
}

// GetWeightBasedRates implements repository.ShippingRateRepository.
func (r *ShippingRateRepository) GetWeightBasedRates(rateID uint) ([]entity.WeightBasedRate, error) {
	var rates []entity.WeightBasedRate
	if err := r.db.Where("rate_id = ?", rateID).Find(&rates).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch weight-based rates: %w", err)
	}
	return rates, nil
}

// GetValueBasedRates implements repository.ShippingRateRepository.
func (r *ShippingRateRepository) GetValueBasedRates(rateID uint) ([]entity.ValueBasedRate, error) {
	var rates []entity.ValueBasedRate
	if err := r.db.Where("rate_id = ?", rateID).Find(&rates).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch value-based rates: %w", err)
	}
	return rates, nil
}

// Update implements repository.ShippingRateRepository.
func (r *ShippingRateRepository) Update(rate *entity.ShippingRate) error {
	return r.db.Save(rate).Error
}

// Delete implements repository.ShippingRateRepository.
func (r *ShippingRateRepository) Delete(rateID uint) error {
	return r.db.Delete(&entity.ShippingRate{}, rateID).Error
}

// NewShippingRateRepository creates a new GORM-based ShippingRateRepository
func NewShippingRateRepository(db *gorm.DB) repository.ShippingRateRepository {
	return &ShippingRateRepository{db: db}
}
