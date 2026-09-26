package service

import (
	"context"
	"fmt"
	"log/slog"
	"mime/multipart"
	"strings"

	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/errors"
	"github.com/assethub/assethub/internal/model"
	"github.com/assethub/assethub/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AssetService handles asset business logic.
type AssetService struct {
	repo         *repository.AssetRepository
	tagRepo      *repository.TagRepository
	categoryRepo *repository.CategoryRepository
	storage      *StorageService
	search       *SearchService
	logger       *slog.Logger
}

// NewAssetService creates an AssetService.
func NewAssetService(repo *repository.AssetRepository, tagRepo *repository.TagRepository, categoryRepo *repository.CategoryRepository, storage *StorageService, search *SearchService, logger *slog.Logger) *AssetService {
	return &AssetService{repo: repo, tagRepo: tagRepo, categoryRepo: categoryRepo, storage: storage, search: search, logger: logger}
}

// Upload stores the file in MinIO and creates the asset metadata.
func (s *AssetService) Upload(ctx context.Context, uploaderID primitive.ObjectID, req UploadMeta, fileHeader *multipart.FileHeader) (*model.Asset, error) {
	if !constants.ValidAssetTypes[constants.AssetType(req.FileType)] {
		return nil, errors.NewBusinessError(constants.CodeValidationFailed, "invalid file_type")
	}
	if !constants.ValidLicenseTypes[constants.LicenseType(req.LicenseType)] {
		return nil, errors.NewBusinessError(constants.CodeValidationFailed, "invalid license_type")
	}
	categoryID, err := primitive.ObjectIDFromHex(req.CategoryID)
	if err != nil {
		return nil, errors.NewBusinessError(constants.CodeValidationFailed, "invalid category_id")
	}
	if _, err := s.categoryRepo.FindByID(ctx, categoryID); err != nil {
		return nil, fmt.Errorf("find category: %w", err)
	}

	objectKey := s.storage.GenObjectKey(uploaderID.Hex(), fileHeader.Filename)
	_, fileURL, thumbnailURL, err := s.storage.UploadObject(ctx, fileHeader, objectKey)
	if err != nil {
		return nil, err
	}

	asset := &model.Asset{
		Title:        req.Title,
		Description:  req.Description,
		FileType:     req.FileType,
		FileFormat:   strings.ToUpper(req.FileFormat),
		FileURL:      fileURL,
		ThumbnailURL: thumbnailURL,
		FileSize:     fileHeader.Size,
		Width:        req.Width,
		Height:       req.Height,
		Tags:         req.Tags,
		CategoryID:   categoryID,
		UploaderID:   uploaderID,
		LicenseType:  req.LicenseType,
		Status:       string(constants.AssetStatusDraft),
		ObjectKey:    objectKey,
	}
	if err := s.repo.Create(ctx, asset); err != nil {
		return nil, err
	}
	for _, tag := range req.Tags {
		if err := s.tagRepo.UpsertByName(ctx, tag, string(constants.TagCategoryOther)); err != nil {
			s.logger.Warn("upsert tag failed", "tag", tag, "error", err.Error())
		}
	}
	s.logger.Info(constants.LogUploadAsset, "asset_id", asset.ID.Hex(), "uploader", uploaderID.Hex())
	return asset, nil
}

// Update edits asset metadata.
func (s *AssetService) Update(ctx context.Context, id string, upd UpdateMeta, operatorID primitive.ObjectID, operatorRole string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewBusinessError(constants.CodeBadRequest, "invalid asset id")
	}
	asset, err := s.repo.FindByID(ctx, oid)
	if err != nil {
		return err
	}
	if asset.UploaderID != operatorID && !canModerate(operatorRole) {
		return errors.NewBusinessError(constants.CodeForbidden, constants.MsgForbidden)
	}
	update := bson.M{}
	if upd.Title != "" {
		update["title"] = upd.Title
	}
	if upd.Description != "" {
		update["description"] = upd.Description
	}
	if upd.Tags != nil {
		update["tags"] = upd.Tags
	}
	if upd.CategoryID != "" {
		categoryID, err := primitive.ObjectIDFromHex(upd.CategoryID)
		if err != nil {
			return errors.NewBusinessError(constants.CodeValidationFailed, "invalid category_id")
		}
		update["category_id"] = categoryID
	}
	if upd.LicenseType != "" {
		if !constants.ValidLicenseTypes[constants.LicenseType(upd.LicenseType)] {
			return errors.NewBusinessError(constants.CodeValidationFailed, "invalid license_type")
		}
		update["license_type"] = upd.LicenseType
	}
	if upd.Status != "" {
		if !constants.ValidAssetStatuses[constants.AssetStatus(upd.Status)] {
			return errors.NewBusinessError(constants.CodeValidationFailed, "invalid status")
		}
		update["status"] = upd.Status
	}
	if len(update) == 0 {
		return nil
	}
	return s.repo.Update(ctx, oid, update)
}

// Publish transitions an asset to Published.
func (s *AssetService) Publish(ctx context.Context, id string, operatorID primitive.ObjectID, operatorRole string) error {
	return s.updateStatus(ctx, id, string(constants.AssetStatusPublished), operatorID, operatorRole)
}

// Archive transitions an asset to Archived.
func (s *AssetService) Archive(ctx context.Context, id string, operatorID primitive.ObjectID, operatorRole string) error {
	return s.updateStatus(ctx, id, string(constants.AssetStatusArchived), operatorID, operatorRole)
}

func (s *AssetService) updateStatus(ctx context.Context, id, status string, operatorID primitive.ObjectID, operatorRole string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewBusinessError(constants.CodeBadRequest, "invalid asset id")
	}
	asset, err := s.repo.FindByID(ctx, oid)
	if err != nil {
		return err
	}
	if asset.UploaderID != operatorID && !canModerate(operatorRole) {
		return errors.NewBusinessError(constants.CodeForbidden, constants.MsgForbidden)
	}
	return s.repo.Update(ctx, oid, bson.M{"status": status})
}

// Get returns an asset by id and increments view count.
func (s *AssetService) Get(ctx context.Context, id string) (*model.Asset, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.NewBusinessError(constants.CodeBadRequest, "invalid asset id")
	}
	asset, err := s.repo.FindByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if err := s.repo.IncViewCount(ctx, oid); err != nil {
		s.logger.Warn("increment view count failed", "error", err.Error())
	}
	asset.ViewCount++
	return asset, nil
}

// List searches assets with filters and returns paginated results.
func (s *AssetService) List(ctx context.Context, query AssetQueryParams) (*AssetListResult, error) {
	assets, total, err := s.search.Search(ctx, query)
	if err != nil {
		return nil, err
	}
	return &AssetListResult{Items: assets, Total: total}, nil
}

// Hot returns the hottest assets by downloads/views.
func (s *AssetService) Hot(ctx context.Context, limit int64) ([]model.Asset, error) {
	if limit <= 0 || limit > 50 {
		limit = 10
	}
	filter := bson.M{"status": string(constants.AssetStatusPublished)}
	assets, err := s.repo.List(ctx, filter, bson.D{{Key: "download_count", Value: -1}, {Key: "view_count", Value: -1}}, 0, limit)
	if err != nil {
		return nil, fmt.Errorf("hot assets: %w", err)
	}
	return assets, nil
}

// UploadMeta carries validated upload metadata.
type UploadMeta struct {
	Title       string
	Description string
	FileType    string
	FileFormat  string
	Tags        []string
	CategoryID  string
	LicenseType string
	Width       int
	Height      int
}

// UpdateMeta carries validated update metadata.
type UpdateMeta struct {
	Title       string
	Description string
	Tags        []string
	CategoryID  string
	LicenseType string
	Status      string
}

// AssetListResult wraps a paginated asset list.
type AssetListResult struct {
	Items []model.Asset `json:"items"`
	Total int64         `json:"total"`
}

// AssetQueryParams are parsed search parameters.
type AssetQueryParams struct {
	Keyword     string
	FileType    string
	LicenseType string
	Status      string
	Tag         string
	CategoryID  string
	Sort        string
	Page        int64
	PageSize    int64
}

// canModerate reports whether the role may moderate (admin/moderator).
func canModerate(role string) bool {
	return role == string(constants.RoleAdmin) || role == string(constants.RoleModerator)
}
