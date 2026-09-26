package service

import (
	"encoding/json"
	"strings"
	"testing"
	"time"

	"github.com/assethub/assethub/internal/constants"
	"github.com/assethub/assethub/internal/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestGenerateShareToken(t *testing.T) {
	t.Parallel()
	seen := make(map[string]bool)
	for i := 0; i < 100; i++ {
		token, err := generateShareToken()
		if err != nil {
			t.Fatalf("generateShareToken() error = %v", err)
		}
		// 32 bytes in raw base64url encode to 43 chars without padding.
		if len(token) != 43 {
			t.Fatalf("token length = %d, want 43", len(token))
		}
		if strings.ContainsAny(token, "+/=") {
			t.Fatalf("token %q is not URL-safe", token)
		}
		if seen[token] {
			t.Fatalf("duplicate token generated: %q", token)
		}
		seen[token] = true
	}
}

func TestShareTokenFingerprint(t *testing.T) {
	t.Parallel()
	// SHA-256 of "abc", well-known test vector.
	const want = "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"
	if got := shareTokenFingerprint("abc"); got != want {
		t.Fatalf("shareTokenFingerprint(abc) = %q, want %q", got, want)
	}
	a, b := shareTokenFingerprint("token-a"), shareTokenFingerprint("token-b")
	if a == b {
		t.Fatal("different tokens must not share a fingerprint")
	}
	if a != shareTokenFingerprint("token-a") {
		t.Fatal("fingerprint must be deterministic")
	}
}

func TestShareLinkStatus(t *testing.T) {
	t.Parallel()
	now := time.Now()
	revokedAt := now.Add(-time.Hour)
	tests := []struct {
		name string
		link *model.ShareLink
		want string
	}{
		{name: "active", link: &model.ShareLink{ExpiresAt: now.Add(time.Hour)}, want: string(constants.ShareLinkActive)},
		{name: "expired", link: &model.ShareLink{ExpiresAt: now.Add(-time.Minute)}, want: string(constants.ShareLinkExpired)},
		{name: "expired at boundary", link: &model.ShareLink{ExpiresAt: now}, want: string(constants.ShareLinkExpired)},
		{name: "revoked", link: &model.ShareLink{ExpiresAt: now.Add(time.Hour), RevokedAt: &revokedAt}, want: string(constants.ShareLinkRevoked)},
		{name: "revoked wins over expired", link: &model.ShareLink{ExpiresAt: now.Add(-time.Hour), RevokedAt: &revokedAt}, want: string(constants.ShareLinkRevoked)},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ShareLinkStatus(tt.link, now); got != tt.want {
				t.Fatalf("ShareLinkStatus() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestToSharedAssetViewsFiltersNonPublished(t *testing.T) {
	t.Parallel()
	assets := []model.Asset{
		{ID: primitive.NewObjectID(), Title: "published hero", Status: string(constants.AssetStatusPublished), FileURL: "https://files/pub.png", ThumbnailURL: "https://files/pub_thumb.png"},
		{ID: primitive.NewObjectID(), Title: "secret draft", Status: string(constants.AssetStatusDraft), FileURL: "https://files/draft.png", ThumbnailURL: "https://files/draft_thumb.png"},
		{ID: primitive.NewObjectID(), Title: "archived asset", Status: string(constants.AssetStatusArchived), FileURL: "https://files/archived.png", ThumbnailURL: "https://files/archived_thumb.png"},
		{ID: primitive.NewObjectID(), Title: "flagged asset", Status: string(constants.AssetStatusFlagged), FileURL: "https://files/flagged.png", ThumbnailURL: "https://files/flagged_thumb.png"},
	}
	views := toSharedAssetViews(assets)
	if len(views) != 1 {
		t.Fatalf("toSharedAssetViews() returned %d views, want 1 (published only)", len(views))
	}
	if views[0].Title != "published hero" {
		t.Fatalf("unexpected asset leaked into public view: %q", views[0].Title)
	}
}

func TestSharedAssetViewJSONExcludesSensitiveFields(t *testing.T) {
	t.Parallel()
	assets := []model.Asset{{
		ID:           primitive.NewObjectID(),
		Title:        "published hero",
		Status:       string(constants.AssetStatusPublished),
		FileURL:      "https://files/pub.png",
		ThumbnailURL: "https://files/pub_thumb.png",
		ObjectKey:    "uploads/user/secret-object-key.png",
		UploaderID:   primitive.NewObjectID(),
	}}
	views := toSharedAssetViews(assets)
	raw, err := json.Marshal(views)
	if err != nil {
		t.Fatalf("marshal views: %v", err)
	}
	for _, forbidden := range []string{"object_key", "uploader_id", "status", "secret-object-key"} {
		if strings.Contains(string(raw), forbidden) {
			t.Fatalf("public view JSON leaks %q: %s", forbidden, raw)
		}
	}
}
