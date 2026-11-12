package storage

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/silentmoon/api/internal/domain"
	"github.com/silentmoon/api/internal/storage"
)

type UploadURLResponse struct {
	UploadURL  string    `json:"uploadUrl"`
	StorageKey string    `json:"storageKey"`
	ExpiresAt  time.Time `json:"expiresAt"`
}

type DownloadURLResponse struct {
	DownloadURL string    `json:"downloadUrl"`
	ExpiresAt   time.Time `json:"expiresAt"`
}

type Service struct {
	storageClient storage.Client
	bookRepo      domain.BookRepository
}

func NewService(storageClient storage.Client, bookRepo domain.BookRepository) *Service {
	return &Service{
		storageClient: storageClient,
		bookRepo:      bookRepo,
	}
}

func (s *Service) GenerateBookUploadURL(ctx context.Context, bookID, userID uuid.UUID) (*UploadURLResponse, error) {
	book, err := s.bookRepo.GetByID(ctx, bookID, userID)
	if err != nil {
		return nil, err
	}

	contentType := "application/octet-stream"
	switch book.Format {
	case domain.BookFormatPDF:
		contentType = "application/pdf"
	case domain.BookFormatEPUB:
		contentType = "application/epub+zip"
	}

	lifetime := 15 * time.Minute
	url, err := s.storageClient.GeneratePresignedUploadURL(ctx, book.StorageKey, contentType, lifetime)
	if err != nil {
		return nil, err
	}

	return &UploadURLResponse{
		UploadURL:  url,
		StorageKey: book.StorageKey,
		ExpiresAt:  time.Now().Add(lifetime),
	}, nil
}

func (s *Service) ConfirmBookUpload(ctx context.Context, bookID, userID uuid.UUID) (*domain.Book, error) {
	book, err := s.bookRepo.GetByID(ctx, bookID, userID)
	if err != nil {
		return nil, err
	}

	meta, err := s.storageClient.HeadObject(ctx, book.StorageKey)
	if err != nil {
		return nil, fmt.Errorf("file verification failed in storage: %w", err)
	}

	if meta.Size <= 0 {
		return nil, fmt.Errorf("uploaded file has invalid size (0 bytes)")
	}

	return s.bookRepo.UpdateStatus(ctx, bookID, userID, domain.BookStatusUploaded, meta.Size, meta.ETag)
}

func (s *Service) GenerateBookDownloadURL(ctx context.Context, bookID, userID uuid.UUID) (*DownloadURLResponse, error) {
	book, err := s.bookRepo.GetByID(ctx, bookID, userID)
	if err != nil {
		return nil, err
	}

	if book.Status != domain.BookStatusUploaded {
		return nil, fmt.Errorf("book is not ready for download (status: %s)", book.Status)
	}

	filename := fmt.Sprintf("%s.%s", sanitizeFilename(book.Title), book.Format)
	lifetime := 15 * time.Minute
	url, err := s.storageClient.GeneratePresignedDownloadURL(ctx, book.StorageKey, filename, lifetime)
	if err != nil {
		return nil, err
	}

	return &DownloadURLResponse{
		DownloadURL: url,
		ExpiresAt:   time.Now().Add(lifetime),
	}, nil
}

func sanitizeFilename(name string) string {
	var b strings.Builder
	for _, r := range name {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '-' || r == '_' || r == ' ' {
			b.WriteRune(r)
		} else {
			b.WriteRune('_')
		}
	}
	res := strings.TrimSpace(b.String())
	if res == "" {
		return "book"
	}
	return res
}
