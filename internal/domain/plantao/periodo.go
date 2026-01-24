package plantao

import (
	"errors"
	"time"
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

var (
	ErrorEndBeforeStart  = errors.New("Data de término anterior à data de início!")
	ErrorPeriodoInvalido = errors.New("Período Inválido!")
)

type Periodo struct {
	Inicio time.Time
	Fim    time.Time
}

type Dia struct {
	Data      time.Time
	DiaSemana DiaDaSemana
	ehFeriado bool
}

func NewPeriodo(inicio, fim time.Time) (*Periodo, error) {
	if inicio.IsZero() || fim.IsZero() {
		return nil, ErrorPeriodoInvalido
	}

	if fim.Before(inicio) {
		return nil, ErrorEndBeforeStart
	}

	return &Periodo{
		Inicio: inicio,
		Fim:    fim,
	}, nil
}

func (p *Periodo) Dias(feriado map[time.Time]bool) []Dia {
	var dias []Dia

	data := p.Inicio
	for !data.After(p.Fim) {
		dia := Dia{
			Data:      data,
			DiaSemana: DiaDaSemana(data.Weekday()),
			ehFeriado: feriado[data],
		}

		dias = append(dias, dia)
		data = data.AddDate(0, 0, 1)
	}

	return dias
}
