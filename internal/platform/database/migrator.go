package database

import (
	"embed"
	"errors"
	"fmt"
	"log/slog"
	"strings"

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

	migrateURL := normalizeMigrationURL(databaseURL)

	m, err := migrate.NewWithSourceInstance("iofs", d, migrateURL)
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

// normalizeMigrationURL ensures URLs with postgres:// or postgresql:// schemes
// use pgx5:// so that golang-migrate dispatches to the registered pgx/v5 driver.
func normalizeMigrationURL(databaseURL string) string {
	if strings.HasPrefix(databaseURL, "postgresql://") {
		return "pgx5://" + strings.TrimPrefix(databaseURL, "postgresql://")
	}
	if strings.HasPrefix(databaseURL, "postgres://") {
		return "pgx5://" + strings.TrimPrefix(databaseURL, "postgres://")
	}
	return databaseURL
}
