package database

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

// RunMigrations runs up migrations from the provided embed.FS against databaseURL.
func RunMigrations(fs embed.FS, migrationsDir, databaseURL string) error {
	d, err := iofs.New(fs, migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to create migration iofs driver: %w", err)
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, databaseURL)
	if err != nil {
		return fmt.Errorf("failed to initialize migration instance: %w", err)
	}
	defer m.Close()

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			slog.Info("database migrations up to date (no change)")
			return nil
		}
		return fmt.Errorf("failed to apply migrations: %w", err)
	}

	slog.Info("database migrations applied successfully")
	return nil
}
