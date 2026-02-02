package api

import (
	"plantao/internal/api/controller"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func NewRouter(
	plantaoController *controller.PlantaoController,
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
