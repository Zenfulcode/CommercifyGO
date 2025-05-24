package handler

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/money"
	"github.com/zenfulcode/commercify/internal/dto"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
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
	var request dto.CalculateShippingOptionsRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to domain address and calculate shipping options
	address := request.Address.ToDomainAddress()
	shippingOptions, err := h.shippingUseCase.CalculateShippingOptions(
		address,
		money.ToCents(request.OrderValue),
		request.OrderWeight,
	)
	if err != nil {
		h.logger.Error("Failed to calculate shipping options: %v", err)
		http.Error(w, "Failed to calculate shipping options", http.StatusInternalServerError)
		return
	}

	// Convert to DTO response
	response := dto.CalculateShippingOptionsResponse{
		Options: dto.ConvertShippingOptionListToDTO(shippingOptions.Options),
	}

	// Return shipping options
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// GetShippingMethodByID handles retrieving a shipping method by ID
func (h *ShippingHandler) GetShippingMethodByID(w http.ResponseWriter, r *http.Request) {
	// Get method ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["shippingMethodId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid shipping method ID", http.StatusBadRequest)
		return
	}

	// Get shipping method
	method, err := h.shippingUseCase.GetShippingMethodByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get shipping method: %v", err)
		http.Error(w, "Shipping method not found", http.StatusNotFound)
		return
	}

	// Convert to DTO and return
	methodDTO := dto.ConvertToShippingMethodDetailDTO(method)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(methodDTO)
}

// ListShippingMethods handles listing all shipping methods
func (h *ShippingHandler) ListShippingMethods(w http.ResponseWriter, r *http.Request) {
	// Get active parameter from query string
	activeOnly := r.URL.Query().Get("active") == "true"

	// Get shipping methods
	methods, err := h.shippingUseCase.ListShippingMethods(activeOnly)
	if err != nil {
		h.logger.Error("Failed to list shipping methods: %v", err)
		http.Error(w, "Failed to list shipping methods", http.StatusInternalServerError)
		return
	}

	// Convert to DTOs and return
	methodDTOs := dto.ConvertShippingMethodListToDTO(methods)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(methodDTOs)
}

// CreateShippingMethod handles creating a new shipping method (admin only)
func (h *ShippingHandler) CreateShippingMethod(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.CreateShippingMethodRequest
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
	methodDTO := dto.ConvertToShippingMethodDetailDTO(method)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(methodDTO)
}

// UpdateShippingMethod handles updating a shipping method (admin only)
func (h *ShippingHandler) UpdateShippingMethod(w http.ResponseWriter, r *http.Request) {
	// Get method ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["shippingMethodId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid shipping method ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request dto.UpdateShippingMethodRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to use case input and update shipping method
	input := request.ToUpdateShippingMethodInput(uint(id))
	method, err := h.shippingUseCase.UpdateShippingMethod(input)
	if err != nil {
		h.logger.Error("Failed to update shipping method: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to DTO and return
	methodDTO := dto.ConvertToShippingMethodDetailDTO(method)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(methodDTO)
}

// CreateShippingZone handles creating a new shipping zone (admin only)
func (h *ShippingHandler) CreateShippingZone(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.CreateShippingZoneRequest
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
	zoneDTO := dto.ConvertToShippingZoneDTO(zone)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(zoneDTO)
}

// GetShippingZoneByID handles retrieving a shipping zone by ID
func (h *ShippingHandler) GetShippingZoneByID(w http.ResponseWriter, r *http.Request) {
	// Get zone ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["shippingZoneId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid shipping zone ID", http.StatusBadRequest)
		return
	}

	// Get shipping zone
	zone, err := h.shippingUseCase.GetShippingZoneByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get shipping zone: %v", err)
		http.Error(w, "Shipping zone not found", http.StatusNotFound)
		return
	}

	// Convert to DTO and return
	zoneDTO := dto.ConvertToShippingZoneDTO(zone)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zoneDTO)
}

// ListShippingZones handles listing all shipping zones
func (h *ShippingHandler) ListShippingZones(w http.ResponseWriter, r *http.Request) {
	// Get active parameter from query string
	activeOnly := r.URL.Query().Get("active") == "true"

	// Get shipping zones
	zones, err := h.shippingUseCase.ListShippingZones(activeOnly)
	if err != nil {
		h.logger.Error("Failed to list shipping zones: %v", err)
		http.Error(w, "Failed to list shipping zones", http.StatusInternalServerError)
		return
	}

	// Convert to DTOs and return
	zoneDTOs := dto.ConvertShippingZoneListToDTO(zones)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zoneDTOs)
}

// UpdateShippingZone handles updating a shipping zone (admin only)
func (h *ShippingHandler) UpdateShippingZone(w http.ResponseWriter, r *http.Request) {
	// Get zone ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["shippingZoneId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid shipping zone ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request dto.UpdateShippingZoneRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to use case input and update shipping zone
	input := request.ToUpdateShippingZoneInput(uint(id))
	zone, err := h.shippingUseCase.UpdateShippingZone(input)
	if err != nil {
		h.logger.Error("Failed to update shipping zone: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to DTO and return
	zoneDTO := dto.ConvertToShippingZoneDTO(zone)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(zoneDTO)
}

// CreateShippingRate handles creating a new shipping rate (admin only)
func (h *ShippingHandler) CreateShippingRate(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.CreateShippingRateRequest
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
	rateDTO := dto.ConvertToShippingRateDTO(rate)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rateDTO)
}

// GetShippingRateByID handles retrieving a shipping rate by ID
func (h *ShippingHandler) GetShippingRateByID(w http.ResponseWriter, r *http.Request) {
	// Get rate ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["shippingRateId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid shipping rate ID", http.StatusBadRequest)
		return
	}

	// Get shipping rate
	rate, err := h.shippingUseCase.GetShippingRateByID(uint(id))
	if err != nil {
		h.logger.Error("Failed to get shipping rate: %v", err)
		http.Error(w, "Shipping rate not found", http.StatusNotFound)
		return
	}

	// Convert to DTO and return
	rateDTO := dto.ConvertToShippingRateDTO(rate)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rateDTO)
}

// UpdateShippingRate handles updating a shipping rate (admin only)
func (h *ShippingHandler) UpdateShippingRate(w http.ResponseWriter, r *http.Request) {
	// Get rate ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["shippingRateId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid shipping rate ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request dto.UpdateShippingRateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Convert to use case input and update shipping rate
	input := request.ToUpdateShippingRateInput(uint(id))
	rate, err := h.shippingUseCase.UpdateShippingRate(input)
	if err != nil {
		h.logger.Error("Failed to update shipping rate: %v", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Convert to DTO and return
	rateDTO := dto.ConvertToShippingRateDTO(rate)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rateDTO)
}

// CreateWeightBasedRate handles creating a new weight-based shipping rate (admin only)
func (h *ShippingHandler) CreateWeightBasedRate(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.CreateWeightBasedRateRequest
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
	rateDTO := dto.ConvertToWeightBasedRateDTO(rate)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rateDTO)
}

// CreateValueBasedRate handles creating a new value-based shipping rate (admin only)
func (h *ShippingHandler) CreateValueBasedRate(w http.ResponseWriter, r *http.Request) {
	// Parse request body
	var request dto.CreateValueBasedRateRequest
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
	rateDTO := dto.ConvertToValueBasedRateDTO(rate)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(rateDTO)
}

// GetShippingCost handles calculating shipping cost for a specific shipping rate
func (h *ShippingHandler) GetShippingCost(w http.ResponseWriter, r *http.Request) {
	// Get rate ID from URL
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["shippingRateId"], 10, 32)
	if err != nil {
		http.Error(w, "Invalid shipping rate ID", http.StatusBadRequest)
		return
	}

	// Parse request body
	var request dto.CalculateShippingCostRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Calculate shipping cost
	cost, err := h.shippingUseCase.GetShippingCost(
		uint(id),
		money.ToCents(request.OrderValue),
		request.OrderWeight,
	)
	if err != nil {
		h.logger.Error("Failed to calculate shipping cost: %v", err)
		http.Error(w, "Failed to calculate shipping cost", http.StatusInternalServerError)
		return
	}

	// Convert to DTO response and return
	response := dto.CalculateShippingCostResponse{
		Cost: money.FromCents(cost),
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
