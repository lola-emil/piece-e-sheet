package expense

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router chi.Router, h *handler) {
	router.Route("/expenses", func(r chi.Router) {
		r.Get("/", h.FindAll)
		r.Get("/{id}", h.FindByID)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}
