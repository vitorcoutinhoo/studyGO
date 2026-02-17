package colaborador

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Interface do repositório de colaboradores, definindo os métodos necessários para manipulação dos dados.
type ColaboradorRepository interface {
	Store(ctx context.Context, colaborador *Colaborador) (*Colaborador, error)
	Update(ctx context.Context, colaborador *Colaborador) error
	Disable(ctx context.Context, colaboradorId uuid.UUID) error
	FindById(ctx context.Context, colaboradorId uuid.UUID) (*Colaborador, error)
	FindByEmail(ctx context.Context, email string) (*Colaborador, error)
	FindByFilter(ctx context.Context, filter ColaboradorFilter) ([]Colaborador, error)
	ExistsEmail(ctx context.Context, email string) bool
	ExistsId(ctx context.Context, colaboradorId uuid.UUID) bool
	ExistsEmailExcludingId(ctx context.Context, email string, colaboradorId uuid.UUID) bool
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
