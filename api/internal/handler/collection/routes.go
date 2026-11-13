package collection

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/collections", func(col chi.Router) {
		col.Use(authMiddleware)
		col.Get("/", h.List)
		col.Post("/", h.Create)
		col.Get("/{id}", h.Get)
		col.Patch("/{id}", h.Update)
		col.Delete("/{id}", h.Delete)

		// Collection Books
		col.Get("/{id}/books", h.ListBooks)
		col.Post("/{id}/books/{bookId}", h.AddBook)
		col.Delete("/{id}/books/{bookId}", h.RemoveBook)
	})
}
