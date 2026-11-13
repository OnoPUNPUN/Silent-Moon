package domain

import (
	"context"

	"github.com/google/uuid"
)

type TagRepository interface {
	Create(ctx context.Context, userID uuid.UUID, name string) (*Tag, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]Tag, error)
	Delete(ctx context.Context, id, userID uuid.UUID) error
	AddTagToBook(ctx context.Context, bookID, tagID uuid.UUID) error
	RemoveTagFromBook(ctx context.Context, bookID, tagID uuid.UUID) error
	ListTagsForBook(ctx context.Context, bookID, userID uuid.UUID) ([]Tag, error)
}
