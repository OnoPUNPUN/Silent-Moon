package main

import (
	"flag"
	"log/slog"
	"os"

	"github.com/silentmoon/api/internal/app"
	"github.com/silentmoon/api/internal/config"
	"github.com/silentmoon/api/internal/platform/logger"
)

func main() {
	runMigrate := flag.Bool("migrate", false, "run database migrations and exit")
	flag.Parse()

	cfg := config.Load()
	appLogger := logger.New(cfg.AppEnv)

	application, err := app.New(cfg, appLogger, *runMigrate)
	if err != nil {
		appLogger.Error("application initialization error", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer application.Close()

	if *runMigrate {
		return
	}

	if err := application.Start(); err != nil {
		appLogger.Error("server terminated with error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
