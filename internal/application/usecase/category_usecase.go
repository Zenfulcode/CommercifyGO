package usecase

import (
	"errors"
	"fmt"
	"strings"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// CategoryUseCase implements category-related use cases
type CategoryUseCase struct {
	categoryRepo repository.CategoryRepository
}

// NewCategoryUseCase creates a new CategoryUseCase
func NewCategoryUseCase(categoryRepo repository.CategoryRepository) *CategoryUseCase {
	return &CategoryUseCase{
		categoryRepo: categoryRepo,
	}
}

type CreateCategory struct {
	Name        string `json:"name" validate:"required,min=1,max=255"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	ParentID    *uint  `json:"parent_id,omitempty"` // Optional parent category ID
}

// CreateCategory creates a new category
func (uc *CategoryUseCase) CreateCategory(input CreateCategory) (*entity.Category, error) {
	// Validate parent category exists if parentID is provided
	if input.ParentID != nil {
		parent, err := uc.categoryRepo.GetByID(*input.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent category not found: %w", err)
		}
		if parent == nil {
			return nil, errors.New("parent category does not exist")
		}
	}

	// Create new category entity
	category, err := entity.NewCategory(input.Name, input.Description, input.ParentID)
	if err != nil {
		return nil, err
	}

	// Save to repository
	if err := uc.categoryRepo.Create(category); err != nil {
		// Check for unique constraint violations
		if strings.Contains(err.Error(), "unique_root_category_name") {
			return nil, fmt.Errorf("a category with the name '%s' already exists at the root level", input.Name)
		}
		if strings.Contains(err.Error(), "unique_child_category_name_parent") {
			return nil, fmt.Errorf("a category with the name '%s' already exists under this parent category", input.Name)
		}
		if strings.Contains(err.Error(), "unique_category_name_parent") {
			return nil, fmt.Errorf("a category with the name '%s' already exists at this level", input.Name)
		}
		return nil, fmt.Errorf("failed to save category: %w", err)
	}

	// Convert to DTO
	return category, nil
}

// GetCategory retrieves a category by ID
func (uc *CategoryUseCase) GetCategory(categoryID uint) (*entity.Category, error) {
	category, err := uc.categoryRepo.GetByID(categoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	return category, nil
}

type UpdateCategory struct {
	CategoryID  uint   `json:"category_id" validate:"required"`
	Name        string `json:"name,omitempty" validate:"omitempty,min=1,max=255"`
	Description string `json:"description,omitempty" validate:"max=1000"`
	ParentID    *uint  `json:"parent_id,omitempty"` // Optional parent category ID
}

// UpdateCategory updates an existing category
func (uc *CategoryUseCase) UpdateCategory(input UpdateCategory) (*entity.Category, error) {
	// Get existing category
	category, err := uc.categoryRepo.GetByID(input.CategoryID)
	if err != nil {
		return nil, fmt.Errorf("failed to get category: %w", err)
	}

	// Validate parent category exists if parentID is provided
	if input.ParentID != nil {
		// Check for circular reference (category cannot be its own parent)
		if *input.ParentID == input.CategoryID {
			return nil, errors.New("category cannot be its own parent")
		}

		parent, err := uc.categoryRepo.GetByID(*input.ParentID)
		if err != nil {
			return nil, fmt.Errorf("parent category not found: %w", err)
		}
		if parent == nil {
			return nil, errors.New("parent category does not exist")
		}
	}

	// Update fields if provided
	if input.Name != "" {
		category.Name = input.Name
	}
	if input.Description != "" {
		category.Description = input.Description
	}
	if input.ParentID != nil {
		category.ParentID = input.ParentID
	}

	// Save updated category
	if err := uc.categoryRepo.Update(category); err != nil {
		// Check for unique constraint violations
		if strings.Contains(err.Error(), "unique_root_category_name") {
			return nil, fmt.Errorf("a category with the name '%s' already exists at the root level", category.Name)
		}
		if strings.Contains(err.Error(), "unique_child_category_name_parent") {
			return nil, fmt.Errorf("a category with the name '%s' already exists under this parent category", category.Name)
		}
		if strings.Contains(err.Error(), "unique_category_name_parent") {
			return nil, fmt.Errorf("a category with the name '%s' already exists at this level", category.Name)
		}
		return nil, fmt.Errorf("failed to update category: %w", err)
	}

	return category, nil
}

// DeleteCategory deletes a category by ID
func (uc *CategoryUseCase) DeleteCategory(categoryID uint) error {
	// Check if category exists
	_, err := uc.categoryRepo.GetByID(categoryID)
	if err != nil {
		return fmt.Errorf("category not found")
	}

	// Check if category has children
	children, err := uc.categoryRepo.GetChildren(categoryID)
	if err != nil {
		return fmt.Errorf("failed to check for child categories: %w", err)
	}

	if len(children) > 0 {
		return errors.New("cannot delete category with child categories")
	}

	// Delete the category
	if err := uc.categoryRepo.Delete(categoryID); err != nil {
		return fmt.Errorf("failed to delete category: %w", err)
	}

	return nil
}

// ListCategories retrieves all categories
func (uc *CategoryUseCase) ListCategories() ([]*entity.Category, error) {
	categories, err := uc.categoryRepo.List()
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, nil
}

// GetChildCategories retrieves all child categories of a parent category
func (uc *CategoryUseCase) GetChildCategories(parentID uint) ([]*entity.Category, error) {
	// Check if parent category exists
	_, err := uc.categoryRepo.GetByID(parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get parent category: %w", err)
	}

	// Get child categories
	children, err := uc.categoryRepo.GetChildren(parentID)
	if err != nil {
		return nil, fmt.Errorf("failed to get child categories: %w", err)
	}

	return children, nil
}
