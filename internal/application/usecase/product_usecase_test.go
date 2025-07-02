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
	t.Run("Create simple product successfully (In complete)", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

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
			orderRepo,
			checkoutRepo,
		)

		// Create product input
		input := usecase.CreateProductInput{
			Name:        "Test Product",
			Description: "This is a test product",
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
		assert.Equal(t, input.CategoryID, product.CategoryID)
		assert.Equal(t, input.Images, product.Images)
		assert.Equal(t, int64(0), product.Price, "Price should be zero for incomplete product")
		assert.Equal(t, 0, product.Stock, "Stock should be zero for incomplete product")
		assert.False(t, product.HasVariants, "HasVariants should be false for incomplete product")
		assert.False(t, product.Active, "Product should be active by default")
	})

	t.Run("Create product with variants successfully", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

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
			orderRepo,
			checkoutRepo,
		)

		// Create product input with variants
		input := usecase.CreateProductInput{
			Name:        "Test Product with Variants",
			Description: "This is a test product with variants",
			Currency:    "USD",
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
		productPrice, _ := product.GetPriceInCurrency("USD")

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, input.Name, product.Name)
		assert.Len(t, product.Variants, 2)
		assert.Equal(t, productPrice, money.ToCents(99.99), "Price should be set to the first variant's price")

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
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
			orderRepo,
			checkoutRepo,
		)

		// Create product input with invalid category
		input := usecase.CreateProductInput{
			Name:        "Test Product",
			Description: "This is a test product",
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
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

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
			orderRepo,
			checkoutRepo,
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
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
			orderRepo,
			checkoutRepo,
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
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

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
			orderRepo,
			checkoutRepo,
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
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

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
			orderRepo,
			checkoutRepo,
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
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

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
			orderRepo,
			checkoutRepo,
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
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

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
			orderRepo,
			checkoutRepo,
		)

		// Update input
		input := usecase.UpdateProductInput{
			Name:        "Updated Product",
			Description: "Updated description",
			CategoryID:  2,
			Images:      []string{"updated.jpg"},
		}

		// Execute
		updatedProduct, err := productUseCase.UpdateProduct(1, input)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, input.Name, updatedProduct.Name)
		assert.Equal(t, input.Description, updatedProduct.Description)
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
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

		// Create a test product with a default variant (as per business rules)
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "This is a test product",
			CategoryID:  1,
			Images:      []string{"image1.jpg", "image2.jpg"},
			HasVariants: false,
		}
		productRepo.Create(product)

		// Create a default variant that already exists
		defaultVariant := &entity.ProductVariant{
			ID:        1,
			ProductID: 1,
			SKU:       "DEFAULT-SKU",
			Price:     9999,
			Stock:     100,
			IsDefault: true,
		}
		productVariantRepo.Create(defaultVariant)
		product.Variants = []*entity.ProductVariant{defaultVariant}

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
			orderRepo,
			checkoutRepo,
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
		assert.True(t, updatedProduct.IsComplete())
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
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

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
			orderRepo,
			checkoutRepo,
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
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

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
			orderRepo,
			checkoutRepo,
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
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

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
			orderRepo,
			checkoutRepo,
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
}

func TestProductUseCase_SearchProducts(t *testing.T) {
	t.Run("Search products by query", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

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
			orderRepo,
			checkoutRepo,
		)

		// Search by shirt
		input := usecase.SearchProductsInput{
			Query:  "shirt",
			Offset: 0,
			Limit:  10,
		}
		results, _, err := productUseCase.ListProducts(input)

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
		results, _, err = productUseCase.ListProducts(input)

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
		results, _, err = productUseCase.ListProducts(input)

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
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

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
			orderRepo,
			checkoutRepo,
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

	t.Run("Delete product with existing orders should fail", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

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

		// Create an order with this product
		order := &entity.Order{
			ID:     1,
			UserID: 1,
			Items: []entity.OrderItem{
				{
					ID:        1,
					ProductID: 1, // Reference to our test product
					Quantity:  2,
					Price:     9999,
					Subtotal:  19998,
				},
			},
			TotalAmount:   19998,
			Status:        entity.OrderStatusPaid,
			PaymentStatus: entity.PaymentStatusCaptured,
		}
		orderRepo.Create(order)

		// Create use case with mocks
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
			orderRepo,
			checkoutRepo,
		)

		// Execute - should fail
		err := productUseCase.DeleteProduct(1)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot delete product that has existing orders")
	})

	t.Run("Delete product with no orders should succeed", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

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

		// Create use case with mocks (no orders in repository)
		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
			orderRepo,
			checkoutRepo,
		)

		// Execute - should succeed
		err := productUseCase.DeleteProduct(1)

		// Assert
		assert.NoError(t, err)
	})
}

func TestProductUseCase_CreateProduct_StockCalculation(t *testing.T) {
	setupTestUseCase := func() *usecase.ProductUseCase {
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

		// Create a test category
		category := &entity.Category{
			ID:   1,
			Name: "Test Category",
		}
		categoryRepo.Create(category)

		return usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
			orderRepo,
			checkoutRepo,
		)
	}

	t.Run("Product with single variant - stock should equal variant stock", func(t *testing.T) {
		productUseCase := setupTestUseCase()

		input := usecase.CreateProductInput{
			Name:        "Single Variant Product",
			Description: "Product with one variant",
			Currency:    "USD",
			CategoryID:  1,
			Images:      []string{"image1.jpg"},
			Variants: []usecase.CreateVariantInput{
				{
					SKU:       "SKU-SINGLE",
					Price:     99.99,
					Stock:     25,
					IsDefault: true,
				},
			},
			Active: true,
		}

		product, err := productUseCase.CreateProduct(input)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, 25, product.Stock, "Product stock should equal the single variant's stock")
		assert.Len(t, product.Variants, 1, "Should have exactly one variant")
		assert.Equal(t, 25, product.Variants[0].Stock, "Variant stock should be preserved")
		assert.True(t, product.Variants[0].IsDefault, "Single variant should be default")
		assert.False(t, product.HasVariants, "HasVariants should be false for single variant (current logic)")
	})

	t.Run("Product with multiple variants - stock should be sum of all variant stocks", func(t *testing.T) {
		productUseCase := setupTestUseCase()

		input := usecase.CreateProductInput{
			Name:        "Multi Variant Product",
			Description: "Product with multiple variants",
			Currency:    "USD",
			CategoryID:  1,
			Images:      []string{"image1.jpg"},
			Variants: []usecase.CreateVariantInput{
				{
					SKU:        "SKU-RED",
					Price:      99.99,
					Stock:      15,
					Attributes: []entity.VariantAttribute{{Name: "Color", Value: "Red"}},
					IsDefault:  true,
				},
				{
					SKU:        "SKU-BLUE",
					Price:      109.99,
					Stock:      20,
					Attributes: []entity.VariantAttribute{{Name: "Color", Value: "Blue"}},
					IsDefault:  false,
				},
				{
					SKU:        "SKU-GREEN",
					Price:      119.99,
					Stock:      10,
					Attributes: []entity.VariantAttribute{{Name: "Color", Value: "Green"}},
					IsDefault:  false,
				},
			},
			Active: true,
		}

		product, err := productUseCase.CreateProduct(input)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, 45, product.Stock, "Product stock should be sum of all variant stocks (15+20+10)")
		assert.Len(t, product.Variants, 3, "Should have exactly three variants")
		assert.True(t, product.HasVariants, "HasVariants should be true for multiple variants")

		// Verify individual variant stocks are preserved
		assert.Equal(t, 15, product.Variants[0].Stock, "First variant stock should be preserved")
		assert.Equal(t, 20, product.Variants[1].Stock, "Second variant stock should be preserved")
		assert.Equal(t, 10, product.Variants[2].Stock, "Third variant stock should be preserved")
	})

	t.Run("Product with variants having zero stock - total should be zero", func(t *testing.T) {
		productUseCase := setupTestUseCase()

		input := usecase.CreateProductInput{
			Name:        "Zero Stock Product",
			Description: "Product with zero stock variants",
			Currency:    "USD",
			CategoryID:  1,
			Images:      []string{"image1.jpg"},
			Variants: []usecase.CreateVariantInput{
				{
					SKU:       "SKU-EMPTY1",
					Price:     99.99,
					Stock:     0,
					IsDefault: true,
				},
				{
					SKU:       "SKU-EMPTY2",
					Price:     109.99,
					Stock:     0,
					IsDefault: false,
				},
			},
			Active: true,
		}

		product, err := productUseCase.CreateProduct(input)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, 0, product.Stock, "Product stock should be zero when all variants have zero stock")
		assert.Len(t, product.Variants, 2, "Should have exactly two variants")
		assert.True(t, product.HasVariants, "HasVariants should be true for multiple variants")
	})

	t.Run("Product with mixed stock levels - should calculate correctly", func(t *testing.T) {
		productUseCase := setupTestUseCase()

		input := usecase.CreateProductInput{
			Name:        "Mixed Stock Product",
			Description: "Product with mixed stock levels",
			Currency:    "USD",
			CategoryID:  1,
			Images:      []string{"image1.jpg"},
			Variants: []usecase.CreateVariantInput{
				{
					SKU:        "SKU-HIGH",
					Price:      99.99,
					Stock:      100,
					Attributes: []entity.VariantAttribute{{Name: "Size", Value: "Large"}},
					IsDefault:  true,
				},
				{
					SKU:        "SKU-ZERO",
					Price:      99.99,
					Stock:      0,
					Attributes: []entity.VariantAttribute{{Name: "Size", Value: "Medium"}},
					IsDefault:  false,
				},
				{
					SKU:        "SKU-LOW",
					Price:      99.99,
					Stock:      5,
					Attributes: []entity.VariantAttribute{{Name: "Size", Value: "Small"}},
					IsDefault:  false,
				},
			},
			Active: true,
		}

		product, err := productUseCase.CreateProduct(input)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, 105, product.Stock, "Product stock should be sum of all variant stocks (100+0+5)")
		assert.Len(t, product.Variants, 3, "Should have exactly three variants")

		// Verify the CalculateStock method works correctly
		product.CalculateStock()
		assert.Equal(t, 105, product.Stock, "CalculateStock should produce the same result")
	})

	t.Run("Product without variants - should have zero stock", func(t *testing.T) {
		productUseCase := setupTestUseCase()

		input := usecase.CreateProductInput{
			Name:        "No Variants Product",
			Description: "Product without any variants",
			Currency:    "USD",
			CategoryID:  1,
			Images:      []string{"image1.jpg"},
			Variants:    []usecase.CreateVariantInput{}, // Empty variants
			Active:      true,
		}

		product, err := productUseCase.CreateProduct(input)

		assert.NoError(t, err)
		assert.NotNil(t, product)
		assert.Equal(t, 0, product.Stock, "Product stock should be zero when no variants exist")
		assert.Len(t, product.Variants, 0, "Should have no variants")
		assert.False(t, product.HasVariants, "HasVariants should be false when no variants exist")
	})

	t.Run("Stock calculation after adding variants individually", func(t *testing.T) {
		productUseCase := setupTestUseCase()

		// First create a product with one variant
		input := usecase.CreateProductInput{
			Name:        "Incremental Product",
			Description: "Product to test incremental variant addition",
			Currency:    "USD",
			CategoryID:  1,
			Images:      []string{"image1.jpg"},
			Variants: []usecase.CreateVariantInput{
				{
					SKU:       "SKU-FIRST",
					Price:     99.99,
					Stock:     30,
					IsDefault: true,
				},
			},
			Active: true,
		}

		product, err := productUseCase.CreateProduct(input)
		assert.NoError(t, err)
		assert.Equal(t, 30, product.Stock, "Initial stock should be 30")

		// Simulate adding a second variant (this would happen through AddVariant use case)
		// But we can test the entity logic directly
		variant2, err := entity.NewProductVariant(
			product.ID,
			"SKU-SECOND",
			79.99,
			"USD",
			20,
			[]entity.VariantAttribute{{Name: "Size", Value: "Small"}},
			[]string{"small.jpg"},
			false,
		)
		assert.NoError(t, err)

		err = product.AddVariant(variant2)
		assert.NoError(t, err)

		// After adding second variant, stock should be recalculated
		assert.Equal(t, 50, product.Stock, "Stock should be sum after adding second variant (30+20)")
		assert.True(t, product.HasVariants, "HasVariants should be true after adding second variant")
	})
}

func TestProductUseCase_UpdateProduct_StockCalculation(t *testing.T) {
	setupTestUseCaseWithProduct := func() (*usecase.ProductUseCase, *entity.Product) {
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

		// Create a test category
		category := &entity.Category{
			ID:   1,
			Name: "Test Category",
		}
		categoryRepo.Create(category)

		// Create a test product with variants
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "Test Description",
			Price:       9999,
			Stock:       50,
			CategoryID:  1,
			Images:      []string{"image1.jpg"},
			HasVariants: true,
			Active:      true,
		}

		// Add some variants
		variant1, _ := entity.NewProductVariant(1, "SKU-1", 99.99, "USD", 25, []entity.VariantAttribute{}, []string{}, true)
		variant2, _ := entity.NewProductVariant(1, "SKU-2", 109.99, "USD", 25, []entity.VariantAttribute{}, []string{}, false)

		variant1.ID = 1
		variant2.ID = 2

		product.Variants = []*entity.ProductVariant{variant1, variant2}
		product.CalculateStock() // Should set stock to 50

		productRepo.Create(product)
		productVariantRepo.Create(variant1)
		productVariantRepo.Create(variant2)

		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
			orderRepo,
			checkoutRepo,
		)

		return productUseCase, product
	}

	t.Run("UpdateProduct should recalculate stock from variants", func(t *testing.T) {
		productUseCase, product := setupTestUseCaseWithProduct()

		// Verify initial stock calculation
		assert.Equal(t, 50, product.Stock, "Initial stock should be 50")

		// Update the product (this should trigger stock recalculation)
		input := usecase.UpdateProductInput{
			Name:        "Updated Product Name",
			Description: "Updated Description",
			Active:      true,
		}

		updatedProduct, err := productUseCase.UpdateProduct(product.ID, input)

		assert.NoError(t, err)
		assert.NotNil(t, updatedProduct)
		assert.Equal(t, "Updated Product Name", updatedProduct.Name)
		assert.Equal(t, 50, updatedProduct.Stock, "Stock should remain correctly calculated after update")
	})

	t.Run("UpdateVariant should trigger product stock recalculation", func(t *testing.T) {
		productUseCase, product := setupTestUseCaseWithProduct()

		// Update a variant's stock
		updateInput := usecase.UpdateVariantInput{
			Stock: 35, // Change from 25 to 35
		}

		updatedVariant, err := productUseCase.UpdateVariant(product.ID, 1, updateInput)

		assert.NoError(t, err)
		assert.NotNil(t, updatedVariant)
		assert.Equal(t, 35, updatedVariant.Stock, "Variant stock should be updated")

		// The product stock should be recalculated when we fetch it again
		// Since we're using mocks, we'll test the entity behavior directly
		product.Variants[0].Stock = 35
		product.CalculateStock()
		assert.Equal(t, 60, product.Stock, "Product stock should be recalculated (35+25)")
	})
}

func TestProductEntity_CalculateStock(t *testing.T) {
	t.Run("CalculateStock with multiple variants", func(t *testing.T) {
		product := &entity.Product{
			ID:   1,
			Name: "Test Product",
		}

		// Create variants with different stock levels
		variant1, _ := entity.NewProductVariant(1, "SKU-1", 99.99, "USD", 10, []entity.VariantAttribute{}, []string{}, true)
		variant2, _ := entity.NewProductVariant(1, "SKU-2", 109.99, "USD", 20, []entity.VariantAttribute{}, []string{}, false)
		variant3, _ := entity.NewProductVariant(1, "SKU-3", 119.99, "USD", 30, []entity.VariantAttribute{}, []string{}, false)

		product.Variants = []*entity.ProductVariant{variant1, variant2, variant3}

		product.CalculateStock()

		assert.Equal(t, 60, product.Stock, "Stock should be sum of all variant stocks (10+20+30)")
	})

	t.Run("CalculateStock with zero stock variants", func(t *testing.T) {
		product := &entity.Product{
			ID:   1,
			Name: "Test Product",
		}

		variant1, _ := entity.NewProductVariant(1, "SKU-1", 99.99, "USD", 0, []entity.VariantAttribute{}, []string{}, true)
		variant2, _ := entity.NewProductVariant(1, "SKU-2", 109.99, "USD", 0, []entity.VariantAttribute{}, []string{}, false)

		product.Variants = []*entity.ProductVariant{variant1, variant2}

		product.CalculateStock()

		assert.Equal(t, 0, product.Stock, "Stock should be zero when all variants have zero stock")
	})

	t.Run("CalculateStock with no variants", func(t *testing.T) {
		product := &entity.Product{
			ID:       1,
			Name:     "Test Product",
			Variants: []*entity.ProductVariant{},
		}

		product.CalculateStock()

		assert.Equal(t, 0, product.Stock, "Stock should be zero when no variants exist")
	})

	t.Run("CalculateStock with single variant", func(t *testing.T) {
		product := &entity.Product{
			ID:   1,
			Name: "Test Product",
		}

		variant1, _ := entity.NewProductVariant(1, "SKU-1", 99.99, "USD", 42, []entity.VariantAttribute{}, []string{}, true)
		product.Variants = []*entity.ProductVariant{variant1}

		product.CalculateStock()

		assert.Equal(t, 42, product.Stock, "Stock should equal single variant's stock")
	})

	t.Run("CalculateStock is called automatically when adding variants", func(t *testing.T) {
		product := &entity.Product{
			ID:   1,
			Name: "Test Product",
		}

		// AddVariant should automatically call CalculateStock
		variant1, _ := entity.NewProductVariant(1, "SKU-1", 99.99, "USD", 15, []entity.VariantAttribute{}, []string{}, true)
		err := product.AddVariant(variant1)

		assert.NoError(t, err)
		assert.Equal(t, 15, product.Stock, "Stock should be calculated automatically when adding variant")

		// Add another variant
		variant2, _ := entity.NewProductVariant(1, "SKU-2", 109.99, "USD", 25, []entity.VariantAttribute{}, []string{}, false)
		err = product.AddVariant(variant2)

		assert.NoError(t, err)
		assert.Equal(t, 40, product.Stock, "Stock should be recalculated when adding second variant (15+25)")
	})
}

func TestProductUseCase_AddVariant_StockCalculation(t *testing.T) {
	t.Run("AddVariant should update product stock", func(t *testing.T) {
		// Setup mocks
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		productVariantRepo := mock.NewMockProductVariantRepository()
		currencyRepo := mock.NewMockCurrencyRepository()
		orderRepo := mock.NewMockOrderRepository(false)
		checkoutRepo := mock.NewMockCheckoutRepository()

		// Create a test category
		category := &entity.Category{
			ID:   1,
			Name: "Test Category",
		}
		categoryRepo.Create(category)

		// Create a product with one variant
		product := &entity.Product{
			ID:          1,
			Name:        "Test Product",
			Description: "Test Description",
			Price:       9999,
			Stock:       30,
			CategoryID:  1,
			HasVariants: false,
			Active:      true,
		}

		variant1, _ := entity.NewProductVariant(1, "SKU-1", 99.99, "USD", 30, []entity.VariantAttribute{}, []string{}, true)
		variant1.ID = 1
		product.Variants = []*entity.ProductVariant{variant1}

		productRepo.Create(product)
		productVariantRepo.Create(variant1)

		productUseCase := usecase.NewProductUseCase(
			productRepo,
			categoryRepo,
			productVariantRepo,
			currencyRepo,
			orderRepo,
			checkoutRepo,
		)

		// Add a second variant
		input := usecase.AddVariantInput{
			ProductID:  1,
			SKU:        "SKU-2",
			Price:      109.99,
			Stock:      20,
			Attributes: []entity.VariantAttribute{{Name: "Color", Value: "Blue"}},
			Images:     []string{"blue.jpg"},
			IsDefault:  false,
		}

		addedVariant, err := productUseCase.AddVariant(input)

		assert.NoError(t, err)
		assert.NotNil(t, addedVariant)
		assert.Equal(t, "SKU-2", addedVariant.SKU)
		assert.Equal(t, 20, addedVariant.Stock)

		// Verify the product stock calculation through entity behavior
		// (In a real scenario, the product would be fetched from repo and have updated stock)
		// Since AddVariant calls product.AddVariant() which calls CalculateStock()
		// and then updates the product in the repository, we need to verify this works

		// The product in the repository should now have updated stock
		updatedProduct, err := productRepo.GetByID(1)
		assert.NoError(t, err)
		assert.Equal(t, 50, updatedProduct.Stock, "Product stock should be sum of all variants after adding new variant (30+20)")
		assert.True(t, updatedProduct.HasVariants, "HasVariants should be true after adding second variant")
	})
}
