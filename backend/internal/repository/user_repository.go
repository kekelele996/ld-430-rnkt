package repository

import (
	"context"
	"fmt"
	"time"

	"github.com/assethub/assethub/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserRepository provides data access for users.
type UserRepository struct {
	coll *mongo.Collection
}

// NewUserRepository creates a UserRepository.
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{coll: db.Collection("users")}
}

// Create inserts a user.
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	if user.CreatedAt.IsZero() {
		now := time.Now()
		user.CreatedAt = now
		user.UpdatedAt = now
	}
	res, err := r.coll.InsertOne(ctx, user)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}
	user.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// FindByEmail returns a user by email.
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.coll.FindOne(ctx, bson.M{"email": email}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by email: %w", err)
	}
	return &user, nil
}

// FindByID returns a user by id.
func (r *UserRepository) FindByID(ctx context.Context, id primitive.ObjectID) (*model.User, error) {
	var user model.User
	err := r.coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user)
	if err == mongo.ErrNoDocuments {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find user by id: %w", err)
	}
	return &user, nil
}

// CountByRole counts users with the given role.
func (r *UserRepository) CountByRole(ctx context.Context, role string) (int64, error) {
	n, err := r.coll.CountDocuments(ctx, bson.M{"role": role})
	if err != nil {
		return 0, fmt.Errorf("count users by role: %w", err)
	}
	return n, nil
}
