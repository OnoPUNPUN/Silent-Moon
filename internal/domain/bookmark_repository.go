package domain

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type BookmarkRepository interface {
	Create(ctx context.Context, b *Bookmark) (*Bookmark, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (*Bookmark, error)
	ListByBookID(ctx context.Context, bookID, userID uuid.UUID) ([]Bookmark, error)
	Update(ctx context.Context, id, userID uuid.UUID, title *string, location json.RawMessage) (*Bookmark, error)
	SoftDelete(ctx context.Context, id, userID uuid.UUID) error
}
