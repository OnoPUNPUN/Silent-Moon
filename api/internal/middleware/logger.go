package middleware

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// RequestLogger returns a middleware that logs incoming HTTP requests using structured slog.
func RequestLogger(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			start := time.Now()

			reqID := middleware.GetReqID(r.Context())

			defer func() {
				latency := time.Since(start)
				status := ww.Status()
				bytes := ww.BytesWritten()

				logLevel := slog.LevelInfo
				if status >= 500 {
					logLevel = slog.LevelError
				} else if status >= 400 {
					logLevel = slog.LevelWarn
				}

				logger.Log(r.Context(), logLevel, "http request",
					slog.String("req_id", reqID),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", status),
					slog.Int("bytes", bytes),
					slog.Duration("latency", latency),
					slog.String("remote_addr", r.RemoteAddr),
					slog.String("user_agent", r.UserAgent()),
				)
			}()

			next.ServeHTTP(ww, r)
		})
	}
}
