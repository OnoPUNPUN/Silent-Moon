package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/platform/response"
	"github.com/silentmoon/api/internal/platform/security"
)

type contextKey string

const (
	userClaimsKey contextKey = "user_claims"
	userIDKey     contextKey = "user_id"
	deviceIDKey   contextKey = "device_id"
)

func Authenticate(tm *security.TokenManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				response.Unauthorized(w, "missing Authorization header")
				return
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				response.Unauthorized(w, "invalid Authorization header format; expected 'Bearer <token>'")
				return
			}

			tokenStr := parts[1]
			claims, err := tm.ValidateAccessToken(tokenStr)
			if err != nil {
				response.Unauthorized(w, "invalid or expired access token")
				return
			}

			ctx := context.WithValue(r.Context(), userClaimsKey, claims)
			ctx = context.WithValue(ctx, userIDKey, claims.UserID)
			ctx = context.WithValue(ctx, deviceIDKey, claims.DeviceID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func GetUserID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(userIDKey).(uuid.UUID)
	return id, ok
}

func GetDeviceID(ctx context.Context) (uuid.UUID, bool) {
	id, ok := ctx.Value(deviceIDKey).(uuid.UUID)
	return id, ok
}

func GetUserClaims(ctx context.Context) (*security.UserClaims, bool) {
	claims, ok := ctx.Value(userClaimsKey).(*security.UserClaims)
	return claims, ok
}
