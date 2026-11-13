package annotation

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	// Root annotations resource
	r.Route("/annotations", func(ann chi.Router) {
		ann.Use(authMiddleware)
		ann.Post("/", h.Create)
		ann.Get("/{id}", h.Get)
		ann.Patch("/{id}", h.Update)
		ann.Delete("/{id}", h.Delete)
	})

	// Subresource: /books/{id}/annotations
	r.With(authMiddleware).Get("/books/{id}/annotations", h.ListByBook)
}
