package user

import (
	"net/http"

	userhandlers "github.com/Cryezidl/go-todo-api/internal/handlers/user"
	"github.com/go-chi/chi/v5"
)

func RegisterUserRoutes(r chi.Router, h *userhandlers.UserHandler, authMW func(http.Handler) http.Handler) {
	r.Route("/users", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authMW)

			r.Get("/me", h.GetMe)
			r.Patch("/me", h.Update)
			r.Patch("/me/email", h.UpdateEmail)
			r.Patch("/me/password", h.UpdatePassword)
			r.Delete("/me", h.Delete)
		})
		r.Get("/{username}", h.GetByName)
	})
}
