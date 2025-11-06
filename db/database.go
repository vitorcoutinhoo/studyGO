package db

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"main.go/configuration"
)

var DB *gorm.DB

func ConnectToDB() {
	var err error
	dsn := connectionStringPostgres()
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database!")
	}
}

func connectionStringPostgres() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", configuration.Config.DBHost, configuration.Config.DBUser, configuration.Config.DBPassword, configuration.Config.DBName, configuration.Config.DBPort)
}

func connectionStringOracle() string {
	return fmt.Sprintf("%s/%s@%s:%s/%s", configuration.Config.DBUser, configuration.Config.DBPassword, configuration.Config.DBHost, configuration.Config.DBPort, configuration.Config.DBName)
}
