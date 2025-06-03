package usecase_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/testutil/mock"
)

func TestProductUseCase_CreateProduct(t *testing.T) {
	t.Run("Create simple product successfully", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()

		// Create a test category
		category := &entity.Category{
			ID:   1,
			Name: "Test Category",
		}
		categoryRepo.Create(category)

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
		)

		// Create product input
		input := usecase.CreateProductInput{
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       99.99,
			Stock:       100,
			CategoryID:  1,
			Images:      []string{"image1.jpg", "image2.jpg"},
		}

		// Execute
		product, err := productUseCase.CreateProduct(input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, input.Name, product.Name)
		assert.Equal(t, input.Description, product.Description)
		assert.Equal(t, money.ToCents(input.Price), product.Price)
		assert.Equal(t, input.Stock, product.Stock)
		assert.Equal(t, input.CategoryID, product.CategoryID)
		assert.Equal(t, input.Images, product.Images)
		assert.True(t, product.HasVariants, "Product should have variants set to true")
		assert.Len(t, product.Variants, 1, "Product should have one default variant")
		assert.Equal(t, product.ProductNumber, product.Variants[0].SKU, "Default variant SKU should match product number")
	})

	t.Run("Create product with variants successfully", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()

		// Create a test category
		category := &entity.Category{
			ID:   1,
			Name: "Test Category",
		}
		categoryRepo.Create(category)

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
		)

		// Create product input with variants
		input := usecase.CreateProductInput{
			Name:        "Test Product with Variants",
			Description: "This is a test product with variants",
			Price:       99.99,
			Stock:       100,
			CategoryID:  1,
			Images:      []string{"image1.jpg", "image2.jpg"},
			Variants: []usecase.CreateVariantInput{
				{
					SKU:        "SKU-1",
					Price:      99.99,
					Stock:      50,
					Attributes: []entity.VariantAttribute{{Name: "Color", Value: "Red"}},
					Images:     []string{"red.jpg"},
					IsDefault:  true,
				},
				{
					SKU:        "SKU-2",
					Price:      109.99,
					Stock:      50,
					Attributes: []entity.VariantAttribute{{Name: "Color", Value: "Blue"}},
					Images:     []string{"blue.jpg"},
					IsDefault:  false,
				},
			},
		}

		// Execute
		product, err := productUseCase.CreateProduct(input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, input.Name, product.Name)
		assert.Len(t, product.Variants, 2)

		// Check variants
		assert.Equal(t, "SKU-1", product.Variants[0].SKU)
		assert.Equal(t, true, product.Variants[0].IsDefault)
		assert.Equal(t, "SKU-2", product.Variants[1].SKU)
	})

	t.Run("Create product with invalid category", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()
		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
		)

		// Create product input with invalid category
		input := usecase.CreateProductInput{
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       99.99,
			Stock:       100,
			CategoryID:  999, // Non-existent category
			Images:      []string{"image1.jpg", "image2.jpg"},
		}

		// Execute
		product, err := productUseCase.CreateProduct(input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, product)
		assert.Contains(t, err.Error(), "category not found")
	})
}

func TestProductUseCase_GetProductByID(t *testing.T) {
	t.Run("Get existing product", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()

		// Create a test product
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       9999,
			Stock:       100,
			CategoryID:  1,
			Images:      []string{"image1.jpg", "image2.jpg"},
		}
		productRepo.Create(product)

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
		)

		// Execute
		result, err := productUseCase.GetProductByID(1, "USD")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, product.ID, result.ID)
		assert.Equal(t, product.Name, result.Name)
	})

	t.Run("Get non-existent product", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()
		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
		)

		// Execute with non-existent ID
		result, err := productUseCase.GetProductByID(999, "USD")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})

	t.Run("Get product by currency", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()

		// Create a test product
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       9999,
			Stock:       100,
			CategoryID:  1,
			Images:      []string{"image1.jpg", "image2.jpg"},
			HasVariants: true,
		}
		productRepo.Create(product)

		// Create a test currency
		currency := &entity.Currency{
			Code:         "USD",
			Name:         "United States Dollar",
			Symbol:       "$",
			ExchangeRate: 1.0,
			IsDefault:    true,
		}
		currencyRepo.Create(currency)

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
		)

		// Execute
		result, err := productUseCase.GetProductByID(1, "USD")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(9999), result.Price)
	})

	t.Run("Get product in different currency with no price", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()

		// Create a test product
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       9999,
			Stock:       100,
			CategoryID:  1,
			Images:      []string{"image1.jpg", "image2.jpg"},
			HasVariants: true,
			Prices: []entity.ProductPrice{
				{
					CurrencyCode: "USD",
					Price:        9999,
				},
			},
		}
		productRepo.Create(product)

		// Create a test currency
		currency := &entity.Currency{
			Code:         "EUR",
			ExchangeRate: 0.85,
			IsDefault:    false,
		}
		currencyRepo.Create(currency)

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
		)

		// Execute
		result, err := productUseCase.GetProductByID(1, "EUR")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(8499), result.Price)
	})

	t.Run("Get product with invalid currency", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()

		// Create a test product
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       9999,
			Stock:       100,
			CategoryID:  1,
			Images:      []string{"image1.jpg", "image2.jpg"},
			HasVariants: true,
		}
		productRepo.Create(product)

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
		)

		// Execute
		result, err := productUseCase.GetProductByID(1, "INVALID")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestProductUseCase_UpdateProduct(t *testing.T) {
	t.Run("Update product successfully", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()

		// Create test category and product
		category := &entity.Category{
			ID:   1,
			Name: "Test Category",
		}
		categoryRepo.Create(category)

		newCategory := &entity.Category{
			ID:   2,
			Name: "New Category",
		}
		categoryRepo.Create(newCategory)

		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       9999,
			Stock:       100,
			CategoryID:  1,
			Images:      []string{"image1.jpg", "image2.jpg"},
		}
		productRepo.Create(product)

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
		)

		// Update input
		input := usecase.UpdateProductInput{
			Name:        "Updated Product",
			Description: "Updated description",
			Price:       12999,
			Stock:       50,
			CategoryID:  2,
			Images:      []string{"updated.jpg"},
		}

		// Execute
		updatedProduct, err := productUseCase.UpdateProduct(1, input)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, input.Name, updatedProduct.Name)
		assert.Equal(t, input.Description, updatedProduct.Description)
		assert.Equal(t, money.ToCents(input.Price), updatedProduct.Price)
		assert.Equal(t, input.Stock, updatedProduct.Stock)
		assert.Equal(t, input.CategoryID, updatedProduct.CategoryID)
		assert.Equal(t, input.Images, updatedProduct.Images)
	})
}

func TestProductUseCase_AddVariant(t *testing.T) {
	t.Run("Add variant to product successfully", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()

		// Create a test product without variants
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       9999,
			Stock:       100,
			CategoryID:  1,
			Images:      []string{"image1.jpg", "image2.jpg"},
		}
		productRepo.Create(product)

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
		)

		// Add variant input
		input := usecase.AddVariantInput{
			ProductID:  1,
			SKU:        "SKU-1",
			Price:      129.99,
			Stock:      50,
			Attributes: []entity.VariantAttribute{{Name: "Color", Value: "Red"}},
			Images:     []string{"red.jpg"},
			IsDefault:  true,
		}

		// Execute
		variant, err := productUseCase.AddVariant(input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, variant)
		assert.Equal(t, input.ProductID, variant.ProductID)
		assert.Equal(t, input.SKU, variant.SKU)
		assert.Equal(t, money.ToCents(input.Price), variant.Price)
		assert.Equal(t, input.Stock, variant.Stock)
		assert.Equal(t, input.Attributes, variant.Attributes)
		assert.Equal(t, input.Images, variant.Images)
		assert.Equal(t, input.IsDefault, variant.IsDefault)

		// Check that product is updated

		updatedProduct, _ := productRepo.GetByID(1)
		updatedProductPrice, err := productUseCase.GetProductByID(1, "USD")

		assert.NoError(t, err)
		assert.True(t, updatedProduct.HasVariants)
		assert.Equal(t, money.ToCents(input.Price), updatedProductPrice.Price)
	})
}

func TestProductUseCase_UpdateVariant(t *testing.T) {
	t.Run("Update variant successfully", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()

		// Create a test product with variants
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       9999,
			Stock:       100,
			CategoryID:  1,
			Images:      []string{"image1.jpg", "image2.jpg"},
			HasVariants: true,
		}
		productRepo.Create(product)

		// Create two variants
		variant1 := &entity.ProductVariant{
			ID:        1,
			ProductID: 1,
			SKU:       "SKU-1",
			Price:     9999,
			Stock:     50,
			Attributes: []entity.VariantAttribute{
				{Name: "Color", Value: "Red"},
			},
			Images:    []string{"red.jpg"},
			IsDefault: true,
		}
		productVariantRepo.Create(variant1)

		variant2 := &entity.ProductVariant{
			ID:        2,
			ProductID: 1,
			SKU:       "SKU-2",
			Price:     10999,
			Stock:     50,
			Attributes: []entity.VariantAttribute{
				{Name: "Color", Value: "Blue"},
			},
			Images:    []string{"blue.jpg"},
			IsDefault: false,
		}
		productVariantRepo.Create(variant2)

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
		)

		// Update variant input
		input := usecase.UpdateVariantInput{
			SKU:        "SKU-2-UPDATED",
			Price:      119.99,
			Stock:      25,
			Attributes: []entity.VariantAttribute{{Name: "Color", Value: "Navy Blue"}},
			Images:     []string{"navy.jpg"},
			IsDefault:  true, // Change default variant
		}

		// Execute
		updatedVariant, err := productUseCase.UpdateVariant(1, 2, input)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, input.SKU, updatedVariant.SKU)
		assert.Equal(t, money.ToCents(input.Price), updatedVariant.Price)
		assert.Equal(t, input.Stock, updatedVariant.Stock)
		assert.Equal(t, input.Attributes, updatedVariant.Attributes)
		assert.Equal(t, input.Images, updatedVariant.Images)
		assert.Equal(t, input.IsDefault, updatedVariant.IsDefault)

		// Check that the previous default variant is no longer default
		formerDefaultVariant, _ := productVariantRepo.GetByID(1)
		assert.False(t, formerDefaultVariant.IsDefault)

		// Check that product price is updated
		// Probably shouldn't be calling GetProductByID here, but it's just for testing
		updatedProduct, _ := productUseCase.GetProductByID(1, "USD")
		assert.Equal(t, money.ToCents(input.Price), updatedProduct.Price)
	})
}

func TestProductUseCase_DeleteVariant(t *testing.T) {
	t.Run("Delete variant successfully", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()

		// Create a test product with variants
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       9999,
			Stock:       100,
			CategoryID:  1,
			Images:      []string{"image1.jpg", "image2.jpg"},
			HasVariants: true,
		}
		productRepo.Create(product)

		// Create two variants
		variant1 := &entity.ProductVariant{
			ID:        1,
			ProductID: 1,
			SKU:       "SKU-1",
			Price:     9999,
			Stock:     50,
			Attributes: []entity.VariantAttribute{
				{Name: "Color", Value: "Red"},
			},
			Images:    []string{"red.jpg"},
			IsDefault: true,
		}
		productVariantRepo.Create(variant1)

		variant2 := &entity.ProductVariant{
			ID:        2,
			ProductID: 1,
			SKU:       "SKU-2",
			Price:     10999,
			Stock:     50,
			Attributes: []entity.VariantAttribute{
				{Name: "Color", Value: "Blue"},
			},
			Images:    []string{"blue.jpg"},
			IsDefault: false,
		}
		productVariantRepo.Create(variant2)

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
		)

		// Execute - delete the non-default variant
		err := productUseCase.DeleteVariant(1, 2)

		// Assert
		assert.NoError(t, err)

		// Check that the variant is deleted
		deletedVariant, err := productVariantRepo.GetByID(2)
		assert.Error(t, err)
		assert.Nil(t, deletedVariant)

		// Default variant should still exist
		defaultVariant, err := productVariantRepo.GetByID(1)
		assert.NoError(t, err)
		assert.NotNil(t, defaultVariant)
	})

	t.Run("Delete default variant should set another as default", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()

		// Create a test product with variants
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       9999,
			Stock:       100,
			CategoryID:  1,
			Images:      []string{"image1.jpg", "image2.jpg"},
			HasVariants: true,
		}
		productRepo.Create(product)

		// Create two variants
		variant1 := &entity.ProductVariant{
			ID:        1,
			ProductID: 1,
			SKU:       "SKU-1",
			Price:     9999,
			Stock:     50,
			Attributes: []entity.VariantAttribute{
				{Name: "Color", Value: "Red"},
			},
			Images:    []string{"red.jpg"},
			IsDefault: true,
		}
		productVariantRepo.Create(variant1)

		variant2 := &entity.ProductVariant{
			ID:        2,
			ProductID: 1,
			SKU:       "SKU-2",
			Price:     10999,
			Stock:     50,
			Attributes: []entity.VariantAttribute{
				{Name: "Color", Value: "Blue"},
			},
			Images:    []string{"blue.jpg"},
			IsDefault: false,
		}
		productVariantRepo.Create(variant2)

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
		)

		// Execute - delete the default variant
		err := productUseCase.DeleteVariant(1, 1)

		// Assert
		assert.NoError(t, err)

		// The other variant should now be default
		newDefaultVariant, err := productVariantRepo.GetByID(2)
		assert.NoError(t, err)
		assert.True(t, newDefaultVariant.IsDefault)

		// Product price should be updated
		updatedProduct, _ := productUseCase.GetProductByID(1, "USD")
		assert.Equal(t, newDefaultVariant.Price, updatedProduct.Price)
	})

	// t.Run("Cannot delete the only variant", func(t *testing.T) {
	// 	// Setup mocks
	// 	productRepo := mock.NewMockProductRepository()
	// 	categoryRepo := mock.NewMockCategoryRepository()
	// 	productVariantRepo := mock.NewMockProductVariantRepository()
	// 	currencyRepo := mock.NewMockCurrencyRepository()

	// 	// Create a test product with one variant
	// 	product := &entity.Product{
	// 		ID:          1,
	// 		Name:        "Test Product",
	// 		Description: "This is a test product",
	// 		Price:       9999,
	// 		Stock:       100,
	// 		CategoryID:  1,
	// 		Images:      []string{"image1.jpg", "image2.jpg"},
	// 		HasVariants: true,
	// 	}
	// 	productRepo.Create(product)

	// 	// Create one variant
	// 	variant := &entity.ProductVariant{
	// 		ID:        1,
	// 		ProductID: 1,
	// 		SKU:       "SKU-1",
	// 		Price:     9999,
	// 		Stock:     50,
	// 		Attributes: []entity.VariantAttribute{
	// 			{Name: "Color", Value: "Red"},
	// 		},
	// 		Images:    []string{"red.jpg"},
	// 		IsDefault: true,
	// 	}
	// 	productVariantRepo.Create(variant)

	// 	// Create use case with mocks
	// 	productUseCase := usecase.NewProductUseCase(
	// 		productRepo,
	// 		categoryRepo,
	// 		productVariantRepo,
	// 		currencyRepo,
	// 	)

	// 	// Execute - try to delete the only variant
	// 	err := productUseCase.DeleteVariant(1, 1)

	// 	// Assert
	// 	assert.Error(t, err)
	// 	assert.Contains(t, err.Error(), "cannot delete the only variant")
	// })
}

func TestProductUseCase_SearchProducts(t *testing.T) {
	t.Run("Search products by query", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()

		// Create test products
		product1 := &entity.Product{
			ID:          1,
			Name:        "Blue Shirt",
			Description: "A nice blue shirt",
			Price:       2999,
			CategoryID:  1,
		}
		productRepo.Create(product1)

		product2 := &entity.Product{
			ID:          2,
			Name:        "Red T-shirt",
			Description: "A comfortable red t-shirt",
			Price:       1999,
			CategoryID:  1,
		}
		productRepo.Create(product2)

		product3 := &entity.Product{
			ID:          3,
			Name:        "Black Jeans",
			Description: "Stylish black jeans",
			Price:       4999,
			CategoryID:  2,
		}
		productRepo.Create(product3)

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
		)

		// Search by shirt
		input := usecase.SearchProductsInput{
			Query:  "shirt",
			Offset: 0,
			Limit:  10,
		}
		results, _, err := productUseCase.SearchProducts(input)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.Equal(t, "Blue Shirt", results[0].Name)
		assert.Equal(t, "Red T-shirt", results[1].Name)

		// Search by category
		input = usecase.SearchProductsInput{
			CategoryID: 2,
			Offset:     0,
			Limit:      10,
		}
		results, _, err = productUseCase.SearchProducts(input)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Black Jeans", results[0].Name)

		// Search by price range
		input = usecase.SearchProductsInput{
			MinPrice: 20.0,
			MaxPrice: 40.0,
			Offset:   0,
			Limit:    10,
		}
		results, _, err = productUseCase.SearchProducts(input)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "Blue Shirt", results[0].Name)
	})
}

func TestProductUseCase_DeleteProduct(t *testing.T) {
	t.Run("Delete product successfully", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()

		// Create a test product
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "This is a test product",
			Price:       9999,
			Stock:       100,
			CategoryID:  1,
			Images:      []string{"image1.jpg", "image2.jpg"},
			HasVariants: true,
		}
		productRepo.Create(product)

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
		)

		// Execute
		err := productUseCase.DeleteProduct(1)

		// Assert
		assert.NoError(t, err)

		// TODO: Verify that product price is deleted and product variants are deleted

		// Verify that product is deleted
		deletedProduct, err := productRepo.GetByID(1)
		assert.Error(t, err)
		assert.Nil(t, deletedProduct)
	})
}
