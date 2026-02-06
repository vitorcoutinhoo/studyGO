package colaborador

import "go.uber.org/fx"

var Module = fx.Module(
	"colaborador",
	fx.Provide(
		NewColaboradorService,
	),
)
