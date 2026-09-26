package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/assethub/assethub/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AssetRepository provides data access for assets.
type AssetRepository struct {
	coll *mongo.Collection
}

// NewAssetRepository creates an AssetRepository.
func NewAssetRepository(db *mongo.Database) *AssetRepository {
	return &AssetRepository{coll: db.Collection("assets")}
}

// Create inserts an asset.
func (r *AssetRepository) Create(ctx context.Context, asset *model.Asset) error {
	now := time.Now()
	asset.CreatedAt = now
	asset.UpdatedAt = now
	res, err := r.coll.InsertOne(ctx, asset)
	if err != nil {
		return fmt.Errorf("create asset: %w", err)
	}
	asset.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// Update updates mutable fields of an asset by id.
func (r *AssetRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("update asset: %w", err)
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// FindByID returns an asset by id.
func (r *AssetRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Asset, error) {
	var asset model.Asset
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&asset)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find asset by id: %w", err)
	}
	return &asset, nil
}

// List queries assets with filter, sort and pagination.
func (r *AssetRepository) List(ctx context.Context, filter bson.M, sort bson.D, skip, limit int64) ([]model.Asset, error) {
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(sort)
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("list assets: %w", err)
	}
	defer cursor.Close(ctx)
	var assets []model.Asset
	if err := cursor.All(ctx, &assets); err != nil {
		return nil, fmt.Errorf("decode assets: %w", err)
	}
	return assets, nil
}

// ListByIDsAndStatus returns assets matching the given ids and status.
// Used by public sharing so non-matching assets never leave the database layer.
func (r *AssetRepository) ListByIDsAndStatus(ctx context.Context, ids []primitive.ObjectID, status string) ([]model.Asset, error) {
	if len(ids) == 0 {
		return []model.Asset{}, nil
	}
	filter := bson.M{"_id": bson.M{"$in": ids}, "status": status}
	cursor, err := r.coll.Find(ctx, filter)
	if err != nil {
		return nil, fmt.Errorf("list assets by ids and status: %w", err)
	}
	defer cursor.Close(ctx)
	var assets []model.Asset
	if err := cursor.All(ctx, &assets); err != nil {
		return nil, fmt.Errorf("decode assets: %w", err)
	}
	return assets, nil
}

// Count returns the number of assets matching the filter.
func (r *AssetRepository) Count(ctx context.Context, filter bson.M) (int64, error) {
	n, err := r.coll.CountDocuments(ctx, filter)
	if err != nil {
		return 0, fmt.Errorf("count assets: %w", err)
	}
	return n, nil
}

// IncViewCount increments the view counter of an asset.
func (r *AssetRepository) IncViewCount(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"view_count": 1}})
	if err != nil {
		return fmt.Errorf("increment asset view count: %w", err)
	}
	return nil
}

// IncDownloadCount increments the download counter of an asset.
func (r *AssetRepository) IncDownloadCount(ctx context.Context, id primitive.ObjectID) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$inc": bson.M{"download_count": 1}})
	if err != nil {
		return fmt.Errorf("increment asset download count: %w", err)
	}
	return nil
}

// IncrementTagUseCount bumps usage count of tags (upsert).
func (r *AssetRepository) IncrementTagUseCount(ctx context.Context, tags []string) error {
	for _, t := range tags {
		_, err := r.coll.Database().Collection("tags").UpdateOne(ctx,
			bson.M{"name": t},
			bson.M{"$inc": bson.M{"use_count": 1}, "$setOnInsert": bson.M{"created_at": time.Now()}},
			options.Update().SetUpsert(true),
		)
		if err != nil {
			return fmt.Errorf("increment tag use count: %w", err)
		}
	}
	return nil
}
