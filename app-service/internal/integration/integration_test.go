//go:build integration
// +build integration

package integration

import (
	"context"
	"os"
	"testing"
	"time"

	appv1 "github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/gen/app/v1"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/config"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/database"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/handler"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/repository"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/service"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/storage"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testHandler *handler.Handler

func TestMain(m *testing.M) {
	// Load configuration
	cfg := &config.Config{
		DatabaseURL:      os.Getenv("TEST_DATABASE_URL"),
		StorageEndpoint:  os.Getenv("STORAGE_ENDPOINT"),
		StorageAccessKey: os.Getenv("STORAGE_ACCESS_KEY"),
		StorageSecretKey: os.Getenv("STORAGE_SECRET_KEY"),
		StorageUseSSL:    false,
	}

	// Initialize database
	db, err := database.New(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	// Initialize storage
	stor, err := storage.New(cfg)
	if err != nil {
		panic(err)
	}

	// Initialize repository
	repo := repository.New(db.Pool)

	// Initialize service
	svc := service.New(repo, stor)

	// Initialize handler
	testHandler = handler.New(svc)

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func TestCreateAndGetApplication(t *testing.T) {
	ctx := context.Background()
	developerID := uuid.New()

	// Create application
	createReq := &appv1.CreateApplicationRequest{
		Name:              "Test App",
		Description:       "Test Description",
		DeveloperId:       developerID.String(),
		Category:          "Games",
		Price:             9.99,
		Size:              1024 * 1024 * 100,
		MinAndroidVersion: "8.0",
		CurrentVersion:    "1.0.0",
		Tags:              []string{"game", "action"},
		Screenshots:       []string{"screenshot1.jpg"},
		ApkFile:           []byte("test apk file"),
	}

	app, err := testHandler.CreateApplication(ctx, createReq)
	require.NoError(t, err)
	assert.NotEmpty(t, app.Id)
	assert.Equal(t, createReq.Name, app.Name)
	assert.Equal(t, createReq.Description, app.Description)
	assert.Equal(t, createReq.DeveloperId, app.DeveloperId)
	assert.Equal(t, createReq.Category, app.Category)
	assert.Equal(t, createReq.Price, app.Price)

	// Get application
	getReq := &appv1.GetApplicationRequest{
		Id: app.Id,
	}

	fetchedApp, err := testHandler.GetApplication(ctx, getReq)
	require.NoError(t, err)
	assert.Equal(t, app.Id, fetchedApp.Id)
	assert.Equal(t, app.Name, fetchedApp.Name)
}

func TestSearchApplications(t *testing.T) {
	ctx := context.Background()
	developerID := uuid.New()

	// Create test applications
	app1, err := testHandler.CreateApplication(ctx, &appv1.CreateApplicationRequest{
		Name:              "Action Game",
		Description:       "An exciting action game",
		DeveloperId:       developerID.String(),
		Category:          "Games",
		Price:             9.99,
		Size:              1024 * 1024 * 100,
		MinAndroidVersion: "8.0",
		CurrentVersion:    "1.0.0",
		Tags:              []string{"game", "action"},
		Screenshots:       []string{"screenshot1.jpg"},
		ApkFile:           []byte("test apk file 1"),
	})
	require.NoError(t, err)

	app2, err := testHandler.CreateApplication(ctx, &appv1.CreateApplicationRequest{
		Name:              "Puzzle Game",
		Description:       "A brain teasing puzzle game",
		DeveloperId:       developerID.String(),
		Category:          "Games",
		Price:             4.99,
		Size:              1024 * 1024 * 50,
		MinAndroidVersion: "7.0",
		CurrentVersion:    "1.0.0",
		Tags:              []string{"game", "puzzle"},
		Screenshots:       []string{"screenshot1.jpg"},
		ApkFile:           []byte("test apk file 2"),
	})
	require.NoError(t, err)

	// Search for action games
	searchReq := &appv1.SearchApplicationsRequest{
		Query: "action",
		Pagination: &appv1.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	}

	results, err := testHandler.SearchApplications(ctx, searchReq)
	require.NoError(t, err)
	assert.NotEmpty(t, results.Applications)
	assert.Equal(t, app1.Id, results.Applications[0].Id)

	// Search for puzzle games
	searchReq.Query = "puzzle"
	results, err = testHandler.SearchApplications(ctx, searchReq)
	require.NoError(t, err)
	assert.NotEmpty(t, results.Applications)
	assert.Equal(t, app2.Id, results.Applications[0].Id)
}

func TestDownloadApplication(t *testing.T) {
	ctx := context.Background()
	developerID := uuid.New()
	userID := uuid.New()

	// Create application
	app, err := testHandler.CreateApplication(ctx, &appv1.CreateApplicationRequest{
		Name:              "Test App",
		Description:       "Test Description",
		DeveloperId:       developerID.String(),
		Category:          "Games",
		Price:             9.99,
		Size:              1024 * 1024 * 100,
		MinAndroidVersion: "8.0",
		CurrentVersion:    "1.0.0",
		Tags:              []string{"game", "action"},
		Screenshots:       []string{"screenshot1.jpg"},
		ApkFile:           []byte("test apk file"),
	})
	require.NoError(t, err)

	// Get download URL
	downloadReq := &appv1.DownloadApplicationRequest{
		ApplicationId: app.Id,
		UserId:        userID.String(),
	}

	downloadResp, err := testHandler.DownloadApplication(ctx, downloadReq)
	require.NoError(t, err)
	assert.NotEmpty(t, downloadResp.DownloadUrl)
	assert.True(t, downloadResp.ExpiresAt.AsTime().After(time.Now()))
}
