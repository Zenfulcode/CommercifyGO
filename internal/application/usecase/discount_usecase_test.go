package usecase_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/testutil/mock"
)

func TestDiscountUseCase_CreateDiscount(t *testing.T) {
	t.Run("Create basket percentage discount successfully", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		now := time.Now()
		startDate := now.Add(-24 * time.Hour)
		endDate := now.Add(30 * 24 * time.Hour)

		// Create discount input
		input := usecase.CreateDiscountInput{
			Code:             "TEST10",
			Type:             string(entity.DiscountTypeBasket),
			Method:           string(entity.DiscountMethodPercentage),
			Value:            10.0,
			MinOrderValue:    50.0,
			MaxDiscountValue: 30.0,
			StartDate:        startDate,
			EndDate:          endDate,
			UsageLimit:       100,
		}

		// Execute
		discount, err := discountUseCase.CreateDiscount(input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, discount)
		assert.Equal(t, input.Code, discount.Code)
		assert.Equal(t, entity.DiscountTypeBasket, discount.Type)
		assert.Equal(t, entity.DiscountMethodPercentage, discount.Method)
		assert.Equal(t, input.Value, discount.Value)
		assert.Equal(t, money.ToCents(input.MinOrderValue), discount.MinOrderValue)
		assert.Equal(t, money.ToCents(input.MaxDiscountValue), discount.MaxDiscountValue)
		assert.Equal(t, input.UsageLimit, discount.UsageLimit)
		assert.Equal(t, 0, discount.CurrentUsage)
		assert.True(t, discount.Active)
	})

	t.Run("Create product fixed discount successfully", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create product
		product := &entity.Product{
			ID:   1,
			Name: "Test Product",
		}
		productRepo.Create(product)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		now := time.Now()
		startDate := now.Add(-24 * time.Hour)
		endDate := now.Add(30 * 24 * time.Hour)

		// Create discount input
		input := usecase.CreateDiscountInput{
			Code:       "PRODUCT10",
			Type:       string(entity.DiscountTypeProduct),
			Method:     string(entity.DiscountMethodFixed),
			Value:      10.0,
			ProductIDs: []uint{1},
			StartDate:  startDate,
			EndDate:    endDate,
		}

		// Execute
		discount, err := discountUseCase.CreateDiscount(input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, discount)
		assert.Equal(t, input.Code, discount.Code)
		assert.Equal(t, entity.DiscountTypeProduct, discount.Type)
		assert.Equal(t, entity.DiscountMethodFixed, discount.Method)
		assert.Equal(t, input.Value, discount.Value)
		assert.Equal(t, input.ProductIDs, discount.ProductIDs)
	})

	t.Run("Create category percentage discount successfully", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create a test category
		category := &entity.Category{
			ID:   1,
			Name: "Test Category",
		}
		categoryRepo.Create(category)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		now := time.Now()
		startDate := now.Add(-24 * time.Hour)
		endDate := now.Add(30 * 24 * time.Hour)

		// Create discount input with category
		input := usecase.CreateDiscountInput{
			Code:        "CATEGORY20",
			Type:        string(entity.DiscountTypeProduct),
			Method:      string(entity.DiscountMethodPercentage),
			Value:       20.0,
			CategoryIDs: []uint{1},
			StartDate:   startDate,
			EndDate:     endDate,
		}

		// Execute
		discount, err := discountUseCase.CreateDiscount(input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, discount)
		assert.Equal(t, input.Code, discount.Code)
		assert.Equal(t, entity.DiscountTypeProduct, discount.Type)
		assert.Equal(t, entity.DiscountMethodPercentage, discount.Method)
		assert.Equal(t, input.Value, discount.Value)
		assert.Equal(t, input.CategoryIDs, discount.CategoryIDs)
		assert.Empty(t, discount.ProductIDs)
		assert.True(t, discount.Active)
	})

	t.Run("Create discount with duplicate code", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create existing discount
		existingDiscount, _ := entity.NewDiscount(
			"DUPLICATE",
			entity.DiscountTypeBasket,
			entity.DiscountMethodPercentage,
			10.0,
			0,
			0,
			[]uint{},
			[]uint{},
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			0,
		)
		discountRepo.Create(existingDiscount)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		now := time.Now()
		startDate := now.Add(-24 * time.Hour)
		endDate := now.Add(30 * 24 * time.Hour)

		// Create discount input with duplicate code
		input := usecase.CreateDiscountInput{
			Code:      "DUPLICATE",
			Type:      string(entity.DiscountTypeBasket),
			Method:    string(entity.DiscountMethodPercentage),
			Value:     10.0,
			StartDate: startDate,
			EndDate:   endDate,
		}

		// Execute
		discount, err := discountUseCase.CreateDiscount(input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, discount)
		assert.Contains(t, err.Error(), "discount code already exists")
	})

	t.Run("Create product discount without products or categories", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		now := time.Now()
		startDate := now.Add(-24 * time.Hour)
		endDate := now.Add(30 * 24 * time.Hour)

		// Create discount input
		input := usecase.CreateDiscountInput{
			Code:        "INVALID",
			Type:        string(entity.DiscountTypeProduct),
			Method:      string(entity.DiscountMethodPercentage),
			Value:       10.0,
			ProductIDs:  []uint{},
			CategoryIDs: []uint{},
			StartDate:   startDate,
			EndDate:     endDate,
		}

		// Execute
		discount, err := discountUseCase.CreateDiscount(input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, discount)
		assert.Contains(t, err.Error(), "product discount must specify at least one product or category")
	})

	t.Run("Create discount with invalid product ID", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		now := time.Now()
		startDate := now.Add(-24 * time.Hour)
		endDate := now.Add(30 * 24 * time.Hour)

		// Create discount input with non-existent product ID
		input := usecase.CreateDiscountInput{
			Code:       "INVALID_PRODUCT",
			Type:       string(entity.DiscountTypeProduct),
			Method:     string(entity.DiscountMethodPercentage),
			Value:      10.0,
			ProductIDs: []uint{999}, // Non-existent product
			StartDate:  startDate,
			EndDate:    endDate,
		}

		// Execute
		discount, err := discountUseCase.CreateDiscount(input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, discount)
		assert.Contains(t, err.Error(), "invalid product ID")
	})

	t.Run("Create discount with invalid category ID", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		now := time.Now()
		startDate := now.Add(-24 * time.Hour)
		endDate := now.Add(30 * 24 * time.Hour)

		// Create discount input with non-existent category ID
		input := usecase.CreateDiscountInput{
			Code:        "INVALID_CATEGORY",
			Type:        string(entity.DiscountTypeProduct),
			Method:      string(entity.DiscountMethodPercentage),
			Value:       10.0,
			CategoryIDs: []uint{999}, // Non-existent category
			StartDate:   startDate,
			EndDate:     endDate,
		}

		// Execute
		discount, err := discountUseCase.CreateDiscount(input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, discount)
		assert.Contains(t, err.Error(), "invalid category ID")
	})
}

func TestDiscountUseCase_ProductSpecificDiscount(t *testing.T) {
	t.Run("Create product-specific fixed amount discount", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create test products
		product1 := &entity.Product{
			ID:    1,
			Name:  "Premium Headphones",
			Price: money.ToCents(200.0),
		}
		product2 := &entity.Product{
			ID:    2,
			Name:  "Budget Headphones",
			Price: money.ToCents(50.0),
		}
		productRepo.Create(product1)
		productRepo.Create(product2)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		now := time.Now()
		startDate := now.Add(-24 * time.Hour)
		endDate := now.Add(30 * 24 * time.Hour)

		// Create discount input for specific products
		input := usecase.CreateDiscountInput{
			Code:       "PREMIUM20",
			Type:       string(entity.DiscountTypeProduct),
			Method:     string(entity.DiscountMethodFixed),
			Value:      20.0,
			ProductIDs: []uint{1}, // Only apply to product ID 1 (Premium Headphones)
			StartDate:  startDate,
			EndDate:    endDate,
		}

		// Execute
		discount, err := discountUseCase.CreateDiscount(input)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, discount)
		assert.Equal(t, input.Code, discount.Code)
		assert.Equal(t, entity.DiscountTypeProduct, discount.Type)
		assert.Equal(t, entity.DiscountMethodFixed, discount.Method)
		assert.Equal(t, input.Value, discount.Value)
		assert.Equal(t, input.ProductIDs, discount.ProductIDs)
		assert.Empty(t, discount.CategoryIDs)
		assert.True(t, discount.Active)
	})

	t.Run("Apply product-specific fixed amount discount to order", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create test products
		product1 := &entity.Product{
			ID:    1,
			Name:  "Premium Headphones",
			Price: money.ToCents(200.0),
		}
		product2 := &entity.Product{
			ID:    2,
			Name:  "Budget Headphones",
			Price: money.ToCents(50.0),
		}
		productRepo.Create(product1)
		productRepo.Create(product2)

		// Create a test discount for the product
		discount, _ := entity.NewDiscount(
			"PREMIUM20",
			entity.DiscountTypeProduct,
			entity.DiscountMethodFixed,
			20.0,
			0,
			0,
			[]uint{1}, // Only apply to product ID 1 (Premium Headphones)
			[]uint{},
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			0,
		)
		discountRepo.Create(discount)

		// Create test order items
		items := []entity.OrderItem{
			{
				ProductID: 1, // Premium Headphones with discount
				Quantity:  2,
				Price:     money.ToCents(200.0),
				Subtotal:  money.ToCents(400.0),
			},
			{
				ProductID: 2, // Budget Headphones without discount
				Quantity:  1,
				Price:     money.ToCents(50.0),
				Subtotal:  money.ToCents(50.0),
			},
		}

		// Create test order
		order, _ := entity.NewOrder(
			1,
			items,
			"USD",
			entity.Address{Street: "123 Main St"},
			entity.Address{Street: "123 Main St"},
			entity.CustomerDetails{
				Email:    "test@example.com",
				Phone:    "1234567890",
				FullName: "John Doe",
			},
		)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Apply discount input
		input := usecase.ApplyDiscountToOrderInput{
			OrderID:      order.ID,
			DiscountCode: "PREMIUM20",
		}

		// Execute
		updatedOrder, err := discountUseCase.ApplyDiscountToOrder(input, order)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedOrder)
		// Fixed discount of $20 is applied once to the product (not per quantity)
		assert.Equal(t, money.ToCents(20.0), updatedOrder.DiscountAmount)
		// Total is $450, discount is $20, so final amount should be $430
		assert.Equal(t, money.ToCents(430.0), updatedOrder.FinalAmount)
		assert.NotNil(t, updatedOrder.AppliedDiscount)
		assert.Equal(t, discount.ID, updatedOrder.AppliedDiscount.DiscountID)
		assert.Equal(t, discount.Code, updatedOrder.AppliedDiscount.DiscountCode)
		assert.Equal(t, money.ToCents(20.0), updatedOrder.AppliedDiscount.DiscountAmount)
	})

	t.Run("Apply product-specific percentage discount to order", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create test products
		product1 := &entity.Product{
			ID:    1,
			Name:  "Premium Headphones",
			Price: money.ToCents(200.0),
		}
		product2 := &entity.Product{
			ID:    2,
			Name:  "Budget Headphones",
			Price: money.ToCents(50.0),
		}
		productRepo.Create(product1)
		productRepo.Create(product2)

		// Create a test discount for the product
		discount, _ := entity.NewDiscount(
			"PREMIUM10PERCENT",
			entity.DiscountTypeProduct,
			entity.DiscountMethodPercentage,
			10.0,
			0,
			0,
			[]uint{1}, // Only apply to product ID 1 (Premium Headphones)
			[]uint{},
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			0,
		)
		discountRepo.Create(discount)

		// Create test order items
		items := []entity.OrderItem{
			{
				ProductID: 1, // Premium Headphones with discount
				Quantity:  2,
				Price:     money.ToCents(200.0),
				Subtotal:  money.ToCents(400.0),
			},
			{
				ProductID: 2, // Budget Headphones without discount
				Quantity:  1,
				Price:     money.ToCents(50.0),
				Subtotal:  money.ToCents(50.0),
			},
		}

		// Create test order
		order, _ := entity.NewOrder(
			1,
			items,
			"USD",
			entity.Address{Street: "123 Main St"},
			entity.Address{Street: "123 Main St"},
			entity.CustomerDetails{
				Email:    "test@example.com",
				Phone:    "1234567890",
				FullName: "John Doe",
			},
		)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Apply discount input
		input := usecase.ApplyDiscountToOrderInput{
			OrderID:      order.ID,
			DiscountCode: "PREMIUM10PERCENT",
		}

		// Execute
		updatedOrder, err := discountUseCase.ApplyDiscountToOrder(input, order)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedOrder)
		// 10% of Premium Headphones total (10% of $400) = $40
		assert.Equal(t, money.ToCents(40.0), updatedOrder.DiscountAmount)
		// Total is $450, discount is $40, so final amount should be $410
		assert.Equal(t, money.ToCents(410.0), updatedOrder.FinalAmount)
		assert.NotNil(t, updatedOrder.AppliedDiscount)
		assert.Equal(t, discount.ID, updatedOrder.AppliedDiscount.DiscountID)
		assert.Equal(t, discount.Code, updatedOrder.AppliedDiscount.DiscountCode)
		assert.Equal(t, money.ToCents(40.0), updatedOrder.AppliedDiscount.DiscountAmount)
	})

	t.Run("Apply product-specific discount with maximum discount cap", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create test products
		product1 := &entity.Product{
			ID:    1,
			Name:  "Premium Headphones",
			Price: money.ToCents(200.0),
		}
		product2 := &entity.Product{
			ID:    2,
			Name:  "Budget Headphones",
			Price: money.ToCents(50.0),
		}
		productRepo.Create(product1)
		productRepo.Create(product2)

		// Create a test discount for multiple products with maximum discount cap
		discount, _ := entity.NewDiscount(
			"HEADPHONES25",
			entity.DiscountTypeProduct,
			entity.DiscountMethodPercentage,
			25.0,
			0,
			money.ToCents(30.0), // Maximum discount of $30
			[]uint{1, 2},        // Apply to both Premium and Budget Headphones
			[]uint{},
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			0,
		)
		discountRepo.Create(discount)

		// Create test order items
		items := []entity.OrderItem{
			{
				ProductID: 1, // Premium Headphones
				Quantity:  1,
				Price:     money.ToCents(200.0),
				Subtotal:  money.ToCents(200.0),
			},
			{
				ProductID: 2, // Budget Headphones
				Quantity:  1,
				Price:     money.ToCents(50.0),
				Subtotal:  money.ToCents(50.0),
			},
		}

		// Create test order
		order, _ := entity.NewOrder(
			1,
			items,
			"USD",
			entity.Address{Street: "123 Main St"},
			entity.Address{Street: "123 Main St"},
			entity.CustomerDetails{
				Email:    "test@example.com",
				Phone:    "1234567890",
				FullName: "John Doe",
			},
		)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Apply discount input
		input := usecase.ApplyDiscountToOrderInput{
			OrderID:      order.ID,
			DiscountCode: "HEADPHONES25",
		}

		// Execute
		updatedOrder, err := discountUseCase.ApplyDiscountToOrder(input, order)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedOrder)
		// 25% of ($200 + $50) = $62.50, but capped at $30
		assert.Equal(t, money.ToCents(30.0), updatedOrder.DiscountAmount)
		// Total is $250, discount is $30, so final amount should be $220
		assert.Equal(t, money.ToCents(220.0), updatedOrder.FinalAmount)
		assert.NotNil(t, updatedOrder.AppliedDiscount)
		assert.Equal(t, discount.ID, updatedOrder.AppliedDiscount.DiscountID)
		assert.Equal(t, discount.Code, updatedOrder.AppliedDiscount.DiscountCode)
		assert.Equal(t, money.ToCents(30.0), updatedOrder.AppliedDiscount.DiscountAmount)
	})
}

func TestDiscountUseCase_GetDiscountByID(t *testing.T) {
	t.Run("Get existing discount", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create a test discount
		discount, _ := entity.NewDiscount(
			"TEST10",
			entity.DiscountTypeBasket,
			entity.DiscountMethodPercentage,
			10.0,
			0,
			0,
			[]uint{},
			[]uint{},
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			0,
		)
		discountRepo.Create(discount)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Execute
		result, err := discountUseCase.GetDiscountByID(discount.ID)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, discount.ID, result.ID)
		assert.Equal(t, discount.Code, result.Code)
	})

	t.Run("Get non-existent discount", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Execute with non-existent ID
		result, err := discountUseCase.GetDiscountByID(999)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestDiscountUseCase_GetDiscountByCode(t *testing.T) {
	t.Run("Get existing discount by code", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create a test discount
		discount, _ := entity.NewDiscount(
			"TESTCODE",
			entity.DiscountTypeBasket,
			entity.DiscountMethodPercentage,
			10.0,
			0,
			0,
			[]uint{},
			[]uint{},
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			0,
		)
		discountRepo.Create(discount)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Execute
		result, err := discountUseCase.GetDiscountByCode("TESTCODE")

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, discount.ID, result.ID)
		assert.Equal(t, discount.Code, result.Code)
	})

	t.Run("Get non-existent discount code", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Execute with non-existent code
		result, err := discountUseCase.GetDiscountByCode("NONEXISTENT")

		// Assert
		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestDiscountUseCase_UpdateDiscount(t *testing.T) {
	t.Run("Update discount successfully", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create a test discount
		discount, _ := entity.NewDiscount(
			"OLD_CODE",
			entity.DiscountTypeBasket,
			entity.DiscountMethodPercentage,
			10.0,
			0,
			0,
			[]uint{},
			[]uint{},
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			100,
		)
		discountRepo.Create(discount)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Update input
		input := usecase.UpdateDiscountInput{
			Code:             "NEW_CODE",
			Value:            20.0,
			MinOrderValue:    50.0,
			MaxDiscountValue: 30.0,
			UsageLimit:       200,
			Active:           true,
		}

		// Execute
		updatedDiscount, err := discountUseCase.UpdateDiscount(discount.ID, input)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, input.Code, updatedDiscount.Code)
		assert.Equal(t, input.Value, updatedDiscount.Value)
		assert.Equal(t, money.ToCents(input.MinOrderValue), updatedDiscount.MinOrderValue)
		assert.Equal(t, money.ToCents(input.MaxDiscountValue), updatedDiscount.MaxDiscountValue)
		assert.Equal(t, input.UsageLimit, updatedDiscount.UsageLimit)
		assert.Equal(t, input.Active, updatedDiscount.Active)
	})

	t.Run("Update non-existent discount", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Update input
		input := usecase.UpdateDiscountInput{
			Code:  "NEW_CODE",
			Value: 20.0,
		}

		// Execute with non-existent ID
		updatedDiscount, err := discountUseCase.UpdateDiscount(999, input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, updatedDiscount)
	})

	t.Run("Update with duplicate code", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create two test discounts
		discount1, _ := entity.NewDiscount(
			"CODE1",
			entity.DiscountTypeBasket,
			entity.DiscountMethodPercentage,
			10.0,
			0,
			0,
			[]uint{},
			[]uint{},
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			0,
		)
		discount2, _ := entity.NewDiscount(
			"CODE2",
			entity.DiscountTypeBasket,
			entity.DiscountMethodPercentage,
			20.0,
			0,
			0,
			[]uint{},
			[]uint{},
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			0,
		)
		discountRepo.Create(discount1)
		discountRepo.Create(discount2)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Update input with duplicate code
		input := usecase.UpdateDiscountInput{
			Code: "CODE1", // Already exists
		}

		// Execute - try to update discount2 to use code1
		updatedDiscount, err := discountUseCase.UpdateDiscount(discount2.ID, input)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, updatedDiscount)
		assert.Contains(t, err.Error(), "discount code already exists")
	})
}

func TestDiscountUseCase_DeleteDiscount(t *testing.T) {
	t.Run("Delete discount successfully", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create a test discount
		discount, _ := entity.NewDiscount(
			"DELETE_ME",
			entity.DiscountTypeBasket,
			entity.DiscountMethodPercentage,
			10.0,
			0,
			0,
			[]uint{},
			[]uint{},
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			0,
		)
		discountRepo.Create(discount)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Execute
		err := discountUseCase.DeleteDiscount(discount.ID)

		// Assert
		assert.NoError(t, err)

		// Verify discount was deleted
		_, err = discountRepo.GetByID(discount.ID)
		assert.Error(t, err)
	})

	t.Run("Delete discount that is in use by an order", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(true)

		// Create a test discount
		discount, _ := entity.NewDiscount(
			"IN_USE",
			entity.DiscountTypeBasket,
			entity.DiscountMethodPercentage,
			10.0,
			0,
			0,
			[]uint{},
			[]uint{},
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			0,
		)
		discountRepo.Create(discount)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Execute
		err := discountUseCase.DeleteDiscount(discount.ID)

		// Assert
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "discount is in use by an order")

		// Verify discount was not deleted
		_, err = discountRepo.GetByID(discount.ID)
		assert.NoError(t, err)
	})
}

func TestDiscountUseCase_ListDiscounts(t *testing.T) {
	t.Run("List discounts with pagination", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create multiple test discounts
		for i := 1; i <= 5; i++ {
			code := "CODE_" + time.Now().Add(time.Duration(i)*time.Hour).Format("150405")
			discount, _ := entity.NewDiscount(
				code,
				entity.DiscountTypeBasket,
				entity.DiscountMethodPercentage,
				float64(i*10),
				0,
				0,
				[]uint{},
				[]uint{},
				time.Now().Add(-24*time.Hour),
				time.Now().Add(30*24*time.Hour),
				0,
			)
			discountRepo.Create(discount)
		}

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Execute - first page
		discounts, err := discountUseCase.ListDiscounts(0, 3)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, discounts, 3)

		// Execute - second page
		discounts, err = discountUseCase.ListDiscounts(3, 3)

		// Assert
		assert.NoError(t, err)
		assert.Len(t, discounts, 2)
	})
}

func TestDiscountUseCase_ApplyDiscountToOrder(t *testing.T) {
	t.Run("Apply valid basket discount to order", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create a test discount
		discount, _ := entity.NewDiscount(
			"BASKET10",
			entity.DiscountTypeBasket,
			entity.DiscountMethodPercentage,
			10.0,
			0,
			0,
			[]uint{},
			[]uint{},
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			1,
		)
		discountRepo.Create(discount)

		// Create test order items
		items := []entity.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
				Price:     5000,
				Subtotal:  10000,
			},
			{
				ProductID: 2,
				Quantity:  1,
				Price:     1000,
				Subtotal:  1000,
			},
		}

		// Create test order
		order, _ := entity.NewOrder(
			1,
			items,
			"USD",
			entity.Address{Street: "123 Main St"},
			entity.Address{Street: "123 Main St"},
			entity.CustomerDetails{
				Email:    "test@example.com",
				Phone:    "1234567890",
				FullName: "John Doe",
			},
		)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Apply discount input
		input := usecase.ApplyDiscountToOrderInput{
			OrderID:      order.ID,
			DiscountCode: "BASKET10",
		}

		// Execute
		updatedOrder, err := discountUseCase.ApplyDiscountToOrder(input, order)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedOrder)
		assert.Equal(t, money.ToCents(11.0), updatedOrder.DiscountAmount) // 10% of 110 = 11
		assert.Equal(t, money.ToCents(99.0), updatedOrder.FinalAmount)    // 110 - 11 = 10
		assert.NotNil(t, updatedOrder.AppliedDiscount)
		assert.Equal(t, discount.ID, updatedOrder.AppliedDiscount.DiscountID)
		assert.Equal(t, discount.Code, updatedOrder.AppliedDiscount.DiscountCode)
		assert.Equal(t, money.ToCents(11.0), updatedOrder.AppliedDiscount.DiscountAmount)
	})

	t.Run("Apply category-specific discount to order", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create a test category
		category := &entity.Category{
			ID:   1,
			Name: "Electronics",
		}
		categoryRepo.Create(category)

		// Create some test products in that category
		product1 := &entity.Product{
			ID:         1,
			Name:       "Phone",
			CategoryID: 1,
			Price:      100.0,
		}
		product2 := &entity.Product{
			ID:         2,
			Name:       "Laptop",
			CategoryID: 1,
			Price:      1000.0,
		}
		productRepo.Create(product1)
		productRepo.Create(product2)

		// Create a test discount for the Electronics category
		discount, _ := entity.NewDiscount(
			"ELECTRONICS25",
			entity.DiscountTypeProduct,
			entity.DiscountMethodPercentage,
			25.0,
			0,
			0,
			[]uint{},  // No specific products
			[]uint{1}, // Category ID 1 (Electronics)
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			0,
		)
		discountRepo.Create(discount)

		// Create test order items including products from the category
		items := []entity.OrderItem{
			{
				ProductID: 1, // Phone (in Electronics category)
				Quantity:  1,
				Price:     10000,
				Subtotal:  10000,
			},
			{
				ProductID: 2, // Laptop (in Electronics category)
				Quantity:  1,
				Price:     100000,
				Subtotal:  100000,
			},
			{
				ProductID: 3, // Some other product not in Electronics
				Quantity:  1,
				Price:     5000,
				Subtotal:  5000,
			},
		}

		// Create test order
		order, _ := entity.NewOrder(
			1,
			items,
			"USD",
			entity.Address{Street: "123 Main St"},
			entity.Address{Street: "123 Main St"},
			entity.CustomerDetails{
				Email:    "test@example.com",
				Phone:    "1234567890",
				FullName: "John Doe",
			},
		)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Apply discount input
		input := usecase.ApplyDiscountToOrderInput{
			OrderID:      order.ID,
			DiscountCode: "ELECTRONICS25",
		}

		// Execute
		updatedOrder, err := discountUseCase.ApplyDiscountToOrder(input, order)

		// Assert
		assert.NoError(t, err)
		assert.NotNil(t, updatedOrder)
		// Should apply 25% discount to the products in Electronics category
		// 25% of (100 + 1000) = 275
		assert.Equal(t, money.ToCents(275.0), updatedOrder.DiscountAmount)
		// Final amount should be: 100 + 1000 + 50 - 275 = 875
		assert.Equal(t, money.ToCents(875.0), updatedOrder.FinalAmount)
		assert.NotNil(t, updatedOrder.AppliedDiscount)
		assert.Equal(t, discount.ID, updatedOrder.AppliedDiscount.DiscountID)
		assert.Equal(t, discount.Code, updatedOrder.AppliedDiscount.DiscountCode)
		assert.Equal(t, money.ToCents(275.0), updatedOrder.AppliedDiscount.DiscountAmount)
	})

	t.Run("Apply invalid discount code", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create test order items
		items := []entity.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
				Price:     5000,
				Subtotal:  10000,
			},
		}

		// Create test order
		order, _ := entity.NewOrder(
			1,
			items,
			"USD",
			entity.Address{Street: "123 Main St"},
			entity.Address{Street: "123 Main St"},
			entity.CustomerDetails{
				Email:    "test@example.com",
				Phone:    "1234567890",
				FullName: "John Doe",
			},
		)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Apply discount input with invalid code
		input := usecase.ApplyDiscountToOrderInput{
			OrderID:      order.ID,
			DiscountCode: "INVALID",
		}

		// Execute
		updatedOrder, err := discountUseCase.ApplyDiscountToOrder(input, order)

		// Assert
		assert.Error(t, err)
		assert.Nil(t, updatedOrder)
		assert.Contains(t, err.Error(), "invalid discount code")
	})
}

func TestDiscountUseCase_RemoveDiscountFromOrder(t *testing.T) {
	t.Run("Remove discount from order", func(t *testing.T) {
		// Setup mocks
		discountRepo := mock.NewMockDiscountRepository()
		productRepo := mock.NewMockProductRepository()
		categoryRepo := mock.NewMockCategoryRepository()
		orderRepo := mock.NewMockOrderRepository(false)

		// Create a test discount
		discount, _ := entity.NewDiscount(
			"BASKET10",
			entity.DiscountTypeBasket,
			entity.DiscountMethodPercentage,
			10.0,
			0,
			0,
			[]uint{},
			[]uint{},
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			0,
		)
		discountRepo.Create(discount)

		// Create test order with discount already applied
		items := []entity.OrderItem{
			{
				ProductID: 1,
				Quantity:  2,
				Price:     5000,
				Subtotal:  1000,
			},
		}

		order, _ := entity.NewOrder(
			1,
			items,
			"USD",
			entity.Address{Street: "123 Main St"},
			entity.Address{Street: "123 Main St"},
			entity.CustomerDetails{
				Email:    "test@example.com",
				Phone:    "1234567890",
				FullName: "John Doe",
			},
		)

		// Apply discount manually
		order.ApplyDiscount(discount)
		assert.NotNil(t, order.AppliedDiscount)
		assert.Greater(t, order.DiscountAmount, money.ToCents(0.0))
		assert.Less(t, order.FinalAmount, order.TotalAmount)

		// Create use case with mocks
		discountUseCase := usecase.NewDiscountUseCase(
			discountRepo,
			productRepo,
			categoryRepo,
			orderRepo,
		)

		// Execute
		discountUseCase.RemoveDiscountFromOrder(order)

		// Assert
		assert.Nil(t, order.AppliedDiscount)
		assert.Zero(t, order.DiscountAmount)
		assert.Equal(t, order.TotalAmount, order.FinalAmount)
	})
}
