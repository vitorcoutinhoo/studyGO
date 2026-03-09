package apihttp

import (
	"plantao/internal/api/controller"
	"plantao/internal/api/middleware"
	midware "plantao/internal/api/middleware"
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
) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	router.Use(middleware.RateLimitMiddleware())

	setupPlantaoRoutes(router, plantaoController, authMidware)
	setupColaboradorRoutes(router, colaboradorController, authMidware)
	setupUsuarioRoutes(router, usuarioController, authMidware)
	setupAuthRoutes(router, authController)
	setupModeloComunicacaoRoutes(router, modeloComunicacaoController, authMidware)

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
		plantaoRoutes.Use(authMidware.AuthenticationMidware(), midware.RoleMidware(ADMIN_ROLE, GERENTE_ROLE, COLABORADOR_ROLE))
		{
			plantaoRoutes.POST("", plantaoController.CreatePlantao)
			plantaoRoutes.GET("", plantaoController.GetPlantoes)
			plantaoRoutes.GET("/:id", plantaoController.GetPlantaoById)
			plantaoRoutes.DELETE("/:id", plantaoController.DeletePlantao)

			plantaoRoutes.GET("/colaborador/:colaborador_id", plantaoController.GetPlantoesByColaboradorId)
			plantaoRoutes.GET("/status/:status", plantaoController.GetPlantoesByStatus)
			plantaoRoutes.GET("/periodo/:start_date/:end_date", plantaoController.GetPlantoesByPeriodo)

			plantaoRoutes.PATCH("/:id/status", plantaoController.UpdateStatusPlantao)
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
		colaboradorRoutes.Use(authMidware.AuthenticationMidware(), midware.RoleMidware(ADMIN_ROLE, GERENTE_ROLE, COLABORADOR_ROLE))
		{
			colaboradorRoutes.POST("", colaboradorController.CreateColaborador)
			colaboradorRoutes.PUT("/:id", colaboradorController.UpdateColaborador)
			colaboradorRoutes.DELETE("/:id", colaboradorController.DisableColaborador)
			colaboradorRoutes.GET("/:id", colaboradorController.GetColaboradorById)
			colaboradorRoutes.GET("", colaboradorController.GetColaboradoresByFilter)
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
		adminRoutes.Use(authMidware.AuthenticationMidware(), midware.RoleMidware(ADMIN_ROLE))
		{
			//Rotas do administrador
		}

		usuarioAuthRoutes := v1.Group("/autheticated/usuarios")
		usuarioAuthRoutes.Use(authMidware.AuthenticationMidware(), midware.RoleMidware(COLABORADOR_ROLE))
		{
			usuarioAuthRoutes.PUT("/:id_usuario", usuarioController.UpdateUsuario)
			usuarioAuthRoutes.GET("/:id_usuario", usuarioController.GetUsuarioById)
			usuarioAuthRoutes.DELETE("/:id_usuario", usuarioController.DisableUsuario)
		}

		usuarioRoutes := v1.Group("/usuarios")
		{
			usuarioRoutes.POST("/colaboradores/:id_colaborador", usuarioController.CreateUsuario)
		}
	}
}

func setupAuthRoutes(
	router *gin.Engine,
	authController *controller.AuthController,
) {
	v1 := router.Group("/api/v1")
	{
		authRoutes := v1.Group("/auth")
		{
			authRoutes.POST("/login", authController.Login)
		}
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
		modeloComunicacao.Use(authMidware.AuthenticationMidware(), midware.RoleMidware(ADMIN_ROLE))
		{
			modeloComunicacao.POST("/", modeloComunicacaoControler.CreateModeloComunicacao)
			modeloComunicacao.PUT("/:id_modelo", modeloComunicacaoControler.UpdateModeloComunicacao)
			modeloComunicacao.DELETE("/:id_modelo", modeloComunicacaoControler.DisableModeloComunicacao)
			modeloComunicacao.GET("/", modeloComunicacaoControler.GetAllModelosComunicacao)
			modeloComunicacao.GET("/:id_modelo", modeloComunicacaoControler.GetModeloComunicacaoById)
		}
	}
}
