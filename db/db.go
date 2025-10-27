package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"main.go/configs"
)

// Função para iniciar a conexão com o banco de dados
func SetupDb() (*sql.DB, error) {
	dbInfo := configs.InitDbConfig()                                                                                                     // Inicializa a configuração com o banco de dados
	strConnection, driverName := makePostgresCofiguration(dbInfo.DbUser, dbInfo.DbPassword, dbInfo.DbHost, dbInfo.DbPort, dbInfo.DbName) // Cria string de conexão para o banco de dados

	db, err := sql.Open(driverName, strConnection) // Cria e configura a conexão

	// Verifica erro na criação da conexão
	if err != nil {
		return nil, err // fmt.Errorf("Failed to prepare database connection: %w", err)
	}

	err = db.Ping() // Faz um ping no banco para verificar a conexão

	// Verifica erro no teste de conxão com o banco
	if err != nil {
		db.Close()
		return nil, err // fmt.Errorf("Failed to verify database connection: %w", err)
	}

	fmt.Println("Database connected successfully") // imprime um aviso de que a conexão foi feita com sucesso no terminal

	return db, nil
}

// Função para criar string de conexão para um banco de dados Oracle e retorna no nome do drive que será usado
func makeOracleCofiguration(dbUser, dbPassword, dbHost, dbPort, dbName string) (string, string) {
	return fmt.Sprintf("%s/%s@%s:%s/%s", dbUser, dbPassword, dbHost, dbPort, dbName), "godror"
}

func makePostgresCofiguration(dbUser, dbPassword, dbHost, dbPort, dbName string) (string, string) {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPassword, dbName), "postgres"
}
