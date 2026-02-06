package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Configuração do aplicativo, incluindo servidor e banco de dados
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

// Configurações do servidor, como porta e host
type ServerConfig struct {
	Port string
	Host string
}

// Configurações do banco de dados, como URL de conexão
type DatabaseConfig struct {
	URL string
}

// Pra teste, seria mais apropriado carregar de um arquivo ou variáveis de ambiente
func LoadConfig() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, fmt.Errorf("erro ao carregar o arquivo .env: %w", err)
	}

	return &Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
			Host: os.Getenv("SERVER_HOST"),
		},
		Database: DatabaseConfig{
			URL: os.Getenv("DB_URL"),
		},
	}, nil
} // Fim LoadConfig
