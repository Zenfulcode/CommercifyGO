package usecase

import (
	"errors"
	"strings"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

// CurrencyUseCase implements currency-related use cases
type CurrencyUseCase struct {
	currencyRepo repository.CurrencyRepository
}

// NewCurrencyUseCase creates a new CurrencyUseCase
func NewCurrencyUseCase(currencyRepo repository.CurrencyRepository) *CurrencyUseCase {
	return &CurrencyUseCase{
		currencyRepo: currencyRepo,
	}
}

// CurrencyInput represents input data for creating or updating a currency
type CurrencyInput struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	ExchangeRate float64 `json:"exchange_rate"`
	IsEnabled    bool    `json:"is_enabled"`
	IsDefault    bool    `json:"is_default"`
}

// CreateCurrency creates a new currency
func (uc *CurrencyUseCase) CreateCurrency(input CurrencyInput) (*entity.Currency, error) {
	// Check if currency with this code already exists
	existingCurrency, err := uc.currencyRepo.GetByCode(input.Code)
	if err == nil && existingCurrency != nil {
		return nil, errors.New("currency with this code already exists")
	}

	// Create new currency entity
	currency, err := entity.NewCurrency(
		input.Code,
		input.Name,
		input.Symbol,
		input.ExchangeRate,
		input.IsEnabled,
		input.IsDefault,
	)
	if err != nil {
		return nil, err
	}

	// Persist the currency
	err = uc.currencyRepo.Create(currency)
	if err != nil {
		return nil, err
	}

	return currency, nil
}

// UpdateCurrency updates an existing currency
func (uc *CurrencyUseCase) UpdateCurrency(code string, input CurrencyInput) (*entity.Currency, error) {
	// Convert code to uppercase for consistency
	code = strings.ToUpper(code)

	// Get the existing currency
	currency, err := uc.currencyRepo.GetByCode(code)
	if err != nil {
		return nil, err
	}

	// Update fields
	if input.Name != "" {
		currency.Name = input.Name
	}

	if input.Symbol != "" {
		currency.Symbol = input.Symbol
	}

	if input.ExchangeRate > 0 {
		if err := currency.SetExchangeRate(input.ExchangeRate); err != nil {
			return nil, err
		}
	}

	// Handle enabled/disabled state
	if currency.IsEnabled != input.IsEnabled {
		if input.IsEnabled {
			currency.Enable()
		} else {
			if err := currency.Disable(); err != nil {
				return nil, err
			}
		}
	}

	// Handle default state
	if input.IsDefault && !currency.IsDefault {
		currency.SetAsDefault()
	} else if !input.IsDefault && currency.IsDefault {
		// If this is the current default and we're trying to unset it,
		// don't allow it - require setting a different currency as default first
		return nil, errors.New("cannot unset the default currency - set another currency as default first")
	}

	// Update in repository
	err = uc.currencyRepo.Update(currency)
	if err != nil {
		return nil, err
	}

	return currency, nil
}

// DeleteCurrency deletes a currency
func (uc *CurrencyUseCase) DeleteCurrency(code string) error {
	return uc.currencyRepo.Delete(code)
}

// GetCurrency gets a currency by its code
func (uc *CurrencyUseCase) GetCurrency(code string) (*entity.Currency, error) {
	return uc.currencyRepo.GetByCode(code)
}

// GetDefaultCurrency gets the default currency
func (uc *CurrencyUseCase) GetDefaultCurrency() (*entity.Currency, error) {
	return uc.currencyRepo.GetDefault()
}

// ListCurrencies lists all currencies
func (uc *CurrencyUseCase) ListCurrencies() ([]*entity.Currency, error) {
	return uc.currencyRepo.List()
}

// ListEnabledCurrencies lists all enabled currencies
func (uc *CurrencyUseCase) ListEnabledCurrencies() ([]*entity.Currency, error) {
	return uc.currencyRepo.ListEnabled()
}

// SetDefaultCurrency sets a currency as the default currency
func (uc *CurrencyUseCase) SetDefaultCurrency(code string) error {
	return uc.currencyRepo.SetDefault(code)
}

// ConvertPrice converts a price from one currency to another
func (uc *CurrencyUseCase) ConvertPrice(amount int64, fromCurrencyCode, toCurrencyCode string) (int64, error) {
	// Get the currencies
	fromCurrency, err := uc.currencyRepo.GetByCode(fromCurrencyCode)
	if err != nil {
		return 0, err
	}

	toCurrency, err := uc.currencyRepo.GetByCode(toCurrencyCode)
	if err != nil {
		return 0, err
	}

	// Convert the amount
	return fromCurrency.ConvertAmount(amount, toCurrency), nil
}
