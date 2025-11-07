package storage

import (
	"context"
	"io"
)

// Storage defines the interface for file storage operations.
type Storage interface {
	Upload(ctx context.Context, bucketName, objectName string, file io.Reader) (string, error)
	Download(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error)
	Delete(ctx context.Context, bucketName, objectName string) error
	GetSignedURL(ctx context.Context, bucketName, objectName string, method string, durationMinutes int) (string, error)
}
