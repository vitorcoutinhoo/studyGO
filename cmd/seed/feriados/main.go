package seed

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
	Data string `json:"data"` // DD/MM/YYYY
	Nome string `json:"nome"`
	Tipo string `json:"tipo"` // NACIONAL | ESTADUAL | MUNICIPAL | FACULTATIVO
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
	ibge := os.Getenv("IBGE_CODE")
	databaseURL := os.Getenv("DATABASE_URL")

	if apiKey == "" || ibge == "" || databaseURL == "" {
		fmt.Println("Variáveis de ambiente obrigatórias: FERIADOS_API_KEY, IBGE_CODE, DATABASE_URL")
		os.Exit(1)
	}

	fmt.Printf("Buscando feriados de %d para o município %s...\n", ano, ibge)

	feriados, err := buscarFeriados(apiKey, ibge, ano)
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
		// API retorna DD/MM/YYYY
		data, err := time.Parse("02/01/2006", f.Data)
		if err != nil {
			fmt.Printf("Data inválida ignorada: %s\n", f.Data)
			continue
		}

		_, err = pool.Exec(context.Background(),
			`INSERT INTO feriados (data, nome, descricao) VALUES ($1, $2, $3) ON CONFLICT (data) DO NOTHING`,
			data, f.Nome, f.Tipo,
		)
		if err != nil {
			fmt.Printf("Erro ao inserir %s (%s): %v\n", f.Data, f.Nome, err)
			continue
		}

		fmt.Printf("[%-11s] %s - %s\n", f.Tipo, f.Data, f.Nome)
		inseridos++
	}

	fmt.Printf("\n%d/%d feriados de %d inseridos com sucesso.\n", inseridos, len(feriados), ano)
}

func buscarFeriados(apiKey, ibge string, ano int) ([]feriadoAPI, error) {
	url := fmt.Sprintf("https://feriadosapi.com/api/v1/feriados/cidade/%s?ano=%d&facultativos=true", ibge, ano)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)

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
