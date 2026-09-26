package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/errors"
	"github.com/assethub/assethub/internal/model"
	"github.com/assethub/assethub/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DownloadService handles download records and permission checks.
type DownloadService struct {
	recordRepo *repository.DownloadRecordRepository
	assetRepo  *repository.AssetRepository
	logger     *slog.Logger
}

// NewDownloadService creates a DownloadService.
func NewDownloadService(recordRepo *repository.DownloadRecordRepository, assetRepo *repository.AssetRepository, logger *slog.Logger) *DownloadService {
	return &DownloadService{recordRepo: recordRepo, assetRepo: assetRepo, logger: logger}
}

// Download records a download after enforcing license permission rules.
func (s *DownloadService) Download(ctx context.Context, assetID string, downloaderID primitive.ObjectID, purpose, licenseVersion, ip string, role string) (*model.Asset, error) {
	oid, err := primitive.ObjectIDFromHex(assetID)
	if err != nil {
		return nil, errors.NewBusinessError(constants.CodeBadRequest, "invalid asset id")
	}
	asset, err := s.assetRepo.FindByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if asset.Status != string(constants.AssetStatusPublished) {
		return nil, errors.NewBusinessError(constants.CodeForbidden, "asset is not published")
	}
	// Viewer can only download free assets; Commercial/Extended require extra permission.
	if asset.LicenseType == string(constants.LicenseCommercial) || asset.LicenseType == string(constants.LicenseExtended) {
		if role == string(constants.RoleViewer) {
			return nil, errors.NewBusinessError(constants.CodeForbidden, constants.MsgCommercialForbidden)
		}
	}

	record := &model.DownloadRecord{
		AssetID:        oid,
		DownloaderID:   downloaderID,
		DownloadedAt:   time.Now(),
		Purpose:        purpose,
		LicenseVersion: licenseVersion,
		IP:             ip,
	}
	if err := s.recordRepo.Create(ctx, record); err != nil {
		return nil, err
	}
	if err := s.assetRepo.IncDownloadCount(ctx, oid); err != nil {
		s.logger.Warn("increment download count failed", "error", err.Error())
	}
	s.logger.Info(constants.LogDownloadAsset, "asset_id", oid.Hex(), "downloader", downloaderID.Hex())
	return asset, nil
}

// List returns download records.
func (s *DownloadService) List(ctx context.Context, skip, limit int64) ([]model.DownloadRecord, error) {
	return s.recordRepo.List(ctx, skip, limit)
}
