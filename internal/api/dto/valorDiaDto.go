package dto

import "time"

type SetValorDiaRequest struct {
	TipoDia        string  `json:"tipo_dia" binding:"required"`
	Valor          float64 `json:"valor" binding:"required"`
	VigenciaInicio string  `json:"vigencia_inicio" binding:"required"`
}

type ValorDiaResponse struct {
	Id             string     `json:"id"`
	TipoDia        string     `json:"tipo_dia"`
	Valor          float64    `json:"valor"`
	VigenciaInicio time.Time  `json:"vigencia_inicio"`
	VigenciaFim    *time.Time `json:"vigencia_fim"`
}
