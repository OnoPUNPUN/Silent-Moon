package postgres

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/silentmoon/api/internal/domain"
	"github.com/silentmoon/api/internal/repository/postgres/db"
)

type BookRepository struct {
	queries *db.Queries
}

func NewBookRepository(queries *db.Queries) *BookRepository {
	return &BookRepository{queries: queries}
}

func (r *BookRepository) Create(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	row, err := r.queries.CreateBook(ctx, db.CreateBookParams{
		UserID:      ToPGUUID(book.UserID),
		Title:       book.Title,
		Author:      book.Author,
		Description: book.Description,
		Format:      string(book.Format),
		StorageKey:  book.StorageKey,
		CoverKey:    book.CoverKey,
		FileSize:    book.FileSize,
		Checksum:    book.Checksum,
		Status:      string(book.Status),
	})
	if err != nil {
		return nil, err
	}

	return toDomainBook(row), nil
}

func (r *BookRepository) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Book, error) {
	row, err := r.queries.GetBookByID(ctx, db.GetBookByIDParams{
		ID:     ToPGUUID(id),
		UserID: ToPGUUID(userID),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainBook(row), nil
}

func (r *BookRepository) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Book, error) {
	rows, err := r.queries.ListBooksByUserID(ctx, ToPGUUID(userID))
	if err != nil {
		return nil, err
	}

	books := make([]domain.Book, 0, len(rows))
	for _, row := range rows {
		books = append(books, *toDomainBook(row))
	}

	return books, nil
}

func (r *BookRepository) Update(ctx context.Context, id, userID uuid.UUID, title string, author, description, coverKey *string) (*domain.Book, error) {
	row, err := r.queries.UpdateBook(ctx, db.UpdateBookParams{
		ID:          ToPGUUID(id),
		UserID:      ToPGUUID(userID),
		Title:       title,
		Author:      author,
		Description: description,
		CoverKey:    coverKey,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainBook(row), nil
}

func (r *BookRepository) UpdateStatus(ctx context.Context, id, userID uuid.UUID, status domain.BookStatus, fileSize int64, checksum string) (*domain.Book, error) {
	row, err := r.queries.UpdateBookStatus(ctx, db.UpdateBookStatusParams{
		ID:       ToPGUUID(id),
		UserID:   ToPGUUID(userID),
		Status:   string(status),
		FileSize: fileSize,
		Checksum: checksum,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrNotFound
		}
		return nil, err
	}

	return toDomainBook(row), nil
}

func (r *BookRepository) SoftDelete(ctx context.Context, id, userID uuid.UUID) error {
	return r.queries.SoftDeleteBook(ctx, db.SoftDeleteBookParams{
		ID:     ToPGUUID(id),
		UserID: ToPGUUID(userID),
	})
}

func toDomainBook(b db.Book) *domain.Book {
	return &domain.Book{
		ID:          FromPGUUID(b.ID),
		UserID:      FromPGUUID(b.UserID),
		Title:       b.Title,
		Author:      b.Author,
		Description: b.Description,
		Format:      domain.BookFormat(b.Format),
		StorageKey:  b.StorageKey,
		CoverKey:    b.CoverKey,
		FileSize:    b.FileSize,
		Checksum:    b.Checksum,
		Status:      domain.BookStatus(b.Status),
		CreatedAt:   b.CreatedAt.Time,
		UpdatedAt:   b.UpdatedAt.Time,
		DeletedAt:   FromPGTimestamp(b.DeletedAt),
	}
}
