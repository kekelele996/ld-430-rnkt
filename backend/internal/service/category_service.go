package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/errors"
	"github.com/assethub/assethub/internal/model"
	"github.com/assethub/assethub/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CategoryService handles category business logic.
type CategoryService struct {
	repo   *repository.CategoryRepository
	logger *slog.Logger
}

// NewCategoryService creates a CategoryService.
func NewCategoryService(repo *repository.CategoryRepository, logger *slog.Logger) *CategoryService {
	return &CategoryService{repo: repo, logger: logger}
}

// Create creates a category.
func (s *CategoryService) Create(ctx context.Context, name, parentID, icon, description string, sortOrder int) (*model.Category, error) {
	category := &model.Category{
		Name:        name,
		Icon:        icon,
		SortOrder:   sortOrder,
		Description: description,
	}
	if parentID != "" {
		pid, err := primitive.ObjectIDFromHex(parentID)
		if err != nil {
			return nil, errors.NewBusinessError(constants.CodeBadRequest, "invalid parent_id")
		}
		if _, err := s.repo.FindByID(ctx, pid); err != nil {
			return nil, fmt.Errorf("find parent category: %w", err)
		}
		category.ParentID = pid
	}
	if err := s.repo.Create(ctx, category); err != nil {
		return nil, err
	}
	return category, nil
}

// Tree returns categories as a nested tree structure.
func (s *CategoryService) Tree(ctx context.Context) ([]*CategoryNode, error) {
	categories, err := s.repo.FindAll(ctx)
	if err != nil {
		return nil, err
	}
	nodes := make(map[primitive.ObjectID]*CategoryNode, len(categories))
	for i := range categories {
		nodes[categories[i].ID] = &CategoryNode{Category: categories[i], Children: []*CategoryNode{}}
	}
	var roots []*CategoryNode
	for _, node := range nodes {
		if node.ParentID != primitive.NilObjectID {
			if parent, ok := nodes[node.ParentID]; ok {
				parent.Children = append(parent.Children, node)
				continue
			}
		}
		roots = append(roots, node)
	}
	return roots, nil
}

// Update updates a category.
func (s *CategoryService) Update(ctx context.Context, id, name, icon, description string, sortOrder int) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewBusinessError(constants.CodeBadRequest, "invalid category id")
	}
	update := bson.M{"name": name, "icon": icon, "description": description, "sort_order": sortOrder}
	if err := s.repo.Update(ctx, oid, update); err != nil {
		return err
	}
	return nil
}

// Delete removes a category.
func (s *CategoryService) Delete(ctx context.Context, id string) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return errors.NewBusinessError(constants.CodeBadRequest, "invalid category id")
	}
	return s.repo.Delete(ctx, oid)
}

// CategoryNode is a tree node with children.
type CategoryNode struct {
	model.Category
	Children []*CategoryNode `json:"children"`
}
