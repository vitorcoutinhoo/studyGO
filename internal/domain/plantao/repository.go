package plantao

import (
	"context"
	"time"

	"plantao/internal/domain/shared"
)

type PlantaoRepository interface {
	Store(ctx context.Context, plantao *Plantao) error
	Update(ctx context.Context, plantao *Plantao) error
	Delete(ctx context.Context, plantaoId string) error
	FindById(ctx context.Context, plantaoId string) (*Plantao, error)
	Find(ctx context.Context, filter *Filtro) ([]Plantao, error)
	StoreDetalhesAndUpdateValorTotal(ctx context.Context, plantaoId string, valorTotal float64, observacoes *string, detalhes []PlantaoDetalhe) error
}

type PlantaoDetalhe struct {
	IdPlantao string
	Data      time.Time
	TipoDia   string
	Valor     float64
}

type Filtro struct {
	ColaboradorID string
	Periodo       *shared.Periodo
	Status        *StatusPlantao
	Limit         *int
	Offset        *int
}
