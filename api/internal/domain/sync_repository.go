package domain

import (
	"context"

	"github.com/google/uuid"
)

type SyncRepository interface {
	CreateChange(ctx context.Context, userID uuid.UUID, deviceID *uuid.UUID, entityType SyncEntityType, entityID uuid.UUID, op SyncOperationType) (*SyncChange, error)
	ListChangesAfterCursor(ctx context.Context, userID uuid.UUID, cursor int64, limit int32) ([]SyncChange, error)
	GetLatestSequence(ctx context.Context, userID uuid.UUID) (int64, error)
	RecordOperation(ctx context.Context, opID, userID, deviceID uuid.UUID, entityType SyncEntityType, entityID uuid.UUID, status OperationResultStatus) error
	GetOperation(ctx context.Context, opID, userID uuid.UUID) (OperationResultStatus, error)
}
