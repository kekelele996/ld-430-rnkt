package util

import "testing"

func TestHashAndCheckPassword(t *testing.T) {
	t.Parallel()
	hash, err := HashPassword("Str0ngPassw0rd!")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if hash == "" || hash == "Str0ngPassw0rd!" {
		t.Fatal("HashPassword() returned unexpected value")
	}
	if !CheckPassword(hash, "Str0ngPassw0rd!") {
		t.Fatal("CheckPassword() = false, want true")
	}
	if CheckPassword(hash, "wrong-password") {
		t.Fatal("CheckPassword() = true, want false")
	}
}

func TestHashPasswordIsUniquePerCall(t *testing.T) {
	t.Parallel()
	a, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	b, err := HashPassword("password123")
	if err != nil {
		t.Fatalf("HashPassword() error = %v", err)
	}
	if a == b {
		t.Fatal("bcrypt hashes should be salted and unique")
	}
}
