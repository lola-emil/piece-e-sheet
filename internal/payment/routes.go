package payment

import "github.com/go-chi/chi/v5"

func RegisterRoutes(router chi.Router, h *handler) {
	router.Route("/payments", func(r chi.Router) {
		r.Get("/debt/{debt_id}", h.GetByDebtID)
		r.Post("/", h.Record)
	})
}
