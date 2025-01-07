package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/config"
	storage "github.com/supabase-community/storage-go"
)

type Storage struct {
	client     *storage.Client
	bucketName string
}

func New(cfg *config.Config) (*Storage, error) {
	// Inicializacija Supabase Storage clienta
	client := storage.NewClient(
		cfg.StorageEndpoint,
		cfg.StorageSecretKey,
		nil, // dodatne opcije lahko pustimo prazne
	)

	// Preveri, če bucket obstaja, če ne, ga ustvari
	_, err := client.GetBucket(cfg.StorageBucketName)
	if err != nil {
		_, err = client.CreateBucket(cfg.StorageBucketName, storage.BucketOptions{
			Public: false, // nastavimo na false za večjo varnost
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}
	}

	return &Storage{
		client:     client,
		bucketName: cfg.StorageBucketName,
	}, nil
}

// UploadFile naloži datoteko v storage
func (s *Storage) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64) error {
	_, err := s.client.UploadFile(s.bucketName, objectName, reader)
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}
	return nil
}

// GetDownloadURL generira URL za prenos datoteke
func (s *Storage) GetDownloadURL(ctx context.Context, objectName string) (string, time.Time, error) {
	// Nastavimo, da URL poteče čez 15 minut
	expiresIn := 60 * 15 // 15 minut v sekundah

	result, err := s.client.CreateSignedUrl(s.bucketName, objectName, expiresIn)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate download URL: %w", err)
	}

	expiry := time.Now().Add(time.Duration(expiresIn) * time.Second)
	return result.SignedURL, expiry, nil
}

// DeleteFile izbriše datoteko iz storage-a
func (s *Storage) DeleteFile(ctx context.Context, objectName string) error {
	_, err := s.client.RemoveFile(s.bucketName, []string{objectName})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}
