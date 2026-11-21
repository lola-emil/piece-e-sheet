package auth

import (
	"github.com/go-chi/chi/v5"
)

func RegisterRoutes(router chi.Router, h *handler) {
	router.Route("/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/refresh", h.Refresh)
		r.Get("/google", h.GoogleLogin)
		r.Get("/callback", h.GoogleCallback)

		r.Post("/logout", h.Logout)

		r.With(AuthMiddleware).Get("/me", h.GetMe)
	})
}
