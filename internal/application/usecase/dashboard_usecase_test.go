package usecase

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zenfulcode/commercify/internal/domain/dto"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/infrastructure/repository/gorm"
	"github.com/zenfulcode/commercify/testutil"
)

// Integration Tests with Real Database
func TestDashboardUseCase_GetDashboardStats_WithRealData(t *testing.T) {
	// Setup test database
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Setup repositories
	orderRepo := gorm.NewOrderRepository(db)
	userRepo := gorm.NewUserRepository(db)
	productRepo := gorm.NewProductRepository(db)
	dashboardUseCase := NewDashboardUseCase(orderRepo, userRepo, productRepo)

	// Create test data
	now := time.Now()

	// Create test users
	user1, err := entity.NewUser("test1@example.com", "password123", "John", "Doe", entity.RoleUser)
	require.NoError(t, err)
	err = userRepo.Create(user1)
	require.NoError(t, err)

	user2, err := entity.NewUser("test2@example.com", "password123", "Jane", "Smith", entity.RoleUser)
	require.NoError(t, err)
	err = userRepo.Create(user2)
	require.NoError(t, err)

	// Create test products and categories for orders
	category := &entity.Category{
		Name:        "Test Category",
		Description: "Test category for integration test",
	}
	err = db.Create(category).Error
	require.NoError(t, err)

	// Create product variant
	variant, err := entity.NewProductVariant(
		"TEST-SKU-001",
		10,
		9999, // 99.99 in cents
		1.0,
		map[string]string{"size": "M"},
		[]string{"test-image.jpg"},
		true,
	)
	require.NoError(t, err)

	// Create product
	product, err := entity.NewProduct(
		"Test Product",
		"Test product description",
		"USD",
		category.ID,
		[]string{"product-image.jpg"},
		[]*entity.ProductVariant{variant},
		true, // isActive
	)
	require.NoError(t, err)
	err = db.Create(product).Error
	require.NoError(t, err)

	// Create test orders
	order1, err := entity.NewOrder(
		&user1.ID,
		[]entity.OrderItem{
			{
				ProductID:        product.ID,
				ProductVariantID: product.Variants[0].ID,
				Quantity:         2,
				Price:            9999,
			},
		},
		"USD",
		nil, nil,
		entity.CustomerDetails{
			Email:    user1.Email,
			Phone:    "123-456-7890",
			FullName: user1.FirstName + " " + user1.LastName,
		},
	)
	require.NoError(t, err)
	order1.Status = entity.OrderStatusPaid
	order1.CreatedAt = now.AddDate(0, 0, -5) // 5 days ago
	err = orderRepo.Create(order1)
	require.NoError(t, err)

	order2, err := entity.NewOrder(
		&user2.ID,
		[]entity.OrderItem{
			{
				ProductID:        product.ID,
				ProductVariantID: product.Variants[0].ID,
				Quantity:         1,
				Price:            9999,
			},
		},
		"USD",
		nil, nil,
		entity.CustomerDetails{
			Email:    user2.Email,
			Phone:    "098-765-4321",
			FullName: user2.FirstName + " " + user2.LastName,
		},
	)
	require.NoError(t, err)
	order2.Status = entity.OrderStatusCompleted
	order2.CreatedAt = now.AddDate(0, 0, -10) // 10 days ago
	err = orderRepo.Create(order2)
	require.NoError(t, err)

	// Test dashboard stats for last 30 days
	request := dto.DashboardStatsRequest{Days: 30}

	result, err := dashboardUseCase.GetDashboardStats(request)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)

	// Should have revenue from both orders (2 * 9999 + 1 * 9999 = 29997)
	assert.Equal(t, int64(29997), result.TotalRevenue)
	assert.Equal(t, int64(2), result.TotalOrders)
	assert.Equal(t, int64(2), result.TotalCustomers)
	assert.Equal(t, int64(2), result.NewCustomers) // Both users created in test

	// Should have product metrics
	assert.Equal(t, int64(1), result.TotalProducts)    // 1 product created
	assert.Equal(t, int64(1), result.LowStockProducts) // Stock is 10, threshold is 10

	// Should have percentage changes (since there's no previous period data, these should be special values)
	require.NotNil(t, result.RevenueChange)
	require.NotNil(t, result.OrdersChange)
	assert.Equal(t, "up", result.RevenueChange.Direction) // No previous data, so 100% up
	assert.Equal(t, "up", result.OrdersChange.Direction)

	// Should have recent orders
	assert.Len(t, result.RecentOrders, 2)
	assert.Contains(t, []string{result.RecentOrders[0].CustomerName, result.RecentOrders[1].CustomerName}, "John Doe")
	assert.Contains(t, []string{result.RecentOrders[0].CustomerName, result.RecentOrders[1].CustomerName}, "Jane Smith")

	// Should have top products
	assert.Len(t, result.TopProducts, 1)
	assert.Equal(t, "Test Product", result.TopProducts[0].ProductName)
	assert.Equal(t, int64(3), result.TopProducts[0].QuantitySold) // 2 + 1
	assert.Equal(t, int64(29997), result.TopProducts[0].Revenue)

	// Verify the date range is approximately correct (30 days ago to now)
	expectedStart := time.Now().AddDate(0, 0, -30)
	assert.WithinDuration(t, expectedStart, result.PeriodStart, time.Hour*24) // Allow 1 day tolerance
	assert.WithinDuration(t, time.Now(), result.PeriodEnd, time.Hour*24)
}

func TestDashboardUseCase_GetDashboardStats_WithDateRange(t *testing.T) {
	// Setup test database
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Setup repositories
	orderRepo := gorm.NewOrderRepository(db)
	userRepo := gorm.NewUserRepository(db)
	productRepo := gorm.NewProductRepository(db)
	dashboardUseCase := NewDashboardUseCase(orderRepo, userRepo, productRepo)

	// Create test data with specific dates
	specificDate := time.Date(2025, 1, 15, 12, 0, 0, 0, time.UTC)

	// Create test user
	user, err := entity.NewUser("test@example.com", "password123", "Test", "User", entity.RoleUser)
	require.NoError(t, err)
	user.CreatedAt = specificDate // Set specific creation date
	err = userRepo.Create(user)
	require.NoError(t, err)

	// Create category and product
	category := &entity.Category{
		Name:        "Test Category",
		Description: "Test category",
	}
	err = db.Create(category).Error
	require.NoError(t, err)

	variant, err := entity.NewProductVariant(
		"TEST-SKU-002",
		5,
		5000, // 50.00 in cents
		0.5,
		map[string]string{"color": "blue"},
		[]string{},
		true,
	)
	require.NoError(t, err)

	product, err := entity.NewProduct(
		"Test Product 2",
		"Another test product",
		"USD",
		category.ID,
		[]string{},
		[]*entity.ProductVariant{variant},
		true,
	)
	require.NoError(t, err)
	err = db.Create(product).Error
	require.NoError(t, err)

	// Create order within the date range
	order, err := entity.NewOrder(
		&user.ID,
		[]entity.OrderItem{
			{
				ProductID:        product.ID,
				ProductVariantID: product.Variants[0].ID,
				Quantity:         1,
				Price:            5000,
			},
		},
		"USD",
		nil, nil,
		entity.CustomerDetails{
			Email:    user.Email,
			Phone:    "555-0123",
			FullName: user.FirstName + " " + user.LastName,
		},
	)
	require.NoError(t, err)
	order.Status = entity.OrderStatusPaid
	order.CreatedAt = specificDate
	err = orderRepo.Create(order)
	require.NoError(t, err)

	// Test with specific date range (January 1-31, 2025)
	startDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 31, 23, 59, 59, 0, time.UTC)

	request := dto.DashboardStatsRequest{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	result, err := dashboardUseCase.GetDashboardStats(request)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, int64(5000), result.TotalRevenue)
	assert.Equal(t, int64(1), result.TotalOrders)
	assert.Equal(t, int64(1), result.TotalCustomers)
	assert.Equal(t, int64(1), result.NewCustomers)
	assert.Len(t, result.RecentOrders, 1)
	assert.Len(t, result.TopProducts, 1)
	assert.Equal(t, startDate, result.PeriodStart)
	assert.Equal(t, endDate, result.PeriodEnd)
}

func TestDashboardUseCase_GetDashboardStats_EmptyData(t *testing.T) {
	// Setup test database with no data
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Setup repositories
	orderRepo := gorm.NewOrderRepository(db)
	userRepo := gorm.NewUserRepository(db)
	productRepo := gorm.NewProductRepository(db)
	dashboardUseCase := NewDashboardUseCase(orderRepo, userRepo, productRepo)

	// Test dashboard stats with no data
	request := dto.DashboardStatsRequest{Days: 30}

	result, err := dashboardUseCase.GetDashboardStats(request)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)

	assert.Equal(t, int64(0), result.TotalRevenue)
	assert.Equal(t, int64(0), result.TotalOrders)
	assert.Equal(t, int64(0), result.TotalCustomers)
	assert.Equal(t, int64(0), result.NewCustomers)
	assert.Len(t, result.RecentOrders, 0)
	assert.Len(t, result.TopProducts, 0)

	// Verify the date range is set correctly
	expectedStart := time.Now().AddDate(0, 0, -30)
	assert.WithinDuration(t, expectedStart, result.PeriodStart, time.Hour*24)
	assert.WithinDuration(t, time.Now(), result.PeriodEnd, time.Hour*24)
}

func TestDashboardUseCase_GetDashboardStats_DefaultRange(t *testing.T) {
	// Setup test database
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Setup repositories
	orderRepo := gorm.NewOrderRepository(db)
	userRepo := gorm.NewUserRepository(db)
	productRepo := gorm.NewProductRepository(db)
	dashboardUseCase := NewDashboardUseCase(orderRepo, userRepo, productRepo)

	// Test with empty request (should default to 30 days)
	request := dto.DashboardStatsRequest{}

	result, err := dashboardUseCase.GetDashboardStats(request)

	// Assertions
	require.NoError(t, err)
	require.NotNil(t, result)

	// Verify the date range is approximately 30 days (default)
	expectedStart := time.Now().AddDate(0, 0, -30)
	assert.WithinDuration(t, expectedStart, result.PeriodStart, time.Hour*24)
	assert.WithinDuration(t, time.Now(), result.PeriodEnd, time.Hour*24)
}

func TestDashboardUseCase_GetDashboardStats_InvalidDateRange(t *testing.T) {
	// Setup test database
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Setup repositories
	orderRepo := gorm.NewOrderRepository(db)
	userRepo := gorm.NewUserRepository(db)
	productRepo := gorm.NewProductRepository(db)
	dashboardUseCase := NewDashboardUseCase(orderRepo, userRepo, productRepo)

	// Test with invalid date range (start after end)
	startDate := time.Date(2025, 2, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)

	request := dto.DashboardStatsRequest{
		StartDate: &startDate,
		EndDate:   &endDate,
	}

	// Execute
	result, err := dashboardUseCase.GetDashboardStats(request)

	// Assertions
	require.Error(t, err)
	assert.Nil(t, result)
	assert.Contains(t, err.Error(), "start date cannot be after end date")
}

func TestDashboardUseCase_GetDashboardStats_OnlyPaidOrders(t *testing.T) {
	// Setup test database
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	// Setup repositories
	orderRepo := gorm.NewOrderRepository(db)
	userRepo := gorm.NewUserRepository(db)
	productRepo := gorm.NewProductRepository(db)
	dashboardUseCase := NewDashboardUseCase(orderRepo, userRepo, productRepo)

	// Create test user
	user, err := entity.NewUser("test@example.com", "password123", "Test", "User", entity.RoleUser)
	require.NoError(t, err)
	err = userRepo.Create(user)
	require.NoError(t, err)

	// Create category and product
	category := &entity.Category{
		Name:        "Test Category",
		Description: "Test category",
	}
	err = db.Create(category).Error
	require.NoError(t, err)

	variant, err := entity.NewProductVariant("SKU-001", 10, 1000, 1.0, map[string]string{}, []string{}, true)
	require.NoError(t, err)

	product, err := entity.NewProduct("Test Product", "Description", "USD", category.ID, []string{}, []*entity.ProductVariant{variant}, true)
	require.NoError(t, err)
	err = db.Create(product).Error
	require.NoError(t, err)

	// Create orders with different statuses
	orderPaid, err := entity.NewOrder(
		&user.ID,
		[]entity.OrderItem{{ProductID: product.ID, ProductVariantID: product.Variants[0].ID, Quantity: 1, Price: 1000}},
		"USD", nil, nil,
		entity.CustomerDetails{Email: user.Email, Phone: "555-0123", FullName: "Test User"},
	)
	require.NoError(t, err)
	orderPaid.Status = entity.OrderStatusPaid
	err = orderRepo.Create(orderPaid)
	require.NoError(t, err)

	// Create second user to avoid order number conflicts
	user2, err := entity.NewUser("test2@example.com", "password456", "Test2", "User2", entity.RoleUser)
	require.NoError(t, err)
	err = userRepo.Create(user2)
	require.NoError(t, err)

	orderPending, err := entity.NewOrder(
		&user2.ID,
		[]entity.OrderItem{{ProductID: product.ID, ProductVariantID: product.Variants[0].ID, Quantity: 1, Price: 1000}},
		"USD", nil, nil,
		entity.CustomerDetails{Email: user2.Email, Phone: "555-0124", FullName: "Test2 User2"},
	)
	require.NoError(t, err)
	orderPending.Status = entity.OrderStatusPending // This should not be included in revenue
	err = orderRepo.Create(orderPending)
	require.NoError(t, err)

	// Test dashboard stats
	request := dto.DashboardStatsRequest{Days: 30}
	result, err := dashboardUseCase.GetDashboardStats(request)

	require.NoError(t, err)
	require.NotNil(t, result)

	// Only paid order should contribute to revenue
	assert.Equal(t, int64(1000), result.TotalRevenue) // Only the paid order
	assert.Equal(t, int64(2), result.TotalOrders)     // Both orders should be counted in total orders

	// Top products should only include quantities from paid orders
	assert.Len(t, result.TopProducts, 1)
	assert.Equal(t, int64(1), result.TopProducts[0].QuantitySold) // Only from paid order
	assert.Equal(t, int64(1000), result.TopProducts[0].Revenue)   // Only from paid order
}
