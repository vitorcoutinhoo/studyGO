package financeiro

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"plantao/internal/domain/log"
	"plantao/internal/domain/shared"
)

var (
	ErrorTipoDiaInvalido  = errors.New("Tipo de dia inválido!")
	ErrorValorDiaInvalido = errors.New("Valor deve ser maior que zero!")
	ErrorValorDiaNotFound = errors.New("Configuração de valor não encontrada!")
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
	log        log.Logger
}

func NewConfigValorDiaService(repository ValorDiaRepository, log log.Logger) *ConfigValorDiaService {
	return &ConfigValorDiaService{repository: repository, log: log}
}

func (s *ConfigValorDiaService) GetVigentes(ctx context.Context) ([]ValorDia, error) {
	s.log.Info("buscando configurações de valor vigentes")

	valores, err := s.repository.FindVigentes(ctx)
	if err != nil {
		s.log.Error("erro ao buscar configurações de valor vigentes", "error", err)
		return nil, err
	}

	s.log.Info("configurações de valor vigentes encontradas com sucesso", "quantidade", len(valores))
	return valores, nil
}

func (s *ConfigValorDiaService) SetValor(ctx context.Context, tipoDia TipoDia, valor float64, vigenciaInicio time.Time) (*ValorDia, error) {
	s.log.Info("iniciando configuração de valor por tipo de dia", "tipo_dia", tipoDia, "valor", valor, "vigencia_inicio", vigenciaInicio)

	if !tiposDiaValidos[tipoDia] {
		s.log.Warn("tipo de dia inválido para configuração de valor", "tipo_dia", tipoDia)
		return nil, ErrorTipoDiaInvalido
	}
	if valor <= 0 {
		s.log.Warn("valor inválido para configuração de tipo de dia", "tipo_dia", tipoDia, "valor", valor)
		return nil, ErrorValorDiaInvalido
	}

	vigenciaInicio = normalizeDate(vigenciaInicio)

	// fecha o vigente anterior, se existir
	s.log.Info("buscando configuração vigente anterior", "tipo_dia", tipoDia)
	vigente, err := s.repository.FindVigenteByTipoDia(ctx, tipoDia)
	if err != nil {
		s.log.Error("erro ao buscar configuração vigente anterior", "tipo_dia", tipoDia, "error", err)
		return nil, err
	}
	if vigente != nil {
		ontem := vigenciaInicio.AddDate(0, 0, -1)
		s.log.Info("fechando vigência anterior", "id_valor_dia", vigente.Id, "tipo_dia", tipoDia, "vigencia_fim", ontem)
		if err := s.repository.CloseVigencia(ctx, vigente.Id, ontem); err != nil {
			s.log.Error("erro ao fechar vigência anterior", "id_valor_dia", vigente.Id, "tipo_dia", tipoDia, "error", err)
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
		s.log.Error("erro ao salvar nova configuração de valor", "tipo_dia", tipoDia, "valor", valor, "error", err)
		return nil, err
	}

	s.log.Info("configuração de valor criada com sucesso", "id_valor_dia", novo.Id, "tipo_dia", novo.TipoDia)
	return novo, nil
}

// ---- CalculoService ----

type CalculoService struct {
	feriadoRepo  FeriadoRepository
	valorDiaRepo ValorDiaRepository
	log          log.Logger
}

func NewCalculoService(feriadoRepo FeriadoRepository, valorDiaRepo ValorDiaRepository, log log.Logger) *CalculoService {
	return &CalculoService{
		feriadoRepo:  feriadoRepo,
		valorDiaRepo: valorDiaRepo,
		log:          log,
	}
}

func (s *CalculoService) Calcular(ctx context.Context, periodo *shared.Periodo) (*ResultadoCalculo, error) {
	s.log.Info("iniciando cálculo financeiro do período", "inicio", periodo.Inicio, "fim", periodo.Fim)

	feriados, err := s.feriadoRepo.FindByPeriodo(ctx, periodo.Inicio, periodo.Fim)
	if err != nil {
		s.log.Error("erro ao buscar feriados para cálculo financeiro", "inicio", periodo.Inicio, "fim", periodo.Fim, "error", err)
		return nil, fmt.Errorf("erro ao buscar feriados: %w", err)
	}

	valores, err := s.valorDiaRepo.FindVigenteByData(ctx, periodo.Inicio)
	if err != nil {
		s.log.Error("erro ao buscar valores vigentes para cálculo financeiro", "data_base", periodo.Inicio, "error", err)
		return nil, fmt.Errorf("erro ao buscar valores: %w", err)
	}

	dias := Dias(periodo, feriados)

	var resultado ResultadoCalculo
	for _, dia := range dias {
		tipoDia := determinaTipoDia(dia)

		valor, ok := valores[tipoDia]
		if !ok {
			s.log.Warn("valor não configurado para tipo de dia no cálculo financeiro", "tipo_dia", tipoDia, "data", dia.Data)
			return nil, fmt.Errorf("valor não configurado para o tipo de dia: %s", tipoDia)
		}

		resultado.Dias = append(resultado.Dias, DiaCalculado{
			Data:    dia.Data,
			TipoDia: tipoDia,
			Valor:   valor,
		})
		resultado.ValorTotal += valor
	}

	s.log.Info("cálculo financeiro concluído com sucesso", "inicio", periodo.Inicio, "fim", periodo.Fim, "quantidade_dias", len(resultado.Dias), "valor_total", resultado.ValorTotal)
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
