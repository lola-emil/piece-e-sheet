package debt

import "github.com/go-chi/chi/v5"

func RegisterRoutes(router chi.Router, h *handler) {
	router.Route("/debts", func(r chi.Router) {
		r.Get("/", h.FindAll)
		r.Post("/", h.Create)
		r.Put("/{id}", h.Update)
		r.Delete("/{id}", h.Delete)
	})
}
