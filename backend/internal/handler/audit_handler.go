package handler

import (
	"strconv"

	"github.com/assethub/assethub/internal/service"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

// AuditHandler exposes audit log endpoints.
type AuditHandler struct {
	auditService *service.AuditService
}

// NewAuditHandler creates an AuditHandler.
func NewAuditHandler(auditService *service.AuditService) *AuditHandler {
	return &AuditHandler{auditService: auditService}
}

// List handles GET /audit-logs.
func (h *AuditHandler) List(c *gin.Context) {
	page, _ := strconv.ParseInt(c.DefaultQuery("page", "1"), 10, 64)
	pageSize, _ := strconv.ParseInt(c.DefaultQuery("page_size", "20"), 10, 64)
	skip := util.Offset(page, pageSize)
	pageSize = util.NormalizePageSize(pageSize)
	logs, err := h.auditService.List(c.Request.Context(), skip, pageSize)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, logs)
}
