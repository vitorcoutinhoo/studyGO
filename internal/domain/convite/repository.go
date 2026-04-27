package convite

import (
	"context"

	"github.com/google/uuid"
)

type ConviteRepository interface {
	Store(ctx context.Context, idColaborador uuid.UUID) (*Convite, error)
	FindByToken(ctx context.Context, token uuid.UUID) (*Convite, error)
	MarkAsUsed(ctx context.Context, token uuid.UUID) error
	Disable(ctx context.Context, token uuid.UUID) error
	FindAll(ctx context.Context) ([]Convite, error)
}
