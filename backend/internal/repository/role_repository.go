package repository

import (
	"context"
	"fmt"

	"github.com/assethub/assethub/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// RoleRepository provides data access for roles.
type RoleRepository struct {
	coll *mongo.Collection
}

// NewRoleRepository creates a RoleRepository.
func NewRoleRepository(db *mongo.Database) *RoleRepository {
	return &RoleRepository{coll: db.Collection("roles")}
}

// Create inserts a role.
func (r *RoleRepository) Create(ctx context.Context, role *model.Role) error {
	_, err := r.coll.InsertOne(ctx, role)
	if err != nil {
		return fmt.Errorf("create role: %w", err)
	}
	return nil
}

// FindByName returns a role by name.
func (r *RoleRepository) FindByName(ctx context.Context, name string) (*model.Role, error) {
	var role model.Role
	err := r.coll.FindOne(ctx, bson.M{"name": name}).Decode(&role)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find role by name: %w", err)
	}
	return &role, nil
}
