package book

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
)

type CreateBookRequest struct {
	Title       string            `json:"title"`
	Author      *string           `json:"author,omitempty"`
	Description *string           `json:"description,omitempty"`
	Format      domain.BookFormat `json:"format"`
	CoverKey    *string           `json:"coverKey,omitempty"`
}

type UpdateBookRequest struct {
	Title       string  `json:"title"`
	Author      *string `json:"author,omitempty"`
	Description *string `json:"description,omitempty"`
	CoverKey    *string `json:"coverKey,omitempty"`
}

type Service struct {
	bookRepo domain.BookRepository
}

func NewService(bookRepo domain.BookRepository) *Service {
	return &Service{bookRepo: bookRepo}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, req CreateBookRequest) (*domain.Book, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, fmt.Errorf("%w: title is required", domain.ErrValidation)
	}

	format := domain.BookFormat(strings.ToLower(string(req.Format)))
	if format != domain.BookFormatPDF && format != domain.BookFormatEPUB {
		return nil, fmt.Errorf("%w: invalid format; must be 'pdf' or 'epub'", domain.ErrValidation)
	}

	// Generate storage key: users/{user_id}/books/{book_id}.{format}
	bookID := uuid.New()
	storageKey := fmt.Sprintf("users/%s/books/%s.%s", userID.String(), bookID.String(), format)

	book := &domain.Book{
		ID:          bookID,
		UserID:      userID,
		Title:       title,
		Author:      req.Author,
		Description: req.Description,
		Format:      format,
		StorageKey:  storageKey,
		CoverKey:    req.CoverKey,
		FileSize:    0,
		Checksum:    "",
		Status:      domain.BookStatusPending,
	}

	return s.bookRepo.Create(ctx, book)
}

func (s *Service) Get(ctx context.Context, id, userID uuid.UUID) (*domain.Book, error) {
	return s.bookRepo.GetByID(ctx, id, userID)
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]domain.Book, error) {
	return s.bookRepo.ListByUserID(ctx, userID)
}

func (s *Service) Update(ctx context.Context, id, userID uuid.UUID, req UpdateBookRequest) (*domain.Book, error) {
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, fmt.Errorf("%w: title cannot be empty", domain.ErrValidation)
	}

	return s.bookRepo.Update(ctx, id, userID, title, req.Author, req.Description, req.CoverKey)
}

func (s *Service) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.bookRepo.SoftDelete(ctx, id, userID)
}

// GenerateCoverStorageKey generates a predictable storage key for a book's cover image.
func GenerateCoverStorageKey(userID, bookID uuid.UUID, ext string) string {
	cleanedExt := strings.TrimPrefix(filepath.Ext(ext), ".")
	if cleanedExt == "" {
		cleanedExt = "jpg"
	}
	return fmt.Sprintf("users/%s/covers/%s.%s", userID.String(), bookID.String(), cleanedExt)
}
