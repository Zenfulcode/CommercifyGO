package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/dto"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
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
		http.Error(w, "Failed to list currencies", http.StatusInternalServerError)
		return
	}

	// Convert to response DTO
	response := dto.CreateListCurrenciesResponse(currencies)

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
		http.Error(w, "Failed to list enabled currencies", http.StatusInternalServerError)
		return
	}

	// Convert to response DTO
	response := dto.CreateListEnabledCurrenciesResponse(currencies)

	// Return currencies
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetCurrency handles retrieving a currency by code
func (h *CurrencyHandler) GetCurrency(w http.ResponseWriter, r *http.Request) {
	// Get currency code from query parameter
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Currency code is required", http.StatusBadRequest)
		return
	}

	// Get currency
	currency, err := h.currencyUseCase.GetCurrency(code)
	if err != nil {
		h.logger.Error("Failed to get currency: %v", err)
		http.Error(w, "Currency not found", http.StatusNotFound)
		return
	}

	// Convert to response DTO
	response := dto.GetCurrencyResponse{
		Currency: dto.FromCurrencyEntityDetail(currency),
	}

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
		http.Error(w, "Default currency not found", http.StatusNotFound)
		return
	}

	// Convert to response DTO
	response := dto.GetDefaultCurrencyResponse{
		Currency: dto.FromCurrencyEntityDetail(currency),
	}

	// Return currency
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateCurrency handles creating a new currency (admin only)
func (h *CurrencyHandler) CreateCurrency(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.CreateCurrencyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to use case input
	input := request.ToUseCaseInput()

	// Create currency
	currency, err := h.currencyUseCase.CreateCurrency(input)
	if err != nil {
		h.logger.Error("Failed to create currency: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to response DTO
	response := dto.CreateCurrencyResponse{
		Currency: dto.FromCurrencyEntityDetail(currency),
	}

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
		http.Error(w, "Currency code is required", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request dto.UpdateCurrencyRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to use case input
	input := request.ToUseCaseInput()

	// Update currency
	currency, err := h.currencyUseCase.UpdateCurrency(code, input)
	if err != nil {
		h.logger.Error("Failed to update currency: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to response DTO
	response := dto.UpdateCurrencyResponse{
		Currency: dto.FromCurrencyEntityDetail(currency),
	}

	// Return updated currency
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteCurrency handles deleting a currency (admin only)
func (h *CurrencyHandler) DeleteCurrency(w http.ResponseWriter, r *http.Request) {
	// Get currency code from query parameter
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "Currency code is required", http.StatusBadRequest)
		return
	}

	// Ensure we're not trying to delete the default currency
	currency, err := h.currencyUseCase.GetCurrency(code)
	if err != nil {
		h.logger.Error("Failed to get currency: %v", err)
		http.Error(w, "Currency not found", http.StatusNotFound)
		return
	}

	if currency.IsDefault {
		http.Error(w, "Cannot delete the default currency", http.StatusBadRequest)
		return
	}

	// Delete currency
	err = h.currencyUseCase.DeleteCurrency(code)
	if err != nil {
		h.logger.Error("Failed to delete currency: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to response DTO
	response := dto.CreateDeleteCurrencyResponse()

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
		http.Error(w, "Currency code is required", http.StatusBadRequest)
		return
	}

	// Set as default
	err := h.currencyUseCase.SetDefaultCurrency(code)
	if err != nil {
		h.logger.Error("Failed to set default currency: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get updated currency
	currency, err := h.currencyUseCase.GetCurrency(code)
	if err != nil {
		h.logger.Error("Failed to get updated currency: %v", err)
		http.Error(w, "Currency not found", http.StatusNotFound)
		return
	}

	// Convert to response DTO
	response := dto.SetDefaultCurrencyResponse{
		Currency: dto.FromCurrencyEntityDetail(currency),
	}

	// Return updated currency
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// ConvertAmount handles converting an amount from one currency to another
func (h *CurrencyHandler) ConvertAmount(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.ConvertAmountRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if request.Amount <= 0 {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(request.FromCurrency) == "" {
		http.Error(w, "From currency is required", http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(request.ToCurrency) == "" {
		http.Error(w, "To currency is required", http.StatusBadRequest)
		return
	}

	// Convert amount
	fromCents := money.ToCents(request.Amount)
	toCents, err := h.currencyUseCase.ConvertPrice(fromCents, request.FromCurrency, request.ToCurrency)
	if err != nil {
		h.logger.Error("Failed to convert amount: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create response DTO
	response := dto.CreateConvertAmountResponse(request.FromCurrency, request.Amount, request.ToCurrency, toCents)

	// Return converted amount
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
