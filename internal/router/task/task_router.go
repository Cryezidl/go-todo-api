package task

import (
	"net/http"

	taskshandlers "github.com/Cryezidl/go-todo-api/internal/handlers/task"
	"github.com/go-chi/chi/v5"
)

func RegisterTaskRoutes(r chi.Router, h *taskshandlers.TaskHandlers, authMW func(http.Handler) http.Handler) {
	r.Route("/tasks", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authMW)

			r.Get("/", h.GetAll)
			r.Post("/", h.Create)
			r.Get("/{id}", h.GetById)
			r.Patch("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)
			r.Patch("/{id}/status", h.UpdateStatus)
		})
	})
}
