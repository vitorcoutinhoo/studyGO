package models

import (
	"time"

	"github.com/google/uuid"
)

type Plantoes struct {
	ID         uuid.UUID `gorm:"type:uuid;primaryKey" json:"id"`
	DataInicio time.Time `gorm:"type:date;not null" json:"data_inicio"`
	DataFim    time.Time `gorm:"type:date;not null" json:"data_fim"`
	Status     string    `gorm:"type:varchar(50);default: 'Agendado'" json:"status"`
	Valor      float64   `gorm:"type:decimal(10,2);not null" json:"valor"`
}
