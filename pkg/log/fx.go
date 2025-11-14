package log

import (
	"study-corner-common/internal/logger/logger"

	"go.uber.org/fx"
)

func FxModule() fx.Option {
	return logger.Module
}