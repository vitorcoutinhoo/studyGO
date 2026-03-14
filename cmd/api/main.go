package main

import (
	"context"
	apihttp "plantao/internal/api/http"
	"plantao/internal/api/middleware"
	appfx "plantao/internal/fx"

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

		fx.Invoke(middleware.StartRateLimitCleanup),

		fx.Invoke(func(lc fx.Lifecycle, server *apihttp.Server) {
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					go server.Start()
					return nil
				},
				OnStop: func(ctx context.Context) error {
					return server.Shutdown(ctx)
				},
			})
		}),
	).Run()
}
