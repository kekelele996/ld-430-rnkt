package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/errors"
	"github.com/assethub/assethub/internal/model"
	"github.com/assethub/assethub/internal/repository"
)

// TagService handles tag business logic.
type TagService struct {
	repo   *repository.TagRepository
	logger *slog.Logger
}

// NewTagService creates a TagService.
func NewTagService(repo *repository.TagRepository, logger *slog.Logger) *TagService {
	return &TagService{repo: repo, logger: logger}
}

// Create creates a tag.
func (s *TagService) Create(ctx context.Context, name, category string) (*model.Tag, error) {
	if _, err := s.repo.FindByName(ctx, name); err == nil {
		return nil, errors.NewBusinessError(constants.CodeConflict, "tag already exists")
	} else if !errors.IsNotFound(err) {
		return nil, fmt.Errorf("find tag: %w", err)
	}
	tag := &model.Tag{Name: name, Category: category}
	if err := s.repo.Create(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

// TagCloud returns tags ordered by usage count.
func (s *TagService) TagCloud(ctx context.Context, limit int64) ([]model.Tag, error) {
	return s.repo.ListTagCloud(ctx, limit)
}
