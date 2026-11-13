package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/silentmoon/api/internal/domain"
	"github.com/silentmoon/api/internal/repository/postgres/db"
)

type RefreshTokenRepository struct {
	queries *db.Queries
}

func NewRefreshTokenRepository(queries *db.Queries) *RefreshTokenRepository {
	return &RefreshTokenRepository{queries: queries}
}

func (r *RefreshTokenRepository) Create(ctx context.Context, userID, deviceID uuid.UUID, tokenHash string, expiresAt time.Time) (*domain.RefreshTokenRecord, error) {
	row, err := r.queries.CreateRefreshToken(ctx, db.CreateRefreshTokenParams{
		UserID:    ToPGUUID(userID),
		DeviceID:  ToPGUUID(deviceID),
		TokenHash: tokenHash,
		ExpiresAt: ToPGTimestamptz(&expiresAt),
	})
	if err != nil {
		return nil, err
	}

	return toDomainRefreshToken(row), nil
}

func (r *RefreshTokenRepository) GetByHash(ctx context.Context, tokenHash string) (*domain.RefreshTokenRecord, error) {
	row, err := r.queries.GetRefreshTokenByHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainRefreshToken(row), nil
}

func (r *RefreshTokenRepository) Revoke(ctx context.Context, tokenHash string) error {
	return r.queries.RevokeRefreshToken(ctx, tokenHash)
}

func (r *RefreshTokenRepository) RevokeAllDeviceTokens(ctx context.Context, userID, deviceID uuid.UUID) error {
	return r.queries.RevokeAllDeviceTokens(ctx, db.RevokeAllDeviceTokensParams{
		UserID:   ToPGUUID(userID),
		DeviceID: ToPGUUID(deviceID),
	})
}

func toDomainRefreshToken(rt db.RefreshToken) *domain.RefreshTokenRecord {
	devID := FromPGUUID(rt.DeviceID)
	var devIDPtr *uuid.UUID
	if devID != uuid.Nil {
		devIDPtr = &devID
	}

	return &domain.RefreshTokenRecord{
		ID:        FromPGUUID(rt.ID),
		UserID:    FromPGUUID(rt.UserID),
		DeviceID:  devIDPtr,
		TokenHash: rt.TokenHash,
		ExpiresAt: rt.ExpiresAt.Time,
		RevokedAt: FromPGTimestamp(rt.RevokedAt),
		CreatedAt: rt.CreatedAt.Time,
	}
}
