package entity

import (
	"encoding/json"
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
		assert.Equal(t, images, variant.Images)
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
		assert.Equal(t, []string{"new-image.jpg"}, variant.Images)
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
}

func TestVariantAttributes(t *testing.T) {
	t.Run("Value method", func(t *testing.T) {
		attrs := VariantAttributes{
			"color": "red",
			"size":  "large",
		}

		value, err := attrs.Value()
		require.NoError(t, err)

		// Should return valid JSON
		// Note: JSON marshaling might change order, so we unmarshal to compare
		var result map[string]string
		err = json.Unmarshal(value.([]byte), &result)
		require.NoError(t, err)
		assert.Equal(t, map[string]string{"color": "red", "size": "large"}, result)
	})

	t.Run("Scan method with valid JSON", func(t *testing.T) {
		var attrs VariantAttributes
		jsonData := []byte(`{"color":"blue","size":"medium"}`)

		err := attrs.Scan(jsonData)
		require.NoError(t, err)
		assert.Equal(t, "blue", attrs["color"])
		assert.Equal(t, "medium", attrs["size"])
	})

	t.Run("Scan method with nil value", func(t *testing.T) {
		var attrs VariantAttributes
		err := attrs.Scan(nil)
		require.NoError(t, err)
		assert.NotNil(t, attrs)
		assert.Empty(t, attrs)
	})

	t.Run("Scan method with invalid type", func(t *testing.T) {
		var attrs VariantAttributes
		err := attrs.Scan("invalid")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "type assertion to []byte failed")
	})

	t.Run("Scan method with invalid JSON", func(t *testing.T) {
		var attrs VariantAttributes
		invalidJSON := []byte(`{"invalid json"}`)
		err := attrs.Scan(invalidJSON)
		assert.Error(t, err)
	})
}
