package entity

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDiscount(t *testing.T) {
	startDate := time.Now()
	endDate := startDate.Add(24 * time.Hour)

	t.Run("NewDiscount success - basket discount", func(t *testing.T) {
		discount, err := NewDiscount(
			"SUMMER20",
			DiscountTypeBasket,
			DiscountMethodPercentage,
			20.0,
			5000, // $50 minimum
			1000, // $10 max discount
			nil,
			nil,
			startDate,
			endDate,
			100,
		)

		require.NoError(t, err)
		assert.Equal(t, "SUMMER20", discount.Code)
		assert.Equal(t, DiscountTypeBasket, discount.Type)
		assert.Equal(t, DiscountMethodPercentage, discount.Method)
		assert.Equal(t, 20.0, discount.Value)
		assert.Equal(t, int64(5000), discount.MinOrderValue)
		assert.Equal(t, int64(1000), discount.MaxDiscountValue)
		assert.Equal(t, startDate, discount.StartDate)
		assert.Equal(t, endDate, discount.EndDate)
		assert.Equal(t, 100, discount.UsageLimit)
		assert.Equal(t, 0, discount.CurrentUsage)
		assert.True(t, discount.Active)
	})

	t.Run("NewDiscount success - product discount", func(t *testing.T) {
		productIDs := []uint{1, 2, 3}
		categoryIDs := []uint{1}

		discount, err := NewDiscount(
			"PROD10",
			DiscountTypeProduct,
			DiscountMethodFixed,
			500, // $5 fixed discount
			0,
			0,
			productIDs,
			categoryIDs,
			startDate,
			endDate,
			50,
		)

		require.NoError(t, err)
		assert.Equal(t, DiscountTypeProduct, discount.Type)
		assert.Equal(t, DiscountMethodFixed, discount.Method)
		assert.Equal(t, productIDs, discount.ProductIDs)
		assert.Equal(t, categoryIDs, discount.CategoryIDs)
	})

	t.Run("NewDiscount validation errors", func(t *testing.T) {
		tests := []struct {
			name          string
			code          string
			discountType  DiscountType
			method        DiscountMethod
			value         float64
			productIDs    []uint
			categoryIDs   []uint
			startDate     time.Time
			endDate       time.Time
			expectedError string
		}{
			{
				name:          "empty code",
				code:          "",
				discountType:  DiscountTypeBasket,
				method:        DiscountMethodPercentage,
				value:         20.0,
				startDate:     startDate,
				endDate:       endDate,
				expectedError: "discount code cannot be empty",
			},
			{
				name:          "zero value",
				code:          "TEST",
				discountType:  DiscountTypeBasket,
				method:        DiscountMethodPercentage,
				value:         0,
				startDate:     startDate,
				endDate:       endDate,
				expectedError: "discount value must be greater than zero",
			},
			{
				name:          "negative value",
				code:          "TEST",
				discountType:  DiscountTypeBasket,
				method:        DiscountMethodPercentage,
				value:         -10,
				startDate:     startDate,
				endDate:       endDate,
				expectedError: "discount value must be greater than zero",
			},
			{
				name:          "percentage over 100",
				code:          "TEST",
				discountType:  DiscountTypeBasket,
				method:        DiscountMethodPercentage,
				value:         150,
				startDate:     startDate,
				endDate:       endDate,
				expectedError: "percentage discount cannot exceed 100%",
			},
			{
				name:          "product discount without products or categories",
				code:          "TEST",
				discountType:  DiscountTypeProduct,
				method:        DiscountMethodFixed,
				value:         500,
				productIDs:    nil,
				categoryIDs:   nil,
				startDate:     startDate,
				endDate:       endDate,
				expectedError: "product discount must specify at least one product or category",
			},
			{
				name:          "end date before start date",
				code:          "TEST",
				discountType:  DiscountTypeBasket,
				method:        DiscountMethodPercentage,
				value:         20,
				startDate:     endDate,
				endDate:       startDate,
				expectedError: "end date cannot be before start date",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				discount, err := NewDiscount(
					tt.code,
					tt.discountType,
					tt.method,
					tt.value,
					0,
					0,
					tt.productIDs,
					tt.categoryIDs,
					tt.startDate,
					tt.endDate,
					0,
				)

				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, discount)
			})
		}
	})

	t.Run("IsValid", func(t *testing.T) {
		// Create a valid discount
		discount, err := NewDiscount(
			"TEST20",
			DiscountTypeBasket,
			DiscountMethodPercentage,
			20.0,
			0,
			0,
			nil,
			nil,
			startDate,
			endDate,
			100,
		)
		require.NoError(t, err)

		// Test valid discount
		assert.True(t, discount.IsValid())

		// Test inactive discount
		discount.Active = false
		assert.False(t, discount.IsValid())

		// Test expired discount
		discount.Active = true
		discount.EndDate = time.Now().Add(-1 * time.Hour) // 1 hour ago
		assert.False(t, discount.IsValid())

		// Test not yet started discount
		discount.StartDate = time.Now().Add(1 * time.Hour) // 1 hour from now
		discount.EndDate = time.Now().Add(2 * time.Hour)   // 2 hours from now
		assert.False(t, discount.IsValid())

		// Test usage limit exceeded
		discount.StartDate = startDate
		discount.EndDate = endDate
		discount.CurrentUsage = 100
		discount.UsageLimit = 100
		assert.False(t, discount.IsValid())

		// Test unlimited usage (0 means no limit)
		discount.UsageLimit = 0
		assert.True(t, discount.IsValid())
	})

	t.Run("IsApplicableToOrder", func(t *testing.T) {
		discount, err := NewDiscount(
			"TEST20",
			DiscountTypeBasket,
			DiscountMethodPercentage,
			20.0,
			5000, // $50 minimum
			0,
			nil,
			nil,
			startDate,
			endDate,
			100,
		)
		require.NoError(t, err)

		// Create a test order
		order := &Order{
			TotalAmount: 10000, // $100 order
		}

		// Test valid usage
		assert.True(t, discount.IsApplicableToOrder(order))

		// Test below minimum order value
		order.TotalAmount = 3000 // $30 order
		assert.False(t, discount.IsApplicableToOrder(order))

		// Test inactive discount
		order.TotalAmount = 10000
		discount.Active = false
		assert.False(t, discount.IsApplicableToOrder(order))

		// Test expired discount
		discount.Active = true
		discount.EndDate = time.Now().Add(-1 * time.Hour)
		assert.False(t, discount.IsApplicableToOrder(order))
	})

	t.Run("CalculateDiscount", func(t *testing.T) {
		t.Run("percentage discount", func(t *testing.T) {
			discount, err := NewDiscount(
				"PERCENT20",
				DiscountTypeBasket,
				DiscountMethodPercentage,
				20.0,
				0,
				1000, // $10 max discount
				nil,
				nil,
				startDate,
				endDate,
				100,
			)
			require.NoError(t, err)

			// Create test order
			order := &Order{
				TotalAmount: 5000, // $50 order
			}

			// Test normal percentage calculation
			amount := discount.CalculateDiscount(order)
			assert.Equal(t, int64(1000), amount) // 20% = $10, capped at max

			// Test without max discount
			discount.MaxDiscountValue = 0
			order.TotalAmount = 10000 // $100 order
			amount = discount.CalculateDiscount(order)
			assert.Equal(t, int64(2000), amount) // 20% = $20
		})

		t.Run("fixed discount", func(t *testing.T) {
			discount, err := NewDiscount(
				"FIXED500",
				DiscountTypeBasket,
				DiscountMethodFixed,
				5.0, // $5 fixed discount (in dollars)
				0,
				0,
				nil,
				nil,
				startDate,
				endDate,
				100,
			)
			require.NoError(t, err)

			order := &Order{
				TotalAmount: 10000, // $100 order
			}
			amount := discount.CalculateDiscount(order)
			assert.Equal(t, int64(500), amount) // Fixed $5

			order.TotalAmount = 300 // $3 order
			amount = discount.CalculateDiscount(order)
			assert.Equal(t, int64(300), amount) // Capped at order amount
		})

		t.Run("product discount", func(t *testing.T) {
			discount, err := NewDiscount(
				"PROD10",
				DiscountTypeProduct,
				DiscountMethodPercentage,
				10.0,
				0,
				0,
				[]uint{1, 2}, // Products 1 and 2
				nil,
				startDate,
				endDate,
				100,
			)
			require.NoError(t, err)

			// Create order with eligible and non-eligible products
			order := &Order{
				TotalAmount: 15000,
				Items: []OrderItem{
					{ProductID: 1, Subtotal: 5000}, // Eligible - $50
					{ProductID: 2, Subtotal: 3000}, // Eligible - $30
					{ProductID: 3, Subtotal: 7000}, // Not eligible - $70
				},
			}

			amount := discount.CalculateDiscount(order)
			// 10% of (5000 + 3000) = 10% of 8000 = 800
			assert.Equal(t, int64(800), amount)
		})
	})

	t.Run("IncrementUsage", func(t *testing.T) {
		discount, err := NewDiscount(
			"TEST20",
			DiscountTypeBasket,
			DiscountMethodPercentage,
			20.0,
			0,
			0,
			nil,
			nil,
			startDate,
			endDate,
			100,
		)
		require.NoError(t, err)

		assert.Equal(t, 0, discount.CurrentUsage)

		discount.IncrementUsage()
		assert.Equal(t, 1, discount.CurrentUsage)

		discount.IncrementUsage()
		assert.Equal(t, 2, discount.CurrentUsage)
	})

	t.Run("ToDiscountDTO", func(t *testing.T) {
		startDate := time.Now()
		endDate := startDate.Add(30 * 24 * time.Hour)

		discount, err := NewDiscount(
			"SUMMER2025",
			DiscountTypeBasket,
			DiscountMethodPercentage,
			15.0,
			5000,  // 50.00 dollars in cents
			10000, // 100.00 dollars in cents
			[]uint{1, 2},
			[]uint{3, 4},
			startDate,
			endDate,
			500,
		)
		require.NoError(t, err)

		// Mock ID that would be set by GORM
		discount.ID = 1
		discount.CurrentUsage = 25

		dto := discount.ToDiscountDTO()
		assert.Equal(t, uint(1), dto.ID)
		assert.Equal(t, "SUMMER2025", dto.Code)
		assert.Equal(t, string(DiscountTypeBasket), dto.Type)
		assert.Equal(t, string(DiscountMethodPercentage), dto.Method)
		assert.Equal(t, 15.0, dto.Value)
		assert.Equal(t, 50.0, dto.MinOrderValue)     // FromCents(5000) = 50.0
		assert.Equal(t, 100.0, dto.MaxDiscountValue) // FromCents(10000) = 100.0
		assert.Equal(t, []uint{1, 2}, dto.ProductIDs)
		assert.Equal(t, []uint{3, 4}, dto.CategoryIDs)
		assert.Equal(t, startDate, dto.StartDate)
		assert.Equal(t, endDate, dto.EndDate)
		assert.Equal(t, 500, dto.UsageLimit)
		assert.Equal(t, 25, dto.CurrentUsage)
		assert.True(t, dto.Active)
	})

}
