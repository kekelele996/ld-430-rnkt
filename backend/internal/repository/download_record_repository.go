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

// DownloadRecordRepository provides data access for download records.
type DownloadRecordRepository struct {
	coll *mongo.Collection
}

// NewDownloadRecordRepository creates a DownloadRecordRepository.
func NewDownloadRecordRepository(db *mongo.Database) *DownloadRecordRepository {
	return &DownloadRecordRepository{coll: db.Collection("download_records")}
}

// Create inserts a download record.
func (r *DownloadRecordRepository) Create(ctx context.Context, record *model.DownloadRecord) error {
	res, err := r.coll.InsertOne(ctx, record)
	if err != nil {
		return fmt.Errorf("create download record: %w", err)
	}
	record.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// List returns download records with pagination.
func (r *DownloadRecordRepository) List(ctx context.Context, skip, limit int64) ([]model.DownloadRecord, error) {
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "downloaded_at", Value: -1}})
	cursor, err := r.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("list download records: %w", err)
	}
	defer cursor.Close(ctx)
	var records []model.DownloadRecord
	if err := cursor.All(ctx, &records); err != nil {
		return nil, fmt.Errorf("decode download records: %w", err)
	}
	return records, nil
}
