package postgres

import (
	"context"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
	"github.com/silentmoon/api/internal/repository/postgres/db"
)

type TagRepository struct {
	queries *db.Queries
}

func NewTagRepository(queries *db.Queries) *TagRepository {
	return &TagRepository{queries: queries}
}

func (r *TagRepository) Create(ctx context.Context, userID uuid.UUID, name string) (*domain.Tag, error) {
	row, err := r.queries.CreateTag(ctx, db.CreateTagParams{
		UserID: ToPGUUID(userID),
		Name:   name,
	})
	if err != nil {
		return nil, err
	}

	return toDomainTag(row), nil
}

func (r *TagRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Tag, error) {
	rows, err := r.queries.ListTagsByUserID(ctx, ToPGUUID(userID))
	if err != nil {
		return nil, err
	}

	tags := make([]domain.Tag, 0, len(rows))
	for _, row := range rows {
		tags = append(tags, *toDomainTag(row))
	}

	return tags, nil
}

func (r *TagRepository) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return r.queries.DeleteTag(ctx, db.DeleteTagParams{
		ID:     ToPGUUID(id),
		UserID: ToPGUUID(userID),
	})
}

func (r *TagRepository) AddTagToBook(ctx context.Context, bookID, tagID uuid.UUID) error {
	return r.queries.AddTagToBook(ctx, db.AddTagToBookParams{
		BookID: ToPGUUID(bookID),
		TagID:  ToPGUUID(tagID),
	})
}

func (r *TagRepository) RemoveTagFromBook(ctx context.Context, bookID, tagID uuid.UUID) error {
	return r.queries.RemoveTagFromBook(ctx, db.RemoveTagFromBookParams{
		BookID: ToPGUUID(bookID),
		TagID:  ToPGUUID(tagID),
	})
}

func (r *TagRepository) ListTagsForBook(ctx context.Context, bookID, userID uuid.UUID) ([]domain.Tag, error) {
	rows, err := r.queries.ListTagsForBook(ctx, db.ListTagsForBookParams{
		BookID: ToPGUUID(bookID),
		UserID: ToPGUUID(userID),
	})
	if err != nil {
		return nil, err
	}

	tags := make([]domain.Tag, 0, len(rows))
	for _, row := range rows {
		tags = append(tags, *toDomainTag(row))
	}

	return tags, nil
}

func toDomainTag(t db.Tag) *domain.Tag {
	return &domain.Tag{
		ID:        FromPGUUID(t.ID),
		UserID:    FromPGUUID(t.UserID),
		Name:      t.Name,
		CreatedAt: t.CreatedAt.Time,
	}
}
