package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"
	"time"

	"github.com/zenfulcode/commercify/internal/infrastructure/logger"
	"github.com/zenfulcode/commercify/internal/interfaces/api/contracts"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db     *sql.DB
	logger logger.Logger
}

// NewHealthHandler creates a new HealthHandler
func NewHealthHandler(db *sql.DB, logger logger.Logger) *HealthHandler {
	return &HealthHandler{
		db:     db,
		logger: logger,
	}
}

// HealthStatus represents the health status of the service
type HealthStatus struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Version   string            `json:"version,omitempty"`
	Services  map[string]string `json:"services"`
	Uptime    string            `json:"uptime,omitempty"`
}

var startTime = time.Now()

// Health performs a health check and returns the service status
func (h *HealthHandler) Health(w http.ResponseWriter, r *http.Request) {
	h.logger.Info("Health check requested")

	status := "healthy"
	httpStatus := http.StatusOK
	services := make(map[string]string)

	// Check database connectivity
	if err := h.checkDatabase(); err != nil {
		h.logger.Error("Database health check failed: %v", err)
		services["database"] = "unhealthy"
		status = "unhealthy"
		httpStatus = http.StatusServiceUnavailable
	} else {
		services["database"] = "healthy"
	}

	// Calculate uptime
	uptime := time.Since(startTime).String()

	healthStatus := HealthStatus{
		Status:    status,
		Timestamp: time.Now(),
		Version:   "1.2.0", // TODO: Make this configurable
		Services:  services,
		Uptime:    uptime,
	}

	response := contracts.ResponseDTO[HealthStatus]{
		Success: status == "healthy",
		Data:    healthStatus,
	}

	if status != "healthy" {
		response.Error = "One or more services are unhealthy"
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	json.NewEncoder(w).Encode(response)
}

// checkDatabase verifies database connectivity
func (h *HealthHandler) checkDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return h.db.PingContext(ctx)
}
