package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/interfaces/api/contracts"
)

// CategoryHandler handles category-related HTTP requests
type CategoryHandler struct {
	categoryUseCase *usecase.CategoryUseCase
	logger          logger.Logger
}

// NewCategoryHandler creates a new CategoryHandler
func NewCategoryHandler(categoryUseCase *usecase.CategoryUseCase, logger logger.Logger) *CategoryHandler {
	return &CategoryHandler{
		categoryUseCase: categoryUseCase,
		logger:          logger,
	}
}

// CreateCategory handles creating a new category (admin only)
func (h *CategoryHandler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var req contracts.CreateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode create category request: %v", err)
		response := contracts.ErrorResponse("Invalid request body")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	input := usecase.CreateCategory{
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
	}

	category, err := h.categoryUseCase.CreateCategory(input)
	if err != nil {
		h.logger.Error("Failed to create category: %v", err)

		// Handle specific error cases
		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to create category"

		if err.Error() == "parent category does not exist" || err.Error() == "parent category not found" {
			statusCode = http.StatusBadRequest
			errorMessage = err.Error()
		}

		if strings.Contains(err.Error(), "a category with the name") {
			statusCode = http.StatusConflict
			errorMessage = err.Error()
		}

		response := contracts.ErrorResponse(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := contracts.CreateCategoryResponse(category.ToCategoryDTO())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetCategory handles retrieving a category by ID
func (h *CategoryHandler) GetCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid category ID: %v", err)
		response := contracts.ErrorResponse("Invalid category ID")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	category, err := h.categoryUseCase.GetCategory(uint(categoryID))
	if err != nil {
		h.logger.Error("Failed to get category: %v", err)

		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to retrieve category"

		if err.Error() == "category not found" {
			statusCode = http.StatusNotFound
			errorMessage = "Category not found"
		}

		response := contracts.ErrorResponse(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := contracts.CreateCategoryResponse(category.ToCategoryDTO())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// UpdateCategory handles updating an existing category (admin only)
func (h *CategoryHandler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid category ID: %v", err)
		response := contracts.ErrorResponse("Invalid category ID")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	var req contracts.UpdateCategoryRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Failed to decode update category request: %v", err)
		response := contracts.ErrorResponse("Invalid request body")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	input := usecase.UpdateCategory{
		CategoryID:  uint(categoryID),
		Name:        req.Name,
		Description: req.Description,
		ParentID:    req.ParentID,
	}

	category, err := h.categoryUseCase.UpdateCategory(input)
	if err != nil {
		h.logger.Error("Failed to update category: %v", err)

		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to update category"

		if err.Error() == "category not found" {
			statusCode = http.StatusNotFound
			errorMessage = "Category not found"
		} else if err.Error() == "category cannot be its own parent" ||
			err.Error() == "parent category does not exist" {
			statusCode = http.StatusBadRequest
			errorMessage = err.Error()
		}

		response := contracts.ErrorResponse(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := contracts.CreateCategoryResponse(category.ToCategoryDTO())

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// DeleteCategory handles deleting a category (admin only)
func (h *CategoryHandler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	categoryID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid category ID: %v", err)
		response := contracts.ErrorResponse("Invalid category ID")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	err = h.categoryUseCase.DeleteCategory(uint(categoryID))
	if err != nil {
		h.logger.Error("Failed to delete category: %v", err)

		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to delete category"

		if err.Error() == "category not found" {
			statusCode = http.StatusNotFound
			errorMessage = "Category not found"
		} else if err.Error() == "cannot delete category with child categories" {
			statusCode = http.StatusBadRequest
			errorMessage = "Cannot delete category with child categories"
		}

		response := contracts.ErrorResponse(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := contracts.ResponseDTO[any]{
		Success: true,
		Message: "Category deleted successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// ListCategories handles listing all categories
func (h *CategoryHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.categoryUseCase.ListCategories()
	if err != nil {
		h.logger.Error("Failed to list categories: %v", err)
		response := contracts.ErrorResponse("Failed to retrieve categories")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := contracts.CreateCategoryListResponse(categories, len(categories), 1, len(categories))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// GetChildCategories handles retrieving child categories of a parent category
func (h *CategoryHandler) GetChildCategories(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	parentID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid parent category ID: %v", err)
		response := contracts.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	categories, err := h.categoryUseCase.GetChildCategories(uint(parentID))
	if err != nil {
		h.logger.Error("Failed to get child categories: %v", err)

		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to retrieve child categories"

		if err.Error() == "parent category not found" {
			statusCode = http.StatusNotFound
			errorMessage = "Parent category not found"
		}

		response := contracts.ErrorResponse(errorMessage)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := contracts.CreateCategoryListResponse(categories, len(categories), 1, len(categories))

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}
