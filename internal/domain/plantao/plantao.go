package plantao

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrorExistingPlantao  = errors.New("Plantao already exists!")
	ErrorPlantaoNotFinded = errors.New("Plantao not found!")
)

type Plantao struct {
	Id            string
	ColaboradorId string
	Periodo       *Periodo
	Status        StatusPlantao
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewPlantao(colaboradorId string, periodo *Periodo) (*Plantao, error) {
	newPeriodo, err := NewPeriodo(periodo.Inicio, periodo.Fim)

	if err != nil {
		return nil, err
	}

	return &Plantao{
		Id:            uuid.NewString(),
		ColaboradorId: colaboradorId,
		Periodo:       newPeriodo,
		Status:        StatusPlantaoAgendado,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}, nil
}

func (p *Plantao) UpdateStatus(newStatus StatusPlantao) error {
	if !p.Status.canStatusPlantaoTransitionTo(newStatus) {
		return ErrorInvalidTransitionStatus
	}

	p.Status = newStatus
	p.UpdatedAt = time.Now()
	return nil
}
