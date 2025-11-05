package database

import (
	"fmt"
	"log"
	"os"

	_ "github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Erro ao abrir conexão com o banco:", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("Erro ao obter objeto sql.DB da conexão gorm:", err)
	}

	if err = sqlDB.Ping(); err != nil {
		log.Fatal("Erro ao conectar ao banco:", err)
	}

	log.Println("Banco de dados conectado com sucesso")
}
