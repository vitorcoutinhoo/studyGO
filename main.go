package main

import (
	"main.go/database"
	"main.go/handlers"

	_ "main.go/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API em Go com Gin e Swagger
// @version 1.0
// @description Exemplo básico de API documentada com Swagger
// @host localhost:8080
// @BasePath /
func main() {
	database.Connect()

	r := gin.Default()

	r.GET("/users", handlers.GetUser)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
