package colaborador

import "context"

type ColaboradorRepository interface {
	Store(ctx context.Context, colaborador *Colaborador) error
	Update(ctx context.Context, colaborador *Colaborador) error
	Delete(ctx context.Context, colaboradorId string) error
	FindById(ctx context.Context, colaboradorId string) (*Colaborador, error)
	FindByEmail(ctx context.Context, email string) (*Colaborador, error)
	Find(ctx context.Context) ([]Colaborador, error)
	ExistsEmail(ctx context.Context, email string) bool
	ExistsId(ctx context.Context, colaboradorId string) bool
	ExistsEmailExcludingId(ctx context.Context, email, colaboradorId string) bool
}
