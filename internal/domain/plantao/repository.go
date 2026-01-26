package plantao

import "context"

type PlantaoRepository interface {
	Store(ctx context.Context, plantao *Plantao) error
	Update(ctx context.Context, plantao *Plantao) error
	Delete(ctx context.Context, plantaoId string) error
	List(ctx context.Context) ([]Plantao, error)
	FindById(ctx context.Context, plantaoId string) (*Plantao, error)
	FindByColaboradorId(ctx context.Context, colaboradorId string) ([]Plantao, error)
	FindByPeriodo(ctx context.Context, periodo Periodo) ([]Plantao, error)
	FindByStatus(ctx context.Context, status StatusPlantao) ([]Plantao, error)
}

type Filter struct {
	ColaboradorId string
	Periodo       *Periodo
	Status        *StatusPlantao
	Limit         *int
	Offset        *int
}
