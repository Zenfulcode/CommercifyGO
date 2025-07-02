package gorm

import (
	"fmt"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/gorm"
)

// ShippingZoneRepository implements repository.ShippingZoneRepository using GORM
type ShippingZoneRepository struct {
	db *gorm.DB
}

// Create implements repository.ShippingZoneRepository.
func (r *ShippingZoneRepository) Create(zone *entity.ShippingZone) error {
	return r.db.Create(zone).Error
}

// GetByID implements repository.ShippingZoneRepository.
func (r *ShippingZoneRepository) GetByID(zoneID uint) (*entity.ShippingZone, error) {
	var zone entity.ShippingZone
	if err := r.db.First(&zone, zoneID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch shipping zone: %w", err)
	}
	return &zone, nil
}

// List implements repository.ShippingZoneRepository.
func (r *ShippingZoneRepository) List(active bool) ([]*entity.ShippingZone, error) {
	var zones []*entity.ShippingZone
	query := r.db
	if active {
		query = query.Where("active = ?", true)
	}
	if err := query.Find(&zones).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch shipping zones: %w", err)
	}
	return zones, nil
}

// Update implements repository.ShippingZoneRepository.
func (r *ShippingZoneRepository) Update(zone *entity.ShippingZone) error {
	return r.db.Save(zone).Error
}

// Delete implements repository.ShippingZoneRepository.
func (r *ShippingZoneRepository) Delete(zoneID uint) error {
	return r.db.Delete(&entity.ShippingZone{}, zoneID).Error
}

// NewShippingZoneRepository creates a new GORM-based ShippingZoneRepository
func NewShippingZoneRepository(db *gorm.DB) repository.ShippingZoneRepository {
	return &ShippingZoneRepository{db: db}
}
