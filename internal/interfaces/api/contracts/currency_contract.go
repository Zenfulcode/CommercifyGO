package contracts

import (
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/dto"
)

// CreateCurrencyRequest represents a request to create a new currency
type CreateCurrencyRequest struct {
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	Symbol       string  `json:"symbol"`
	ExchangeRate float64 `json:"exchange_rate"`
	IsEnabled    bool    `json:"is_enabled"`
	IsDefault    bool    `json:"is_default,omitempty"`
}

// UpdateCurrencyRequest represents a request to update an existing currency
type UpdateCurrencyRequest struct {
	Name         string  `json:"name,omitempty"`
	Symbol       string  `json:"symbol,omitempty"`
	ExchangeRate float64 `json:"exchange_rate,omitempty"`
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

// CreateConvertAmountResponse creates a ConvertAmountResponse from conversion data
func CreateConvertAmountResponse(fromCurrency string, fromAmount float64, toCurrency string, toAmountCents int64) ConvertAmountResponse {
	fromCents := money.ToCents(fromAmount)

	return ConvertAmountResponse{
		From: createConvertedAmountDTO(fromCurrency, fromCents),
		To:   createConvertedAmountDTO(toCurrency, toAmountCents),
	}
}

// CreateListCurrenciesResponse creates a response for listing currencies
func CreateCurrenciesListResponse(currencies []*entity.Currency, page, pageSize, total int) ListResponseDTO[dto.CurrencyDTO] {
	var currencyDTOs []dto.CurrencyDTO
	for _, currency := range currencies {
		currencyDTOs = append(currencyDTOs, *currency.ToCurrencyDTO())
	}

	if len(currencyDTOs) == 0 {
		return ListResponseDTO[dto.CurrencyDTO]{
			Success:    true,
			Data:       []dto.CurrencyDTO{},
			Pagination: PaginationDTO{Page: page, PageSize: pageSize, Total: 0},
			Message:    "No currencies found",
		}
	}

	return ListResponseDTO[dto.CurrencyDTO]{
		Success: true,
		Data:    currencyDTOs,
		Pagination: PaginationDTO{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
		Message: "Currencies retrieved successfully",
	}
}

func CreateCurrencyResponse(currency *dto.CurrencyDTO) ResponseDTO[dto.CurrencyDTO] {
	return SuccessResponse(*currency)
}

// CreateDeleteCurrencyResponse creates a standard delete response
func CreateDeleteCurrencyResponse() ResponseDTO[DeleteCurrencyResponse] {
	return SuccessResponse(DeleteCurrencyResponse{
		Status:  "success",
		Message: "Currency deleted successfully",
	})
}

func createConvertedAmountDTO(currency string, amountCents int64) ConvertedAmountDTO {
	return ConvertedAmountDTO{
		Currency: currency,
		Amount:   money.FromCents(amountCents),
		Cents:    amountCents,
	}
}
