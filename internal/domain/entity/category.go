package entity

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/dto"
	"gorm.io/gorm"
)

// Category represents a product category
type Category struct {
	gorm.Model
	Name        string     `gorm:"not null;size:255"`
	Description string     `gorm:"type:text"`
	ParentID    *uint      `gorm:"index"` // Nullable for top-level categories
	Parent      *Category  `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Children    []Category `gorm:"foreignKey:ParentID;constraint:OnDelete:SET NULL,OnUpdate:CASCADE"`
	Products    []Product  `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT,OnUpdate:CASCADE"`
}

// NewCategory creates a new category
func NewCategory(name, description string, parentID *uint) (*Category, error) {
	if name == "" {
		return nil, errors.New("category name cannot be empty")
	}

	if len(name) > 255 {
		return nil, errors.New("category name cannot exceed 255 characters")
	}

	if parentID != nil && *parentID == 0 {
		return nil, errors.New("parent ID cannot be zero")
	}

	if len(description) > 65535 {
		return nil, errors.New("category description cannot exceed 65535 characters")
	}

	return &Category{
		Name:        name,
		Description: description,
		ParentID:    parentID,
	}, nil
}

func (c *Category) ToCategoryDTO() *dto.CategoryDTO {
	return &dto.CategoryDTO{
		ID:          c.ID,
		Name:        c.Name,
		Description: c.Description,
		ParentID:    c.ParentID,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}
