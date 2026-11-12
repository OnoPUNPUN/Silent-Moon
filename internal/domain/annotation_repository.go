package domain

import (
	"context"
	"encoding/json"

	"github.com/google/uuid"
)

type AnnotationRepository interface {
	Create(ctx context.Context, ann *Annotation) (*Annotation, error)
	GetByID(ctx context.Context, id, userID uuid.UUID) (*Annotation, error)
	ListByBookID(ctx context.Context, bookID, userID uuid.UUID) ([]Annotation, error)
	Update(ctx context.Context, id, userID uuid.UUID, annType *AnnotationType, selectedText, color, note *string, location json.RawMessage) (*Annotation, error)
	SoftDelete(ctx context.Context, id, userID uuid.UUID) error
}
