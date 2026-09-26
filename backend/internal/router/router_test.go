package router

import (
	"log/slog"
	"testing"
	"time"

	"github.com/assethub/assethub/internal/client"
	"github.com/assethub/assethub/internal/config"
	"github.com/assethub/assethub/internal/util"
	"github.com/gin-gonic/gin"
)

// TestNewRegistersShareRoutes ensures the share-link routes register without
// conflicting with the existing /collections/:id routes (gin panics on
// conflicting wildcards at registration time).
func TestNewRegistersShareRoutes(t *testing.T) {
	t.Parallel()
	gin.SetMode(gin.TestMode)
	cfg := &config.Config{}
	cfg.CORS.AllowedOrigins = []string{"*"}
	jwtManager := util.NewJWTManager("test-secret", "test-issuer", time.Hour)

	engine := New(Handlers{}, cfg, jwtManager, nil, &client.MongoClient{}, nil, slog.Default())

	want := map[string]bool{
		"GET /api/v1/shared/:token":                          false,
		"POST /api/v1/collections/:id/share-links":           false,
		"GET /api/v1/collections/:id/share-links":            false,
		"DELETE /api/v1/collections/:id/share-links/:linkId": false,
	}
	for _, route := range engine.Routes() {
		key := route.Method + " " + route.Path
		if _, ok := want[key]; ok {
			want[key] = true
		}
	}
	for key, found := range want {
		if !found {
			t.Errorf("route %s not registered", key)
		}
	}
}
