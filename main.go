package main

import (
	"log"

	"main.go/cmd/api"
	"main.go/db"
	_ "main.go/docs"
)

// @title API em Go com Gin e Gorm
// @description API documentada com Swagger
// @host localhost:8080
func main() {
	db, err := db.NewConnectToDB()

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	log.Println("Sucesso ao se conectar com o banco de dados")

	server := api.NewServerAPI(db)

	if err := server.RunServer(); err != nil {
		log.Fatal(err)
	}
}
