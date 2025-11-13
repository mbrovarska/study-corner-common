package logger

import (
	"context"
	"study-corner-common/domain"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type ContextFields func(ctx context.Context) []domain.Field

type AdapterParams struct {
	fx.In
	Zap     *zap.Logger
	Extract ContextFields `optional:"true"`
}

// zapDomainLogger adapts zap.SugaredLogger to domain.Logger.
type zapDomainLogger struct {
	s       *zap.SugaredLogger
	extract ContextFields
}

func ProvideDomainLogger(p AdapterParams) domain.Logger {
	extract := p.Extract
	if extract == nil {
		extract = func(context.Context) []domain.Field { return nil }
	}

	return zapDomainLogger{
		s:       p.Zap.Sugar(),
		extract: extract,
	}
}


func (z zapDomainLogger) With(fields ...domain.Field) domain.Logger {
	args := make([]any, 0, len(fields)*2)
	for _, f := range fields {
		args = append(args, f.Key, f.Value)
	}
	return zapDomainLogger{s: z.s.With(args...), extract: z.extract}
}

func (z zapDomainLogger) Named(name string) domain.Logger {
	return zapDomainLogger{s: z.s.Named(name), extract: z.extract}
}

func (z zapDomainLogger) Debug(ctx context.Context, msg string, kv ...any) {
	z.s.Debugw(msg, z.merge(ctx, kv...)...)
}

func (z zapDomainLogger) Info(ctx context.Context, msg string, kv ...any) {
	z.s.Infow(msg, z.merge(ctx, kv...)...)
}

func (z zapDomainLogger) Warn(ctx context.Context, msg string, kv ...any) {
	z.s.Warnw(msg, z.merge(ctx, kv...)...)
}

func (z zapDomainLogger) Error(ctx context.Context, msg string, kv ...any) {
	z.s.Errorw(msg, z.merge(ctx, kv...)...)
}

func (z zapDomainLogger) Sync() error { return z.s.Sync() }

func (z zapDomainLogger) merge(ctx context.Context, kv ...any) []any {
	if z.extract == nil {
		return kv
	}

	cf := z.extract(ctx)
	if len(cf) == 0 {
		return kv
	}

	args := make([]any, 0, len(kv)+len(cf)*2)
	args = append(args, kv...)
	for _, f := range cf {
		args = append(args, f.Key, f.Value)
	}
	return args
}

