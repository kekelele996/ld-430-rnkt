package handler

import (
	"strconv"

	"github.com/assethub/assethub/internal/dto"
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/service"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

// CollectionHandler exposes collection endpoints.
type CollectionHandler struct {
	collectionService *service.CollectionService
}

// NewCollectionHandler creates a CollectionHandler.
func NewCollectionHandler(collectionService *service.CollectionService) *CollectionHandler {
	return &CollectionHandler{collectionService: collectionService}
}

// Create handles POST /collections.
func (h *CollectionHandler) Create(c *gin.Context) {
	var req dto.CollectionRequest
	if err := middleware.ValidateJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	userID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	collection, err := h.collectionService.Create(c.Request.Context(), userID, req.Name, req.Description, req.CoverImageURL, req.IsPublic)
	if err != nil {
		_ = c.Error(err)
		return
	}
	middleware.SetAuditAction(c, "collection.create", collection.Name)
	util.Created(c, collection)
}

// List handles GET /collections.
func (h *CollectionHandler) List(c *gin.Context) {
	userID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	onlyPublic, _ := strconv.ParseBool(c.DefaultQuery("only_public", "false"))
	collections, err := h.collectionService.List(c.Request.Context(), userID, onlyPublic)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, collections)
}

// Get handles GET /collections/:id.
func (h *CollectionHandler) Get(c *gin.Context) {
	userID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	collection, err := h.collectionService.Get(c.Request.Context(), c.Param("id"), userID)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, collection)
}

// Update handles PUT /collections/:id.
func (h *CollectionHandler) Update(c *gin.Context) {
	var req dto.CollectionRequest
	if err := middleware.ValidateJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	userID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	if err := h.collectionService.Update(c.Request.Context(), c.Param("id"), userID, req.Name, req.Description, req.CoverImageURL, req.IsPublic); err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"updated": true})
}

// AddAsset handles POST /collections/:id/assets.
func (h *CollectionHandler) AddAsset(c *gin.Context) {
	var req dto.AddAssetRequest
	if err := middleware.ValidateJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	userID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	if err := h.collectionService.AddAsset(c.Request.Context(), c.Param("id"), req.AssetID, userID); err != nil {
		_ = c.Error(err)
		return
	}
	middleware.SetAuditAction(c, "collection.add_asset", req.AssetID)
	util.OK(c, gin.H{"added": true})
}

// RemoveAsset handles DELETE /collections/:id/assets/:assetId.
func (h *CollectionHandler) RemoveAsset(c *gin.Context) {
	userID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	if err := h.collectionService.RemoveAsset(c.Request.Context(), c.Param("id"), c.Param("assetId"), userID); err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"removed": true})
}

// AddCollaborator handles POST /collections/:id/members.
func (h *CollectionHandler) AddCollaborator(c *gin.Context) {
	var req dto.AddCollaboratorRequest
	if err := middleware.ValidateJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	userID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	if err := h.collectionService.AddCollaborator(c.Request.Context(), c.Param("id"), req.UserID, userID); err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"added": true})
}
