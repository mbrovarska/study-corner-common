package domain

import "context"

type Logger interface{
	With(fields ...Field) Logger
	Named(name string) Logger
    Debug(ctx context.Context, msg string, kv ...any)
	Info(ctx context.Context, msg string,  kv ...any)
	Warn(ctx context.Context, msg string,  kv ...any)
	Error(ctx context.Context, msg string, kv ...any)

	Sync() error

}