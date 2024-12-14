package handler

import (
	"context"
	"testing"
	"time"

	"connectrpc.com/connect"
	appv1 "github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/gen/app/v1"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/database/generated"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/pgutil"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// MockService je mock implementacija types.ApplicationService interface-a
type MockService struct {
	mock.Mock
}

func (m *MockService) CreateApplication(ctx context.Context, params types.CreateApplicationParams) (*database.CreateApplicationRow, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.CreateApplicationRow), args.Error(1)
}

func (m *MockService) GetApplication(ctx context.Context, id uuid.UUID) (*database.GetApplicationRow, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.GetApplicationRow), args.Error(1)
}

func (m *MockService) UpdateApplication(ctx context.Context, params types.UpdateApplicationParams) (*database.UpdateApplicationRow, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.UpdateApplicationRow), args.Error(1)
}

func (m *MockService) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}

func (m *MockService) SearchApplications(ctx context.Context, params types.SearchApplicationsParams) ([]*database.SearchApplicationsRow, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*database.SearchApplicationsRow), args.Error(1)
}

func (m *MockService) GetDownloadURL(ctx context.Context, appID, userID uuid.UUID, ipAddress string) (string, time.Time, error) {
	args := m.Called(ctx, appID, userID, ipAddress)
	return args.String(0), args.Get(1).(time.Time), args.Error(2)
}

func (m *MockService) ListCategories(ctx context.Context, params types.ListCategoriesParams) ([]*database.Category, error) {
	args := m.Called(ctx, params)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*database.Category), args.Error(1)
}

// TestCreateApplication testira ustvarjanje nove aplikacije
func TestCreateApplication(t *testing.T) {
	mockSvc := new(MockService)
	h := New(mockSvc)
	ctx := context.Background()

	// Pripravimo testne podatke
	req := connect.NewRequest(&appv1.CreateApplicationRequest{
		Name:              "Test App",
		Description:       "Test Description",
		DeveloperId:       uuid.New().String(),
		Category:          "Games",
		Price:             1.99,
		Size:              1024,
		MinAndroidVersion: "8.0",
		CurrentVersion:    "1.0.0",
		Tags:              []string{"game", "arcade"},
		Screenshots:       []string{"screenshot1.jpg"},
		ApkFile:           []byte("fake-apk-content"),
	})

	// Pripravimo pričakovan odgovor iz service layerja
	now := pgutil.NowTime()
	expectedApp := &database.Application{
		ID:                pgutil.UUID(uuid.New()),
		Name:              "Test App",
		Description:       pgutil.String("Test Description"),
		DeveloperID:       pgutil.UUID(uuid.MustParse(req.Msg.DeveloperId)),
		Category:          "Games",
		Price:             pgutil.Float64(1.99),
		Size:              1024,
		MinAndroidVersion: "8.0",
		CurrentVersion:    "1.0.0",
		Tags:              []string{"game", "arcade"},
		Screenshots:       []string{"screenshot1.jpg"},
		StorageUrl:        "apps/test/app.apk",
		Rating:            pgutil.Float64(0),
		Downloads:         0,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	// Nastavimo pričakovanja za mock
	mockSvc.On("CreateApplication", ctx, mock.AnythingOfType("types.CreateApplicationParams")).Return(expectedApp, nil)

	// Izvedemo test
	resp, err := h.CreateApplication(ctx, req)

	// Preverimo rezultate
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Msg)
	assert.NotNil(t, resp.Msg.Application)
	assert.Equal(t, pgutil.StringFromUUID(expectedApp.ID), resp.Msg.Application.Id)
	assert.Equal(t, expectedApp.Name, resp.Msg.Application.Name)
	assert.Equal(t, pgutil.FromString(expectedApp.Description), resp.Msg.Application.Description)
	assert.Equal(t, expectedApp.Category, resp.Msg.Application.Category)
	assert.Equal(t, pgutil.FromFloat64(expectedApp.Price), resp.Msg.Application.Price)
	assert.Equal(t, expectedApp.Size, resp.Msg.Application.Size)
	assert.Equal(t, expectedApp.MinAndroidVersion, resp.Msg.Application.MinAndroidVersion)
	assert.Equal(t, expectedApp.CurrentVersion, resp.Msg.Application.CurrentVersion)
	assert.Equal(t, expectedApp.Tags, resp.Msg.Application.Tags)
	assert.Equal(t, expectedApp.Screenshots, resp.Msg.Application.Screenshots)
	assert.Equal(t, timestamppb.New(pgutil.FromTime(expectedApp.CreatedAt)), resp.Msg.Application.CreatedAt)
	assert.Equal(t, timestamppb.New(pgutil.FromTime(expectedApp.UpdatedAt)), resp.Msg.Application.UpdatedAt)
}

// TestSearchApplications testira iskanje aplikacij
func TestSearchApplications(t *testing.T) {
	mockSvc := new(MockService)
	h := New(mockSvc)
	ctx := context.Background()

	// Pripravimo testne podatke
	req := connect.NewRequest(&appv1.SearchApplicationsRequest{
		Query:             "test",
		Category:          proto.String("Games"),
		MinPrice:          proto.Float64(0.99),
		MaxPrice:          proto.Float64(9.99),
		MinAndroidVersion: proto.String("8.0"),
		Tags:              []string{"game"},
		Pagination: &appv1.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
		Sort: []*appv1.SortOrder{
			{Field: "downloads", Ascending: true},
		},
	})

	// Pripravimo pričakovan odgovor iz service layerja
	app1 := createTestApplication("Test App 1", "Games", 1.99)
	app2 := createTestApplication("Test App 2", "Games", 2.99)
	expectedApps := []*database.Application{app1, app2}

	// Nastavimo pričakovanja za mock
	mockSvc.On("SearchApplications", ctx, mock.AnythingOfType("types.SearchApplicationsParams")).Return(expectedApps, nil)

	// Izvedemo test
	resp, err := h.SearchApplications(ctx, req)

	// Preverimo rezultate
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Msg)
	assert.Len(t, resp.Msg.Applications, 2)
	assert.Equal(t, pgutil.StringFromUUID(app1.ID), resp.Msg.Applications[0].Id)
	assert.Equal(t, pgutil.StringFromUUID(app2.ID), resp.Msg.Applications[1].Id)
	assert.NotNil(t, resp.Msg.Pagination)
	assert.Equal(t, int32(2), resp.Msg.Pagination.TotalItems)
	assert.Equal(t, int32(1), resp.Msg.Pagination.TotalPages)
	assert.Equal(t, int32(1), resp.Msg.Pagination.CurrentPage)
	assert.Equal(t, int32(10), resp.Msg.Pagination.PageSize)
}

// TestListCategories testira pridobivanje seznama kategorij
func TestListCategories(t *testing.T) {
	mockSvc := new(MockService)
	h := New(mockSvc)
	ctx := context.Background()

	// Pripravimo testne podatke
	req := connect.NewRequest(&appv1.ListCategoriesRequest{
		Pagination: &appv1.PaginationRequest{
			Page:     1,
			PageSize: 10,
		},
	})

	// Pripravimo pričakovan odgovor iz service layerja
	now := pgutil.NowTime()
	cat1 := &database.Category{
		ID:          pgutil.UUID(uuid.New()),
		Name:        "Games",
		Description: pgutil.String("Mobile games"),
		CreatedAt:   now,
	}
	cat2 := &database.Category{
		ID:          pgutil.UUID(uuid.New()),
		Name:        "Productivity",
		Description: pgutil.String("Productivity apps"),
		CreatedAt:   now,
	}
	expectedCategories := []*database.Category{cat1, cat2}

	// Nastavimo pričakovanja za mock
	mockSvc.On("ListCategories", ctx, mock.AnythingOfType("types.ListCategoriesParams")).Return(expectedCategories, nil)

	// Izvedemo test
	resp, err := h.ListCategories(ctx, req)

	// Preverimo rezultate
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Msg)
	assert.Len(t, resp.Msg.Categories, 2)
	assert.Equal(t, pgutil.StringFromUUID(cat1.ID), resp.Msg.Categories[0].Id)
	assert.Equal(t, cat1.Name, resp.Msg.Categories[0].Name)
	assert.Equal(t, pgutil.FromString(cat1.Description), resp.Msg.Categories[0].Description)
	assert.Equal(t, timestamppb.New(pgutil.FromTime(cat1.CreatedAt)), resp.Msg.Categories[0].CreatedAt)
}

// TestGetApplication testira pridobivanje aplikacije po ID-ju
func TestGetApplication(t *testing.T) {
	mockSvc := new(MockService)
	h := New(mockSvc)
	ctx := context.Background()

	// Pripravimo testne podatke
	appID := uuid.New()
	req := connect.NewRequest(&appv1.GetApplicationRequest{
		Id: appID.String(),
	})

	// Pripravimo pričakovan odgovor iz service layerja
	expectedApp := createTestApplication("Test App", "Games", 1.99)
	expectedApp.ID = pgutil.UUID(appID)

	// Nastavimo pričakovanja za mock
	mockSvc.On("GetApplication", ctx, appID).Return(expectedApp, nil)

	// Izvedemo test
	resp, err := h.GetApplication(ctx, req)

	// Preverimo rezultate
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Msg)
	assert.NotNil(t, resp.Msg.Application)
	assert.Equal(t, appID.String(), resp.Msg.Application.Id)
	assert.Equal(t, expectedApp.Name, resp.Msg.Application.Name)
	assert.Equal(t, pgutil.FromString(expectedApp.Description), resp.Msg.Application.Description)
}

// TestUpdateApplication testira posodabljanje aplikacije
func TestUpdateApplication(t *testing.T) {
	mockSvc := new(MockService)
	h := New(mockSvc)
	ctx := context.Background()

	// Pripravimo testne podatke
	appID := uuid.New()
	newName := "Updated App"
	newPrice := 2.99
	req := connect.NewRequest(&appv1.UpdateApplicationRequest{
		Id:    appID.String(),
		Name:  &newName,
		Price: &newPrice,
	})

	// Pripravimo pričakovan odgovor iz service layerja
	expectedApp := createTestApplication(newName, "Games", newPrice)
	expectedApp.ID = pgutil.UUID(appID)

	// Nastavimo pričakovanja za mock
	mockSvc.On("UpdateApplication", ctx, mock.AnythingOfType("types.UpdateApplicationParams")).Return(expectedApp, nil)

	// Izvedemo test
	resp, err := h.UpdateApplication(ctx, req)

	// Preverimo rezultate
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Msg)
	assert.NotNil(t, resp.Msg.Application)
	assert.Equal(t, appID.String(), resp.Msg.Application.Id)
	assert.Equal(t, newName, resp.Msg.Application.Name)
	assert.Equal(t, newPrice, resp.Msg.Application.Price)
}

// TestDeleteApplication testira brisanje aplikacije
func TestDeleteApplication(t *testing.T) {
	mockSvc := new(MockService)
	h := New(mockSvc)
	ctx := context.Background()

	// Pripravimo testne podatke
	appID := uuid.New()
	req := connect.NewRequest(&appv1.DeleteApplicationRequest{
		Id: appID.String(),
	})

	// Pripravimo aplikacijo, ki bo vrnjena pred brisanjem
	expectedApp := createTestApplication("Test App", "Games", 1.99)
	expectedApp.ID = pgutil.UUID(appID)

	// Nastavimo pričakovanja za mock
	mockSvc.On("GetApplication", ctx, appID).Return(expectedApp, nil)
	mockSvc.On("DeleteApplication", ctx, appID).Return(nil)

	// Izvedemo test
	resp, err := h.DeleteApplication(ctx, req)

	// Preverimo rezultate
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Msg)
	assert.NotNil(t, resp.Msg.Application)
	assert.Equal(t, appID.String(), resp.Msg.Application.Id)
}

// TestDownloadApplication testira generiranje URL-ja za prenos
func TestDownloadApplication(t *testing.T) {
	mockSvc := new(MockService)
	h := New(mockSvc)
	ctx := context.Background()

	// Pripravimo testne podatke
	appID := uuid.New()
	userID := uuid.New()
	req := connect.NewRequest(&appv1.DownloadApplicationRequest{
		ApplicationId: appID.String(),
		UserId:        userID.String(),
	})

	// Pripravimo pričakovan odgovor
	expectedURL := "https://storage.example.com/download/app.apk"
	expiresAt := time.Now().Add(time.Hour)

	// Nastavimo pričakovanja za mock
	mockSvc.On("GetDownloadURL", ctx, appID, userID, mock.AnythingOfType("string")).Return(expectedURL, expiresAt, nil)

	// Izvedemo test
	resp, err := h.DownloadApplication(ctx, req)

	// Preverimo rezultate
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.NotNil(t, resp.Msg)
	assert.Equal(t, expectedURL, resp.Msg.DownloadUrl)
	assert.Equal(t, timestamppb.New(expiresAt), resp.Msg.ExpiresAt)
}

// Popravimo helper funkcijo za ustvarjanje testne aplikacije
func createTestApplication(name, category string, price float64) *database.Application {
	now := time.Now()
	return &database.Application{
		ID:                pgutil.UUID(uuid.New()),
		Name:              name,
		Description:       pgutil.String("Test Description"),
		DeveloperID:       pgutil.UUID(uuid.New()),
		Category:          category,
		Price:             pgutil.Float64(price),
		Size:              1024,
		MinAndroidVersion: "8.0",
		CurrentVersion:    "1.0.0",
		Tags:              []string{"test"},
		Screenshots:       []string{"test.jpg"},
		StorageUrl:        "apps/test/app.apk",
		Rating:            pgutil.Float64(4.5),
		Downloads:         100,
		CreatedAt:         pgutil.Time(now),
		UpdatedAt:         pgutil.Time(now),
	}
}
