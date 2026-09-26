package config

import "github.com/caarlos0/env/v11"

// RedisConfig holds Redis connection settings.
type RedisConfig struct {
	URL      string `env:"REDIS_URL" envDefault:"redis://localhost:6379/0"`
	Password string `env:"REDIS_PASSWORD" envDefault:""`
}

// RedisConfigFromEnv parses Redis settings from environment variables.
func RedisConfigFromEnv() (RedisConfig, error) {
	var c RedisConfig
	if err := env.Parse(&c); err != nil {
		return c, err
	}
	return c, nil
}
