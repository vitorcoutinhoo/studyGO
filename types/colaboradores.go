package types

import (
	"time"

	"github.com/google/uuid"
)

type ColaboradorRequest struct {
	Nome             string     `json:"nome"`
	Email            string     `json:"email"`
	Telefone         *string    `json:"telefone"`
	Cargo            *string    `json:"cargo"`
	Departamento     *string    `json:"departamento"`
	FotoURL          *string    `json:"foto_url"`
	DataAdmissao     *time.Time `json:"data_admissao"`
	DataDesligamento *time.Time `json:"data_desligamento"`
}

type ColaboradorResponse struct {
	ID               uuid.UUID  `json:"id"`
	Nome             string     `json:"nome"`
	Email            string     `json:"email"`
	Telefone         *string    `json:"telefone"`
	Cargo            *string    `json:"cargo"`
	Departamento     *string    `json:"departamento"`
	FotoURL          *string    `json:"foto_url"`
	Ativo            string     `json:"ativo"`
	DataAdmissao     *time.Time `json:"data_admissao"`
	DataDesligamento *time.Time `json:"data_desligamento"`
	CreatedAt        *time.Time `json:"created_at"`
	UpdatedAt        *time.Time `json:"updated_at"`
}

type ColaboradorRepository interface {
	CreateColaborador(colaborador ColaboradorRequest) (*ColaboradorResponse, error)
	GetColaboradores() ([]*ColaboradorResponse, error)
	GetColaboradoresById(id uuid.UUID) (*ColaboradorResponse, error)
	UpdateColaborador(id uuid.UUID, colaborador ColaboradorRequest) error
	DeleteColaboradorById(id uuid.UUID) error
}
