package handler

import (
	"github.com/go-chi/chi/v5"
	"github.com/user/kanban-saas/pkg/auth"
)

func SetupRoutes(r chi.Router, h *AuthHandler, jwtCfg auth.JWTConfig) {
	r.Route("/api/v1/auth", func(r chi.Router) {
		r.Post("/register", h.Register)
		r.Post("/login", h.Login)
		r.Post("/refresh", h.RefreshToken)

		r.Group(func(r chi.Router) {
			r.Use(auth.RequireAuth(jwtCfg))
			r.Post("/logout", h.Logout)
			r.Get("/me", h.GetProfile)
			r.Put("/me", h.UpdateProfile)
		})
	})
}
