package gorm

import (
	"errors"
	"fmt"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
	"gorm.io/gorm"
)

// UserRepository implements repository.UserRepository using GORM
type UserRepository struct {
	db *gorm.DB
}

// Create implements repository.UserRepository.
func (u *UserRepository) Create(user *entity.User) error {
	return u.db.Create(user).Error
}

// Delete implements repository.UserRepository.
func (u *UserRepository) Delete(id uint) error {
	return u.db.Delete(&entity.User{}, id).Error
}

// GetByEmail implements repository.UserRepository.
func (u *UserRepository) GetByEmail(email string) (*entity.User, error) {
	var user entity.User
	if err := u.db.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with email %s not found", email)
		}
		return nil, fmt.Errorf("failed to fetch user by email: %w", err)
	}
	return &user, nil
}

// GetByID implements repository.UserRepository.
func (u *UserRepository) GetByID(id uint) (*entity.User, error) {
	var user entity.User
	if err := u.db.First(&user, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("user with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch user: %w", err)
	}
	return &user, nil
}

// List implements repository.UserRepository.
func (u *UserRepository) List(offset int, limit int) ([]*entity.User, error) {
	var users []*entity.User
	if err := u.db.Offset(offset).Limit(limit).Order("created_at DESC").Find(&users).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch users: %w", err)
	}
	return users, nil
}

// Update implements repository.UserRepository.
func (u *UserRepository) Update(user *entity.User) error {
	return u.db.Save(user).Error
}

// NewUserRepository creates a new GORM-based UserRepository
func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &UserRepository{db: db}
}
