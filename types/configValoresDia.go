package types

import (
	"time"

	"github.com/google/uuid"
)

type ConfigValoresDiaRequest struct {
	TipoDia        string     `json:"tipo_dia"`
	Valor          float32    `json:"valor"`
	Descricao      string     `json:"descricao"`
	VigenciaInicio time.Time  `json:"vigencia_inicio"`
	VigenciaFim    *time.Time `json:"vigencia_fim"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type ConfigValoresDiaResponse struct {
	ID             uuid.UUID  `json:"id"`
	TipoDia        string     `json:"tipo_dia"`
	Valor          float32    `json:"valor"`
	Descricao      string     `json:"descricao"`
	VigenciaInicio time.Time  `json:"vigencia_inicio"`
	VigenciaFim    *time.Time `json:"vigencia_fim"`
	CreatedAt      *time.Time `json:"created_at"`
	UpdatedAt      *time.Time `json:"updated_at"`
}

type ConfigValoresDiaInterface interface {
	CreateConfigValoresDia(valoresDia ConfigValoresDiaRequest) (*ConfigValoresDiaResponse, error)
	GetConfigValoresDia() ([]*ConfigValoresDiaRequest, error)
	GetConfigValoresDiaById(id uuid.UUID) (*ConfigValoresDiaRequest, error)
	UpdateConfigValoresDia(id uuid.UUID, valoresDia ConfigValoresDiaRequest) error
	DeleteConfigValoresDiaById(id uuid.UUID) error
}
