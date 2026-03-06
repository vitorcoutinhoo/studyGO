package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Configuração do aplicativo, incluindo servidor e banco de dados
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
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

type JWTConfig struct {
	PrivateKeyPath string
	PublicKeyPath  string
	ExpireTime     int64
}

// Pra teste, seria mais apropriado carregar de um arquivo ou variáveis de ambiente
func LoadConfig() (*Config, error) {
	err := godotenv.Load()

	if err != nil {
		return nil, fmt.Errorf("erro ao carregar o arquivo .env: %w", err)
	}

	expireTime, err := strconv.ParseInt(os.Getenv("EXPIRATION_TIME_MINUTES"), 10, 64)

	if err != nil {
		return nil, fmt.Errorf("erro ao converter EXPIRATION_TIME_MINUTES para int64: %w", err)
	}

	return &Config{
		Server: ServerConfig{
			Port: os.Getenv("SERVER_PORT"),
			Host: os.Getenv("SERVER_HOST"),
		},
		Database: DatabaseConfig{
			URL: os.Getenv("DB_URL"),
		},
		JWT: JWTConfig{
			PrivateKeyPath: os.Getenv("JWT_PRIVATE_KEY_PATH"),
			PublicKeyPath:  os.Getenv("JWT_PUBLIC_KEY_PATH"),
			ExpireTime:     expireTime,
		},
	}, nil
} // Fim LoadConfig
