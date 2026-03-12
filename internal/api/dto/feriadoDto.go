package dto

import "time"

type UpdateDataFeriadoRequest struct {
	NovaData string `json:"nova_data" binding:"required"`
}

type FeriadoResponse struct {
	Id        string    `json:"id"`
	Data      time.Time `json:"data"`
	Nome      string    `json:"nome"`
	Descricao string    `json:"descricao"`
}
