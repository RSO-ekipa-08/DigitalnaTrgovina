package service

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"time"

	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/database/generated"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/pgutil"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/repository"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

// serviceImpl implements Service interface
type serviceImpl struct {
	repo    repository.ApplicationRepository
	storage types.StorageService
}

// Make sure serviceImpl implements types.ApplicationService interface
var _ types.ApplicationService = (*serviceImpl)(nil)

// New creates a new service instance
func New(repo repository.ApplicationRepository, storage types.StorageService) types.ApplicationService {
	return &serviceImpl{
		repo:    repo,
		storage: storage,
	}
}

// CreateApplication creates a new application
func (s *serviceImpl) CreateApplication(ctx context.Context, params types.CreateApplicationParams) (*database.CreateApplicationRow, error) {
	// Generate storage URL for APK file
	objectName := fmt.Sprintf("apps/%s/%s.apk", params.DeveloperID, uuid.New())

	// Upload APK file
	reader := bytes.NewReader(params.APKFile)
	if err := s.storage.UploadFile(ctx, objectName, reader, int64(len(params.APKFile))); err != nil {
		return nil, fmt.Errorf("failed to upload APK file: %w", err)
	}

	// Parse developer ID
	developerID, err := pgutil.UUIDFromString(params.DeveloperID)
	if err != nil {
		// Cleanup storage if parsing fails
		if err := s.storage.DeleteFile(ctx, objectName); err != nil {
			log.Error().Err(err).Str("object", objectName).Msg("failed to cleanup storage after failed ID parsing")
		}
		return nil, fmt.Errorf("invalid developer ID: %w", err)
	}

	// Create application in database
	dbParams := &database.CreateApplicationParams{
		Name:              params.Name,
		Description:       pgutil.String(params.Description),
		DeveloperID:       developerID,
		Category:          params.Category,
		Price:             pgutil.Float64(params.Price),
		Size:              params.Size,
		MinAndroidVersion: params.MinAndroidVersion,
		CurrentVersion:    params.CurrentVersion,
		Tags:              params.Tags,
		Screenshots:       params.Screenshots,
		StorageUrl:        path.Join("apps", params.DeveloperID, objectName),
	}

	app, err := s.repo.CreateApplication(ctx, dbParams)
	if err != nil {
		// Cleanup storage if database operation fails
		if err := s.storage.DeleteFile(ctx, objectName); err != nil {
			log.Error().Err(err).Str("object", objectName).Msg("failed to cleanup storage after failed application creation")
		}
		return nil, fmt.Errorf("failed to create application: %w", err)
	}

	return app, nil
}

// GetApplication retrieves an application by ID
func (s *serviceImpl) GetApplication(ctx context.Context, id uuid.UUID) (*database.GetApplicationRow, error) {
	app, err := s.repo.GetApplication(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}
	return app, nil
}

// UpdateApplication updates an existing application
func (s *serviceImpl) UpdateApplication(ctx context.Context, params types.UpdateApplicationParams) (*database.UpdateApplicationRow, error) {
	// Get existing application
	app, err := s.repo.GetApplication(ctx, params.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to get application: %w", err)
	}

	// Handle APK file update if provided
	var storageURL *string
	if len(params.APKFile) > 0 {
		objectName := path.Join("apps", pgutil.StringFromUUID(app.DeveloperID), uuid.New().String()+".apk")
		reader := bytes.NewReader(params.APKFile)
		if err := s.storage.UploadFile(ctx, objectName, reader, int64(len(params.APKFile))); err != nil {
			return nil, fmt.Errorf("failed to upload APK file: %w", err)
		}
		url := path.Join("apps", pgutil.StringFromUUID(app.DeveloperID), objectName)
		storageURL = &url

		// Delete old APK file
		if err := s.storage.DeleteFile(ctx, app.StorageUrl); err != nil {
			log.Error().Err(err).Str("object", app.StorageUrl).Msg("failed to delete old APK file")
		}
	}

	// Update application in database
	dbParams := &database.UpdateApplicationParams{
		ID:                pgutil.UUID(params.ID),
		Name:              pgutil.StringPtr(params.Name),
		Description:       pgutil.StringPtr(params.Description),
		Category:          pgutil.StringPtr(params.Category),
		Price:             pgutil.Float64Ptr(params.Price),
		MinAndroidVersion: pgutil.StringPtr(params.MinAndroidVersion),
		CurrentVersion:    pgutil.StringPtr(params.CurrentVersion),
		Tags:              *params.Tags,
		Screenshots:       *params.Screenshots,
		StorageUrl:        pgutil.StringPtr(storageURL),
	}

	app1, err := s.repo.UpdateApplication(ctx, dbParams)
	if err != nil {
		return nil, fmt.Errorf("failed to update application: %w", err)
	}

	return app1, nil
}

// DeleteApplication deletes an application
func (s *serviceImpl) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	// Get application to get storage URL
	app, err := s.repo.GetApplication(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to get application: %w", err)
	}

	// Delete from database
	if err := s.repo.DeleteApplication(ctx, id); err != nil {
		return fmt.Errorf("failed to delete application: %w", err)
	}

	// Delete APK file
	if err := s.storage.DeleteFile(ctx, app.StorageUrl); err != nil {
		log.Error().Err(err).Str("object", app.StorageUrl).Msg("failed to delete APK file")
	}

	return nil
}

// SearchApplications searches for applications based on criteria
func (s *serviceImpl) SearchApplications(ctx context.Context, params types.SearchApplicationsParams) ([]*database.SearchApplicationsRow, error) {
	dbParams := &database.SearchApplicationsParams{
		Query:             params.Query,
		Category:          *params.Category,
		MinPrice:          pgutil.Float64Ptr(params.MinPrice),
		MaxPrice:          pgutil.Float64Ptr(params.MaxPrice),
		MinAndroidVersion: *params.MinAndroidVersion,
		Tags:              params.Tags,
		Limit:             params.Limit,
		Offset:            params.Offset,
		SortByDownloads:   params.SortByDownloads,
		SortByRating:      params.SortByRating,
	}

	apps, err := s.repo.SearchApplications(ctx, dbParams)
	if err != nil {
		return nil, fmt.Errorf("failed to search applications: %w", err)
	}

	return apps, nil
}

// GetDownloadURL generates a download URL for an application
func (s *serviceImpl) GetDownloadURL(ctx context.Context, appID, userID uuid.UUID, ipAddress string) (string, time.Time, error) {
	// Get application
	app, err := s.repo.GetApplication(ctx, appID)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to get application: %w", err)
	}

	// Generate presigned URL
	url, expiresAt, err := s.storage.GetDownloadURL(ctx, app.StorageUrl)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to generate download URL: %w", err)
	}

	// Record download
	downloadParams := &database.CreateDownloadParams{
		UserID:        pgutil.UUID(userID),
		ApplicationID: pgutil.UUID(appID),
		IpAddress:     ipAddress,
		Success:       true,
	}

	if _, err := s.repo.CreateDownload(ctx, downloadParams); err != nil {
		log.Error().Err(err).Msg("failed to record download")
	}

	// Increment download count
	if _, err := s.repo.IncrementDownloads(ctx, appID); err != nil {
		log.Error().Err(err).Msg("failed to increment download count")
	}

	return url, expiresAt, nil
}

// ListCategories lists all available categories
func (s *serviceImpl) ListCategories(ctx context.Context, params types.ListCategoriesParams) ([]*database.Category, error) {
	dbParams := &database.ListCategoriesParams{
		Limit:  params.Limit,
		Offset: params.Offset,
	}

	categories, err := s.repo.ListCategories(ctx, dbParams)
	if err != nil {
		return nil, fmt.Errorf("failed to list categories: %w", err)
	}

	return categories, nil
}
