package progress

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/books/{id}/progress", func(prog chi.Router) {
		prog.Use(authMiddleware)
		prog.Get("/", h.Get)
		prog.Put("/", h.Update)
	})
}
