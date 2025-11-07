package storage

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// LocalStorageConfig holds configuration for local storage.
type LocalStorageConfig struct {
	BasePath string `mapstructure:"base_path"`
}

// LocalStorage implements the Storage interface for local file system.
type LocalStorage struct {
	basePath string
}

// NewLocalStorage creates a new LocalStorage instance.
func NewLocalStorage(cfg LocalStorageConfig) (*LocalStorage, error) {
	if cfg.BasePath == "" {
		return nil, fmt.Errorf("base_path cannot be empty for local storage")
	}

	// Ensure the base path exists
	if err := os.MkdirAll(cfg.BasePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create base path for local storage: %w", err)
	}

	return &LocalStorage{
			basePath: cfg.BasePath,
		},
		nil
}

// Upload saves a file to the local file system.
func (l *LocalStorage) Upload(ctx context.Context, bucketName, objectName string, file io.Reader) (string, error) {
	// For local storage, bucketName can be used as a sub-directory or ignored.
	dirPath := filepath.Join(l.basePath, bucketName)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return "", fmt.Errorf("failed to create directory %s: %w", dirPath, err)
	}

	filePath := filepath.Join(dirPath, objectName)
	out, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file %s: %w", filePath, err)
	}
	defer out.Close()

	_, err = io.Copy(out, file)
	if err != nil {
		return "", fmt.Errorf("failed to write file %s: %w", filePath, err)
	}

	return filePath, nil
}

// Download retrieves a file from the local file system.
func (l *LocalStorage) Download(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	filePath := filepath.Join(l.basePath, bucketName, objectName)

	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}

	return file, nil
}

// Delete removes a file from the local file system.
func (l *LocalStorage) Delete(ctx context.Context, bucketName, objectName string) error {
	filePath := filepath.Join(l.basePath, bucketName, objectName)

	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("failed to delete file %s: %w", filePath, err)
	}

	return nil
}

// GetSignedURL generates a signed URL for local storage. This is a placeholder as local storage doesn't inherently support signed URLs.
func (l *LocalStorage) GetSignedURL(ctx context.Context, bucketName, objectName string, method string, durationMinutes int) (string, error) {
	// For local storage, we might just return the direct path or an error if signed URLs are not applicable.
	// Depending on the use case, you might want to serve these files via an HTTP server and generate a URL to that server.
	// For simplicity, returning a direct path here.
	filePath := filepath.Join(l.basePath, bucketName, objectName)
	// In a real application, you might want to return an error or a URL to a local web server that serves this file.
	return filePath, nil
}
