package plantao

import "go.uber.org/fx"

// The Plantao module for dependency injection
var Module = fx.Module(
	"plantao",
	fx.Provide(
		NewPlantaoService,
	),
)
