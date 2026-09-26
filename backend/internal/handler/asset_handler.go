package handler

import (
	"strconv"

	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/dto"
	"github.com/assethub/assethub/internal/errors"
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/service"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// AssetHandler exposes asset endpoints.
type AssetHandler struct {
	assetService *service.AssetService
}

// NewAssetHandler creates an AssetHandler.
func NewAssetHandler(assetService *service.AssetService) *AssetHandler {
	return &AssetHandler{assetService: assetService}
}

// Upload handles POST /assets/upload.
func (h *AssetHandler) Upload(c *gin.Context) {
	var req dto.UploadAssetRequest
	if err := middleware.ValidateForm(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	fileHeader, err := util.ValidateUploadedFile(c, "file")
	if err != nil {
		_ = c.Error(errors.NewBusinessError(constants.CodeValidationFailed, err.Error()))
		return
	}
	uploaderID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	asset, err := h.assetService.Upload(c.Request.Context(), uploaderID, service.UploadMeta{
		Title:       req.Title,
		Description: req.Description,
		FileType:    req.FileType,
		FileFormat:  req.FileFormat,
		Tags:        req.Tags,
		CategoryID:  req.CategoryID,
		LicenseType: req.LicenseType,
		Width:       req.Width,
		Height:      req.Height,
	}, fileHeader)
	if err != nil {
		_ = c.Error(err)
		return
	}
	middleware.SetAuditAction(c, "asset.upload", asset.Title)
	util.Created(c, asset)
}

// Update handles PUT /assets/:id.
func (h *AssetHandler) Update(c *gin.Context) {
	var req dto.UpdateAssetRequest
	if err := middleware.ValidateJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	operatorID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	role := c.GetString(middleware.ContextRoleKey)
	if err := h.assetService.Update(c.Request.Context(), c.Param("id"), service.UpdateMeta{
		Title:       req.Title,
		Description: req.Description,
		Tags:        req.Tags,
		CategoryID:  req.CategoryID,
		LicenseType: req.LicenseType,
		Status:      req.Status,
	}, operatorID, role); err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"updated": true})
}

// Get handles GET /assets/:id.
func (h *AssetHandler) Get(c *gin.Context) {
	asset, err := h.assetService.Get(c.Request.Context(), c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, asset)
}

// List handles GET /assets.
func (h *AssetHandler) List(c *gin.Context) {
	var q dto.AssetQuery
	if err := c.ShouldBindQuery(&q); err != nil {
		_ = c.Error(errors.NewValidationError(err))
		return
	}
	page, _ := strconv.ParseInt(q.Page, 10, 64)
	pageSize, _ := strconv.ParseInt(q.PageSize, 10, 64)
	page = util.NormalizePage(page)
	pageSize = util.NormalizePageSize(pageSize)
	result, err := h.assetService.List(c.Request.Context(), service.AssetQueryParams{
		Keyword:     q.Keyword,
		FileType:    q.FileType,
		LicenseType: q.LicenseType,
		Status:      q.Status,
		Tag:         q.Tag,
		CategoryID:  q.CategoryID,
		Sort:        q.Sort,
		Page:        page,
		PageSize:    pageSize,
	})
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, result)
}

// Publish handles POST /assets/:id/publish.
func (h *AssetHandler) Publish(c *gin.Context) {
	operatorID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	role := c.GetString(middleware.ContextRoleKey)
	if err := h.assetService.Publish(c.Request.Context(), c.Param("id"), operatorID, role); err != nil {
		_ = c.Error(err)
		return
	}
	middleware.SetAuditAction(c, "asset.publish", c.Param("id"))
	util.OK(c, gin.H{"status": "published"})
}

// Archive handles POST /assets/:id/archive.
func (h *AssetHandler) Archive(c *gin.Context) {
	operatorID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	role := c.GetString(middleware.ContextRoleKey)
	if err := h.assetService.Archive(c.Request.Context(), c.Param("id"), operatorID, role); err != nil {
		_ = c.Error(err)
		return
	}
	middleware.SetAuditAction(c, "asset.archive", c.Param("id"))
	util.OK(c, gin.H{"status": "archived"})
}

// Hot handles GET /assets/hot.
func (h *AssetHandler) Hot(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "10"), 10, 64)
	assets, err := h.assetService.Hot(c.Request.Context(), limit)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, assets)
}

func mustObjectID(hex string) primitive.ObjectID {
	oid, _ := primitive.ObjectIDFromHex(hex)
	return oid
}
