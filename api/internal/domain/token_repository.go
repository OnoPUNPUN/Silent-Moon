package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type RefreshTokenRepository interface {
	Create(ctx context.Context, userID, deviceID uuid.UUID, tokenHash string, expiresAt time.Time) (*RefreshTokenRecord, error)
	GetByHash(ctx context.Context, tokenHash string) (*RefreshTokenRecord, error)
	Revoke(ctx context.Context, tokenHash string) error
	RevokeAllDeviceTokens(ctx context.Context, userID, deviceID uuid.UUID) error
}
