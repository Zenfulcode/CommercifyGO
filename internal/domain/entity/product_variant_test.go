package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProductVariant(t *testing.T) {
	t.Run("NewProductVariant success", func(t *testing.T) {
		attributes := VariantAttributes{
			"color": "red",
			"size":  "large",
		}
		images := []string{"image1.jpg", "image2.jpg"}

		variant, err := NewProductVariant(
			"TEST-SKU-001",
			10,
			9999,
			1.5,
			attributes,
			images,
			true,
		)

		require.NoError(t, err)
		assert.Equal(t, "TEST-SKU-001", variant.SKU)
		assert.Equal(t, 10, variant.Stock)
		assert.Equal(t, int64(9999), variant.Price)
		assert.Equal(t, 1.5, variant.Weight)
		assert.Equal(t, attributes, variant.Attributes)
		assert.Equal(t, images, []string(variant.Images))
		assert.True(t, variant.IsDefault)
	})

	t.Run("NewProductVariant validation errors", func(t *testing.T) {
		tests := []struct {
			name          string
			sku           string
			stock         int
			price         int64
			weight        float64
			expectedError string
		}{
			{
				name:          "empty SKU",
				sku:           "",
				stock:         10,
				price:         9999,
				weight:        1.5,
				expectedError: "SKU cannot be empty",
			},
			{
				name:          "negative stock",
				sku:           "TEST-SKU",
				stock:         -1,
				price:         9999,
				weight:        1.5,
				expectedError: "stock cannot be negative",
			},
			{
				name:          "negative price",
				sku:           "TEST-SKU",
				stock:         10,
				price:         -1,
				weight:        1.5,
				expectedError: "price cannot be negative",
			},
			{
				name:          "negative weight",
				sku:           "TEST-SKU",
				stock:         10,
				price:         9999,
				weight:        -1.0,
				expectedError: "weight cannot be negative",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				variant, err := NewProductVariant(
					tt.sku,
					tt.stock,
					tt.price,
					tt.weight,
					nil,
					nil,
					false,
				)

				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, variant)
			})
		}
	})

	t.Run("NewProductVariant with nil attributes", func(t *testing.T) {
		variant, err := NewProductVariant(
			"TEST-SKU-001",
			10,
			9999,
			1.5,
			nil,
			nil,
			false,
		)

		require.NoError(t, err)
		assert.NotNil(t, variant.Attributes)
		assert.Empty(t, variant.Attributes)
	})

	t.Run("Update method", func(t *testing.T) {
		variant, err := NewProductVariant(
			"TEST-SKU-001",
			10,
			9999,
			1.5,
			VariantAttributes{"color": "red"},
			[]string{"image1.jpg"},
			false,
		)
		require.NoError(t, err)

		// Test successful update
		updated, err := variant.Update(
			"NEW-SKU",
			20,
			19999,
			2.5,
			[]string{"new-image.jpg"},
			VariantAttributes{"color": "blue"},
		)

		require.NoError(t, err)
		assert.True(t, updated)
		assert.Equal(t, "NEW-SKU", variant.SKU)
		assert.Equal(t, 20, variant.Stock)
		assert.Equal(t, int64(19999), variant.Price)
		assert.Equal(t, 2.5, variant.Weight)
		assert.Equal(t, []string{"new-image.jpg"}, []string(variant.Images))
		assert.Equal(t, VariantAttributes{"color": "blue"}, variant.Attributes)
	})

	t.Run("Update method with no changes", func(t *testing.T) {
		variant, err := NewProductVariant(
			"TEST-SKU-001",
			10,
			9999,
			1.5,
			nil,
			nil,
			false,
		)
		require.NoError(t, err)

		// Test no update
		updated, err := variant.Update(
			"TEST-SKU-001", // same SKU
			10,             // same stock
			9999,           // same price
			1.5,            // same weight
			nil,            // same images
			nil,            // same attributes
		)

		require.NoError(t, err)
		assert.False(t, updated)
	})

	t.Run("ToVariantDTO", func(t *testing.T) {
		attributes := VariantAttributes{
			"color": "red",
			"size":  "large",
		}

		variant, err := NewProductVariant("SKU-001", 10, 9999, 1.5, attributes, nil, true)
		require.NoError(t, err)

		// Mock IDs that would be set by GORM
		variant.ID = 1
		variant.ProductID = 2

		dto := variant.ToVariantDTO()
		assert.Equal(t, uint(1), dto.ID)
		assert.Equal(t, uint(2), dto.ProductID)
		assert.Equal(t, "SKU-001", dto.SKU)
		assert.Equal(t, 10, dto.Stock)
		assert.Equal(t, float64(99.99), dto.Price) // Converted from cents to dollars
		assert.Equal(t, 1.5, dto.Weight)
		assert.True(t, dto.IsDefault)
		assert.NotNil(t, dto.Attributes)
		assert.Equal(t, "red", dto.Attributes["color"])
		assert.Equal(t, "large", dto.Attributes["size"])
		// VariantName is generated from attributes, check it contains both values
		assert.Contains(t, dto.VariantName, "red")
		assert.Contains(t, dto.VariantName, "large")
		assert.Contains(t, dto.VariantName, " / ")
	})

	t.Run("ToVariantDTO_VariantName", func(t *testing.T) {
		// Test that VariantName is properly generated from attributes
		attributes := VariantAttributes{
			"color": "blue",
			"size":  "medium",
		}

		variant, err := NewProductVariant("SKU-002", 5, 1999, 0.8, attributes, nil, false)
		require.NoError(t, err)

		variant.ID = 10
		variant.ProductID = 20

		dto := variant.ToVariantDTO()

		// VariantName should contain both attribute values separated by " / "
		// Order may vary due to map iteration, so check both possibilities
		expectedName1 := "blue / medium"
		expectedName2 := "medium / blue"

		actualName := dto.VariantName
		isValidName := actualName == expectedName1 || actualName == expectedName2
		assert.True(t, isValidName, "VariantName should be '%s' or '%s', got '%s'", expectedName1, expectedName2, actualName)

		// Also verify it contains the expected components
		assert.Contains(t, actualName, "blue")
		assert.Contains(t, actualName, "medium")
		assert.Contains(t, actualName, " / ")
	})
}
