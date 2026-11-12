package book

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/books", func(bk chi.Router) {
		bk.Use(authMiddleware)
		bk.Get("/", h.List)
		bk.Post("/", h.Create)
		bk.Get("/{id}", h.Get)
		bk.Patch("/{id}", h.Update)
		bk.Delete("/{id}", h.Delete)
	})
}
