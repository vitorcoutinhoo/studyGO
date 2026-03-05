package postgres

import (
	"plantao/internal/domain/colaborador"
	"plantao/internal/domain/plantao"
	"plantao/internal/domain/usuario"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"postgres",
	fx.Provide(
		NewPool,
		fx.Annotate(
			NewPlantaoRepository,
			fx.As(new(plantao.PlantaoRepository)),
		),
		fx.Annotate(
			NewColaboradorRepository,
			fx.As(new(colaborador.ColaboradorRepository)),
		),
		fx.Annotate(
			NewUsuarioRepository,
			fx.As(new(usuario.UsuarioRepository)),
		),
	),
)
