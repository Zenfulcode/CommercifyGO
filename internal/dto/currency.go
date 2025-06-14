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
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	ExchangeRate float64 `json:"exchange_rate"`
	IsEnabled    bool    `json:"is_enabled"`
	IsDefault    bool    `json:"is_default"`
}

// UpdateCurrencyRequest represents a request to update an existing currency
type UpdateCurrencyRequest struct {
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	ExchangeRate float64 `json:"exchange_rate"`
	IsEnabled    *bool   `json:"is_enabled,omitempty"`
	IsDefault    *bool   `json:"is_default,omitempty"`
}

// ConvertAmountRequest represents a request to convert an amount between currencies
type ConvertAmountRequest struct {
	Amount       float64 `json:"amount"`
	FromCurrency string  `json:"from_currency"`
	ToCurrency   string  `json:"to_currency"`
}

// SetDefaultCurrencyRequest represents a request to set a currency as default
type SetDefaultCurrencyRequest struct {
	Code string `json:"code"`
}

// =================================================================================================
// RESPONSE DTOs
// =================================================================================================

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

// toCurrencyDTO converts a Currency entity to CurrencyDTO
func toCurrencyDTO(currency *entity.Currency) CurrencyDTO {
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

// FromCurrencyEntitySummary converts a Currency entity to CurrencySummaryDTO
func toCurrencySummary(currency *entity.Currency) CurrencySummaryDTO {
	return CurrencySummaryDTO{
		Code:         currency.Code,
		Name:         currency.Name,
		Symbol:       currency.Symbol,
		ExchangeRate: currency.ExchangeRate,
		IsDefault:    currency.IsDefault,
	}
}

// toCurrencyDTOList converts a slice of Currency entities to CurrencyDTOs
func toCurrencyDTOList(currencies []*entity.Currency) []CurrencyDTO {
	dtos := make([]CurrencyDTO, len(currencies))
	for i, currency := range currencies {
		dtos[i] = toCurrencyDTO(currency)
	}
	return dtos
}

// toCurrencySummaryDTOList converts a slice of Currency entities to CurrencySummaryDTOs
func toCurrencySummaryDTOList(currencies []*entity.Currency) []CurrencySummaryDTO {
	dtos := make([]CurrencySummaryDTO, len(currencies))
	for i, currency := range currencies {
		dtos[i] = toCurrencySummary(currency)
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
func createConvertedAmountDTO(currency string, amountCents int64) ConvertedAmountDTO {
	return ConvertedAmountDTO{
		Currency: currency,
		Amount:   money.FromCents(amountCents),
		Cents:    amountCents,
	}
}

// =================================================================================================
// UTILITY FUNCTIONS
// =================================================================================================

// CreateConvertAmountResponse creates a ConvertAmountResponse from conversion data
func CreateConvertAmountResponse(fromCurrency string, fromAmount float64, toCurrency string, toAmountCents int64) ConvertAmountResponse {
	fromCents := money.ToCents(fromAmount)

	return ConvertAmountResponse{
		From: createConvertedAmountDTO(fromCurrency, fromCents),
		To:   createConvertedAmountDTO(toCurrency, toAmountCents),
	}
}

// CreateListCurrenciesResponse creates a response for listing currencies
func CreateCurrenciesListResponse(currencies []*entity.Currency, page, pageSize, total int) ListResponseDTO[CurrencyDTO] {
	dtos := toCurrencyDTOList(currencies)
	return ListResponseDTO[CurrencyDTO]{
		Success: true,
		Data:    dtos,
		Pagination: PaginationDTO{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	}
}

func CreateCurrencySummaryResponse(currencies []*entity.Currency, page, size, total int) ListResponseDTO[CurrencySummaryDTO] {
	dtos := toCurrencySummaryDTOList(currencies)
	return ListResponseDTO[CurrencySummaryDTO]{
		Success: true,
		Data:    dtos,
		Pagination: PaginationDTO{
			Page:     page,
			PageSize: size,
			Total:    total,
		},
	}
}

func CreateCurrencyResponse(currency *entity.Currency) ResponseDTO[CurrencyDTO] {
	return SuccessResponse(toCurrencyDTO(currency))
}

// CreateDeleteCurrencyResponse creates a standard delete response
func CreateDeleteCurrencyResponse() ResponseDTO[DeleteCurrencyResponse] {
	return SuccessResponse(DeleteCurrencyResponse{
		Status:  "success",
		Message: "Currency deleted successfully",
	})
}
