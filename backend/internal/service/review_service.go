package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/errors"
	"github.com/assethub/assethub/internal/model"
	"github.com/assethub/assethub/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReviewService handles asset moderation reviews.
type ReviewService struct {
	reviewRepo *repository.ReviewRecordRepository
	assetRepo  *repository.AssetRepository
	logger     *slog.Logger
}

// NewReviewService creates a ReviewService.
func NewReviewService(reviewRepo *repository.ReviewRecordRepository, assetRepo *repository.AssetRepository, logger *slog.Logger) *ReviewService {
	return &ReviewService{reviewRepo: reviewRepo, assetRepo: assetRepo, logger: logger}
}

// Review submits a review result for an asset.
func (s *ReviewService) Review(ctx context.Context, assetID string, reviewerID primitive.ObjectID, result, comment string) (*model.ReviewRecord, error) {
	oid, err := primitive.ObjectIDFromHex(assetID)
	if err != nil {
		return nil, errors.NewBusinessError(constants.CodeBadRequest, "invalid asset id")
	}
	if _, err := s.assetRepo.FindByID(ctx, oid); err != nil {
		return nil, err
	}
	record := &model.ReviewRecord{
		AssetID:    oid,
		ReviewerID: reviewerID,
		Result:     result,
		Comment:    comment,
		ReviewedAt: time.Now(),
	}
	if err := s.reviewRepo.Create(ctx, record); err != nil {
		return nil, err
	}
	// Approved reviews publish the asset; Rejected marks it Flagged.
	status := string(constants.AssetStatusPublished)
	if result == string(constants.ReviewRejected) {
		status = string(constants.AssetStatusFlagged)
	}
	if err := s.assetRepo.Update(ctx, oid, bson.M{"status": status}); err != nil {
		s.logger.Warn("update asset status on review failed", "error", err.Error())
	}
	s.logger.Info(constants.LogReviewAsset, "asset_id", oid.Hex(), "result", result)
	return record, nil
}

// ListByAsset returns reviews for an asset.
func (s *ReviewService) ListByAsset(ctx context.Context, assetID string) ([]model.ReviewRecord, error) {
	oid, err := primitive.ObjectIDFromHex(assetID)
	if err != nil {
		return nil, errors.NewBusinessError(constants.CodeBadRequest, "invalid asset id")
	}
	return s.reviewRepo.ListByAsset(ctx, oid)
}
