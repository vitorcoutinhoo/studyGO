package plantao

import (
	"context"
	"errors"
	"fmt"
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
	tx              *transactionFake
	committed       bool
	stored          *Plantao
	storeErr        error
	relatorio       []RelatorioItem
	relatorioFiltro *RelatorioFiltro
	relatorioErr    error
}

func (r *repositoryFake) Store(_ context.Context, p *Plantao) error {
	if r.storeErr != nil {
		return r.storeErr
	}
	r.stored = p
	return nil
}
func (r *repositoryFake) Update(context.Context, *Plantao) error { return nil }
func (r *repositoryFake) Delete(context.Context, string) error   { return nil }
func (r *repositoryFake) FindById(context.Context, string) (*Plantao, error) {
	return r.tx.plantao, nil
}
func (r *repositoryFake) Find(context.Context, *Filtro) ([]Plantao, error) { return nil, nil }
func (r *repositoryFake) FindRelatorio(_ context.Context, filtro *RelatorioFiltro) ([]RelatorioItem, error) {
	r.relatorioFiltro = filtro
	return r.relatorio, r.relatorioErr
}
func (r *repositoryFake) WithTransaction(_ context.Context, fn func(PlantaoTransaction) error) error {
	if err := fn(r.tx); err != nil {
		return err
	}
	r.committed = true
	return nil
}

type transactionFake struct {
	plantao                 *Plantao
	actor                   *Actor
	colaboradoresExistem    bool
	detalhesCount           int
	pagamentosCount         int
	feriados                map[string]bool
	valores                 map[financeiro.TipoDia]int64
	pagamento               *Pagamento
	resumo                  *DetalhesResumo
	detalhesInseridos       []PlantaoDetalhe
	statusAtualizado        *StatusPlantao
	historicos              int
	pagamentoCriado         bool
	pagamentoPago           bool
	plantoesAgendados       []*Plantao
	instanteRecebido        time.Time
	consultasAgendados      int
	plantaoIniciado         []string
	historicosAutomaticos   int
	historicoAutomaticoErr  error
	colaboradoresBloqueados []string
	colaboradoresAusentes   map[string]bool
	sobreposto              bool
	sobreposicaoErr         error
	excludePlantaoID        string
	agendaAtualizada        bool
}

func (t *transactionFake) LockPlantao(context.Context, string) (*Plantao, error) {
	return t.plantao, nil
}
func (t *transactionFake) LockColaborador(_ context.Context, colaboradorID string) error {
	if t.colaboradoresAusentes[colaboradorID] {
		return ErrorColaboradorNotFound
	}
	t.colaboradoresBloqueados = append(t.colaboradoresBloqueados, colaboradorID)
	return nil
}
func (t *transactionFake) HasOverlappingPlantao(_ context.Context, _ string, _, _ time.Time, excludePlantaoID string) (bool, error) {
	t.excludePlantaoID = excludePlantaoID
	return t.sobreposto, t.sobreposicaoErr
}
func (t *transactionFake) UpdateSchedule(_ context.Context, _ string, colaboradorID string, inicio, fim time.Time) (*Plantao, error) {
	copia := *t.plantao
	copia.Periodo = &shared.Periodo{Inicio: inicio, Fim: fim}
	copia.ColaboradorId = colaboradorID
	copia.UpdatedAt = time.Now()
	t.agendaAtualizada = true
	return &copia, nil
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
func (t *transactionFake) FindValorDiaCentavos(_ context.Context, tipo financeiro.TipoDia) (int64, error) {
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
func (t *transactionFake) LockPlantoesAgendadosAte(_ context.Context, instante time.Time, limite int) ([]*Plantao, error) {
	t.instanteRecebido = instante
	t.consultasAgendados++
	quantidade := len(t.plantoesAgendados)
	if quantidade > limite {
		quantidade = limite
	}
	resultado := t.plantoesAgendados[:quantidade]
	t.plantoesAgendados = t.plantoesAgendados[quantidade:]
	return resultado, nil
}
func (t *transactionFake) StartPlantao(_ context.Context, plantaoID string) error {
	t.plantaoIniciado = append(t.plantaoIniciado, plantaoID)
	return nil
}
func (t *transactionFake) InsertHistoricoAutomatico(context.Context, string, StatusPlantao, StatusPlantao, string) error {
	if t.historicoAutomaticoErr != nil {
		return t.historicoAutomaticoErr
	}
	t.historicosAutomaticos++
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

func TestCreatePlantaoDelegaCriacaoAtomicaAoRepositorio(t *testing.T) {
	repo := &repositoryFake{}
	service := NewPlantaoService(repo, nil, testLogger{})
	periodo := &shared.Periodo{
		Inicio: time.Date(2026, 9, 1, 8, 0, 0, 0, time.UTC),
		Fim:    time.Date(2026, 9, 1, 18, 0, 0, 0, time.UTC),
	}

	criado, err := service.CreatePlantao(context.Background(), "colaborador-1", periodo)
	if err != nil {
		t.Fatal(err)
	}
	if repo.stored == nil || criado != repo.stored || criado.ColaboradorId != "colaborador-1" {
		t.Fatalf("plantão não delegado ao repositório: criado=%+v armazenado=%+v", criado, repo.stored)
	}
}

func TestCreatePlantaoPropagaConflitoSemCriacao(t *testing.T) {
	repo := &repositoryFake{storeErr: ErrorExistingPlantao}
	service := NewPlantaoService(repo, nil, testLogger{})
	periodo := &shared.Periodo{
		Inicio: time.Date(2026, 9, 1, 8, 0, 0, 0, time.UTC),
		Fim:    time.Date(2026, 9, 1, 18, 0, 0, 0, time.UTC),
	}

	criado, err := service.CreatePlantao(context.Background(), "colaborador-1", periodo)
	if !errors.Is(err, ErrorExistingPlantao) || criado != nil || repo.stored != nil {
		t.Fatalf("resultado inesperado: criado=%+v armazenado=%+v erro=%v", criado, repo.stored, err)
	}
}

func TestEditarPlantaoAtualizaCamposParcialmente(t *testing.T) {
	tests := []struct {
		nome             string
		atualizacao      func() *AtualizacaoPlantao
		colaboradorFinal string
		inicioFinal      int
		fimFinal         int
	}{
		{
			nome: "somente colaborador",
			atualizacao: func() *AtualizacaoPlantao {
				id := "colaborador-2"
				return &AtualizacaoPlantao{ColaboradorID: &id}
			},
			colaboradorFinal: "colaborador-2", inicioFinal: 8, fimFinal: 18,
		},
		{
			nome: "somente início",
			atualizacao: func() *AtualizacaoPlantao {
				inicio := time.Date(2026, 9, 1, 9, 0, 0, 0, time.UTC)
				return &AtualizacaoPlantao{DataInicio: &inicio}
			},
			colaboradorFinal: "colaborador-1", inicioFinal: 9, fimFinal: 18,
		},
		{
			nome: "somente fim",
			atualizacao: func() *AtualizacaoPlantao {
				fim := time.Date(2026, 9, 1, 20, 0, 0, 0, time.UTC)
				return &AtualizacaoPlantao{DataFim: &fim}
			},
			colaboradorFinal: "colaborador-1", inicioFinal: 8, fimFinal: 20,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			service, repo := novoServicoFake(StatusPlantaoAgendado, roleAdmin, "ator")
			repo.tx.plantao.Periodo = &shared.Periodo{
				Inicio: time.Date(2026, 9, 1, 8, 0, 0, 0, time.UTC),
				Fim:    time.Date(2026, 9, 1, 18, 0, 0, 0, time.UTC),
			}

			resultado, err := service.EditarPlantao(context.Background(), repo.tx.plantao.Id, "usuario-1", tt.atualizacao())
			if err != nil {
				t.Fatal(err)
			}
			if !repo.committed || !repo.tx.agendaAtualizada || resultado.ColaboradorId != tt.colaboradorFinal ||
				resultado.Periodo.Inicio.Hour() != tt.inicioFinal || resultado.Periodo.Fim.Hour() != tt.fimFinal ||
				repo.tx.excludePlantaoID != repo.tx.plantao.Id {
				t.Fatalf("edição parcial incorreta: resultado=%+v tx=%+v", resultado, repo.tx)
			}
		})
	}
}

func TestEditarPlantaoAutorizaAdminEGerente(t *testing.T) {
	for _, tt := range []struct {
		role     string
		esperado error
	}{
		{roleAdmin, nil},
		{roleGerente, nil},
		{roleColaborador, ErrorUsuarioSemPermissao},
	} {
		t.Run(tt.role, func(t *testing.T) {
			service, repo := novoServicoFake(StatusPlantaoAgendado, tt.role, "colaborador-1")
			fim := repo.tx.plantao.Periodo.Fim.Add(time.Hour)
			_, err := service.EditarPlantao(context.Background(), repo.tx.plantao.Id, "usuario-1", &AtualizacaoPlantao{DataFim: &fim})
			if !errors.Is(err, tt.esperado) {
				t.Fatalf("erro = %v, esperado %v", err, tt.esperado)
			}
			if tt.esperado != nil && (repo.committed || repo.tx.agendaAtualizada) {
				t.Fatal("edição não autorizada foi confirmada")
			}
		})
	}
}

func TestEditarPlantaoValidaEstadoPeriodoSobreposicaoEColaborador(t *testing.T) {
	tests := []struct {
		nome     string
		preparar func(*repositoryFake) *AtualizacaoPlantao
		esperado error
	}{
		{
			nome:     "patch vazio",
			preparar: func(*repositoryFake) *AtualizacaoPlantao { return &AtualizacaoPlantao{} },
			esperado: ErrorAtualizacaoPlantaoVazia,
		},
		{
			nome: "status não editável",
			preparar: func(repo *repositoryFake) *AtualizacaoPlantao {
				repo.tx.plantao.Status = StatusPlantaoEmAndamento
				fim := repo.tx.plantao.Periodo.Fim.Add(time.Hour)
				return &AtualizacaoPlantao{DataFim: &fim}
			},
			esperado: ErrorPlantaoNaoEditavel,
		},
		{
			nome: "período invertido após mesclagem",
			preparar: func(repo *repositoryFake) *AtualizacaoPlantao {
				inicio := repo.tx.plantao.Periodo.Fim.Add(time.Hour)
				return &AtualizacaoPlantao{DataInicio: &inicio}
			},
			esperado: shared.ErrorEndBeforeStart,
		},
		{
			nome: "sobreposição",
			preparar: func(repo *repositoryFake) *AtualizacaoPlantao {
				repo.tx.sobreposto = true
				fim := repo.tx.plantao.Periodo.Fim.Add(time.Hour)
				return &AtualizacaoPlantao{DataFim: &fim}
			},
			esperado: ErrorExistingPlantao,
		},
		{
			nome: "colaborador inexistente",
			preparar: func(repo *repositoryFake) *AtualizacaoPlantao {
				id := "colaborador-2"
				repo.tx.colaboradoresAusentes = map[string]bool{id: true}
				return &AtualizacaoPlantao{ColaboradorID: &id}
			},
			esperado: ErrorColaboradorNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			service, repo := novoServicoFake(StatusPlantaoAgendado, roleAdmin, "ator")
			_, err := service.EditarPlantao(context.Background(), repo.tx.plantao.Id, "usuario-1", tt.preparar(repo))
			if !errors.Is(err, tt.esperado) {
				t.Fatalf("erro = %v, esperado %v", err, tt.esperado)
			}
			if repo.committed || repo.tx.agendaAtualizada {
				t.Fatal("edição inválida foi confirmada")
			}
		})
	}
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

func TestIniciarPlantoesAgendadosAtualizaStatusEHistorico(t *testing.T) {
	service, repo := novoServicoFake(StatusPlantaoAgendado, roleAdmin, "colaborador-1")
	instante := time.Date(2026, 7, 31, 12, 0, 0, 0, time.UTC)
	service.now = func() time.Time { return instante }
	primeiro := plantaoDeTeste(StatusPlantaoAgendado)
	primeiro.Id = "plantao-automatico-1"
	segundo := plantaoDeTeste(StatusPlantaoAgendado)
	segundo.Id = "plantao-automatico-2"
	repo.tx.plantoesAgendados = []*Plantao{primeiro, segundo}

	quantidade, err := service.IniciarPlantoesAgendados(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if quantidade != 2 || !repo.committed || repo.tx.historicosAutomaticos != 2 {
		t.Fatalf("quantidade/commit/históricos = %d/%v/%d", quantidade, repo.committed, repo.tx.historicosAutomaticos)
	}
	if len(repo.tx.plantaoIniciado) != 2 || primeiro.Status != StatusPlantaoEmAndamento || segundo.Status != StatusPlantaoEmAndamento {
		t.Fatalf("plantões não foram iniciados: %+v", repo.tx.plantaoIniciado)
	}
	if !repo.tx.instanteRecebido.Equal(instante) {
		t.Fatalf("instante recebido = %s, esperado %s", repo.tx.instanteRecebido, instante)
	}
}

func TestIniciarPlantoesAgendadosProcessaLotesSucessivos(t *testing.T) {
	service, repo := novoServicoFake(StatusPlantaoAgendado, roleAdmin, "colaborador-1")
	repo.tx.plantoesAgendados = make([]*Plantao, loteInicioAutomatico+1)
	for i := range repo.tx.plantoesAgendados {
		p := plantaoDeTeste(StatusPlantaoAgendado)
		p.Id = fmt.Sprintf("plantao-automatico-%d", i)
		repo.tx.plantoesAgendados[i] = p
	}

	quantidade, err := service.IniciarPlantoesAgendados(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if quantidade != loteInicioAutomatico+1 || repo.tx.consultasAgendados != 2 {
		t.Fatalf("quantidade/consultas = %d/%d", quantidade, repo.tx.consultasAgendados)
	}
}

func TestIniciarPlantoesAgendadosRejeitaStatusInesperado(t *testing.T) {
	service, repo := novoServicoFake(StatusPlantaoAgendado, roleAdmin, "colaborador-1")
	repo.tx.plantoesAgendados = []*Plantao{plantaoDeTeste(StatusPlantaoCancelado)}

	_, err := service.IniciarPlantoesAgendados(context.Background())
	if !errors.Is(err, ErrorInvalidTransitionStatus) {
		t.Fatalf("erro = %v, esperado transição inválida", err)
	}
	if repo.committed || len(repo.tx.plantaoIniciado) != 0 || repo.tx.historicosAutomaticos != 0 {
		t.Fatal("status inesperado deixou alterações")
	}
}

func TestIniciarPlantoesAgendadosFazRollbackSeHistoricoFalhar(t *testing.T) {
	service, repo := novoServicoFake(StatusPlantaoAgendado, roleAdmin, "colaborador-1")
	sentinel := errors.New("falha no histórico")
	repo.tx.plantoesAgendados = []*Plantao{plantaoDeTeste(StatusPlantaoAgendado)}
	repo.tx.historicoAutomaticoErr = sentinel

	quantidade, err := service.IniciarPlantoesAgendados(context.Background())
	if !errors.Is(err, sentinel) || quantidade != 0 {
		t.Fatalf("quantidade/erro = %d/%v", quantidade, err)
	}
	if repo.committed {
		t.Fatal("transação com falha foi confirmada")
	}
}

func TestGetRelatorioValidaEPropagaFiltros(t *testing.T) {
	service, repo := novoServicoFake(StatusPlantaoAgendado, roleAdmin, "colaborador-1")
	inicio := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	fim := time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC)
	status := StatusPlantaoConcluido
	repo.relatorio = []RelatorioItem{{PlantaoID: "plantao-relatorio"}}

	resultado, err := service.GetRelatorio(context.Background(), &RelatorioFiltro{
		DataInicio:    inicio,
		DataFim:       fim,
		ColaboradorID: "colaborador-1",
		Status:        &status,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(resultado) != 1 || resultado[0].PlantaoID != "plantao-relatorio" {
		t.Fatalf("resultado inesperado: %+v", resultado)
	}
	if repo.relatorioFiltro == nil ||
		!repo.relatorioFiltro.DataInicio.Equal(inicio) ||
		!repo.relatorioFiltro.DataFim.Equal(fim) ||
		repo.relatorioFiltro.ColaboradorID != "colaborador-1" ||
		repo.relatorioFiltro.Status == nil || *repo.relatorioFiltro.Status != status {
		t.Fatalf("filtro não propagado: %+v", repo.relatorioFiltro)
	}
}

func TestGetRelatorioValidaPeriodoEStatus(t *testing.T) {
	inicio := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	fim := time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC)
	statusNegativo := StatusPlantao(-1)
	statusAcimaDoLimite := StatusPlantaoPago + 1

	tests := []struct {
		nome     string
		filtro   *RelatorioFiltro
		esperado error
	}{
		{"filtro ausente", nil, shared.ErrorPeriodoInvalido},
		{"início ausente", &RelatorioFiltro{DataFim: fim}, shared.ErrorPeriodoInvalido},
		{"fim ausente", &RelatorioFiltro{DataInicio: inicio}, shared.ErrorPeriodoInvalido},
		{"fim anterior", &RelatorioFiltro{DataInicio: fim, DataFim: inicio}, shared.ErrorEndBeforeStart},
		{"status negativo", &RelatorioFiltro{DataInicio: inicio, DataFim: fim, Status: &statusNegativo}, ErrorInvalidStatusPlantao},
		{"status acima do limite", &RelatorioFiltro{DataInicio: inicio, DataFim: fim, Status: &statusAcimaDoLimite}, ErrorInvalidStatusPlantao},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			service, repo := novoServicoFake(StatusPlantaoAgendado, roleAdmin, "colaborador-1")
			_, err := service.GetRelatorio(context.Background(), tt.filtro)
			if !errors.Is(err, tt.esperado) {
				t.Fatalf("erro = %v, esperado %v", err, tt.esperado)
			}
			if repo.relatorioFiltro != nil {
				t.Fatal("repositório foi consultado com filtro inválido")
			}
		})
	}
}

func TestGetRelatorioPropagaErroDoRepositorio(t *testing.T) {
	service, repo := novoServicoFake(StatusPlantaoAgendado, roleAdmin, "colaborador-1")
	sentinel := errors.New("falha no relatório")
	repo.relatorioErr = sentinel

	_, err := service.GetRelatorio(context.Background(), &RelatorioFiltro{
		DataInicio: time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC),
		DataFim:    time.Date(2026, 8, 31, 0, 0, 0, 0, time.UTC),
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("erro = %v, esperado sentinel", err)
	}
}
