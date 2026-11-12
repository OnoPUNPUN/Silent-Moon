package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/silentmoon/api/internal/domain"
	"github.com/silentmoon/api/internal/repository/postgres/db"
)

type SyncRepository struct {
	queries *db.Queries
}

func NewSyncRepository(queries *db.Queries) *SyncRepository {
	return &SyncRepository{queries: queries}
}

func (r *SyncRepository) CreateChange(ctx context.Context, userID uuid.UUID, deviceID *uuid.UUID, entityType domain.SyncEntityType, entityID uuid.UUID, op domain.SyncOperationType) (*domain.SyncChange, error) {
	row, err := r.queries.CreateSyncChange(ctx, db.CreateSyncChangeParams{
		UserID:     ToPGUUID(userID),
		DeviceID:   ToPGNullableUUID(deviceID),
		EntityType: string(entityType),
		EntityID:   ToPGUUID(entityID),
		Operation:  string(op),
	})
	if err != nil {
		return nil, err
	}

	return toDomainSyncChange(row), nil
}

func (r *SyncRepository) ListChangesAfterCursor(ctx context.Context, userID uuid.UUID, cursor int64, limit int32) ([]domain.SyncChange, error) {
	rows, err := r.queries.ListSyncChangesAfterCursor(ctx, db.ListSyncChangesAfterCursorParams{
		UserID: ToPGUUID(userID),
		ID:     cursor,
		Limit:  limit,
	})
	if err != nil {
		return nil, err
	}

	changes := make([]domain.SyncChange, 0, len(rows))
	for _, row := range rows {
		changes = append(changes, *toDomainSyncChange(row))
	}

	return changes, nil
}

func (r *SyncRepository) GetLatestSequence(ctx context.Context, userID uuid.UUID) (int64, error) {
	return r.queries.GetLatestSyncSequenceForUser(ctx, ToPGUUID(userID))
}

func (r *SyncRepository) RecordOperation(ctx context.Context, opID, userID, deviceID uuid.UUID, entityType domain.SyncEntityType, entityID uuid.UUID, status domain.OperationResultStatus) error {
	_, err := r.queries.RecordSyncOperation(ctx, db.RecordSyncOperationParams{
		OperationID: ToPGUUID(opID),
		UserID:      ToPGUUID(userID),
		DeviceID:    ToPGUUID(deviceID),
		EntityType:  string(entityType),
		EntityID:    ToPGUUID(entityID),
		Status:      string(status),
	})
	return err
}

func (r *SyncRepository) GetOperation(ctx context.Context, opID, userID uuid.UUID) (domain.OperationResultStatus, error) {
	row, err := r.queries.GetSyncOperationByID(ctx, db.GetSyncOperationByIDParams{
		OperationID: ToPGUUID(opID),
		UserID:      ToPGUUID(userID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrNotFound
		}
		return "", err
	}

	return domain.OperationResultStatus(row.Status), nil
}

func toDomainSyncChange(c db.SyncChange) *domain.SyncChange {
	return &domain.SyncChange{
		ID:         c.ID,
		UserID:     FromPGUUID(c.UserID),
		DeviceID:   FromPGNullableUUID(c.DeviceID),
		EntityType: domain.SyncEntityType(c.EntityType),
		EntityID:   FromPGUUID(c.EntityID),
		Operation:  domain.SyncOperationType(c.Operation),
		CreatedAt:  c.CreatedAt.Time,
	}
}
