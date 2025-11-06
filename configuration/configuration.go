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
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "PLANTAO"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		ServerHost: getEnv("SERVER_HOST", "http://localhost"),
		ServerPort: getEnv("SERVER_PORT", ":8080"),
	}
}

func getEnv(key string, defaultValue string) string {
	value := os.Getenv(key)

	if value != "" {
		return value
	}

	return defaultValue
}
