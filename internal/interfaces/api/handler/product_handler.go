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
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/dto"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
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

// --- Helper Functions --- //

func toVariantDTO(variant *entity.ProductVariant) dto.VariantDTO {
	if variant == nil {
		return dto.VariantDTO{}
	}

	attributesDTO := make([]dto.VariantAttributeDTO, len(variant.Attributes))
	for i, a := range variant.Attributes {
		attributesDTO[i] = dto.VariantAttributeDTO{
			Name:  a.Name,
			Value: a.Value,
		}
	}

	return dto.VariantDTO{
		ID:         variant.ID,
		ProductID:  variant.ProductID,
		SKU:        variant.SKU,
		Price:      money.FromCents(variant.Price),
		Currency:   variant.CurrencyCode,
		Stock:      variant.Stock,
		Attributes: attributesDTO,
		Images:     variant.Images,
		IsDefault:  variant.IsDefault,
		CreatedAt:  variant.CreatedAt,
		UpdatedAt:  variant.UpdatedAt,
	}
}

func toProductDTO(product *entity.Product) dto.ProductDTO {
	if product == nil {
		return dto.ProductDTO{}
	}
	variantsDTO := make([]dto.VariantDTO, len(product.Variants))
	for i, v := range product.Variants {
		variantsDTO[i] = toVariantDTO(v)
	}

	return dto.ProductDTO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		SKU:         product.ProductNumber,
		Price:       money.FromCents(product.Price),
		Currency:    product.CurrencyCode,
		Stock:       product.Stock,
		Weight:      product.Weight,
		CategoryID:  product.CategoryID,
		Images:      product.Images,
		HasVariants: product.HasVariants,
		Variants:    variantsDTO,
		CreatedAt:   product.CreatedAt,
		UpdatedAt:   product.UpdatedAt,
		Active:      product.Active,
	}
}

// --- Handlers --- //

// CreateProduct handles product creation
func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)

	if !ok || userID == 0 {
		h.logger.Error("Unauthorized access attempt in CreateProduct")
		response := dto.ErrorResponse("Unauthorized")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse request body
	var request dto.CreateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Invalid request body in CreateProduct: %v", err)
		response := dto.ErrorResponse("Invalid request body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	variantInputs := make([]usecase.CreateVariantInput, len(request.Variants))
	for i, v := range request.Variants {
		attributes := make([]entity.VariantAttribute, len(v.Attributes))
		for j, a := range v.Attributes {
			attributes[j] = entity.VariantAttribute{
				Name:  a.Name,
				Value: a.Value,
			}
		}

		variantInputs[i] = usecase.CreateVariantInput{
			SKU:        v.SKU,
			Price:      v.Price,
			Stock:      v.Stock,
			Attributes: attributes,
			Images:     v.Images,
			IsDefault:  v.IsDefault,
		}
	}

	// Convert DTO to usecase input
	input := usecase.CreateProductInput{
		Name:        request.Name,
		Description: request.Description,
		Price:       request.Price,
		Stock:       request.Stock,
		Weight:      request.Weight,
		CategoryID:  request.CategoryID,
		Images:      request.Images,
		Variants:    variantInputs,
	}

	// Create product
	product, err := h.productUseCase.CreateProduct(input)
	if err != nil {
		h.logger.Error("Failed to create product: %v", err)

		// Handle specific error cases
		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to create product"

		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
			statusCode = http.StatusConflict
			errorMessage = "Product with this SKU already exists"
		} else if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "validation") {
			statusCode = http.StatusBadRequest
			errorMessage = "Invalid product data"
		} else if strings.Contains(err.Error(), "category") && strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusBadRequest
			errorMessage = "Category not found"
		} else if strings.Contains(err.Error(), "unauthorized") {
			statusCode = http.StatusForbidden
		}

		response := dto.ErrorResponse(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to DTO
	productDTO := toProductDTO(product)

	response := dto.SuccessResponseWithMessage(productDTO, "Product created successfully")

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
		h.logger.Error("Invalid product ID in GetProduct: %v", err)
		response := dto.ErrorResponse("Invalid product ID")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get product
	currencyCode := &h.config.DefaultCurrency
	if currencyCodeStr := r.URL.Query().Get("currency"); currencyCodeStr != "" {
		currencyCode = &currencyCodeStr
	}

	var product *entity.Product
	product, err = h.productUseCase.GetProductByID(uint(id), *currencyCode)

	if err != nil {
		h.logger.Error("Failed to get product: %v", err)

		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to retrieve product"

		if err.Error() == errors.ProductNotFoundError {
			statusCode = http.StatusNotFound
			errorMessage = "Product not found"
		} else if strings.Contains(err.Error(), "currency") {
			statusCode = http.StatusBadRequest
			errorMessage = "Invalid currency code"
		}

		response := dto.ErrorResponse(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to DTO
	productDTO := toProductDTO(product)

	response := dto.SuccessResponse(productDTO)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// UpdateProduct handles updating a product
func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		h.logger.Error("Unauthorized access attempt in UpdateProduct")
		response := dto.ErrorResponse("Unauthorized")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get product ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid product ID in UpdateProduct: %v", err)
		response := dto.ErrorResponse("Invalid product ID")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse request body
	var request dto.UpdateProductRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Invalid request body in UpdateProduct: %v", err)
		response := dto.ErrorResponse("Invalid request body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert DTO to usecase input
	input := usecase.UpdateProductInput{
		Name:        request.Name,
		Description: request.Description,
		Images:      request.Images,
		Active:      request.Active,
	}

	if request.Weight != nil {
		input.Weight = *request.Weight
	}
	if request.CategoryID != nil {
		input.CategoryID = *request.CategoryID
	}

	// Update product
	product, err := h.productUseCase.UpdateProduct(uint(id), input)
	if err != nil {
		h.logger.Error("Failed to update product: %v", err)

		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to update product"

		if err.Error() == "unauthorized: not the seller of this product" {
			statusCode = http.StatusForbidden
			errorMessage = "Not authorized to update this product"
		} else if err.Error() == errors.ProductNotFoundError {
			statusCode = http.StatusNotFound
			errorMessage = "Product not found"
		} else if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
			statusCode = http.StatusConflict
			errorMessage = "Product with this SKU already exists"
		} else if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "validation") {
			statusCode = http.StatusBadRequest
			errorMessage = "Invalid product data"
		} else if strings.Contains(err.Error(), "category") && strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusBadRequest
			errorMessage = "Category not found"
		}

		response := dto.ErrorResponse(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to DTO
	productDTO := toProductDTO(product)

	response := dto.SuccessResponseWithMessage(productDTO, "Product updated successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteProduct handles deleting a product
func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		h.logger.Error("Unauthorized access attempt in DeleteProduct")
		response := dto.ErrorResponse("Unauthorized")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get product ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid product ID in DeleteProduct: %v", err)
		response := dto.ErrorResponse("Invalid product ID")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Delete product
	err = h.productUseCase.DeleteProduct(uint(id))
	if err != nil {
		h.logger.Error("Failed to delete product: %v", err)

		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to delete product"

		if err.Error() == "unauthorized: not the seller of this product" {
			statusCode = http.StatusForbidden
			errorMessage = "Not authorized to delete this product"
		} else if err.Error() == errors.ProductNotFoundError {
			statusCode = http.StatusNotFound
			errorMessage = "Product not found"
		} else if strings.Contains(err.Error(), "has orders") || strings.Contains(err.Error(), "cannot delete") {
			statusCode = http.StatusConflict
			errorMessage = "Cannot delete product with existing orders"
		}

		response := dto.ErrorResponse(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.SuccessResponseMessage("Product deleted successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListProducts handles listing all products
func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)

	if !ok || userID == 0 {
		h.logger.Error("Unauthorized access attempt in ListProducts")
		response := dto.ErrorResponse("Unauthorized")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse pagination parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page <= 0 {
		page = 1 // Default page
	}
	pageSize, _ := strconv.Atoi(r.URL.Query().Get("page_size"))
	if pageSize <= 0 {
		pageSize = 10 // Default page size
	}

	offset := (page - 1) * pageSize
	products, total, err := h.productUseCase.ListProducts(offset, pageSize)

	if err != nil {
		h.logger.Error("Failed to list products: %v", err)
		response := dto.ErrorResponse("Failed to list products")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to DTOs
	productDTOs := make([]dto.ProductDTO, len(products))
	for i, product := range products {
		productDTOs[i] = toProductDTO(product)
	}

	response := dto.ProductListResponse{
		ListResponseDTO: dto.ListResponseDTO[dto.ProductDTO]{
			Success: true,
			Data:    productDTOs,
			Message: "Products retrieved successfully",
			Pagination: dto.PaginationDTO{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
		},
	}

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
		Offset:       offset,
		Limit:        pageSize,
		CurrencyCode: currencyCode,
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

	products, total, err := h.productUseCase.SearchProducts(input)
	if err != nil {
		h.logger.Error("Failed to search products: %v", err)

		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to search products"

		if strings.Contains(err.Error(), "currency") {
			statusCode = http.StatusBadRequest
			errorMessage = "Invalid currency code"
		} else if strings.Contains(err.Error(), "category") && strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusBadRequest
			errorMessage = "Category not found"
		} else if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "validation") {
			statusCode = http.StatusBadRequest
			errorMessage = "Invalid search parameters"
		}

		response := dto.ErrorResponse(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to DTOs
	productDTOs := make([]dto.ProductDTO, len(products))
	for i, product := range products {
		productDTOs[i] = toProductDTO(product)
	}

	response := dto.ProductListResponse{
		ListResponseDTO: dto.ListResponseDTO[dto.ProductDTO]{
			Success: true,
			Data:    productDTOs,
			Message: "Products search completed successfully",
			Pagination: dto.PaginationDTO{
				Page:     page,
				PageSize: pageSize,
				Total:    total,
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListCategories handles listing all product categories
func (h *ProductHandler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.productUseCase.ListCategories()
	if err != nil {
		h.logger.Error("Failed to list categories: %v", err)
		response := dto.ErrorResponse("Failed to list categories")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.SuccessResponseWithMessage(categories, "Categories retrieved successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// AddVariant handles adding a new variant to a product
func (h *ProductHandler) AddVariant(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		h.logger.Error("Unauthorized access attempt in AddVariant")
		response := dto.ErrorResponse("Unauthorized")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse request body
	var request dto.CreateVariantRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Invalid request body in AddVariant: %v", err)
		response := dto.ErrorResponse("Invalid request body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get product ID from URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["productId"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid product ID in AddVariant: %v", err)
		response := dto.ErrorResponse("Invalid product ID")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	attributesDTO := make([]entity.VariantAttribute, len(request.Attributes))
	for i, a := range request.Attributes {
		attributesDTO[i] = entity.VariantAttribute{
			Name:  a.Name,
			Value: a.Value,
		}

	}

	// Convert DTO to usecase input
	input := usecase.AddVariantInput{
		ProductID:  uint(productID),
		SKU:        request.SKU,
		Price:      request.Price,
		Stock:      request.Stock,
		Attributes: attributesDTO,
		Images:     request.Images,
		IsDefault:  request.IsDefault,
	}

	// Add variant
	variant, err := h.productUseCase.AddVariant(input)
	if err != nil {
		h.logger.Error("Failed to add variant: %v", err)

		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to add variant"

		if err.Error() == "unauthorized: not the seller of this product" {
			statusCode = http.StatusForbidden
			errorMessage = "Not authorized to add variant to this product"
		} else if err.Error() == errors.ProductNotFoundError {
			statusCode = http.StatusNotFound
			errorMessage = "Product not found"
		} else if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
			statusCode = http.StatusConflict
			errorMessage = "Variant with this SKU already exists"
		} else if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "validation") {
			statusCode = http.StatusBadRequest
			errorMessage = "Invalid variant data"
		}

		response := dto.ErrorResponse(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to DTO
	variantDTO := toVariantDTO(variant)

	response := dto.SuccessResponseWithMessage(variantDTO, "Variant added successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateVariant handles updating a product variant
func (h *ProductHandler) UpdateVariant(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		h.logger.Error("Unauthorized access attempt in UpdateVariant")
		response := dto.ErrorResponse("Unauthorized")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get IDs from URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid product ID in UpdateVariant: %v", err)
		response := dto.ErrorResponse("Invalid product ID")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	variantID, err := strconv.ParseUint(vars["variant_id"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid variant ID in UpdateVariant: %v", err)
		response := dto.ErrorResponse("Invalid variant ID")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Parse request body
	var request dto.CreateVariantRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Invalid request body in UpdateVariant: %v", err)
		response := dto.ErrorResponse("Invalid request body")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	attributesDTO := make([]entity.VariantAttribute, len(request.Attributes))
	for i, a := range request.Attributes {
		attributesDTO[i] = entity.VariantAttribute{
			Name:  a.Name,
			Value: a.Value,
		}
	}

	// Convert DTO to usecase input
	input := usecase.UpdateVariantInput{
		SKU:        request.SKU,
		Price:      request.Price,
		Stock:      request.Stock,
		Attributes: attributesDTO,
		Images:     request.Images,
		IsDefault:  request.IsDefault,
	}

	// Update variant
	variant, err := h.productUseCase.UpdateVariant(uint(productID), uint(variantID), input)
	if err != nil {
		h.logger.Error("Failed to update variant: %v", err)

		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to update variant"

		if err.Error() == "unauthorized: not the seller of this product" {
			statusCode = http.StatusForbidden
			errorMessage = "Not authorized to update this variant"
		} else if err.Error() == errors.ProductNotFoundError {
			statusCode = http.StatusNotFound
			errorMessage = "Product not found"
		} else if strings.Contains(err.Error(), "variant") && strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
			errorMessage = "Variant not found"
		} else if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
			statusCode = http.StatusConflict
			errorMessage = "Variant with this SKU already exists"
		} else if strings.Contains(err.Error(), "invalid") || strings.Contains(err.Error(), "validation") {
			statusCode = http.StatusBadRequest
			errorMessage = "Invalid variant data"
		}

		response := dto.ErrorResponse(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to DTO
	variantDTO := toVariantDTO(variant)

	response := dto.SuccessResponseWithMessage(variantDTO, "Variant updated successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteVariant handles deleting a product variant
func (h *ProductHandler) DeleteVariant(w http.ResponseWriter, r *http.Request) {
	// Get user ID from context
	userID, ok := r.Context().Value(middleware.UserIDKey).(uint)
	if !ok || userID == 0 {
		h.logger.Error("Unauthorized access attempt in DeleteVariant")
		response := dto.ErrorResponse("Unauthorized")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get IDs from URL
	vars := mux.Vars(r)
	productID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid product ID in DeleteVariant: %v", err)
		response := dto.ErrorResponse("Invalid product ID")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	variantID, err := strconv.ParseUint(vars["variant_id"], 10, 32)
	if err != nil {
		h.logger.Error("Invalid variant ID in DeleteVariant: %v", err)
		response := dto.ErrorResponse("Invalid variant ID")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Delete variant
	err = h.productUseCase.DeleteVariant(uint(productID), uint(variantID))

	if err != nil {
		h.logger.Error("Failed to delete variant: %v", err)

		statusCode := http.StatusInternalServerError
		errorMessage := "Failed to delete variant"

		if err.Error() == "unauthorized: not the seller of this product" {
			statusCode = http.StatusForbidden
			errorMessage = "Not authorized to delete this variant"
		} else if err.Error() == errors.ProductNotFoundError {
			statusCode = http.StatusNotFound
			errorMessage = "Product not found"
		} else if strings.Contains(err.Error(), "variant") && strings.Contains(err.Error(), "not found") {
			statusCode = http.StatusNotFound
			errorMessage = "Variant not found"
		} else if strings.Contains(err.Error(), "last variant") || strings.Contains(err.Error(), "cannot delete") {
			statusCode = http.StatusConflict
			errorMessage = "Cannot delete the last variant of a product"
		} else if strings.Contains(err.Error(), "has orders") {
			statusCode = http.StatusConflict
			errorMessage = "Cannot delete variant with existing orders"
		}

		response := dto.ErrorResponse(errorMessage)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(response)
		return
	}

	response := dto.SuccessResponseMessage("Variant deleted successfully")

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
