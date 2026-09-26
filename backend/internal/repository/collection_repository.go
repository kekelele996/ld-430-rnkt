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

// CollectionRepository provides data access for collections.
type CollectionRepository struct {
	coll *mongo.Collection
}

// NewCollectionRepository creates a CollectionRepository.
func NewCollectionRepository(db *mongo.Database) *CollectionRepository {
	return &CollectionRepository{coll: db.Collection("collections")}
}

// Create inserts a collection.
func (r *CollectionRepository) Create(ctx context.Context, collection *model.Collection) error {
	now := time.Now()
	collection.CreatedAt = now
	collection.UpdatedAt = now
	res, err := r.coll.InsertOne(ctx, collection)
	if err != nil {
		return fmt.Errorf("create collection: %w", err)
	}
	collection.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByID returns a collection by id.
func (r *CollectionRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.Collection, error) {
	var collection model.Collection
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&collection)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find collection by id: %w", err)
	}
	return &collection, nil
}

// List returns collections visible to a user (own or public).
func (r *CollectionRepository) List(ctx context.Context, userID primitive.ObjectID, onlyPublic bool) ([]model.Collection, error) {
	filter := bson.M{}
	if onlyPublic {
		filter["is_public"] = true
	} else {
		filter["$or"] = bson.A{
			bson.M{"creator_id": userID},
			bson.M{"collaborator_ids": userID},
			bson.M{"is_public": true},
		}
	}
	opts := options.Find().SetSort(bson.D{{Key: "updated_at", Value: -1}})
	cursor, err := r.coll.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("list collections: %w", err)
	}
	defer cursor.Close(ctx)
	var collections []model.Collection
	if err := cursor.All(ctx, &collections); err != nil {
		return nil, fmt.Errorf("decode collections: %w", err)
	}
	return collections, nil
}

// Update updates a collection by id.
func (r *CollectionRepository) Update(ctx context.Context, id primitive.ObjectID, update bson.M) error {
	update["updated_at"] = time.Now()
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("update collection: %w", err)
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// PushAsset adds an asset id to a collection.
func (r *CollectionRepository) PushAsset(ctx context.Context, id, assetID primitive.ObjectID) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$addToSet": bson.M{"asset_ids": assetID}, "$set": bson.M{"updated_at": time.Now()}})
	if err != nil {
		return fmt.Errorf("push asset to collection: %w", err)
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// PullAsset removes an asset id from a collection.
func (r *CollectionRepository) PullAsset(ctx context.Context, id, assetID primitive.ObjectID) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$pull": bson.M{"asset_ids": assetID}, "$set": bson.M{"updated_at": time.Now()}})
	if err != nil {
		return fmt.Errorf("pull asset from collection: %w", err)
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// PushCollaborator adds a collaborator to a collection.
func (r *CollectionRepository) PushCollaborator(ctx context.Context, id, userID primitive.ObjectID) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$addToSet": bson.M{"collaborator_ids": userID}, "$set": bson.M{"updated_at": time.Now()}})
	if err != nil {
		return fmt.Errorf("push collaborator: %w", err)
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}

// Delete removes a collection by id.
func (r *CollectionRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.coll.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return fmt.Errorf("delete collection: %w", err)
	}
	if res.DeletedCount == 0 {
		return ErrNotFound
	}
	return nil
}
