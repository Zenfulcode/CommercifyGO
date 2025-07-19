package gorm

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/testutil"
)

func TestCategoryRepository_UpdateParentID(t *testing.T) {
	// Setup
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewCategoryRepository(db)

	t.Run("should update ParentID from nil to valid ID", func(t *testing.T) {
		// Create parent category
		parentCategory, err := entity.NewCategory("Parent Category", "Parent description", nil)
		require.NoError(t, err)
		err = repo.Create(parentCategory)
		require.NoError(t, err)

		// Create child category without parent
		childCategory, err := entity.NewCategory("Child Category", "Child description", nil)
		require.NoError(t, err)
		err = repo.Create(childCategory)
		require.NoError(t, err)

		// Verify initial state
		assert.Nil(t, childCategory.ParentID)

		// Update child to have parent
		childCategory.ParentID = &parentCategory.ID
		err = repo.Update(childCategory)
		require.NoError(t, err)

		// Fetch from database to verify
		updated, err := repo.GetByID(childCategory.ID)
		require.NoError(t, err)
		assert.NotNil(t, updated.ParentID)
		assert.Equal(t, parentCategory.ID, *updated.ParentID)
	})

	t.Run("should update ParentID from valid ID to nil using 0", func(t *testing.T) {
		// Create parent category
		parentCategory, err := entity.NewCategory("Parent Category 2", "Parent description", nil)
		require.NoError(t, err)
		err = repo.Create(parentCategory)
		require.NoError(t, err)

		// Create child category with parent
		childCategory, err := entity.NewCategory("Child Category 2", "Child description", &parentCategory.ID)
		require.NoError(t, err)
		err = repo.Create(childCategory)
		require.NoError(t, err)

		// Verify initial state
		assert.NotNil(t, childCategory.ParentID)
		assert.Equal(t, parentCategory.ID, *childCategory.ParentID)

		// Update child to remove parent (simulate sending parent_id: 0 from API)
		childCategory.ParentID = nil
		err = repo.Update(childCategory)
		require.NoError(t, err)

		// Fetch from database to verify
		updated, err := repo.GetByID(childCategory.ID)
		require.NoError(t, err)
		assert.Nil(t, updated.ParentID)
	})

	t.Run("should update ParentID from one valid ID to another", func(t *testing.T) {
		// Create parent categories
		parentCategory1, err := entity.NewCategory("Parent Category 3", "Parent description", nil)
		require.NoError(t, err)
		err = repo.Create(parentCategory1)
		require.NoError(t, err)

		parentCategory2, err := entity.NewCategory("Parent Category 4", "Parent description", nil)
		require.NoError(t, err)
		err = repo.Create(parentCategory2)
		require.NoError(t, err)

		// Create child category with first parent
		childCategory, err := entity.NewCategory("Child Category 3", "Child description", &parentCategory1.ID)
		require.NoError(t, err)
		err = repo.Create(childCategory)
		require.NoError(t, err)

		// Verify initial state
		assert.NotNil(t, childCategory.ParentID)
		assert.Equal(t, parentCategory1.ID, *childCategory.ParentID)

		// Update child to have second parent
		childCategory.ParentID = &parentCategory2.ID
		err = repo.Update(childCategory)
		require.NoError(t, err)

		// Fetch from database to verify
		updated, err := repo.GetByID(childCategory.ID)
		require.NoError(t, err)
		assert.NotNil(t, updated.ParentID)
		assert.Equal(t, parentCategory2.ID, *updated.ParentID)
	})
}

func TestCategoryRepository_UpdateParentIDToNil(t *testing.T) {
	// Setup
	db := testutil.SetupTestDB(t)
	defer testutil.CleanupTestDB(t, db)

	repo := NewCategoryRepository(db)

	t.Run("should update ParentID to nil when explicitly set", func(t *testing.T) {
		// Create parent category
		parentCategory, err := entity.NewCategory("Parent Category", "Parent description", nil)
		require.NoError(t, err)
		err = repo.Create(parentCategory)
		require.NoError(t, err)

		// Create child category with parent
		childCategory, err := entity.NewCategory("Child Category", "Child description", &parentCategory.ID)
		require.NoError(t, err)
		err = repo.Create(childCategory)
		require.NoError(t, err)

		// Verify initial state
		initial, err := repo.GetByID(childCategory.ID)
		require.NoError(t, err)
		assert.NotNil(t, initial.ParentID)
		assert.Equal(t, parentCategory.ID, *initial.ParentID)

		// Update child to remove parent by setting ParentID to nil
		childCategory.ParentID = nil
		t.Logf("Before update: childCategory.ParentID = %v", childCategory.ParentID)

		err = repo.Update(childCategory)
		require.NoError(t, err)

		// Verify the update worked by fetching from database
		updated, err := repo.GetByID(childCategory.ID)
		require.NoError(t, err)
		t.Logf("After update: updated.ParentID = %v", updated.ParentID)
		assert.Nil(t, updated.ParentID, "ParentID should be nil after update")
	})
}
