package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"main.go/configuration"
)

func NewConnectToDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", connectionStringPostgres())

	if err != nil {
		return nil, fmt.Errorf("falha ao abrir conexão: %v", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("não foi possível conectar ao banco: %v", err)
	}

	return db, nil
}

func connectionStringPostgres() string {
	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", configuration.Config.DBHost, configuration.Config.DBUser, configuration.Config.DBPassword, configuration.Config.DBName, configuration.Config.DBPort)
}

func connectionStringOracle() string {
	return fmt.Sprintf("%s/%s@%s:%s/%s", configuration.Config.DBUser, configuration.Config.DBPassword, configuration.Config.DBHost, configuration.Config.DBPort, configuration.Config.DBName)
}
