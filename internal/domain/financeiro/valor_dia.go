package financeiro

import (
	"context"
	"errors"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"plantao/internal/domain/log"
	"plantao/internal/domain/shared"
)

var (
	ErrorTipoDiaInvalido      = errors.New("Tipo de dia inválido!")
	ErrorValorDiaInvalido     = errors.New("Valor deve ser maior que zero!")
	ErrorValorDiaNotFound     = errors.New("Configuração de valor não encontrada!")
	ErrorValorDiaForaVigencia = errors.New("Configuração de valor fora da vigência!")
	ErrorVigenciaInvalida     = errors.New("Período de vigência inválido!")
	ErrorPrecisaoValorDia     = errors.New("Valor deve possuir no máximo duas casas decimais!")
	ErrorLimiteValorDia       = errors.New("Valor excede o limite monetário suportado!")
	ErrorAtualizacaoVazia     = errors.New("Informe ao menos um campo para atualização!")
	ErrorConflitoValorDia     = errors.New("Conflito de concorrência ao atualizar configuração de valor!")
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
	Descricao      *string
	VigenciaInicio time.Time
	VigenciaFim    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type AtualizacaoValorDia struct {
	Valor                *float64
	DescricaoInformada   bool
	Descricao            *string
	VigenciaInicio       *time.Time
	VigenciaFimInformada bool
	VigenciaFim          *time.Time
}

type DiaCalculado struct {
	Data          time.Time
	TipoDia       TipoDia
	ValorCentavos int64
}

type ResultadoCalculo struct {
	ValorTotalCentavos int64
	Dias               []DiaCalculado
}

type CalculoFonte interface {
	FindFeriados(ctx context.Context, inicio, fim time.Time) (map[string]bool, error)
	FindValorDiaCentavos(ctx context.Context, tipoDia TipoDia, data time.Time) (int64, error)
}

type ValorDiaRepository interface {
	FindVigentes(ctx context.Context) ([]ValorDia, error)
	FindVigenteByTipoDia(ctx context.Context, tipoDia TipoDia) (*ValorDia, error)
	FindVigenteByData(ctx context.Context, data time.Time) (map[TipoDia]float64, error)
	Store(ctx context.Context, valorDia *ValorDia) error
	CloseVigencia(ctx context.Context, id uuid.UUID, vigenciaFim time.Time) error
	WithTransaction(ctx context.Context, fn func(ValorDiaTransaction) error) error
}

type ValorDiaTransaction interface {
	LockByTipoDia(ctx context.Context, tipoDia TipoDia) (*ValorDia, error)
	Update(ctx context.Context, valorDia *ValorDia) (*ValorDia, error)
}

type ConfigValorDiaService struct {
	repository ValorDiaRepository
	log        log.Logger
}

func NewConfigValorDiaService(repository ValorDiaRepository, log log.Logger) *ConfigValorDiaService {
	return &ConfigValorDiaService{
		repository: repository,
		log:        log,
	}
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

func (s *ConfigValorDiaService) UpdateVigente(ctx context.Context, tipoDia TipoDia, atualizacao *AtualizacaoValorDia) (*ValorDia, error) {
	s.log.Info("iniciando atualização de configuração de valor", "tipo_dia", tipoDia)

	if !tiposDiaValidos[tipoDia] {
		return nil, ErrorTipoDiaInvalido
	}
	if atualizacao == nil || !atualizacao.temCampos() {
		return nil, ErrorAtualizacaoVazia
	}
	if atualizacao.Valor != nil {
		if err := validarValorMonetario(*atualizacao.Valor); err != nil {
			return nil, err
		}
	}

	var resultado *ValorDia

	err := s.repository.WithTransaction(ctx, func(tx ValorDiaTransaction) error {
		configuracao, err := tx.LockByTipoDia(ctx, tipoDia)
		if err != nil {
			return err
		}

		if atualizacao.Valor != nil {
			configuracao.Valor = *atualizacao.Valor
		}
		if atualizacao.DescricaoInformada {
			configuracao.Descricao = atualizacao.Descricao
		}
		if atualizacao.VigenciaInicio != nil {
			data := normalizeDate(*atualizacao.VigenciaInicio)
			configuracao.VigenciaInicio = data
		}
		if atualizacao.VigenciaFimInformada {
			if atualizacao.VigenciaFim == nil {
				configuracao.VigenciaFim = nil
			} else {
				data := normalizeDate(*atualizacao.VigenciaFim)
				configuracao.VigenciaFim = &data
			}
		}

		if configuracao.VigenciaFim != nil && configuracao.VigenciaFim.Before(configuracao.VigenciaInicio) {
			return ErrorVigenciaInvalida
		}

		resultado, err = tx.Update(ctx, configuracao)
		return err
	})
	if err != nil {
		s.log.Warn("falha ao atualizar configuração de valor", "tipo_dia", tipoDia, "error", err)
		return nil, err
	}

	s.log.Info("configuração de valor atualizada", "tipo_dia", tipoDia, "id_valor_dia", resultado.Id)
	return resultado, nil
}

func (a *AtualizacaoValorDia) temCampos() bool {
	return a.Valor != nil ||
		a.DescricaoInformada ||
		a.VigenciaInicio != nil ||
		a.VigenciaFimInformada
}

func validarValorMonetario(valor float64) error {
	if valor <= 0 || math.IsNaN(valor) || math.IsInf(valor, 0) {
		return ErrorValorDiaInvalido
	}
	if valor > 99_999_999.99 {
		return ErrorLimiteValorDia
	}
	if math.Abs(valor*100-math.Round(valor*100)) > 0.0000001 {
		return ErrorPrecisaoValorDia
	}
	return nil
}

// ---- CalculoService ----

type CalculoService struct {
	log log.Logger
}

func NewCalculoService(log log.Logger) *CalculoService {
	return &CalculoService{log: log}
}

func (s *CalculoService) Calcular(ctx context.Context, periodo *shared.Periodo, fonte CalculoFonte, loc *time.Location) (*ResultadoCalculo, error) {
	s.log.Info("iniciando cálculo financeiro do período", "inicio", periodo.Inicio, "fim", periodo.Fim)

	if periodo == nil || fonte == nil {
		return nil, shared.ErrorPeriodoInvalido
	}
	if loc == nil {
		loc = time.UTC
	}

	inicio := normalizeDateInLocation(periodo.Inicio, loc)
	fim := normalizeDateInLocation(periodo.Fim, loc)
	if fim.Before(inicio) {
		return nil, shared.ErrorEndBeforeStart
	}

	feriados, err := fonte.FindFeriados(ctx, inicio, fim)
	if err != nil {
		s.log.Error("erro ao buscar feriados para cálculo financeiro", "inicio", inicio, "fim", fim, "error", err)
		return nil, fmt.Errorf("erro ao buscar feriados: %w", err)
	}

	var resultado ResultadoCalculo
	for data := inicio; !data.After(fim); data = data.AddDate(0, 0, 1) {
		dia := Dia{
			Data:      data,
			DiaSemana: DiaDaSemana(data.Weekday()),
			EhFeriado: feriados[data.Format("2006-01-02")],
		}
		tipoDia := determinaTipoDia(dia)
		valorCentavos, err := fonte.FindValorDiaCentavos(ctx, tipoDia, data)
		if err != nil {
			s.log.Warn("valor não configurado para data do cálculo financeiro", "tipo_dia", tipoDia, "data", data, "error", err)
			return nil, err
		}
		if valorCentavos <= 0 {
			return nil, ErrorValorDiaInvalido
		}
		if resultado.ValorTotalCentavos > math.MaxInt64-valorCentavos {
			return nil, ErrorValorDiaInvalido
		}

		resultado.Dias = append(resultado.Dias, DiaCalculado{
			Data:          dia.Data,
			TipoDia:       tipoDia,
			ValorCentavos: valorCentavos,
		})
		resultado.ValorTotalCentavos += valorCentavos
	}

	s.log.Info("cálculo financeiro concluído com sucesso", "inicio", inicio, "fim", fim, "quantidade_dias", len(resultado.Dias), "valor_total_centavos", resultado.ValorTotalCentavos)
	return &resultado, nil
}

func normalizeDateInLocation(t time.Time, loc *time.Location) time.Time {
	local := t.In(loc)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, loc)
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
