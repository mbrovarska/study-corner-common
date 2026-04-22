package fxlogger

import (
	"context"
	"strings"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"study-corner-common/pkg/config"
	log "study-corner-common/pkg/logger"
)

type Params struct {
	fx.In

	Config    *config.AppConfig
	Lifecycle fx.Lifecycle
}

type zapLogger struct {
	l *zap.Logger
}

func (z *zapLogger) Debug(msg string, fields ...log.Field) {
	z.l.Debug(msg, toZapFields(fields)...)
}

func (z *zapLogger) Info(msg string, fields ...log.Field) {
	z.l.Info(msg, toZapFields(fields)...)
}

func (z *zapLogger) Warn(msg string, fields ...log.Field) {
	z.l.Warn(msg, toZapFields(fields)...)
}

func (z *zapLogger) Error(msg string, fields ...log.Field) {
	z.l.Error(msg, toZapFields(fields)...)
}

func (z *zapLogger) With(fields ...log.Field) log.Logger {
	return &zapLogger{l: z.l.With(toZapFields(fields)...)}
}

func toZapFields(fields []log.Field) []zap.Field {
	out := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		out = append(out, zap.Any(f.Key, f.Value))
	}
	return out
}

func parseLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func New(p Params) (log.Logger, error) {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}
	cfg.Level = zap.NewAtomicLevelAt(parseLevel(p.Config.LogLevel))
	cfg.InitialFields = map[string]interface{}{
		"service": p.Config.ServiceName,
		"env":     p.Config.Env,
	}

	l, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	p.Lifecycle.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			return l.Sync()
		},
	})

	return &zapLogger{l: l}, nil
}

var Module = fx.Module(
	"logger",
	fx.Provide(New),
)