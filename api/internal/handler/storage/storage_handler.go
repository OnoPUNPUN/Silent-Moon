package storage

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
	"github.com/silentmoon/api/internal/middleware"
	"github.com/silentmoon/api/internal/platform/response"
	storageservice "github.com/silentmoon/api/internal/service/storage"
)

type Handler struct {
	service *storageservice.Service
}

func NewHandler(service *storageservice.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetUploadURL(w http.ResponseWriter, r *http.Request) {
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

	res, err := h.service.GenerateBookUploadURL(r.Context(), bookID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "book not found")
			return
		}
		response.InternalError(w, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, res)
}

func (h *Handler) ConfirmUpload(w http.ResponseWriter, r *http.Request) {
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

	book, err := h.service.ConfirmBookUpload(r.Context(), bookID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "book not found")
			return
		}
		response.Error(w, http.StatusBadRequest, response.CodeStorageError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, book)
}

func (h *Handler) GetDownloadURL(w http.ResponseWriter, r *http.Request) {
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

	res, err := h.service.GenerateBookDownloadURL(r.Context(), bookID, userID)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			response.NotFound(w, "book not found")
			return
		}
		response.Error(w, http.StatusBadRequest, response.CodeStorageError, err.Error())
		return
	}

	response.JSON(w, http.StatusOK, res)
}
