package tag

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
)

type CreateTagRequest struct {
	Name string `json:"name"`
}

type Service struct {
	tagRepo  domain.TagRepository
	bookRepo domain.BookRepository
}

func NewService(tagRepo domain.TagRepository, bookRepo domain.BookRepository) *Service {
	return &Service{
		tagRepo:  tagRepo,
		bookRepo: bookRepo,
	}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, req CreateTagRequest) (*domain.Tag, error) {
	name := strings.TrimSpace(strings.ToLower(req.Name))
	if name == "" {
		return nil, fmt.Errorf("%w: tag name is required", domain.ErrValidation)
	}

	return s.tagRepo.Create(ctx, userID, name)
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]domain.Tag, error) {
	return s.tagRepo.ListByUserID(ctx, userID)
}

func (s *Service) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.tagRepo.Delete(ctx, id, userID)
}

func (s *Service) AddTagToBook(ctx context.Context, bookID, tagID, userID uuid.UUID) error {
	// Verify user owns book
	if _, err := s.bookRepo.GetByID(ctx, bookID, userID); err != nil {
		return err
	}

	return s.tagRepo.AddTagToBook(ctx, bookID, tagID)
}

func (s *Service) RemoveTagFromBook(ctx context.Context, bookID, tagID, userID uuid.UUID) error {
	// Verify user owns book
	if _, err := s.bookRepo.GetByID(ctx, bookID, userID); err != nil {
		return err
	}

	return s.tagRepo.RemoveTagFromBook(ctx, bookID, tagID)
}

func (s *Service) ListTagsForBook(ctx context.Context, bookID, userID uuid.UUID) ([]domain.Tag, error) {
	// Verify user owns book
	if _, err := s.bookRepo.GetByID(ctx, bookID, userID); err != nil {
		return nil, err
	}

	return s.tagRepo.ListTagsForBook(ctx, bookID, userID)
}
