package main

import (
	"context"
	"plantao/internal/api"
	"plantao/internal/domain/auth"
	"plantao/internal/domain/colaborador"
	"plantao/internal/domain/comunicacao"
	"plantao/internal/domain/plantao"
	"plantao/internal/domain/usuario"
	"plantao/internal/infra/config"
	"plantao/internal/infra/persistence/postgres"
	"plantao/internal/infra/security"

	"go.uber.org/fx"
)

func main() {
	fx.New(
		config.Module,
		postgres.Module,
		security.Module,
		plantao.Module,
		colaborador.Module,
		usuario.Module,
		auth.Module,
		api.Module,
		comunicacao.Module,

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
