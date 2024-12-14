package repository

import (
	"context"
	"os"
	"path"
	"testing"

	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/config"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/database"
	dbgen "github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/database/generated"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/pgutil"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testDB *database.DB
var testRepo *Repository

func TestMain(m *testing.M) {
	// Set up test database connection
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	testDB, err = database.New(context.Background(), cfg)
	if err != nil {
		panic(err)
	}
	defer testDB.Close()

	testRepo = New(testDB.Pool)

	// Run tests
	code := m.Run()

	os.Exit(code)
}

func TestCreateAndGetApplication(t *testing.T) {
	ctx := context.Background()

	// Create test application
	developerID := uuid.New()
	app, err := testRepo.CreateApplication(ctx, &dbgen.CreateApplicationParams{
		Name:              "Test App",
		Description:       pgutil.String("Test Description"),
		DeveloperID:       pgutil.UUID(developerID),
		Category:          "Games",
		Price:             pgutil.Float64(9.99),
		Size:              1024 * 1024 * 100, // 100MB
		MinAndroidVersion: "8.0",
		CurrentVersion:    "1.0.0",
		Tags:              []string{"game", "action"},
		Screenshots:       []string{"screenshot1.jpg", "screenshot2.jpg"},
		StorageUrl:        path.Join("apps", developerID.String(), "test.apk"),
	})

	require.NoError(t, err)
	assert.NotNil(t, app)
	assert.NotEqual(t, uuid.Nil, app.ID)
	assert.Equal(t, "Test App", app.Name)
	assert.Equal(t, developerID, pgutil.FromUUID(app.DeveloperID))
	assert.Equal(t, float64(9.99), pgutil.FromFloat64(app.Price))
	assert.Equal(t, []string{"game", "action"}, app.Tags)

	// Get application
	fetchedApp, err := testRepo.GetApplication(ctx, app.ID.Bytes)
	require.NoError(t, err)
	assert.NotNil(t, fetchedApp)
	assert.Equal(t, app.ID, fetchedApp.ID)
	assert.Equal(t, app.Name, fetchedApp.Name)
}

func TestUpdateApplication(t *testing.T) {
	ctx := context.Background()

	// Create test application
	developerID := uuid.New()
	app, err := testRepo.CreateApplication(ctx, &dbgen.CreateApplicationParams{
		Name:              "Test App",
		Description:       pgutil.String("Test Description"),
		DeveloperID:       pgutil.UUID(developerID),
		Category:          "Games",
		Price:             pgutil.Float64(9.99),
		Size:              1024 * 1024 * 100,
		MinAndroidVersion: "8.0",
		CurrentVersion:    "1.0.0",
		Tags:              []string{"game", "action"},
		Screenshots:       []string{"screenshot1.jpg"},
		StorageUrl:        path.Join("apps", developerID.String(), "test.apk"),
	})
	require.NoError(t, err)

	// Update application
	newName := "Updated App"
	newPrice := 19.99
	updatedApp, err := testRepo.UpdateApplication(ctx, &dbgen.UpdateApplicationParams{
		ID:    app.ID,
		Name:  pgutil.String(newName),
		Price: pgutil.Float64(newPrice),
	})

	require.NoError(t, err)
	assert.NotNil(t, updatedApp)
	assert.Equal(t, newName, updatedApp.Name)
	assert.Equal(t, float64(19.99), pgutil.FromFloat64(updatedApp.Price))
	assert.Equal(t, app.Description, updatedApp.Description) // Unchanged fields should remain the same
}

func TestSearchApplications(t *testing.T) {
	ctx := context.Background()

	// Create test applications
	developerID := uuid.New()
	_, err := testRepo.CreateApplication(ctx, &dbgen.CreateApplicationParams{
		Name:              "Action Game 1",
		Description:       pgutil.String("An exciting action game"),
		DeveloperID:       pgutil.UUID(developerID),
		Category:          "Games",
		Price:             pgutil.Float64(9.99),
		Size:              1024 * 1024 * 100,
		MinAndroidVersion: "8.0",
		CurrentVersion:    "1.0.0",
		Tags:              []string{"game", "action"},
		Screenshots:       []string{"screenshot1.jpg"},
		StorageUrl:        path.Join("apps", developerID.String(), "game1.apk"),
	})
	require.NoError(t, err)

	_, err = testRepo.CreateApplication(ctx, &dbgen.CreateApplicationParams{
		Name:              "Puzzle Game",
		Description:       pgutil.String("A brain teasing puzzle game"),
		DeveloperID:       pgutil.UUID(developerID),
		Category:          "Games",
		Price:             pgutil.Float64(4.99),
		Size:              1024 * 1024 * 50,
		MinAndroidVersion: "7.0",
		CurrentVersion:    "1.0.0",
		Tags:              []string{"game", "puzzle"},
		Screenshots:       []string{"screenshot1.jpg"},
		StorageUrl:        path.Join("apps", developerID.String(), "game2.apk"),
	})
	require.NoError(t, err)

	// Search for action games
	results, err := testRepo.SearchApplications(ctx, &dbgen.SearchApplicationsParams{
		Query:  "action",
		Limit:  10,
		Offset: 0,
	})

	require.NoError(t, err)
	assert.NotEmpty(t, results)
	assert.Contains(t, results[0].Name, "Action")
}

func TestDeleteApplication(t *testing.T) {
	ctx := context.Background()

	// Create test application
	developerID := uuid.New()
	app, err := testRepo.CreateApplication(ctx, &dbgen.CreateApplicationParams{
		Name:              "Test App",
		Description:       pgutil.String("Test Description"),
		DeveloperID:       pgutil.UUID(developerID),
		Category:          "Games",
		Price:             pgutil.Float64(9.99),
		Size:              1024 * 1024 * 100,
		MinAndroidVersion: "8.0",
		CurrentVersion:    "1.0.0",
		Tags:              []string{"game", "action"},
		Screenshots:       []string{"screenshot1.jpg"},
		StorageUrl:        path.Join("apps", developerID.String(), "test.apk"),
	})
	require.NoError(t, err)

	// Delete application
	err = testRepo.DeleteApplication(ctx, app.ID.Bytes)
	require.NoError(t, err)

	// Try to get deleted application
	_, err = testRepo.GetApplication(ctx, app.ID.Bytes)
	assert.Error(t, err)
}

func TestDownloads(t *testing.T) {
	ctx := context.Background()

	// Create test application
	developerID := uuid.New()
	app, err := testRepo.CreateApplication(ctx, &dbgen.CreateApplicationParams{
		Name:              "Test App",
		Description:       pgutil.String("Test Description"),
		DeveloperID:       pgutil.UUID(developerID),
		Category:          "Games",
		Price:             pgutil.Float64(9.99),
		Size:              1024 * 1024 * 100,
		MinAndroidVersion: "8.0",
		CurrentVersion:    "1.0.0",
		Tags:              []string{"game", "action"},
		Screenshots:       []string{"screenshot1.jpg"},
		StorageUrl:        path.Join("apps", developerID.String(), "test.apk"),
	})
	require.NoError(t, err)

	// Create download record
	userID := uuid.New()
	download, err := testRepo.CreateDownload(ctx, &dbgen.CreateDownloadParams{
		UserID:        pgutil.UUID(userID),
		ApplicationID: app.ID,
		IpAddress:     "127.0.0.1",
		Success:       true,
	})
	require.NoError(t, err)
	assert.NotNil(t, download)
	assert.Equal(t, userID, pgutil.FromUUID(download.UserID))
	assert.Equal(t, app.ID, download.ApplicationID)

	// Get downloads by user
	userDownloads, err := testRepo.GetDownloadsByUser(ctx, userID, 10, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, userDownloads)
	assert.Equal(t, userID, pgutil.FromUUID(userDownloads[0].UserID))

	// Get downloads by application
	appDownloads, err := testRepo.GetDownloadsByApplication(ctx, app.ID.Bytes, 10, 0)
	require.NoError(t, err)
	assert.NotEmpty(t, appDownloads)
	assert.Equal(t, app.ID, appDownloads[0].ApplicationID)

	// Get download stats
	stats, err := testRepo.GetDownloadStats(ctx, app.ID.Bytes)
	require.NoError(t, err)
	assert.Equal(t, int64(1), stats.TotalDownloads)
	assert.Equal(t, int64(1), stats.SuccessfulDownloads)
	assert.Equal(t, int64(0), stats.FailedDownloads)
}
