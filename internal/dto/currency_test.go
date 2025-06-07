package dto

import (
	"testing"
	"time"

	"github.com/zenfulcode/commercify/internal/domain/entity"
	"github.com/zenfulcode/commercify/internal/domain/money"
)

func TestFromCurrencyEntity(t *testing.T) {
	now := time.Now()
	currency := &entity.Currency{
		Code:         "USD",
		Name:         "US Dollar",
		Symbol:       "$",
		ExchangeRate: 1.0,
		IsEnabled:    true,
		IsDefault:    true,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	dto := toCurrencyDTO(currency)

	if dto.Code != currency.Code {
		t.Errorf("Expected Code %s, got %s", currency.Code, dto.Code)
	}
	if dto.Name != currency.Name {
		t.Errorf("Expected Name %s, got %s", currency.Name, dto.Name)
	}
	if dto.Symbol != currency.Symbol {
		t.Errorf("Expected Symbol %s, got %s", currency.Symbol, dto.Symbol)
	}
	if dto.ExchangeRate != currency.ExchangeRate {
		t.Errorf("Expected ExchangeRate %f, got %f", currency.ExchangeRate, dto.ExchangeRate)
	}
	if dto.IsEnabled != currency.IsEnabled {
		t.Errorf("Expected IsEnabled %t, got %t", currency.IsEnabled, dto.IsEnabled)
	}
	if dto.IsDefault != currency.IsDefault {
		t.Errorf("Expected IsDefault %t, got %t", currency.IsDefault, dto.IsDefault)
	}
	if !dto.CreatedAt.Equal(currency.CreatedAt) {
		t.Errorf("Expected CreatedAt %v, got %v", currency.CreatedAt, dto.CreatedAt)
	}
	if !dto.UpdatedAt.Equal(currency.UpdatedAt) {
		t.Errorf("Expected UpdatedAt %v, got %v", currency.UpdatedAt, dto.UpdatedAt)
	}
}

func TestFromCurrencyEntityDetail(t *testing.T) {
	now := time.Now()
	currency := &entity.Currency{
		Code:         "EUR",
		Name:         "Euro",
		Symbol:       "€",
		ExchangeRate: 0.85,
		IsEnabled:    true,
		IsDefault:    false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	dto := toCurrencyDTO(currency)

	if dto.Code != currency.Code {
		t.Errorf("Expected Code %s, got %s", currency.Code, dto.Code)
	}
	if dto.Name != currency.Name {
		t.Errorf("Expected Name %s, got %s", currency.Name, dto.Name)
	}
	if dto.Symbol != currency.Symbol {
		t.Errorf("Expected Symbol %s, got %s", currency.Symbol, dto.Symbol)
	}
	if dto.ExchangeRate != currency.ExchangeRate {
		t.Errorf("Expected ExchangeRate %f, got %f", currency.ExchangeRate, dto.ExchangeRate)
	}
	if dto.IsEnabled != currency.IsEnabled {
		t.Errorf("Expected IsEnabled %t, got %t", currency.IsEnabled, dto.IsEnabled)
	}
	if dto.IsDefault != currency.IsDefault {
		t.Errorf("Expected IsDefault %t, got %t", currency.IsDefault, dto.IsDefault)
	}
}

func TestFromCurrencyEntitySummary(t *testing.T) {
	currency := &entity.Currency{
		Code:         "GBP",
		Name:         "British Pound",
		Symbol:       "£",
		ExchangeRate: 0.76,
		IsEnabled:    true,
		IsDefault:    false,
	}

	dto := toCurrencyDTO(currency)

	if dto.Code != currency.Code {
		t.Errorf("Expected Code %s, got %s", currency.Code, dto.Code)
	}
	if dto.Name != currency.Name {
		t.Errorf("Expected Name %s, got %s", currency.Name, dto.Name)
	}
	if dto.Symbol != currency.Symbol {
		t.Errorf("Expected Symbol %s, got %s", currency.Symbol, dto.Symbol)
	}
	if dto.ExchangeRate != currency.ExchangeRate {
		t.Errorf("Expected ExchangeRate %f, got %f", currency.ExchangeRate, dto.ExchangeRate)
	}
	if dto.IsDefault != currency.IsDefault {
		t.Errorf("Expected IsDefault %t, got %t", currency.IsDefault, dto.IsDefault)
	}
}

func TestFromCurrencyEntities(t *testing.T) {
	currencies := []*entity.Currency{
		{
			Code:         "USD",
			Name:         "US Dollar",
			Symbol:       "$",
			ExchangeRate: 1.0,
			IsEnabled:    true,
			IsDefault:    true,
		},
		{
			Code:         "EUR",
			Name:         "Euro",
			Symbol:       "€",
			ExchangeRate: 0.85,
			IsEnabled:    true,
			IsDefault:    false,
		},
	}

	dtos := toCurrencyDTOList(currencies)

	if len(dtos) != len(currencies) {
		t.Errorf("Expected %d DTOs, got %d", len(currencies), len(dtos))
	}

	for i, dto := range dtos {
		if dto.Code != currencies[i].Code {
			t.Errorf("Expected Code %s, got %s", currencies[i].Code, dto.Code)
		}
		if dto.Name != currencies[i].Name {
			t.Errorf("Expected Name %s, got %s", currencies[i].Name, dto.Name)
		}
	}
}

func TestFromCurrencyEntitiesSummary(t *testing.T) {
	currencies := []*entity.Currency{
		{
			Code:         "USD",
			Name:         "US Dollar",
			Symbol:       "$",
			ExchangeRate: 1.0,
			IsDefault:    true,
		},
		{
			Code:         "EUR",
			Name:         "Euro",
			Symbol:       "€",
			ExchangeRate: 0.85,
			IsDefault:    false,
		},
	}

	dtos := toCurrencySummaryDTOList(currencies)

	if len(dtos) != len(currencies) {
		t.Errorf("Expected %d DTOs, got %d", len(currencies), len(dtos))
	}

	for i, dto := range dtos {
		if dto.Code != currencies[i].Code {
			t.Errorf("Expected Code %s, got %s", currencies[i].Code, dto.Code)
		}
		if dto.Name != currencies[i].Name {
			t.Errorf("Expected Name %s, got %s", currencies[i].Name, dto.Name)
		}
	}
}

func TestCreateCurrencyRequestToUseCaseInput(t *testing.T) {
	request := CreateCurrencyRequest{
		Code:         "CAD",
		Name:         "Canadian Dollar",
		Symbol:       "C$",
		ExchangeRate: 1.25,
		IsEnabled:    true,
		IsDefault:    false,
	}

	input := request.ToUseCaseInput()

	if input.Code != request.Code {
		t.Errorf("Expected Code %s, got %s", request.Code, input.Code)
	}
	if input.Name != request.Name {
		t.Errorf("Expected Name %s, got %s", request.Name, input.Name)
	}
	if input.Symbol != request.Symbol {
		t.Errorf("Expected Symbol %s, got %s", request.Symbol, input.Symbol)
	}
	if input.ExchangeRate != request.ExchangeRate {
		t.Errorf("Expected ExchangeRate %f, got %f", request.ExchangeRate, input.ExchangeRate)
	}
	if input.IsEnabled != request.IsEnabled {
		t.Errorf("Expected IsEnabled %t, got %t", request.IsEnabled, input.IsEnabled)
	}
	if input.IsDefault != request.IsDefault {
		t.Errorf("Expected IsDefault %t, got %t", request.IsDefault, input.IsDefault)
	}
}

func TestUpdateCurrencyRequestToUseCaseInput(t *testing.T) {
	isEnabled := true
	isDefault := false

	request := UpdateCurrencyRequest{
		Name:         "Updated Dollar",
		Symbol:       "$$$",
		ExchangeRate: 1.1,
		IsEnabled:    &isEnabled,
		IsDefault:    &isDefault,
	}

	input := request.ToUseCaseInput()

	if input.Name != request.Name {
		t.Errorf("Expected Name %s, got %s", request.Name, input.Name)
	}
	if input.Symbol != request.Symbol {
		t.Errorf("Expected Symbol %s, got %s", request.Symbol, input.Symbol)
	}
	if input.ExchangeRate != request.ExchangeRate {
		t.Errorf("Expected ExchangeRate %f, got %f", request.ExchangeRate, input.ExchangeRate)
	}
	if input.IsEnabled != *request.IsEnabled {
		t.Errorf("Expected IsEnabled %t, got %t", *request.IsEnabled, input.IsEnabled)
	}
	if input.IsDefault != *request.IsDefault {
		t.Errorf("Expected IsDefault %t, got %t", *request.IsDefault, input.IsDefault)
	}
}

func TestUpdateCurrencyRequestToUseCaseInputWithNilValues(t *testing.T) {
	request := UpdateCurrencyRequest{
		Name:         "Updated Dollar",
		Symbol:       "$$$",
		ExchangeRate: 1.1,
		IsEnabled:    nil,
		IsDefault:    nil,
	}

	input := request.ToUseCaseInput()

	if input.Name != request.Name {
		t.Errorf("Expected Name %s, got %s", request.Name, input.Name)
	}
	if input.Symbol != request.Symbol {
		t.Errorf("Expected Symbol %s, got %s", request.Symbol, input.Symbol)
	}
	if input.ExchangeRate != request.ExchangeRate {
		t.Errorf("Expected ExchangeRate %f, got %f", request.ExchangeRate, input.ExchangeRate)
	}
	// IsEnabled and IsDefault should have default values when nil
	if input.IsEnabled != false {
		t.Errorf("Expected IsEnabled false (default), got %t", input.IsEnabled)
	}
	if input.IsDefault != false {
		t.Errorf("Expected IsDefault false (default), got %t", input.IsDefault)
	}
}

func TestCreateConvertedAmountDTO(t *testing.T) {
	currency := "USD"
	amountCents := int64(12345) // $123.45

	dto := createConvertedAmountDTO(currency, amountCents)

	if dto.Currency != currency {
		t.Errorf("Expected Currency %s, got %s", currency, dto.Currency)
	}
	if dto.Cents != amountCents {
		t.Errorf("Expected Cents %d, got %d", amountCents, dto.Cents)
	}
	expectedAmount := money.FromCents(amountCents)
	if dto.Amount != expectedAmount {
		t.Errorf("Expected Amount %f, got %f", expectedAmount, dto.Amount)
	}
}

func TestCreateConvertAmountResponse(t *testing.T) {
	fromCurrency := "USD"
	fromAmount := 100.0
	toCurrency := "EUR"
	toAmountCents := int64(8500) // 85.00 EUR

	response := CreateConvertAmountResponse(fromCurrency, fromAmount, toCurrency, toAmountCents)

	// Test from currency
	if response.From.Currency != fromCurrency {
		t.Errorf("Expected From.Currency %s, got %s", fromCurrency, response.From.Currency)
	}
	if response.From.Amount != fromAmount {
		t.Errorf("Expected From.Amount %f, got %f", fromAmount, response.From.Amount)
	}
	expectedFromCents := money.ToCents(fromAmount)
	if response.From.Cents != expectedFromCents {
		t.Errorf("Expected From.Cents %d, got %d", expectedFromCents, response.From.Cents)
	}

	// Test to currency
	if response.To.Currency != toCurrency {
		t.Errorf("Expected To.Currency %s, got %s", toCurrency, response.To.Currency)
	}
	if response.To.Cents != toAmountCents {
		t.Errorf("Expected To.Cents %d, got %d", toAmountCents, response.To.Cents)
	}
	expectedToAmount := money.FromCents(toAmountCents)
	if response.To.Amount != expectedToAmount {
		t.Errorf("Expected To.Amount %f, got %f", expectedToAmount, response.To.Amount)
	}
}

func TestCreateListCurrenciesResponse(t *testing.T) {
	currencies := []*entity.Currency{
		{
			Code:         "USD",
			Name:         "US Dollar",
			Symbol:       "$",
			ExchangeRate: 1.0,
			IsEnabled:    true,
			IsDefault:    true,
		},
		{
			Code:         "EUR",
			Name:         "Euro",
			Symbol:       "€",
			ExchangeRate: 0.85,
			IsEnabled:    true,
			IsDefault:    false,
		},
	}

	response := CreateCurrenciesListResponse(currencies, 1, 10, len(currencies))

	if response.Pagination.Total != len(currencies) {
		t.Errorf("Expected Total %d, got %d", len(currencies), response.Pagination.Total)
	}
	if len(response.Data) != len(currencies) {
		t.Errorf("Expected %d currencies, got %d", len(currencies), len(response.Data))
	}

	for i, dto := range response.Data {
		if dto.Code != currencies[i].Code {
			t.Errorf("Expected Currency[%d].Code %s, got %s", i, currencies[i].Code, dto.Code)
		}
	}
}

func TestCreateListEnabledCurrenciesResponse(t *testing.T) {
	currencies := []*entity.Currency{
		{
			Code:         "USD",
			Name:         "US Dollar",
			Symbol:       "$",
			ExchangeRate: 1.0,
			IsDefault:    true,
		},
		{
			Code:         "EUR",
			Name:         "Euro",
			Symbol:       "€",
			ExchangeRate: 0.85,
			IsDefault:    false,
		},
	}

	response := CreateCurrenciesListResponse(currencies, 1, 10, len(currencies))

	if response.Pagination.Total != len(currencies) {
		t.Errorf("Expected Total %d, got %d", len(currencies), response.Pagination.Total)
	}
	if len(response.Data) != len(currencies) {
		t.Errorf("Expected %d currencies, got %d", len(currencies), len(response.Data))
	}

	for i, dto := range response.Data {
		if dto.Code != currencies[i].Code {
			t.Errorf("Expected Currency[%d].Code %s, got %s", i, currencies[i].Code, dto.Code)
		}
	}
}

func TestCreateDeleteCurrencyResponse(t *testing.T) {
	response := CreateDeleteCurrencyResponse()

	expectedStatus := "success"
	expectedMessage := "Currency deleted successfully"

	if response.Data.Status != expectedStatus {
		t.Errorf("Expected Status %s, got %s", expectedStatus, response.Data.Status)
	}
	if response.Data.Message != expectedMessage {
		t.Errorf("Expected Message %s, got %s", expectedMessage, response.Data.Message)
	}
}
