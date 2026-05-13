package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type feriadoAPI struct {
	Name        string `json:"name"`
	Date        string `json:"date"`
	Type        string `json:"type"`
	Description string `json:"description"`
}

type apiResponse struct {
	Data []feriadoAPI `json:"data"`
}

func main() {
	ano := time.Now().Year()
	if len(os.Args) > 1 {
		fmt.Sscanf(os.Args[1], "%d", &ano)
	}

	apiKey := os.Getenv("FERIADOS_API_KEY")
	estado := os.Getenv("ESTADO")
	cidade := os.Getenv("CIDADE")
	databaseURL := os.Getenv("DATABASE_URL")

	if apiKey == "" || databaseURL == "" {
		fmt.Println("Variáveis de ambiente obrigatórias: FERIADOS_API_KEY, DATABASE_URL")
		os.Exit(1)
	}

	fmt.Printf("Buscando feriados de %d para %s/%s...\n", ano, cidade, estado)

	feriados, err := buscarFeriados(apiKey, estado, cidade, ano)
	if err != nil {
		fmt.Println("Erro ao buscar feriados:", err)
		os.Exit(1)
	}

	fmt.Printf("%d feriados encontrados. Inserindo no banco...\n\n", len(feriados))

	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		fmt.Println("Erro ao conectar ao banco:", err)
		os.Exit(1)
	}
	defer pool.Close()

	inseridos := 0
	for _, f := range feriados {
		data, err := time.Parse(time.RFC3339, f.Date)
		if err != nil {
			fmt.Printf("Data inválida ignorada: %s\n", f.Date)
			continue
		}

		descricao := f.Description
		if descricao == "" {
			descricao = f.Type
		}

		_, err = pool.Exec(context.Background(),
			`INSERT INTO feriados (data, nome, descricao) VALUES ($1, $2, $3) ON CONFLICT (data) DO NOTHING`,
			data, f.Name, descricao,
		)
		if err != nil {
			fmt.Printf("Erro ao inserir %s (%s): %v\n", f.Date, f.Name, err)
			continue
		}

		fmt.Printf("[%-10s] %s - %s\n", f.Type, data.Format("02/01/2006"), f.Name)
		inseridos++
	}

	fmt.Printf("\n%d/%d feriados de %d inseridos com sucesso.\n", inseridos, len(feriados), ano)
}

func buscarFeriados(apiKey, estado, cidade string, ano int) ([]feriadoAPI, error) {
	url := fmt.Sprintf("https://api.feriados.dev/v1/holidays?year=%d", ano)
	if estado != "" {
		url += "&state=" + estado
	}
	if cidade != "" {
		url += "&city=" + cidade
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("X-API-Key", apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API retornou status %d", resp.StatusCode)
	}

	var result apiResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Data, nil
}
