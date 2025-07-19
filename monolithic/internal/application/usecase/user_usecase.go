package usecase

import (
	"errors"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// UserUseCase implements user-related use cases
type UserUseCase struct {
	userRepo repository.UserRepository
}

// NewUserUseCase creates a new UserUseCase
func NewUserUseCase(userRepo repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

// RegisterInput contains the data needed to register a new user
type RegisterInput struct {
	Email     string `json:"email"`
	Password  string `json:"password"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// Register registers a new user
func (uc *UserUseCase) Register(input RegisterInput) (*entity.User, error) {
	// Check if user already exists
	existingUser, err := uc.userRepo.GetByEmail(input.Email)
	if err == nil && existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}

	// Create new user
	user, err := entity.NewUser(input.Email, input.Password, input.FirstName, input.LastName, entity.RoleUser)
	if err != nil {
		return nil, err
	}

	// Save user to repository
	if err := uc.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

// LoginInput contains the data needed for user login
type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Login authenticates a user
func (uc *UserUseCase) Login(input LoginInput) (*entity.User, error) {
	// Get user by email
	user, err := uc.userRepo.GetByEmail(input.Email)
	if err != nil {
		return nil, errors.New("invalid email or password")
	}

	// Compare password
	if err := user.ComparePassword(input.Password); err != nil {
		return nil, errors.New("invalid email or password")
	}

	return user, nil
}

// GetUserByID retrieves a user by ID
func (uc *UserUseCase) GetUserByID(id uint) (*entity.User, error) {
	return uc.userRepo.GetByID(id)
}

// UpdateUserInput contains the data needed to update a user
type UpdateUserInput struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

// UpdateUser updates a user's information
func (uc *UserUseCase) UpdateUser(id uint, input UpdateUserInput) (*entity.User, error) {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	user.Update(input.FirstName, input.LastName)

	if err := uc.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangePasswordInput contains the data needed to change a password
type ChangePasswordInput struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

// ChangePassword changes a user's password
func (uc *UserUseCase) ChangePassword(id uint, input ChangePasswordInput) error {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	// Verify current password
	if err := user.ComparePassword(input.CurrentPassword); err != nil {
		return errors.New("current password is incorrect")
	}

	// Update password
	if err := user.UpdatePassword(input.NewPassword); err != nil {
		return err
	}

	return uc.userRepo.Update(user)
}

func (uc *UserUseCase) ListUsers(offset, limit int) ([]*entity.User, error) {
	users, err := uc.userRepo.List(offset, limit)
	if err != nil {
		return nil, err
	}
	return users, nil
}
