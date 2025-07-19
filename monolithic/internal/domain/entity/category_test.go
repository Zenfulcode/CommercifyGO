package entity

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCategory(t *testing.T) {
	t.Run("NewCategory success", func(t *testing.T) {
		parentID := uint(1)

		category, err := NewCategory(
			"Electronics",
			"Electronic devices and accessories",
			&parentID,
		)

		require.NoError(t, err)
		assert.Equal(t, "Electronics", category.Name)
		assert.Equal(t, "Electronic devices and accessories", category.Description)
		assert.NotNil(t, category.ParentID)
		assert.Equal(t, uint(1), *category.ParentID)
	})

	t.Run("NewCategory success - no parent", func(t *testing.T) {
		category, err := NewCategory(
			"Root Category",
			"Top level category",
			nil,
		)

		require.NoError(t, err)
		assert.Equal(t, "Root Category", category.Name)
		assert.Equal(t, "Top level category", category.Description)
		assert.Nil(t, category.ParentID)
	})

	t.Run("NewCategory validation errors", func(t *testing.T) {
		tests := []struct {
			name         string
			categoryName string
			description  string
			parentID     *uint
			expectedErr  string
		}{
			{
				name:         "empty name",
				categoryName: "",
				description:  "Description",
				parentID:     nil,
				expectedErr:  "category name cannot be empty",
			},
			{
				name:         "name too long",
				categoryName: strings.Repeat("a", 256),
				description:  "Description",
				parentID:     nil,
				expectedErr:  "category name cannot exceed 255 characters",
			},
			{
				name:         "zero parent ID",
				categoryName: "Electronics",
				description:  "Description",
				parentID:     func() *uint { id := uint(0); return &id }(),
				expectedErr:  "parent ID cannot be zero",
			},
			{
				name:         "description too long",
				categoryName: "Electronics",
				description:  strings.Repeat("a", 65536),
				parentID:     nil,
				expectedErr:  "category description cannot exceed 65535 characters",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				category, err := NewCategory(tt.categoryName, tt.description, tt.parentID)
				assert.Nil(t, category)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			})
		}
	})

	t.Run("ToCategoryDTO", func(t *testing.T) {
		parentID := uint(2)
		category, err := NewCategory("Test Category", "Test description", &parentID)
		require.NoError(t, err)

		// Mock some fields that would be set by GORM
		category.ID = 1

		dto := category.ToCategoryDTO()
		assert.Equal(t, uint(1), dto.ID)
		assert.Equal(t, "Test Category", dto.Name)
		assert.Equal(t, "Test description", dto.Description)
		assert.NotNil(t, dto.ParentID)
		assert.Equal(t, uint(2), *dto.ParentID)
		assert.NotNil(t, dto.CreatedAt)
		assert.NotNil(t, dto.UpdatedAt)
	})
}
