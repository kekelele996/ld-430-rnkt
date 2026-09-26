package service

import (
	"encoding/base64"
	"testing"

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
		raw, err := base64.RawURLEncoding.DecodeString(token)
		if err != nil {
			t.Fatalf("token %q is not valid raw base64url: %v", token, err)
		}
		if len(raw) != shareTokenBytes {
			t.Fatalf("token entropy = %d bytes, want %d", len(raw), shareTokenBytes)
		}
		if seen[token] {
			t.Fatalf("duplicate token generated: %q", token)
		}
		seen[token] = true
	}
}

func TestShareTokenFingerprint(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name  string
		token string
		want  string
	}{
		// Well-known SHA-256 test vectors.
		{name: "empty", token: "", want: "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"},
		{name: "abc", token: "abc", want: "ba7816bf8f01cfea414140de5dae2223b00361a396177a9cb410ff61f20015ad"},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := shareTokenFingerprint(tt.token); got != tt.want {
				t.Fatalf("shareTokenFingerprint(%q) = %q, want %q", tt.token, got, tt.want)
			}
		})
	}
}

func TestOrderAssetsByCollection(t *testing.T) {
	t.Parallel()
	a, b, c := primitive.NewObjectID(), primitive.NewObjectID(), primitive.NewObjectID()
	assets := []model.Asset{
		{ID: c, Title: "c"},
		{ID: a, Title: "a"},
	}
	tests := []struct {
		name string
		ids  []primitive.ObjectID
		want []string
	}{
		{name: "follows collection order", ids: []primitive.ObjectID{a, b, c}, want: []string{"a", "c"}},
		{name: "missing assets dropped", ids: []primitive.ObjectID{b}, want: []string{}},
		{name: "empty collection", ids: nil, want: []string{}},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ordered := orderAssetsByCollection(tt.ids, assets)
			if len(ordered) != len(tt.want) {
				t.Fatalf("orderAssetsByCollection() returned %d assets, want %d", len(ordered), len(tt.want))
			}
			for i, title := range tt.want {
				if ordered[i].Title != title {
					t.Fatalf("ordered[%d].Title = %q, want %q", i, ordered[i].Title, title)
				}
			}
		})
	}
}
