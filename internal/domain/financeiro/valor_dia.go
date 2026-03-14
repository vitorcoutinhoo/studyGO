package financeiro

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"plantao/internal/domain/shared"
)

var (
	ErrorTipoDiaInvalido    = errors.New("Tipo de dia inválido!")
	ErrorValorDiaInvalido   = errors.New("Valor deve ser maior que zero!")
	ErrorValorDiaNotFound   = errors.New("Configuração de valor não encontrada!")
)

type TipoDia string

const (
	TipoDiaUtil    TipoDia = "UTIL"
	TipoDiaSabado  TipoDia = "SABADO"
	TipoDiaDomingo TipoDia = "DOMINGO"
	TipoDiaFeriado TipoDia = "FERIADO"
)

var tiposDiaValidos = map[TipoDia]bool{
	TipoDiaUtil: true, TipoDiaSabado: true,
	TipoDiaDomingo: true, TipoDiaFeriado: true,
}

type ValorDia struct {
	Id             uuid.UUID
	TipoDia        TipoDia
	Valor          float64
	VigenciaInicio time.Time
	VigenciaFim    *time.Time
}

type DiaCalculado struct {
	Data    time.Time
	TipoDia TipoDia
	Valor   float64
}

type ResultadoCalculo struct {
	ValorTotal float64
	Dias       []DiaCalculado
}

type ValorDiaRepository interface {
	FindVigentes(ctx context.Context) ([]ValorDia, error)
	FindVigenteByTipoDia(ctx context.Context, tipoDia TipoDia) (*ValorDia, error)
	FindVigenteByData(ctx context.Context, data time.Time) (map[TipoDia]float64, error)
	Store(ctx context.Context, valorDia *ValorDia) error
	CloseVigencia(ctx context.Context, id uuid.UUID, vigenciaFim time.Time) error
}

type ConfigValorDiaService struct {
	repository ValorDiaRepository
}

func NewConfigValorDiaService(repository ValorDiaRepository) *ConfigValorDiaService {
	return &ConfigValorDiaService{repository: repository}
}

func (s *ConfigValorDiaService) GetVigentes(ctx context.Context) ([]ValorDia, error) {
	return s.repository.FindVigentes(ctx)
}

func (s *ConfigValorDiaService) SetValor(ctx context.Context, tipoDia TipoDia, valor float64, vigenciaInicio time.Time) (*ValorDia, error) {
	if !tiposDiaValidos[tipoDia] {
		return nil, ErrorTipoDiaInvalido
	}
	if valor <= 0 {
		return nil, ErrorValorDiaInvalido
	}

	vigenciaInicio = normalizeDate(vigenciaInicio)

	// fecha o vigente anterior, se existir
	vigente, err := s.repository.FindVigenteByTipoDia(ctx, tipoDia)
	if err != nil {
		return nil, err
	}
	if vigente != nil {
		ontem := vigenciaInicio.AddDate(0, 0, -1)
		if err := s.repository.CloseVigencia(ctx, vigente.Id, ontem); err != nil {
			return nil, err
		}
	}

	novo := &ValorDia{
		Id:             uuid.New(),
		TipoDia:        tipoDia,
		Valor:          valor,
		VigenciaInicio: vigenciaInicio,
	}

	if err := s.repository.Store(ctx, novo); err != nil {
		return nil, err
	}

	return novo, nil
}

// ---- CalculoService ----

type CalculoService struct {
	feriadoRepo  FeriadoRepository
	valorDiaRepo ValorDiaRepository
}

func NewCalculoService(feriadoRepo FeriadoRepository, valorDiaRepo ValorDiaRepository) *CalculoService {
	return &CalculoService{
		feriadoRepo:  feriadoRepo,
		valorDiaRepo: valorDiaRepo,
	}
}

func (s *CalculoService) Calcular(ctx context.Context, periodo *shared.Periodo) (*ResultadoCalculo, error) {
	feriados, err := s.feriadoRepo.FindByPeriodo(ctx, periodo.Inicio, periodo.Fim)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar feriados: %w", err)
	}

	valores, err := s.valorDiaRepo.FindVigenteByData(ctx, periodo.Inicio)
	if err != nil {
		return nil, fmt.Errorf("erro ao buscar valores: %w", err)
	}

	dias := Dias(periodo, feriados)

	var resultado ResultadoCalculo
	for _, dia := range dias {
		tipoDia := determinaTipoDia(dia)

		valor, ok := valores[tipoDia]
		if !ok {
			return nil, fmt.Errorf("valor não configurado para o tipo de dia: %s", tipoDia)
		}

		resultado.Dias = append(resultado.Dias, DiaCalculado{
			Data:    dia.Data,
			TipoDia: tipoDia,
			Valor:   valor,
		})
		resultado.ValorTotal += valor
	}

	return &resultado, nil
}

func determinaTipoDia(dia Dia) TipoDia {
	if dia.EhFeriado {
		return TipoDiaFeriado
	}
	switch dia.DiaSemana {
	case Sabado:
		return TipoDiaSabado
	case Domingo:
		return TipoDiaDomingo
	default:
		return TipoDiaUtil
	}
}
