package device

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/devices", func(dev chi.Router) {
		dev.Use(authMiddleware)
		dev.Get("/", h.List)
		dev.Post("/", h.Register)
		dev.Patch("/{id}", h.Update)
		dev.Delete("/{id}", h.Revoke)
	})
}
