package gorm

import (
	"errors"
	"fmt"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/gorm"
)

// CategoryRepository implements repository.CategoryRepository using GORM
type CategoryRepository struct {
	db *gorm.DB
}

// Create implements repository.CategoryRepository.
func (c *CategoryRepository) Create(category *entity.Category) error {
	return c.db.Create(category).Error
}

// Delete implements repository.CategoryRepository.
func (c *CategoryRepository) Delete(categoryID uint) error {
	// Note: This will fail if there are products in this category due to RESTRICT constraint
	// which is the intended behavior for data integrity
	return c.db.Unscoped().Delete(&entity.Category{}, categoryID).Error
}

// GetByID implements repository.CategoryRepository.
func (c *CategoryRepository) GetByID(categoryID uint) (*entity.Category, error) {
	var category entity.Category
	if err := c.db.Preload("Parent").Preload("Children").First(&category, categoryID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("category with ID %d not found", categoryID)
		}
		return nil, fmt.Errorf("failed to fetch category: %w", err)
	}
	return &category, nil
}

// GetChildren implements repository.CategoryRepository.
func (c *CategoryRepository) GetChildren(parentID uint) ([]*entity.Category, error) {
	var children []*entity.Category
	if err := c.db.Preload("Parent").Preload("Children").
		Where("parent_id = ?", parentID).
		Order("name ASC").
		Find(&children).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch children for category %d: %w", parentID, err)
	}
	return children, nil
}

// List implements repository.CategoryRepository.
func (c *CategoryRepository) List() ([]*entity.Category, error) {
	var categories []*entity.Category
	if err := c.db.Preload("Parent").Preload("Children").Order("name ASC").Find(&categories).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch categories: %w", err)
	}
	return categories, nil
}

// Update implements repository.CategoryRepository.
func (c *CategoryRepository) Update(category *entity.Category) error {
	// Use explicit field updates to ensure parent_id is properly updated when nil
	return c.db.Model(category).Select("name", "description", "parent_id").Updates(category).Error
}

// NewCategoryRepository creates a new GORM-based CategoryRepository
func NewCategoryRepository(db *gorm.DB) repository.CategoryRepository {
	return &CategoryRepository{db: db}
}
