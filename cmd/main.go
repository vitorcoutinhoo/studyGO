package main

import (
	"acommerce-api/cmd/api"
	"acommerce-api/configurations"
	"acommerce-api/db"
	"database/sql"
	"log"
)

func main() {
	db, err := db.NewPostgresDB()

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	initializeDb(db)

	server := api.NewApiServer(configurations.Config.ServerPort, db)

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initializeDb(db *sql.DB) {
	err := db.Ping()

	if err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	log.Println("Successfully connected to database")
}
