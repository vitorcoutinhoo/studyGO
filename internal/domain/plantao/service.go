package plantao

import (
	"context"

	"plantao/internal/domain/financeiro"
	"plantao/internal/domain/shared"
)

type PlantaoService struct {
	repository    PlantaoRepository
	calculoService *financeiro.CalculoService
}

func NewPlantaoService(repository PlantaoRepository, calculoService *financeiro.CalculoService) *PlantaoService {
	return &PlantaoService{
		repository:    repository,
		calculoService: calculoService,
	}
}

func (s *PlantaoService) CreatePlantao(ctx context.Context, colaboradorId string, periodo *shared.Periodo) (*Plantao, error) {
	existingPlantoes, err := s.repository.Find(ctx, &Filtro{
		ColaboradorID: colaboradorId,
		Periodo:       periodo,
	})

	if err != nil {
		return nil, err
	}

	if existingPlantoes != nil {
		return nil, ErrorExistingPlantao
	}

	plantao, err := NewPlantao(colaboradorId, periodo)
	if err != nil {
		return nil, err
	}

	if err := s.repository.Store(ctx, plantao); err != nil {
		return nil, err
	}

	return plantao, nil
}

func (s *PlantaoService) UpdatePlantaoStatus(ctx context.Context, plantaoId string, newStatus StatusPlantao, observacoes *string) (*Plantao, error) {
	plantao, err := s.repository.FindById(ctx, plantaoId)
	if err != nil {
		return nil, err
	}

	if plantao == nil {
		return nil, ErrorPlantaoNotFinded
	}

	if err := plantao.UpdateStatus(newStatus); err != nil {
		return nil, err
	}

	if newStatus == StatusPlantaoConcluido {
		resultado, err := s.calculoService.Calcular(ctx, plantao.Periodo)
		if err != nil {
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
			return nil, err
		}

		plantao.ValorTotal = resultado.ValorTotal
		plantao.Observacoes = observacoes
	} else {
		if err := s.repository.Update(ctx, plantao); err != nil {
			return nil, err
		}
	}

	return plantao, nil
}

func (s *PlantaoService) DeletePlantao(ctx context.Context, plantaoId string) error {
	plantao, err := s.repository.FindById(ctx, plantaoId)
	if err != nil {
		return err
	}

	if plantao == nil {
		return ErrorPlantaoNotFinded
	}

	return s.repository.Delete(ctx, plantaoId)
}

func (s *PlantaoService) GetPlantaoById(ctx context.Context, plantaoId string) (*Plantao, error) {
	plantao, err := s.repository.FindById(ctx, plantaoId)
	if err != nil {
		return nil, err
	}

	if plantao == nil {
		return nil, ErrorPlantaoNotFinded
	}

	return plantao, nil
}

func (s *PlantaoService) GetPlantoes(ctx context.Context, filter *Filtro) ([]Plantao, error) {
	return s.repository.Find(ctx, filter)
}

func (s *PlantaoService) GetPlantoesByColaboradorId(ctx context.Context, colaboradorId string) ([]Plantao, error) {
	return s.repository.Find(ctx, &Filtro{
		ColaboradorID: colaboradorId,
	})
}

func (s *PlantaoService) GetPlantoesByPeriodo(ctx context.Context, periodo *shared.Periodo) ([]Plantao, error) {
	return s.repository.Find(ctx, &Filtro{
		Periodo: periodo,
	})
}

func (s *PlantaoService) GetPlantoesByStatus(ctx context.Context, status StatusPlantao) ([]Plantao, error) {
	return s.repository.Find(ctx, &Filtro{
		Status: &status,
	})
}
