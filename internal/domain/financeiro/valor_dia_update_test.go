package financeiro

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
)

type valorDiaRepositoryFake struct {
	atual     *ValorDia
	updateErr error
	committed bool
}

func (r *valorDiaRepositoryFake) FindVigentes(context.Context) ([]ValorDia, error) {
	return nil, nil
}
func (r *valorDiaRepositoryFake) FindVigenteByTipoDia(context.Context, TipoDia) (*ValorDia, error) {
	return nil, nil
}
func (r *valorDiaRepositoryFake) FindVigenteByData(context.Context, time.Time) (map[TipoDia]float64, error) {
	return nil, nil
}
func (r *valorDiaRepositoryFake) Store(context.Context, *ValorDia) error {
	return nil
}
func (r *valorDiaRepositoryFake) CloseVigencia(context.Context, uuid.UUID, time.Time) error {
	return nil
}
func (r *valorDiaRepositoryFake) WithTransaction(_ context.Context, fn func(ValorDiaTransaction) error) error {
	tx := &valorDiaTransactionFake{repository: r}
	if err := fn(tx); err != nil {
		return err
	}
	r.committed = true
	r.atual = tx.atualizado
	return nil
}

type valorDiaTransactionFake struct {
	repository *valorDiaRepositoryFake
	atualizado *ValorDia
}

func (t *valorDiaTransactionFake) LockByTipoDia(context.Context, TipoDia) (*ValorDia, error) {
	if t.repository.atual == nil {
		return nil, ErrorValorDiaNotFound
	}
	copia := *t.repository.atual
	if copia.Descricao != nil {
		descricao := *copia.Descricao
		copia.Descricao = &descricao
	}
	if copia.VigenciaFim != nil {
		fim := *copia.VigenciaFim
		copia.VigenciaFim = &fim
	}
	return &copia, nil
}

func (t *valorDiaTransactionFake) Update(_ context.Context, valor *ValorDia) (*ValorDia, error) {
	if t.repository.updateErr != nil {
		return nil, t.repository.updateErr
	}
	copia := *valor
	copia.UpdatedAt = time.Date(2026, 7, 30, 15, 0, 0, 0, time.UTC)
	t.atualizado = &copia
	return &copia, nil
}

func configVigenteTeste() *ValorDia {
	descricao := "descrição original"
	return &ValorDia{
		Id:             uuid.New(),
		TipoDia:        TipoDiaUtil,
		Valor:          150.75,
		Descricao:      &descricao,
		VigenciaInicio: time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
	}
}

func serviceValorDiaTeste(repo ValorDiaRepository) *ConfigValorDiaService {
	service := NewConfigValorDiaService(repo, noopLogger{})
	service.now = func() time.Time {
		return time.Date(2026, 7, 30, 12, 0, 0, 0, service.location)
	}
	return service
}

func TestUpdateVigenteAtualizaParcialmente(t *testing.T) {
	repo := &valorDiaRepositoryFake{atual: configVigenteTeste()}
	service := serviceValorDiaTeste(repo)
	novoValor := 175.50

	resultado, err := service.UpdateVigente(context.Background(), TipoDiaUtil, &AtualizacaoValorDia{Valor: &novoValor})
	if err != nil {
		t.Fatal(err)
	}
	if resultado.Valor != novoValor || resultado.Descricao == nil || *resultado.Descricao != "descrição original" {
		t.Fatalf("atualização parcial incorreta: %+v", resultado)
	}
	if !repo.committed || resultado.UpdatedAt.IsZero() {
		t.Fatal("transação não foi confirmada ou updated_at não foi atualizado")
	}
}

func TestUpdateVigentePermiteLimparCamposAnulaveis(t *testing.T) {
	fim := time.Date(2026, 8, 30, 0, 0, 0, 0, time.UTC)
	config := configVigenteTeste()
	config.VigenciaFim = &fim
	repo := &valorDiaRepositoryFake{atual: config}
	service := serviceValorDiaTeste(repo)

	resultado, err := service.UpdateVigente(context.Background(), TipoDiaUtil, &AtualizacaoValorDia{
		DescricaoInformada:   true,
		Descricao:            nil,
		VigenciaFimInformada: true,
		VigenciaFim:          nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resultado.Descricao != nil || resultado.VigenciaFim != nil {
		t.Fatalf("campos não foram limpos: %+v", resultado)
	}
}

func TestUpdateVigenteAtualizaDescricaoEVigencias(t *testing.T) {
	repo := &valorDiaRepositoryFake{atual: configVigenteTeste()}
	service := serviceValorDiaTeste(repo)
	descricao := "novo valor para dia útil"
	inicio := time.Date(2026, 6, 1, 15, 0, 0, 0, time.FixedZone("offset", -3*60*60))
	fim := time.Date(2026, 12, 31, 15, 0, 0, 0, time.FixedZone("offset", -3*60*60))

	resultado, err := service.UpdateVigente(context.Background(), TipoDiaUtil, &AtualizacaoValorDia{
		DescricaoInformada:   true,
		Descricao:            &descricao,
		VigenciaInicio:       &inicio,
		VigenciaFimInformada: true,
		VigenciaFim:          &fim,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resultado.Descricao == nil || *resultado.Descricao != descricao ||
		resultado.VigenciaInicio.Format("2006-01-02") != "2026-06-01" ||
		resultado.VigenciaFim == nil || resultado.VigenciaFim.Format("2006-01-02") != "2026-12-31" {
		t.Fatalf("descrição/vigências incorretas: %+v", resultado)
	}
}

func TestUpdateVigenteValidaEntradaEVigencia(t *testing.T) {
	zero := 0.0
	negativo := -1.0
	precisaoInvalida := 10.001
	inicioFuturo := time.Date(2026, 8, 1, 0, 0, 0, 0, time.UTC)
	fimAnterior := time.Date(2026, 4, 1, 0, 0, 0, 0, time.UTC)
	configEncerrada := configVigenteTeste()
	configFim := time.Date(2026, 7, 29, 0, 0, 0, 0, time.UTC)
	configEncerrada.VigenciaFim = &configFim

	tests := []struct {
		nome        string
		tipo        TipoDia
		atual       *ValorDia
		atualizacao *AtualizacaoValorDia
		esperado    error
	}{
		{"tipo inválido", TipoDia("INVALIDO"), configVigenteTeste(), &AtualizacaoValorDia{Valor: &zero}, ErrorTipoDiaInvalido},
		{"patch vazio", TipoDiaUtil, configVigenteTeste(), &AtualizacaoValorDia{}, ErrorAtualizacaoVazia},
		{"valor zero", TipoDiaUtil, configVigenteTeste(), &AtualizacaoValorDia{Valor: &zero}, ErrorValorDiaInvalido},
		{"valor negativo", TipoDiaUtil, configVigenteTeste(), &AtualizacaoValorDia{Valor: &negativo}, ErrorValorDiaInvalido},
		{"precisão inválida", TipoDiaUtil, configVigenteTeste(), &AtualizacaoValorDia{Valor: &precisaoInvalida}, ErrorPrecisaoValorDia},
		{"configuração ausente", TipoDiaUtil, nil, &AtualizacaoValorDia{DescricaoInformada: true}, ErrorValorDiaNotFound},
		{"configuração não vigente", TipoDiaUtil, configEncerrada, &AtualizacaoValorDia{DescricaoInformada: true}, ErrorValorDiaNaoVigente},
		{"início futuro", TipoDiaUtil, configVigenteTeste(), &AtualizacaoValorDia{VigenciaInicio: &inicioFuturo}, ErrorValorDiaNaoVigente},
		{"fim anterior ao início", TipoDiaUtil, configVigenteTeste(), &AtualizacaoValorDia{VigenciaFimInformada: true, VigenciaFim: &fimAnterior}, ErrorVigenciaInvalida},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			repo := &valorDiaRepositoryFake{atual: tt.atual}
			_, err := serviceValorDiaTeste(repo).UpdateVigente(context.Background(), tt.tipo, tt.atualizacao)
			if !errors.Is(err, tt.esperado) {
				t.Fatalf("erro = %v, esperado %v", err, tt.esperado)
			}
			if repo.committed {
				t.Fatal("operação inválida confirmou a transação")
			}
		})
	}
}

func TestUpdateVigenteFazRollbackQuandoPersistenciaFalha(t *testing.T) {
	sentinel := errors.New("falha simulada")
	original := configVigenteTeste()
	repo := &valorDiaRepositoryFake{atual: original, updateErr: sentinel}
	novoValor := 200.0

	_, err := serviceValorDiaTeste(repo).UpdateVigente(
		context.Background(),
		TipoDiaUtil,
		&AtualizacaoValorDia{Valor: &novoValor},
	)
	if !errors.Is(err, sentinel) {
		t.Fatalf("erro = %v, esperado sentinel", err)
	}
	if repo.committed || repo.atual.Valor != original.Valor {
		t.Fatal("falha de persistência deixou alteração parcial")
	}
}
