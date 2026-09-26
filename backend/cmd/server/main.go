package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/assethub/assethub/database/seeds"
	"github.com/assethub/assethub/internal/client"
	"github.com/assethub/assethub/internal/config"
	"github.com/assethub/assethub/internal/handler"
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/repository"
	"github.com/assethub/assethub/internal/router"
	"github.com/assethub/assethub/internal/service"
	"github.com/assethub/assethub/internal/util"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load()
	if err != nil {
		logger.Error("load config", "error", err.Error())
		os.Exit(1)
	}
	if err := run(cfg, logger); err != nil {
		logger.Error("server stopped with error", "error", err.Error())
		os.Exit(1)
	}
}

func run(cfg *config.Config, logger *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	// Connect dependencies.
	mongoClient, err := client.NewMongoClient(ctx, cfg.Mongo.URI, cfg.Mongo.Database)
	if err != nil {
		return err
	}
	defer mongoClient.Close(context.Background())

	redisClient, err := client.NewRedisClient(ctx, cfg.Redis.URL, cfg.Redis.Password)
	if err != nil {
		return err
	}
	defer redisClient.Close()

	minioClient, err := client.NewMinioClient(ctx, cfg.Minio.Endpoint, cfg.Minio.AccessKey, cfg.Minio.SecretKey, cfg.Minio.Bucket, cfg.Minio.UseSSL)
	if err != nil {
		return err
	}

	// Seed default data using validated configuration.
	if err := seeds.Seed(ctx, mongoClient.DB, logger, cfg.Admin.Email, cfg.Admin.Password); err != nil {
		return err
	}

	// Assemble services via constructor injection.
	db := mongoClient.DB
	userRepo := repository.NewUserRepository(db)
	roleRepo := repository.NewRoleRepository(db)
	assetRepo := repository.NewAssetRepository(db)
	categoryRepo := repository.NewCategoryRepository(db)
	collectionRepo := repository.NewCollectionRepository(db)
	shareLinkRepo := repository.NewShareLinkRepository(db)
	downloadRepo := repository.NewDownloadRecordRepository(db)
	tagRepo := repository.NewTagRepository(db)
	reviewRepo := repository.NewReviewRecordRepository(db)
	auditRepo := repository.NewAuditLogRepository(db)

	jwtManager := util.NewJWTManager(cfg.JWT.Secret, cfg.JWT.Issuer, cfg.JWT.ExpiresIn)
	authService := service.NewAuthService(userRepo, roleRepo, jwtManager, logger)
	cache := service.NewCacheClient(redisClient)
	storageService := service.NewStorageService(minioClient, logger)
	searchService := service.NewSearchService(assetRepo, cache, logger)
	assetService := service.NewAssetService(assetRepo, tagRepo, categoryRepo, storageService, searchService, logger)
	categoryService := service.NewCategoryService(categoryRepo, logger)
	collectionService := service.NewCollectionService(collectionRepo, logger)
	shareLinkService := service.NewShareLinkService(shareLinkRepo, collectionRepo, assetRepo, logger)
	downloadService := service.NewDownloadService(downloadRepo, assetRepo, logger)
	tagService := service.NewTagService(tagRepo, logger)
	reviewService := service.NewReviewService(reviewRepo, assetRepo, logger)
	auditService := service.NewAuditService(auditRepo, logger)

	rateLimiter := middleware.NewRateLimiter(redisClient, logger)

	handlers := router.Handlers{
		Health:     handler.NewHealthHandler(mongoClient.Client, redisClient.Client),
		Auth:       handler.NewAuthHandler(authService),
		Asset:      handler.NewAssetHandler(assetService),
		Category:   handler.NewCategoryHandler(categoryService),
		Collection: handler.NewCollectionHandler(collectionService),
		ShareLink:  handler.NewShareLinkHandler(shareLinkService),
		Download:   handler.NewDownloadHandler(downloadService),
		Tag:        handler.NewTagHandler(tagService),
		Review:     handler.NewReviewHandler(reviewService),
		Audit:      handler.NewAuditHandler(auditService),
	}

	engine := router.New(handlers, cfg, jwtManager, rateLimiter, mongoClient, redisClient, logger)
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      engine,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	errCh := make(chan error, 1)
	go func() {
		logger.Info("asset library api listening", "port", cfg.Port, "env", cfg.AppEnv)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		logger.Info("shutting down server")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}
