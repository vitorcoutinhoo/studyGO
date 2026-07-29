package financeiro

import (
	"context"
	"errors"
	"testing"
	"time"

	"plantao/internal/domain/log"
	"plantao/internal/domain/shared"
)

type calculoFonteFake struct {
	feriados map[string]bool
	valores  map[TipoDia]int64
	err      error
}

func (f *calculoFonteFake) FindFeriados(context.Context, time.Time, time.Time) (map[string]bool, error) {
	return f.feriados, nil
}

func (f *calculoFonteFake) FindValorDiaCentavos(_ context.Context, tipo TipoDia, _ time.Time) (int64, error) {
	if f.err != nil {
		return 0, f.err
	}
	valor, ok := f.valores[tipo]
	if !ok {
		return 0, ErrorValorDiaNotFound
	}
	return valor, nil
}

type noopLogger struct{}

func (noopLogger) Info(string, ...any)  {}
func (noopLogger) Warn(string, ...any)  {}
func (noopLogger) Error(string, ...any) {}
func (noopLogger) Debug(string, ...any) {}
func (noopLogger) Fatal(string, ...any) {}
func (n noopLogger) With(...any) log.Logger {
	return n
}
func (noopLogger) Sync() error { return nil }

func TestCalcularClassificaDiasEPriorizaFeriado(t *testing.T) {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		t.Fatal(err)
	}
	inicio := time.Date(2026, 7, 3, 0, 0, 0, 0, loc) // sexta
	fim := time.Date(2026, 7, 5, 23, 0, 0, 0, loc)   // domingo
	fonte := &calculoFonteFake{
		feriados: map[string]bool{"2026-07-04": true},
		valores: map[TipoDia]int64{
			TipoDiaUtil: 10000, TipoDiaSabado: 20000,
			TipoDiaDomingo: 30000, TipoDiaFeriado: 40000,
		},
	}

	resultado, err := NewCalculoService(noopLogger{}).Calcular(
		context.Background(),
		&shared.Periodo{Inicio: inicio, Fim: fim},
		fonte,
		loc,
	)
	if err != nil {
		t.Fatal(err)
	}
	if resultado.ValorTotalCentavos != 80000 {
		t.Fatalf("total = %d, esperado 80000", resultado.ValorTotalCentavos)
	}
	tipos := []TipoDia{TipoDiaUtil, TipoDiaFeriado, TipoDiaDomingo}
	for i, esperado := range tipos {
		if resultado.Dias[i].TipoDia != esperado {
			t.Fatalf("dia %d = %s, esperado %s", i, resultado.Dias[i].TipoDia, esperado)
		}
	}
}

func TestCalcularUsaDataCivilDeSaoPaulo(t *testing.T) {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		t.Fatal(err)
	}
	// Em UTC já é sábado; em São Paulo ainda é sexta-feira.
	instante := time.Date(2026, 7, 4, 1, 30, 0, 0, time.UTC)
	fonte := &calculoFonteFake{
		valores: map[TipoDia]int64{TipoDiaUtil: 12345},
	}

	resultado, err := NewCalculoService(noopLogger{}).Calcular(
		context.Background(),
		&shared.Periodo{Inicio: instante, Fim: instante},
		fonte,
		loc,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(resultado.Dias) != 1 || resultado.Dias[0].Data.Format("2006-01-02") != "2026-07-03" {
		t.Fatalf("data civil inesperada: %+v", resultado.Dias)
	}
	if resultado.Dias[0].TipoDia != TipoDiaUtil {
		t.Fatalf("tipo = %s, esperado UTIL", resultado.Dias[0].TipoDia)
	}
}

func TestCalcularPropagaAusenciaDeConfiguracao(t *testing.T) {
	loc := time.UTC
	data := time.Date(2026, 7, 6, 0, 0, 0, 0, loc)
	fonte := &calculoFonteFake{err: ErrorValorDiaNotFound}

	_, err := NewCalculoService(noopLogger{}).Calcular(
		context.Background(),
		&shared.Periodo{Inicio: data, Fim: data},
		fonte,
		loc,
	)
	if !errors.Is(err, ErrorValorDiaNotFound) {
		t.Fatalf("erro = %v, esperado ErrorValorDiaNotFound", err)
	}
}
