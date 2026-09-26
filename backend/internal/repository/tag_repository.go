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

// TagRepository provides data access for tags.
type TagRepository struct {
	coll *mongo.Collection
}

// NewTagRepository creates a TagRepository.
func NewTagRepository(db *mongo.Database) *TagRepository {
	return &TagRepository{coll: db.Collection("tags")}
}

// Create inserts a tag.
func (r *TagRepository) Create(ctx context.Context, tag *model.Tag) error {
	now := time.Now()
	tag.CreatedAt = now
	tag.UpdatedAt = now
	res, err := r.coll.InsertOne(ctx, tag)
	if err != nil {
		return fmt.Errorf("create tag: %w", err)
	}
	tag.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByName returns a tag by name.
func (r *TagRepository) FindByName(ctx context.Context, name string) (*model.Tag, error) {
	var tag model.Tag
	err := r.coll.FindOne(ctx, bson.M{"name": name}).Decode(&tag)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find tag by name: %w", err)
	}
	return &tag, nil
}

// ListTagCloud returns tags ordered by usage count descending (optionally limited).
func (r *TagRepository) ListTagCloud(ctx context.Context, limit int64) ([]model.Tag, error) {
	opts := options.Find().SetSort(bson.D{{Key: "use_count", Value: -1}})
	if limit > 0 {
		opts.SetLimit(limit)
	}
	cursor, err := r.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("list tags: %w", err)
	}
	defer cursor.Close(ctx)
	var tags []model.Tag
	if err := cursor.All(ctx, &tags); err != nil {
		return nil, fmt.Errorf("decode tags: %w", err)
	}
	return tags, nil
}

// UpsertByName increments use_count of a tag, creating it if absent.
func (r *TagRepository) UpsertByName(ctx context.Context, name, category string) error {
	_, err := r.coll.UpdateOne(ctx, bson.M{"name": name}, bson.M{
		"$inc":         bson.M{"use_count": 1},
		"$setOnInsert": bson.M{"category": category, "created_at": time.Now(), "updated_at": time.Now()},
	}, options.Update().SetUpsert(true))
	if err != nil {
		return fmt.Errorf("upsert tag: %w", err)
	}
	return nil
}
