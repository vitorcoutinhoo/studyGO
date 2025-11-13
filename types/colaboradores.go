package types

import (
	"time"

	"github.com/google/uuid"
)

type ColaboradorRequest struct {
	Nome         string
	Email        string
	Telefone     string
	Cargo        string
	Departamento string
	FotoURL      string
	DataAdmissao time.Time
}

type ColaboradorResponse struct {
	ID               uuid.UUID
	Nome             string
	Email            string
	Telefone         string
	Cargo            string
	Departamento     string
	FotoURL          string
	DataAdmissao     time.Time
	DataDesligamento time.Time
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type ColaboradorRepository interface {
	CreateColaborador(colaborador ColaboradorRequest) (*ColaboradorResponse, error)
	GetColaboradores() ([]*ColaboradorResponse, error)
	GetColaboradoresById(id uuid.UUID) (*ColaboradorResponse, error)
	UpdateColaborador(id uuid.UUID, colaborador ColaboradorRequest) error
	DeleteColaboradorById(id uuid.UUID) error
}
