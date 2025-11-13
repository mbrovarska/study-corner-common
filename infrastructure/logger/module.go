package logger

import "go.uber.org/fx"

var Module = fx.Module("loger",
	fx.Provide(
		ProvideLogger,
		ProvideDomainLogger,
	),
)

