package database

import (
	"fmt"
	"log"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	dsn := connectionStringPostgres()

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao abrir conexão com o banco:", err)
	}

	log.Println("Banco de dados conectado com sucesso")
}

func connectionStringPostgres() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		Config.DB_HOST,
		Config.DB_PORT,
		Config.DB_USER,
		Config.DB_PASSWORD,
		Config.DB_NAME,
	)
}
