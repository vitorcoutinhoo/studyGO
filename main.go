package main

import (
	"main.go/controllers"
	"main.go/database"
	"main.go/models"

	_ "main.go/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title API em Go com Gin e Gorm
// @description API documentada com Swagger
// @host localhost:8080
func main() {
	database.Connect()
	database.DB.AutoMigrate(&models.User{})

	r := gin.Default()

	// Users routes
	r.POST("/users", controllers.UserCreate)
	r.GET("/users", controllers.UserGetAll)
	r.GET("/users/:id", controllers.UserGetById)
	r.PUT("/users/:id", controllers.UserUpdate)
	r.DELETE("/users/:id", controllers.UserDelete)

	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	r.Run(":8080")
}
