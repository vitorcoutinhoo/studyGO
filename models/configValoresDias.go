package models

import (
	"time"

	"github.com/google/uuid"
)

type ConfigValoresDias struct {
	ID             uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	TipoDia        string    `gorm:"type:varchar(100);not null;unique" json:"tipo_dia"`
	Valor          float64   `gorm:"type:decimal(10,2);not null" json:"valor"`
	Descricao      string    `gorm:"type:text" json:"descricao"`
	VigenciaInicio time.Time `gorm:"type:date;not null" json:"vigencia_inicio"`
	VigenciaFim    time.Time `gorm:"type:date" json:"vigencia_fim"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}
