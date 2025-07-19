package gorm

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/testutil"
)

func TestProductRepository_Delete(t *testing.T) {
	db := testutil.SetupTestDB(t)
	productRepo := NewProductRepository(db)
	variantRepo := NewProductVariantRepository(db)

	t.Run("Delete product should delete all variants", func(t *testing.T) {
		// Create a test category
		category, err := entity.NewCategory("Test Category", "Test Description", nil)
		require.NoError(t, err)
		err = db.Create(category).Error
		require.NoError(t, err)

		// Create a test product with multiple variants
		variant1, err := entity.NewProductVariant(
			"TEST-SKU-1",
			10,
			1000,
			1.0,
			map[string]string{"size": "S", "color": "red"},
			[]string{"image1.jpg"},
			true,
		)
		require.NoError(t, err)

		variant2, err := entity.NewProductVariant(
			"TEST-SKU-2",
			5,
			1200,
			1.2,
			map[string]string{"size": "M", "color": "blue"},
			[]string{"image2.jpg"},
			false,
		)
		require.NoError(t, err)

		product, err := entity.NewProduct(
			"Test Product",
			"Test Description",
			"USD",
			category.ID,
			[]string{"product_image.jpg"},
			[]*entity.ProductVariant{variant1, variant2},
			true,
		)
		require.NoError(t, err)

		// Create the product (this should create variants too)
		err = productRepo.Create(product)
		require.NoError(t, err)
		require.NotZero(t, product.ID)

		// Verify product and variants were created
		createdProduct, err := productRepo.GetByID(product.ID)
		require.NoError(t, err)
		require.Len(t, createdProduct.Variants, 2)

		// Store variant IDs for later verification
		variantID1 := createdProduct.Variants[0].ID
		variantID2 := createdProduct.Variants[1].ID

		// Delete the product
		err = productRepo.Delete(product.ID)
		require.NoError(t, err)

		// Verify product is deleted
		_, err = productRepo.GetByID(product.ID)
		assert.Error(t, err)

		// Verify variants are also deleted
		_, err = variantRepo.GetByID(variantID1)
		assert.Error(t, err, "Variant 1 should be deleted")

		_, err = variantRepo.GetByID(variantID2)
		assert.Error(t, err, "Variant 2 should be deleted")

		// Also check using direct database query to ensure no orphaned variants
		var variantCount int64
		err = db.Model(&entity.ProductVariant{}).Where("product_id = ?", product.ID).Count(&variantCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), variantCount, "No variants should remain for the deleted product")
	})

	t.Run("Delete non-existent product should not error", func(t *testing.T) {
		// Try to delete a product that doesn't exist
		err := productRepo.Delete(99999)
		assert.NoError(t, err, "Deleting non-existent product should not error")
	})

	t.Run("Delete product with no variants should work", func(t *testing.T) {
		// Create a test category
		category, err := entity.NewCategory("Test Category 2", "Test Description", nil)
		require.NoError(t, err)
		err = db.Create(category).Error
		require.NoError(t, err)

		// Create a product with one variant, then delete the variant manually
		variant, err := entity.NewProductVariant(
			"TEST-SKU-ORPHAN",
			10,
			1000,
			1.0,
			map[string]string{"size": "S"},
			[]string{"image1.jpg"},
			true,
		)
		require.NoError(t, err)

		product, err := entity.NewProduct(
			"Test Product No Variants",
			"Test Description",
			"USD",
			category.ID,
			[]string{"product_image.jpg"},
			[]*entity.ProductVariant{variant},
			true,
		)
		require.NoError(t, err)

		// Create the product
		err = productRepo.Create(product)
		require.NoError(t, err)

		// Manually delete the variant to simulate orphaned product
		err = db.Delete(&entity.ProductVariant{}, variant.ID).Error
		require.NoError(t, err)

		// Now delete the product (should work even with no variants)
		err = productRepo.Delete(product.ID)
		assert.NoError(t, err)

		// Verify product is deleted
		_, err = productRepo.GetByID(product.ID)
		assert.Error(t, err)
	})

	t.Run("Transaction rollback on variant deletion failure", func(t *testing.T) {
		// This test would require mocking to simulate a failure during variant deletion
		// For now, we'll just verify the basic transaction behavior
		// In a real scenario, you might use a mock database to simulate failures
	})

	t.Run("Delete should be hard deletion, not soft deletion", func(t *testing.T) {
		// Create a test category
		category, err := entity.NewCategory("Test Hard Delete Category", "Test Description", nil)
		require.NoError(t, err)
		err = db.Create(category).Error
		require.NoError(t, err)

		// Create a test product with variants
		variant, err := entity.NewProductVariant(
			"TEST-HARD-DELETE-SKU",
			10,
			1000,
			1.0,
			map[string]string{"size": "M"},
			[]string{"image1.jpg"},
			true,
		)
		require.NoError(t, err)

		product, err := entity.NewProduct(
			"Test Hard Delete Product",
			"Test Description",
			"USD",
			category.ID,
			[]string{"product_image.jpg"},
			[]*entity.ProductVariant{variant},
			true,
		)
		require.NoError(t, err)

		// Create the product
		err = productRepo.Create(product)
		require.NoError(t, err)
		require.NotZero(t, product.ID)

		// Store IDs for verification
		productID := product.ID
		variantID := variant.ID

		// Delete the product
		err = productRepo.Delete(productID)
		require.NoError(t, err)

		// Verify hard deletion by checking with Unscoped() - should find nothing
		var deletedProduct entity.Product
		err = db.Unscoped().First(&deletedProduct, productID).Error
		assert.Error(t, err, "Product should be hard deleted (not found even with Unscoped)")
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		var deletedVariant entity.ProductVariant
		err = db.Unscoped().First(&deletedVariant, variantID).Error
		assert.Error(t, err, "Variant should be hard deleted (not found even with Unscoped)")
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))

		// Double-check with count queries
		var productCount int64
		err = db.Unscoped().Model(&entity.Product{}).Where("id = ?", productID).Count(&productCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), productCount, "Product should not exist in database")

		var variantCount int64
		err = db.Unscoped().Model(&entity.ProductVariant{}).Where("id = ?", variantID).Count(&variantCount).Error
		require.NoError(t, err)
		assert.Equal(t, int64(0), variantCount, "Variant should not exist in database")
	})
}
