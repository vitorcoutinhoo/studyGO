package financeiro

import (
	"time"

	"plantao/internal/domain/shared"
)

type DiaDaSemana int

const (
	Domingo DiaDaSemana = iota
	SegundaFeira
	TercaFeira
	QuartaFeira
	QuintaFeira
	SextaFeira
	Sabado
)

type Dia struct {
	Data      time.Time
	DiaSemana DiaDaSemana
	EhFeriado bool
}

func Dias(p *shared.Periodo, feriados map[time.Time]bool) []Dia {
	var dias []Dia

	data := p.Inicio
	for !data.After(p.Fim) {
		dia := Dia{
			Data:      data,
			DiaSemana: DiaDaSemana(data.Weekday()),
			EhFeriado: feriados[normalizeDate(data)],
		}

		dias = append(dias, dia)
		data = data.AddDate(0, 0, 1)
	}

	return dias
}

func normalizeDate(t time.Time) time.Time {
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}