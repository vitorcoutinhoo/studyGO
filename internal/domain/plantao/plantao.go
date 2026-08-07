package plantao

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"plantao/internal/domain/shared"
)

var (
	ErrorExistingPlantao              = errors.New("Plantao already exists!")
	ErrorPlantaoNotFinded             = errors.New("Plantao not found!")
	ErrorColaboradorNotFound          = errors.New("Colaborador do plantão não encontrado!")
	ErrorUsuarioSemPermissao          = errors.New("Usuário sem permissão para esta operação!")
	ErrorPlantaoJaFechado             = errors.New("Plantão já fechado!")
	ErrorPlantaoJaPago                = errors.New("Plantão já pago!")
	ErrorDetalhesExistentes           = errors.New("Plantão já possui detalhes calculados!")
	ErrorPagamentoExistente           = errors.New("Plantão já possui pagamento!")
	ErrorPagamentoNotFound            = errors.New("Pagamento do plantão não encontrado!")
	ErrorPagamentoInconsistente       = errors.New("Pagamento inconsistente com o plantão!")
	ErrorDetalhesInconsistentes       = errors.New("Detalhes calculados inconsistentes com o plantão!")
	ErrorConflitoConcorrencia         = errors.New("Conflito de concorrência ao processar o plantão!")
	ErrorValorTotalInvalido           = errors.New("Valor total do plantão inválido!")
	ErrorOperacaoPagamentoObrigatoria = errors.New("Status pago exige a operação de pagamento!")
	ErrorAtualizacaoPlantaoVazia      = errors.New("Informe ao menos um campo para atualização do plantão!")
	ErrorPlantaoNaoEditavel           = errors.New("Somente plantões agendados podem ser editados!")
)

type AtualizacaoPlantao struct {
	ColaboradorID *string
	DataInicio    *time.Time
	DataFim       *time.Time
}

func (a *AtualizacaoPlantao) TemCampos() bool {
	return a != nil && (a.ColaboradorID != nil || a.DataInicio != nil || a.DataFim != nil)
}

type Plantao struct {
	Id            string
	ColaboradorId string
	Periodo       *shared.Periodo
	Status        StatusPlantao
	ValorTotal    float64
	Observacoes   *string
	shared.Auditoria
}

func NewPlantao(colaboradorId string, periodo *shared.Periodo) (*Plantao, error) {
	newPeriodo, err := shared.NewPeriodo(periodo.Inicio, periodo.Fim)

	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &Plantao{
		Id:            uuid.NewString(),
		ColaboradorId: colaboradorId,
		Periodo:       newPeriodo,
		Status:        StatusPlantaoAgendado,
		Auditoria:     shared.Auditoria{CreatedAt: now, UpdatedAt: now},
	}, nil
}

func (p *Plantao) UpdateStatus(newStatus StatusPlantao) error {
	if !p.Status.canStatusPlantaoTransitionTo(newStatus) {
		return ErrorInvalidTransitionStatus
	}

	p.Status = newStatus
	p.UpdatedAt = time.Now()
	return nil
}
