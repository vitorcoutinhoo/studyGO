package api

import (
	"acommerce-api/controllers/user"
	"database/sql"
	"log"

	"github.com/gin-gonic/gin"
)

type APIServer struct {
	address string
	db      *sql.DB
}

func NewApiServer(address string, db *sql.DB) *APIServer {
	return &APIServer{
		address: address,
		db:      db,
	}
}

func (s *APIServer) Start() error {
	r := gin.Default()

	userHandler := user.NewHandler()
	userHandler.RegisterRoutes(r)

	log.Println("Server runing in ", configurations.Config.ServerHost, s.address)

	return r.Run(s.address)
}
