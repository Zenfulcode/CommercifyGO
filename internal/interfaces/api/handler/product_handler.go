package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/config"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	errors "github.com/zenfulcode/commercify/internal/domain/error"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/interfaces/api/contracts"
	"github.com/zenfulcode/commercify/internal/interfaces/api/middleware"
)

// ProductHandler handles product-related HTTP requests
type ProductHandler struct {
	productUseCase *usecase.ProductUseCase
	logger         logger.Logger
	config         *config.Config
}

// NewProductHandler creates a new ProductHandler
func NewProductHandler(productUseCase *usecase.ProductUseCase, logger logger.Logger, config *config.Config) *ProductHandler {
	return &ProductHandler{
		productUseCase: productUseCase,
		logger:         logger,
		config:         config,
	}
}

// handleError processes errors and returns appropriate HTTP responses
func (h *ProductHandler) handleError(w http.ResponseWriter, err error, operation string) {
	h.logger.Error("Failed to %s: %v", operation, err)

	statusCode := http.StatusInternalServerError
	errorMessage := "Failed to " + operation

	// Handle specific error types
	switch {
	case err.Error() == errors.ProductNotFoundError:
		statusCode = http.StatusNotFound
		errorMessage = err.Error()
	case strings.Contains(err.Error(), "unauthorized") || strings.Contains(err.Error(), "not authorized"):
		statusCode = http.StatusForbidden
		errorMessage = "Not authorized to perform this operation"
	case strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists"):
		statusCode = http.StatusConflict
		if strings.Contains(err.Error(), "variant") {
			errorMessage = "Variant with this SKU already exists"
		} else {
			errorMessage = "Product with this SKU already exists"
		}
	case strings.Contains(err.Error(), "category") && strings.Contains(err.Error(), "not found"):
		statusCode = http.StatusBadRequest
		errorMessage = "Category not found"
	case strings.Contains(err.Error(), "variant") && strings.Contains(err.Error(), "not found"):
		statusCode = http.StatusNotFound
		errorMessage = "Variant not found"
	case strings.Contains(err.Error(), "last variant") || (strings.Contains(err.Error(), "cannot delete") && strings.Contains(err.Error(), "variant")):
		statusCode = http.StatusConflict
		errorMessage = "Cannot delete the last variant of a product"
	case strings.Contains(err.Error(), "has orders") || strings.Contains(err.Error(), "cannot delete"):
		statusCode = http.StatusConflict
		if strings.Contains(err.Error(), "variant") {
			errorMessage = "Cannot delete variant with existing orders"
		} else {
			errorMessage = "Cannot delete product with existing orders"
		}
	case strings.Contains(err.Error(), "currency"):
		statusCode = http.StatusBadRequest
		errorMessage = "Invalid currency code"
	case strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "validation"):
		statusCode = http.StatusBadRequest
		if strings.Contains(err.Error(), "variant") {
			errorMessage = "Invalid variant data"
		} else if strings.Contains(err.Error(), "search") || strings.Contains(err.Error(), "parameters") {
			errorMessage = "Invalid search parameters"
		} else {
			errorMessage = "Invalid product data"
		}
	}

	h.writeErrorResponse(w, statusCode, errorMessage)
}

// handleValidationError handles request validation errors
func (h *ProductHandler) handleValidationError(w http.ResponseWriter, err error, context string) {
	h.logger.Error("Validation error in %s: %v", context, err)
	h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
}

// handleAuthorizationError handles authorization errors
func (h *ProductHandler) handleAuthorizationError(w http.ResponseWriter, context string) {
	h.logger.Error("Unauthorized access attempt in %s - admin required", context)
	h.writeErrorResponse(w, http.StatusForbidden, "Unauthorized - admin access required")
}

// handleIDParsingError handles URL parameter parsing errors
func (h *ProductHandler) handleIDParsingError(w http.ResponseWriter, err error, idType, context string) {
	h.logger.Error("Invalid %s ID in %s: %v", idType, context, err)
	h.writeErrorResponse(w, http.StatusBadRequest, "Invalid "+idType+" ID")
}

// writeErrorResponse is a helper to write error responses consistently
func (h *ProductHandler) writeErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	response := contracts.ErrorResponse(message)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)
}

// checkAdminAuthorization checks if the user has admin role
func (h *ProductHandler) checkAdminAuthorization(r *http.Request) bool {
	role, ok := r.Context().Value(middleware.RoleKey).(string)
	return ok && role == string(entity.RoleAdmin)
}

// --- Handlers --- //

// CreateProduct handles product creation (admin only)
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Check admin authorization
	if !h.checkAdminAuthorization(r) {
		h.handleAuthorizationError(w, "CreateProduct")
		return
	}

	// Parse request body
	var request contracts.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.handleValidationError(w, err, "CreateProduct")
		return
	}

	h.logger.Info("Creating product:", request)

	input := request.ToUseCaseInput()

	// Create product
	product, err := h.productUseCase.CreateProduct(input)
	if err != nil {
		h.handleError(w, err, "create product")
		return
	}

	// Convert to DTO
	response := contracts.SuccessResponseWithMessage(product.ToProductDTO(), "Product created successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// GetProduct handles getting a product by ID
func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	// Get product ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		h.handleIDParsingError(w, err, "product", "GetProduct")
		return
	}

	// Get product - no currency filtering needed since each product has its own currency
	product, err := h.productUseCase.GetProductByID(uint(id))
	if err != nil {
		h.handleError(w, err, "retrieve product")
		return
	}

	// Convert to DTO
	response := contracts.SuccessResponse(product.ToProductDTO())

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateProduct handles updating a product (admin only)
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Check admin authorization
	if !h.checkAdminAuthorization(r) {
		h.handleAuthorizationError(w, "UpdateProduct")
		return
	}

	// Get product ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		h.handleIDParsingError(w, err, "product", "UpdateProduct")
		return
	}

	// Parse request body
	var request contracts.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.handleValidationError(w, err, "UpdateProduct")
		return
	}

	// Convert DTO to usecase input
	input := request.ToUseCaseInput()

	// Update product
	product, err := h.productUseCase.UpdateProduct(uint(id), input)
	if err != nil {
		h.handleError(w, err, "update product")
		return
	}

	// Convert to DTO
	response := contracts.SuccessResponseWithMessage(product.ToProductDTO(), "Product updated successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteProduct handles deleting a product (admin only)
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Check admin authorization
	if !h.checkAdminAuthorization(r) {
		h.handleAuthorizationError(w, "DeleteProduct")
		return
	}

	// Get product ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		h.handleIDParsingError(w, err, "product", "DeleteProduct")
		return
	}

	// Delete product
	err = h.productUseCase.DeleteProduct(uint(id))
	if err != nil {
		h.handleError(w, err, "delete product")
		return
	}

	response := contracts.SuccessResponseMessage("Product deleted successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListProducts handles listing all products (admin only)
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Check admin authorization
	if !h.checkAdminAuthorization(r) {
		h.handleAuthorizationError(w, "ListProducts")
		return
	}

	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1 // Default page
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}

	// Parse optional parameters
	var query *string
	if queryStr := r.URL.Query().Get("query"); queryStr != "" {
		query = &queryStr
	}

	var categoryID *uint
	if catIDStr := r.URL.Query().Get("category_id"); catIDStr != "" {
		if catID, err := strconv.ParseUint(catIDStr, 10, 32); err == nil {
			catIDUint := uint(catID)
			categoryID = &catIDUint
		}
	}

	var minPrice *float64
	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		if minPriceVal, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			minPrice = &minPriceVal
		}
	}

	var maxPrice *float64
	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if maxPriceVal, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			maxPrice = &maxPriceVal
		}
	}

	var currencyCode string
	if currencyCodeStr := r.URL.Query().Get("currency"); currencyCodeStr != "" {
		currencyCode = currencyCodeStr
	}

	// Parse active parameter - defaults to true for admin (show active products)
	activeOnly := true // Default to showing active products for admin
	if activeStr := r.URL.Query().Get("active"); activeStr != "" {
		switch activeStr {
		case "false", "0":
			activeOnly = false
		case "true", "1":
			activeOnly = true
		}
		// If the query parameter is "all", we want to show all products regardless of status
		if activeStr == "all" {
			// We'll handle this case in the repository by modifying the logic
			activeOnly = false // For now, this will need repository changes
		}
	}

	offset := (page - 1) * pageSize

	// Convert to usecase input
	input := usecase.SearchProductsInput{
		Offset:       uint(offset),
		Limit:        uint(pageSize),
		CurrencyCode: currencyCode,
		ActiveOnly:   activeOnly,
	}

	// Handle optional fields
	if query != nil {
		input.Query = *query
	}
	if categoryID != nil {
		input.CategoryID = *categoryID
	}
	if minPrice != nil {
		input.MinPrice = *minPrice
	}
	if maxPrice != nil {
		input.MaxPrice = *maxPrice
	}

	products, total, err := h.productUseCase.ListProducts(input)
	if err != nil {
		h.handleError(w, err, "search products")
		return
	}

	response := contracts.CreateProductListResponse(products, total, page, pageSize)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// SearchProducts handles searching products
func (h *ProductHandler) SearchProducts(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1 // Default page
	}

	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}

	// Parse optional parameters
	var query *string
	if queryStr := r.URL.Query().Get("query"); queryStr != "" {
		query = &queryStr
	}

	var categoryID *uint
	if catIDStr := r.URL.Query().Get("category_id"); catIDStr != "" {
		if catID, err := strconv.ParseUint(catIDStr, 10, 32); err == nil {
			catIDUint := uint(catID)
			categoryID = &catIDUint
		}
	}

	var minPrice *float64
	if minPriceStr := r.URL.Query().Get("min_price"); minPriceStr != "" {
		if minPriceVal, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
			minPrice = &minPriceVal
		}
	}

	var maxPrice *float64
	if maxPriceStr := r.URL.Query().Get("max_price"); maxPriceStr != "" {
		if maxPriceVal, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
			maxPrice = &maxPriceVal
		}
	}

	var currencyCode string
	if currencyCodeStr := r.URL.Query().Get("currency"); currencyCodeStr != "" {
		currencyCode = currencyCodeStr
	}

	offset := (page - 1) * pageSize

	// Convert to usecase input
	input := usecase.SearchProductsInput{
		Offset:       uint(offset),
		Limit:        uint(pageSize),
		CurrencyCode: currencyCode,
		ActiveOnly:   true, // Only active products by default
	}

	// Handle optional fields
	if query != nil {
		input.Query = *query
	}
	if categoryID != nil {
		input.CategoryID = *categoryID
	}
	if minPrice != nil {
		input.MinPrice = *minPrice
	}
	if maxPrice != nil {
		input.MaxPrice = *maxPrice
	}

	products, total, err := h.productUseCase.ListProducts(input)
	if err != nil {
		h.handleError(w, err, "search products")
		return
	}

	// Convert to DTOs
	response := contracts.CreateProductListResponse(products, total, page, pageSize)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListCategories handles listing all product categories
func (h *ProductHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.productUseCase.ListCategories()
	if err != nil {
		h.logger.Error("Failed to list categories: %v", err)
		response := contracts.ErrorResponse("Failed to list categories")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := contracts.SuccessResponseWithMessage(categories, "Categories retrieved successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AddVariant handles adding a new variant to a product (admin only)
func (h *ProductHandler) AddVariant(w http.ResponseWriter, r *http.Request) {
	// Check admin authorization
	if !h.checkAdminAuthorization(r) {
		h.handleAuthorizationError(w, "AddVariant")
		return
	}

	// Parse request body
	var request contracts.CreateVariantRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.handleValidationError(w, err, "AddVariant")
		return
	}

	// Get product ID from URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		h.handleIDParsingError(w, err, "product", "AddVariant")
		return
	}

	// Convert DTO to usecase input
	input := request.ToUseCaseInput()

	// Add variant
	variant, err := h.productUseCase.AddVariant(uint(productID), input)
	if err != nil {
		h.handleError(w, err, "add variant")
		return
	}

	// Convert to DTO
	response := contracts.SuccessResponseWithMessage(variant.ToVariantDTO(), "Variant added successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateVariant handles updating a product variant (admin only)
func (h *ProductHandler) UpdateVariant(w http.ResponseWriter, r *http.Request) {
	// Check admin authorization
	if !h.checkAdminAuthorization(r) {
		h.handleAuthorizationError(w, "UpdateVariant")
		return
	}

	// Get IDs from URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		h.handleIDParsingError(w, err, "product", "UpdateVariant")
		return
	}

	variantID, err := strconv.ParseUint(vars["variantId"], 10, 32)
	if err != nil {
		h.handleIDParsingError(w, err, "variant", "UpdateVariant")
		return
	}

	// Parse request body
	var request contracts.UpdateVariantRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.handleValidationError(w, err, "UpdateVariant")
		return
	}

	// Convert DTO to usecase input
	input := request.ToUseCaseInput()

	// Update variant
	variant, err := h.productUseCase.UpdateVariant(uint(productID), uint(variantID), input)
	if err != nil {
		h.handleError(w, err, "update variant")
		return
	}

	response := contracts.SuccessResponseWithMessage(variant.ToVariantDTO(), "Variant updated successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteVariant handles deleting a product variant (admin only)
func (h *ProductHandler) DeleteVariant(w http.ResponseWriter, r *http.Request) {
	// Check admin authorization
	if !h.checkAdminAuthorization(r) {
		h.handleAuthorizationError(w, "DeleteVariant")
		return
	}

	// Get IDs from URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		h.handleIDParsingError(w, err, "product", "DeleteVariant")
		return
	}

	variantID, err := strconv.ParseUint(vars["variantId"], 10, 32)
	if err != nil {
		h.handleIDParsingError(w, err, "variant", "DeleteVariant")
		return
	}

	// Delete variant
	err = h.productUseCase.DeleteVariant(uint(productID), uint(variantID))
	if err != nil {
		h.handleError(w, err, "delete variant")
		return
	}

	response := contracts.SuccessResponseMessage("Variant deleted successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
