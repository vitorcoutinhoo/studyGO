package plantao

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"plantao/internal/domain/shared"
)

var (
	ErrorExistingPlantao  = errors.New("Plantao already exists!")
	ErrorPlantaoNotFinded = errors.New("Plantao not found!")
)

type Plantao struct {
	Id            string
	ColaboradorId string
	Periodo       *shared.Periodo
	Status        StatusPlantao
	ValorTotal    float64
	Observacoes   *string
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func NewPlantao(colaboradorId string, periodo *shared.Periodo) (*Plantao, error) {
	newPeriodo, err := shared.NewPeriodo(periodo.Inicio, periodo.Fim)

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
