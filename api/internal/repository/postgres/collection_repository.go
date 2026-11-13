package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/silentmoon/api/internal/domain"
	"github.com/silentmoon/api/internal/repository/postgres/db"
)

type CollectionRepository struct {
	queries *db.Queries
}

func NewCollectionRepository(queries *db.Queries) *CollectionRepository {
	return &CollectionRepository{queries: queries}
}

func (r *CollectionRepository) Create(ctx context.Context, userID uuid.UUID, name string, description *string) (*domain.Collection, error) {
	row, err := r.queries.CreateCollection(ctx, db.CreateCollectionParams{
		UserID:      ToPGUUID(userID),
		Name:        name,
		Description: description,
	})
	if err != nil {
		return nil, err
	}

	return toDomainCollection(row), nil
}

func (r *CollectionRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Collection, error) {
	row, err := r.queries.GetCollectionByID(ctx, db.GetCollectionByIDParams{
		ID:     ToPGUUID(id),
		UserID: ToPGUUID(userID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainCollection(row), nil
}

func (r *CollectionRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Collection, error) {
	rows, err := r.queries.ListCollectionsByUserID(ctx, ToPGUUID(userID))
	if err != nil {
		return nil, err
	}

	collections := make([]domain.Collection, 0, len(rows))
	for _, row := range rows {
		collections = append(collections, *toDomainCollection(row))
	}

	return collections, nil
}

func (r *CollectionRepository) Update(ctx context.Context, id, userID uuid.UUID, name string, description *string) (*domain.Collection, error) {
	row, err := r.queries.UpdateCollection(ctx, db.UpdateCollectionParams{
		ID:          ToPGUUID(id),
		UserID:      ToPGUUID(userID),
		Name:        name,
		Description: description,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainCollection(row), nil
}

func (r *CollectionRepository) SoftDelete(ctx context.Context, id, userID uuid.UUID) error {
	return r.queries.SoftDeleteCollection(ctx, db.SoftDeleteCollectionParams{
		ID:     ToPGUUID(id),
		UserID: ToPGUUID(userID),
	})
}

func (r *CollectionRepository) AddBook(ctx context.Context, collectionID, bookID uuid.UUID) error {
	return r.queries.AddBookToCollection(ctx, db.AddBookToCollectionParams{
		CollectionID: ToPGUUID(collectionID),
		BookID:       ToPGUUID(bookID),
	})
}

func (r *CollectionRepository) RemoveBook(ctx context.Context, collectionID, bookID uuid.UUID) error {
	return r.queries.RemoveBookFromCollection(ctx, db.RemoveBookFromCollectionParams{
		CollectionID: ToPGUUID(collectionID),
		BookID:       ToPGUUID(bookID),
	})
}

func (r *CollectionRepository) ListBooks(ctx context.Context, collectionID, userID uuid.UUID) ([]domain.Book, error) {
	rows, err := r.queries.ListBooksInCollection(ctx, db.ListBooksInCollectionParams{
		CollectionID: ToPGUUID(collectionID),
		UserID:       ToPGUUID(userID),
	})
	if err != nil {
		return nil, err
	}

	books := make([]domain.Book, 0, len(rows))
	for _, row := range rows {
		books = append(books, *toDomainBook(row))
	}

	return books, nil
}

func toDomainCollection(c db.Collection) *domain.Collection {
	return &domain.Collection{
		ID:          FromPGUUID(c.ID),
		UserID:      FromPGUUID(c.UserID),
		Name:        c.Name,
		Description: c.Description,
		CreatedAt:   c.CreatedAt.Time,
		UpdatedAt:   c.UpdatedAt.Time,
		DeletedAt:   FromPGTimestamp(c.DeletedAt),
	}
}
