package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type DeviceRepository interface {
	Create(ctx context.Context, userID uuid.UUID, name, platform string) (*Device, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (*Device, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]Device, error)
	Update(ctx context.Context, id, userID uuid.UUID, name string) (*Device, error)
	UpdateSyncTime(ctx context.Context, id, userID uuid.UUID, syncTime time.Time) error
	Revoke(ctx context.Context, id, userID uuid.UUID) error
}
