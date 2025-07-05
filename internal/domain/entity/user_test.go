package entity

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUser(t *testing.T) {
	t.Run("NewUser success", func(t *testing.T) {
		user, err := NewUser(
			"test@example.com",
			"password123",
			"John",
			"Doe",
			RoleUser,
		)

		require.NoError(t, err)
		assert.Equal(t, "test@example.com", user.Email)
		assert.NotEmpty(t, user.Password)
		assert.NotEqual(t, "password123", user.Password) // Should be hashed
		assert.Equal(t, "John", user.FirstName)
		assert.Equal(t, "Doe", user.LastName)
		assert.Equal(t, string(RoleUser), user.Role)
	})

	t.Run("NewUser with admin role", func(t *testing.T) {
		user, err := NewUser(
			"admin@example.com",
			"adminpass",
			"Jane",
			"Admin",
			RoleAdmin,
		)

		require.NoError(t, err)
		assert.Equal(t, "admin@example.com", user.Email)
		assert.Equal(t, string(RoleAdmin), user.Role)
		assert.True(t, user.IsAdmin())
	})

	t.Run("NewUser validation errors", func(t *testing.T) {
		tests := []struct {
			name        string
			email       string
			password    string
			firstName   string
			lastName    string
			role        UserRole
			expectedErr string
		}{
			{
				name:        "empty email",
				email:       "",
				password:    "password123",
				firstName:   "John",
				lastName:    "Doe",
				role:        RoleUser,
				expectedErr: "email cannot be empty",
			},
			{
				name:        "empty password",
				email:       "test@example.com",
				password:    "",
				firstName:   "John",
				lastName:    "Doe",
				role:        RoleUser,
				expectedErr: "password cannot be empty",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				user, err := NewUser(tt.email, tt.password, tt.firstName, tt.lastName, tt.role)
				assert.Nil(t, user)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedErr)
			})
		}
	})

	t.Run("ComparePassword", func(t *testing.T) {
		user, err := NewUser("test@example.com", "password123", "John", "Doe", RoleUser)
		require.NoError(t, err)

		// Correct password
		err = user.ComparePassword("password123")
		assert.NoError(t, err)

		// Incorrect password
		err = user.ComparePassword("wrongpassword")
		assert.Error(t, err)

		// Empty password
		err = user.ComparePassword("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password cannot be empty")
	})

	t.Run("UpdatePassword", func(t *testing.T) {
		user, err := NewUser("test@example.com", "oldpassword", "John", "Doe", RoleUser)
		require.NoError(t, err)

		oldPasswordHash := user.Password

		// Update password
		err = user.UpdatePassword("newpassword123")
		assert.NoError(t, err)
		assert.NotEqual(t, oldPasswordHash, user.Password)

		// Verify new password works
		err = user.ComparePassword("newpassword123")
		assert.NoError(t, err)

		// Verify old password doesn't work
		err = user.ComparePassword("oldpassword")
		assert.Error(t, err)

		// Empty password
		err = user.UpdatePassword("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "password cannot be empty")
	})

	t.Run("Update", func(t *testing.T) {
		user, err := NewUser("test@example.com", "password123", "John", "Doe", RoleUser)
		require.NoError(t, err)

		// Valid update
		err = user.Update("Jane", "Smith")
		assert.NoError(t, err)
		assert.Equal(t, "Jane", user.FirstName)
		assert.Equal(t, "Smith", user.LastName)

		// Invalid updates
		err = user.Update("", "Smith")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "first name cannot be empty")

		err = user.Update("Jane", "")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "last name cannot be empty")
	})

	t.Run("FullName", func(t *testing.T) {
		user, err := NewUser("test@example.com", "password123", "John", "Doe", RoleUser)
		require.NoError(t, err)

		assert.Equal(t, "John Doe", user.FullName())
	})

	t.Run("IsAdmin", func(t *testing.T) {
		// Regular user
		user, err := NewUser("user@example.com", "password123", "John", "Doe", RoleUser)
		require.NoError(t, err)
		assert.False(t, user.IsAdmin())

		// Admin user
		admin, err := NewUser("admin@example.com", "password123", "Jane", "Admin", RoleAdmin)
		require.NoError(t, err)
		assert.True(t, admin.IsAdmin())
	})

	t.Run("ToUserDTO", func(t *testing.T) {
		user, err := NewUser("test@example.com", "password123", "John", "Doe", RoleUser)
		require.NoError(t, err)

		// Mock ID that would be set by GORM
		user.ID = 1

		dto := user.ToUserDTO()
		assert.Equal(t, uint(1), dto.ID)
		assert.Equal(t, "test@example.com", dto.Email)
		assert.Equal(t, "John", dto.FirstName)
		assert.Equal(t, "Doe", dto.LastName)
		assert.Equal(t, string(RoleUser), dto.Role)
	})
}

func TestUserRole(t *testing.T) {
	t.Run("UserRole constants", func(t *testing.T) {
		assert.Equal(t, UserRole("admin"), RoleAdmin)
		assert.Equal(t, UserRole("user"), RoleUser)
	})
}
