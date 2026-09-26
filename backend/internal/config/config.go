package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

// Config holds all runtime configuration injected through environment variables.
type Config struct {
	AppEnv string `env:"APP_ENV" envDefault:"development"`
	Port   string `env:"SERVER_PORT" envDefault:"8080"`

	Mongo   MongoConfig
	JWT     JWTConfig
	Redis   RedisConfig
	Minio   MinioConfig
	CORS    CORSConfig
	Rate    RateConfig
	Storage StorageConfig
	Admin   AdminConfig
}

type CORSConfig struct {
	AllowedOrigins []string `env:"CORS_ALLOWED_ORIGINS" envDefault:"*" envSeparator:","`
}

type RateConfig struct {
	DownloadPerMinute int `env:"RATE_DOWNLOAD_PER_MINUTE" envDefault:"10"`
	UploadPerHour     int `env:"RATE_UPLOAD_PER_HOUR" envDefault:"20"`
}

type StorageConfig struct {
	AssetBucket string `env:"MINIO_BUCKET" envDefault:"asset-library"`
}

// AdminConfig holds seed administrator settings.
type AdminConfig struct {
	Email    string `env:"ADMIN_EMAIL" envDefault:"admin@assethub.local"`
	Password string `env:"ADMIN_PASSWORD" envDefault:"admin123456"`
}

// Load parses environment variables into Config and validates security-critical values.
func Load() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate rejects insecure settings, especially in production.
func (c *Config) Validate() error {
	if c.AppEnv == "production" {
		if len(c.JWT.Secret) < 32 || c.JWT.Secret == "change_me_to_a_long_random_secret" {
			return fmt.Errorf("load config: JWT_SECRET must be at least 32 characters and must not use the default in production")
		}
		if c.Minio.AccessKey == "minioadmin" || c.Minio.SecretKey == "minioadmin" {
			return fmt.Errorf("load config: MINIO_ACCESS_KEY and MINIO_SECRET_KEY must not use defaults in production")
		}
		if len(c.CORS.AllowedOrigins) == 1 && c.CORS.AllowedOrigins[0] == "*" {
			return fmt.Errorf("load config: CORS_ALLOWED_ORIGINS must not be * in production")
		}
		if len(c.Admin.Password) < 8 {
			return fmt.Errorf("load config: ADMIN_PASSWORD must be at least 8 characters in production")
		}
	} else if len(c.JWT.Secret) < 16 {
		return fmt.Errorf("load config: JWT_SECRET must be at least 16 characters")
	}
	if c.Port == "" {
		return fmt.Errorf("load config: SERVER_PORT is required")
	}
	if c.Mongo.URI == "" {
		return fmt.Errorf("load config: MONGO_URI is required")
	}
	if c.Redis.URL == "" {
		return fmt.Errorf("load config: REDIS_URL is required")
	}
	return nil
}
