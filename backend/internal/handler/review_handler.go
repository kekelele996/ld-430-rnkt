package handler

import (
	"github.com/assethub/assethub/internal/dto"
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/service"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

// ReviewHandler exposes review endpoints.
type ReviewHandler struct {
	reviewService *service.ReviewService
}

// NewReviewHandler creates a ReviewHandler.
func NewReviewHandler(reviewService *service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewService: reviewService}
}

// Review handles POST /assets/:id/reviews.
func (h *ReviewHandler) Review(c *gin.Context) {
	var req dto.ReviewRequest
	if err := middleware.ValidateJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	reviewerID := mustObjectID(c.GetString(middleware.ContextUserIDKey))
	record, err := h.reviewService.Review(c.Request.Context(), c.Param("id"), reviewerID, req.Result, req.Comment)
	if err != nil {
		_ = c.Error(err)
		return
	}
	middleware.SetAuditAction(c, "asset.review", req.Result)
	util.Created(c, record)
}

// ListByAsset handles GET /assets/:id/reviews.
func (h *ReviewHandler) ListByAsset(c *gin.Context) {
	records, err := h.reviewService.ListByAsset(c.Request.Context(), c.Param("id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, records)
}
