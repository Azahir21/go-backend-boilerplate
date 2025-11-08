package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)

// GCSStorageConfig holds configuration for GCS storage.
type GCSStorageConfig struct {
	ProjectID       string `mapstructure:"project_id"`
	Bucket          string `mapstructure:"bucket"`
	CredentialsFile string `mapstructure:"credentials_file"`
}

// GCSStorage implements the Storage interface for Google Cloud Storage.
type GCSStorage struct {
	client *storage.Client
	bucket string
}

// NewGCSStorage creates a new GCSStorage instance.
func NewGCSStorage(ctx context.Context, cfg GCSStorageConfig) (*GCSStorage, error) {
	if cfg.ProjectID == "" || cfg.Bucket == "" || cfg.CredentialsFile == "" {
		return nil, fmt.Errorf("GCS configuration (ProjectID, Bucket, CredentialsFile) cannot be empty")
	}

	client, err := storage.NewClient(ctx, option.WithCredentialsFile(cfg.CredentialsFile))
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS client: %w", err)
	}

	return &GCSStorage{
		client: client,
		bucket: cfg.Bucket,
	}, nil
}

// Upload saves a file to GCS.
func (g *GCSStorage) Upload(ctx context.Context, bucketName, objectName string, file io.Reader) (string, error) {
	wc := g.client.Bucket(bucketName).Object(objectName).NewWriter(ctx)

	if _, err := io.Copy(wc, file); err != nil {
		return "", fmt.Errorf("failed to write file to GCS: %w", err)
	}

	if err := wc.Close(); err != nil {
		return "", fmt.Errorf("failed to close GCS writer: %w", err)
	}

	return fmt.Sprintf("gs://%s/%s", bucketName, objectName), nil
}

// Download retrieves a file from GCS.
func (g *GCSStorage) Download(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	rc, err := g.client.Bucket(bucketName).Object(objectName).NewReader(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create GCS reader: %w", err)
	}

	return rc, nil
}

// Delete removes a file from GCS.
func (g *GCSStorage) Delete(ctx context.Context, bucketName, objectName string) error {
	if err := g.client.Bucket(bucketName).Object(objectName).Delete(ctx); err != nil {
		return fmt.Errorf("failed to delete file from GCS: %w", err)
	}

	return nil
}

// GetSignedURL generates a signed URL for GCS.
func (g *GCSStorage) GetSignedURL(ctx context.Context, bucketName, objectName string, method string, durationMinutes int) (string, error) {
	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  method,
		Expires: time.Now().Add(time.Duration(durationMinutes) * time.Minute),
	}

	url, err := g.client.Bucket(bucketName).SignedURL(objectName, opts)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL for GCS: %w", err)
	}

	return url, nil
}
