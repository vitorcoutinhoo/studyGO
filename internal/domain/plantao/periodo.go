package plantao

import (
	"errors"
	"time"
)

var (
	ErrorEndBeforeStart  = errors.New("Data de término anterior à data de início!")
	ErrorPeriodoInvalido = errors.New("Período Inválido!")
)

type Periodo struct {
	Inicio time.Time
	Fim    time.Time
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
