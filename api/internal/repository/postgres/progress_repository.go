package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/silentmoon/api/internal/domain"
	"github.com/silentmoon/api/internal/repository/postgres/db"
)

type ReadingProgressRepository struct {
	queries *db.Queries
}

func NewReadingProgressRepository(queries *db.Queries) *ReadingProgressRepository {
	return &ReadingProgressRepository{queries: queries}
}

func (r *ReadingProgressRepository) Upsert(ctx context.Context, progress *domain.ReadingProgress) (*domain.ReadingProgress, error) {
	row, err := r.queries.UpsertReadingProgress(ctx, db.UpsertReadingProgressParams{
		UserID:             ToPGUUID(progress.UserID),
		BookID:             ToPGUUID(progress.BookID),
		DeviceID:           ToPGUUID(progress.DeviceID),
		Location:           []byte(progress.Location),
		ProgressPercentage: ToPGNumeric(progress.ProgressPercentage),
		LastReadAt:         ToPGTimestamptz(&progress.LastReadAt),
	})
	if err != nil {
		return nil, err
	}

	return toDomainProgress(row), nil
}

func (r *ReadingProgressRepository) GetLatestForBook(ctx context.Context, bookID, userID uuid.UUID) (*domain.ReadingProgress, error) {
	row, err := r.queries.GetLatestReadingProgressForBook(ctx, db.GetLatestReadingProgressForBookParams{
		BookID: ToPGUUID(bookID),
		UserID: ToPGUUID(userID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainProgress(row), nil
}

func (r *ReadingProgressRepository) GetByDevice(ctx context.Context, bookID, deviceID, userID uuid.UUID) (*domain.ReadingProgress, error) {
	row, err := r.queries.GetReadingProgressByDevice(ctx, db.GetReadingProgressByDeviceParams{
		BookID:   ToPGUUID(bookID),
		DeviceID: ToPGUUID(deviceID),
		UserID:   ToPGUUID(userID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainProgress(row), nil
}

func (r *ReadingProgressRepository) ListForBook(ctx context.Context, bookID, userID uuid.UUID) ([]domain.ReadingProgress, error) {
	rows, err := r.queries.ListReadingProgressForBook(ctx, db.ListReadingProgressForBookParams{
		BookID: ToPGUUID(bookID),
		UserID: ToPGUUID(userID),
	})
	if err != nil {
		return nil, err
	}

	progressList := make([]domain.ReadingProgress, 0, len(rows))
	for _, row := range rows {
		progressList = append(progressList, *toDomainProgress(row))
	}

	return progressList, nil
}

func toDomainProgress(rp db.ReadingProgress) *domain.ReadingProgress {
	return &domain.ReadingProgress{
		ID:                 FromPGUUID(rp.ID),
		UserID:             FromPGUUID(rp.UserID),
		BookID:             FromPGUUID(rp.BookID),
		DeviceID:           FromPGUUID(rp.DeviceID),
		Location:           json.RawMessage(rp.Location),
		ProgressPercentage: FromPGNumeric(rp.ProgressPercentage),
		LastReadAt:         rp.LastReadAt.Time,
		UpdatedAt:          rp.UpdatedAt.Time,
	}
}
