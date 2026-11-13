package storage

import (
	"context"
	"time"
)

type ObjectMeta struct {
	Size         int64
	ETag         string
	ContentType  string
	LastModified time.Time
}

type Client interface {
	GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, lifetime time.Duration) (string, error)
	GeneratePresignedDownloadURL(ctx context.Context, key string, filename string, lifetime time.Duration) (string, error)
	HeadObject(ctx context.Context, key string) (*ObjectMeta, error)
	DeleteObject(ctx context.Context, key string) error
}
