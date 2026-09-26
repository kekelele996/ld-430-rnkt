package service

import (
	"context"
	"log/slog"

	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/errors"
	"github.com/assethub/assethub/internal/model"
	"github.com/assethub/assethub/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CollectionService handles collection business logic.
type CollectionService struct {
	repo   *repository.CollectionRepository
	logger *slog.Logger
}

// NewCollectionService creates a CollectionService.
func NewCollectionService(repo *repository.CollectionRepository, logger *slog.Logger) *CollectionService {
	return &CollectionService{repo: repo, logger: logger}
}

// Create creates a collection for the user.
func (s *CollectionService) Create(ctx context.Context, userID primitive.ObjectID, name, description, coverImageURL string, isPublic bool) (*model.Collection, error) {
	collection := &model.Collection{
		Name:            name,
		Description:     description,
		CreatorID:       userID,
		CoverImageURL:   coverImageURL,
		IsPublic:        isPublic,
		AssetIDs:        []primitive.ObjectID{},
		CollaboratorIDs: []primitive.ObjectID{},
	}
	if err := s.repo.Create(ctx, collection); err != nil {
		return nil, err
	}
	return collection, nil
}

// List returns collections visible to the user.
func (s *CollectionService) List(ctx context.Context, userID primitive.ObjectID, onlyPublic bool) ([]model.Collection, error) {
	return s.repo.List(ctx, userID, onlyPublic)
}

// Get returns a collection if the user can access it.
func (s *CollectionService) Get(ctx context.Context, id string, userID primitive.ObjectID) (*model.Collection, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, errors.NewBusinessError(constants.CodeBadRequest, "invalid collection id")
	}
	collection, err := s.repo.FindByID(ctx, oid)
	if err != nil {
		return nil, err
	}
	if !collection.IsPublic && collection.CreatorID != userID && !containsObjectID(collection.CollaboratorIDs, userID) {
		return nil, errors.NewBusinessError(constants.CodeForbidden, constants.MsgForbidden)
	}
	return collection, nil
}

// Update updates collection metadata (owner or collaborator).
func (s *CollectionService) Update(ctx context.Context, id string, userID primitive.ObjectID, name, description, coverImageURL string, isPublic bool) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewBusinessError(constants.CodeBadRequest, "invalid collection id")
	}
	collection, err := s.repo.FindByID(ctx, oid)
	if err != nil {
		return err
	}
	if collection.CreatorID != userID && !containsObjectID(collection.CollaboratorIDs, userID) {
		return errors.NewBusinessError(constants.CodeForbidden, constants.MsgForbidden)
	}
	update := bson.M{"name": name, "description": description, "cover_image_url": coverImageURL, "is_public": isPublic}
	return s.repo.Update(ctx, oid, update)
}

// AddAsset adds an asset to a collection.
func (s *CollectionService) AddAsset(ctx context.Context, id, assetID string, userID primitive.ObjectID) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewBusinessError(constants.CodeBadRequest, "invalid collection id")
	}
	assetOID, err := primitive.ObjectIDFromHex(assetID)
	if err != nil {
		return errors.NewBusinessError(constants.CodeBadRequest, "invalid asset id")
	}
	collection, err := s.repo.FindByID(ctx, oid)
	if err != nil {
		return err
	}
	if collection.CreatorID != userID && !containsObjectID(collection.CollaboratorIDs, userID) {
		return errors.NewBusinessError(constants.CodeForbidden, constants.MsgForbidden)
	}
	return s.repo.PushAsset(ctx, oid, assetOID)
}

// RemoveAsset removes an asset from a collection.
func (s *CollectionService) RemoveAsset(ctx context.Context, id, assetID string, userID primitive.ObjectID) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewBusinessError(constants.CodeBadRequest, "invalid collection id")
	}
	assetOID, err := primitive.ObjectIDFromHex(assetID)
	if err != nil {
		return errors.NewBusinessError(constants.CodeBadRequest, "invalid asset id")
	}
	collection, err := s.repo.FindByID(ctx, oid)
	if err != nil {
		return err
	}
	if collection.CreatorID != userID && !containsObjectID(collection.CollaboratorIDs, userID) {
		return errors.NewBusinessError(constants.CodeForbidden, constants.MsgForbidden)
	}
	return s.repo.PullAsset(ctx, oid, assetOID)
}

// AddCollaborator adds a member to a collection.
func (s *CollectionService) AddCollaborator(ctx context.Context, id, collaboratorID string, userID primitive.ObjectID) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewBusinessError(constants.CodeBadRequest, "invalid collection id")
	}
	memberOID, err := primitive.ObjectIDFromHex(collaboratorID)
	if err != nil {
		return errors.NewBusinessError(constants.CodeBadRequest, "invalid user id")
	}
	collection, err := s.repo.FindByID(ctx, oid)
	if err != nil {
		return err
	}
	if collection.CreatorID != userID {
		return errors.NewBusinessError(constants.CodeForbidden, constants.MsgForbidden)
	}
	return s.repo.PushCollaborator(ctx, oid, memberOID)
}

func containsObjectID(list []primitive.ObjectID, target primitive.ObjectID) bool {
	for _, id := range list {
		if id == target {
			return true
		}
	}
	return false
}
