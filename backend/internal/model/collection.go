package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Collection is a user-curated asset collection (favorites/folders).
type Collection struct {
	ID              primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name            string               `bson:"name" json:"name"`
	Description     string               `bson:"description" json:"description"`
	CreatorID       primitive.ObjectID   `bson:"creator_id" json:"creator_id"`
	CoverImageURL   string               `bson:"cover_image_url" json:"cover_image_url"`
	AssetIDs        []primitive.ObjectID `bson:"asset_ids" json:"asset_ids"`
	IsPublic        bool                 `bson:"is_public" json:"is_public"`
	CollaboratorIDs []primitive.ObjectID `bson:"collaborator_ids" json:"collaborator_ids"`
	CreatedAt       time.Time            `bson:"created_at" json:"created_at"`
	UpdatedAt       time.Time            `bson:"updated_at" json:"updated_at"`
}
