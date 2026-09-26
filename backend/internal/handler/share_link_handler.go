package handler

import (
	"time"

	"github.com/assethub/assethub/internal/dto"
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/model"
	"github.com/assethub/assethub/internal/service"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

// ShareLinkHandler exposes collection share link endpoints.
type ShareLinkHandler struct {
	shareLinkService *service.ShareLinkService
}

// NewShareLinkHandler creates a ShareLinkHandler.
func NewShareLinkHandler(shareLinkService *service.ShareLinkService) *ShareLinkHandler {
	return &ShareLinkHandler{shareLinkService: shareLinkService}
}

// Create handles POST /collections/:id/share-links.
func (h *ShareLinkHandler) Create(c *gin.Context) {
	var req dto.CreateShareLinkRequest
	if err := middleware.ValidateJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	userID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	link, token, err := h.shareLinkService.Create(c.Request.Context(), c.Param("id"), userID, req.ExpiresInHours)
	if err != nil {
		_ = c.Error(err)
		return
	}
	middleware.SetAuditAction(c, "collection.share_link.create", link.ID.Hex())
	util.Created(c, dto.ShareLinkCreatedResponse{
		ID:        link.ID.Hex(),
		Token:     token,
		ShareURL:  "/api/v1/shared/" + token,
		ExpiresAt: link.ExpiresAt,
	})
}

// List handles GET /collections/:id/share-links.
func (h *ShareLinkHandler) List(c *gin.Context) {
	userID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	links, err := h.shareLinkService.List(c.Request.Context(), c.Param("id"), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, toShareLinkResponses(links))
}

// Revoke handles DELETE /collections/:id/share-links/:linkId.
func (h *ShareLinkHandler) Revoke(c *gin.Context) {
	userID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	if err := h.shareLinkService.Revoke(c.Request.Context(), c.Param("id"), c.Param("linkId"), userID); err != nil {
		_ = c.Error(err)
		return
	}
	middleware.SetAuditAction(c, "collection.share_link.revoke", c.Param("linkId"))
	util.OK(c, gin.H{"revoked": true})
}

// Resolve handles GET /shared/:token. Public endpoint, no login required.
func (h *ShareLinkHandler) Resolve(c *gin.Context) {
	link, collection, assets, err := h.shareLinkService.ResolvePublished(c.Request.Context(), c.Param("token"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, toSharedCollectionView(link, collection, assets))
}

func toShareLinkResponses(links []model.ShareLink) []dto.ShareLinkResponse {
	now := time.Now()
	out := make([]dto.ShareLinkResponse, 0, len(links))
	for _, l := range links {
		out = append(out, dto.ShareLinkResponse{
			ID:        l.ID.Hex(),
			ExpiresAt: l.ExpiresAt,
			RevokedAt: l.RevokedAt,
			Active:    l.IsActive(now),
			CreatedAt: l.CreatedAt,
		})
	}
	return out
}

// toSharedCollectionView maps models to the public view. assets must already be
// restricted to published ones by the service/repository layer.
func toSharedCollectionView(link *model.ShareLink, collection *model.Collection, assets []model.Asset) dto.SharedCollectionView {
	views := make([]dto.SharedAssetView, 0, len(assets))
	for _, a := range assets {
		views = append(views, dto.SharedAssetView{
			ID:           a.ID.Hex(),
			Title:        a.Title,
			Description:  a.Description,
			FileType:     a.FileType,
			FileFormat:   a.FileFormat,
			FileURL:      a.FileURL,
			ThumbnailURL: a.ThumbnailURL,
			FileSize:     a.FileSize,
			Width:        a.Width,
			Height:       a.Height,
			Tags:         a.Tags,
			LicenseType:  a.LicenseType,
		})
	}
	return dto.SharedCollectionView{
		Name:          collection.Name,
		Description:   collection.Description,
		CoverImageURL: collection.CoverImageURL,
		AssetCount:    len(views),
		Assets:        views,
		ExpiresAt:     link.ExpiresAt,
	}
}
