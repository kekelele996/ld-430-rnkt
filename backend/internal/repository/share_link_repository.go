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

// ShareLinkRepository provides data access for collection share links.
type ShareLinkRepository struct {
	coll *mongo.Collection
}

// NewShareLinkRepository creates a ShareLinkRepository.
func NewShareLinkRepository(db *mongo.Database) *ShareLinkRepository {
	return &ShareLinkRepository{coll: db.Collection("collection_share_links")}
}

// Create inserts a share link.
func (r *ShareLinkRepository) Create(ctx context.Context, link *model.ShareLink) error {
	link.CreatedAt = time.Now()
	res, err := r.coll.InsertOne(ctx, link)
	if err != nil {
		return fmt.Errorf("create share link: %w", err)
	}
	link.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByTokenHash returns a share link by its token fingerprint.
func (r *ShareLinkRepository) FindByTokenHash(ctx context.Context, tokenHash string) (*model.ShareLink, error) {
	var link model.ShareLink
	err := r.coll.FindOne(ctx, bson.M{"token_hash": tokenHash}).Decode(&link)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find share link by token hash: %w", err)
	}
	return &link, nil
}

// FindByID returns a share link by id.
func (r *ShareLinkRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.ShareLink, error) {
	var link model.ShareLink
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&link)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find share link by id: %w", err)
	}
	return &link, nil
}

// ListByCollection returns the share links of a collection, newest first.
func (r *ShareLinkRepository) ListByCollection(ctx context.Context, collectionID primitive.ObjectID) ([]model.ShareLink, error) {
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.coll.Find(ctx, bson.M{"collection_id": collectionID}, opts)
	if err != nil {
		return nil, fmt.Errorf("list share links: %w", err)
	}
	defer cursor.Close(ctx)
	var links []model.ShareLink
	if err := cursor.All(ctx, &links); err != nil {
		return nil, fmt.Errorf("decode share links: %w", err)
	}
	return links, nil
}

// Revoke marks a share link as revoked. Revoking twice is idempotent.
func (r *ShareLinkRepository) Revoke(ctx context.Context, id primitive.ObjectID) error {
	res, err := r.coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": bson.M{"revoked_at": time.Now()}})
	if err != nil {
		return fmt.Errorf("revoke share link: %w", err)
	}
	if res.MatchedCount == 0 {
		return ErrNotFound
	}
	return nil
}
