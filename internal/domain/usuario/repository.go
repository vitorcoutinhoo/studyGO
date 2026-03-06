package usuario

import (
	"context"

	"github.com/google/uuid"
)

type UsuarioRepository interface {
	Store(ctx context.Context, usuario *Usuario) (*Usuario, error)
	Update(ctx context.Context, u *Usuario) error
	Disable(ctx context.Context, usuarioId uuid.UUID) error
	FindById(ctx context.Context, usuarioId uuid.UUID) (*Usuario, error)
	FindByEmail(ctx context.Context, email string) (*Usuario, error)
	ExistsEmail(ctx context.Context, email string) (bool, error)
	ExistsId(ctx context.Context, usuarioId uuid.UUID) (bool, error)
	ExistsEmailExcludingId(ctx context.Context, email string, usuarioId uuid.UUID) (bool, error)
}
