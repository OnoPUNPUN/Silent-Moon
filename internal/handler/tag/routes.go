package tag

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) RegisterRoutes(r chi.Router, authMiddleware func(http.Handler) http.Handler) {
	// Root tags resource
	r.Route("/tags", func(tg chi.Router) {
		tg.Use(authMiddleware)
		tg.Get("/", h.List)
		tg.Post("/", h.Create)
		tg.Delete("/{id}", h.Delete)
	})

	// Book tags subresource
	r.Route("/books/{id}/tags", func(bt chi.Router) {
		bt.Use(authMiddleware)
		bt.Get("/", h.ListTagsForBook)
		bt.Post("/{tagId}", h.AddTagToBook)
		bt.Delete("/{tagId}", h.RemoveTagFromBook)
	})
}
