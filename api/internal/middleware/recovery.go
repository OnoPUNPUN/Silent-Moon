package middleware

import (
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"

	"github.com/silentmoon/api/internal/platform/response"
)

// Recoverer catches panics, logs stack traces, and returns a JSON 500 error.
func Recoverer(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rvr := recover(); rvr != nil {
					if rvr == http.ErrAbortHandler {
						panic(rvr)
					}

					stack := string(debug.Stack())
					logger.Error("panic recovered in HTTP handler",
						slog.Any("error", rvr),
						slog.String("path", r.URL.Path),
						slog.String("method", r.Method),
						slog.String("stack", stack),
					)

					response.InternalError(w, fmt.Sprintf("Internal Server Error"))
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}
