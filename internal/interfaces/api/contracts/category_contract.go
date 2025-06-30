package contracts

import (
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/interfaces/api/contracts/dto"
)

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

func CreateCategoryResponse(category *dto.CategoryDTO) ResponseDTO[dto.CategoryDTO] {
	return SuccessResponse(*category)
}

func CreateCategoryListResponse(categories []*entity.Category, totalCount, page, pageSize int) ListResponseDTO[dto.CategoryDTO] {
	var categoryDTOs []dto.CategoryDTO
	for _, category := range categories {
		categoryDTOs = append(categoryDTOs, *category.ToCategoryDTO())
	}

	if len(categoryDTOs) == 0 {
		return ListResponseDTO[dto.CategoryDTO]{
			Success:    true,
			Data:       []dto.CategoryDTO{},
			Pagination: PaginationDTO{Page: page, PageSize: pageSize, Total: 0},
			Message:    "No categories found",
		}
	}

	return ListResponseDTO[dto.CategoryDTO]{
		Success: true,
		Data:    categoryDTOs,
		Pagination: PaginationDTO{
			Page:     page,
			PageSize: pageSize,
			Total:    totalCount,
		},
		Message: "Categories retrieved successfully",
	}
}
