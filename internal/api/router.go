package api

import (
	"plantao/internal/api/controller"
	"plantao/internal/api/midware"
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
) *gin.Engine {
	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	setupPlantaoRoutes(router, plantaoController)
	setupColaboradorRoutes(router, colaboradorController)
	setupUsuarioRoutes(router, usuarioController, authMidware)
	setupAuthRoutes(router, authController)

	return router
}

func setupPlantaoRoutes(
	router *gin.Engine,
	plantaoController *controller.PlantaoController,
) {
	v1 := router.Group("/api/v1")
	{
		plantaoRoutes := v1.Group("/plantoes")
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
) {
	v1 := router.Group("/api/v1")
	{
		colaboradorRoutes := v1.Group("/colaboradores")
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
		adminRoutes.Use(midware.RoleMidware("admin"))
		{
			//adminRoutes.GET("/usuarios", usuarioController.GetAllUsuarios)
		}

		usuarioAuthRoutes := v1.Group("/autheticated/usuarios")
		usuarioAuthRoutes.Use(authMidware.AuthenticationMidware(), midware.RoleMidware(ADMIN_ROLE, GERENTE_ROLE, COLABORADOR_ROLE))
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
