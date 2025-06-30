package gorm

import (
	"fmt"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/gorm"
)

// ShippingMethodRepository implements repository.ShippingMethodRepository using GORM
type ShippingMethodRepository struct {
	db *gorm.DB
}

// Create implements repository.ShippingMethodRepository.
func (r *ShippingMethodRepository) Create(method *entity.ShippingMethod) error {
	return r.db.Create(method).Error
}

// GetByID implements repository.ShippingMethodRepository.
func (r *ShippingMethodRepository) GetByID(methodID uint) (*entity.ShippingMethod, error) {
	var method entity.ShippingMethod
	if err := r.db.First(&method, methodID).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch shipping method: %w", err)
	}
	return &method, nil
}

// List implements repository.ShippingMethodRepository.
func (r *ShippingMethodRepository) List(active bool) ([]*entity.ShippingMethod, error) {
	var methods []*entity.ShippingMethod
	query := r.db
	if active {
		query = query.Where("active = ?", true)
	}
	if err := query.Find(&methods).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch shipping methods: %w", err)
	}
	return methods, nil
}

// Update implements repository.ShippingMethodRepository.
func (r *ShippingMethodRepository) Update(method *entity.ShippingMethod) error {
	return r.db.Save(method).Error
}

// Delete implements repository.ShippingMethodRepository.
func (r *ShippingMethodRepository) Delete(methodID uint) error {
	return r.db.Delete(&entity.ShippingMethod{}, methodID).Error
}

// NewShippingMethodRepository creates a new GORM-based ShippingMethodRepository
func NewShippingMethodRepository(db *gorm.DB) repository.ShippingMethodRepository {
	return &ShippingMethodRepository{db: db}
}
