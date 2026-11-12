package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/silentmoon/api/internal/config"
	"github.com/silentmoon/api/internal/handler"
	"github.com/silentmoon/api/internal/platform/database"
	"github.com/silentmoon/api/internal/platform/logger"
	"github.com/silentmoon/api/internal/platform/server"
	"github.com/silentmoon/api/migrations"
)

func main() {
	runMigrate := flag.Bool("migrate", false, "run database migrations and exit")
	flag.Parse()

	cfg := config.Load()
	appLogger := logger.New(cfg.AppEnv)

	appLogger.Info("initializing personal reading library api",
		slog.String("env", cfg.AppEnv),
		slog.String("port", cfg.HTTPPort),
	)

	// Attempt database connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var db *database.DB
	var err error

	db, err = database.Connect(ctx, cfg.DatabaseURL)
	if err != nil {
		appLogger.Warn("database connection could not be established; starting with degraded state",
			slog.String("error", err.Error()),
		)
	} else {
		appLogger.Info("database connection established successfully")
		defer db.Close()

		// Run migrations if connected
		if err := database.RunMigrations(migrations.FS, ".", cfg.DatabaseURL); err != nil {
			appLogger.Error("failed to run database migrations", slog.String("error", err.Error()))
			if *runMigrate {
				os.Exit(1)
			}
		}
	}

	if *runMigrate {
		appLogger.Info("migrations completed, exiting")
		return
	}

	srv := server.New(cfg, appLogger)
	r := srv.Router()

	// Health check handlers
	healthHandler := handler.NewHealthHandler(db)
	r.Get("/health", healthHandler.Check)

	// API v1 group
	r.Route("/api/v1", func(api chi.Router) {
		api.Get("/health", healthHandler.Check)
	})

	if err := srv.Start(); err != nil {
		appLogger.Error("server terminated with error", slog.String("error", err.Error()))
		os.Exit(1)
	}
}
