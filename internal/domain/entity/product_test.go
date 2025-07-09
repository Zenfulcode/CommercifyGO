package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProduct(t *testing.T) {
	t.Run("NewProduct success", func(t *testing.T) {
		variant, err := NewProductVariant("SKU-001", 10, 9999, 1.5, nil, nil, true)
		require.NoError(t, err)

		variants := []*ProductVariant{variant}
		images := []string{"image1.jpg", "image2.jpg"}

		product, err := NewProduct(
			"Test Product",
			"A test product description",
			"USD",
			1,
			images,
			variants,
			true,
		)

		require.NoError(t, err)
		assert.Equal(t, "Test Product", product.Name)
		assert.Equal(t, "A test product description", product.Description)
		assert.Equal(t, "USD", product.Currency)
		assert.Equal(t, uint(1), product.CategoryID)
		assert.Equal(t, images, []string(product.Images))
		assert.True(t, product.Active)
		assert.NotNil(t, product.Variants)
		assert.Len(t, product.Variants, 1) // One variant was provided in constructor
		assert.Equal(t, "SKU-001", product.Variants[0].SKU)
	})

	t.Run("NewProduct validation errors", func(t *testing.T) {
		variant, err := NewProductVariant("SKU-001", 10, 9999, 1.5, nil, nil, true)
		require.NoError(t, err)

		tests := []struct {
			name          string
			productName   string
			categoryID    uint
			variants      []*ProductVariant
			expectedError string
		}{
			{
				name:          "empty name",
				productName:   "",
				categoryID:    1,
				variants:      []*ProductVariant{variant},
				expectedError: "product name cannot be empty",
			},
			{
				name:          "zero category ID",
				productName:   "Test Product",
				categoryID:    0,
				variants:      []*ProductVariant{variant},
				expectedError: "category ID cannot be zero",
			},
			{
				name:          "no variants",
				productName:   "Test Product",
				categoryID:    1,
				variants:      []*ProductVariant{},
				expectedError: "at least one variant must be provided",
			},
			{
				name:          "nil variants",
				productName:   "Test Product",
				categoryID:    1,
				variants:      nil,
				expectedError: "at least one variant must be provided",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				product, err := NewProduct(
					tt.productName,
					"Description",
					"USD",
					tt.categoryID,
					nil,
					tt.variants,
					true,
				)

				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err.Error())
				assert.Nil(t, product)
			})
		}
	})

	t.Run("AddVariant", func(t *testing.T) {
		// Create initial product with one variant
		variant1, err := NewProductVariant("SKU-001", 10, 9999, 1.5, nil, nil, true)
		require.NoError(t, err)

		product, err := NewProduct(
			"Test Product",
			"Description",
			"USD",
			1,
			[]string{},
			[]*ProductVariant{variant1},
			true,
		)
		require.NoError(t, err)

		// Add another variant
		variant2, err := NewProductVariant("SKU-002", 5, 19999, 2.0, nil, nil, false)
		require.NoError(t, err)

		product.AddVariant(variant2)
		assert.Len(t, product.Variants, 2) // Started with 1 variant, added 1 more
		assert.Equal(t, "SKU-002", product.Variants[1].SKU)
	})

	t.Run("RemoveVariant", func(t *testing.T) {
		variant1, err := NewProductVariant("SKU-001", 10, 9999, 1.5, nil, nil, true)
		require.NoError(t, err)
		variant1.ID = 1

		variant2, err := NewProductVariant("SKU-002", 5, 19999, 2.0, nil, nil, false)
		require.NoError(t, err)
		variant2.ID = 2

		product, err := NewProduct(
			"Test Product",
			"Description",
			"USD",
			1,
			nil,
			[]*ProductVariant{variant1},
			true,
		)
		require.NoError(t, err)

		// Add second variant
		err = product.AddVariant(variant2)
		require.NoError(t, err)

		// Remove variant
		err = product.RemoveVariant(1)
		require.NoError(t, err)
		assert.Len(t, product.Variants, 1)
		assert.Equal(t, uint(2), product.Variants[0].ID)

		// Try to remove non-existent variant
		err = product.RemoveVariant(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "variant with ID 999 not found")
	})

	t.Run("GetVariantBySKU", func(t *testing.T) {
		variant, err := NewProductVariant("SKU-001", 10, 9999, 1.5, nil, nil, true)
		require.NoError(t, err)

		product, err := NewProduct(
			"Test Product",
			"Description",
			"USD",
			1,
			nil,
			[]*ProductVariant{variant},
			true,
		)
		require.NoError(t, err)

		// Add the variant to the product
		err = product.AddVariant(variant)
		require.NoError(t, err)

		// Get variant by SKU
		foundVariant := product.GetVariantBySKU("SKU-001")
		assert.NotNil(t, foundVariant)
		assert.Equal(t, "SKU-001", foundVariant.SKU)

		// Get non-existent variant
		notFound := product.GetVariantBySKU("NON-EXISTENT")
		assert.Nil(t, notFound)
	})

	t.Run("GetDefaultVariant", func(t *testing.T) {
		variant, err := NewProductVariant("SKU-001", 10, 9999, 1.5, nil, nil, true)
		require.NoError(t, err)

		product, err := NewProduct(
			"Test Product",
			"Description",
			"USD",
			1,
			nil,
			[]*ProductVariant{variant},
			true,
		)
		require.NoError(t, err)

		// Add the variant to the product
		err = product.AddVariant(variant)
		require.NoError(t, err)

		// Get default variant
		defaultVariant := product.GetDefaultVariant()
		assert.NotNil(t, defaultVariant)
		assert.True(t, defaultVariant.IsDefault)
		assert.Equal(t, "SKU-001", defaultVariant.SKU)
	})

	t.Run("Active status", func(t *testing.T) {
		variant, err := NewProductVariant("SKU-001", 10, 9999, 1.5, nil, nil, true)
		require.NoError(t, err)

		product, err := NewProduct(
			"Test Product",
			"Description",
			"USD",
			1,
			nil,
			[]*ProductVariant{variant},
			true,
		)
		require.NoError(t, err)

		assert.True(t, product.Active)

		// Test inactive product
		product.Active = false
		assert.False(t, product.Active)
	})

	t.Run("IsAvailable", func(t *testing.T) {
		variant1, err := NewProductVariant("SKU-001", 10, 9999, 1.5, nil, nil, true)
		require.NoError(t, err)

		variant2, err := NewProductVariant("SKU-002", 0, 19999, 2.0, nil, nil, false)
		require.NoError(t, err)

		product, err := NewProduct(
			"Test Product",
			"Description",
			"USD",
			1,
			nil,
			[]*ProductVariant{variant1},
			true,
		)
		require.NoError(t, err)

		// Add variants
		err = product.AddVariant(variant1)
		require.NoError(t, err)
		err = product.AddVariant(variant2)
		require.NoError(t, err)

		// Test availability
		assert.True(t, product.IsAvailable(5))   // variant1 has stock
		assert.True(t, product.IsAvailable(10))  // variant1 has exactly 10
		assert.False(t, product.IsAvailable(15)) // no variant has 15+ stock
	})

	t.Run("GetTotalStock", func(t *testing.T) {
		variant1, err := NewProductVariant("SKU-001", 10, 9999, 1.5, nil, nil, true)
		require.NoError(t, err)

		variant2, err := NewProductVariant("SKU-002", 5, 19999, 2.0, nil, nil, false)
		require.NoError(t, err)

		product, err := NewProduct(
			"Test Product",
			"Description",
			"USD",
			1,
			nil,
			[]*ProductVariant{variant1}, // Start with variant1
			true,
		)
		require.NoError(t, err)

		// Add variant2 only (variant1 is already included from constructor)
		err = product.AddVariant(variant2)
		require.NoError(t, err)

		assert.Equal(t, 15, product.GetTotalStock()) // 10 + 5
	})

	t.Run("GetStockForVariant", func(t *testing.T) {
		variant, err := NewProductVariant("SKU-001", 10, 9999, 1.5, nil, nil, true)
		require.NoError(t, err)
		variant.ID = 1

		product, err := NewProduct(
			"Test Product",
			"Description",
			"USD",
			1,
			nil,
			[]*ProductVariant{variant},
			true,
		)
		require.NoError(t, err)

		// Add variant
		err = product.AddVariant(variant)
		require.NoError(t, err)

		// Get stock for existing variant
		stock, err := product.GetStockForVariant(1)
		require.NoError(t, err)
		assert.Equal(t, 10, stock)

		// Get stock for non-existent variant
		_, err = product.GetStockForVariant(999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "variant with ID 999 not found")
	})

	t.Run("GetTotalWeight", func(t *testing.T) {
		variant, err := NewProductVariant("SKU-001", 10, 9999, 2.5, nil, nil, true)
		require.NoError(t, err)

		product, err := NewProduct(
			"Test Product",
			"Description",
			"USD",
			1,
			nil,
			[]*ProductVariant{variant},
			true,
		)
		require.NoError(t, err)

		// Add variant
		err = product.AddVariant(variant)
		require.NoError(t, err)

		// Test weight calculation
		assert.Equal(t, 2.5, product.GetTotalWeight(1))  // 2.5 * 1
		assert.Equal(t, 5.0, product.GetTotalWeight(2))  // 2.5 * 2
		assert.Equal(t, 0.0, product.GetTotalWeight(0))  // 0 quantity
		assert.Equal(t, 0.0, product.GetTotalWeight(-1)) // negative quantity
	})

	t.Run("HasVariants", func(t *testing.T) {
		variant, err := NewProductVariant("SKU-001", 10, 9999, 1.5, nil, nil, true)
		require.NoError(t, err)

		product, err := NewProduct(
			"Test Product",
			"Description",
			"USD",
			1,
			nil,
			[]*ProductVariant{variant},
			true,
		)
		require.NoError(t, err)

		// Initially has the variant from constructor
		assert.True(t, product.HasVariants())

		// Add another variant
		variant2, err := NewProductVariant("SKU-002", 5, 19999, 2.0, nil, nil, false)
		require.NoError(t, err)
		err = product.AddVariant(variant2)
		require.NoError(t, err)
		assert.True(t, product.HasVariants())
	})

	t.Run("Update", func(t *testing.T) {
		variant, err := NewProductVariant("SKU-001", 10, 9999, 1.5, nil, nil, true)
		require.NoError(t, err)

		product, err := NewProduct(
			"Test Product",
			"Description",
			"USD",
			1,
			[]string{"old-image.jpg"},
			[]*ProductVariant{variant},
			true,
		)
		require.NoError(t, err)

		// Test successful update
		newName := "Updated Product"
		newDescription := "Updated Description"
		newImages := []string{"new-image1.jpg", "new-image2.jpg"}
		newActive := false

		updated := product.Update(&newName, &newDescription, nil, &newImages, &newActive, nil)
		assert.True(t, updated)
		assert.Equal(t, "Updated Product", product.Name)
		assert.Equal(t, "Updated Description", product.Description)
		assert.Equal(t, []string{"new-image1.jpg", "new-image2.jpg"}, []string(product.Images))
		assert.False(t, product.Active)

		// Test no update (same values)
		updated = product.Update(&newName, &newDescription, nil, &newImages, &newActive, nil)
		assert.False(t, updated)

		// Test empty name (should not update)
		emptyName := ""
		updated = product.Update(&emptyName, nil, nil, nil, nil, nil)
		assert.False(t, updated)
		assert.Equal(t, "Updated Product", product.Name) // unchanged
	})

	t.Run("ToProductDTO", func(t *testing.T) {
		variant, err := NewProductVariant("SKU-001", 10, 9999, 1.5, nil, nil, true)
		require.NoError(t, err)

		product, err := NewProduct(
			"Test Product",
			"A test product description",
			"USD",
			1,
			[]string{"image1.jpg", "image2.jpg"},
			[]*ProductVariant{variant},
			true,
		)
		require.NoError(t, err)

		// Mock ID that would be set by GORM
		product.ID = 1
		product.CategoryID = 2

		dto := product.ToProductDTO()
		assert.Equal(t, uint(1), dto.ID)
		assert.Equal(t, "Test Product", dto.Name)
		assert.Equal(t, "A test product description", dto.Description)
		assert.Equal(t, "USD", dto.Currency)
		assert.Equal(t, uint(2), dto.CategoryID)
		assert.Equal(t, []string{"image1.jpg", "image2.jpg"}, dto.Images)
		assert.True(t, dto.Active)
		assert.Equal(t, float64(99.99), dto.Price) // Default variant price converted from cents to dollars
		assert.Equal(t, 10, dto.TotalStock)        // Total stock across all variants
		assert.NotEmpty(t, dto.Variants)
		assert.Len(t, dto.Variants, 1)
		assert.Equal(t, "SKU-001", dto.Variants[0].SKU)
	})

	t.Run("ToProductDTO_MultipleVariants", func(t *testing.T) {
		// Test with multiple variants to verify TotalStock calculation
		variant1, err := NewProductVariant("SKU-001", 10, 9999, 1.5, nil, nil, true) // default variant
		require.NoError(t, err)

		variant2, err := NewProductVariant("SKU-002", 15, 12999, 2.0, nil, nil, false)
		require.NoError(t, err)

		product, err := NewProduct(
			"Multi-Variant Product",
			"A product with multiple variants",
			"USD",
			1,
			[]string{"image1.jpg"},
			[]*ProductVariant{variant1, variant2},
			true,
		)
		require.NoError(t, err)

		// Mock ID that would be set by GORM
		product.ID = 5
		product.CategoryID = 3

		dto := product.ToProductDTO()
		assert.Equal(t, uint(5), dto.ID)
		assert.Equal(t, "Multi-Variant Product", dto.Name)
		assert.Equal(t, float64(99.99), dto.Price) // Price from default variant (variant1)
		assert.Equal(t, 25, dto.TotalStock)        // 10 + 15 = 25 total stock
		assert.True(t, dto.HasVariants)
		assert.Len(t, dto.Variants, 2)

		// Verify both variants are present
		skus := []string{dto.Variants[0].SKU, dto.Variants[1].SKU}
		assert.Contains(t, skus, "SKU-001")
		assert.Contains(t, skus, "SKU-002")
	})
}
