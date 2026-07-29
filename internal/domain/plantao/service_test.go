package plantao

import (
	"context"
	"errors"
	"testing"
	"time"

	"plantao/internal/domain/financeiro"
	"plantao/internal/domain/log"
	"plantao/internal/domain/shared"
)

type testLogger struct{}

func (testLogger) Info(string, ...any)  {}
func (testLogger) Warn(string, ...any)  {}
func (testLogger) Error(string, ...any) {}
func (testLogger) Debug(string, ...any) {}
func (testLogger) Fatal(string, ...any) {}
func (l testLogger) With(...any) log.Logger {
	return l
}
func (testLogger) Sync() error { return nil }

type repositoryFake struct {
	tx        *transactionFake
	committed bool
}

func (r *repositoryFake) Store(context.Context, *Plantao) error  { return nil }
func (r *repositoryFake) Update(context.Context, *Plantao) error { return nil }
func (r *repositoryFake) Delete(context.Context, string) error   { return nil }
func (r *repositoryFake) FindById(context.Context, string) (*Plantao, error) {
	return r.tx.plantao, nil
}
func (r *repositoryFake) Find(context.Context, *Filtro) ([]Plantao, error) { return nil, nil }
func (r *repositoryFake) WithTransaction(_ context.Context, fn func(PlantaoTransaction) error) error {
	if err := fn(r.tx); err != nil {
		return err
	}
	r.committed = true
	return nil
}

type transactionFake struct {
	plantao              *Plantao
	actor                *Actor
	colaboradoresExistem bool
	detalhesCount        int
	pagamentosCount      int
	feriados             map[string]bool
	valores              map[financeiro.TipoDia]int64
	pagamento            *Pagamento
	resumo               *DetalhesResumo
	detalhesInseridos    []PlantaoDetalhe
	statusAtualizado     *StatusPlantao
	historicos           int
	pagamentoCriado      bool
	pagamentoPago        bool
}

func (t *transactionFake) LockPlantao(context.Context, string) (*Plantao, error) {
	return t.plantao, nil
}
func (t *transactionFake) FindActor(context.Context, string) (*Actor, error) {
	return t.actor, nil
}
func (t *transactionFake) ColaboradorExists(context.Context, string) (bool, error) {
	return t.colaboradoresExistem, nil
}
func (t *transactionFake) CountDetalhes(context.Context, string) (int, error) {
	return t.detalhesCount, nil
}
func (t *transactionFake) CountPagamentos(context.Context, string) (int, error) {
	return t.pagamentosCount, nil
}
func (t *transactionFake) FindFeriados(context.Context, time.Time, time.Time) (map[string]bool, error) {
	return t.feriados, nil
}
func (t *transactionFake) FindValorDiaCentavos(_ context.Context, tipo financeiro.TipoDia, _ time.Time) (int64, error) {
	valor, ok := t.valores[tipo]
	if !ok {
		return 0, financeiro.ErrorValorDiaNotFound
	}
	return valor, nil
}
func (t *transactionFake) InsertDetalhes(_ context.Context, _ string, detalhes []PlantaoDetalhe) error {
	t.detalhesInseridos = detalhes
	return nil
}
func (t *transactionFake) ClosePlantao(context.Context, string, int64, *string) error {
	status := StatusPlantaoConcluido
	t.statusAtualizado = &status
	return nil
}
func (t *transactionFake) UpdateStatus(_ context.Context, _ string, status StatusPlantao) error {
	t.statusAtualizado = &status
	return nil
}
func (t *transactionFake) InsertHistorico(context.Context, string, StatusPlantao, StatusPlantao, string, *string) error {
	t.historicos++
	return nil
}
func (t *transactionFake) InsertPagamentoPendente(context.Context, string, string, int64) error {
	t.pagamentoCriado = true
	return nil
}
func (t *transactionFake) FindPagamento(context.Context, string) (*Pagamento, error) {
	if t.pagamento == nil {
		return nil, ErrorPagamentoNotFound
	}
	return t.pagamento, nil
}
func (t *transactionFake) SummarizeDetalhes(context.Context, string) (*DetalhesResumo, error) {
	return t.resumo, nil
}
func (t *transactionFake) PayPagamento(context.Context, string, time.Time, *string) error {
	t.pagamentoPago = true
	return nil
}

func plantaoDeTeste(status StatusPlantao) *Plantao {
	loc, _ := time.LoadLocation("America/Sao_Paulo")
	data := time.Date(2026, 7, 6, 0, 0, 0, 0, loc)
	return &Plantao{
		Id:            "plantao-1",
		ColaboradorId: "colaborador-1",
		Periodo:       &shared.Periodo{Inicio: data, Fim: data},
		Status:        status,
	}
}

func novoServicoFake(status StatusPlantao, role, actorColaborador string) (*PlantaoService, *repositoryFake) {
	tx := &transactionFake{
		plantao:              plantaoDeTeste(status),
		actor:                &Actor{UsuarioID: "usuario-1", ColaboradorID: actorColaborador, Role: role},
		colaboradoresExistem: true,
		valores: map[financeiro.TipoDia]int64{
			financeiro.TipoDiaUtil:    10000,
			financeiro.TipoDiaSabado:  20000,
			financeiro.TipoDiaDomingo: 30000,
			financeiro.TipoDiaFeriado: 40000,
		},
	}
	repo := &repositoryFake{tx: tx}
	return NewPlantaoService(repo, financeiro.NewCalculoService(testLogger{}), testLogger{}), repo
}

func TestFecharPlantaoAutorizacaoEAtomicidade(t *testing.T) {
	tests := []struct {
		nome             string
		role             string
		actorColaborador string
		erroEsperado     error
	}{
		{"proprio colaborador", roleColaborador, "colaborador-1", nil},
		{"outro colaborador", roleColaborador, "colaborador-2", ErrorUsuarioSemPermissao},
		{"gerente", roleGerente, "colaborador-2", nil},
		{"admin", roleAdmin, "colaborador-2", nil},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			service, repo := novoServicoFake(StatusPlantaoEmAndamento, tt.role, tt.actorColaborador)
			resultado, err := service.FecharPlantao(context.Background(), "plantao-1", "usuario-1", nil)
			if !errors.Is(err, tt.erroEsperado) {
				t.Fatalf("erro = %v, esperado %v", err, tt.erroEsperado)
			}
			if tt.erroEsperado != nil {
				if repo.committed || repo.tx.pagamentoCriado || len(repo.tx.detalhesInseridos) != 0 {
					t.Fatal("operação inválida deixou efeitos persistentes")
				}
				return
			}
			if resultado.Status != StatusPlantaoConcluido || resultado.ValorTotal != 100 {
				t.Fatalf("resultado inesperado: %+v", resultado)
			}
			if !repo.committed || !repo.tx.pagamentoCriado || repo.tx.historicos != 1 || len(repo.tx.detalhesInseridos) != 1 {
				t.Fatal("fechamento não executou todas as etapas")
			}
		})
	}
}

func TestFecharPlantaoRejeitaRepeticaoEDadosExistentes(t *testing.T) {
	service, _ := novoServicoFake(StatusPlantaoConcluido, roleAdmin, "colaborador-2")
	if _, err := service.FecharPlantao(context.Background(), "plantao-1", "usuario-1", nil); !errors.Is(err, ErrorPlantaoJaFechado) {
		t.Fatalf("erro = %v, esperado plantão já fechado", err)
	}

	service, repo := novoServicoFake(StatusPlantaoEmAndamento, roleAdmin, "colaborador-2")
	repo.tx.detalhesCount = 1
	if _, err := service.FecharPlantao(context.Background(), "plantao-1", "usuario-1", nil); !errors.Is(err, ErrorDetalhesExistentes) {
		t.Fatalf("erro = %v, esperado detalhes existentes", err)
	}
}

func TestPagarPlantao(t *testing.T) {
	service, repo := novoServicoFake(StatusPlantaoConcluido, roleGerente, "colaborador-2")
	repo.tx.plantao.ValorTotal = 100
	repo.tx.pagamento = &Pagamento{
		ID: "pagamento-1", PlantaoID: "plantao-1", ColaboradorID: "colaborador-1",
		ValorTotalCentavos: 10000, Status: statusPagamentoPendente,
	}
	repo.tx.resumo = &DetalhesResumo{Quantidade: 1, DatasDistintas: 1, ValorTotalCentavos: 10000}

	if err := service.PagarPlantao(context.Background(), "plantao-1", "usuario-1", nil); err != nil {
		t.Fatal(err)
	}
	if !repo.committed || !repo.tx.pagamentoPago || repo.tx.statusAtualizado == nil ||
		*repo.tx.statusAtualizado != StatusPlantaoPago || repo.tx.historicos != 1 {
		t.Fatal("pagamento não executou todas as etapas")
	}
}

func TestPagarPlantaoRejeitaColaboradorEInconsistencia(t *testing.T) {
	service, repo := novoServicoFake(StatusPlantaoConcluido, roleColaborador, "colaborador-1")
	repo.tx.plantao.ValorTotal = 100
	if err := service.PagarPlantao(context.Background(), "plantao-1", "usuario-1", nil); !errors.Is(err, ErrorUsuarioSemPermissao) {
		t.Fatalf("erro = %v, esperado sem permissão", err)
	}

	service, repo = novoServicoFake(StatusPlantaoConcluido, roleAdmin, "colaborador-2")
	repo.tx.plantao.ValorTotal = 100
	repo.tx.pagamento = &Pagamento{
		ID: "pagamento-1", ColaboradorID: "colaborador-1",
		ValorTotalCentavos: 9999, Status: statusPagamentoPendente,
	}
	repo.tx.resumo = &DetalhesResumo{Quantidade: 1, DatasDistintas: 1, ValorTotalCentavos: 10000}
	if err := service.PagarPlantao(context.Background(), "plantao-1", "usuario-1", nil); !errors.Is(err, ErrorPagamentoInconsistente) {
		t.Fatalf("erro = %v, esperado pagamento inconsistente", err)
	}
	if repo.committed || repo.tx.pagamentoPago {
		t.Fatal("pagamento inconsistente deixou efeitos persistentes")
	}
}
