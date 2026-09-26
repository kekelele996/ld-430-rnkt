package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/model"
	"github.com/assethub/assethub/internal/repository"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// SearchService combines filters, keyword search, sorting and pagination for assets.
type SearchService struct {
	repo   *repository.AssetRepository
	cache  *CacheClient
	logger *slog.Logger
}

// NewSearchService creates a SearchService.
func NewSearchService(repo *repository.AssetRepository, cache *CacheClient, logger *slog.Logger) *SearchService {
	return &SearchService{repo: repo, cache: cache, logger: logger}
}

// Search returns assets matching the query with total count.
func (s *SearchService) Search(ctx context.Context, q AssetQueryParams) ([]model.Asset, int64, error) {
	filter := bson.M{}
	if q.Status != "" {
		filter["status"] = q.Status
	} else {
		filter["status"] = string(constants.AssetStatusPublished)
	}
	if q.FileType != "" {
		filter["file_type"] = q.FileType
	}
	if q.LicenseType != "" {
		filter["license_type"] = q.LicenseType
	}
	if q.CategoryID != "" {
		if oid, err := primitive.ObjectIDFromHex(q.CategoryID); err == nil {
			filter["category_id"] = oid
		}
	}
	if q.Keyword != "" {
		pattern := primitive.Regex{Pattern: q.Keyword, Options: "i"}
		filter["$or"] = bson.A{
			bson.M{"title": pattern},
			bson.M{"description": pattern},
			bson.M{"tags": pattern},
		}
	}
	if q.Tag != "" {
		filter["tags"] = q.Tag
	}

	sort := bson.D{{Key: "created_at", Value: -1}}
	switch q.Sort {
	case "downloads":
		sort = bson.D{{Key: "download_count", Value: -1}}
	case "views":
		sort = bson.D{{Key: "view_count", Value: -1}}
	case "oldest":
		sort = bson.D{{Key: "created_at", Value: 1}}
	}

	total, err := s.repo.Count(ctx, filter)
	if err != nil {
		return nil, 0, fmt.Errorf("count assets: %w", err)
	}
	page := q.Page
	if page < 1 {
		page = 1
	}
	size := q.PageSize
	if size < 1 {
		size = 20
	}
	if size > 100 {
		size = 100
	}
	assets, err := s.repo.List(ctx, filter, sort, (page-1)*size, size)
	if err != nil {
		return nil, 0, err
	}
	return assets, total, nil
}
