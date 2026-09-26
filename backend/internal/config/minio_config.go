package config

import "github.com/caarlos0/env/v11"

// MinioConfig holds MinIO object storage settings.
type MinioConfig struct {
	Endpoint  string `env:"MINIO_ENDPOINT" envDefault:"localhost:9000"`
	AccessKey string `env:"MINIO_ACCESS_KEY" envDefault:"minioadmin"`
	SecretKey string `env:"MINIO_SECRET_KEY" envDefault:"minioadmin"`
	UseSSL    bool   `env:"MINIO_USE_SSL" envDefault:"false"`
	Bucket    string `env:"MINIO_BUCKET" envDefault:"asset-library"`
}

// MinioConfigFromEnv parses MinIO settings from environment variables.
func MinioConfigFromEnv() (MinioConfig, error) {
	var c MinioConfig
	if err := env.Parse(&c); err != nil {
		return c, err
	}
	return c, nil
}
