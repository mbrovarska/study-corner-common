package logging

import (
	"fmt"
	"os"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//use Params to inject runtime information like service name
type Params struct {
	fx.In
	ServiceName string `name:"service-name"`
}

// ProvideLogger creates and returns a configured *zap.Logger instance
func ProvideLogger(p Params) (*zap.Logger, error) {
	// safe default for microservices
	cfg := zap.NewProductionConfig()

	//cfg JSON encoding
	cfg.Encoding = "json"

	//machine-readable time encoder 
	cfg.EncoderConfig.TimeKey = "timeestamp"
	cfg.EncoderConfig.EncodeTime =zapcore.ISO8601TimeEncoder


    //allow debug/development mode if the environment variable is set.
	// In development, we use console output for readability.
	if os.Getenv("APP_ENV") == "development" {
		// Use colored console output for local development ease
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	//build logger
	logger, err := cfg.Build()
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}

	//add immutable fields to the logger
	logger = logger.With(
		zap.String("service_name", p.ServiceName),
		zap.String("environment", os.Getenv("APP_ENV")),
	)

	return logger, nil
}