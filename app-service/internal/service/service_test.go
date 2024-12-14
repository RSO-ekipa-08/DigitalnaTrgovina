package service

import (
	"context"
	"io"
	"testing"
	"time"

	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/database/generated"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/pgutil"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/repository"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// Mock repository
type mockRepository struct {
	mock.Mock
	repository.ApplicationRepository
}

// Ensure mockRepository implements repository.ApplicationRepository interface
var _ repository.ApplicationRepository = (*mockRepository)(nil)

// Mock storage
type mockStorage struct {
	mock.Mock
}

// Ensure mockStorage implements types.StorageService interface
var _ types.StorageService = (*mockStorage)(nil)

func (m *mockStorage) UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64) error {
	args := m.Called(ctx, objectName, reader, size)
	return args.Error(0)
}

func (m *mockStorage) GetDownloadURL(ctx context.Context, objectName string) (string, time.Time, error) {
	args := m.Called(ctx, objectName)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *mockStorage) DeleteFile(ctx context.Context, objectName string) error {
	args := m.Called(ctx, objectName)
	return args.Error(0)
}

// Repository mock methods
func (m *mockRepository) GetApplication(ctx context.Context, id uuid.UUID) (*database.GetApplicationRow, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.GetApplicationRow), args.Error(1)
}

func (m *mockRepository) ListApplications(ctx context.Context, limit, offset int32) ([]*database.ListApplicationsRow, error) {
	args := m.Called(ctx, limit, offset)
	return args.Get(0).([]*database.ListApplicationsRow), args.Error(1)
}

func (m *mockRepository) SearchApplications(ctx context.Context, params *database.SearchApplicationsParams) ([]*database.SearchApplicationsRow, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]*database.SearchApplicationsRow), args.Error(1)
}

func (m *mockRepository) CreateApplication(ctx context.Context, params *database.CreateApplicationParams) (*database.CreateApplicationRow, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*database.CreateApplicationRow), args.Error(1)
}

func (m *mockRepository) UpdateApplication(ctx context.Context, params *database.UpdateApplicationParams) (*database.UpdateApplicationRow, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*database.UpdateApplicationRow), args.Error(1)
}

func (m *mockRepository) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *mockRepository) IncrementDownloads(ctx context.Context, id uuid.UUID) (*database.IncrementDownloadsRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(*database.IncrementDownloadsRow), args.Error(1)
}

func (m *mockRepository) ListCategories(ctx context.Context, params *database.ListCategoriesParams) ([]*database.Category, error) {
	args := m.Called(ctx, params)
	return args.Get(0).([]*database.Category), args.Error(1)
}

func (m *mockRepository) CreateDownload(ctx context.Context, params *database.CreateDownloadParams) (*database.Download, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(*database.Download), args.Error(1)
}

func (m *mockRepository) GetDownloadsByUser(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*database.Download, error) {
	args := m.Called(ctx, userID, limit, offset)
	return args.Get(0).([]*database.Download), args.Error(1)
}

func TestCreateApplication(t *testing.T) {
	repo := new(mockRepository)
	stor := new(mockStorage)
	svc := New(repo, stor)

	ctx := context.Background()
	developerID := uuid.New()

	// Setup mock expectations
	expectedApp := &database.CreateApplicationRow{
		ID:          pgutil.UUID(uuid.New()),
		Name:        "Test App",
		Description: pgutil.String("Test Description"),
		DeveloperID: pgutil.UUID(developerID),
		Category:    "Games",
		Price:       pgutil.Float64(9.99),
		Size:        1024,
		StorageUrl:  "apps/" + developerID.String() + "/test.apk",
	}

	stor.On("UploadFile", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return(nil)
	repo.On("CreateApplication", mock.Anything, mock.Anything).Return(expectedApp, nil)

	// Test
	app, err := svc.CreateApplication(ctx, types.CreateApplicationParams{
		Name:              "Test App",
		Description:       "Test Description",
		DeveloperID:       developerID.String(),
		Category:          "Games",
		Price:             9.99,
		Size:              1024,
		MinAndroidVersion: "8.0",
		CurrentVersion:    "1.0.0",
		Tags:              []string{"game", "arcade"},
		Screenshots:       []string{"screen1.jpg", "screen2.jpg"},
		APKFile:           []byte("test"),
	})

	require.NoError(t, err)
	assert.Equal(t, expectedApp.Name, app.Name)
	assert.Equal(t, expectedApp.DeveloperID, app.DeveloperID)
	assert.Equal(t, expectedApp.Category, app.Category)
	assert.Equal(t, expectedApp.Price, app.Price)

	stor.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestGetDownloadURL(t *testing.T) {
	repo := new(mockRepository)
	stor := new(mockStorage)
	svc := New(repo, stor)

	ctx := context.Background()
	appID := uuid.New()
	userID := uuid.New()

	// Setup mock expectations
	app := &database.Application{
		ID:         pgutil.UUID(appID),
		StorageUrl: "apps/test.apk",
	}
	download := &database.Download{
		ID:            pgutil.UUID(uuid.New()),
		UserID:        pgutil.UUID(userID),
		ApplicationID: pgutil.UUID(appID),
	}
	downloadURL := "https://storage.com/download/test.apk"
	expiresAt := time.Now().Add(time.Hour)

	repo.On("GetApplication", mock.Anything, appID).Return(app, nil)
	stor.On("GetDownloadURL", mock.Anything, app.StorageUrl).Return(downloadURL, expiresAt, nil)
	repo.On("CreateDownload", mock.Anything, mock.Anything).Return(download, nil)
	repo.On("IncrementDownloads", mock.Anything, appID).Return(app, nil)

	// Test
	url, expires, err := svc.GetDownloadURL(ctx, appID, userID, "127.0.0.1")

	require.NoError(t, err)
	assert.Equal(t, downloadURL, url)
	assert.Equal(t, expiresAt, expires)

	stor.AssertExpectations(t)
	repo.AssertExpectations(t)
}

func TestSearchApplications(t *testing.T) {
	repo := new(mockRepository)
	stor := new(mockStorage)
	svc := New(repo, stor)

	ctx := context.Background()

	// Setup mock expectations
	expectedApps := []*database.Application{
		{
			ID:          pgutil.UUID(uuid.New()),
			Name:        "Test App 1",
			Description: pgutil.String("Description 1"),
			Price:       pgutil.Float64(9.99),
			Rating:      pgutil.Float64(4.5),
			Downloads:   1000,
		},
		{
			ID:          pgutil.UUID(uuid.New()),
			Name:        "Test App 2",
			Description: pgutil.String("Description 2"),
			Price:       pgutil.Float64(19.99),
			Rating:      pgutil.Float64(4.8),
			Downloads:   2000,
		},
	}

	repo.On("SearchApplications", mock.Anything, mock.Anything).Return(expectedApps, nil)

	// Test
	apps, err := svc.SearchApplications(ctx, types.SearchApplicationsParams{
		Query:           "test",
		Category:        stringPtr("Games"),
		MinPrice:        float64Ptr(5.0),
		MaxPrice:        float64Ptr(20.0),
		SortByDownloads: true,
		Limit:           10,
		Offset:          0,
	})

	require.NoError(t, err)
	assert.Len(t, apps, 2)
	assert.Equal(t, expectedApps[0].Name, apps[0].Name)
	assert.Equal(t, expectedApps[1].Name, apps[1].Name)

	repo.AssertExpectations(t)
}

// Helper functions
func stringPtr(s string) *string {
	return &s
}

func float64Ptr(f float64) *float64 {
	return &f
}
