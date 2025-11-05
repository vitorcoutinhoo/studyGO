package db

import (
	"acommerce-api/configurations"
	"database/sql"

	_ "github.com/lib/pq"
)

func NewPostgresDB() (*sql.DB, error) {
	db, err := sql.Open("postgres", stringConnectionPostgres())

	if err != nil {
		return nil, err
	}

	return db, nil
}

func stringConnectionPostgres() string {
	return "user=" + configurations.Config.DBUser +
		" password=" + configurations.Config.DBPassword +
		" dbname=" + configurations.Config.DBName +
		" host=" + configurations.Config.DBHost +
		" port=" + configurations.Config.DBPort +
		" sslmode=disable"
}
