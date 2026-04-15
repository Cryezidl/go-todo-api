package router

import (
	"net/http"

	"time"

	authhandlers "github.com/Cryezidl/go-todo-api/internal/handlers/auth"
	taskhandlers "github.com/Cryezidl/go-todo-api/internal/handlers/task"
	tasklisthandlers "github.com/Cryezidl/go-todo-api/internal/handlers/tasklist"
	userhandlers "github.com/Cryezidl/go-todo-api/internal/handlers/user"
	authrouters "github.com/Cryezidl/go-todo-api/internal/router/auth"
	taskrouters "github.com/Cryezidl/go-todo-api/internal/router/task"
	tasklistrouters "github.com/Cryezidl/go-todo-api/internal/router/tasklist"
	userrouters "github.com/Cryezidl/go-todo-api/internal/router/user"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Handlers struct {
	AuthHandlers     *authhandlers.AuthHandler
	UserHandlers     *userhandlers.UserHandler
	TaskListHandlers *tasklisthandlers.TaskListHandlers
	TaskHandlers     *taskhandlers.TaskHandlers
}

type Middlewares struct {
	Auth func(http.Handler) http.Handler
}

func SetupRouters(h *Handlers, mw Middlewares) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/api/v1", func(r chi.Router) {
		userrouters.RegisterUserRoutes(r, h.UserHandlers, mw.Auth)
		authrouters.RegisterAuthRoutes(r, h.AuthHandlers)
		tasklistrouters.RegisterTaskListRoutes(r, h.TaskListHandlers, mw.Auth)
		taskrouters.RegisterTaskRoutes(r, h.TaskHandlers, mw.Auth)
	})
	return r
}
