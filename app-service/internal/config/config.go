package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	// Server configuration
	ServerPort int `envconfig:"SERVER_PORT" default:"8080"`
	Environment string `envconfig:"ENVIRONMENT" default:"development"`

	// Database configuration
	DatabaseURL string `envconfig:"DATABASE_URL" required:"true"`
	DatabaseMaxConns int `envconfig:"DATABASE_MAX_CONNS" default:"25"`
	DatabaseMinConns int `envconfig:"DATABASE_MIN_CONNS" default:"5"`
	DatabaseMaxConnLifetime time.Duration `envconfig:"DATABASE_MAX_CONN_LIFETIME" default:"5m"`

	// Object storage configuration
	StorageEndpoint string `envconfig:"STORAGE_ENDPOINT" required:"true"`
	StorageAccessKey string `envconfig:"STORAGE_ACCESS_KEY" required:"true"`
	StorageSecretKey string `envconfig:"STORAGE_SECRET_KEY" required:"true"`
	StorageBucketName string `envconfig:"STORAGE_BUCKET_NAME" default:"applications"`
	StorageUseSSL bool `envconfig:"STORAGE_USE_SSL" default:"true"`

	// Application configuration
	DownloadURLExpiration time.Duration `envconfig:"DOWNLOAD_URL_EXPIRATION" default:"15m"`
	MaxFileSize int64 `envconfig:"MAX_FILE_SIZE" default:"2147483648"` // 2GB
	AllowedAndroidVersions []string `envconfig:"ALLOWED_ANDROID_VERSIONS" default:"5.0,6.0,7.0,8.0,9.0,10.0,11.0,12.0,13.0,14.0"`
}

// Load loads the configuration from environment variables
func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}
	return &cfg, nil
} 