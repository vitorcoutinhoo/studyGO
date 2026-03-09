package comunicacao

import "go.uber.org/fx"

var Module = fx.Module(
	"comunicacao",
	fx.Provide(
		NewModeloComunicacaoService,
		NewEnvioService,
	),
)
