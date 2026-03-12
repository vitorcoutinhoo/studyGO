package appfx

import (
	"net/http"

	"plantao/internal/api/controller"
	midware "plantao/internal/api/middleware"
	apihttp "plantao/internal/api/http"
	"plantao/internal/domain/colaborador"
	"plantao/internal/domain/comunicacao"
	"plantao/internal/domain/financeiro"
	"plantao/internal/domain/plantao"
	"plantao/internal/domain/usuario"
	"plantao/internal/infra/config"
	"plantao/internal/infra/mail"
	pgstore "plantao/internal/infra/persistence/postgres"
	"plantao/internal/infra/security"

	"go.uber.org/fx"
)

var ConfigModule = fx.Module("config",
	fx.Provide(config.LoadConfig),
)

var PostgresModule = fx.Module("postgres",
	fx.Provide(
		pgstore.NewPool,
		fx.Annotate(pgstore.NewPlantaoRepository, fx.As(new(plantao.PlantaoRepository))),
		fx.Annotate(pgstore.NewFeriadoRepository, fx.As(new(financeiro.FeriadoRepository))),
		fx.Annotate(pgstore.NewColaboradorRepository, fx.As(new(colaborador.ColaboradorRepository))),
		fx.Annotate(pgstore.NewUsuarioRepository, fx.As(new(usuario.UsuarioRepository))),
		fx.Annotate(pgstore.NewModeloRepository, fx.As(new(comunicacao.ModeloComunicaRepository))),
		fx.Annotate(pgstore.NewEnvioRepository, fx.As(new(comunicacao.EnvioComunicacaoRepository))),
	),
)

var SecurityModule = fx.Module("security",
	fx.Provide(
		fx.Annotate(security.NewBcryptHasher, fx.As(new(usuario.PasswordHasher))),
		security.NewJWTService,
		func(j *security.JWTService) usuario.TokenGenerator { return j },
		fx.Annotate(mail.NewSMTPMailer, fx.As(new(comunicacao.Mailer))),
	),
)

var DomainModule = fx.Module("domain",
	fx.Provide(
		financeiro.NewFeriadoService,
		plantao.NewPlantaoService,
		colaborador.NewColaboradorService,
		usuario.NewUsuarioService,
		usuario.NewAuthService,
		comunicacao.NewModeloComunicacaoService,
		comunicacao.NewEnvioService,
	),
)

var APIModule = fx.Module("api",
	fx.Provide(
		controller.NewPlantaoController,
		controller.NewFeriadoController,
		controller.NewColaboradorController,
		controller.NewUsuarioController,
		controller.NewAuthController,
		controller.NewModeloComunicacaoController,
		midware.NewAuthMidware,
		fx.Annotate(apihttp.NewRouter, fx.As(new(http.Handler))),
		apihttp.NewServer,
	),
)