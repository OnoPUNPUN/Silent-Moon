package device

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
	"github.com/silentmoon/api/internal/middleware"
	"github.com/silentmoon/api/internal/platform/response"
	deviceservice "github.com/silentmoon/api/internal/service/device"
)

type Handler struct {
	service *deviceservice.Service
}

func NewHandler(service *deviceservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	devices, err := h.service.List(r.Context(), userID)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, devices)
}

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	var req deviceservice.RegisterDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid request body")
		return
	}

	dev, err := h.service.Register(r.Context(), userID, req)
	if err != nil {
		if errors.Is(err, domain.ErrValidation) {
			response.ValidationError(w, err.Error())
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusCreated, dev)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	idStr := chi.URLParam(r, "id")
	deviceID, err := uuid.Parse(idStr)
	if err != nil {
		response.ValidationError(w, "invalid device id")
		return
	}

	var req deviceservice.UpdateDeviceRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid request body")
		return
	}

	dev, err := h.service.Update(r.Context(), userID, deviceID, req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "device not found")
			return
		}
		if errors.Is(err, domain.ErrValidation) {
			response.ValidationError(w, err.Error())
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, dev)
}

func (h *Handler) Revoke(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	idStr := chi.URLParam(r, "id")
	deviceID, err := uuid.Parse(idStr)
	if err != nil {
		response.ValidationError(w, "invalid device id")
		return
	}

	if err := h.service.Revoke(r.Context(), userID, deviceID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "device not found")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "device revoked successfully"})
}
