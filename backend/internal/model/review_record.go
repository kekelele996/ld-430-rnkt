package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReviewRecord captures a moderator's review of an asset.
type ReviewRecord struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	AssetID    primitive.ObjectID `bson:"asset_id" json:"asset_id"`
	ReviewerID primitive.ObjectID `bson:"reviewer_id" json:"reviewer_id"`
	Result     string             `bson:"result" json:"result"`
	Comment    string             `bson:"comment" json:"comment"`
	ReviewedAt time.Time          `bson:"reviewed_at" json:"reviewed_at"`
}
