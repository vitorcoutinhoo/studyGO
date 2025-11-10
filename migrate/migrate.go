package migrate

import (
	"main.go/db"
	"main.go/models"
)

func Init() {
	db.ConnectToDB()
	db.DB.AutoMigrate(&models.User{})
	db.DB.AutoMigrate(&models.Colaboradores{})
	db.DB.AutoMigrate(&models.ConfigValoresDias{})
	db.DB.AutoMigrate(&models.ModeloComunicacao{})
	db.DB.AutoMigrate(&models.Feriados{})
	db.DB.AutoMigrate(&models.Plantoes{})
}
