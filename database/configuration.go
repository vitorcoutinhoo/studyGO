package database

import (
	"os"

	"github.com/joho/godotenv"
)

type DBObject struct {
	User     string
	Password string
	DB       string
}

type Configuration struct {
	DB_USER     string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	DB          DBObject
}

var Config = loadEnvConfig()

func loadEnvConfig() Configuration {
	godotenv.Load()

	return Configuration{
		DB_USER:     os.Getenv("DB_USER"),
		DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		DB_HOST:     os.Getenv("DB_HOST"),
		DB_PORT:     os.Getenv("DB_PORT"),
		DB_NAME:     os.Getenv("DB_NAME"),
		DB: DBObject{
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			DB:       os.Getenv("POSTGRES_DB"),
		},
	}
}
