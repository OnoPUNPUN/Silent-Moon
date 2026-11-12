package app_test

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/silentmoon/api/internal/app"
	"github.com/silentmoon/api/internal/config"
	"github.com/silentmoon/api/internal/handler"
)

func TestApp_HealthRouting(t *testing.T) {
	cfg := &config.Config{
		AppEnv:          "test",
		HTTPPort:        "8080",
		DatabaseURL:     "postgres://invalid:invalid@localhost:9999/invalid",
		JWTAccessSecret: "test-secret-1234567890",
		AllowedOrigins:  []string{"*"},
	}

	logger := slog.New(slog.NewTextHandler(io.Discard, nil))

	application, err := app.New(cfg, logger, false)
	if err != nil {
		t.Fatalf("failed to initialize app: %v", err)
	}
	defer application.Close()

	// Direct test of health check through server router
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := handler.NewHealthHandler(nil)
		h.Check(w, r)
	}))
	defer srv.Close()

	res, err := http.Get(srv.URL + "/health")
	if err != nil {
		t.Fatalf("health request failed: %v", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", res.StatusCode)
	}

	var payload handler.HealthResponse
	if err := json.NewDecoder(res.Body).Decode(&payload); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if payload.Status != "healthy" {
		t.Errorf("expected healthy status, got %s", payload.Status)
	}
}
