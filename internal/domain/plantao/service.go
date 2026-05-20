package plantao

import (
	"context"

	"plantao/internal/domain/financeiro"
	"plantao/internal/domain/log"
	"plantao/internal/domain/shared"
)

type PlantaoService struct {
	repository     PlantaoRepository
	calculoService *financeiro.CalculoService
	log            log.Logger
}

func NewPlantaoService(repository PlantaoRepository, calculoService *financeiro.CalculoService, log log.Logger) *PlantaoService {
	return &PlantaoService{
		repository:     repository,
		calculoService: calculoService,
		log:            log,
	}
}

func (s *PlantaoService) CreatePlantao(ctx context.Context, colaboradorId string, periodo *shared.Periodo) (*Plantao, error) {
	s.log.Info("iniciando criação de plantão", "id_colaborador", colaboradorId, "periodo", periodo)

	existingPlantoes, err := s.repository.Find(ctx, &Filtro{
		ColaboradorID: colaboradorId,
		Periodo:       periodo,
	})

	if err != nil {
		s.log.Error("erro ao verificar existência de plantão no período", "id_colaborador", colaboradorId, "periodo", periodo, "error", err)
		return nil, err
	}

	if existingPlantoes != nil {
		s.log.Warn("plantão já existe para colaborador no período", "id_colaborador", colaboradorId, "periodo", periodo, "quantidade", len(existingPlantoes))
		return nil, ErrorExistingPlantao
	}

	s.log.Info("criando objeto de plantão", "id_colaborador", colaboradorId)
	plantao, err := NewPlantao(colaboradorId, periodo)
	if err != nil {
		s.log.Error("erro ao criar objeto de plantão", "id_colaborador", colaboradorId, "error", err)
		return nil, err
	}

	if err := s.repository.Store(ctx, plantao); err != nil {
		s.log.Error("erro ao salvar plantão no banco de dados", "id_plantao", plantao.Id, "id_colaborador", colaboradorId, "error", err)
		return nil, err
	}

	s.log.Info("plantão criado com sucesso", "id_plantao", plantao.Id, "id_colaborador", colaboradorId)
	return plantao, nil
}

func (s *PlantaoService) UpdatePlantaoStatus(ctx context.Context, plantaoId string, newStatus StatusPlantao, observacoes *string) (*Plantao, error) {
	s.log.Info("iniciando atualização de status do plantão", "id_plantao", plantaoId, "novo_status", newStatus)

	plantao, err := s.repository.FindById(ctx, plantaoId)
	if err != nil {
		s.log.Error("erro ao buscar plantão para atualização de status", "id_plantao", plantaoId, "error", err)
		return nil, err
	}

	if plantao == nil {
		s.log.Warn("plantão não encontrado para atualização de status", "id_plantao", plantaoId)
		return nil, ErrorPlantaoNotFinded
	}

	statusAnterior := plantao.Status
	if err := plantao.UpdateStatus(newStatus); err != nil {
		s.log.Warn("transição de status do plantão inválida", "id_plantao", plantaoId, "status_atual", statusAnterior, "novo_status", newStatus, "error", err)
		return nil, err
	}

	if newStatus == StatusPlantaoConcluido {
		s.log.Info("calculando valor do plantão concluído", "id_plantao", plantaoId)
		resultado, err := s.calculoService.Calcular(ctx, plantao.Periodo)
		if err != nil {
			s.log.Error("erro ao calcular valor do plantão concluído", "id_plantao", plantaoId, "error", err)
			return nil, err
		}

		detalhes := make([]PlantaoDetalhe, 0, len(resultado.Dias))
		for _, d := range resultado.Dias {
			detalhes = append(detalhes, PlantaoDetalhe{
				IdPlantao: plantaoId,
				Data:      d.Data,
				TipoDia:   string(d.TipoDia),
				Valor:     d.Valor,
			})
		}

		if err := s.repository.StoreDetalhesAndUpdateValorTotal(ctx, plantaoId, resultado.ValorTotal, observacoes, detalhes); err != nil {
			s.log.Error("erro ao salvar detalhes e valor total do plantão", "id_plantao", plantaoId, "valor_total", resultado.ValorTotal, "error", err)
			return nil, err
		}

		plantao.ValorTotal = resultado.ValorTotal
		plantao.Observacoes = observacoes
	} else {
		if err := s.repository.Update(ctx, plantao); err != nil {
			s.log.Error("erro ao atualizar plantão no banco de dados", "id_plantao", plantaoId, "novo_status", newStatus, "error", err)
			return nil, err
		}
	}

	s.log.Info("status do plantão atualizado com sucesso", "id_plantao", plantaoId, "status_anterior", statusAnterior, "novo_status", newStatus)
	return plantao, nil
}

func (s *PlantaoService) DeletePlantao(ctx context.Context, plantaoId string) error {
	s.log.Info("iniciando exclusão de plantão", "id_plantao", plantaoId)

	plantao, err := s.repository.FindById(ctx, plantaoId)
	if err != nil {
		s.log.Error("erro ao buscar plantão para exclusão", "id_plantao", plantaoId, "error", err)
		return err
	}

	if plantao == nil {
		s.log.Warn("plantão não encontrado para exclusão", "id_plantao", plantaoId)
		return ErrorPlantaoNotFinded
	}

	if err := s.repository.Delete(ctx, plantaoId); err != nil {
		s.log.Error("erro ao excluir plantão", "id_plantao", plantaoId, "error", err)
		return err
	}

	s.log.Info("plantão excluído com sucesso", "id_plantao", plantaoId)
	return nil
}

func (s *PlantaoService) GetPlantaoById(ctx context.Context, plantaoId string) (*Plantao, error) {
	s.log.Info("buscando plantão por ID", "id_plantao", plantaoId)

	plantao, err := s.repository.FindById(ctx, plantaoId)
	if err != nil {
		s.log.Error("erro ao buscar plantão por ID", "id_plantao", plantaoId, "error", err)
		return nil, err
	}

	if plantao == nil {
		s.log.Warn("plantão não encontrado", "id_plantao", plantaoId)
		return nil, ErrorPlantaoNotFinded
	}

	s.log.Info("plantão encontrado com sucesso", "id_plantao", plantaoId)
	return plantao, nil
}

func (s *PlantaoService) GetPlantoes(ctx context.Context, filter *Filtro) ([]Plantao, error) {
	s.log.Info("buscando plantões por filtro", "filter", filter)

	plantoes, err := s.repository.Find(ctx, filter)
	if err != nil {
		s.log.Error("erro ao buscar plantões por filtro", "filter", filter, "error", err)
		return nil, err
	}

	s.log.Info("plantões encontrados com sucesso", "quantidade", len(plantoes))
	return plantoes, nil
}

func (s *PlantaoService) GetPlantoesByColaboradorId(ctx context.Context, colaboradorId string) ([]Plantao, error) {
	s.log.Info("buscando plantões por colaborador", "id_colaborador", colaboradorId)

	plantoes, err := s.repository.Find(ctx, &Filtro{
		ColaboradorID: colaboradorId,
	})
	if err != nil {
		s.log.Error("erro ao buscar plantões por colaborador", "id_colaborador", colaboradorId, "error", err)
		return nil, err
	}

	s.log.Info("plantões por colaborador encontrados com sucesso", "id_colaborador", colaboradorId, "quantidade", len(plantoes))
	return plantoes, nil
}

func (s *PlantaoService) GetPlantoesByPeriodo(ctx context.Context, periodo *shared.Periodo) ([]Plantao, error) {
	s.log.Info("buscando plantões por período", "periodo", periodo)

	plantoes, err := s.repository.Find(ctx, &Filtro{
		Periodo: periodo,
	})
	if err != nil {
		s.log.Error("erro ao buscar plantões por período", "periodo", periodo, "error", err)
		return nil, err
	}

	s.log.Info("plantões por período encontrados com sucesso", "periodo", periodo, "quantidade", len(plantoes))
	return plantoes, nil
}

func (s *PlantaoService) GetPlantoesByStatus(ctx context.Context, status StatusPlantao) ([]Plantao, error) {
	s.log.Info("buscando plantões por status", "status", status)

	plantoes, err := s.repository.Find(ctx, &Filtro{
		Status: &status,
	})
	if err != nil {
		s.log.Error("erro ao buscar plantões por status", "status", status, "error", err)
		return nil, err
	}

	s.log.Info("plantões por status encontrados com sucesso", "status", status, "quantidade", len(plantoes))
	return plantoes, nil
}
