package server

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

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/silentmoon/api/internal/config"
	custommiddleware "github.com/silentmoon/api/internal/middleware"
)

type Server struct {
	httpServer *http.Server
	router     *chi.Mux
	logger     *slog.Logger
	cfg        *config.Config
}

func New(cfg *config.Config, logger *slog.Logger) *Server {
	r := chi.NewRouter()

	// Base middlewares
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(custommiddleware.RequestLogger(logger))
	r.Use(custommiddleware.Recoverer(logger))
	r.Use(custommiddleware.CORS(cfg.AllowedOrigins))
	r.Use(chimiddleware.Compress(5))

	httpSrv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return &Server{
		httpServer: httpSrv,
		router:     r,
		logger:     logger,
		cfg:        cfg,
	}
}

func (s *Server) Router() *chi.Mux {
	return s.router
}

func (s *Server) Start() error {
	shutdownErr := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		sig := <-quit

		s.logger.Info("shutting down server...", slog.String("signal", sig.String()))

		ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
		defer cancel()

		shutdownErr <- s.httpServer.Shutdown(ctx)
	}()

	s.logger.Info("starting HTTP server", slog.String("port", s.cfg.HTTPPort), slog.String("env", s.cfg.AppEnv))

	err := s.httpServer.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("HTTP server listen failed: %w", err)
	}

	err = <-shutdownErr
	if err != nil {
		return fmt.Errorf("HTTP server graceful shutdown error: %w", err)
	}

	s.logger.Info("server stopped gracefully")
	return nil
}
