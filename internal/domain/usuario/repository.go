package usuario

import (
	"context"

	"github.com/google/uuid"
)

type UsuarioRepostory interface {
	Store(ctx context.Context, usuario *Usuario) (*Usuario, error)
	Update(ctx context.Context, usuario *Usuario) error
	Disable(ctx context.Context, usuarioId uuid.UUID) error
	FindById(ctx context.Context, usuarioId uuid.UUID) (*Usuario, error)
	ExistsEmail(ctx context.Context, email string) bool
	ExistsId(ctx context.Context, usuarioId uuid.UUID) bool
	ExistsEmailExcludingId(ctx context.Context, email string, usuarioId uuid.UUID) bool
}
