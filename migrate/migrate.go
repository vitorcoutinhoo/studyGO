package migrate

import (
	"main.go/db"
	"main.go/models"
)

func Init() {
	db.ConnectToDB()
	db.DB.AutoMigrate(&models.User{})
	db.DB.AutoMigrate(&models.Colaboradores{})
}
