package storage_test

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
	storageservice "github.com/silentmoon/api/internal/service/storage"
	"github.com/silentmoon/api/internal/storage"
)

type mockStorageClient struct {
	objects map[string]*storage.ObjectMeta
}

func newMockStorageClient() *mockStorageClient {
	return &mockStorageClient{objects: make(map[string]*storage.ObjectMeta)}
}

func (m *mockStorageClient) GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, lifetime time.Duration) (string, error) {
	return "https://mock-r2.cloudflarestorage.com/upload?key=" + key, nil
}

func (m *mockStorageClient) GeneratePresignedDownloadURL(ctx context.Context, key string, filename string, lifetime time.Duration) (string, error) {
	return "https://mock-r2.cloudflarestorage.com/download?key=" + key + "&filename=" + filename, nil
}

func (m *mockStorageClient) HeadObject(ctx context.Context, key string) (*storage.ObjectMeta, error) {
	if obj, ok := m.objects[key]; ok {
		return obj, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockStorageClient) DeleteObject(ctx context.Context, key string) error {
	delete(m.objects, key)
	return nil
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
	if b, ok := m.books[id]; ok && b.UserID == userID {
		b.Status = status
		b.FileSize = fileSize
		b.Checksum = checksum
		return b, nil
	}
	return nil, domain.ErrNotFound
}

func (m *mockBookRepo) SoftDelete(ctx context.Context, id, userID uuid.UUID) error {
	return nil
}

func TestStorageService_UploadAndDownloadFlow(t *testing.T) {
	userID := uuid.New()
	bookID := uuid.New()
	storageKey := "users/" + userID.String() + "/books/" + bookID.String() + ".epub"

	book := &domain.Book{
		ID:         bookID,
		UserID:     userID,
		Title:      "Hyperion",
		Format:     domain.BookFormatEPUB,
		StorageKey: storageKey,
		Status:     domain.BookStatusPending,
	}

	bookRepo := &mockBookRepo{books: map[uuid.UUID]*domain.Book{bookID: book}}
	storageClient := newMockStorageClient()
	svc := storageservice.NewService(storageClient, bookRepo)
	ctx := context.Background()

	// 1. Generate Presigned Upload URL
	uploadResp, err := svc.GenerateBookUploadURL(ctx, bookID, userID)
	if err != nil {
		t.Fatalf("failed to generate upload url: %v", err)
	}
	if uploadResp.UploadURL == "" || uploadResp.StorageKey != storageKey {
		t.Errorf("unexpected upload response: %+v", uploadResp)
	}

	// 2. Download should fail before upload is confirmed
	_, err = svc.GenerateBookDownloadURL(ctx, bookID, userID)
	if err == nil {
		t.Fatal("expected download to fail before confirmation")
	}

	// 3. Simulate file uploaded to storage
	storageClient.objects[storageKey] = &storage.ObjectMeta{
		Size: 1048576,
		ETag: "abc123etaghash",
	}

	// 4. Confirm upload
	confirmedBook, err := svc.ConfirmBookUpload(ctx, bookID, userID)
	if err != nil {
		t.Fatalf("failed to confirm upload: %v", err)
	}
	if confirmedBook.Status != domain.BookStatusUploaded {
		t.Errorf("expected status 'uploaded', got '%s'", confirmedBook.Status)
	}
	if confirmedBook.FileSize != 1048576 {
		t.Errorf("expected file size 1048576, got %d", confirmedBook.FileSize)
	}

	// 5. Generate Presigned Download URL
	downloadResp, err := svc.GenerateBookDownloadURL(ctx, bookID, userID)
	if err != nil {
		t.Fatalf("failed to generate download url: %v", err)
	}
	if downloadResp.DownloadURL == "" {
		t.Fatal("download URL must not be empty")
	}
}
