package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ShareLink is a time-limited public share link for a collection.
// Only the SHA-256 fingerprint of the token is persisted; the raw token
// is returned to the owner exactly once at creation time.
type ShareLink struct {
	ID           primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	CollectionID primitive.ObjectID `bson:"collection_id" json:"collection_id"`
	TokenHash    string             `bson:"token_hash" json:"-"`
	CreatorID    primitive.ObjectID `bson:"creator_id" json:"creator_id"`
	ExpiresAt    time.Time          `bson:"expires_at" json:"expires_at"`
	RevokedAt    *time.Time         `bson:"revoked_at,omitempty" json:"revoked_at,omitempty"`
	CreatedAt    time.Time          `bson:"created_at" json:"created_at"`
}
