package apihttp

import (
	"plantao/internal/api/controller"
	"plantao/internal/api/middleware"
	midware "plantao/internal/api/middleware"
	"plantao/internal/domain/log"
	"plantao/internal/infra/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const (
	ADMIN_ROLE       = "admin"
	COLABORADOR_ROLE = "colaborador"
	GERENTE_ROLE     = "gerente"
)

func NewRouter(
	plantaoController *controller.PlantaoController,
	colaboradorController *controller.ColaboradorController,
	usuarioController *controller.UsuarioController,
	authController *controller.AuthController,
	authMidware *midware.AuthMidware,
	modeloComunicacaoController *controller.ModeloComunicacaoController,
	feriadoController *controller.FeriadoController,
	valorDiaController *controller.ValorDiaController,
	cfg *config.Config,
	log log.Logger,
	globalLimiter *middleware.RateLimiter,
	loginLimiter *middleware.RateLimiter,
	conviteController *controller.ConviteController,
	cargoController *controller.CargoController,
	setorController *controller.SetorController,
	roleController *controller.RoleController,
) *gin.Engine {
	router := gin.Default()

	router.Static("/uploads", cfg.FalePath.Path)

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173", "http://127.0.0.1:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(globalLimiter.Middleware())
	router.Use(middleware.LoggingMiddleware(log))

	setupPlantaoRoutes(router, plantaoController, authMidware)
	setupColaboradorRoutes(router, colaboradorController, authMidware)
	setupUsuarioRoutes(router, usuarioController, authMidware)
	setupAuthRoutes(router, authController, loginLimiter)
	setupModeloComunicacaoRoutes(router, modeloComunicacaoController, authMidware)
	setupFeriadoRoutes(router, feriadoController, authMidware)
	setupValorDiaRoutes(router, valorDiaController, authMidware)
	setupConviteRoutes(router, conviteController, authMidware)
	setupCargoRoutes(router, cargoController, authMidware)
	setupSetorRoutes(router, setorController, authMidware)
	setupRoleRoutes(router, roleController)

	return router
}

func setupPlantaoRoutes(
	router *gin.Engine,
	plantaoController *controller.PlantaoController,
	authMidware *midware.AuthMidware,
) {
	v1 := router.Group("/api/v1")
	{
		plantaoRoutes := v1.Group("/plantoes")
		plantaoRoutes.Use(authMidware.AuthenticationMiddleware(), midware.RoleMidware(ADMIN_ROLE, GERENTE_ROLE, COLABORADOR_ROLE))
		{
			plantaoRoutes.POST("", plantaoController.CreatePlantao)
			plantaoRoutes.PATCH("/:id", midware.RoleMidware(ADMIN_ROLE, GERENTE_ROLE), plantaoController.UpdatePlantao)
			plantaoRoutes.GET("", plantaoController.GetPlantoes)
			plantaoRoutes.GET("/relatorio", midware.RoleMidware(ADMIN_ROLE, GERENTE_ROLE), plantaoController.GetRelatorio)
			plantaoRoutes.GET("/:id", plantaoController.GetPlantaoById)
			plantaoRoutes.DELETE("/:id", plantaoController.DeletePlantao)

			plantaoRoutes.GET("/colaborador/:colaborador_id", plantaoController.GetPlantoesByColaboradorId)
			plantaoRoutes.GET("/status/:status", plantaoController.GetPlantoesByStatus)
			plantaoRoutes.GET("/periodo/:start_date/:end_date", plantaoController.GetPlantoesByPeriodo)

			plantaoRoutes.PATCH("/:id/status", plantaoController.UpdateStatusPlantao)
			plantaoRoutes.POST("/:id/pagamento", midware.RoleMidware(ADMIN_ROLE, GERENTE_ROLE), plantaoController.PagarPlantao)
		}
	}
}

func setupColaboradorRoutes(
	router *gin.Engine,
	colaboradorController *controller.ColaboradorController,
	authMidware *midware.AuthMidware,
) {
	v1 := router.Group("/api/v1")
	{
		colaboradorRoutes := v1.Group("/colaboradores")
		colaboradorRoutes.Use(authMidware.AuthenticationMiddleware(), midware.RoleMidware(ADMIN_ROLE, GERENTE_ROLE, COLABORADOR_ROLE))
		{
			colaboradorRoutes.POST("", colaboradorController.CreateColaborador)
			colaboradorRoutes.PATCH("/:id", colaboradorController.UpdateColaborador)
			colaboradorRoutes.PATCH("/:id/foto", colaboradorController.UploadFotoColaborador)
			colaboradorRoutes.DELETE("/:id", colaboradorController.DisableColaborador)
			colaboradorRoutes.GET("/:id", colaboradorController.GetColaboradorById)
			colaboradorRoutes.GET("", colaboradorController.GetColaboradoresByFilter)
		}
	}
}

func setupConviteRoutes(
	router *gin.Engine,
	conviteController *controller.ConviteController,
	authMidware *midware.AuthMidware,
) {
	v1 := router.Group("/api/v1")
	{
		conviteRoutes := v1.Group("/convites")
		conviteRoutes.Use(authMidware.AuthenticationMiddleware(), middleware.RoleMidware(ADMIN_ROLE))
		{
			conviteRoutes.POST("", conviteController.CreateConvite)
			conviteRoutes.GET("", conviteController.GetAllConvites)
			conviteRoutes.DELETE("/:token", conviteController.DisableConvite)
		}
	}
}

func setupValorDiaRoutes(
	router *gin.Engine,
	valorDiaController *controller.ValorDiaController,
	authMidware *midware.AuthMidware,
) {
	v1 := router.Group("/api/v1")
	{
		valorDiaRoutes := v1.Group("/admin/config-valores")
		valorDiaRoutes.Use(authMidware.AuthenticationMiddleware(), midware.RoleMidware(ADMIN_ROLE))
		{
			valorDiaRoutes.GET("", valorDiaController.GetAll)
			valorDiaRoutes.POST("", valorDiaController.SetValor)
			valorDiaRoutes.PATCH("", valorDiaController.UpdateValor)
			valorDiaRoutes.PATCH("/:tipo_dia", valorDiaController.UpdateValor)
		}
	}
}

func setupFeriadoRoutes(
	router *gin.Engine,
	feriadoController *controller.FeriadoController,
	authMidware *midware.AuthMidware,
) {
	v1 := router.Group("/api/v1")
	{
		feriadoRoutes := v1.Group("/admin/feriados")
		feriadoRoutes.Use(authMidware.AuthenticationMiddleware(), midware.RoleMidware(ADMIN_ROLE))
		{
			feriadoRoutes.GET("", feriadoController.GetFeriadosByAno)
			feriadoRoutes.PATCH("/:id/data", feriadoController.UpdateDataFeriado)
		}
	}
}

func setupUsuarioRoutes(
	router *gin.Engine,
	usuarioController *controller.UsuarioController,
	authMidware *midware.AuthMidware,
) {
	v1 := router.Group("/api/v1")
	{
		adminRoutes := v1.Group("/admin")
		adminRoutes.Use(authMidware.AuthenticationMiddleware(), midware.RoleMidware(ADMIN_ROLE))
		{
			adminRoutes.GET("/all", usuarioController.GetAll)
			adminRoutes.PATCH("/usuarios/:id/role", usuarioController.UpdateRole)
		}

		usuarioAuthRoutes := v1.Group("/authenticated/usuarios")
		usuarioAuthRoutes.Use(authMidware.AuthenticationMiddleware(), midware.RoleMidware(COLABORADOR_ROLE, GERENTE_ROLE, ADMIN_ROLE))
		{
			usuarioAuthRoutes.PUT("", usuarioController.UpdateUsuario)
			usuarioAuthRoutes.GET("", usuarioController.GetUsuarioById)
			usuarioAuthRoutes.DELETE("", usuarioController.DeleteUsuario)
		}

		usuarioRoutes := v1.Group("/usuarios")
		{
			usuarioRoutes.POST("/cadastro", usuarioController.CreateUsuarioByToken)
		}
	}
}

func setupAuthRoutes(
	router *gin.Engine,
	authController *controller.AuthController,
	loginLimiter *middleware.RateLimiter,
) {
	v1 := router.Group("/api/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/login", loginLimiter.Middleware(), authController.Login)
		}
	}
}

func setupCargoRoutes(
	router *gin.Engine,
	cargoController *controller.CargoController,
	authMidware *midware.AuthMidware,
) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/cargos", authMidware.AuthenticationMiddleware(), midware.RoleMidware(ADMIN_ROLE, GERENTE_ROLE, COLABORADOR_ROLE), cargoController.GetAll)

		adminCargos := v1.Group("/admin/cargos")
		adminCargos.Use(authMidware.AuthenticationMiddleware(), midware.RoleMidware(ADMIN_ROLE))
		{
			adminCargos.POST("", cargoController.Create)
			adminCargos.DELETE("/:id", cargoController.Delete)
		}
	}
}

func setupSetorRoutes(
	router *gin.Engine,
	setorController *controller.SetorController,
	authMidware *midware.AuthMidware,
) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/setores", authMidware.AuthenticationMiddleware(), midware.RoleMidware(ADMIN_ROLE, GERENTE_ROLE, COLABORADOR_ROLE), setorController.GetAll)

		adminSetores := v1.Group("/admin/setores")
		adminSetores.Use(authMidware.AuthenticationMiddleware(), midware.RoleMidware(ADMIN_ROLE))
		{
			adminSetores.POST("", setorController.Create)
			adminSetores.DELETE("/:id", setorController.Delete)
		}
	}
}

func setupRoleRoutes(
	router *gin.Engine,
	roleController *controller.RoleController,
) {
	v1 := router.Group("/api/v1")
	{
		v1.GET("/roles", roleController.GetAll)
	}
}

func setupModeloComunicacaoRoutes(
	router *gin.Engine,
	modeloComunicacaoControler *controller.ModeloComunicacaoController,
	authMidware *midware.AuthMidware,
) {
	v1 := router.Group("/api/v1")
	{
		modeloComunicacao := v1.Group("/auth/admin/modelo-comunicacao")
		modeloComunicacao.Use(authMidware.AuthenticationMiddleware(), midware.RoleMidware(ADMIN_ROLE))
		{
			modeloComunicacao.POST("/", modeloComunicacaoControler.CreateModeloComunicacao)
			modeloComunicacao.PUT("/:id_modelo", modeloComunicacaoControler.UpdateModeloComunicacao)
			modeloComunicacao.DELETE("/:id_modelo", modeloComunicacaoControler.DisableModeloComunicacao)
			modeloComunicacao.GET("/", modeloComunicacaoControler.GetAllModelosComunicacao)
			modeloComunicacao.GET("/:id_modelo", modeloComunicacaoControler.GetModeloComunicacaoById)
		}
	}
}
