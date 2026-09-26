package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"log/slog"
	"time"

	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/errors"
	"github.com/assethub/assethub/internal/model"
	"github.com/assethub/assethub/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// shareTokenBytes is the entropy size of a share token before base64url encoding.
const shareTokenBytes = 32

// ShareLinkService handles collection share link business logic.
type ShareLinkService struct {
	links       *repository.ShareLinkRepository
	collections *repository.CollectionRepository
	assets      *repository.AssetRepository
	logger      *slog.Logger
}

// NewShareLinkService creates a ShareLinkService.
func NewShareLinkService(links *repository.ShareLinkRepository, collections *repository.CollectionRepository, assets *repository.AssetRepository, logger *slog.Logger) *ShareLinkService {
	return &ShareLinkService{links: links, collections: collections, assets: assets, logger: logger}
}

// Create mints a share link for a collection. Only the collection owner may create one.
// The raw token is returned to the caller exactly once; only its fingerprint is stored.
func (s *ShareLinkService) Create(ctx context.Context, collectionID string, ownerID primitive.ObjectID, expiresInHours int) (*model.ShareLink, string, error) {
	collection, err := s.ownedCollection(ctx, collectionID, ownerID)
	if err != nil {
		return nil, "", err
	}
	token, err := generateShareToken()
	if err != nil {
		return nil, "", err
	}
	link := &model.ShareLink{
		CollectionID: collection.ID,
		CreatorID:    ownerID,
		TokenHash:    shareTokenFingerprint(token),
		ExpiresAt:    time.Now().Add(time.Duration(expiresInHours) * time.Hour),
	}
	if err := s.links.Create(ctx, link); err != nil {
		return nil, "", err
	}
	return link, token, nil
}

// List returns the share links of a collection. Owner only.
func (s *ShareLinkService) List(ctx context.Context, collectionID string, ownerID primitive.ObjectID) ([]model.ShareLink, error) {
	collection, err := s.ownedCollection(ctx, collectionID, ownerID)
	if err != nil {
		return nil, err
	}
	return s.links.ListByCollection(ctx, collection.ID)
}

// Revoke revokes a share link of a collection. Owner only.
func (s *ShareLinkService) Revoke(ctx context.Context, collectionID, linkID string, ownerID primitive.ObjectID) error {
	collection, err := s.ownedCollection(ctx, collectionID, ownerID)
	if err != nil {
		return err
	}
	linkOID, err := primitive.ObjectIDFromHex(linkID)
	if err != nil {
		return errors.NewBusinessError(constants.CodeBadRequest, "invalid share link id")
	}
	link, err := s.links.FindByID(ctx, linkOID)
	if err != nil {
		if errors.IsNotFound(err) {
			return errors.NewBusinessError(constants.CodeNotFound, constants.MsgNotFound)
		}
		return err
	}
	if link.CollectionID != collection.ID {
		return errors.NewBusinessError(constants.CodeNotFound, constants.MsgNotFound)
	}
	return s.links.Revoke(ctx, linkOID)
}

// ResolvePublished resolves a share token into the collection and its published assets.
// Unknown, expired or revoked tokens all yield 403, as does a deleted collection.
func (s *ShareLinkService) ResolvePublished(ctx context.Context, token string) (*model.ShareLink, *model.Collection, []model.Asset, error) {
	link, err := s.links.FindByTokenHash(ctx, shareTokenFingerprint(token))
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil, nil, shareLinkForbidden()
		}
		return nil, nil, nil, err
	}
	if !link.IsActive(time.Now()) {
		return nil, nil, nil, shareLinkForbidden()
	}
	collection, err := s.collections.FindByID(ctx, link.CollectionID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, nil, nil, shareLinkForbidden()
		}
		return nil, nil, nil, err
	}
	assets, err := s.assets.ListPublishedByIDs(ctx, collection.AssetIDs)
	if err != nil {
		return nil, nil, nil, err
	}
	return link, collection, orderAssetsByCollection(collection.AssetIDs, assets), nil
}

// ownedCollection loads a collection and requires the caller to be its owner.
func (s *ShareLinkService) ownedCollection(ctx context.Context, collectionID string, ownerID primitive.ObjectID) (*model.Collection, error) {
	oid, err := primitive.ObjectIDFromHex(collectionID)
	if err != nil {
		return nil, errors.NewBusinessError(constants.CodeBadRequest, "invalid collection id")
	}
	collection, err := s.collections.FindByID(ctx, oid)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NewBusinessError(constants.CodeNotFound, constants.MsgNotFound)
		}
		return nil, err
	}
	if collection.CreatorID != ownerID {
		return nil, errors.NewBusinessError(constants.CodeForbidden, constants.MsgForbidden)
	}
	return collection, nil
}

// shareLinkForbidden is the uniform error for any unusable share token.
func shareLinkForbidden() *errors.BusinessError {
	return errors.NewBusinessError(constants.CodeForbidden, constants.MsgShareLinkInvalid)
}

// generateShareToken returns a URL-safe random token with 256 bits of entropy.
func generateShareToken() (string, error) {
	buf := make([]byte, shareTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate share token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// shareTokenFingerprint returns the SHA-256 hex fingerprint of a share token.
// Only the fingerprint is persisted, never the raw token.
func shareTokenFingerprint(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// orderAssetsByCollection orders assets to match the collection's asset id list.
// Assets missing from the list (e.g. filtered out as non-published) are dropped.
func orderAssetsByCollection(ids []primitive.ObjectID, assets []model.Asset) []model.Asset {
	byID := make(map[primitive.ObjectID]model.Asset, len(assets))
	for _, a := range assets {
		byID[a.ID] = a
	}
	ordered := make([]model.Asset, 0, len(assets))
	for _, id := range ids {
		if a, ok := byID[id]; ok {
			ordered = append(ordered, a)
		}
	}
	return ordered
}
