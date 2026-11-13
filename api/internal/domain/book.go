package domain

import (
	"time"

	"github.com/google/uuid"
)

type BookFormat string

const (
	BookFormatPDF  BookFormat = "pdf"
	BookFormatEPUB BookFormat = "epub"
)

type BookStatus string

const (
	BookStatusPending  BookStatus = "pending"
	BookStatusUploaded BookStatus = "uploaded"
	BookStatusFailed   BookStatus = "failed"
)

type Book struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"userId"`
	Title       string     `json:"title"`
	Author      *string    `json:"author,omitempty"`
	Description *string    `json:"description,omitempty"`
	Format      BookFormat `json:"format"`
	StorageKey  string     `json:"storageKey"`
	CoverKey    *string    `json:"coverKey,omitempty"`
	FileSize    int64      `json:"fileSize"`
	Checksum    string     `json:"checksum"`
	Status      BookStatus `json:"status"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}
