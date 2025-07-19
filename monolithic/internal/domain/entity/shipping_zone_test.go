package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShippingZoneDTOConversions(t *testing.T) {
	t.Run("ToShippingZoneDTO", func(t *testing.T) {
		countries := []string{"US", "CA", "MX"}
		shippingZone, err := NewShippingZone("North America", "North American shipping zone", countries)
		require.NoError(t, err)

		// Mock ID that would be set by GORM
		shippingZone.ID = 1

		dto := shippingZone.ToShippingZoneDTO()
		assert.Equal(t, uint(1), dto.ID)
		assert.Equal(t, "North America", dto.Name)
		assert.Equal(t, "North American shipping zone", dto.Description)
		assert.Equal(t, []string{"US", "CA", "MX"}, dto.Countries)
		assert.True(t, dto.Active)
		assert.NotNil(t, dto.CreatedAt)
		assert.NotNil(t, dto.UpdatedAt)
	})
}
