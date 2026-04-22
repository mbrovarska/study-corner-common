package fxconfig

import (
	"fmt"
	"os"
	"strconv"

	"go.uber.org/fx"

	"study-corner-common/pkg/config"
	"study-corner-common/pkg/db"
)

func NewAppConfig() (*config.AppConfig, error) {
	portStr := os.Getenv("HTTP_PORT")
	if portStr == "" {
		portStr = "8080"
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid HTTP_PORT: %w", err)
	}

	cfg := &config.AppConfig{
		ServiceName: getEnv("SERVICE_NAME", "user-service"),
		Env:         getEnv("APP_ENV", "local"),
		HTTPPort:    port,
		DBDSN:       getEnv("DB_DSN", ""),
		LogLevel:    getEnv("LOG_LEVEL", "info"),
	}

	if cfg.DBDSN == "" {
		return nil, fmt.Errorf("DB_DSN is required")
	}

	return cfg, nil
}

func ProvideDBConfig(cfg *config.AppConfig) db.Config {
	return db.Config{
		DSN:                    cfg.DBDSN,
		MaxOpenConns:           10,
		MaxIdleConns:           5,
		ConnMaxLifetimeSeconds: 300,
	}
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

var Module = fx.Module(
	"config",
	fx.Provide(
		NewAppConfig,
		ProvideDBConfig,
	),
)