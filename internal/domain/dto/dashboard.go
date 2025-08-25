package dto

import (
	"time"
)

// DashboardStatsRequest represents a request for dashboard statistics
type DashboardStatsRequest struct {
	StartDate *time.Time `json:"start_date,omitempty" form:"start_date"`
	EndDate   *time.Time `json:"end_date,omitempty" form:"end_date"`
	Days      int        `json:"days,omitempty" form:"days"` // Alternative to date range, defaults to 30
}

// PercentageChange represents a percentage change with value and direction
type PercentageChange struct {
	Value     float64 `json:"value"`     // percentage change (e.g., 15.5 for +15.5%)
	Direction string  `json:"direction"` // "up", "down", or "stable"
}

// DashboardStats represents aggregated dashboard statistics
type DashboardStats struct {
	TotalRevenue     int64                `json:"total_revenue"` // in cents
	TotalOrders      int64                `json:"total_orders"`
	TotalCustomers   int64                `json:"total_customers"`
	NewCustomers     int64                `json:"new_customers"`
	TotalProducts    int64                `json:"total_products"`
	LowStockProducts int64                `json:"low_stock_products"`
	RevenueChange    *PercentageChange    `json:"revenue_change"` // vs previous period
	OrdersChange     *PercentageChange    `json:"orders_change"`  // vs previous period
	RecentOrders     []RecentOrderSummary `json:"recent_orders"`
	TopProducts      []TopProductSummary  `json:"top_products"`
	PeriodStart      time.Time            `json:"period_start"`
	PeriodEnd        time.Time            `json:"period_end"`
}

// RecentOrderSummary represents a summary of recent orders for dashboard
type RecentOrderSummary struct {
	ID            uint      `json:"id"`
	OrderNumber   string    `json:"order_number"`
	CustomerName  string    `json:"customer_name"`
	CustomerEmail string    `json:"customer_email"`
	TotalAmount   int64     `json:"total_amount"` // in cents
	Status        string    `json:"status"`
	CreatedAt     time.Time `json:"created_at"`
}

// TopProductSummary represents top selling products for dashboard
type TopProductSummary struct {
	ProductID    uint   `json:"product_id"`
	ProductName  string `json:"product_name"`
	VariantID    *uint  `json:"variant_id,omitempty"`
	VariantName  string `json:"variant_name,omitempty"`
	QuantitySold int64  `json:"quantity_sold"`
	Revenue      int64  `json:"revenue"` // in cents
}
