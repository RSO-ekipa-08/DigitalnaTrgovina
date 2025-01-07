package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/config"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/types"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// Storage represents an object storage client
type Storage struct {
	client     *minio.Client
	bucketName string
}

// Make sure Storage implements types.StorageService interface
var _ types.StorageService = (*Storage)(nil)

// New creates a new object storage client
func New(cfg *config.Config) (*Storage, error) {
	// Initialize MinIO client
	client, err := minio.New(cfg.StorageEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.StorageAccessKey, cfg.StorageSecretKey, ""),
		Secure: cfg.StorageUseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create storage client: %w", err)
	}

	// Check if bucket exists
	exists, err := client.BucketExists(context.Background(), cfg.StorageBucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	// Create bucket if it doesn't exist
	if !exists {
		err = client.MakeBucket(context.Background(), cfg.StorageBucketName, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &Storage{
		client:     client,
		bucketName: cfg.StorageBucketName,
	}, nil
}

// UploadFile uploads a file to object storage
func (s *Storage) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64) error {
	_, err := s.client.PutObject(ctx, s.bucketName, objectName, reader, size, minio.PutObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

// GetDownloadURL generates a presigned URL for downloading a file
func (s *Storage) GetDownloadURL(ctx context.Context, objectName string) (string, time.Time, error) {
	// Get presigned URL
	url, err := s.client.PresignedGetObject(ctx, s.bucketName, objectName, 15*time.Minute, nil)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate download URL: %w", err)
	}
	return url.String(), time.Now().Add(15 * time.Minute), nil
}

// DeleteFile deletes a file from object storage
func (s *Storage) DeleteFile(ctx context.Context, objectName string) error {
	err := s.client.RemoveObject(ctx, s.bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
