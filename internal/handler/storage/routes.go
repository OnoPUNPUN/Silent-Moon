package storage

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/books/{id}", func(bk chi.Router) {
		bk.Use(authMiddleware)
		bk.Post("/upload-url", h.GetUploadURL)
		bk.Post("/upload-complete", h.ConfirmUpload)
		bk.Get("/download-url", h.GetDownloadURL)
	})
}
