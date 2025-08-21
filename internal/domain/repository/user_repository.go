package repository

import (
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
)

// UserRepository defines the interface for user data access
type UserRepository interface {
	Create(user *entity.User) error
	GetByID(id uint) (*entity.User, error)
	GetByEmail(email string) (*entity.User, error)
	Update(user *entity.User) error
	Delete(id uint) error
	List(offset, limit int) ([]*entity.User, error)

	// Dashboard statistics methods
	GetTotalCustomersCount() (int64, error)
	GetNewCustomersCount(startDate, endDate time.Time) (int64, error)
}
