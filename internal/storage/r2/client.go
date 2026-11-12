package r2

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/silentmoon/api/internal/storage"
)

type Config struct {
	AccountID string
	AccessKey string
	SecretKey string
	Bucket    string
	Endpoint  string
}

type Client struct {
	s3Client      *s3.Client
	presignClient *s3.PresignClient
	bucket        string
}

func New(cfg Config) (*Client, error) {
	if cfg.Bucket == "" {
		return nil, fmt.Errorf("R2 bucket name is required")
	}

	endpoint := cfg.Endpoint
	if endpoint == "" {
		if cfg.AccountID == "" {
			return nil, fmt.Errorf("either R2_ENDPOINT or R2_ACCOUNT_ID must be specified")
		}
		endpoint = fmt.Sprintf("https://%s.r2.cloudflarestorage.com", cfg.AccountID)
	}

	customResolver := s3.EndpointResolverV2(nil)
	_ = customResolver

	s3Client := s3.New(s3.Options{
		Region:       "auto",
		BaseEndpoint: aws.String(endpoint),
		Credentials:  credentials.NewStaticCredentialsProvider(cfg.AccessKey, cfg.SecretKey, ""),
	})

	presignClient := s3.NewPresignClient(s3Client)

	return &Client{
		s3Client:      s3Client,
		presignClient: presignClient,
		bucket:        cfg.Bucket,
	}, nil
}

func (c *Client) GeneratePresignedUploadURL(ctx context.Context, key string, contentType string, lifetime time.Duration) (string, error) {
	input := &s3.PutObjectInput{
		Bucket:      aws.String(c.bucket),
		Key:         aws.String(key),
		ContentType: aws.String(contentType),
	}

	req, err := c.presignClient.PresignPutObject(ctx, input, s3.WithPresignExpires(lifetime))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned upload url: %w", err)
	}

	return req.URL, nil
}

func (c *Client) GeneratePresignedDownloadURL(ctx context.Context, key string, filename string, lifetime time.Duration) (string, error) {
	contentDisposition := fmt.Sprintf("attachment; filename=%q", filename)

	input := &s3.GetObjectInput{
		Bucket:                     aws.String(c.bucket),
		Key:                        aws.String(key),
		ResponseContentDisposition: aws.String(contentDisposition),
	}

	req, err := c.presignClient.PresignGetObject(ctx, input, s3.WithPresignExpires(lifetime))
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned download url: %w", err)
	}

	return req.URL, nil
}

func (c *Client) HeadObject(ctx context.Context, key string) (*storage.ObjectMeta, error) {
	input := &s3.HeadObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	output, err := c.s3Client.HeadObject(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to head object %s: %w", key, err)
	}

	etag := ""
	if output.ETag != nil {
		etag = strings.Trim(*output.ETag, `"`)
	}

	contentType := ""
	if output.ContentType != nil {
		contentType = *output.ContentType
	}

	var lastModified time.Time
	if output.LastModified != nil {
		lastModified = *output.LastModified
	}

	var size int64
	if output.ContentLength != nil {
		size = *output.ContentLength
	}

	return &storage.ObjectMeta{
		Size:         size,
		ETag:         etag,
		ContentType:  contentType,
		LastModified: lastModified,
	}, nil
}

func (c *Client) DeleteObject(ctx context.Context, key string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(c.bucket),
		Key:    aws.String(key),
	}

	_, err := c.s3Client.DeleteObject(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete object %s: %w", key, err)
	}

	return nil
}
