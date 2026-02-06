package colaborador

import (
	"context"
	"time"
)

// Interface do repositório de colaboradores, definindo os métodos necessários para manipulação dos dados.
type ColaboradorRepository interface {
	Store(ctx context.Context, colaborador *Colaborador) error
	Update(ctx context.Context, colaborador *Colaborador) error
	Disable(ctx context.Context, colaboradorId string) error
	FindById(ctx context.Context, colaboradorId string) (*Colaborador, error)
	FindByEmail(ctx context.Context, email string) (*Colaborador, error)
	FindByFilter(ctx context.Context, filter ColaboradorFilter) ([]Colaborador, error)
	ExistsEmail(ctx context.Context, email string) bool
	ExistsId(ctx context.Context, colaboradorId string) bool
	ExistsEmailExcludingId(ctx context.Context, email, colaboradorId string) bool
}

// Filtro para busca de colaboradores, permitindo filtrar por nome, email, telefone, cargo, departamento e data de admissão.
type ColaboradorFilter struct {
	Nome         *string
	Email        *string
	Telefone     *string
	Cargo        *string
	Departamento *string
	DataAdmissao *time.Time
}
