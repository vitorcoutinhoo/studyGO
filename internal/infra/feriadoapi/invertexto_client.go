package feriadoapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"plantao/internal/domain/financeiro"
	"plantao/internal/infra/config"
)

const baseURL = "https://api.invertexto.com/v1/holidays"

type invertextoHoliday struct {
	Date  string `json:"date"`
	Name  string `json:"name"`
	Type  string `json:"type"`
	Level string `json:"level"`
}

type InvertextoClient struct {
	apiKey     string
	estado     string
	httpClient *http.Client
}

func NewInvertextoClient(cfg *config.Config) *InvertextoClient {
	return &InvertextoClient{
		apiKey: cfg.Feriados.APIKey,
		estado: cfg.Feriados.Estado,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *InvertextoClient) FetchFeriados(ano int) ([]financeiro.Feriado, error) {
	reqURL := fmt.Sprintf("%s/%d?token=%s&state=%s",
		baseURL, ano, url.QueryEscape(c.apiKey), url.QueryEscape(c.estado))

	resp, err := c.httpClient.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch feriados: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code from feriados API: %d", resp.StatusCode)
	}

	var holidays []invertextoHoliday
	if err := json.NewDecoder(resp.Body).Decode(&holidays); err != nil {
		return nil, fmt.Errorf("failed to decode feriados response: %w", err)
	}

	feriados := make([]financeiro.Feriado, 0, len(holidays))
	for _, h := range holidays {
		data, err := time.Parse("2006-01-02", h.Date)
		if err != nil {
			continue
		}
		feriados = append(feriados, financeiro.Feriado{
			Data:      data,
			Nome:      h.Name,
			Descricao: h.Level,
		})
	}

	return feriados, nil
}
