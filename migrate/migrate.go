package main

import (
	"main.go/db"
	"main.go/models"
)

func init() {
	db.ConnectToDB()
}

func main() {
	db.DB.AutoMigrate(&models.User{})
	db.DB.AutoMigrate(&models.Colaboradores{})
}
