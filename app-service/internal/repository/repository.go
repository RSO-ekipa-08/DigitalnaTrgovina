package repository

import (
	"context"

	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/database/generated"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/pgutil"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository represents the application's data access layer
type Repository struct {
	db *pgxpool.Pool
	q  *database.Queries
}

// New creates a new repository instance
func New(db *pgxpool.Pool) *Repository {
	return &Repository{
		db: db,
		q:  database.New(db),
	}
}

// Application represents the application repository interface
type ApplicationRepository interface {
	GetApplication(ctx context.Context, id uuid.UUID) (*database.GetApplicationRow, error)
	ListApplications(ctx context.Context, limit, offset int32) ([]*database.ListApplicationsRow, error)
	SearchApplications(ctx context.Context, params *database.SearchApplicationsParams) ([]*database.SearchApplicationsRow, error)
	CreateApplication(ctx context.Context, params *database.CreateApplicationParams) (*database.CreateApplicationRow, error)
	UpdateApplication(ctx context.Context, params *database.UpdateApplicationParams) (*database.UpdateApplicationRow, error)
	DeleteApplication(ctx context.Context, id uuid.UUID) error
	IncrementDownloads(ctx context.Context, id uuid.UUID) (*database.IncrementDownloadsRow, error)
	ListCategories(ctx context.Context, params *database.ListCategoriesParams) ([]*database.Category, error)
	CreateDownload(ctx context.Context, params *database.CreateDownloadParams) (*database.Download, error)
	GetDownloadsByUser(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*database.Download, error)
}

// Download represents the download repository interface
type DownloadRepository interface {
	CreateDownload(ctx context.Context, params *database.CreateDownloadParams) (*database.Download, error)
	GetDownloadsByUser(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*database.Download, error)
	GetDownloadsByApplication(ctx context.Context, appID uuid.UUID, limit, offset int32) ([]*database.Download, error)
	GetDownloadStats(ctx context.Context, appID uuid.UUID) (*database.GetDownloadStatsRow, error)
}

// Ensure Repository implements both interfaces
var _ ApplicationRepository = (*Repository)(nil)
var _ DownloadRepository = (*Repository)(nil)

// Application repository methods
func (r *Repository) GetApplication(ctx context.Context, id uuid.UUID) (*database.GetApplicationRow, error) {
	return r.q.GetApplication(ctx, pgutil.UUID(id))
}

func (r *Repository) ListApplications(ctx context.Context, limit, offset int32) ([]*database.ListApplicationsRow, error) {
	params := &database.ListApplicationsParams{
		Limit:  limit,
		Offset: offset,
	}
	return r.q.ListApplications(ctx, params)
}

func (r *Repository) SearchApplications(ctx context.Context, params *database.SearchApplicationsParams) ([]*database.SearchApplicationsRow, error) {
	return r.q.SearchApplications(ctx, params)
}

func (r *Repository) CreateApplication(ctx context.Context, params *database.CreateApplicationParams) (*database.CreateApplicationRow, error) {
	return r.q.CreateApplication(ctx, params)
}

func (r *Repository) UpdateApplication(ctx context.Context, params *database.UpdateApplicationParams) (*database.UpdateApplicationRow, error) {
	return r.q.UpdateApplication(ctx, params)
}

func (r *Repository) DeleteApplication(ctx context.Context, id uuid.UUID) error {
	return r.q.DeleteApplication(ctx, pgutil.UUID(id))
}

func (r *Repository) IncrementDownloads(ctx context.Context, id uuid.UUID) (*database.IncrementDownloadsRow, error) {
	return r.q.IncrementDownloads(ctx, pgutil.UUID(id))
}

func (r *Repository) ListCategories(ctx context.Context, params *database.ListCategoriesParams) ([]*database.Category, error) {
	return r.q.ListCategories(ctx, params)
}

// Download repository methods
func (r *Repository) CreateDownload(ctx context.Context, params *database.CreateDownloadParams) (*database.Download, error) {
	return r.q.CreateDownload(ctx, params)
}

func (r *Repository) GetDownloadsByUser(ctx context.Context, userID uuid.UUID, limit, offset int32) ([]*database.Download, error) {
	params := &database.GetDownloadsByUserParams{
		UserID: pgutil.UUID(userID),
		Limit:  limit,
		Offset: offset,
	}
	return r.q.GetDownloadsByUser(ctx, params)
}

func (r *Repository) GetDownloadsByApplication(ctx context.Context, appID uuid.UUID, limit, offset int32) ([]*database.Download, error) {
	params := &database.GetDownloadsByApplicationParams{
		ApplicationID: pgutil.UUID(appID),
		Limit:         limit,
		Offset:        offset,
	}
	return r.q.GetDownloadsByApplication(ctx, params)
}

func (r *Repository) GetDownloadStats(ctx context.Context, appID uuid.UUID) (*database.GetDownloadStatsRow, error) {
	return r.q.GetDownloadStats(ctx, pgutil.UUID(appID))
}
