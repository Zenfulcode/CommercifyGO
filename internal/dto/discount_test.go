package dto

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
)

func TestConvertToDiscountDTO(t *testing.T) {
	t.Run("Convert discount entity to DTO successfully", func(t *testing.T) {
		// Create a test discount entity
		discount, err := entity.NewDiscount(
			"TEST10",
			entity.DiscountTypeBasket,
			entity.DiscountMethodPercentage,
			10.0,
			money.ToCents(50.0), // MinOrderValue
			money.ToCents(30.0), // MaxDiscountValue
			[]uint{1, 2},
			[]uint{3, 4},
			time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC),
			time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC),
			100,
		)
		assert.NoError(t, err)
		discount.ID = 1

		// Convert to DTO
		dto := toDiscountDTO(discount)

		// Assert all fields are correctly converted
		assert.Equal(t, uint(1), dto.ID)
		assert.Equal(t, "TEST10", dto.Code)
		assert.Equal(t, "basket", dto.Type)
		assert.Equal(t, "percentage", dto.Method)
		assert.Equal(t, 10.0, dto.Value)
		assert.Equal(t, 50.0, dto.MinOrderValue)
		assert.Equal(t, 30.0, dto.MaxDiscountValue)
		assert.Equal(t, []uint{1, 2}, dto.ProductIDs)
		assert.Equal(t, []uint{3, 4}, dto.CategoryIDs)
		assert.Equal(t, time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC), dto.StartDate)
		assert.Equal(t, time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC), dto.EndDate)
		assert.Equal(t, 100, dto.UsageLimit)
		assert.Equal(t, 0, dto.CurrentUsage)
		assert.True(t, dto.Active)
	})

	t.Run("Convert nil discount returns empty DTO", func(t *testing.T) {
		dto := toDiscountDTO(nil)
		assert.Equal(t, DiscountDTO{}, dto)
	})
}

func TestConvertToAppliedDiscountDTO(t *testing.T) {
	t.Run("Convert applied discount entity to DTO successfully", func(t *testing.T) {
		appliedDiscount := &entity.AppliedDiscount{
			DiscountID:     1,
			DiscountCode:   "TEST10",
			DiscountAmount: money.ToCents(15.0),
		}

		dto := ConvertToAppliedDiscountDTO(appliedDiscount)

		assert.Equal(t, uint(1), dto.ID)
		assert.Equal(t, "TEST10", dto.Code)
		assert.Equal(t, 15.0, dto.Amount)
		// Type, Method, Value are empty as noted in the conversion function
		assert.Equal(t, "", dto.Type)
		assert.Equal(t, "", dto.Method)
		assert.Equal(t, 0.0, dto.Value)
	})

	t.Run("Convert nil applied discount returns empty DTO", func(t *testing.T) {
		dto := ConvertToAppliedDiscountDTO(nil)
		assert.Equal(t, AppliedDiscountDTO{}, dto)
	})
}

func TestConvertDiscountListToDTO(t *testing.T) {
	t.Run("Convert list of discounts to DTOs", func(t *testing.T) {
		// Create test discounts
		discount1, _ := entity.NewDiscount(
			"FIRST10",
			entity.DiscountTypeBasket,
			entity.DiscountMethodPercentage,
			10.0, 0, 0, []uint{}, []uint{},
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			0,
		)
		discount1.ID = 1

		discount2, _ := entity.NewDiscount(
			"SECOND20",
			entity.DiscountTypeProduct,
			entity.DiscountMethodFixed,
			20.0, 0, 0, []uint{1}, []uint{},
			time.Now().Add(-24*time.Hour),
			time.Now().Add(30*24*time.Hour),
			0,
		)
		discount2.ID = 2

		discounts := []*entity.Discount{discount1, discount2}

		// Convert to DTOs
		dtos := ConvertDiscountListToDTO(discounts)

		// Assert
		assert.Len(t, dtos, 2)
		assert.Equal(t, "FIRST10", dtos[0].Code)
		assert.Equal(t, "basket", dtos[0].Type)
		assert.Equal(t, "SECOND20", dtos[1].Code)
		assert.Equal(t, "product", dtos[1].Type)
	})

	t.Run("Convert empty list returns empty slice", func(t *testing.T) {
		dtos := ConvertDiscountListToDTO([]*entity.Discount{})
		assert.Empty(t, dtos)
	})
}
