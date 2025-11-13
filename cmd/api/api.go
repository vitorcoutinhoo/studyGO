package api

import (
	"database/sql"
	"fmt"

	"github.com/gin-gonic/gin"
	"main.go/configuration"
	"main.go/controllers"
	"main.go/service"
)

type ServerAPI struct {
	db *sql.DB
}

func NewServerAPI(db *sql.DB) *ServerAPI {
	return &ServerAPI{
		db: db,
	}
}

func (s *ServerAPI) RunServer() error {
	r := gin.Default()

	userService := service.NewUserService(s.db)
	userController := controllers.NewUserController(userService)
	userController.RegisterRoutes(r)

	fmt.Println(configuration.Config)

	return r.Run(configuration.Config.ServerPort)
}
