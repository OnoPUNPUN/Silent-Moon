package book_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
	bookservice "github.com/silentmoon/api/internal/service/book"
)

type mockBookRepo struct {
	books map[uuid.UUID]*domain.Book
}

func newMockBookRepo() *mockBookRepo {
	return &mockBookRepo{books: make(map[uuid.UUID]*domain.Book)}
}

func (m *mockBookRepo) Create(ctx context.Context, book *domain.Book) (*domain.Book, error) {
	b := *book
	b.CreatedAt = time.Now()
	b.UpdatedAt = time.Now()
	m.books[b.ID] = &b
	return &b, nil
}

func (m *mockBookRepo) GetByID(ctx context.Context, id, userID uuid.UUID) (*domain.Book, error) {
	if b, ok := m.books[id]; ok && b.UserID == userID && b.DeletedAt == nil {
		return b, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockBookRepo) ListByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Book, error) {
	var list []domain.Book
	for _, b := range m.books {
		if b.UserID == userID && b.DeletedAt == nil {
			list = append(list, *b)
		}
	}
	return list, nil
}

func (m *mockBookRepo) Update(ctx context.Context, id, userID uuid.UUID, title string, author, description, coverKey *string) (*domain.Book, error) {
	if b, ok := m.books[id]; ok && b.UserID == userID && b.DeletedAt == nil {
		b.Title = title
		b.Author = author
		b.Description = description
		b.CoverKey = coverKey
		b.UpdatedAt = time.Now()
		return b, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockBookRepo) UpdateStatus(ctx context.Context, id, userID uuid.UUID, status domain.BookStatus, fileSize int64, checksum string) (*domain.Book, error) {
	if b, ok := m.books[id]; ok && b.UserID == userID && b.DeletedAt == nil {
		b.Status = status
		b.FileSize = fileSize
		b.Checksum = checksum
		b.UpdatedAt = time.Now()
		return b, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockBookRepo) SoftDelete(ctx context.Context, id, userID uuid.UUID) error {
	if b, ok := m.books[id]; ok && b.UserID == userID && b.DeletedAt == nil {
		now := time.Now()
		b.DeletedAt = &now
		return nil
	}
	return domain.ErrNotFound
}

func TestBookService_CRUD(t *testing.T) {
	repo := newMockBookRepo()
	svc := bookservice.NewService(repo)
	ctx := context.Background()
	userID := uuid.New()

	author := "Frank Herbert"
	desc := "Classic sci-fi masterpiece"

	// 1. Create Book
	book, err := svc.Create(ctx, userID, bookservice.CreateBookRequest{
		Title:       "Dune",
		Author:      &author,
		Description: &desc,
		Format:      domain.BookFormatEPUB,
	})
	if err != nil {
		t.Fatalf("failed to create book: %v", err)
	}
	if book.Title != "Dune" {
		t.Errorf("expected title 'Dune', got '%s'", book.Title)
	}
	if book.Status != domain.BookStatusPending {
		t.Errorf("expected status 'pending', got '%s'", book.Status)
	}

	// 2. Get Book
	fetched, err := svc.Get(ctx, book.ID, userID)
	if err != nil {
		t.Fatalf("failed to get book: %v", err)
	}
	if fetched.ID != book.ID {
		t.Errorf("book ID mismatch")
	}

	// 3. Update Book
	newTitle := "Dune: Deluxe Edition"
	updated, err := svc.Update(ctx, book.ID, userID, bookservice.UpdateBookRequest{
		Title:  newTitle,
		Author: &author,
	})
	if err != nil {
		t.Fatalf("failed to update book: %v", err)
	}
	if updated.Title != newTitle {
		t.Errorf("expected updated title '%s', got '%s'", newTitle, updated.Title)
	}

	// 4. Soft Delete Book
	if err := svc.Delete(ctx, book.ID, userID); err != nil {
		t.Fatalf("failed to delete book: %v", err)
	}

	// 5. Must not be accessible after soft delete
	_, err = svc.Get(ctx, book.ID, userID)
	if err == nil {
		t.Fatal("expected error getting soft deleted book")
	}
}
