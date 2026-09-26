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
	"github.com/assethub/assethub/internal/dto"
	"github.com/assethub/assethub/internal/errors"
	"github.com/assethub/assethub/internal/model"
	"github.com/assethub/assethub/internal/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// shareTokenBytes is the entropy of a share token before base64url encoding.
const shareTokenBytes = 32

// ShareService handles collection share-link business logic.
type ShareService struct {
	shareRepo      *repository.ShareLinkRepository
	collectionRepo *repository.CollectionRepository
	assetRepo      *repository.AssetRepository
	logger         *slog.Logger
}

// NewShareService creates a ShareService.
func NewShareService(shareRepo *repository.ShareLinkRepository, collectionRepo *repository.CollectionRepository, assetRepo *repository.AssetRepository, logger *slog.Logger) *ShareService {
	return &ShareService{shareRepo: shareRepo, collectionRepo: collectionRepo, assetRepo: assetRepo, logger: logger}
}

// CreateShareLink issues a share link for a collection owned by userID.
// The raw token is returned exactly once; only its fingerprint is persisted.
func (s *ShareService) CreateShareLink(ctx context.Context, collectionID string, userID primitive.ObjectID, expiresInHours int) (*model.ShareLink, string, error) {
	collection, err := s.findOwnedCollection(ctx, collectionID, userID)
	if err != nil {
		return nil, "", err
	}
	token, err := generateShareToken()
	if err != nil {
		return nil, "", err
	}
	link := &model.ShareLink{
		CollectionID: collection.ID,
		TokenHash:    shareTokenFingerprint(token),
		CreatorID:    userID,
		ExpiresAt:    time.Now().Add(time.Duration(expiresInHours) * time.Hour),
	}
	if err := s.shareRepo.Create(ctx, link); err != nil {
		return nil, "", err
	}
	s.logger.Info(constants.LogShareLinkCreated, "collection_id", collection.ID.Hex(), "share_link_id", link.ID.Hex())
	return link, token, nil
}

// ListShareLinks returns the share links of a collection owned by userID.
func (s *ShareService) ListShareLinks(ctx context.Context, collectionID string, userID primitive.ObjectID) ([]model.ShareLink, error) {
	collection, err := s.findOwnedCollection(ctx, collectionID, userID)
	if err != nil {
		return nil, err
	}
	return s.shareRepo.ListByCollectionID(ctx, collection.ID)
}

// RevokeShareLink revokes a share link of a collection owned by userID.
func (s *ShareService) RevokeShareLink(ctx context.Context, collectionID, linkID string, userID primitive.ObjectID) error {
	collection, err := s.findOwnedCollection(ctx, collectionID, userID)
	if err != nil {
		return err
	}
	linkOID, err := primitive.ObjectIDFromHex(linkID)
	if err != nil {
		return errors.NewBusinessError(constants.CodeBadRequest, "invalid share link id")
	}
	link, err := s.shareRepo.FindByIDAndCollection(ctx, linkOID, collection.ID)
	if err != nil {
		if errors.IsNotFound(err) {
			return errors.NewBusinessError(constants.CodeNotFound, constants.MsgNotFound)
		}
		return err
	}
	if link.RevokedAt != nil {
		return nil // already revoked: idempotent
	}
	if err := s.shareRepo.Revoke(ctx, linkOID, collection.ID, time.Now()); err != nil {
		return err
	}
	s.logger.Info(constants.LogShareLinkRevoked, "collection_id", collection.ID.Hex(), "share_link_id", linkOID.Hex())
	return nil
}

// ResolveSharedCollection returns the public view of a collection behind a
// share token. Expired or revoked links are rejected with 403, and only
// published assets are included in the result.
func (s *ShareService) ResolveSharedCollection(ctx context.Context, token string) (*dto.SharedCollectionView, error) {
	if token == "" {
		return nil, errors.NewBusinessError(constants.CodeBadRequest, "missing share token")
	}
	link, err := s.shareRepo.FindByTokenHash(ctx, shareTokenFingerprint(token))
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NewBusinessError(constants.CodeNotFound, constants.MsgNotFound)
		}
		return nil, err
	}
	now := time.Now()
	if link.RevokedAt != nil {
		return nil, errors.NewBusinessError(constants.CodeForbidden, constants.MsgShareLinkRevoked)
	}
	if !now.Before(link.ExpiresAt) {
		return nil, errors.NewBusinessError(constants.CodeForbidden, constants.MsgShareLinkExpired)
	}
	collection, err := s.collectionRepo.FindByID(ctx, link.CollectionID)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NewBusinessError(constants.CodeNotFound, constants.MsgNotFound)
		}
		return nil, err
	}
	// Only published assets are fetched; drafts and archived assets never
	// leave the database layer for public consumers.
	assets, err := s.assetRepo.ListByIDsAndStatus(ctx, collection.AssetIDs, string(constants.AssetStatusPublished))
	if err != nil {
		return nil, err
	}
	views := toSharedAssetViews(assets)
	return &dto.SharedCollectionView{
		Name:          collection.Name,
		Description:   collection.Description,
		CoverImageURL: collection.CoverImageURL,
		AssetCount:    len(views),
		ExpiresAt:     link.ExpiresAt,
		Assets:        views,
	}, nil
}

// ShareLinkStatus computes the lifecycle status of a share link at a given time.
func ShareLinkStatus(link *model.ShareLink, now time.Time) string {
	if link.RevokedAt != nil {
		return string(constants.ShareLinkRevoked)
	}
	if !now.Before(link.ExpiresAt) {
		return string(constants.ShareLinkExpired)
	}
	return string(constants.ShareLinkActive)
}

// findOwnedCollection loads a collection and requires userID to be its creator.
func (s *ShareService) findOwnedCollection(ctx context.Context, collectionID string, userID primitive.ObjectID) (*model.Collection, error) {
	oid, err := primitive.ObjectIDFromHex(collectionID)
	if err != nil {
		return nil, errors.NewBusinessError(constants.CodeBadRequest, "invalid collection id")
	}
	collection, err := s.collectionRepo.FindByID(ctx, oid)
	if err != nil {
		if errors.IsNotFound(err) {
			return nil, errors.NewBusinessError(constants.CodeNotFound, constants.MsgNotFound)
		}
		return nil, err
	}
	if collection.CreatorID != userID {
		return nil, errors.NewBusinessError(constants.CodeForbidden, constants.MsgForbidden)
	}
	return collection, nil
}

// generateShareToken returns a cryptographically random URL-safe token.
func generateShareToken() (string, error) {
	buf := make([]byte, shareTokenBytes)
	if _, err := rand.Read(buf); err != nil {
		return "", fmt.Errorf("generate share token: %w", err)
	}
	return base64.RawURLEncoding.EncodeToString(buf), nil
}

// shareTokenFingerprint returns the SHA-256 fingerprint of a share token.
// Only this fingerprint is persisted, never the raw token.
func shareTokenFingerprint(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

// toSharedAssetViews maps published assets to their public view, skipping any
// asset that is not published as a second line of defense.
func toSharedAssetViews(assets []model.Asset) []dto.SharedAssetView {
	views := make([]dto.SharedAssetView, 0, len(assets))
	for _, asset := range assets {
		if asset.Status != string(constants.AssetStatusPublished) {
			continue
		}
		views = append(views, dto.SharedAssetView{
			ID:            asset.ID.Hex(),
			Title:         asset.Title,
			Description:   asset.Description,
			FileType:      asset.FileType,
			FileFormat:    asset.FileFormat,
			FileURL:       asset.FileURL,
			ThumbnailURL:  asset.ThumbnailURL,
			FileSize:      asset.FileSize,
			Width:         asset.Width,
			Height:        asset.Height,
			Tags:          asset.Tags,
			LicenseType:   asset.LicenseType,
			DownloadCount: asset.DownloadCount,
			ViewCount:     asset.ViewCount,
		})
	}
	return views
}
