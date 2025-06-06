package dto

import (
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// CategoryDTO represents a category in the system
type CategoryDTO struct {
	ID          uint      `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	ParentID    *uint     `json:"parent_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateCategoryRequest represents the data needed to create a new category
type CreateCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id,omitempty"`
}

// UpdateCategoryRequest represents the data needed to update an existing category
type UpdateCategoryRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ParentID    *uint  `json:"parent_id,omitempty"`
}

func toCategoryDTOList(categories []*entity.Category) []CategoryDTO {
	var categoryDTOs []CategoryDTO
	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, CategoryDTO{
			ID:          category.ID,
			Name:        category.Name,
			Description: category.Description,
			ParentID:    category.ParentID,
			CreatedAt:   category.CreatedAt,
			UpdatedAt:   category.UpdatedAt,
		})
	}
	return categoryDTOs
}

func CreateCategoryResponse(category *entity.Category) ResponseDTO[CategoryDTO] {
	return SuccessResponse(CategoryDTO{
		ID:          category.ID,
		Name:        category.Name,
		Description: category.Description,
		ParentID:    category.ParentID,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	})
}

func CreateCategoryListResponse(categories []*entity.Category, totalCount, page, pageSize int) ListResponseDTO[CategoryDTO] {
	return ListResponseDTO[CategoryDTO]{
		Success: true,
		Data:    toCategoryDTOList(categories),
		Pagination: PaginationDTO{
			Page:     page,
			PageSize: pageSize,
			Total:    totalCount,
		},
	}
}
