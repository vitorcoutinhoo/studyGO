package usuario

import "go.uber.org/fx"

var Module = fx.Module(
	"usuario",
	fx.Provide(
		NewUsuarioService,
	),
)
