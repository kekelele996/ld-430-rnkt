package repository

import (
	"context"
	"fmt"

	"github.com/assethub/assethub/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ReviewRecordRepository provides data access for review records.
type ReviewRecordRepository struct {
	coll *mongo.Collection
}

// NewReviewRecordRepository creates a ReviewRecordRepository.
func NewReviewRecordRepository(db *mongo.Database) *ReviewRecordRepository {
	return &ReviewRecordRepository{coll: db.Collection("review_records")}
}

// Create inserts a review record.
func (r *ReviewRecordRepository) Create(ctx context.Context, record *model.ReviewRecord) error {
	res, err := r.coll.InsertOne(ctx, record)
	if err != nil {
		return fmt.Errorf("create review record: %w", err)
	}
	record.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// ListByAsset returns reviews for an asset.
func (r *ReviewRecordRepository) ListByAsset(ctx context.Context, assetID primitive.ObjectID) ([]model.ReviewRecord, error) {
	opts := options.Find().SetSort(bson.D{{Key: "reviewed_at", Value: -1}})
	cursor, err := r.coll.Find(ctx, bson.M{"asset_id": assetID}, opts)
	if err != nil {
		return nil, fmt.Errorf("list review records: %w", err)
	}
	defer cursor.Close(ctx)
	var records []model.ReviewRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("decode review records: %w", err)
	}
	return records, nil
}
