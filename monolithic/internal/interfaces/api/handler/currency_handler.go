package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/interfaces/api/contracts"
)

// CurrencyHandler handles currency-related HTTP requests
type CurrencyHandler struct {
	currencyUseCase *usecase.CurrencyUseCase
	logger          logger.Logger
}

// NewCurrencyHandler creates a new CurrencyHandler
func NewCurrencyHandler(currencyUseCase *usecase.CurrencyUseCase, logger logger.Logger) *CurrencyHandler {
	return &CurrencyHandler{
		currencyUseCase: currencyUseCase,
		logger:          logger,
	}
}

// ListCurrencies handles listing all currencies
func (h *CurrencyHandler) ListCurrencies(w http.ResponseWriter, r *http.Request) {
	// Get currencies
	currencies, err := h.currencyUseCase.ListCurrencies()
	if err != nil {
		h.logger.Error("Failed to list currencies: %v", err)
		response := contracts.ErrorResponse("Failed to list currencies")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to response DTO
	response := contracts.CreateCurrenciesListResponse(currencies, 1, len(currencies), len(currencies))

	// Return currencies
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ListEnabledCurrencies handles listing all enabled currencies
func (h *CurrencyHandler) ListEnabledCurrencies(w http.ResponseWriter, r *http.Request) {
	// Get enabled currencies
	currencies, err := h.currencyUseCase.ListEnabledCurrencies()
	if err != nil {
		h.logger.Error("Failed to list enabled currencies: %v", err)
		response := contracts.ErrorResponse("Failed to list enabled currencies")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to response DTO
	response := contracts.CreateCurrenciesListResponse(currencies, 1, len(currencies), len(currencies))

	// Return currencies
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCurrency handles retrieving a currency by code
func (h *CurrencyHandler) GetCurrency(w http.ResponseWriter, r *http.Request) {
	// Get currency code from query parameter
	code := r.URL.Query().Get("code")
	if code == "" {
		h.logger.Error("Currency code is required")
		http.Error(w, "Currency code is required", http.StatusBadRequest)
		return
	}

	// Get currency
	currency, err := h.currencyUseCase.GetCurrency(code)
	if err != nil {
		h.logger.Error("Failed to get currency: %v", err)
		response := contracts.ErrorResponse("Currency not found")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to response DTO
	response := contracts.CreateCurrencyResponse(currency.ToCurrencyDTO())

	// Return currency
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetDefaultCurrency handles retrieving the default currency
func (h *CurrencyHandler) GetDefaultCurrency(w http.ResponseWriter, r *http.Request) {
	// Get default currency
	currency, err := h.currencyUseCase.GetDefaultCurrency()
	if err != nil {
		h.logger.Error("Failed to get default currency: %v", err)
		response := contracts.ErrorResponse("Default currency not found")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to response DTO
	response := contracts.CreateCurrencyResponse(currency.ToCurrencyDTO())

	// Return currency
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateCurrency handles creating a new currency (admin only)
func (h *CurrencyHandler) CreateCurrency(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request contracts.CreateCurrencyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to decode create currency request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to use case input
	input := request.ToUseCaseInput()

	// Create currency
	currency, err := h.currencyUseCase.CreateCurrency(input)
	if err != nil {
		h.logger.Error("Failed to create currency: %v", err)
		response := contracts.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to response DTO
	response := contracts.CreateCurrencyResponse(currency.ToCurrencyDTO())

	// Return created currency
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

// UpdateCurrency handles updating a currency (admin only)
func (h *CurrencyHandler) UpdateCurrency(w http.ResponseWriter, r *http.Request) {
	// Get currency code from query parameter
	code := r.URL.Query().Get("code")
	if code == "" {
		h.logger.Error("Currency code is required")
		http.Error(w, "Currency code is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request contracts.UpdateCurrencyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to decode update currency request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to use case input
	input := request.ToUseCaseInput()

	// Update currency
	currency, err := h.currencyUseCase.UpdateCurrency(code, input)
	if err != nil {
		h.logger.Error("Failed to update currency: %v", err)
		response := contracts.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to response DTO
	response := contracts.CreateCurrencyResponse(currency.ToCurrencyDTO())

	// Return updated currency
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteCurrency handles deleting a currency (admin only)
func (h *CurrencyHandler) DeleteCurrency(w http.ResponseWriter, r *http.Request) {
	// Get currency code from query parameter
	code := r.URL.Query().Get("code")
	if code == "" {
		h.logger.Error("Currency code is required")
		http.Error(w, "Currency code is required", http.StatusBadRequest)
		return
	}

	// Ensure we're not trying to delete the default currency
	currency, err := h.currencyUseCase.GetCurrency(code)
	if err != nil {
		h.logger.Error("Failed to get currency: %v", err)
		response := contracts.ErrorResponse("Currency not found")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	if currency.IsDefault {
		h.logger.Error("Cannot delete default currency")
		response := contracts.ErrorResponse("Cannot delete default currency")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Delete currency
	err = h.currencyUseCase.DeleteCurrency(code)
	if err != nil {
		h.logger.Error("Failed to delete currency: %v", err)
		response := contracts.ErrorResponse(err.Error())

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to response DTO
	response := contracts.CreateDeleteCurrencyResponse()

	// Return success
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

// SetDefaultCurrency handles setting a currency as the default (admin only)
func (h *CurrencyHandler) SetDefaultCurrency(w http.ResponseWriter, r *http.Request) {
	// Get currency code from query parameter
	code := r.URL.Query().Get("code")
	if code == "" {
		h.logger.Error("Currency code is required")
		http.Error(w, "Currency code is required", http.StatusBadRequest)
		return
	}

	// Set as default
	err := h.currencyUseCase.SetDefaultCurrency(code)
	if err != nil {
		h.logger.Error("Failed to set default currency: %v", err)
		response := contracts.ErrorResponse(err.Error())
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Get updated currency
	currency, err := h.currencyUseCase.GetCurrency(code)
	if err != nil {
		h.logger.Error("Failed to get updated currency: %v", err)
		response := contracts.ErrorResponse("Currency not found")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert to response DTO
	response := contracts.CreateCurrencyResponse(currency.ToCurrencyDTO())

	// Return updated currency
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ConvertAmount handles converting an amount from one currency to another
func (h *CurrencyHandler) ConvertAmount(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request contracts.ConvertAmountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		h.logger.Error("Failed to decode convert amount request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if request.Amount <= 0 {
		response := contracts.ErrorResponse("Amount must be greater than zero")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if strings.TrimSpace(request.FromCurrency) == "" {
		response := contracts.ErrorResponse("From currency is required")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	if strings.TrimSpace(request.ToCurrency) == "" {
		response := contracts.ErrorResponse("To currency is required")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Convert amount
	fromCents := money.ToCents(request.Amount)
	toCents, err := h.currencyUseCase.ConvertPrice(fromCents, request.FromCurrency, request.ToCurrency)
	if err != nil {
		h.logger.Error("Failed to convert amount: %v", err)
		response := contracts.ErrorResponse("Failed to convert amount")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Create response DTO
	response := contracts.CreateConvertAmountResponse(request.FromCurrency, request.Amount, request.ToCurrency, toCents)

	// Return converted amount
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
