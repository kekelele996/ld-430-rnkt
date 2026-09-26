package config

import (
	"fmt"
	"time"

	"github.com/caarlos0/env/v11"
)

// JWTConfig holds JWT signing settings.
type JWTConfig struct {
	Secret    string        `env:"JWT_SECRET" envDefault:"change_me_to_a_long_random_secret"`
	ExpiresIn time.Duration `env:"JWT_EXPIRES_IN" envDefault:"24h"`
	Issuer    string        `env:"JWT_ISSUER" envDefault:"assethub"`
}

// JWTConfigFromEnv parses JWT settings from environment variables.
func JWTConfigFromEnv() (JWTConfig, error) {
	var c JWTConfig
	if err := env.Parse(&c); err != nil {
		return c, fmt.Errorf("parse jwt config: %w", err)
	}
	if len(c.Secret) < 16 {
		return c, fmt.Errorf("parse jwt config: JWT_SECRET must be at least 16 characters")
	}
	return c, nil
}
