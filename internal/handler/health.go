package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/silentmoon/api/internal/platform/database"
	"github.com/silentmoon/api/internal/platform/response"
)

type HealthHandler struct {
	db *database.DB
}

func NewHealthHandler(db *database.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

type HealthResponse struct {
	Status    string            `json:"status"`
	Timestamp time.Time         `json:"timestamp"`
	Services  map[string]string `json:"services"`
}

func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	services := make(map[string]string)
	isHealthy := true

	if h.db != nil {
		if err := h.db.Ping(ctx); err != nil {
			services["database"] = "unhealthy: " + err.Error()
			isHealthy = false
		} else {
			services["database"] = "healthy"
		}
	} else {
		services["database"] = "not_configured"
	}

	statusStr := "healthy"
	httpStatus := http.StatusOK
	if !isHealthy {
		statusStr = "degraded"
		httpStatus = http.StatusServiceUnavailable
	}

	response.JSON(w, httpStatus, HealthResponse{
		Status:    statusStr,
		Timestamp: time.Now().UTC(),
		Services:  services,
	})
}
