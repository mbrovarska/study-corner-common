package fxconfig

import (
	"os"
	"strconv"
	"study-corner-common/pkg/config"

	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Provide(func() (*config.AppConfig, error) {
		port, _ := strconv.Atoi(os.Getenv("HTTP_PORT"))

		return &config.AppConfig{
			ServiceName: os.Getenv("SERVICE_NAME"),
			ENV: os.Getenv("APP_ENV"),
			HTTPPort: port,
			DB_DSN:   os.Getenv("DB_DSN"),
			LogLevel: os.Getenv("LOG_LEVEL"),
		}, nil
	})
}