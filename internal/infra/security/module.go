package security

import (
	"plantao/internal/domain/email"
	"plantao/internal/domain/usuario"
	"plantao/internal/infra/mail"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"security",
	fx.Provide(
		fx.Annotate(
			NewBcryptHasher,
			fx.As(new(usuario.PasswordHasher)),
		),
		NewJWTService,
		fx.Annotate(
			mail.NewSMTPMailer,
			fx.As(new(email.EmailSender)),
		),
	),
)
