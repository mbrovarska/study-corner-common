package logger

import (
	"os"
	"study-corner-common/pkg/log"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

//use Params to inject runtime information like service name
type Params struct {
	fx.In
	// Injected via fx.Provide(fx.Annotate(func() string { return "my-service" }, fx.ResultTags(`name:"service-name"`)))
	ServiceName string `name:"service-name"`
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
	out := make([]zap.Field, len(fields))
	for i, f := range fields {
		out[i] = zap.Any(f.Key, f.Value)
	}
	return out
}

func Module() fx.Option {
	return fx.Provide(
		func (p Params) (log.Logger, error) {
			cfg := zap.NewProductionConfig()
			cfg.OutputPaths = []string{"stdout"}
			cfg.InitialFields = map[string]interface{}{
				"service": p.ServiceName,
				"env": os.Getenv("APP_ENV"),
			}
			l, err := cfg.Build()
			if err != nil {
				return nil, err
			}
			return &zapLogger{l: l}, nil

		},
	)
}