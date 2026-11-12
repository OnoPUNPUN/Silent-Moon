package bookmark

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	// Root bookmarks resource
	r.Route("/bookmarks", func(bm chi.Router) {
		bm.Use(authMiddleware)
		bm.Post("/", h.Create)
		bm.Get("/{id}", h.Get)
		bm.Patch("/{id}", h.Update)
		bm.Delete("/{id}", h.Delete)
	})

	// Subresource: /books/{id}/bookmarks
	r.With(authMiddleware).Get("/books/{id}/bookmarks", h.ListByBook)
}
