package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id uint) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id uint) error
	List(offset, limit int) ([]*entity.User, error)
}
