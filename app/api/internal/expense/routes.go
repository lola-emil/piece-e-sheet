package expense

import "github.com/go-chi/chi/v5"

func RegisterRoutes(r chi.Router, h *ExpenseHandler) {
	r.Get("/", h.List)
	r.Get("/{id}", h.Get)
}
