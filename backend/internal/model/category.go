package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Category is a hierarchical asset category.
type Category struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	ParentID    primitive.ObjectID `bson:"parent_id,omitempty" json:"parent_id,omitempty"`
	Icon        string             `bson:"icon" json:"icon"`
	SortOrder   int                `bson:"sort_order" json:"sort_order"`
	Description string             `bson:"description" json:"description"`
	CreatedAt   time.Time          `bson:"created_at" json:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at" json:"updated_at"`
}
