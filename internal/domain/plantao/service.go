package plantao

import (
	"context"
	"time"

	"plantao/internal/domain/financeiro"
	"plantao/internal/domain/log"
	"plantao/internal/domain/shared"
)

const (
	roleAdmin       = "admin"
	roleGerente     = "gerente"
	roleColaborador = "colaborador"

	statusPagamentoPendente    = "pendente"
	maxPagamentoCentavos       = int64(9_999_999_999)
	loteInicioAutomatico       = 100
	observacaoInicioAutomatico = "Status alterado automaticamente pelo sistema"
)

type PlantaoService struct {
	repository     PlantaoRepository
	calculoService *financeiro.CalculoService
	log            log.Logger
	location       *time.Location
	now            func() time.Time
}

func (s *PlantaoService) IniciarPlantoesAgendados(ctx context.Context) (int, error) {
	instante := s.now()
	totalAtualizado := 0

	for {
		atualizadosNoLote := 0
		err := s.repository.WithTransaction(ctx, func(tx PlantaoTransaction) error {
			plantoes, err := tx.LockPlantoesAgendadosAte(ctx, instante, loteInicioAutomatico)
			if err != nil {
				return err
			}

			for _, p := range plantoes {
				statusAnterior := p.Status
				if err := p.UpdateStatus(StatusPlantaoEmAndamento); err != nil {
					return err
				}
				if err := tx.StartPlantao(ctx, p.Id); err != nil {
					return err
				}
				if err := tx.InsertHistoricoAutomatico(
					ctx,
					p.Id,
					statusAnterior,
					StatusPlantaoEmAndamento,
					observacaoInicioAutomatico,
				); err != nil {
					return err
				}
				atualizadosNoLote++
			}
			return nil
		})
		if err != nil {
			s.log.Error("erro ao iniciar plantões automaticamente", "quantidade_processada", totalAtualizado, "error", err)
			return totalAtualizado, err
		}

		totalAtualizado += atualizadosNoLote
		if atualizadosNoLote < loteInicioAutomatico {
			if totalAtualizado > 0 {
				s.log.Info("plantões iniciados automaticamente", "quantidade", totalAtualizado, "instante_limite", instante)
			}
			return totalAtualizado, nil
		}
	}
}

func NewPlantaoService(repository PlantaoRepository, calculoService *financeiro.CalculoService, log log.Logger) *PlantaoService {
	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		location = time.UTC
	}
	return &PlantaoService{
		repository:     repository,
		calculoService: calculoService,
		log:            log,
		location:       location,
		now:            time.Now,
	}
}

func (s *PlantaoService) CreatePlantao(ctx context.Context, colaboradorID string, periodo *shared.Periodo) (*Plantao, error) {
	s.log.Info("iniciando criação de plantão", "id_colaborador", colaboradorID, "periodo", periodo)
	existing, err := s.repository.Find(ctx, &Filtro{ColaboradorID: colaboradorID, Periodo: periodo})
	if err != nil {
		return nil, err
	}
	if len(existing) > 0 {
		return nil, ErrorExistingPlantao
	}

	p, err := NewPlantao(colaboradorID, periodo)
	if err != nil {
		return nil, err
	}
	if err := s.repository.Store(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *PlantaoService) UpdatePlantaoStatus(ctx context.Context, plantaoID, usuarioID string, newStatus StatusPlantao, observacoes *string) (*Plantao, error) {
	if newStatus == StatusPlantaoConcluido {
		return s.FecharPlantao(ctx, plantaoID, usuarioID, observacoes)
	}
	if newStatus == StatusPlantaoPago {
		return nil, ErrorOperacaoPagamentoObrigatoria
	}

	var result *Plantao
	err := s.repository.WithTransaction(ctx, func(tx PlantaoTransaction) error {
		p, err := tx.LockPlantao(ctx, plantaoID)
		if err != nil {
			return err
		}
		if _, err := tx.FindActor(ctx, usuarioID); err != nil {
			return err
		}

		oldStatus := p.Status
		if err := p.UpdateStatus(newStatus); err != nil {
			return err
		}
		if err := tx.UpdateStatus(ctx, plantaoID, newStatus); err != nil {
			return err
		}
		if err := tx.InsertHistorico(ctx, plantaoID, oldStatus, newStatus, usuarioID, observacoes); err != nil {
			return err
		}
		result = p
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PlantaoService) FecharPlantao(ctx context.Context, plantaoID, usuarioID string, observacoes *string) (*Plantao, error) {
	s.log.Info("iniciando fechamento de plantão", "id_plantao", plantaoID, "id_usuario", usuarioID)

	var result *Plantao
	err := s.repository.WithTransaction(ctx, func(tx PlantaoTransaction) error {
		p, err := tx.LockPlantao(ctx, plantaoID)
		if err != nil {
			return err
		}
		actor, err := tx.FindActor(ctx, usuarioID)
		if err != nil {
			return err
		}
		if err := s.validateActorAndColaborador(ctx, tx, actor, p); err != nil {
			return err
		}
		if actor.Role == roleColaborador && actor.ColaboradorID != p.ColaboradorId {
			return ErrorUsuarioSemPermissao
		}
		if actor.Role != roleColaborador && actor.Role != roleGerente && actor.Role != roleAdmin {
			return ErrorUsuarioSemPermissao
		}

		switch p.Status {
		case StatusPlantaoConcluido:
			return ErrorPlantaoJaFechado
		case StatusPlantaoPago:
			return ErrorPlantaoJaPago
		case StatusPlantaoEmAndamento:
			// estado válido
		default:
			return ErrorInvalidTransitionStatus
		}

		detalhes, err := tx.CountDetalhes(ctx, plantaoID)
		if err != nil {
			return err
		}
		if detalhes != 0 {
			return ErrorDetalhesExistentes
		}
		pagamentos, err := tx.CountPagamentos(ctx, plantaoID)
		if err != nil {
			return err
		}
		if pagamentos != 0 {
			return ErrorPagamentoExistente
		}

		calculo, err := s.calculoService.Calcular(ctx, p.Periodo, tx, s.location)
		if err != nil {
			return err
		}
		if calculo.ValorTotalCentavos <= 0 || calculo.ValorTotalCentavos > maxPagamentoCentavos {
			return ErrorValorTotalInvalido
		}

		dias := make([]PlantaoDetalhe, 0, len(calculo.Dias))
		for _, dia := range calculo.Dias {
			dias = append(dias, PlantaoDetalhe{
				Data:          dia.Data,
				TipoDia:       string(dia.TipoDia),
				ValorCentavos: dia.ValorCentavos,
			})
		}
		if err := tx.InsertDetalhes(ctx, plantaoID, dias); err != nil {
			return err
		}
		if err := tx.ClosePlantao(ctx, plantaoID, calculo.ValorTotalCentavos, observacoes); err != nil {
			return err
		}
		if err := tx.InsertHistorico(ctx, plantaoID, StatusPlantaoEmAndamento, StatusPlantaoConcluido, usuarioID, observacoes); err != nil {
			return err
		}
		if err := tx.InsertPagamentoPendente(ctx, plantaoID, p.ColaboradorId, calculo.ValorTotalCentavos); err != nil {
			return err
		}

		p.Status = StatusPlantaoConcluido
		p.ValorTotal = float64(calculo.ValorTotalCentavos) / 100
		p.Observacoes = observacoes
		result = p
		return nil
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (s *PlantaoService) PagarPlantao(ctx context.Context, plantaoID, usuarioID string, observacoes *string) error {
	s.log.Info("iniciando pagamento de plantão", "id_plantao", plantaoID, "id_usuario", usuarioID)

	return s.repository.WithTransaction(ctx, func(tx PlantaoTransaction) error {
		p, err := tx.LockPlantao(ctx, plantaoID)
		if err != nil {
			return err
		}
		actor, err := tx.FindActor(ctx, usuarioID)
		if err != nil {
			return err
		}
		if err := s.validateActorAndColaborador(ctx, tx, actor, p); err != nil {
			return err
		}
		if actor.Role != roleGerente && actor.Role != roleAdmin {
			return ErrorUsuarioSemPermissao
		}

		switch p.Status {
		case StatusPlantaoPago:
			return ErrorPlantaoJaPago
		case StatusPlantaoConcluido:
			// estado válido
		default:
			return ErrorInvalidTransitionStatus
		}
		if p.ValorTotal <= 0 {
			return ErrorValorTotalInvalido
		}

		pagamento, err := tx.FindPagamento(ctx, plantaoID)
		if err != nil {
			return err
		}
		if pagamento.Status != statusPagamentoPendente ||
			pagamento.ColaboradorID != p.ColaboradorId {
			return ErrorPagamentoInconsistente
		}

		resumo, err := tx.SummarizeDetalhes(ctx, plantaoID)
		if err != nil {
			return err
		}
		diasEsperados, err := quantidadeDiasCivis(p.Periodo, s.location)
		if err != nil {
			return err
		}
		valorPlantaoCentavos := int64(p.ValorTotal*100 + 0.5)
		if resumo.Quantidade != diasEsperados ||
			resumo.DatasDistintas != diasEsperados ||
			resumo.ValorTotalCentavos != valorPlantaoCentavos {
			return ErrorDetalhesInconsistentes
		}
		if pagamento.ValorTotalCentavos != valorPlantaoCentavos {
			return ErrorPagamentoInconsistente
		}

		agora := s.now().In(s.location)
		dataPagamento := time.Date(agora.Year(), agora.Month(), agora.Day(), 0, 0, 0, 0, s.location)
		if err := tx.PayPagamento(ctx, pagamento.ID, dataPagamento, observacoes); err != nil {
			return err
		}
		if err := tx.UpdateStatus(ctx, plantaoID, StatusPlantaoPago); err != nil {
			return err
		}
		return tx.InsertHistorico(ctx, plantaoID, StatusPlantaoConcluido, StatusPlantaoPago, usuarioID, observacoes)
	})
}

func (s *PlantaoService) validateActorAndColaborador(ctx context.Context, tx PlantaoTransaction, actor *Actor, p *Plantao) error {
	actorExists, err := tx.ColaboradorExists(ctx, actor.ColaboradorID)
	if err != nil {
		return err
	}
	plantaoColaboradorExists, err := tx.ColaboradorExists(ctx, p.ColaboradorId)
	if err != nil {
		return err
	}
	if !actorExists || !plantaoColaboradorExists {
		return ErrorColaboradorNotFound
	}
	return nil
}

func quantidadeDiasCivis(periodo *shared.Periodo, loc *time.Location) (int, error) {
	if periodo == nil {
		return 0, shared.ErrorPeriodoInvalido
	}
	inicioLocal := periodo.Inicio.In(loc)
	fimLocal := periodo.Fim.In(loc)
	inicio := time.Date(inicioLocal.Year(), inicioLocal.Month(), inicioLocal.Day(), 0, 0, 0, 0, loc)
	fim := time.Date(fimLocal.Year(), fimLocal.Month(), fimLocal.Day(), 0, 0, 0, 0, loc)
	if fim.Before(inicio) {
		return 0, shared.ErrorEndBeforeStart
	}
	total := 0
	for data := inicio; !data.After(fim); data = data.AddDate(0, 0, 1) {
		total++
	}
	return total, nil
}

func (s *PlantaoService) DeletePlantao(ctx context.Context, plantaoID string) error {
	p, err := s.repository.FindById(ctx, plantaoID)
	if err != nil {
		return err
	}
	if p == nil {
		return ErrorPlantaoNotFinded
	}
	return s.repository.Delete(ctx, plantaoID)
}

func (s *PlantaoService) GetPlantaoById(ctx context.Context, plantaoID string) (*Plantao, error) {
	p, err := s.repository.FindById(ctx, plantaoID)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, ErrorPlantaoNotFinded
	}
	return p, nil
}

func (s *PlantaoService) GetPlantoes(ctx context.Context, filter *Filtro) ([]Plantao, error) {
	return s.repository.Find(ctx, filter)
}

func (s *PlantaoService) GetPlantoesByColaboradorId(ctx context.Context, colaboradorID string) ([]Plantao, error) {
	return s.repository.Find(ctx, &Filtro{ColaboradorID: colaboradorID})
}

func (s *PlantaoService) GetPlantoesByPeriodo(ctx context.Context, periodo *shared.Periodo) ([]Plantao, error) {
	return s.repository.Find(ctx, &Filtro{Periodo: periodo})
}

func (s *PlantaoService) GetPlantoesByStatus(ctx context.Context, status StatusPlantao) ([]Plantao, error) {
	return s.repository.Find(ctx, &Filtro{Status: &status})
}
