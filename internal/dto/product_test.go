package dto

import (
	"testing"
	"time"
)

func TestProductDTO(t *testing.T) {
	now := time.Now()
	variants := []VariantDTO{
		{
			ID:        1,
			ProductID: 1,
			SKU:       "VAR-001",
			Price:     29.99,
			Currency:  "USD",
			Stock:     50,
			IsDefault: true,
			CreatedAt: now,
			UpdatedAt: now,
		},
	}

	product := ProductDTO{
		ID:          1,
		Name:        "Test Product",
		Description: "A test product description",
		SKU:         "PROD-001",
		Price:       99.99,
		Currency:    "USD",
		Stock:       100,
		Weight:      2.5,
		CategoryID:  5,
		CreatedAt:   now,
		UpdatedAt:   now,
		Images:      []string{"image1.jpg", "image2.jpg"},
		HasVariants: true,
		Variants:    variants,
		Active:      true,
	}

	if product.ID != 1 {
		t.Errorf("Expected ID 1, got %d", product.ID)
	}
	if product.Name != "Test Product" {
		t.Errorf("Expected Name 'Test Product', got %s", product.Name)
	}
	if product.Description != "A test product description" {
		t.Errorf("Expected Description 'A test product description', got %s", product.Description)
	}
	if product.SKU != "PROD-001" {
		t.Errorf("Expected SKU 'PROD-001', got %s", product.SKU)
	}
	if product.Price != 99.99 {
		t.Errorf("Expected Price 99.99, got %f", product.Price)
	}
	if product.Currency != "USD" {
		t.Errorf("Expected Currency 'USD', got %s", product.Currency)
	}
	if product.Stock != 100 {
		t.Errorf("Expected Stock 100, got %d", product.Stock)
	}
	if product.Weight != 2.5 {
		t.Errorf("Expected Weight 2.5, got %f", product.Weight)
	}
	if product.CategoryID != 5 {
		t.Errorf("Expected CategoryID 5, got %d", product.CategoryID)
	}
	if !product.HasVariants {
		t.Errorf("Expected HasVariants true, got %t", product.HasVariants)
	}
	if !product.Active {
		t.Errorf("Expected Active true, got %t", product.Active)
	}
	if len(product.Images) != 2 {
		t.Errorf("Expected Images length 2, got %d", len(product.Images))
	}
	if product.Images[0] != "image1.jpg" {
		t.Errorf("Expected Images[0] 'image1.jpg', got %s", product.Images[0])
	}
	if len(product.Variants) != 1 {
		t.Errorf("Expected Variants length 1, got %d", len(product.Variants))
	}
	if product.Variants[0].SKU != "VAR-001" {
		t.Errorf("Expected Variants[0].SKU 'VAR-001', got %s", product.Variants[0].SKU)
	}
}

func TestVariantDTO(t *testing.T) {
	now := time.Now()
	attributes := []VariantAttributeDTO{
		{Name: "Color", Value: "Red"},
		{Name: "Size", Value: "Large"},
	}

	variant := VariantDTO{
		ID:         1,
		ProductID:  1,
		SKU:        "VAR-001",
		Price:      29.99,
		Currency:   "USD",
		Stock:      50,
		Attributes: attributes,
		Images:     []string{"variant1.jpg"},
		IsDefault:  true,
		CreatedAt:  now,
		UpdatedAt:  now,
	}

	if variant.ID != 1 {
		t.Errorf("Expected ID 1, got %d", variant.ID)
	}
	if variant.ProductID != 1 {
		t.Errorf("Expected ProductID 1, got %d", variant.ProductID)
	}
	if variant.SKU != "VAR-001" {
		t.Errorf("Expected SKU 'VAR-001', got %s", variant.SKU)
	}
	if variant.Price != 29.99 {
		t.Errorf("Expected Price 29.99, got %f", variant.Price)
	}
	if variant.Currency != "USD" {
		t.Errorf("Expected Currency 'USD', got %s", variant.Currency)
	}
	if variant.Stock != 50 {
		t.Errorf("Expected Stock 50, got %d", variant.Stock)
	}
	if !variant.IsDefault {
		t.Errorf("Expected IsDefault true, got %t", variant.IsDefault)
	}
	if len(variant.Attributes) != 2 {
		t.Errorf("Expected Attributes length 2, got %d", len(variant.Attributes))
	}
	if variant.Attributes[0].Name != "Color" {
		t.Errorf("Expected Attributes[0].Name 'Color', got %s", variant.Attributes[0].Name)
	}
	if variant.Attributes[0].Value != "Red" {
		t.Errorf("Expected Attributes[0].Value 'Red', got %s", variant.Attributes[0].Value)
	}
	if len(variant.Images) != 1 {
		t.Errorf("Expected Images length 1, got %d", len(variant.Images))
	}
	if variant.Images[0] != "variant1.jpg" {
		t.Errorf("Expected Images[0] 'variant1.jpg', got %s", variant.Images[0])
	}
}

func TestVariantAttributeDTO(t *testing.T) {
	attribute := VariantAttributeDTO{
		Name:  "Color",
		Value: "Blue",
	}

	if attribute.Name != "Color" {
		t.Errorf("Expected Name 'Color', got %s", attribute.Name)
	}
	if attribute.Value != "Blue" {
		t.Errorf("Expected Value 'Blue', got %s", attribute.Value)
	}
}

func TestCreateProductRequest(t *testing.T) {
	variants := []CreateVariantRequest{
		{
			SKU:   "VAR-001",
			Price: 19.99,
			Stock: 30,
			Attributes: []VariantAttributeDTO{
				{Name: "Size", Value: "Small"},
			},
			IsDefault: true,
		},
	}

	request := CreateProductRequest{
		Name:        "New Product",
		Description: "New product description",
		Price:       49.99,
		Stock:       75,
		Weight:      1.5,
		CategoryID:  3,
		Images:      []string{"new1.jpg", "new2.jpg"},
		Variants:    variants,
	}

	if request.Name != "New Product" {
		t.Errorf("Expected Name 'New Product', got %s", request.Name)
	}
	if request.Description != "New product description" {
		t.Errorf("Expected Description 'New product description', got %s", request.Description)
	}
	if request.Price != 49.99 {
		t.Errorf("Expected Price 49.99, got %f", request.Price)
	}
	if request.Stock != 75 {
		t.Errorf("Expected Stock 75, got %d", request.Stock)
	}
	if request.Weight != 1.5 {
		t.Errorf("Expected Weight 1.5, got %f", request.Weight)
	}
	if request.CategoryID != 3 {
		t.Errorf("Expected CategoryID 3, got %d", request.CategoryID)
	}
	if len(request.Images) != 2 {
		t.Errorf("Expected Images length 2, got %d", len(request.Images))
	}
	if len(request.Variants) != 1 {
		t.Errorf("Expected Variants length 1, got %d", len(request.Variants))
	}
	if request.Variants[0].SKU != "VAR-001" {
		t.Errorf("Expected Variants[0].SKU 'VAR-001', got %s", request.Variants[0].SKU)
	}
}

func TestCreateVariantRequest(t *testing.T) {
	attributes := []VariantAttributeDTO{
		{Name: "Color", Value: "Green"},
		{Name: "Size", Value: "Medium"},
	}

	request := CreateVariantRequest{
		SKU:        "VAR-002",
		Price:      24.99,
		Stock:      40,
		Attributes: attributes,
		Images:     []string{"variant2.jpg"},
		IsDefault:  false,
	}

	if request.SKU != "VAR-002" {
		t.Errorf("Expected SKU 'VAR-002', got %s", request.SKU)
	}
	if request.Price != 24.99 {
		t.Errorf("Expected Price 24.99, got %f", request.Price)
	}
	if request.Stock != 40 {
		t.Errorf("Expected Stock 40, got %d", request.Stock)
	}
	if request.IsDefault {
		t.Errorf("Expected IsDefault false, got %t", request.IsDefault)
	}
	if len(request.Attributes) != 2 {
		t.Errorf("Expected Attributes length 2, got %d", len(request.Attributes))
	}
	if len(request.Images) != 1 {
		t.Errorf("Expected Images length 1, got %d", len(request.Images))
	}
}

func TestUpdateProductRequest(t *testing.T) {
	categoryID := uint(7)

	request := UpdateProductRequest{
		Name:        "Updated Product",
		Description: "Updated description",
		CategoryID:  &categoryID,
		Images:      []string{"updated1.jpg"},
		Active:      true,
	}

	if request.Name != "Updated Product" {
		t.Errorf("Expected Name 'Updated Product', got %s", request.Name)
	}
	if request.Description != "Updated description" {
		t.Errorf("Expected Description 'Updated description', got %s", request.Description)
	}

	if request.CategoryID == nil || *request.CategoryID != 7 {
		t.Errorf("Expected CategoryID 7, got %v", request.CategoryID)
	}
	if !request.Active {
		t.Errorf("Expected Active true, got %t", request.Active)
	}
	if len(request.Images) != 1 {
		t.Errorf("Expected Images length 1, got %d", len(request.Images))
	}
}

func TestUpdateProductRequestWithNilValues(t *testing.T) {
	request := UpdateProductRequest{
		Name:        "Only Name Updated",
		Description: "Only Description Updated",
		Active:      false,
	}

	if request.Name != "Only Name Updated" {
		t.Errorf("Expected Name 'Only Name Updated', got %s", request.Name)
	}
	if request.Description != "Only Description Updated" {
		t.Errorf("Expected Description 'Only Description Updated', got %s", request.Description)
	}
	if request.CategoryID != nil {
		t.Errorf("Expected CategoryID nil, got %v", request.CategoryID)
	}
	if request.Active {
		t.Errorf("Expected Active false, got %t", request.Active)
	}
}

func TestProductListResponse(t *testing.T) {
	products := []ProductDTO{
		{
			ID:       1,
			Name:     "Product 1",
			Price:    19.99,
			Currency: "USD",
			Active:   true,
		},
		{
			ID:       2,
			Name:     "Product 2",
			Price:    29.99,
			Currency: "USD",
			Active:   false,
		},
	}

	pagination := PaginationDTO{
		Page:     1,
		PageSize: 10,
		Total:    2,
	}

	response := ProductListResponse{
		ListResponseDTO: ListResponseDTO[ProductDTO]{
			Success:    true,
			Message:    "Products retrieved successfully",
			Data:       products,
			Pagination: pagination,
		},
	}

	if !response.Success {
		t.Errorf("Expected Success true, got %t", response.Success)
	}
	if len(response.Data) != 2 {
		t.Errorf("Expected Data length 2, got %d", len(response.Data))
	}
	if response.Data[0].Name != "Product 1" {
		t.Errorf("Expected Data[0].Name 'Product 1', got %s", response.Data[0].Name)
	}
	if response.Data[1].Active {
		t.Errorf("Expected Data[1].Active false, got %t", response.Data[1].Active)
	}
	if response.Pagination.Total != 2 {
		t.Errorf("Expected Pagination.Total 2, got %d", response.Pagination.Total)
	}
}
