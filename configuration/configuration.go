package configuration

import (
	"os"

	"github.com/joho/godotenv"
)

type Configuration struct {
	DBUser     string
	DBPassword string
	DBName     string
	DBHost     string
	DBPort     string
	ServerHost string
	ServerPort string
}

var Config = LoadEnvConfig()

func LoadEnvConfig() Configuration {
	godotenv.Load()

	return Configuration{
		DBUser:     os.Getenv("DB_USER"),
		DBPassword: os.Getenv("DB_PASSWORD"),
		DBName:     os.Getenv("DB_NAME"),
		DBHost:     os.Getenv("DB_HOST"),
		DBPort:     os.Getenv("DB_PORT"),
	}
}
