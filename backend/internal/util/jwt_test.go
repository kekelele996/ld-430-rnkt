package util

import (
	"testing"
	"time"
)

func TestJWTManagerGenerateAndVerify(t *testing.T) {
	t.Parallel()
	m := NewJWTManager("0123456789abcdef0123456789abcdef", "assethub-test", time.Hour)
	token, expiresAt, err := m.Generate("507f1f77bcf86cd799439011", "Admin", "admin@assethub.local")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if token == "" {
		t.Fatal("Generate() returned empty token")
	}
	if !expiresAt.After(time.Now()) {
		t.Fatalf("expiresAt should be in the future, got %s", expiresAt)
	}

	claims, err := m.Verify(token)
	if err != nil {
		t.Fatalf("Verify() error = %v", err)
	}
	if claims.UserID != "507f1f77bcf86cd799439011" {
		t.Fatalf("UserID = %q", claims.UserID)
	}
	if claims.Role != "Admin" {
		t.Fatalf("Role = %q", claims.Role)
	}
	if claims.Email != "admin@assethub.local" {
		t.Fatalf("Email = %q", claims.Email)
	}
}

func TestJWTManagerVerifyRejectsBadToken(t *testing.T) {
	t.Parallel()
	m := NewJWTManager("0123456789abcdef0123456789abcdef", "assethub-test", time.Hour)
	for _, tt := range []struct {
		name  string
		token string
	}{
		{name: "empty", token: ""},
		{name: "malformed", token: "a.b.c"},
		{name: "wrong-secret", token: func() string {
			other := NewJWTManager("9999999999abcdef0123456789abcdef", "assethub-test", time.Hour)
			s, _, _ := other.Generate("507f1f77bcf86cd799439011", "Admin", "admin@assethub.local")
			return s
		}()},
	} {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if _, err := m.Verify(tt.token); err == nil {
				t.Fatal("Verify() expected error, got nil")
			}
		})
	}
}

func TestJWTManagerExpiredToken(t *testing.T) {
	t.Parallel()
	m := NewJWTManager("0123456789abcdef0123456789abcdef", "assethub-test", -time.Minute)
	token, _, err := m.Generate("507f1f77bcf86cd799439011", "Admin", "admin@assethub.local")
	if err != nil {
		t.Fatalf("Generate() error = %v", err)
	}
	if _, err := m.Verify(token); err == nil {
		t.Fatal("Verify() expected error for expired token")
	}
}
