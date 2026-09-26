package config

import "github.com/caarlos0/env/v11"

// MongoConfig holds MongoDB connection settings.
type MongoConfig struct {
	URI      string `env:"MONGO_URI" envDefault:"mongodb://localhost:27017"`
	Database string `env:"MONGO_DATABASE" envDefault:"asset_library"`
	Username string `env:"MONGO_USERNAME" envDefault:""`
	Password string `env:"MONGO_PASSWORD" envDefault:""`
}

// MongoConfigFromEnv parses MongoDB settings from environment variables.
func MongoConfigFromEnv() (MongoConfig, error) {
	var c MongoConfig
	if err := env.Parse(&c); err != nil {
		return c, err
	}
	return c, nil
}
