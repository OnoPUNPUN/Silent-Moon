package domain

import (
	"context"

	"github.com/google/uuid"
)

type CollectionRepository interface {
	Create(ctx context.Context, userID uuid.UUID, name string, description *string) (*Collection, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (*Collection, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]Collection, error)
	Update(ctx context.Context, id, userID uuid.UUID, name string, description *string) (*Collection, error)
	SoftDelete(ctx context.Context, id, userID uuid.UUID) error
	AddBook(ctx context.Context, collectionID, bookID uuid.UUID) error
	RemoveBook(ctx context.Context, collectionID, bookID uuid.UUID) error
	ListBooks(ctx context.Context, collectionID, userID uuid.UUID) ([]Book, error)
}
