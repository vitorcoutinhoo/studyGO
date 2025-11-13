package api

import (
	"database/sql"

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

	// Dependencias dos usuário
	userService := service.NewUserService(s.db)
	userController := controllers.NewUserController(userService)
	userController.RegisterRoutes(r)

	// Dependencias do colaborador
	colaboradorService := service.NewColaboradorService(s.db)
	colaboradorController := controllers.NewColaboradorController(colaboradorService)
	colaboradorController.RegisterRoutes(r)

	return r.Run(configuration.Config.ServerPort)
}
