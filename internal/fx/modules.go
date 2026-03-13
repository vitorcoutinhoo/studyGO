package appfx

import (
	"net/http"

	"plantao/internal/api/controller"
	apihttp "plantao/internal/api/http"
	midware "plantao/internal/api/middleware"
	"plantao/internal/domain/colaborador"
	"plantao/internal/domain/comunicacao"
	"plantao/internal/domain/plantao"
	"plantao/internal/domain/usuario"
	"plantao/internal/infra/config"
	"plantao/internal/infra/mail"
	pgstore "plantao/internal/infra/persistence/postgres"
	"plantao/internal/infra/security"
	"plantao/internal/infra/storage"

	"go.uber.org/fx"
)

var ConfigModule = fx.Module("config",
	fx.Provide(config.LoadConfig),
)

var PostgresModule = fx.Module("postgres",
	fx.Provide(
		pgstore.NewPool,
		fx.Annotate(pgstore.NewPlantaoRepository, fx.As(new(plantao.PlantaoRepository))),
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
		controller.NewColaboradorController,
		controller.NewUsuarioController,
		controller.NewAuthController,
		controller.NewModeloComunicacaoController,
		midware.NewAuthMidware,
		fx.Annotate(apihttp.NewRouter, fx.As(new(http.Handler))),
		apihttp.NewServer,
	),
)

var FileModlule = fx.Module("file",
	fx.Provide(
		fx.Annotate(storage.NewLocalStorage, fx.As(new(colaborador.FileStorage))),
	),
)
