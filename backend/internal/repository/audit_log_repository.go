package repository

import (
	"context"
	"fmt"

	"github.com/assethub/assethub/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AuditLogRepository provides data access for audit logs.
type AuditLogRepository struct {
	coll *mongo.Collection
}

// NewAuditLogRepository creates an AuditLogRepository.
func NewAuditLogRepository(db *mongo.Database) *AuditLogRepository {
	return &AuditLogRepository{coll: db.Collection("audit_logs")}
}

// Create inserts an audit log entry.
func (r *AuditLogRepository) Create(ctx context.Context, log *model.AuditLog) error {
	res, err := r.coll.InsertOne(ctx, log)
	if err != nil {
		return fmt.Errorf("create audit log: %w", err)
	}
	if oid, ok := res.InsertedID.(interface{ String() string }); ok {
		log.ID.Hex()
		_ = oid
	}
	return nil
}

// List returns audit logs with pagination.
func (r *AuditLogRepository) List(ctx context.Context, skip, limit int64) ([]model.AuditLog, error) {
	opts := options.Find().SetSkip(skip).SetLimit(limit).SetSort(bson.D{{Key: "created_at", Value: -1}})
	cursor, err := r.coll.Find(ctx, bson.M{}, opts)
	if err != nil {
		return nil, fmt.Errorf("list audit logs: %w", err)
	}
	defer cursor.Close(ctx)
	var logs []model.AuditLog
	if err := cursor.All(ctx, &logs); err != nil {
		return nil, fmt.Errorf("decode audit logs: %w", err)
	}
	return logs, nil
}
