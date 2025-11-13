package logger

import (
	"context"
	"fmt"
	"os"
	"strings"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

//use Params to inject runtime information like service name
type Params struct {
	fx.In
	// Injected via fx.Provide(fx.Annotate(func() string { return "my-service" }, fx.ResultTags(`name:"service-name"`)))
	ServiceName string `name:"service-name"`

	//for shutdown sync
	Lifecycle fx.Lifecycle
}

// ProvideLogger creates and returns a configured *zap.Logger instance
func ProvideLogger(p Params) (*zap.Logger, error) {
	// safe default for microservices
	cfg := zap.NewProductionConfig()

	//encoding
	switch strings.ToLower(os.Getenv("LOG_ENCODING")) {
	case "console":
		cfg.Encoding = "console"
	default:
		cfg.Encoding = "json"
	}
	
	//level
	if lvl := strings.ToLower(os.Getenv("LOG_LEVEL")); lvl != ""{
		var level zapcore.Level
		if err := level.Set(lvl); err == nil {
			al := zap.NewAtomicLevel()
			al.SetLevel(level)
			cfg.Level = al
		}
	}

	//machine-readable time encoder 
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime =zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.LevelKey = "level"
	cfg.EncoderConfig.MessageKey = "msg"
	cfg.EncoderConfig.CallerKey = "caller"


	if cfg.Encoding == "console" {
		// Use colored console output for local development ease
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	//build logger
	logger, err := cfg.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zap.ErrorLevel),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to build zap logger: %w", err)
	}


// Immutable contextual fields
	appEnv := os.Getenv("APP_ENV")
	if appEnv == "" {
		appEnv = "production"
	}
	logger = logger.With(
		zap.String("service_name", p.ServiceName),
		zap.String("environment", appEnv),
	)


	// Ensure logs flush on shutdown
	p.Lifecycle.Append(fx.Hook{
		OnStop: func(_ context.Context) error {
			_ = logger.Sync() // ignore EINVAL on some environments
			return nil
		},
	})

	return logger, nil
}