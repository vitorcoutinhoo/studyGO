package dto

import (
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

type CreatePlantaoResponse struct {
	Id            string                `json:"id"`
	ColaboradorId string                `json:"colaborador_id"`
	Periodo       shared.Periodo        `json:"periodo"`
	Status        plantao.StatusPlantao `json:"status"`
	ValorTotal    float64               `json:"valor_total"`
	Observacoes   *string               `json:"observacoes,omitempty"`
}