package types

import (
	"context"
	"io"
	"time"

	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/database/generated"
	"github.com/google/uuid"
)

// StorageService represents the interface for object storage operations
type StorageService interface {
	UploadFile(ctx context.Context, objectName string, reader io.Reader, size int64) error
	GetDownloadURL(ctx context.Context, objectName string) (string, time.Time, error)
	DeleteFile(ctx context.Context, objectName string) error
}

// ApplicationService represents the application service interface
type ApplicationService interface {
	CreateApplication(ctx context.Context, params CreateApplicationParams) (*database.CreateApplicationRow, error)
	GetApplication(ctx context.Context, id uuid.UUID) (*database.GetApplicationRow, error)
	UpdateApplication(ctx context.Context, params UpdateApplicationParams) (*database.UpdateApplicationRow, error)
	DeleteApplication(ctx context.Context, id uuid.UUID) error
	SearchApplications(ctx context.Context, params SearchApplicationsParams) ([]*database.SearchApplicationsRow, error)
	GetDownloadURL(ctx context.Context, appID, userID uuid.UUID, ipAddress string) (string, time.Time, error)
	ListCategories(ctx context.Context, params ListCategoriesParams) ([]*database.Category, error)
}

// CreateApplicationParams represents the parameters for creating a new application
type CreateApplicationParams struct {
	Name              string
	Description       string
	DeveloperID       string
	Category          string
	Price             float64
	Size              int64
	MinAndroidVersion string
	CurrentVersion    string
	Tags              []string
	Screenshots       []string
	APKFile           []byte
}

// UpdateApplicationParams represents the parameters for updating an application
type UpdateApplicationParams struct {
	ID                uuid.UUID
	Name              *string
	Description       *string
	Category          *string
	Price             *float64
	MinAndroidVersion *string
	CurrentVersion    *string
	Tags              *[]string
	Screenshots       *[]string
	APKFile           []byte
}

// SearchApplicationsParams represents the parameters for searching applications
type SearchApplicationsParams struct {
	Query             string
	Category          *string
	MinPrice          *float64
	MaxPrice          *float64
	MinAndroidVersion *string
	Tags              []string
	Limit             int32
	Offset            int32
	SortByDownloads   bool
	SortByRating      bool
}

// ListCategoriesParams represents the parameters for listing categories
type ListCategoriesParams struct {
	Limit  int32
	Offset int32
}
