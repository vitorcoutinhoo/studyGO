package appfx

import (
	"net/http"

	"plantao/internal/api/controller"
	apihttp "plantao/internal/api/http"
	"plantao/internal/api/middleware"
	midware "plantao/internal/api/middleware"
	"plantao/internal/domain/cargo"
	"plantao/internal/domain/colaborador"
	"plantao/internal/domain/comunicacao"
	"plantao/internal/domain/convite"
	"plantao/internal/domain/financeiro"
	"plantao/internal/domain/log"
	"plantao/internal/domain/plantao"
	"plantao/internal/domain/role"
	"plantao/internal/domain/setor"
	"plantao/internal/domain/usuario"
	"plantao/internal/infra/config"
	"plantao/internal/infra/logger"
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
		fx.Annotate(pgstore.NewFeriadoRepository, fx.As(new(financeiro.FeriadoRepository))),
		fx.Annotate(pgstore.NewValorDiaRepository, fx.As(new(financeiro.ValorDiaRepository))),
		fx.Annotate(pgstore.NewColaboradorRepository, fx.As(new(colaborador.ColaboradorRepository))),
		fx.Annotate(pgstore.NewUsuarioRepository, fx.As(new(usuario.UsuarioRepository))),
		fx.Annotate(pgstore.NewModeloRepository, fx.As(new(comunicacao.ModeloComunicaRepository))),
		fx.Annotate(pgstore.NewEnvioRepository, fx.As(new(comunicacao.EnvioComunicacaoRepository))),
		fx.Annotate(pgstore.NewConviteRepository, fx.As(new(convite.ConviteRepository))),
		fx.Annotate(pgstore.NewCargoRepository, fx.As(new(cargo.CargoRepository))),
		fx.Annotate(pgstore.NewSetorRepository, fx.As(new(setor.SetorRepository))),
		fx.Annotate(pgstore.NewRoleRepository, fx.As(new(role.RoleRepository))),
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
		financeiro.NewCalculoService,
		financeiro.NewConfigValorDiaService,
		plantao.NewPlantaoService,
		colaborador.NewColaboradorService,
		usuario.NewUsuarioService,
		usuario.NewAuthService,
		comunicacao.NewModeloComunicacaoService,
		comunicacao.NewEnvioService,
		convite.NewConviteService,
		cargo.NewCargoService,
		setor.NewSetorService,
		role.NewRoleService,
	),
)

var APIModule = fx.Module("api",
	fx.Provide(
		controller.NewPlantaoController,
		controller.NewFeriadoController,
		controller.NewValorDiaController,
		controller.NewColaboradorController,
		controller.NewUsuarioController,
		controller.NewAuthController,
		controller.NewModeloComunicacaoController,
		controller.NewConviteController,
		controller.NewCargoController,
		controller.NewSetorController,
		controller.NewRoleController,
		midware.NewAuthMidware,
		fx.Annotate(
			apihttp.NewRouter,
			fx.As(new(http.Handler)),
			fx.ParamTags(
				``,
				``,
				``,
				``,
				``,
				``,
				``,
				``,
				``,
				``,
				`name:"globalLimiter"`,
				`name:"loginLimiter"`,
				``,
				``,
				``,
				``,
			),
		),
		apihttp.NewServer,
	),
)

var FileModlule = fx.Module("file",
	fx.Provide(
		fx.Annotate(storage.NewLocalStorage, fx.As(new(colaborador.FileStorage))),
	),
)

var LogglerModule = fx.Module("logger",
	fx.Provide(
		fx.Annotate(logger.NewLogger, fx.As(new(log.Logger))),
	),
)

var RateLimitModule = fx.Module("ratelimit",
	fx.Provide(
		fx.Annotate(
			middleware.NewGlobalRateLimiter,
			fx.ResultTags(`name:"globalLimiter"`),
		),
		fx.Annotate(
			middleware.NewLoginRateLimiter,
			fx.ResultTags(`name:"loginLimiter"`),
		),
	),
)
