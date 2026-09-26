package router

import (
	"io"
	"log/slog"
	"testing"
	"time"

	"github.com/assethub/assethub/internal/client"
	"github.com/assethub/assethub/internal/config"
	"github.com/assethub/assethub/internal/handler"
	"github.com/assethub/assethub/internal/middleware"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

// TestNewRegistersRoutes ensures the whole route tree assembles without gin
// panicking on conflicting routes (e.g. /collections/:id vs /collections/:id/share-links).
func TestNewRegistersRoutes(t *testing.T) {
	gin.SetMode(gin.TestMode)
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	cfg := &config.Config{}
	cfg.CORS.AllowedOrigins = []string{"*"}
	jwtManager := util.NewJWTManager("test-secret", "test-issuer", time.Hour)
	rateLimiter := middleware.NewRateLimiter(nil, logger)
	handlers := Handlers{
		Health:     handler.NewHealthHandler(nil, nil),
		Auth:       handler.NewAuthHandler(nil),
		Asset:      handler.NewAssetHandler(nil),
		Category:   handler.NewCategoryHandler(nil),
		Collection: handler.NewCollectionHandler(nil),
		ShareLink:  handler.NewShareLinkHandler(nil),
		Download:   handler.NewDownloadHandler(nil),
		Tag:        handler.NewTagHandler(nil),
		Review:     handler.NewReviewHandler(nil),
		Audit:      handler.NewAuditHandler(nil),
	}
	engine := New(handlers, cfg, jwtManager, rateLimiter, &client.MongoClient{}, nil, logger)
	if engine == nil {
		t.Fatal("New() returned nil engine")
	}

	wantRoutes := map[string]bool{
		"POST /api/v1/collections/:id/share-links":           false,
		"GET /api/v1/collections/:id/share-links":            false,
		"DELETE /api/v1/collections/:id/share-links/:linkId": false,
		"GET /api/v1/shared/:token":                          false,
	}
	for _, r := range engine.Routes() {
		key := r.Method + " " + r.Path
		if _, ok := wantRoutes[key]; ok {
			wantRoutes[key] = true
		}
	}
	for route, found := range wantRoutes {
		if !found {
			t.Errorf("route %s not registered", route)
		}
	}
}
