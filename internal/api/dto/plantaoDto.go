package dto

import (
	"time"

	"plantao/internal/domain/plantao"
	"plantao/internal/domain/shared"
)

type PeriodoRequest struct {
	Inicio string `json:"inicio" binding:"required"`
	Fim    string `json:"fim" binding:"required"`
}

type CreatePlantaoRequest struct {
	Periodo       PeriodoRequest `json:"periodo" binding:"required"`
	ColaboradorId string         `json:"colaborador_id" binding:"required"`
}

type UpdateStatusPlantaoRequest struct {
	NewStatus   string  `json:"new_status" binding:"required"`
	Observacoes *string `json:"observacoes"`
}

type PagamentoPlantaoRequest struct {
	Observacoes *string `json:"observacoes"`
}

type CreatePlantaoResponse struct {
	Id            string                `json:"id"`
	ColaboradorId string                `json:"colaborador_id"`
	Periodo       shared.Periodo        `json:"periodo"`
	Status        plantao.StatusPlantao `json:"status"`
	ValorTotal    float64               `json:"valor_total"`
	Observacoes   *string               `json:"observacoes,omitempty"`
}

type RelatorioPlantaoResponse struct {
	ColaboradorID   string                `json:"colaborador_id"`
	PlantaoID       string                `json:"plantao_id"`
	Status          plantao.StatusPlantao `json:"status"`
	DataInicio      time.Time             `json:"data_inicio"`
	DataFim         time.Time             `json:"data_fim"`
	Data            string                `json:"data"`
	NomeColaborador string                `json:"nome_colaborador"`
	ValorTotal      float64               `json:"valor_total"`
	Valor           float64               `json:"valor"`
	Observacoes     *string               `json:"observacoes"`
}
