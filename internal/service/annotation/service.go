package annotation

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
)

type CreateAnnotationRequest struct {
	BookID       uuid.UUID             `json:"bookId"`
	Type         domain.AnnotationType `json:"type"`
	SelectedText *string               `json:"selectedText,omitempty"`
	Color        *string               `json:"color,omitempty"`
	Note         *string               `json:"note,omitempty"`
	Location     json.RawMessage       `json:"location"`
}

type UpdateAnnotationRequest struct {
	Type         *domain.AnnotationType `json:"type,omitempty"`
	SelectedText *string                `json:"selectedText,omitempty"`
	Color        *string                `json:"color,omitempty"`
	Note         *string                `json:"note,omitempty"`
	Location     json.RawMessage        `json:"location,omitempty"`
}

type Service struct {
	annRepo  domain.AnnotationRepository
	bookRepo domain.BookRepository
}

func NewService(annRepo domain.AnnotationRepository, bookRepo domain.BookRepository) *Service {
	return &Service{
		annRepo:  annRepo,
		bookRepo: bookRepo,
	}
}

func (s *Service) Create(ctx context.Context, userID, deviceID uuid.UUID, req CreateAnnotationRequest) (*domain.Annotation, error) {
	// Verify user owns book
	if _, err := s.bookRepo.GetByID(ctx, req.BookID, userID); err != nil {
		return nil, err
	}

	annType := domain.AnnotationType(strings.ToLower(string(req.Type)))
	switch annType {
	case domain.AnnotationTypeHighlight, domain.AnnotationTypeUnderline, domain.AnnotationTypeStrikethrough, domain.AnnotationTypeNote:
	default:
		return nil, fmt.Errorf("%w: invalid annotation type", domain.ErrValidation)
	}

	if err := domain.ValidateLocation(req.Location); err != nil {
		return nil, fmt.Errorf("%w: %s", domain.ErrValidation, err.Error())
	}

	var devIDPtr *uuid.UUID
	if deviceID != uuid.Nil {
		devIDPtr = &deviceID
	}

	ann := &domain.Annotation{
		ID:           uuid.New(),
		UserID:       userID,
		BookID:       req.BookID,
		DeviceID:     devIDPtr,
		Type:         annType,
		SelectedText: req.SelectedText,
		Color:        req.Color,
		Note:         req.Note,
		Location:     req.Location,
	}

	return s.annRepo.Create(ctx, ann)
}

func (s *Service) Get(ctx context.Context, id, userID uuid.UUID) (*domain.Annotation, error) {
	return s.annRepo.GetByID(ctx, id, userID)
}

func (s *Service) ListByBook(ctx context.Context, bookID, userID uuid.UUID) ([]domain.Annotation, error) {
	if _, err := s.bookRepo.GetByID(ctx, bookID, userID); err != nil {
		return nil, err
	}
	return s.annRepo.ListByBookID(ctx, bookID, userID)
}

func (s *Service) Update(ctx context.Context, id, userID uuid.UUID, req UpdateAnnotationRequest) (*domain.Annotation, error) {
	if len(req.Location) > 0 {
		if err := domain.ValidateLocation(req.Location); err != nil {
			return nil, fmt.Errorf("%w: %s", domain.ErrValidation, err.Error())
		}
	}

	return s.annRepo.Update(ctx, id, userID, req.Type, req.SelectedText, req.Color, req.Note, req.Location)
}

func (s *Service) Delete(ctx context.Context, id, userID uuid.UUID) error {
	return s.annRepo.SoftDelete(ctx, id, userID)
}
