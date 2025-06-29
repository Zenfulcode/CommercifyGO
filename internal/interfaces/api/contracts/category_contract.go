package contracts

import (
	"github.com/zenfulcode/commercify/internal/dto"
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

func CreateCategoryResponse(category dto.CategoryDTO) ResponseDTO[dto.CategoryDTO] {
	return SuccessResponse(category)
}

func CreateCategoryListResponse(categories []dto.CategoryDTO, totalCount, page, pageSize int) ListResponseDTO[dto.CategoryDTO] {
	return ListResponseDTO[dto.CategoryDTO]{
		Success: true,
		Data:    categories,
		Pagination: PaginationDTO{
			Page:     page,
			PageSize: pageSize,
			Total:    totalCount,
		},
	}
}
