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

type DeviceRepository struct {
	queries *db.Queries
}

func NewDeviceRepository(queries *db.Queries) *DeviceRepository {
	return &DeviceRepository{queries: queries}
}

func (r *DeviceRepository) Create(ctx context.Context, userID uuid.UUID, name, platform string) (*domain.Device, error) {
	row, err := r.queries.CreateDevice(ctx, db.CreateDeviceParams{
		UserID:   ToPGUUID(userID),
		Name:     name,
		Platform: platform,
	})
	if err != nil {
		return nil, err
	}

	return toDomainDevice(row), nil
}

func (r *DeviceRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Device, error) {
	row, err := r.queries.GetDeviceByID(ctx, db.GetDeviceByIDParams{
		ID:     ToPGUUID(id),
		UserID: ToPGUUID(userID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainDevice(row), nil
}

func (r *DeviceRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Device, error) {
	rows, err := r.queries.ListDevicesByUserID(ctx, ToPGUUID(userID))
	if err != nil {
		return nil, err
	}

	devices := make([]domain.Device, 0, len(rows))
	for _, row := range rows {
		devices = append(devices, *toDomainDevice(row))
	}

	return devices, nil
}

func (r *DeviceRepository) Update(ctx context.Context, id, userID uuid.UUID, name string) (*domain.Device, error) {
	row, err := r.queries.UpdateDevice(ctx, db.UpdateDeviceParams{
		ID:     ToPGUUID(id),
		UserID: ToPGUUID(userID),
		Name:   name,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainDevice(row), nil
}

func (r *DeviceRepository) UpdateSyncTime(ctx context.Context, id, userID uuid.UUID, syncTime time.Time) error {
	return r.queries.UpdateDeviceSyncTime(ctx, db.UpdateDeviceSyncTimeParams{
		ID:           ToPGUUID(id),
		UserID:       ToPGUUID(userID),
		LastSyncedAt: ToPGTimestamptz(&syncTime),
	})
}

func (r *DeviceRepository) Revoke(ctx context.Context, id, userID uuid.UUID) error {
	return r.queries.RevokeDevice(ctx, db.RevokeDeviceParams{
		ID:     ToPGUUID(id),
		UserID: ToPGUUID(userID),
	})
}

func toDomainDevice(d db.Device) *domain.Device {
	return &domain.Device{
		ID:           FromPGUUID(d.ID),
		UserID:       FromPGUUID(d.UserID),
		Name:         d.Name,
		Platform:     d.Platform,
		LastSyncedAt: FromPGTimestamp(d.LastSyncedAt),
		CreatedAt:    d.CreatedAt.Time,
		RevokedAt:    FromPGTimestamp(d.RevokedAt),
	}
}
