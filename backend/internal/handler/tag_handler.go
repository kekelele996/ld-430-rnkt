package handler

import (
	"strconv"

	"github.com/assethub/assethub/internal/dto"
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/service"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

// TagHandler exposes tag endpoints.
type TagHandler struct {
	tagService *service.TagService
}

// NewTagHandler creates a TagHandler.
func NewTagHandler(tagService *service.TagService) *TagHandler {
	return &TagHandler{tagService: tagService}
}

// Create handles POST /tags.
func (h *TagHandler) Create(c *gin.Context) {
	var req dto.TagRequest
	if err := middleware.ValidateJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	tag, err := h.tagService.Create(c.Request.Context(), req.Name, req.Category)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.Created(c, tag)
}

// TagCloud handles GET /tags.
func (h *TagHandler) TagCloud(c *gin.Context) {
	limit, _ := strconv.ParseInt(c.DefaultQuery("limit", "50"), 10, 64)
	tags, err := h.tagService.TagCloud(c.Request.Context(), limit)
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, tags)
}
