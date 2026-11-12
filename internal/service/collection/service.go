package collection

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
)

type CreateCollectionRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type UpdateCollectionRequest struct {
	Name        string  `json:"name"`
	Description *string `json:"description,omitempty"`
}

type Service struct {
	collRepo domain.CollectionRepository
	bookRepo domain.BookRepository
}

func NewService(collRepo domain.CollectionRepository, bookRepo domain.BookRepository) *Service {
	return &Service{
		collRepo: collRepo,
		bookRepo: bookRepo,
	}
}

func (s *Service) Create(ctx context.Context, userID uuid.UUID, req CreateCollectionRequest) (*domain.Collection, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: collection name is required", domain.ErrValidation)
	}

	return s.collRepo.Create(ctx, userID, name, req.Description)
}

func (s *Service) Get(ctx context.Context, id, userID uuid.UUID) (*domain.Collection, error) {
	return s.collRepo.GetByID(ctx, id, userID)
}

func (s *Service) List(ctx context.Context, userID uuid.UUID) ([]domain.Collection, error) {
	return s.collRepo.ListByUserID(ctx, userID)
}

func (s *Service) Update(ctx context.Context, id, userID uuid.UUID, req UpdateCollectionRequest) (*domain.Collection, error) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, fmt.Errorf("%w: collection name is required", domain.ErrValidation)
	}

	return s.collRepo.Update(ctx, id, userID, name, req.Description)
}

func (s *Service) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.collRepo.SoftDelete(ctx, id, userID)
}

func (s *Service) AddBook(ctx context.Context, collectionID, bookID, userID uuid.UUID) error {
	// Verify user owns collection
	if _, err := s.collRepo.GetByID(ctx, collectionID, userID); err != nil {
		return err
	}
	// Verify user owns book
	if _, err := s.bookRepo.GetByID(ctx, bookID, userID); err != nil {
		return err
	}

	return s.collRepo.AddBook(ctx, collectionID, bookID)
}

func (s *Service) RemoveBook(ctx context.Context, collectionID, bookID, userID uuid.UUID) error {
	// Verify user owns collection
	if _, err := s.collRepo.GetByID(ctx, collectionID, userID); err != nil {
		return err
	}

	return s.collRepo.RemoveBook(ctx, collectionID, bookID)
}

func (s *Service) ListBooks(ctx context.Context, collectionID, userID uuid.UUID) ([]domain.Book, error) {
	// Verify user owns collection
	if _, err := s.collRepo.GetByID(ctx, collectionID, userID); err != nil {
		return nil, err
	}

	return s.collRepo.ListBooks(ctx, collectionID, userID)
}
