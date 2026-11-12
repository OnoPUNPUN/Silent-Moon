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

type BookmarkRepository struct {
	queries *db.Queries
}

func NewBookmarkRepository(queries *db.Queries) *BookmarkRepository {
	return &BookmarkRepository{queries: queries}
}

func (r *BookmarkRepository) Create(ctx context.Context, b *domain.Bookmark) (*domain.Bookmark, error) {
	var idParam any
	if b.ID != uuid.Nil {
		idParam = ToPGUUID(b.ID)
	}

	row, err := r.queries.CreateBookmark(ctx, db.CreateBookmarkParams{
		UserID:   ToPGUUID(b.UserID),
		BookID:   ToPGUUID(b.BookID),
		DeviceID: ToPGNullableUUID(b.DeviceID),
		Title:    b.Title,
		Location: []byte(b.Location),
		ID:       idParam,
	})
	if err != nil {
		return nil, err
	}

	return toDomainBookmark(row), nil
}

func (r *BookmarkRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Bookmark, error) {
	row, err := r.queries.GetBookmarkByID(ctx, db.GetBookmarkByIDParams{
		ID:     ToPGUUID(id),
		UserID: ToPGUUID(userID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainBookmark(row), nil
}

func (r *BookmarkRepository) ListByBookID(ctx context.Context, bookID, userID uuid.UUID) ([]domain.Bookmark, error) {
	rows, err := r.queries.ListBookmarksByBookID(ctx, db.ListBookmarksByBookIDParams{
		BookID: ToPGUUID(bookID),
		UserID: ToPGUUID(userID),
	})
	if err != nil {
		return nil, err
	}

	bookmarks := make([]domain.Bookmark, 0, len(rows))
	for _, row := range rows {
		bookmarks = append(bookmarks, *toDomainBookmark(row))
	}

	return bookmarks, nil
}

func (r *BookmarkRepository) Update(ctx context.Context, id, userID uuid.UUID, title *string, location json.RawMessage) (*domain.Bookmark, error) {
	var locBytes []byte
	if len(location) > 0 {
		locBytes = []byte(location)
	}

	row, err := r.queries.UpdateBookmark(ctx, db.UpdateBookmarkParams{
		ID:       ToPGUUID(id),
		UserID:   ToPGUUID(userID),
		Title:    title,
		Location: locBytes,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainBookmark(row), nil
}

func (r *BookmarkRepository) SoftDelete(ctx context.Context, id, userID uuid.UUID) error {
	return r.queries.SoftDeleteBookmark(ctx, db.SoftDeleteBookmarkParams{
		ID:     ToPGUUID(id),
		UserID: ToPGUUID(userID),
	})
}

func toDomainBookmark(b db.Bookmark) *domain.Bookmark {
	return &domain.Bookmark{
		ID:        FromPGUUID(b.ID),
		UserID:    FromPGUUID(b.UserID),
		BookID:    FromPGUUID(b.BookID),
		DeviceID:  FromPGNullableUUID(b.DeviceID),
		Title:     b.Title,
		Location:  json.RawMessage(b.Location),
		Version:   b.Version,
		CreatedAt: b.CreatedAt.Time,
		UpdatedAt: b.UpdatedAt.Time,
		DeletedAt: FromPGTimestamp(b.DeletedAt),
	}
}
