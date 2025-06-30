package dto

import (
	"time"
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

// CreateCategoryRequest represents the request to create a category
type CreateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
}

// UpdateCategoryRequest represents the request to update a category
type UpdateCategoryRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
	ParentID    *uint  `json:"parent_id"`
}

// CreateCategoryResponse represents the response after creating a category
type CreateCategoryResponse struct {
	Category CategoryDTO `json:"category"`
}
