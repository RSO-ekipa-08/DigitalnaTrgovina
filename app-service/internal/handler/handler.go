package handler

import (
	"context"
	"fmt"

	"connectrpc.com/connect"
	appv1 "github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/gen/app/v1"
	database "github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/database/generated"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/pgutil"
	"github.com/RSO-ekipa-08/DigitalnaTrgovina/app-service/internal/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type application = database.CreateApplicationRow

// Handler represents the Connect-RPC handler
type Handler struct {
	svc types.ApplicationService
}

// New creates a new handler instance
func New(svc types.ApplicationService) *Handler {
	return &Handler{
		svc: svc,
	}
}

// CreateApplication creates a new application
func (h *Handler) CreateApplication(ctx context.Context, req *connect.Request[appv1.CreateApplicationRequest]) (*connect.Response[appv1.CreateApplicationResponse], error) {
	developerID := req.Msg.DeveloperId
	if developerID == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("developer_id is required"))
	}

	params := types.CreateApplicationParams{
		Name:              req.Msg.Name,
		Description:       req.Msg.Description,
		DeveloperID:       developerID,
		Category:          req.Msg.Category,
		Price:             req.Msg.Price,
		Size:              req.Msg.Size,
		MinAndroidVersion: req.Msg.MinAndroidVersion,
		CurrentVersion:    req.Msg.CurrentVersion,
		Tags:              req.Msg.Tags,
		Screenshots:       req.Msg.Screenshots,
		APKFile:           req.Msg.ApkFile,
	}

	app, err := h.svc.CreateApplication(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to create application")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to create application"))
	}

	return connect.NewResponse(&appv1.CreateApplicationResponse{
		Application: convertApplicationToProto(app),
	}), nil
}

// GetApplication gets an application by ID
func (h *Handler) GetApplication(ctx context.Context, req *connect.Request[appv1.GetApplicationRequest]) (*connect.Response[appv1.GetApplicationResponse], error) {
	id, err := uuid.Parse(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid application ID"))
	}

	app, err := h.svc.GetApplication(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get application")
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("application not found"))
	}

	return connect.NewResponse(&appv1.GetApplicationResponse{
		Application: convertApplicationToProto((*application)(app)),
	}), nil
}

// UpdateApplication updates an existing application
func (h *Handler) UpdateApplication(ctx context.Context, req *connect.Request[appv1.UpdateApplicationRequest]) (*connect.Response[appv1.UpdateApplicationResponse], error) {
	id, err := uuid.Parse(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid application ID"))
	}

	params := types.UpdateApplicationParams{
		ID: id,
	}

	if req.Msg.Name != nil {
		params.Name = req.Msg.Name
	}
	if req.Msg.Description != nil {
		params.Description = req.Msg.Description
	}
	if req.Msg.Category != nil {
		params.Category = req.Msg.Category
	}
	if req.Msg.Price != nil {
		params.Price = req.Msg.Price
	}
	if req.Msg.MinAndroidVersion != nil {
		params.MinAndroidVersion = req.Msg.MinAndroidVersion
	}
	if req.Msg.CurrentVersion != nil {
		params.CurrentVersion = req.Msg.CurrentVersion
	}
	if req.Msg.Tags != nil {
		params.Tags = &req.Msg.Tags
	}
	if req.Msg.Screenshots != nil {
		params.Screenshots = &req.Msg.Screenshots
	}
	if req.Msg.ApkFile != nil {
		params.APKFile = req.Msg.ApkFile
	}

	app, err := h.svc.UpdateApplication(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to update application")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to update application"))
	}

	return connect.NewResponse(&appv1.UpdateApplicationResponse{
		Application: convertApplicationToProto((*application)(app)),
	}), nil
}

// DeleteApplication deletes an application
func (h *Handler) DeleteApplication(ctx context.Context, req *connect.Request[appv1.DeleteApplicationRequest]) (*connect.Response[appv1.DeleteApplicationResponse], error) {
	id, err := uuid.Parse(req.Msg.Id)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid application ID"))
	}

	// Get application before deletion
	app, err := h.svc.GetApplication(ctx, id)
	if err != nil {
		log.Error().Err(err).Msg("failed to get application")
		return nil, connect.NewError(connect.CodeNotFound, fmt.Errorf("application not found"))
	}

	if err := h.svc.DeleteApplication(ctx, id); err != nil {
		log.Error().Err(err).Msg("failed to delete application")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to delete application"))
	}

	return connect.NewResponse(&appv1.DeleteApplicationResponse{
		Application: convertApplicationToProto((*application)(app)),
	}), nil
}

// SearchApplications searches for applications
func (h *Handler) SearchApplications(ctx context.Context, req *connect.Request[appv1.SearchApplicationsRequest]) (*connect.Response[appv1.SearchApplicationsResponse], error) {
	params := types.SearchApplicationsParams{
		Query:  req.Msg.Query,
		Limit:  req.Msg.Pagination.PageSize,
		Offset: (req.Msg.Pagination.Page - 1) * req.Msg.Pagination.PageSize,
	}

	if req.Msg.Category != nil {
		params.Category = req.Msg.Category
	}
	if req.Msg.MinPrice != nil {
		params.MinPrice = req.Msg.MinPrice
	}
	if req.Msg.MaxPrice != nil {
		params.MaxPrice = req.Msg.MaxPrice
	}
	if req.Msg.MinAndroidVersion != nil {
		params.MinAndroidVersion = req.Msg.MinAndroidVersion
	}
	if len(req.Msg.Tags) > 0 {
		params.Tags = req.Msg.Tags
	}

	// Handle sorting
	for _, sort := range req.Msg.Sort {
		switch sort.Field {
		case "downloads":
			params.SortByDownloads = sort.Ascending
		case "rating":
			params.SortByRating = sort.Ascending
		}
	}

	apps, err := h.svc.SearchApplications(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to search applications")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to search applications"))
	}

	// Convert applications to proto
	protoApps := make([]*appv1.Application, len(apps))
	for i, app := range apps {
		protoApps[i] = convertApplicationToProto((*application)(app))
	}

	// Calculate pagination
	totalItems := int32(len(apps))
	totalPages := (totalItems + params.Limit - 1) / params.Limit
	currentPage := (params.Offset / params.Limit) + 1

	resp := &appv1.SearchApplicationsResponse{
		Applications: protoApps,
		Pagination: &appv1.PaginationResponse{
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: currentPage,
			PageSize:    params.Limit,
		},
	}

	return connect.NewResponse(resp), nil
}

// DownloadApplication generates a download URL for an application
func (h *Handler) DownloadApplication(ctx context.Context, req *connect.Request[appv1.DownloadApplicationRequest]) (*connect.Response[appv1.DownloadApplicationResponse], error) {
	appID, err := uuid.Parse(req.Msg.ApplicationId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid application ID"))
	}

	userID, err := uuid.Parse(req.Msg.UserId)
	if err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("invalid user ID"))
	}

	// Get client IP from context
	ipAddress := "unknown" // TODO: Get from context

	url, expiresAt, err := h.svc.GetDownloadURL(ctx, appID, userID, ipAddress)
	if err != nil {
		log.Error().Err(err).Msg("failed to get download URL")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to get download URL"))
	}

	resp := &appv1.DownloadApplicationResponse{
		DownloadUrl: url,
		ExpiresAt:   timestamppb.New(expiresAt),
	}

	return connect.NewResponse(resp), nil
}

// ListCategories lists all available categories
func (h *Handler) ListCategories(ctx context.Context, req *connect.Request[appv1.ListCategoriesRequest]) (*connect.Response[appv1.ListCategoriesResponse], error) {
	params := types.ListCategoriesParams{
		Limit:  req.Msg.Pagination.PageSize,
		Offset: (req.Msg.Pagination.Page - 1) * req.Msg.Pagination.PageSize,
	}

	categories, err := h.svc.ListCategories(ctx, params)
	if err != nil {
		log.Error().Err(err).Msg("failed to list categories")
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed to list categories"))
	}

	// Convert categories to proto
	protoCategories := make([]*appv1.Category, len(categories))
	for i, cat := range categories {
		protoCategories[i] = &appv1.Category{
			Id:          pgutil.StringFromUUID(cat.ID),
			Name:        cat.Name,
			Description: pgutil.FromString(cat.Description),
			CreatedAt:   timestamppb.New(cat.CreatedAt.Time),
		}
	}

	// Calculate pagination
	totalItems := int32(len(categories))
	totalPages := (totalItems + params.Limit - 1) / params.Limit
	currentPage := (params.Offset / params.Limit) + 1

	resp := &appv1.ListCategoriesResponse{
		Categories: protoCategories,
		Pagination: &appv1.PaginationResponse{
			TotalItems:  totalItems,
			TotalPages:  totalPages,
			CurrentPage: currentPage,
			PageSize:    params.Limit,
		},
	}

	return connect.NewResponse(resp), nil
}

// Helper function to convert database Application to proto Application
func convertApplicationToProto(app *application) *appv1.Application {
	return &appv1.Application{
		Id:                pgutil.StringFromUUID(app.ID),
		Name:              app.Name,
		Description:       pgutil.FromString(app.Description),
		DeveloperId:       pgutil.StringFromUUID(app.DeveloperID),
		Category:          app.Category,
		Price:             pgutil.FromFloat64(app.Price),
		Size:              app.Size,
		MinAndroidVersion: app.MinAndroidVersion,
		CurrentVersion:    app.CurrentVersion,
		Tags:              app.Tags,
		Screenshots:       app.Screenshots,
		Rating:            pgutil.FromFloat64(app.Rating),
		Downloads:         app.Downloads,
		CreatedAt:         timestamppb.New(app.CreatedAt.Time),
		UpdatedAt:         timestamppb.New(app.UpdatedAt.Time),
		StorageUrl:        app.StorageUrl,
	}
}
