package repository

import "github.com/zenfulcode/commercify/internal/domain/entity"

// CurrencyRepository defines the contract for currency operations
type CurrencyRepository interface {
	// Currency operations
	Create(currency *entity.Currency) error
	Update(currency *entity.Currency) error
	Delete(code string) error
	GetByCode(code string) (*entity.Currency, error)
	GetDefault() (*entity.Currency, error)
	List() ([]*entity.Currency, error)
	ListEnabled() ([]*entity.Currency, error)
	SetDefault(code string) error
}
