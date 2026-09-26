package router

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/assethub/assethub/internal/client"
	"github.com/assethub/assethub/internal/config"
	"github.com/assethub/assethub/internal/handler"
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// Handlers aggregates every HTTP handler.
type Handlers struct {
	Health     *handler.HealthHandler
	Auth       *handler.AuthHandler
	Asset      *handler.AssetHandler
	Category   *handler.CategoryHandler
	Collection *handler.CollectionHandler
	Share      *handler.ShareHandler
	Download   *handler.DownloadHandler
	Tag        *handler.TagHandler
	Review     *handler.ReviewHandler
	Audit      *handler.AuditHandler
}

// New builds the gin engine with all routes and middlewares.
func New(h Handlers, cfg *config.Config, jwtManager *util.JWTManager, rateLimiter *middleware.RateLimiter, db *client.MongoClient, redisClient *client.RedisClient, logger *slog.Logger) *gin.Engine {
	g := gin.New()
	g.Use(gin.Recovery())
	g.Use(middleware.ErrorHandler(logger))
	g.Use(middleware.RequestLogger(logger))
	g.Use(cors.New(corsConfig(cfg)))

	g.GET("/healthz", h.Health.Check)
	g.GET("/readyz", h.Health.Ready)
	g.GET("/api-docs", serveSwaggerHTML)
	g.GET("/swagger/doc.json", serveOpenAPIJSON)

	api := g.Group("/api/v1")
	registerAPI(api, h, cfg, jwtManager, rateLimiter, db)
	// Nginx strips /api/ when proxy_pass ends with a slash; keep /v1 for parity.
	proxyAPI := g.Group("/v1")
	registerAPI(proxyAPI, h, cfg, jwtManager, rateLimiter, db)
	return g
}

func corsConfig(cfg *config.Config) cors.Config {
	allowAll := len(cfg.CORS.AllowedOrigins) == 1 && cfg.CORS.AllowedOrigins[0] == "*"
	origins := cfg.CORS.AllowedOrigins
	if allowAll {
		origins = []string{}
	}
	return cors.Config{
		AllowAllOrigins:  allowAll,
		AllowOrigins:     origins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"X-Request-ID"},
		AllowCredentials: !allowAll,
		MaxAge:           12 * time.Hour,
	}
}

func registerAPI(api *gin.RouterGroup, h Handlers, cfg *config.Config, jwtManager *util.JWTManager, rateLimiter *middleware.RateLimiter, db *client.MongoClient) {
	// Audit logging applies to all protected routes.
	api.Use(middleware.AuditLogger(db.DB))
	registerAuth(api, h, cfg, jwtManager, rateLimiter)
	registerAsset(api, h, cfg, jwtManager, rateLimiter)
	registerCategory(api, h, jwtManager)
	registerCollection(api, h, jwtManager)
	registerShare(api, h, jwtManager)
	registerDownload(api, h, cfg, jwtManager, rateLimiter)
	registerTag(api, h, jwtManager)
	registerReview(api, h, jwtManager)
}

func serveSwaggerHTML(c *gin.Context) {
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="utf-8"/>
  <title>Asset Library API - Swagger UI</title>
  <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css"/>
</head>
<body>
<div id="swagger-ui"></div>
<script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
<script>
  window.onload = function () {
    window.ui = SwaggerUIBundle({ url: '/swagger/doc.json', dom_id: '#swagger-ui' });
  };
</script>
</body>
</html>`)
}

func serveOpenAPIJSON(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	c.File("api/openapi.json")
}
