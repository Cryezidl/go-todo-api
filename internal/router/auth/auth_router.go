package auth

import (
	authhandlers "github.com/Cryezidl/go-todo-api/internal/handlers/auth"
	"github.com/go-chi/chi/v5"
)

func RegisterAuthRoutes(r chi.Router, h *authhandlers.AuthHandler) {
	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/register", h.Register)
	})
}
