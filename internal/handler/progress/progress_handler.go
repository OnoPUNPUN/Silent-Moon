package progress

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
	"github.com/silentmoon/api/internal/middleware"
	"github.com/silentmoon/api/internal/platform/response"
	progressservice "github.com/silentmoon/api/internal/service/progress"
)

type Handler struct {
	service *progressservice.Service
}

func NewHandler(service *progressservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
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

	bookID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid book id")
		return
	}

	var req progressservice.UpdateProgressRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid request body")
		return
	}

	progress, err := h.service.Update(r.Context(), bookID, userID, deviceID, req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "book not found")
			return
		}
		if errors.Is(err, domain.ErrValidation) {
			response.ValidationError(w, err.Error())
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, progress)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	bookID, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid book id")
		return
	}

	progress, err := h.service.GetLatestForBook(r.Context(), bookID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "reading progress not found")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, progress)
}
