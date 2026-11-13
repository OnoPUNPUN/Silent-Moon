package domain

import (
	"context"

	"github.com/google/uuid"
)

type ReadingProgressRepository interface {
	Upsert(ctx context.Context, progress *ReadingProgress) (*ReadingProgress, error)
	GetLatestForBook(ctx context.Context, bookID, userID uuid.UUID) (*ReadingProgress, error)
	GetByDevice(ctx context.Context, bookID, deviceID, userID uuid.UUID) (*ReadingProgress, error)
	ListForBook(ctx context.Context, bookID, userID uuid.UUID) ([]ReadingProgress, error)
}
