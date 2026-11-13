package sync

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/sync", func(sync chi.Router) {
		sync.Use(authMiddleware)
		sync.Post("/push", h.Push)
		sync.Get("/changes", h.Pull)
		sync.Get("/status", h.Status)
	})
}
