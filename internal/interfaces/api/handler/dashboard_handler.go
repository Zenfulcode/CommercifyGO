package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/zenfulcode/commercify/internal/application/usecase"
	"github.com/zenfulcode/commercify/internal/domain/dto"
	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/interfaces/api/contracts"
)

// DashboardHandler handles dashboard-related HTTP requests
type DashboardHandler struct {
	dashboardUseCase *usecase.DashboardUseCase
	logger           logger.Logger
}

// NewDashboardHandler creates a new DashboardHandler
func NewDashboardHandler(dashboardUseCase *usecase.DashboardUseCase, logger logger.Logger) *DashboardHandler {
	return &DashboardHandler{
		dashboardUseCase: dashboardUseCase,
		logger:           logger,
	}
}

// GetStats handles GET /admin/dashboard/stats
func (h *DashboardHandler) GetStats(w http.ResponseWriter, r *http.Request) {
	var request dto.DashboardStatsRequest

	// Parse query parameters
	query := r.URL.Query()

	// Parse start_date
	if startDateStr := query.Get("start_date"); startDateStr != "" {
		if startDate, err := time.Parse("2006-01-02", startDateStr); err == nil {
			request.StartDate = &startDate
		} else {
			h.logger.Error("Invalid start_date format: %s", startDateStr)
			errorResponse := contracts.ErrorResponse("Invalid start_date format. Use YYYY-MM-DD")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	// Parse end_date
	if endDateStr := query.Get("end_date"); endDateStr != "" {
		if endDate, err := time.Parse("2006-01-02", endDateStr); err == nil {
			// Set to end of day
			endOfDay := time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())
			request.EndDate = &endOfDay
		} else {
			h.logger.Error("Invalid end_date format: %s", endDateStr)
			errorResponse := contracts.ErrorResponse("Invalid end_date format. Use YYYY-MM-DD")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	// Parse days
	if daysStr := query.Get("days"); daysStr != "" {
		if days, err := strconv.Atoi(daysStr); err == nil && days > 0 {
			request.Days = days
		} else {
			h.logger.Error("Invalid days parameter: %s", daysStr)
			errorResponse := contracts.ErrorResponse("Invalid days parameter. Must be a positive integer")
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(errorResponse)
			return
		}
	}

	// Get dashboard stats
	stats, err := h.dashboardUseCase.GetDashboardStats(request)
	if err != nil {
		h.logger.Error("Failed to get dashboard stats: %v", err)
		errorResponse := contracts.ErrorResponse("Failed to retrieve dashboard statistics")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	// Return success response
	successResponse := contracts.SuccessResponseWithMessage(stats, "Dashboard statistics retrieved successfully")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(successResponse)
}
