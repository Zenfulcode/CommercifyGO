package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShippingMethodDTOConversions(t *testing.T) {
	t.Run("ToShippingMethodDTO", func(t *testing.T) {
		shippingMethod, err := NewShippingMethod("Standard Delivery", "Reliable standard delivery", 5)
		require.NoError(t, err)

		// Mock ID that would be set by GORM
		shippingMethod.ID = 1

		dto := shippingMethod.ToShippingMethodDTO()
		assert.Equal(t, uint(1), dto.ID)
		assert.Equal(t, "Standard Delivery", dto.Name)
		assert.Equal(t, "Reliable standard delivery", dto.Description)
		assert.Equal(t, 5, dto.EstimatedDeliveryDays)
		assert.True(t, dto.Active)
		assert.NotNil(t, dto.CreatedAt)
		assert.NotNil(t, dto.UpdatedAt)
	})
}
