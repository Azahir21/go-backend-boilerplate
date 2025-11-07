package storage

import (
	"context"
	"fmt"

	"github.com/azahir21/go-backend-boilerplate/pkg/config"
	"github.com/sirupsen/logrus"
)

// NewStorage creates a new Storage implementation based on the provided configuration.
func NewStorage(ctx context.Context, log *logrus.Logger, cfg config.StorageConfig) (Storage, error) {
	switch cfg.Type {
	case "local":
		log.Info("Initializing Local Storage...")
		return NewLocalStorage(LocalStorageConfig{
			BasePath: cfg.Local.BasePath,
		})
	case "s3":
		log.Info("Initializing S3 Storage...")
		return NewS3Storage(S3StorageConfig{
			Region:          cfg.S3.Region,
			AccessKeyID:     cfg.S3.AccessKeyID,
			SecretAccessKey: cfg.S3.SecretAccessKey,
			Bucket:          cfg.S3.Bucket,
		})
	case "gcs":
		log.Info("Initializing GCS Storage...")
		return NewGCSStorage(ctx, GCSStorageConfig{
			ProjectID:       cfg.GCS.ProjectID,
			Bucket:          cfg.GCS.Bucket,
			CredentialsFile: cfg.GCS.CredentialsFile,
		})
	default:
		return nil, fmt.Errorf("unsupported storage type: %s", cfg.Type)
	}
}
