package comunicacao

import (
	"context"

	"github.com/google/uuid"
)

type ModeloComunicaRepository interface {
	Store(ctx context.Context, com *Comunicacao) (*Comunicacao, error)
	Update(ctx context.Context, com *Comunicacao) error
	Disable(ctx context.Context, modeloId uuid.UUID) error
	FindById(ctx context.Context, modeloId uuid.UUID) (*Comunicacao, error)
	FindByTipo(ctx context.Context, tipoComunicacao string) (*Comunicacao, error)
	FindAll(ctx context.Context) ([]*Comunicacao, error)
	ExistsName(ctx context.Context, nome string) (bool, error)
	ExistsNameExcludingId(ctx context.Context, nome string, id uuid.UUID) (bool, error)
}

type EnvioComunicacaoRepository interface {
	Store(ctx context.Context, com *Envio) error
}
