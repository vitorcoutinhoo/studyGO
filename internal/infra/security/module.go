package security

import (
	"plantao/internal/domain/usuario"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"security",
	fx.Provide(
		fx.Annotate(
			NewBcryptHasher,
			fx.As(new(usuario.PasswordHasher)),
		),
	),
)
