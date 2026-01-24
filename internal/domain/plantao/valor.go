package plantao

import "errors"

var (
	ErrorNegtiveMoneyValue = errors.New("Valor em dinheiro não pode ser negativo!")
)

type Valor struct {
	Quantidade int64
}

type DiaValor struct {
	Dia   Dia
	Valor Valor
}

type Dias []Dia

func newValor(quantidade float64) (*Valor, error) {
	if quantidade < 0 {
		return nil, ErrorNegtiveMoneyValue
	}

	return &Valor{
		Quantidade: int64(quantidade * 100),
	}, nil
}

func (d Dias) ValorDiaSemana() []DiaValor {
	var diasValor []DiaValor

	for _, dia := range d {
		var valorDia float64

		switch dia.DiaSemana {
		case Domingo:
			valorDia = 150.00
		case Sabado:
			valorDia = 150.00
		default:
			valorDia = 120.00
		}

		if dia.ehFeriado {
			valorDia *= 2.0
		}

		valor, _ := newValor(valorDia)
		diasValor = append(diasValor, DiaValor{
			Dia:   dia,
			Valor: *valor,
		})
	}

	return diasValor
}
