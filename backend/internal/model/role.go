package model

import "go.mongodb.org/mongo-driver/bson/primitive"

// Role is the persisted RBAC role definition.
type Role struct {
	ID          primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Name        string             `bson:"name" json:"name"`
	Description string             `bson:"description" json:"description"`
	Permissions []string           `bson:"permissions" json:"permissions"`
}
