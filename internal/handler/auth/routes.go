package auth

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	r.Route("/auth", func(auth chi.Router) {
		auth.Post("/register", h.Register)
		auth.Post("/login", h.Login)
		auth.Post("/refresh", h.Refresh)
		auth.Post("/logout", h.Logout)

		// Protected auth routes
		auth.Group(func(protected chi.Router) {
			protected.Use(authMiddleware)
			protected.Get("/me", h.Me)
		})
	})
}
