package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// DiscountRepository defines the interface for discount data access
type DiscountRepository interface {
	Create(discount *entity.Discount) error
	GetByID(discountID uint) (*entity.Discount, error)
	GetByCode(code string) (*entity.Discount, error)
	Update(discount *entity.Discount) error
	Delete(discountID uint) error
	List(offset, limit int) ([]*entity.Discount, error)
	ListActive(offset, limit int) ([]*entity.Discount, error)
	IncrementUsage(discountID uint) error
}
