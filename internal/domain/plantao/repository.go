package plantao

import (
	"context"
	"time"

	"plantao/internal/domain/financeiro"
	"plantao/internal/domain/shared"
)

type PlantaoRepository interface {
	Store(ctx context.Context, plantao *Plantao) error
	Update(ctx context.Context, plantao *Plantao) error
	Delete(ctx context.Context, plantaoID string) error
	FindById(ctx context.Context, plantaoID string) (*Plantao, error)
	Find(ctx context.Context, filter *Filtro) ([]Plantao, error)
	FindRelatorio(ctx context.Context, filter *RelatorioFiltro) ([]RelatorioItem, error)
	WithTransaction(ctx context.Context, fn func(PlantaoTransaction) error) error
}

type PlantaoTransaction interface {
	financeiro.CalculoFonte

	LockPlantao(ctx context.Context, plantaoID string) (*Plantao, error)
	LockColaborador(ctx context.Context, colaboradorID string) error
	HasOverlappingPlantao(ctx context.Context, colaboradorID string, inicio, fim time.Time, excludePlantaoID string) (bool, error)
	UpdateSchedule(ctx context.Context, plantaoID, colaboradorID string, inicio, fim time.Time) (*Plantao, error)
	FindActor(ctx context.Context, usuarioID string) (*Actor, error)
	ColaboradorExists(ctx context.Context, colaboradorID string) (bool, error)
	CountDetalhes(ctx context.Context, plantaoID string) (int, error)
	CountPagamentos(ctx context.Context, plantaoID string) (int, error)
	InsertDetalhes(ctx context.Context, plantaoID string, detalhes []PlantaoDetalhe) error
	ClosePlantao(ctx context.Context, plantaoID string, valorTotalCentavos int64, observacoes *string) error
	UpdateStatus(ctx context.Context, plantaoID string, status StatusPlantao) error
	InsertHistorico(ctx context.Context, plantaoID string, statusAntigo, statusNovo StatusPlantao, usuarioID string, observacoes *string) error
	InsertPagamentoPendente(ctx context.Context, plantaoID, colaboradorID string, valorTotalCentavos int64) error
	FindPagamento(ctx context.Context, plantaoID string) (*Pagamento, error)
	SummarizeDetalhes(ctx context.Context, plantaoID string) (*DetalhesResumo, error)
	PayPagamento(ctx context.Context, pagamentoID string, dataPagamento time.Time, observacoes *string) error
	LockPlantoesAgendadosAte(ctx context.Context, instante time.Time, limite int) ([]*Plantao, error)
	StartPlantao(ctx context.Context, plantaoID string) error
	InsertHistoricoAutomatico(ctx context.Context, plantaoID string, statusAntigo, statusNovo StatusPlantao, observacoes string) error
}

type Actor struct {
	UsuarioID     string
	ColaboradorID string
	Role          string
}

type PlantaoDetalhe struct {
	Data          time.Time
	TipoDia       string
	ValorCentavos int64
}

type Pagamento struct {
	ID                 string
	PlantaoID          string
	ColaboradorID      string
	ValorTotalCentavos int64
	Status             string
}

type DetalhesResumo struct {
	Quantidade         int
	DatasDistintas     int
	ValorTotalCentavos int64
}

type Filtro struct {
	ColaboradorID string
	Periodo       *shared.Periodo
	Status        *StatusPlantao
	Limit         *int
	Offset        *int
}
