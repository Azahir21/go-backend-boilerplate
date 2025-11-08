package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// S3StorageConfig holds configuration for S3 storage.
type S3StorageConfig struct {
	Region          string `mapstructure:"region"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	SecretAccessKey string `mapstructure:"secret_access_key"`
	Bucket          string `mapstructure:"bucket"`
}

// S3Storage implements the Storage interface for AWS S3.
type S3Storage struct {
	s3Client *s3.S3
	bucket   string
}

// NewS3Storage creates a new S3Storage instance.
func NewS3Storage(cfg S3StorageConfig) (*S3Storage, error) {
	if cfg.Region == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" || cfg.Bucket == "" {
		return nil, fmt.Errorf("S3 configuration (Region, AccessKeyID, SecretAccessKey, Bucket) cannot be empty")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.Region),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create AWS session: %w", err)
	}

	return &S3Storage{
		s3Client: s3.New(sess),
		bucket:   cfg.Bucket,
	}, nil
}

// Upload saves a file to S3.
func (s *S3Storage) Upload(ctx context.Context, bucketName, objectName string, file io.Reader) (string, error) {
	uploader := s3manager.NewUploaderWithClient(s.s3Client)

	input := &s3manager.UploadInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
		Body:   file,
	}

	result, err := uploader.UploadWithContext(ctx, input)
	if err != nil {
		return "", fmt.Errorf("failed to upload file to S3: %w", err)
	}

	return result.Location, nil
}

// Download retrieves a file from S3.
func (s *S3Storage) Download(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	input := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}

	result, err := s.s3Client.GetObjectWithContext(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to download file from S3: %w", err)
	}

	return result.Body, nil
}

// Delete removes a file from S3.
func (s *S3Storage) Delete(ctx context.Context, bucketName, objectName string) error {
	input := &s3.DeleteObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	}

	_, err := s.s3Client.DeleteObjectWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to delete file from S3: %w", err)
	}

	return nil
}

// GetSignedURL generates a signed URL for S3.
func (s *S3Storage) GetSignedURL(ctx context.Context, bucketName, objectName string, method string, durationMinutes int) (string, error) {
	req, _ := s.s3Client.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectName),
	})

	urlStr, err := req.Presign(time.Duration(durationMinutes) * time.Minute)
	if err != nil {
		return "", fmt.Errorf("failed to generate signed URL for S3: %w", err)
	}

	return urlStr, nil
}
