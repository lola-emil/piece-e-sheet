package auth

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router chi.Router, h *handler) {
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
	})
}
