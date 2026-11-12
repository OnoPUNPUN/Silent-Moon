package sync

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/silentmoon/api/internal/domain"
	"github.com/silentmoon/api/internal/middleware"
	"github.com/silentmoon/api/internal/platform/response"
	syncservice "github.com/silentmoon/api/internal/service/sync"
)

type Handler struct {
	service *syncservice.Service
}

func NewHandler(service *syncservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Push(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	var req domain.SyncPushRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid request body")
		return
	}

	// Default to token's device ID if omitted in payload
	if req.DeviceID == [16]byte{} {
		if tokenDeviceID, ok := middleware.GetDeviceID(r.Context()); ok {
			req.DeviceID = tokenDeviceID
		}
	}

	res, err := h.service.Push(r.Context(), userID, req)
	if err != nil {
		if errors.Is(err, domain.ErrDeviceRevoked) {
			response.Forbidden(w, "device has been revoked")
			return
		}
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "device not found")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Pull(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	deviceID, ok := middleware.GetDeviceID(r.Context())
	if !ok {
		response.Unauthorized(w, "device context required")
		return
	}

	cursorStr := r.URL.Query().Get("cursor")
	cursor, _ := strconv.ParseInt(cursorStr, 10, 64)

	limitStr := r.URL.Query().Get("limit")
	limitVal, _ := strconv.ParseInt(limitStr, 10, 32)
	limit := int32(limitVal)

	res, err := h.service.Pull(r.Context(), userID, deviceID, cursor, limit)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) Status(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	deviceID, ok := middleware.GetDeviceID(r.Context())
	if !ok {
		response.Unauthorized(w, "device context required")
		return
	}

	res, err := h.service.Status(r.Context(), userID, deviceID)
	if err != nil {
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, res)
}
