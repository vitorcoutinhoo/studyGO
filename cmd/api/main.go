package main

import (
	"context"
	apihttp "plantao/internal/api/http"
	"plantao/internal/api/middleware"
	"plantao/internal/domain/log"
	appfx "plantao/internal/fx"
	"plantao/internal/worker"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		appfx.ConfigModule,
		appfx.PostgresModule,
		appfx.SecurityModule,
		appfx.DomainModule,
		appfx.APIModule,
		appfx.FileModlule,
		appfx.LogglerModule,
		appfx.RateLimitModule,

		fx.Invoke(middleware.StartRateLimitCleanup),
		fx.Provide(worker.NewPlantaoStatusWorker),
		fx.Invoke(worker.RegisterPlantaoStatusWorker),

		fx.Invoke(func(lc fx.Lifecycle, server *apihttp.Server, log log.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go server.Start()
					log.Info("Servidor iniciando...")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					log.Info("Enserrando aplicação...")
					defer log.Sync()
					return server.Shutdown(ctx)
				},
			})
		}),
	).Run()
}
