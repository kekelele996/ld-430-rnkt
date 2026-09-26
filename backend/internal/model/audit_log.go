package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AuditLog records critical operations for traceability.
type AuditLog struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	OperatorID primitive.ObjectID `bson:"operator_id,omitempty" json:"operator_id,omitempty"`
	Action     string             `bson:"action" json:"action"`
	Resource   string             `bson:"resource" json:"resource"`
	ResourceID string             `bson:"resource_id" json:"resource_id"`
	Detail     string             `bson:"detail" json:"detail"`
	IP         string             `bson:"ip" json:"ip"`
	CreatedAt  time.Time          `bson:"created_at" json:"created_at"`
}
