package financeiro

import (
	"context"
	"errors"
	"plantao/internal/domain/log"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorFeriadoNotFound     = errors.New("Feriado não encontrado!")
	ErrorFeriadoDataInvalid  = errors.New("Data do feriado inválida!")
	ErrorFeriadoNotMunicipal = errors.New("Apenas feriados municipais podem ter a data alterada!")
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
	log        log.Logger
}

func NewFeriadoService(repository FeriadoRepository, log log.Logger) *FeriadoService {
	return &FeriadoService{repository: repository, log: log}
}

func (s *FeriadoService) GetFeriadosByAno(ctx context.Context, ano int) ([]Feriado, error) {
	s.log.Info("buscando feriados por ano", "ano", ano)

	feriados, err := s.repository.FindByAno(ctx, ano)
	if err != nil {
		s.log.Error("erro ao buscar feriados por ano", "ano", ano, "error", err)
		return nil, err
	}

	s.log.Info("feriados encontrados com sucesso", "ano", ano, "quantidade", len(feriados))
	return feriados, nil
}

func (s *FeriadoService) UpdateDataFeriado(ctx context.Context, id uuid.UUID, novaData time.Time) (*Feriado, error) {
	s.log.Info("iniciando atualização da data do feriado", "id_feriado", id, "nova_data", novaData)

	if novaData.IsZero() {
		s.log.Warn("data do feriado inválida", "id_feriado", id)
		return nil, ErrorFeriadoDataInvalid
	}

	feriado, err := s.repository.FindById(ctx, id)
	if err != nil {
		s.log.Error("erro ao buscar feriado para atualização", "id_feriado", id, "error", err)
		return nil, err
	}
	if feriado == nil {
		s.log.Warn("feriado não encontrado para atualização", "id_feriado", id)
		return nil, ErrorFeriadoNotFound
	}
	if feriado.Descricao != "MUNICIPAL" {
		s.log.Warn("tentativa de atualizar feriado não municipal", "id_feriado", id, "descricao", feriado.Descricao)
		return nil, ErrorFeriadoNotMunicipal
	}

	novaData = normalizeDate(novaData)

	if err := s.repository.UpdateData(ctx, id, novaData); err != nil {
		s.log.Error("erro ao atualizar data do feriado", "id_feriado", id, "nova_data", novaData, "error", err)
		return nil, err
	}

	feriado.Data = novaData
	s.log.Info("data do feriado atualizada com sucesso", "id_feriado", id, "nova_data", novaData)
	return feriado, nil
}
