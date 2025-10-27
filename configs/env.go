package configs

import (
	"os"

	"github.com/joho/godotenv"
)

// Struct para guarda a configuração de conexão do banco de dados
type DbConfig struct {
	DbUser     string
	DbPassword string
	DbHost     string
	DbPort     string
	DbName     string
}

// Função para inicializar e guardar as informações do banco de dados
func InitDbConfig() DbConfig {
	godotenv.Load()

	return DbConfig{
		DbUser:     getEnv("DB_USER", "root"),
		DbPassword: getEnv("DB_PASSWORD", ""),
		DbHost:     getEnv("DB_HOST", "localhost"),
		DbPort:     getEnv("DB_PORT", "8080"),
		DbName:     getEnv("DB_NAME", "PLANTAO"),
	}
}

// Função para verificar se a variavel de ambiente não está vazia, e caso esteja vazia preenche com um valor padrão
func getEnv(envVar string, defaultValue string) string {
	value := os.Getenv(envVar)

	if value == "" {
		return defaultValue
	}

	return value
}
