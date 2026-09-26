package util

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims is the JWT payload for authenticated users.
type Claims struct {
	UserID string `json:"uid"`
	Role   string `json:"role"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// JWTManager signs and verifies JWT tokens.
type JWTManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

// NewJWTManager creates a JWTManager.
func NewJWTManager(secret, issuer string, ttl time.Duration) *JWTManager {
	return &JWTManager{secret: []byte(secret), issuer: issuer, ttl: ttl}
}

// Generate signs a token for the given user.
func (m *JWTManager) Generate(userID, role, email string) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.ttl)
	claims := Claims{
		UserID: userID,
		Role:   role,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			Subject:   userID,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign token: %w", err)
	}
	return signed, expiresAt, nil
}

// Verify parses and validates a token string.
func (m *JWTManager) Verify(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, fmt.Errorf("verify token: %w", err)
	}
	if !token.Valid {
		return nil, fmt.Errorf("verify token: invalid token")
	}
	return claims, nil
}
