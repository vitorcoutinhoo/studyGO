package plantao

import "context"

type PlantaoRepository interface {
	Store(ctx context.Context, plantao *Plantao) error
	Update(ctx context.Context, plantao *Plantao) error
	Delete(ctx context.Context, plantaoId string) error
	FindById(ctx context.Context, plantaoId string) (*Plantao, error)
	Find(ctx context.Context, filter *Filter) ([]Plantao, error)
}

type Filter struct {
	ColaboradorID string
	Periodo       *Periodo
	Status        *StatusPlantao
	Limit         *int
	Offset        *int
}
