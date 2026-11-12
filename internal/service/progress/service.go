package progress

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
)

type UpdateProgressRequest struct {
	Location           json.RawMessage `json:"location"`
	ProgressPercentage float64         `json:"progressPercentage"`
	LastReadAt         *time.Time      `json:"lastReadAt,omitempty"`
}

type Service struct {
	progressRepo domain.ReadingProgressRepository
	bookRepo     domain.BookRepository
}

func NewService(progressRepo domain.ReadingProgressRepository, bookRepo domain.BookRepository) *Service {
	return &Service{
		progressRepo: progressRepo,
		bookRepo:     bookRepo,
	}
}

func (s *Service) Update(ctx context.Context, bookID, userID, deviceID uuid.UUID, req UpdateProgressRequest) (*domain.ReadingProgress, error) {
	// Verify user owns book
	if _, err := s.bookRepo.GetByID(ctx, bookID, userID); err != nil {
		return nil, err
	}

	if err := domain.ValidateLocation(req.Location); err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrValidation, err.Error())
	}

	if req.ProgressPercentage < 0 || req.ProgressPercentage > 100 {
		return nil, fmt.Errorf("%w: progress percentage must be between 0 and 100", domain.ErrValidation)
	}

	lastReadAt := time.Now().UTC()
	if req.LastReadAt != nil && !req.LastReadAt.IsZero() {
		lastReadAt = *req.LastReadAt
	}

	p := &domain.ReadingProgress{
		UserID:             userID,
		BookID:             bookID,
		DeviceID:           deviceID,
		Location:           req.Location,
		ProgressPercentage: req.ProgressPercentage,
		LastReadAt:         lastReadAt,
	}

	return s.progressRepo.Upsert(ctx, p)
}

func (s *Service) GetLatestForBook(ctx context.Context, bookID, userID uuid.UUID) (*domain.ReadingProgress, error) {
	if _, err := s.bookRepo.GetByID(ctx, bookID, userID); err != nil {
		return nil, err
	}

	return s.progressRepo.GetLatestForBook(ctx, bookID, userID)
}

func (s *Service) ListForBook(ctx context.Context, bookID, userID uuid.UUID) ([]domain.ReadingProgress, error) {
	if _, err := s.bookRepo.GetByID(ctx, bookID, userID); err != nil {
		return nil, err
	}

	return s.progressRepo.ListForBook(ctx, bookID, userID)
}
