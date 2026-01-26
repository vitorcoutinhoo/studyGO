package plantao

import "context"

type PlantaoService struct {
	reposiotory PlantaoRepository
}

func NewPlantaoService(repository PlantaoRepository) *PlantaoService {
	return &PlantaoService{
		reposiotory: repository,
	}
}
func (s *PlantaoService) CreatePlantao(ctx context.Context, colaboradorId string, periodo *Periodo) (*Plantao, error) {
	existingPlantoes, _ := s.reposiotory.Find(ctx, &Filtro{
		ColaboradorID: colaboradorId,
		Periodo:       periodo,
	})

	if existingPlantoes != nil {
		return nil, ErrorExistingPlantao
	}

	plantao, err := NewPlantao(colaboradorId, periodo)
	if err != nil {
		return nil, err
	}

	if err := s.reposiotory.Store(ctx, plantao); err != nil {
		return nil, err
	}

	return plantao, nil
}

func (s *PlantaoService) UpdatePlantaoStatus(ctx context.Context, plantaoId string, newStatus StatusPlantao) (*Plantao, error) {
	plantao, _ := s.reposiotory.FindById(ctx, plantaoId)
	if plantao == nil {
		return nil, ErrorPlantaoNotFinded
	}

	if err := plantao.UpdateStatus(newStatus); err != nil {
		return nil, err
	}

	if err := s.reposiotory.Update(ctx, plantao); err != nil {
		return nil, err
	}

	return plantao, nil
}

func (s *PlantaoService) DeletePlantao(ctx context.Context, plantaoId string) error {
	plantao, _ := s.reposiotory.FindById(ctx, plantaoId)
	if plantao == nil {
		return ErrorPlantaoNotFinded
	}

	if err := s.reposiotory.Delete(ctx, plantaoId); err != nil {
		return err
	}

	return nil
}

func (s *PlantaoService) GetPlantaoById(ctx context.Context, plantaoId string) (*Plantao, error) {
	plantao, _ := s.reposiotory.FindById(ctx, plantaoId)
	if plantao == nil {
		return nil, ErrorPlantaoNotFinded
	}

	return plantao, nil
}

func (s *PlantaoService) GetPlantoes(ctx context.Context, filter *Filtro) ([]Plantao, error) {
	plantoes, err := s.reposiotory.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	return plantoes, nil
}

func (s *PlantaoService) GetPlantoesByColaboradorId(ctx context.Context, colaboradorId string) ([]Plantao, error) {
	plantoes, err := s.reposiotory.Find(ctx, &Filtro{
		ColaboradorID: colaboradorId,
	})

	if err != nil {
		return nil, err
	}

	return plantoes, nil
}

func (s *PlantaoService) GetPlantoesByPeriodo(ctx context.Context, periodo *Periodo) ([]Plantao, error) {
	plantoes, err := s.reposiotory.Find(ctx, &Filtro{
		Periodo: periodo,
	})

	if err != nil {
		return nil, err
	}

	return plantoes, nil
}

func (s *PlantaoService) GetPlantoesByStatus(ctx context.Context, status StatusPlantao) ([]Plantao, error) {
	plantoes, err := s.reposiotory.Find(ctx, &Filtro{
		Status: &status,
	})

	if err != nil {
		return nil, err
	}

	return plantoes, nil
}
