package plantao

import (
	"context"

	"plantao/internal/domain/shared"
)

type PlantaoRepository interface {
	Store(ctx context.Context, plantao *Plantao) error
	Update(ctx context.Context, plantao *Plantao) error
	Delete(ctx context.Context, plantaoId string) error
	FindById(ctx context.Context, plantaoId string) (*Plantao, error)
	Find(ctx context.Context, filter *Filtro) ([]Plantao, error)
}

type Filtro struct {
	ColaboradorID string
	Periodo       *shared.Periodo
	Status        *StatusPlantao
	Limit         *int
	Offset        *int
}
