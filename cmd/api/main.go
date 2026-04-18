package main

import (
	"log/slog"
	"os"

	"github.com/Cryezidl/go-todo-api/internal/config"
	"github.com/Cryezidl/go-todo-api/pkg/logger"
)

func main() {
	cfg := config.MustLoad()
	log := logger.SetupLogger("local")
	if err := RunAPI(cfg, log); err != nil {
		log.Error("failed to run API", slog.Any("error", err))
		os.Exit(1)
	}

	log.Info("application stopped successfully")
}
