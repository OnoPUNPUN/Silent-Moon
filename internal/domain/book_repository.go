package domain

import (
	"context"

	"github.com/google/uuid"
)

type BookRepository interface {
	Create(ctx context.Context, book *Book) (*Book, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (*Book, error)
	ListByUserID(ctx context.Context, userID uuid.UUID) ([]Book, error)
	Update(ctx context.Context, id, userID uuid.UUID, title string, author, description, coverKey *string) (*Book, error)
	UpdateStatus(ctx context.Context, id, userID uuid.UUID, status BookStatus, fileSize int64, checksum string) (*Book, error)
	SoftDelete(ctx context.Context, id, userID uuid.UUID) error
}
