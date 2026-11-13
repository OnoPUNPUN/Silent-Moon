package annotation_test

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
	annservice "github.com/silentmoon/api/internal/service/annotation"
)

type mockAnnRepo struct {
	annotations map[uuid.UUID]*domain.Annotation
}

func newMockAnnRepo() *mockAnnRepo {
	return &mockAnnRepo{annotations: make(map[uuid.UUID]*domain.Annotation)}
}

func (m *mockAnnRepo) Create(ctx context.Context, ann *domain.Annotation) (*domain.Annotation, error) {
	a := *ann
	a.Version = 1
	a.CreatedAt = time.Now()
	a.UpdatedAt = time.Now()
	m.annotations[a.ID] = &a
	return &a, nil
}

func (m *mockAnnRepo) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Annotation, error) {
	if a, ok := m.annotations[id]; ok && a.UserID == userID && a.DeletedAt == nil {
		return a, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockAnnRepo) ListByBookID(ctx context.Context, bookID, userID uuid.UUID) ([]domain.Annotation, error) {
	var list []domain.Annotation
	for _, a := range m.annotations {
		if a.BookID == bookID && a.UserID == userID && a.DeletedAt == nil {
			list = append(list, *a)
		}
	}
	return list, nil
}

func (m *mockAnnRepo) Update(ctx context.Context, id, userID uuid.UUID, annType *domain.AnnotationType, selectedText, color, note *string, location json.RawMessage) (*domain.Annotation, error) {
	if a, ok := m.annotations[id]; ok && a.UserID == userID && a.DeletedAt == nil {
		if annType != nil {
			a.Type = *annType
		}
		if selectedText != nil {
			a.SelectedText = selectedText
		}
		if color != nil {
			a.Color = color
		}
		if note != nil {
			a.Note = note
		}
		if len(location) > 0 {
			a.Location = location
		}
		a.Version++
		a.UpdatedAt = time.Now()
		return a, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockAnnRepo) SoftDelete(ctx context.Context, id, userID uuid.UUID) error {
	if a, ok := m.annotations[id]; ok && a.UserID == userID && a.DeletedAt == nil {
		now := time.Now()
		a.DeletedAt = &now
		a.Version++
		return nil
	}
	return domain.ErrNotFound
}

type mockBookRepo struct {
	books map[uuid.UUID]*domain.Book
}

func (m *mockBookRepo) Create(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	return book, nil
}
func (m *mockBookRepo) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Book, error) {
	if b, ok := m.books[id]; ok && b.UserID == userID {
		return b, nil
	}
	return nil, domain.ErrNotFound
}
func (m *mockBookRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Book, error) {
	return nil, nil
}
func (m *mockBookRepo) Update(ctx context.Context, id, userID uuid.UUID, title string, author, description, coverKey *string) (*domain.Book, error) {
	return nil, nil
}
func (m *mockBookRepo) UpdateStatus(ctx context.Context, id, userID uuid.UUID, status domain.BookStatus, fileSize int64, checksum string) (*domain.Book, error) {
	return nil, nil
}
func (m *mockBookRepo) SoftDelete(ctx context.Context, id, userID uuid.UUID) error {
	return nil
}

func TestAnnotationService_LocationValidation_And_CRUD(t *testing.T) {
	userID := uuid.New()
	deviceID := uuid.New()
	bookID := uuid.New()

	bookRepo := &mockBookRepo{books: map[uuid.UUID]*domain.Book{
		bookID: {ID: bookID, UserID: userID, Title: "Test Book"},
	}}
	annRepo := newMockAnnRepo()
	svc := annservice.NewService(annRepo, bookRepo)
	ctx := context.Background()

	// 1. Invalid location format
	_, err := svc.Create(ctx, userID, deviceID, annservice.CreateAnnotationRequest{
		BookID:   bookID,
		Type:     domain.AnnotationTypeHighlight,
		Location: json.RawMessage(`{"format":"unknown"}`),
	})
	if err == nil {
		t.Fatal("expected error for unsupported location format")
	}

	// 2. Valid PDF location
	pdfLoc := `{"format":"pdf","page":42,"rectangles":[{"x":120.5,"y":340.2,"width":240.0,"height":18.5}],"pageWidth":595.0,"pageHeight":842.0}`
	text := "An important passage"
	color := "#FFE082"

	ann, err := svc.Create(ctx, userID, deviceID, annservice.CreateAnnotationRequest{
		BookID:       bookID,
		Type:         domain.AnnotationTypeHighlight,
		SelectedText: &text,
		Color:        &color,
		Location:     json.RawMessage(pdfLoc),
	})
	if err != nil {
		t.Fatalf("failed to create pdf annotation: %v", err)
	}
	if ann.Version != 1 {
		t.Errorf("expected version 1, got %d", ann.Version)
	}

	// 3. Update note
	newNote := "Review this in chapter 3"
	updated, err := svc.Update(ctx, ann.ID, userID, annservice.UpdateAnnotationRequest{
		Note: &newNote,
	})
	if err != nil {
		t.Fatalf("failed to update annotation: %v", err)
	}
	if updated.Version != 2 {
		t.Errorf("expected version 2 after update, got %d", updated.Version)
	}

	// 4. Soft Delete
	if err := svc.Delete(ctx, ann.ID, userID); err != nil {
		t.Fatalf("failed to delete annotation: %v", err)
	}

	_, err = svc.Get(ctx, ann.ID, userID)
	if err == nil {
		t.Fatal("expected error getting soft deleted annotation")
	}
}
