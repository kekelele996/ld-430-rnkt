package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Tag is an asset label with usage statistics.
type Tag struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name      string             `bson:"name" json:"name"`
	Category  string             `bson:"category" json:"category"`
	UseCount  int64              `bson:"use_count" json:"use_count"`
	CreatedAt time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at" json:"updated_at"`
}
