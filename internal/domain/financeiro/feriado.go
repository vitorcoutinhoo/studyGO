package financeiro

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorFeriadoNotFound       = errors.New("Feriado não encontrado!")
	ErrorFeriadoDataInvalid    = errors.New("Data do feriado inválida!")
	ErrorFeriadoNotMunicipal   = errors.New("Apenas feriados municipais podem ter a data alterada!")
)

type Feriado struct {
	Id        uuid.UUID
	Data      time.Time
	Nome      string
	Descricao string
	CreatedAt time.Time
}

type FeriadoRepository interface {
	FindById(ctx context.Context, id uuid.UUID) (*Feriado, error)
	FindByAno(ctx context.Context, ano int) ([]Feriado, error)
	FindByPeriodo(ctx context.Context, inicio, fim time.Time) (map[time.Time]bool, error)
	UpdateData(ctx context.Context, id uuid.UUID, novaData time.Time) error
}

type FeriadoService struct {
	repository FeriadoRepository
}

func NewFeriadoService(repository FeriadoRepository) *FeriadoService {
	return &FeriadoService{repository: repository}
}

func (s *FeriadoService) GetFeriadosByAno(ctx context.Context, ano int) ([]Feriado, error) {
	return s.repository.FindByAno(ctx, ano)
}

func (s *FeriadoService) UpdateDataFeriado(ctx context.Context, id uuid.UUID, novaData time.Time) (*Feriado, error) {
	if novaData.IsZero() {
		return nil, ErrorFeriadoDataInvalid
	}

	feriado, err := s.repository.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	if feriado == nil {
		return nil, ErrorFeriadoNotFound
	}
	if feriado.Descricao != "MUNICIPAL" {
		return nil, ErrorFeriadoNotMunicipal
	}

	novaData = normalizeDate(novaData)

	if err := s.repository.UpdateData(ctx, id, novaData); err != nil {
		return nil, err
	}

	feriado.Data = novaData
	return feriado, nil
}
