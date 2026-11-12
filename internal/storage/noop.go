package storage

import (
	"context"
	"errors"
	"time"
)

type NoopClient struct{}

func NewNoopClient() *NoopClient {
	return &NoopClient{}
}

func (n *NoopClient) GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, lifetime time.Duration) (string, error) {
	return "", errors.New("cloud storage is not configured; please set R2 environment variables")
}

func (n *NoopClient) GeneratePresignedDownloadURL(ctx context.Context, key string, filename string, lifetime time.Duration) (string, error) {
	return "", errors.New("cloud storage is not configured; please set R2 environment variables")
}

func (n *NoopClient) HeadObject(ctx context.Context, key string) (*ObjectMeta, error) {
	return nil, errors.New("cloud storage is not configured; please set R2 environment variables")
}

func (n *NoopClient) DeleteObject(ctx context.Context, key string) error {
	return errors.New("cloud storage is not configured; please set R2 environment variables")
}
