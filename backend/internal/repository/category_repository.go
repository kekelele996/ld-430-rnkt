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

// CategoryRepository provides data access for categories.
type CategoryRepository struct {
	coll *mongo.Collection
}

// NewCategoryRepository creates a CategoryRepository.
func NewCategoryRepository(db *mongo.Database) *CategoryRepository {
	return &CategoryRepository{coll: db.Collection("categories")}
}

// Create inserts a category.
func (r *CategoryRepository) Create(ctx context.Context, category *model.Category) error {
	now := time.Now()
	category.CreatedAt = now
	category.UpdatedAt = now
	res, err := r.coll.InsertOne(ctx, category)
	if err != nil {
		return fmt.Errorf("create category: %w", err)
	}
	category.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// FindAll returns all categories sorted by sort_order.
func (r *CategoryRepository) FindAll(ctx context.Context) ([]model.Category, error) {
	opts := options.Find().SetSort(bson.D{{Key: "sort_order", Value: 1}})
	cursor, err := r.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("find categories: %w", err)
	}
	defer cursor.Close(ctx)
	var categories []model.Category
	if err := cursor.All(ctx, &categories); err != nil {
		return nil, fmt.Errorf("decode categories: %w", err)
	}
	return categories, nil
}

// FindByID returns a category by id.
func (r *CategoryRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Category, error) {
	var category model.Category
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&category)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find category by id: %w", err)
	}
	return &category, nil
}

// Update updates a category by id.
func (r *CategoryRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("update category: %w", err)
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete removes a category by id.
func (r *CategoryRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete category: %w", err)
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}
