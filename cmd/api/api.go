package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Cryezidl/go-todo-api/internal/config"
	authhandlers "github.com/Cryezidl/go-todo-api/internal/handlers/auth"
	taskhandlers "github.com/Cryezidl/go-todo-api/internal/handlers/task"
	tasklisthandlers "github.com/Cryezidl/go-todo-api/internal/handlers/tasklist"
	userhandlers "github.com/Cryezidl/go-todo-api/internal/handlers/user"
	authMW "github.com/Cryezidl/go-todo-api/internal/middleware/auth"
	repo "github.com/Cryezidl/go-todo-api/internal/repository/postgres"
	"github.com/Cryezidl/go-todo-api/internal/router"
	authservice "github.com/Cryezidl/go-todo-api/internal/service/auth"
	taskservice "github.com/Cryezidl/go-todo-api/internal/service/task"
	tasklistservice "github.com/Cryezidl/go-todo-api/internal/service/tasklist"
	userservice "github.com/Cryezidl/go-todo-api/internal/service/user"
)

func RunAPI(cfg *config.Config, log *slog.Logger) error {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", cfg.DB.Host, cfg.DB.Port, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.SSLMode)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return err
	}
	//инициализировать репозитории
	userRepository := repo.NewUserRepository(db, log)
	taskRepository := repo.NewTaskRepository(db, log)
	taskListRepository := repo.NewTaskListRepository(db, log)

	//инициализировать service
	userService := userservice.NewUserService(userRepository, log)
	authService := authservice.NewAuthService(userRepository, log, cfg.JWT.Secret, cfg.JWT.Expiration)
	taskService := taskservice.NewTaskService(taskRepository, taskListRepository, log)
	taskListService := tasklistservice.NewTaskListService(taskListRepository, log)

	//иннициализировать handlers
	userHandlers := userhandlers.NewUserHandler(userService, log)
	authHandlers := authhandlers.NewAuthHandler(authService, log, cfg.JWT.Expiration)
	taskHandlers := taskhandlers.NewTaskHandlers(taskService, taskListService, log)
	taskListHandlers := tasklisthandlers.NewTaskListHandlers(taskListService, log)

	//собрать handlers и middlewares для роутера
	h := &router.Handlers{UserHandlers: userHandlers, AuthHandlers: authHandlers,
		TaskListHandlers: taskListHandlers, TaskHandlers: taskHandlers}
	mw := router.Middlewares{Auth: authMW.AuthMiddleware(cfg.JWT.Secret, log)}
	//настроить роутер
	r := router.SetupRouters(h, mw)
	//запустить сервер
	addr := ":" + cfg.App.Port
	srv := http.Server{
		Addr:         addr,
		Handler:      r,
		ReadTimeout:  cfg.HTTP.ReadTimeout,
		WriteTimeout: cfg.HTTP.WriteTimeout,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("server started", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("listen and serve error", slog.Any("err", err))
		}
	}()
	<-stop
	log.Info("stopping server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("server forced to shutdown: %w", err)
	}
	if err := db.Close(); err != nil {
		log.Error("failed to close db", slog.Any("err", err))
	}

	log.Info("server exited properly")
	return nil
}
