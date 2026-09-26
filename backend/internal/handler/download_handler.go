package handler

import (
	"strconv"

	"github.com/assethub/assethub/internal/dto"
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/service"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

// DownloadHandler exposes download endpoints.
type DownloadHandler struct {
	downloadService *service.DownloadService
}

// NewDownloadHandler creates a DownloadHandler.
func NewDownloadHandler(downloadService *service.DownloadService) *DownloadHandler {
	return &DownloadHandler{downloadService: downloadService}
}

// Download handles POST /assets/:id/download.
func (h *DownloadHandler) Download(c *gin.Context) {
	var req dto.DownloadRequest
	if err := middleware.ValidateJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	userID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	role := c.GetString(middleware.ContextRoleKey)
	asset, err := h.downloadService.Download(c.Request.Context(), c.Param("id"), userID, req.Purpose, req.LicenseVersion, c.ClientIP(), role)
	if err != nil {
		_ = c.Error(err)
		return
	}
	middleware.SetAuditAction(c, "asset.download", asset.Title)
	util.OK(c, gin.H{"download_url": asset.FileURL, "asset_id": asset.ID.Hex()})
}

// List handles GET /downloads.
func (h *DownloadHandler) List(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "20"), 10, 64)
	skip := util.Offset(page, pageSize)
	pageSize = util.NormalizePageSize(pageSize)
	records, err := h.downloadService.List(c.Request.Context(), skip, pageSize)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, records)
}
