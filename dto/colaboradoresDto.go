package dto

import "time"

type ColabotadoresRequestDTO struct {
	Nome             string     `json:"nome" binding:"required"`
	Email            string     `json:"email" binding:"required,email"`
	Telefone         *string    `json:"telefone"`
	Cargo            *string    `json:"cargo"`
	Departamento     *string    `json:"departamento" `
	FotoURL          *string    `json:"foto_url"`
	DataAdmissao     *time.Time `json:"data_admissao"`
	DataDesligamento *time.Time `json:"data_desligamento"`
}

type ColabotadoresResponseDTO struct {
	ID               string
	Nome             string
	Email            string
	Telefone         *string
	Cargo            *string
	Departamento     *string
	FotoURL          *string
	DataAdmissao     *time.Time
	DataDesligamento *time.Time
}
