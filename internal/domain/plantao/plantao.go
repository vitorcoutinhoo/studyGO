package plantao

import (
	"github.com/google/uuid"
)

type Plantao struct {
	Id            string
	ColaboradorId string
	Periodo       *Periodo
	Status        StatusPlantao
}

func NewPlantao(colaboradorId string, periodo Periodo) (*Plantao, error) {
	newPeriodo, err := NewPeriodo(periodo.Inicio, periodo.Fim)

	if err != nil {
		return nil, err
	}

	return &Plantao{
		Id:            uuid.NewString(),
		ColaboradorId: colaboradorId,
		Periodo:       newPeriodo,
		Status:        StatusPlantaoAgendado,
	}, nil
}

func ValorTotalPlantao(dias []Dia) (*Valor, error) {
	diasValor := Dias(dias).ValorDiaSemana()

	var total int64
	for _, diaValor := range diasValor {
		total += diaValor.Valor.Quantidade
	}

	return &Valor{
		Quantidade: total,
	}, nil
}

func (p *Plantao) UpdateStatus(newStatus StatusPlantao) {
	if p.Status.canStatusPlantaoTransitionTo(newStatus) {
		p.Status = newStatus
	}
}
