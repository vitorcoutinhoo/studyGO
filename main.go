package main

import (
	"main.go/controllers"
	"main.go/migrate"

	_ "main.go/docs"

	"github.com/gin-gonic/gin"
)

// @title API em Go com Gin e Gorm
// @description API documentada com Swagger
// @host localhost:8080
func main() {
	migrate.Init()
	r := gin.Default()

	// Users routes
	r.POST("/users", controllers.UserCreate)
	r.GET("/users", controllers.UserGetAll)
	r.GET("/users/:id", controllers.UserGetById)
	r.PUT("/users/:id", controllers.UserUpdate)
	r.DELETE("/users/:id", controllers.UserDelete)

	// Colaboradores routes
	r.POST("/colaboradores", controllers.ColaboradorCreate)
	r.GET("/colaboradores", controllers.ColaboradoresGetAll)
	r.GET("/colaboradores/:id", controllers.ColaboradorGetById)
	r.PUT("/colaboradores/:id", controllers.ColaboradorUpdate)
	r.DELETE("/colaboradores/:id", controllers.ColaboradorDelete)

	//r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
