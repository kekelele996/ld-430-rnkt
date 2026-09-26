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

// ShareHandler exposes collection share-link endpoints, including the
// unauthenticated public view consumed by external collaborators.
type ShareHandler struct {
	shareService *service.ShareService
}

// NewShareHandler creates a ShareHandler.
func NewShareHandler(shareService *service.ShareService) *ShareHandler {
	return &ShareHandler{shareService: shareService}
}

// Create handles POST /collections/:id/share-links.
func (h *ShareHandler) Create(c *gin.Context) {
	var req dto.CreateShareLinkRequest
	if err := middleware.ValidateJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	userID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	link, token, err := h.shareService.CreateShareLink(c.Request.Context(), c.Param("id"), userID, req.ExpiresInHours)
	if err != nil {
		_ = c.Error(err)
		return
	}
	middleware.SetAuditAction(c, "collection.share_link.create", c.Param("id"))
	util.Created(c, toShareLinkResponse(link, token))
}

// List handles GET /collections/:id/share-links.
func (h *ShareHandler) List(c *gin.Context) {
	userID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	links, err := h.shareService.ListShareLinks(c.Request.Context(), c.Param("id"), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	responses := make([]dto.ShareLinkResponse, 0, len(links))
	for i := range links {
		responses = append(responses, toShareLinkResponse(&links[i], ""))
	}
	util.OK(c, responses)
}

// Revoke handles DELETE /collections/:id/share-links/:linkId.
func (h *ShareHandler) Revoke(c *gin.Context) {
	userID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	if err := h.shareService.RevokeShareLink(c.Request.Context(), c.Param("id"), c.Param("linkId"), userID); err != nil {
		_ = c.Error(err)
		return
	}
	middleware.SetAuditAction(c, "collection.share_link.revoke", c.Param("linkId"))
	util.OK(c, gin.H{"revoked": true})
}

// Resolve handles GET /shared/:token (no authentication required).
func (h *ShareHandler) Resolve(c *gin.Context) {
	view, err := h.shareService.ResolveSharedCollection(c.Request.Context(), c.Param("token"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, view)
}

// toShareLinkResponse builds the owner-facing response; token is only set at
// creation time because the raw token is never persisted.
func toShareLinkResponse(link *model.ShareLink, token string) dto.ShareLinkResponse {
	return dto.ShareLinkResponse{
		ID:           link.ID.Hex(),
		CollectionID: link.CollectionID.Hex(),
		Token:        token,
		Status:       service.ShareLinkStatus(link, time.Now()),
		ExpiresAt:    link.ExpiresAt,
		RevokedAt:    link.RevokedAt,
		CreatedAt:    link.CreatedAt,
	}
}
