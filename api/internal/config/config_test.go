package config

import (
	"os"
	"testing"
)

func TestLoad_ReadsEnv(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	defer os.Unsetenv("APP_ENV")

	cfg := Load()
	if cfg.AppEnv != "test" {
		t.Fatalf("expected APP_ENV to be 'test', got '%s'", cfg.AppEnv)
	}
}
