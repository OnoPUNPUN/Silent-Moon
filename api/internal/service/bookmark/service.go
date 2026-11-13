package bookmark

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
)

type CreateBookmarkRequest struct {
	BookID   uuid.UUID       `json:"bookId"`
	Title    *string         `json:"title,omitempty"`
	Location json.RawMessage `json:"location"`
}

type UpdateBookmarkRequest struct {
	Title    *string         `json:"title,omitempty"`
	Location json.RawMessage `json:"location,omitempty"`
}

type Service struct {
	bookmarkRepo domain.BookmarkRepository
	bookRepo     domain.BookRepository
}

func NewService(bookmarkRepo domain.BookmarkRepository, bookRepo domain.BookRepository) *Service {
	return &Service{
		bookmarkRepo: bookmarkRepo,
		bookRepo:     bookRepo,
	}
}

func (s *Service) Create(ctx context.Context, userID, deviceID uuid.UUID, req CreateBookmarkRequest) (*domain.Bookmark, error) {
	// Verify user owns book
	if _, err := s.bookRepo.GetByID(ctx, req.BookID, userID); err != nil {
		return nil, err
	}

	if err := domain.ValidateLocation(req.Location); err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrValidation, err.Error())
	}

	var devIDPtr *uuid.UUID
	if deviceID != uuid.Nil {
		devIDPtr = &deviceID
	}

	bookmark := &domain.Bookmark{
		ID:       uuid.New(),
		UserID:   userID,
		BookID:   req.BookID,
		DeviceID: devIDPtr,
		Title:    req.Title,
		Location: req.Location,
	}

	return s.bookmarkRepo.Create(ctx, bookmark)
}

func (s *Service) Get(ctx context.Context, id, userID uuid.UUID) (*domain.Bookmark, error) {
	return s.bookmarkRepo.GetByID(ctx, id, userID)
}

func (s *Service) ListByBook(ctx context.Context, bookID, userID uuid.UUID) ([]domain.Bookmark, error) {
	if _, err := s.bookRepo.GetByID(ctx, bookID, userID); err != nil {
		return nil, err
	}
	return s.bookmarkRepo.ListByBookID(ctx, bookID, userID)
}

func (s *Service) Update(ctx context.Context, id, userID uuid.UUID, req UpdateBookmarkRequest) (*domain.Bookmark, error) {
	if len(req.Location) > 0 {
		if err := domain.ValidateLocation(req.Location); err != nil {
			return nil, fmt.Errorf("%w: %s", domain.ErrValidation, err.Error())
		}
	}

	return s.bookmarkRepo.Update(ctx, id, userID, req.Title, req.Location)
}

func (s *Service) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.bookmarkRepo.SoftDelete(ctx, id, userID)
}
