package bookmark

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
	"github.com/silentmoon/api/internal/middleware"
	"github.com/silentmoon/api/internal/platform/response"
	bookmarkservice "github.com/silentmoon/api/internal/service/bookmark"
)

type Handler struct {
	service *bookmarkservice.Service
}

func NewHandler(service *bookmarkservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	deviceID, _ := middleware.GetDeviceID(r.Context())

	var req bookmarkservice.CreateBookmarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid request body")
		return
	}

	b, err := h.service.Create(r.Context(), userID, deviceID, req)
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

	response.JSON(w, http.StatusCreated, b)
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid bookmark id")
		return
	}

	b, err := h.service.Get(r.Context(), id, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "bookmark not found")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, b)
}

func (h *Handler) ListByBook(w http.ResponseWriter, r *http.Request) {
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

	bookmarks, err := h.service.ListByBook(r.Context(), bookID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "book not found")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, bookmarks)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid bookmark id")
		return
	}

	var req bookmarkservice.UpdateBookmarkRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.ValidationError(w, "invalid request body")
		return
	}

	b, err := h.service.Update(r.Context(), id, userID, req)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "bookmark not found")
			return
		}
		if errors.Is(err, domain.ErrValidation) {
			response.ValidationError(w, err.Error())
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, b)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		response.Unauthorized(w, "unauthenticated")
		return
	}

	id, err := uuid.Parse(chi.URLParam(r, "id"))
	if err != nil {
		response.ValidationError(w, "invalid bookmark id")
		return
	}

	if err := h.service.Delete(r.Context(), id, userID); err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "bookmark not found")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"message": "bookmark deleted successfully"})
}
