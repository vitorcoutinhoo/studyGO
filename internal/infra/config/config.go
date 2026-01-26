package config

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type ServerConfig struct {
	Port string
	Host string
}

type DatabaseConfig struct {
	URL string
}

// Pra teste, seria mais apropriado carregar de um arquivo ou variáveis de ambiente
func LoadConfig() (*Config, error) {
	return &Config{
		Server: ServerConfig{
			Port: "8080",
			Host: "localhost",
		},
		Database: DatabaseConfig{
			URL: "postgres://postgres:password@localhost:5432/repfiles?sslmode=disable",
		},
	}, nil
}