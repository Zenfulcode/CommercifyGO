package usecase

import (
	"errors"
	"math"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/dto"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// DashboardUseCase handles dashboard-related business logic
type DashboardUseCase struct {
	orderRepo   repository.OrderRepository
	userRepo    repository.UserRepository
	productRepo repository.ProductRepository
}

// NewDashboardUseCase creates a new DashboardUseCase
func NewDashboardUseCase(orderRepo repository.OrderRepository, userRepo repository.UserRepository, productRepo repository.ProductRepository) *DashboardUseCase {
	return &DashboardUseCase{
		orderRepo:   orderRepo,
		userRepo:    userRepo,
		productRepo: productRepo,
	}
}

// GetDashboardStats retrieves dashboard statistics for a given time period
func (d *DashboardUseCase) GetDashboardStats(request dto.DashboardStatsRequest) (*dto.DashboardStats, error) {
	// Calculate time range
	endDate := time.Now()
	var startDate time.Time

	if request.StartDate != nil && request.EndDate != nil {
		startDate = *request.StartDate
		endDate = *request.EndDate
	} else if request.Days > 0 {
		startDate = endDate.AddDate(0, 0, -request.Days)
	} else {
		// Default to 30 days if no range specified
		startDate = endDate.AddDate(0, 0, -30)
	}

	// Validate date range
	if startDate.After(endDate) {
		return nil, errors.New("start date cannot be after end date")
	}

	// Get total revenue
	totalRevenue, err := d.orderRepo.GetTotalRevenueByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get total orders
	totalOrders, err := d.orderRepo.GetTotalOrdersByDateRange(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get total customers
	totalCustomers, err := d.userRepo.GetTotalCustomersCount()
	if err != nil {
		return nil, err
	}

	// Get new customers
	newCustomers, err := d.userRepo.GetNewCustomersCount(startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Get recent orders (limit to 10 for dashboard)
	recentOrders, err := d.orderRepo.GetRecentOrdersSummary(startDate, endDate, 10)
	if err != nil {
		return nil, err
	}

	// Get top products (limit to 10 for dashboard)
	topProducts, err := d.orderRepo.GetTopProductsSummary(startDate, endDate, 10)
	if err != nil {
		return nil, err
	}

	// Get total products count
	totalProducts, err := d.productRepo.GetTotalProductsCount()
	if err != nil {
		return nil, err
	}

	// Get low stock products count (threshold of 10 or less)
	lowStockProducts, err := d.productRepo.GetLowStockProductsCount(10)
	if err != nil {
		return nil, err
	}

	// Calculate previous period for comparison
	periodDuration := endDate.Sub(startDate)
	previousStartDate := startDate.Add(-periodDuration)
	previousEndDate := startDate

	// Get previous period revenue for comparison
	revenueChange, err := d.calculatePercentageChange(
		func() (int64, error) {
			return d.orderRepo.GetTotalRevenueByDateRange(previousStartDate, previousEndDate)
		},
		func() (int64, error) { return d.orderRepo.GetTotalRevenueByDateRange(startDate, endDate) },
	)
	if err != nil {
		return nil, err
	}

	// Get previous period orders for comparison
	ordersChange, err := d.calculatePercentageChange(
		func() (int64, error) {
			return d.orderRepo.GetTotalOrdersByDateRange(previousStartDate, previousEndDate)
		},
		func() (int64, error) { return d.orderRepo.GetTotalOrdersByDateRange(startDate, endDate) },
	)
	if err != nil {
		return nil, err
	}

	return &dto.DashboardStats{
		TotalRevenue:     totalRevenue,
		TotalOrders:      totalOrders,
		TotalCustomers:   totalCustomers,
		NewCustomers:     newCustomers,
		TotalProducts:    totalProducts,
		LowStockProducts: lowStockProducts,
		RevenueChange:    revenueChange,
		OrdersChange:     ordersChange,
		RecentOrders:     recentOrders,
		TopProducts:      topProducts,
		PeriodStart:      startDate,
		PeriodEnd:        endDate,
	}, nil
}

// calculatePercentageChange calculates the percentage change between two values
func (d *DashboardUseCase) calculatePercentageChange(getPreviousValue, getCurrentValue func() (int64, error)) (*dto.PercentageChange, error) {
	previousValue, err := getPreviousValue()
	if err != nil {
		return nil, err
	}

	currentValue, err := getCurrentValue()
	if err != nil {
		return nil, err
	}

	if previousValue == 0 {
		// If previous value is 0, we can't calculate a percentage
		if currentValue > 0 {
			return &dto.PercentageChange{Value: 100.0, Direction: "up"}, nil
		}
		return &dto.PercentageChange{Value: 0.0, Direction: "stable"}, nil
	}

	// Calculate percentage change
	change := float64(currentValue-previousValue) / float64(previousValue) * 100

	direction := "stable"
	if change > 0 {
		direction = "up"
	} else if change < 0 {
		direction = "down"
	}

	return &dto.PercentageChange{
		Value:     math.Abs(change),
		Direction: direction,
	}, nil
}
