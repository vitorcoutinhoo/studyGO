package main

import (
	"context"
	"plantao/internal/api"
	"plantao/internal/domain/plantao"
	"plantao/internal/infra/config"
	"plantao/internal/infra/persistence/postgres"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		postgres.Module,
		plantao.Module,
		api.Module,

		fx.Invoke(func(lc fx.Lifecycle, server *api.Server) {
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
