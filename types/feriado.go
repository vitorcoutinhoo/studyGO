package types

import "time"

type FeriadoRequest struct {
	Data      time.Time `json:"data" binding:"required"`
	Nome      string    `json:"nome" binding:"required"`
	Descricao *string   `json:"descricao"`
}

type FeriadoResponse struct {
	ID        string    `json:"id"`
	Data      time.Time `json:"data"`
	Nome      string    `json:"nome"`
	Descricao *string   `json:"descricao"`
	CreatedAt time.Time `json:"created_at"`
}
