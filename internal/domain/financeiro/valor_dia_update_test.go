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
	valores   []ValorDia
	stored    *ValorDia
	storeErr  error
	updateErr error
	committed bool
}

func (r *valorDiaRepositoryFake) FindAll(context.Context) ([]ValorDia, error) {
	return r.valores, nil
}
func (r *valorDiaRepositoryFake) Store(_ context.Context, valor *ValorDia) error {
	if r.storeErr != nil {
		return r.storeErr
	}
	copia := *valor
	r.stored = &copia
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
	copia.UpdatedAt = time.Date(2026, 8, 5, 15, 0, 0, 0, time.UTC)
	t.atualizado = &copia
	return &copia, nil
}

func configGlobalTeste() *ValorDia {
	descricao := "descrição original"
	fim := time.Date(2026, 6, 30, 0, 0, 0, 0, time.UTC)
	return &ValorDia{
		Id:             uuid.New(),
		TipoDia:        TipoDiaUtil,
		Valor:          150.75,
		Descricao:      &descricao,
		VigenciaInicio: time.Date(2026, 5, 1, 0, 0, 0, 0, time.UTC),
		VigenciaFim:    &fim,
	}
}

func serviceValorDiaTeste(repo ValorDiaRepository) *ConfigValorDiaService {
	return NewConfigValorDiaService(repo, noopLogger{})
}

func TestSetValorCriaConfiguracaoGlobal(t *testing.T) {
	repo := &valorDiaRepositoryFake{}
	descricao := "valor de feriado"
	resultado, err := serviceValorDiaTeste(repo).SetValor(context.Background(), TipoDiaFeriado, 300, &descricao)
	if err != nil {
		t.Fatal(err)
	}
	if repo.stored == nil || resultado.Id != repo.stored.Id ||
		resultado.TipoDia != TipoDiaFeriado || resultado.Valor != 300 ||
		resultado.Descricao == nil || *resultado.Descricao != descricao ||
		!resultado.VigenciaInicio.Equal(configValorDiaInicioGlobal) || resultado.VigenciaFim != nil {
		t.Fatalf("configuração global incorreta: %+v", resultado)
	}
}

func TestSetValorValidaEntradaEPropagaDuplicidade(t *testing.T) {
	zero := 0.0
	negativo := -1.0
	precisaoInvalida := 10.001
	limiteExcedido := 100_000_000.0
	tests := []struct {
		nome     string
		tipo     TipoDia
		valor    float64
		storeErr error
		esperado error
	}{
		{"tipo inválido", TipoDia("INVALIDO"), 10, nil, ErrorTipoDiaInvalido},
		{"valor zero", TipoDiaUtil, zero, nil, ErrorValorDiaInvalido},
		{"valor negativo", TipoDiaUtil, negativo, nil, ErrorValorDiaInvalido},
		{"precisão inválida", TipoDiaUtil, precisaoInvalida, nil, ErrorPrecisaoValorDia},
		{"limite excedido", TipoDiaUtil, limiteExcedido, nil, ErrorLimiteValorDia},
		{"tipo duplicado", TipoDiaUtil, 10, ErrorValorDiaAlreadyExists, ErrorValorDiaAlreadyExists},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			repo := &valorDiaRepositoryFake{storeErr: tt.storeErr}
			_, err := serviceValorDiaTeste(repo).SetValor(context.Background(), tt.tipo, tt.valor, nil)
			if !errors.Is(err, tt.esperado) {
				t.Fatalf("erro = %v, esperado %v", err, tt.esperado)
			}
		})
	}
}

func TestGetAllRetornaConfiguracoesIndependentementeDasDatasInternas(t *testing.T) {
	config := *configGlobalTeste()
	repo := &valorDiaRepositoryFake{valores: []ValorDia{config}}
	resultado, err := serviceValorDiaTeste(repo).GetAll(context.Background())
	if err != nil || len(resultado) != 1 || resultado[0].Id != config.Id {
		t.Fatalf("resultado/erro = %+v/%v", resultado, err)
	}
}

func TestUpdateAtualizaSomenteValorEDescricao(t *testing.T) {
	original := configGlobalTeste()
	inicioOriginal := original.VigenciaInicio
	fimOriginal := *original.VigenciaFim
	repo := &valorDiaRepositoryFake{atual: original}
	novoValor := 225.50
	novaDescricao := "configuração global"

	resultado, err := serviceValorDiaTeste(repo).Update(context.Background(), TipoDiaUtil, &AtualizacaoValorDia{
		Valor:              &novoValor,
		DescricaoInformada: true,
		Descricao:          &novaDescricao,
	})
	if err != nil {
		t.Fatal(err)
	}
	if !repo.committed || resultado.Valor != novoValor || resultado.Descricao == nil || *resultado.Descricao != novaDescricao {
		t.Fatalf("configuração não atualizada: %+v", resultado)
	}
	if !resultado.VigenciaInicio.Equal(inicioOriginal) || resultado.VigenciaFim == nil || !resultado.VigenciaFim.Equal(fimOriginal) {
		t.Fatalf("datas internas foram modificadas: %+v", resultado)
	}
}

func TestUpdatePermiteLimparDescricao(t *testing.T) {
	repo := &valorDiaRepositoryFake{atual: configGlobalTeste()}
	resultado, err := serviceValorDiaTeste(repo).Update(context.Background(), TipoDiaUtil, &AtualizacaoValorDia{
		DescricaoInformada: true,
		Descricao:          nil,
	})
	if err != nil {
		t.Fatal(err)
	}
	if resultado.Descricao != nil {
		t.Fatalf("descrição não foi removida: %+v", resultado)
	}
}

func TestUpdateValidaEntrada(t *testing.T) {
	zero := 0.0
	precisaoInvalida := 10.001
	tests := []struct {
		nome        string
		tipo        TipoDia
		atual       *ValorDia
		atualizacao *AtualizacaoValorDia
		esperado    error
	}{
		{"tipo inválido", TipoDia("INVALIDO"), configGlobalTeste(), &AtualizacaoValorDia{Valor: &zero}, ErrorTipoDiaInvalido},
		{"patch vazio", TipoDiaUtil, configGlobalTeste(), &AtualizacaoValorDia{}, ErrorAtualizacaoVazia},
		{"valor zero", TipoDiaUtil, configGlobalTeste(), &AtualizacaoValorDia{Valor: &zero}, ErrorValorDiaInvalido},
		{"precisão inválida", TipoDiaUtil, configGlobalTeste(), &AtualizacaoValorDia{Valor: &precisaoInvalida}, ErrorPrecisaoValorDia},
		{"configuração ausente", TipoDiaUtil, nil, &AtualizacaoValorDia{DescricaoInformada: true}, ErrorValorDiaNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.nome, func(t *testing.T) {
			repo := &valorDiaRepositoryFake{atual: tt.atual}
			_, err := serviceValorDiaTeste(repo).Update(context.Background(), tt.tipo, tt.atualizacao)
			if !errors.Is(err, tt.esperado) {
				t.Fatalf("erro = %v, esperado %v", err, tt.esperado)
			}
			if repo.committed {
				t.Fatal("operação inválida confirmou a transação")
			}
		})
	}
}

func TestUpdateFazRollbackQuandoPersistenciaFalha(t *testing.T) {
	sentinel := errors.New("falha simulada")
	original := configGlobalTeste()
	repo := &valorDiaRepositoryFake{atual: original, updateErr: sentinel}
	novoValor := 200.0

	_, err := serviceValorDiaTeste(repo).Update(context.Background(), TipoDiaUtil, &AtualizacaoValorDia{Valor: &novoValor})
	if !errors.Is(err, sentinel) {
		t.Fatalf("erro = %v, esperado sentinel", err)
	}
	if repo.committed || repo.atual.Valor != original.Valor {
		t.Fatal("falha de persistência deixou alteração parcial")
	}
}
