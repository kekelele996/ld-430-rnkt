package model

import (
	"testing"
	"time"
)

func TestShareLinkIsActive(t *testing.T) {
	t.Parallel()
	now := time.Date(2026, 9, 26, 12, 0, 0, 0, time.UTC)
	revokedAt := now.Add(-time.Hour)
	tests := []struct {
		name string
		link ShareLink
		want bool
	}{
		{name: "active", link: ShareLink{ExpiresAt: now.Add(time.Hour)}, want: true},
		{name: "expires at boundary", link: ShareLink{ExpiresAt: now}, want: false},
		{name: "expired", link: ShareLink{ExpiresAt: now.Add(-time.Minute)}, want: false},
		{name: "revoked", link: ShareLink{ExpiresAt: now.Add(time.Hour), RevokedAt: &revokedAt}, want: false},
		{name: "revoked and expired", link: ShareLink{ExpiresAt: now.Add(-time.Minute), RevokedAt: &revokedAt}, want: false},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := tt.link.IsActive(now); got != tt.want {
				t.Fatalf("IsActive() = %v, want %v", got, tt.want)
			}
		})
	}
}
