package seeds

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/model"
	"github.com/assethub/assethub/internal/repository"
	"github.com/assethub/assethub/internal/util"
	"go.mongodb.org/mongo-driver/mongo"
)

// Seed populates roles, default categories, tags and the admin account if missing.
func Seed(ctx context.Context, db *mongo.Database, logger *slog.Logger, adminEmail, adminPassword string) error {
	roleRepo := repository.NewRoleRepository(db)
	for _, role := range constants.AllRoles {
		if _, err := roleRepo.FindByName(ctx, string(role)); err != nil {
			if !repository.IsNotFound(err) {
				return fmt.Errorf("seed roles: %w", err)
			}
			if err := roleRepo.Create(ctx, &model.Role{
				Name:        string(role),
				Description: string(role) + " role",
				Permissions: defaultPermissions(role),
			}); err != nil {
				return fmt.Errorf("seed role %s: %w", role, err)
			}
		}
	}

	categoryRepo := repository.NewCategoryRepository(db)
	categories, err := categoryRepo.FindAll(ctx)
	if err != nil {
		return fmt.Errorf("seed categories: %w", err)
	}
	if len(categories) == 0 {
		defaultCategories := []struct {
			name string
			icon string
			sort int
			desc string
		}{
			{"图片素材", "image", 1, "图片类数字素材"},
			{"矢量素材", "vector", 2, "矢量图形与插画"},
			{"字体素材", "font", 3, "中英文字体文件"},
			{"模板素材", "template", 4, "演示、设计模板"},
			{"视频素材", "video", 5, "视频与动态素材"},
			{"音频素材", "audio", 6, "音效与背景音乐"},
			{"3D模型", "3d", 7, "三维模型素材"},
		}
		for _, c := range defaultCategories {
			if err := categoryRepo.Create(ctx, &model.Category{Name: c.name, Icon: c.icon, SortOrder: c.sort, Description: c.desc}); err != nil {
				return fmt.Errorf("seed category %s: %w", c.name, err)
			}
		}
		logger.Info("seeded default categories")
	}

	tagRepo := repository.NewTagRepository(db)
	tags, err := tagRepo.ListTagCloud(ctx, 1)
	if err != nil {
		return fmt.Errorf("seed tags: %w", err)
	}
	if len(tags) == 0 {
		defaultTags := []struct {
			name     string
			category string
		}{
			{"极简", string(constants.TagCategoryStyle)},
			{"扁平", string(constants.TagCategoryStyle)},
			{"暖色", string(constants.TagCategoryColor)},
			{"自然", string(constants.TagCategorySubject)},
			{"商务", string(constants.TagCategoryMood)},
			{"科技感", string(constants.TagCategoryMood)},
			{"PNG", string(constants.TagCategoryTechnical)},
			{"SVG", string(constants.TagCategoryTechnical)},
		}
		for _, t := range defaultTags {
			if err := tagRepo.Create(ctx, &model.Tag{Name: t.name, Category: t.category}); err != nil {
				return fmt.Errorf("seed tag %s: %w", t.name, err)
			}
		}
		logger.Info("seeded default tags")
	}

	userRepo := repository.NewUserRepository(db)
	if _, err := userRepo.FindByEmail(ctx, adminEmail); err != nil {
		if !repository.IsNotFound(err) {
			return fmt.Errorf("seed admin: %w", err)
		}
		hash, err := util.HashPassword(adminPassword)
		if err != nil {
			return fmt.Errorf("seed admin hash: %w", err)
		}
		if err := userRepo.Create(ctx, &model.User{
			Email:        adminEmail,
			PasswordHash: hash,
			Nickname:     "管理员",
			Role:         string(constants.RoleAdmin),
		}); err != nil {
			return fmt.Errorf("seed admin user: %w", err)
		}
		logger.Info("seeded admin user", "email", adminEmail)
	}
	return nil
}

func defaultPermissions(role constants.Role) []string {
	switch role {
	case constants.RoleAdmin:
		return []string{"*"}
	case constants.RoleModerator:
		return []string{"asset:review", "asset:publish", "category:manage", "tag:manage"}
	case constants.RoleUploader:
		return []string{"asset:upload", "asset:edit", "collection:manage"}
	default:
		return []string{"asset:view", "asset:download_free"}
	}
}
