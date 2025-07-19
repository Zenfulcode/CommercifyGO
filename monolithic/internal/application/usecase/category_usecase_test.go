package usecase

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/infrastructure/repository/gorm"
	"github.com/zenfulcode/commercify/testutil"
)

func TestCategoryUseCase_DeleteCategory_WithProducts(t *testing.T) {
	t.Run("should prevent deletion of category with products", func(t *testing.T) {
		// Setup
		db := testutil.SetupTestDB(t)
		defer testutil.CleanupTestDB(t, db)

		categoryRepo := gorm.NewCategoryRepository(db)
		productRepo := gorm.NewProductRepository(db)
		categoryUseCase := NewCategoryUseCase(categoryRepo, productRepo)

		// Create test category
		category, err := categoryUseCase.CreateCategory(CreateCategory{
			Name:        "Test Category",
			Description: "Test category for deletion test",
			ParentID:    nil,
		})
		require.NoError(t, err)

		// Create a product in this category
		variant, err := entity.NewProductVariant(
			"TEST-SKU-001",
			10,
			9999, // 99.99 in cents
			1.0,
			map[string]string{"size": "M"},
			[]string{"test-image.jpg"},
			true,
		)
		require.NoError(t, err)

		product, err := entity.NewProduct(
			"Test Product",
			"Test product description",
			"USD",
			category.ID,
			[]string{"product-image.jpg"},
			[]*entity.ProductVariant{variant},
			true,
		)
		require.NoError(t, err)

		err = productRepo.Create(product)
		require.NoError(t, err)

		// Act: Try to delete category with products
		err = categoryUseCase.DeleteCategory(category.ID)

		// Assert: Should fail with specific error message
		assert.Error(t, err)
		assert.Equal(t, "cannot delete category with products", err.Error())

		// Verify category still exists
		existingCategory, err := categoryRepo.GetByID(category.ID)
		assert.NoError(t, err)
		assert.NotNil(t, existingCategory)
	})

	t.Run("should allow deletion of category without products", func(t *testing.T) {
		// Setup
		db := testutil.SetupTestDB(t)
		defer testutil.CleanupTestDB(t, db)

		categoryRepo := gorm.NewCategoryRepository(db)
		productRepo := gorm.NewProductRepository(db)
		categoryUseCase := NewCategoryUseCase(categoryRepo, productRepo)

		// Create test category
		category, err := categoryUseCase.CreateCategory(CreateCategory{
			Name:        "Empty Category",
			Description: "Category without products",
			ParentID:    nil,
		})
		require.NoError(t, err)

		// Act: Delete category without products
		err = categoryUseCase.DeleteCategory(category.ID)

		// Assert: Should succeed
		assert.NoError(t, err)

		// Verify category was deleted
		_, err = categoryRepo.GetByID(category.ID)
		assert.Error(t, err)
	})

	t.Run("should still prevent deletion of category with child categories", func(t *testing.T) {
		// Setup
		db := testutil.SetupTestDB(t)
		defer testutil.CleanupTestDB(t, db)

		categoryRepo := gorm.NewCategoryRepository(db)
		productRepo := gorm.NewProductRepository(db)
		categoryUseCase := NewCategoryUseCase(categoryRepo, productRepo)

		// Create parent category
		parentCategory, err := categoryUseCase.CreateCategory(CreateCategory{
			Name:        "Parent Category",
			Description: "Parent category",
			ParentID:    nil,
		})
		require.NoError(t, err)

		// Create child category
		_, err = categoryUseCase.CreateCategory(CreateCategory{
			Name:        "Child Category",
			Description: "Child category",
			ParentID:    &parentCategory.ID,
		})
		require.NoError(t, err)

		// Act: Try to delete parent category with children
		err = categoryUseCase.DeleteCategory(parentCategory.ID)

		// Assert: Should fail with specific error message
		assert.Error(t, err)
		assert.Equal(t, "cannot delete category with child categories", err.Error())
	})
}

func TestCategoryUseCase_UpdateCategory_ParentIDZero(t *testing.T) {
	t.Run("should remove parent when parent_id is 0", func(t *testing.T) {
		// Setup
		db := testutil.SetupTestDB(t)
		defer testutil.CleanupTestDB(t, db)

		categoryRepo := gorm.NewCategoryRepository(db)
		productRepo := gorm.NewProductRepository(db)
		categoryUseCase := NewCategoryUseCase(categoryRepo, productRepo)

		// Create parent category
		parentCategory, err := categoryUseCase.CreateCategory(CreateCategory{
			Name:        "Parent Category",
			Description: "Parent category for test",
			ParentID:    nil,
		})
		require.NoError(t, err)

		// Create child category with parent
		childCategory, err := categoryUseCase.CreateCategory(CreateCategory{
			Name:        "Child Category",
			Description: "Child category for test",
			ParentID:    &parentCategory.ID,
		})
		require.NoError(t, err)

		// Verify initial state
		assert.NotNil(t, childCategory.ParentID)
		assert.Equal(t, parentCategory.ID, *childCategory.ParentID)

		// Update child to remove parent using parent_id: 0
		zeroID := uint(0)
		_, err = categoryUseCase.UpdateCategory(UpdateCategory{
			CategoryID: childCategory.ID,
			ParentID:   &zeroID,
		})
		require.NoError(t, err)

		// Verify in database directly (this is the most reliable test)
		fetchedCategory, err := categoryRepo.GetByID(childCategory.ID)
		require.NoError(t, err)
		assert.Nil(t, fetchedCategory.ParentID, "ParentID should be nil after setting to 0")
	})
}
