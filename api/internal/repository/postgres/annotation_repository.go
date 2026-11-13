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

type AnnotationRepository struct {
	queries *db.Queries
}

func NewAnnotationRepository(queries *db.Queries) *AnnotationRepository {
	return &AnnotationRepository{queries: queries}
}

func (r *AnnotationRepository) Create(ctx context.Context, ann *domain.Annotation) (*domain.Annotation, error) {
	var idParam any
	if ann.ID != uuid.Nil {
		idParam = ToPGUUID(ann.ID)
	}

	row, err := r.queries.CreateAnnotation(ctx, db.CreateAnnotationParams{
		UserID:       ToPGUUID(ann.UserID),
		BookID:       ToPGUUID(ann.BookID),
		DeviceID:     ToPGNullableUUID(ann.DeviceID),
		Type:         string(ann.Type),
		SelectedText: ann.SelectedText,
		Color:        ann.Color,
		Note:         ann.Note,
		Location:     []byte(ann.Location),
		ID:           idParam,
	})
	if err != nil {
		return nil, err
	}

	return toDomainAnnotation(row), nil
}

func (r *AnnotationRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Annotation, error) {
	row, err := r.queries.GetAnnotationByID(ctx, db.GetAnnotationByIDParams{
		ID:     ToPGUUID(id),
		UserID: ToPGUUID(userID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainAnnotation(row), nil
}

func (r *AnnotationRepository) ListByBookID(ctx context.Context, bookID, userID uuid.UUID) ([]domain.Annotation, error) {
	rows, err := r.queries.ListAnnotationsByBookID(ctx, db.ListAnnotationsByBookIDParams{
		BookID: ToPGUUID(bookID),
		UserID: ToPGUUID(userID),
	})
	if err != nil {
		return nil, err
	}

	annotations := make([]domain.Annotation, 0, len(rows))
	for _, row := range rows {
		annotations = append(annotations, *toDomainAnnotation(row))
	}

	return annotations, nil
}

func (r *AnnotationRepository) Update(ctx context.Context, id, userID uuid.UUID, annType *domain.AnnotationType, selectedText, color, note *string, location json.RawMessage) (*domain.Annotation, error) {
	var typeStr *string
	if annType != nil {
		s := string(*annType)
		typeStr = &s
	}

	var locBytes []byte
	if len(location) > 0 {
		locBytes = []byte(location)
	}

	row, err := r.queries.UpdateAnnotation(ctx, db.UpdateAnnotationParams{
		ID:           ToPGUUID(id),
		UserID:       ToPGUUID(userID),
		Type:         typeStr,
		SelectedText: selectedText,
		Color:        color,
		Note:         note,
		Location:     locBytes,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainAnnotation(row), nil
}

func (r *AnnotationRepository) SoftDelete(ctx context.Context, id, userID uuid.UUID) error {
	return r.queries.SoftDeleteAnnotation(ctx, db.SoftDeleteAnnotationParams{
		ID:     ToPGUUID(id),
		UserID: ToPGUUID(userID),
	})
}

func toDomainAnnotation(a db.Annotation) *domain.Annotation {
	return &domain.Annotation{
		ID:           FromPGUUID(a.ID),
		UserID:       FromPGUUID(a.UserID),
		BookID:       FromPGUUID(a.BookID),
		DeviceID:     FromPGNullableUUID(a.DeviceID),
		Type:         domain.AnnotationType(a.Type),
		SelectedText: a.SelectedText,
		Color:        a.Color,
		Note:         a.Note,
		Location:     json.RawMessage(a.Location),
		Version:      a.Version,
		CreatedAt:    a.CreatedAt.Time,
		UpdatedAt:    a.UpdatedAt.Time,
		DeletedAt:    FromPGTimestamp(a.DeletedAt),
	}
}
