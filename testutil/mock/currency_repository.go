package mock

import (
	"fmt"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/domain/repository"
)

type MockCurrencyRepository struct {
	currencies      map[string]*entity.Currency
	defaultCurrency *entity.Currency
}

func NewMockCurrencyRepository() repository.CurrencyRepository {
	currencies := make(map[string]*entity.Currency)
	currencies["USD"] = &entity.Currency{
		Code:         "USD",
		Name:         "US Dollar",
		Symbol:       "$",
		ExchangeRate: 1.0,
		IsEnabled:    true,
		IsDefault:    true,
	}
	return &MockCurrencyRepository{
		currencies:      currencies,
		defaultCurrency: currencies["USD"],
	}
}

func (r *MockCurrencyRepository) Create(currency *entity.Currency) error {
	if _, exists := r.currencies[currency.Code]; exists {
		return fmt.Errorf("currency with code %s already exists", currency.Code)
	}
	r.currencies[currency.Code] = currency
	return nil
}

func (r *MockCurrencyRepository) Update(currency *entity.Currency) error {
	if _, exists := r.currencies[currency.Code]; !exists {
		return fmt.Errorf("currency with code %s does not exist", currency.Code)
	}
	r.currencies[currency.Code] = currency
	return nil
}

func (r *MockCurrencyRepository) Delete(code string) error {
	if _, exists := r.currencies[code]; !exists {
		return fmt.Errorf("currency with code %s does not exist", code)
	}
	delete(r.currencies, code)
	return nil
}
func (r *MockCurrencyRepository) GetByCode(code string) (*entity.Currency, error) {
	if code == "" {
		return r.defaultCurrency, nil
	}

	if currency, exists := r.currencies[code]; exists {
		return currency, nil
	}

	return nil, fmt.Errorf("currency with code %s does not exist", code)
}

func (r *MockCurrencyRepository) GetDefault() (*entity.Currency, error) {
	return r.defaultCurrency, nil
}

func (r *MockCurrencyRepository) List() ([]*entity.Currency, error) {
	var currencies []*entity.Currency
	for _, currency := range r.currencies {
		currencies = append(currencies, currency)
	}
	return currencies, nil
}
func (r *MockCurrencyRepository) ListEnabled() ([]*entity.Currency, error) {
	var currencies []*entity.Currency
	for _, currency := range r.currencies {
		if currency.IsEnabled {
			currencies = append(currencies, currency)
		}
	}
	return currencies, nil
}
func (r *MockCurrencyRepository) SetDefault(code string) error {
	if _, exists := r.currencies[code]; !exists {
		return fmt.Errorf("currency with code %s does not exist", code)
	}

	for _, currency := range r.currencies {
		currency.IsDefault = false
	}

	r.currencies[code].IsDefault = true
	r.defaultCurrency = r.currencies[code]
	return nil
}

// Product price operations
func (r *MockCurrencyRepository) GetProductPrices(productID uint) ([]entity.ProductPrice, error) {
	if productID == 0 {
		return nil, fmt.Errorf("product ID cannot be zero")
	}
	var prices []entity.ProductPrice
	for _, currency := range r.currencies {
		price := entity.ProductPrice{
			ProductID:    productID,
			CurrencyCode: currency.Code,
			Price:        money.ToCents(100.0),
		}
		prices = append(prices, price)
	}
	return prices, nil
}

// SetProductPrices(productID uint, prices []entity.ProductPrice) error
func (r *MockCurrencyRepository) DeleteProductPrice(productID uint, currencyCode string) error {
	if productID == 0 {
		return fmt.Errorf("product ID cannot be zero")
	}
	if _, exists := r.currencies[currencyCode]; !exists {
		return fmt.Errorf("currency with code %s does not exist", currencyCode)
	}
	return nil
}

// SetProductPrice(price *entity.ProductPrice) error

// Product variant price operations
func (r *MockCurrencyRepository) GetVariantPrices(variantID uint) ([]entity.ProductPrice, error) {
	if variantID == 0 {
		return nil, fmt.Errorf("variant ID cannot be zero")
	}
	var prices []entity.ProductPrice
	for _, currency := range r.currencies {
		price := entity.ProductPrice{
			VariantID:    variantID,
			CurrencyCode: currency.Code,
			Price:        money.ToCents(100.0),
		}
		prices = append(prices, price)
	}
	return prices, nil
}

// SetVariantPrices(variantID uint, prices []entity.ProductPrice) error
// SetVariantPrice(prices *entity.ProductPrice) error
func (r *MockCurrencyRepository) DeleteVariantPrice(variantID uint, currencyCode string) error {
	if variantID == 0 {
		return fmt.Errorf("variant ID cannot be zero")
	}
	if _, exists := r.currencies[currencyCode]; !exists {
		return fmt.Errorf("currency with code %s does not exist", currencyCode)
	}
	return nil
}
