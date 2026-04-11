package tasklist

import (
	"net/http"

	tasklisthandlers "github.com/Cryezidl/go-todo-api/internal/handlers/tasklist"
	"github.com/go-chi/chi/v5"
)

func RegisterTaskListRoutes(r chi.Router, h *tasklisthandlers.TaskListHandlers, authMW func(http.Handler) http.Handler) {
	r.Route("/task-lists", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authMW)

			r.Post("/", h.Create)
			r.Get("/", h.GetAll)
			r.Get("/{id}", h.GetById)
			r.Patch("/{id}", h.Update)
			r.Delete("/{id}", h.Delete)

		})
	})
}
