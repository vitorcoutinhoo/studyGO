package plantao

import "errors"

type StatusPlantao int

const (
	StatusPlantaoAgendado StatusPlantao = iota
	StatusPlantaoEmAndamento
	StatusPlantaoConcluido
	StatusPlantaoCancelado
	StatusPlantaoPago
)

var (
	ErrorInvalidStatusPlantao    = errors.New("Status do Plantão Inválido!")
	ErrorInvalidTransitionStatus = errors.New("Transição de Status Inválida!")
)

func (s StatusPlantao) canStatusPlantaoTransitionTo(new StatusPlantao) bool {
	switch s {
	case StatusPlantaoAgendado:
		return new == StatusPlantaoEmAndamento || new == StatusPlantaoCancelado
	case StatusPlantaoEmAndamento:
		return new == StatusPlantaoConcluido || new == StatusPlantaoCancelado
	case StatusPlantaoConcluido:
		return new == StatusPlantaoPago
	case StatusPlantaoCancelado:
		return false
	case StatusPlantaoPago:
		return false
	default:
		return false
	}
}
