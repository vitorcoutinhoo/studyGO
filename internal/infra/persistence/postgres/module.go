package postgres

import (
	"plantao/internal/domain/plantao"

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
	),
)
