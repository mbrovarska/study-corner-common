package application

import (
	"context"
	"study-corner-common/domain"
	"study-corner-common/infrastructure/logger"

	"go.uber.org/fx"
)


func NewApp(serviceName string, extra ...fx.Option) *fx.App {
	base := fx.Options(
		logger.Module,
		//service name for log fields
		fx.Supply(
			fx.Annotate(serviceName, fx.ResultTags(`name:"service-name"`)),
		),
		//lifecycle
		fx.Invoke(func(lc fx.Lifecycle, log domain.Logger){
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
                      log.Info(ctx, "starting")
					  return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Info(ctx, "stopping")
					_ =  log.Sync()
					return nil
				},
			})

		}),
	)
		return fx.New(append([]fx.Option{base}, extra...)...)
	
}

 