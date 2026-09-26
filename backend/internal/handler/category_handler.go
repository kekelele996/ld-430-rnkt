package handler

import (
	"github.com/assethub/assethub/internal/dto"
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/service"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

// CategoryHandler exposes category endpoints.
type CategoryHandler struct {
	categoryService *service.CategoryService
}

// NewCategoryHandler creates a CategoryHandler.
func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categoryService: categoryService}
}

// Create handles POST /categories.
func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CategoryRequest
	if err := middleware.ValidateJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	category, err := h.categoryService.Create(c.Request.Context(), req.Name, req.ParentID, req.Icon, req.Description, req.SortOrder)
	if err != nil {
		_ = c.Error(err)
		return
	}
	middleware.SetAuditAction(c, "category.create", category.Name)
	util.Created(c, category)
}

// Tree handles GET /categories.
func (h *CategoryHandler) Tree(c *gin.Context) {
	tree, err := h.categoryService.Tree(c.Request.Context())
	if err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, tree)
}

// Update handles PUT /categories/:id.
func (h *CategoryHandler) Update(c *gin.Context) {
	var req dto.CategoryRequest
	if err := middleware.ValidateJSON(c, &req); err != nil {
		_ = c.Error(err)
		return
	}
	if err := h.categoryService.Update(c.Request.Context(), c.Param("id"), req.Name, req.Icon, req.Description, req.SortOrder); err != nil {
		_ = c.Error(err)
		return
	}
	util.OK(c, gin.H{"updated": true})
}

// Delete handles DELETE /categories/:id.
func (h *CategoryHandler) Delete(c *gin.Context) {
	if err := h.categoryService.Delete(c.Request.Context(), c.Param("id")); err != nil {
		_ = c.Error(err)
		return
	}
	middleware.SetAuditAction(c, "category.delete", c.Param("id"))
	util.OK(c, gin.H{"deleted": true})
}
