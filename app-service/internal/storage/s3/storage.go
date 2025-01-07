package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Storage represents an object storage client
type Storage struct {
	client     *s3.Client
	bucketName string
}

// New creates a new object storage client
func New(cfg *config.Config) (*Storage, error) {
	// Initialize S3 client
	s3Client := s3.New(s3.Options{
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(
			cfg.StorageAccessKey,
			cfg.StorageSecretKey,
			"", // session token
		)),
		BaseEndpoint: aws.String(cfg.StorageEndpoint),
	})

	// Check if bucket exists
	_, err := s3Client.HeadBucket(context.Background(), &s3.HeadBucketInput{
		Bucket: aws.String(cfg.StorageBucketName),
	})
	if err != nil {
		// Create bucket if it doesn't exist
		_, err = s3Client.CreateBucket(context.Background(), &s3.CreateBucketInput{
			Bucket: aws.String(cfg.StorageBucketName),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &Storage{
		client:     s3Client,
		bucketName: cfg.StorageBucketName,
	}, nil
}

// UploadFile uploads a file to object storage
func (s *Storage) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectName),
		Body:   reader,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

// GetDownloadURL generates a presigned URL for downloading a file
func (s *Storage) GetDownloadURL(ctx context.Context, objectName string) (string, time.Time, error) {
	presignClient := s3.NewPresignClient(s.client)
	presignDuration := 15 * time.Minute

	presignResult, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectName),
	}, s3.WithPresignExpires(presignDuration))

	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate download URL: %w", err)
	}

	expiry := time.Now().Add(presignDuration)
	return presignResult.URL, expiry, nil
}

// DeleteFile deletes a file from object storage
func (s *Storage) DeleteFile(ctx context.Context, objectName string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucketName),
		Key:    aws.String(objectName),
	})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
