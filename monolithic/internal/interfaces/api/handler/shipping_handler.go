package handler

import (
	"encoding/json"
	"net/http"

	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/interfaces/api/contracts"
)

// ShippingHandler handles shipping-related HTTP requests
type ShippingHandler struct {
	shippingUseCase *usecase.ShippingUseCase
	logger          logger.Logger
}

// NewShippingHandler creates a new ShippingHandler
func NewShippingHandler(shippingUseCase *usecase.ShippingUseCase, logger logger.Logger) *ShippingHandler {
	return &ShippingHandler{
		shippingUseCase: shippingUseCase,
		logger:          logger,
	}
}

// CalculateShippingOptions handles calculating available shipping options for an address and order details
func (h *ShippingHandler) CalculateShippingOptions(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request contracts.CalculateShippingOptionsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	input := request.ToUseCaseInput()

	// Convert to domain address and calculate shipping options
	shippingOptions, err := h.shippingUseCase.CalculateShippingOptions(input)
	if err != nil {
		h.logger.Error("Failed to calculate shipping options: %v", err)
		http.Error(w, "Failed to calculate shipping options", http.StatusInternalServerError)
		return
	}

	// Convert to DTO response
	response := contracts.CreateShippingOptionsListResponse(shippingOptions.Options, len(shippingOptions.Options), 1, len(shippingOptions.Options))

	// Return shipping options
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// CreateShippingMethod handles creating a new shipping method (admin only)
func (h *ShippingHandler) CreateShippingMethod(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request contracts.CreateShippingMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to use case input and create shipping method
	input := request.ToCreateShippingMethodInput()
	method, err := h.shippingUseCase.CreateShippingMethod(input)
	if err != nil {
		h.logger.Error("Failed to create shipping method: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to DTO and return
	methodDTO := method.ToShippingMethodDTO()
	if methodDTO == nil {
		http.Error(w, "Failed to convert shipping method to DTO", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(methodDTO)
}

// CreateShippingZone handles creating a new shipping zone (admin only)
func (h *ShippingHandler) CreateShippingZone(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request contracts.CreateShippingZoneRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to use case input and create shipping zone
	input := request.ToCreateShippingZoneInput()
	zone, err := h.shippingUseCase.CreateShippingZone(input)
	if err != nil {
		h.logger.Error("Failed to create shipping zone: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to DTO and return
	zoneDTO := zone.ToShippingZoneDTO()
	if zoneDTO == nil {
		http.Error(w, "Failed to convert shipping zone to DTO", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(zoneDTO)
}

// CreateShippingRate handles creating a new shipping rate (admin only)
func (h *ShippingHandler) CreateShippingRate(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request contracts.CreateShippingRateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to use case input and create shipping rate
	input := request.ToCreateShippingRateInput()
	rate, err := h.shippingUseCase.CreateShippingRate(input)
	if err != nil {
		h.logger.Error("Failed to create shipping rate: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to DTO and return
	rateDTO := rate.ToShippingRateDTO()
	if rateDTO == nil {
		http.Error(w, "Failed to convert shipping rate to DTO", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rateDTO)
}

// CreateWeightBasedRate handles creating a new weight-based shipping rate (admin only)
func (h *ShippingHandler) CreateWeightBasedRate(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request contracts.CreateWeightBasedRateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to use case input and create weight-based rate
	input := request.ToCreateWeightBasedRateInput()
	rate, err := h.shippingUseCase.CreateWeightBasedRate(input)
	if err != nil {
		h.logger.Error("Failed to create weight-based rate: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to DTO and return
	rateDTO := rate.ToWeightBasedRateDTO()
	if rateDTO == nil {
		http.Error(w, "Failed to convert weight-based rate to DTO", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rateDTO)
}

// CreateValueBasedRate handles creating a new value-based shipping rate (admin only)
func (h *ShippingHandler) CreateValueBasedRate(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request contracts.CreateValueBasedRateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to use case input and create value-based rate
	input := request.ToCreateValueBasedRateInput()
	rate, err := h.shippingUseCase.CreateValueBasedRate(input)
	if err != nil {
		h.logger.Error("Failed to create value-based rate: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to DTO and return
	rateDTO := rate.ToValueBasedRateDTO()
	if rateDTO == nil {
		http.Error(w, "Failed to convert value-based rate to DTO", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rateDTO)
}
