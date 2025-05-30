package dto

import (
	"time"

	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
)

// =================================================================================================
// CURRENCY DTOs
// =================================================================================================

// CurrencyDTO represents a currency entity
type CurrencyDTO struct {
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Symbol       string    `json:"symbol"`
	ExchangeRate float64   `json:"exchange_rate"`
	IsEnabled    bool      `json:"is_enabled"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CurrencyDetailDTO represents detailed currency information
type CurrencyDetailDTO struct {
	Code         string    `json:"code"`
	Name         string    `json:"name"`
	Symbol       string    `json:"symbol"`
	ExchangeRate float64   `json:"exchange_rate"`
	IsEnabled    bool      `json:"is_enabled"`
	IsDefault    bool      `json:"is_default"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// CurrencySummaryDTO represents a simplified currency view
type CurrencySummaryDTO struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	ExchangeRate float64 `json:"exchange_rate"`
	IsDefault    bool    `json:"is_default"`
}

// =================================================================================================
// REQUEST DTOs
// =================================================================================================

// CreateCurrencyRequest represents a request to create a new currency
type CreateCurrencyRequest struct {
	Code         string  `json:"code" validate:"required,min=3,max=3,alpha"`
	Name         string  `json:"name" validate:"required,min=1,max=100"`
	Symbol       string  `json:"symbol" validate:"required,min=1,max=10"`
	ExchangeRate float64 `json:"exchange_rate" validate:"required,gt=0"`
	IsEnabled    bool    `json:"is_enabled"`
	IsDefault    bool    `json:"is_default"`
}

// UpdateCurrencyRequest represents a request to update an existing currency
type UpdateCurrencyRequest struct {
	Name         string  `json:"name" validate:"omitempty,min=1,max=100"`
	Symbol       string  `json:"symbol" validate:"omitempty,min=1,max=10"`
	ExchangeRate float64 `json:"exchange_rate" validate:"omitempty,gt=0"`
	IsEnabled    *bool   `json:"is_enabled,omitempty"`
	IsDefault    *bool   `json:"is_default,omitempty"`
}

// ConvertAmountRequest represents a request to convert an amount between currencies
type ConvertAmountRequest struct {
	Amount       float64 `json:"amount" validate:"required,gt=0"`
	FromCurrency string  `json:"from_currency" validate:"required,min=3,max=3,alpha"`
	ToCurrency   string  `json:"to_currency" validate:"required,min=3,max=3,alpha"`
}

// SetDefaultCurrencyRequest represents a request to set a currency as default
type SetDefaultCurrencyRequest struct {
	Code string `json:"code" validate:"required,min=3,max=3,alpha"`
}

// =================================================================================================
// RESPONSE DTOs
// =================================================================================================

// CreateCurrencyResponse represents the response after creating a currency
type CreateCurrencyResponse struct {
	Currency CurrencyDetailDTO `json:"currency"`
}

// UpdateCurrencyResponse represents the response after updating a currency
type UpdateCurrencyResponse struct {
	Currency CurrencyDetailDTO `json:"currency"`
}

// GetCurrencyResponse represents the response for getting a currency
type GetCurrencyResponse struct {
	Currency CurrencyDetailDTO `json:"currency"`
}

// ListCurrenciesResponse represents the response for listing currencies
type ListCurrenciesResponse struct {
	Currencies []CurrencyDTO `json:"currencies"`
	Total      int           `json:"total"`
}

// ListEnabledCurrenciesResponse represents the response for listing enabled currencies
type ListEnabledCurrenciesResponse struct {
	Currencies []CurrencySummaryDTO `json:"currencies"`
	Total      int                  `json:"total"`
}

// GetDefaultCurrencyResponse represents the response for getting the default currency
type GetDefaultCurrencyResponse struct {
	Currency CurrencyDetailDTO `json:"currency"`
}

// SetDefaultCurrencyResponse represents the response after setting default currency
type SetDefaultCurrencyResponse struct {
	Currency CurrencyDetailDTO `json:"currency"`
}

// ConvertAmountResponse represents the response for currency conversion
type ConvertAmountResponse struct {
	From ConvertedAmountDTO `json:"from"`
	To   ConvertedAmountDTO `json:"to"`
}

// ConvertedAmountDTO represents an amount in a specific currency
type ConvertedAmountDTO struct {
	Currency string  `json:"currency"`
	Amount   float64 `json:"amount"`
	Cents    int64   `json:"cents"`
}

// DeleteCurrencyResponse represents the response after deleting a currency
type DeleteCurrencyResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

// =================================================================================================
// CONVERSION FUNCTIONS - Entity to DTO
// =================================================================================================

// FromCurrencyEntity converts a Currency entity to CurrencyDTO
func FromCurrencyEntity(currency *entity.Currency) CurrencyDTO {
	return CurrencyDTO{
		Code:         currency.Code,
		Name:         currency.Name,
		Symbol:       currency.Symbol,
		ExchangeRate: currency.ExchangeRate,
		IsEnabled:    currency.IsEnabled,
		IsDefault:    currency.IsDefault,
		CreatedAt:    currency.CreatedAt,
		UpdatedAt:    currency.UpdatedAt,
	}
}

// FromCurrencyEntityDetail converts a Currency entity to CurrencyDetailDTO
func FromCurrencyEntityDetail(currency *entity.Currency) CurrencyDetailDTO {
	return CurrencyDetailDTO{
		Code:         currency.Code,
		Name:         currency.Name,
		Symbol:       currency.Symbol,
		ExchangeRate: currency.ExchangeRate,
		IsEnabled:    currency.IsEnabled,
		IsDefault:    currency.IsDefault,
		CreatedAt:    currency.CreatedAt,
		UpdatedAt:    currency.UpdatedAt,
	}
}

// FromCurrencyEntitySummary converts a Currency entity to CurrencySummaryDTO
func FromCurrencyEntitySummary(currency *entity.Currency) CurrencySummaryDTO {
	return CurrencySummaryDTO{
		Code:         currency.Code,
		Name:         currency.Name,
		Symbol:       currency.Symbol,
		ExchangeRate: currency.ExchangeRate,
		IsDefault:    currency.IsDefault,
	}
}

// FromCurrencyEntities converts a slice of Currency entities to CurrencyDTOs
func FromCurrencyEntities(currencies []*entity.Currency) []CurrencyDTO {
	dtos := make([]CurrencyDTO, len(currencies))
	for i, currency := range currencies {
		dtos[i] = FromCurrencyEntity(currency)
	}
	return dtos
}

// FromCurrencyEntitiesSummary converts a slice of Currency entities to CurrencySummaryDTOs
func FromCurrencyEntitiesSummary(currencies []*entity.Currency) []CurrencySummaryDTO {
	dtos := make([]CurrencySummaryDTO, len(currencies))
	for i, currency := range currencies {
		dtos[i] = FromCurrencyEntitySummary(currency)
	}
	return dtos
}

// =================================================================================================
// CONVERSION FUNCTIONS - DTO to Use Case Input
// =================================================================================================

// ToUseCaseInput converts CreateCurrencyRequest to usecase.CurrencyInput
func (r CreateCurrencyRequest) ToUseCaseInput() usecase.CurrencyInput {
	return usecase.CurrencyInput{
		Code:         r.Code,
		Name:         r.Name,
		Symbol:       r.Symbol,
		ExchangeRate: r.ExchangeRate,
		IsEnabled:    r.IsEnabled,
		IsDefault:    r.IsDefault,
	}
}

// ToUseCaseInput converts UpdateCurrencyRequest to usecase.CurrencyInput
func (r UpdateCurrencyRequest) ToUseCaseInput() usecase.CurrencyInput {
	input := usecase.CurrencyInput{
		Name:         r.Name,
		Symbol:       r.Symbol,
		ExchangeRate: r.ExchangeRate,
	}

	// Handle optional boolean fields
	if r.IsEnabled != nil {
		input.IsEnabled = *r.IsEnabled
	}
	if r.IsDefault != nil {
		input.IsDefault = *r.IsDefault
	}

	return input
}

// =================================================================================================
// CONVERSION FUNCTIONS - Amount Conversion
// =================================================================================================

// CreateConvertedAmountDTO creates a ConvertedAmountDTO from currency and amount in cents
func CreateConvertedAmountDTO(currency string, amountCents int64) ConvertedAmountDTO {
	return ConvertedAmountDTO{
		Currency: currency,
		Amount:   money.FromCents(amountCents),
		Cents:    amountCents,
	}
}

// CreateConvertAmountResponse creates a ConvertAmountResponse from conversion data
func CreateConvertAmountResponse(fromCurrency string, fromAmount float64, toCurrency string, toAmountCents int64) ConvertAmountResponse {
	fromCents := money.ToCents(fromAmount)

	return ConvertAmountResponse{
		From: CreateConvertedAmountDTO(fromCurrency, fromCents),
		To:   CreateConvertedAmountDTO(toCurrency, toAmountCents),
	}
}

// =================================================================================================
// UTILITY FUNCTIONS
// =================================================================================================

// CreateListCurrenciesResponse creates a response for listing currencies
func CreateListCurrenciesResponse(currencies []*entity.Currency) ListCurrenciesResponse {
	dtos := FromCurrencyEntities(currencies)
	return ListCurrenciesResponse{
		Currencies: dtos,
		Total:      len(dtos),
	}
}

// CreateListEnabledCurrenciesResponse creates a response for listing enabled currencies
func CreateListEnabledCurrenciesResponse(currencies []*entity.Currency) ListEnabledCurrenciesResponse {
	dtos := FromCurrencyEntitiesSummary(currencies)
	return ListEnabledCurrenciesResponse{
		Currencies: dtos,
		Total:      len(dtos),
	}
}

// CreateDeleteCurrencyResponse creates a standard delete response
func CreateDeleteCurrencyResponse() DeleteCurrencyResponse {
	return DeleteCurrencyResponse{
		Status:  "success",
		Message: "Currency deleted successfully",
	}
}
